package transport

// Лимиты JSON API: 45 единиц за 3 секунды на пользователя или решение, не более 5 параллельных
// запросов от него и 20 от аккаунта; вес запроса зависит от способа авторизации и даты.
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_restrictions.md

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Profile — кто авторизует запросы. По Bearer-строке его не определить, поэтому он задаётся явно;
// нулевое значение NewLimiter не принимает.
type Profile int

const (
	_ Profile = iota
	// ProfileUnknown считается как самый тяжёлый вес обычного запроса — 4 единицы.
	ProfileUnknown
	ProfileApp      // токен публичного или приватного решения
	ProfileDraftApp // токен решения в статусе «Черновик» или «На модерации»
	ProfileEmployee // логин и пароль или токен пользователя
)

const (
	DefaultBudget               = 45
	DefaultWindow               = 3 * time.Second
	DefaultPrincipalConcurrency = 5
	DefaultAccountConcurrency   = 20
)

// Clock подменяется в тестах; nil — системное время.
type Clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}

// Account ограничивает параллельные запросы аккаунта и делится между Limiter разных
// пользователей и решений одного процесса. Другие процессы и токены он не видит:
// общий лимит аккаунта между процессами требует внешней координации.
type Account struct {
	slots chan struct{}
}

func NewAccount(concurrency int) (*Account, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("moysklad: account concurrency must be positive")
	}
	return &Account{slots: make(chan struct{}, concurrency)}, nil
}

type LimiterConfig struct {
	Profile Profile
	// Account — общий семафор аккаунта; nil — собственный на DefaultAccountConcurrency.
	Account *Account
	Clock   Clock
	// Weight переопределяет вес запроса; nil — DefaultWeight.
	Weight func(profile Profile, method, path string, now time.Time) int
	// Нулевые значения — Default*.
	Budget      int
	Window      time.Duration
	Concurrency int
}

// Limiter — бюджет одного пользователя или решения. Один экземпляр делится между всеми Client
// этого principal в процессе. X-RateLimit-* ответа он не использует: их единицы не уточнены.
type Limiter struct {
	cfg   LimiterConfig
	slots chan struct{}
	turn  chan struct{} // очередь ожидания бюджета: ожидающие на канале обслуживаются по порядку

	mu   sync.Mutex
	used []admission
}

type admission struct {
	at    time.Time
	units int
}

func NewLimiter(cfg LimiterConfig) (*Limiter, error) {
	if cfg.Budget == 0 {
		cfg.Budget = DefaultBudget
	}
	if cfg.Window == 0 {
		cfg.Window = DefaultWindow
	}
	if cfg.Concurrency == 0 {
		cfg.Concurrency = DefaultPrincipalConcurrency
	}
	if cfg.Budget < 0 || cfg.Window < 0 || cfg.Concurrency < 0 {
		return nil, fmt.Errorf("moysklad: limiter budget, window and concurrency must not be negative")
	}
	if cfg.Profile < ProfileUnknown || cfg.Profile > ProfileEmployee {
		return nil, fmt.Errorf("moysklad: limiter profile must be set explicitly")
	}
	if cfg.Account == nil {
		cfg.Account, _ = NewAccount(DefaultAccountConcurrency)
	}
	if cfg.Clock == nil {
		cfg.Clock = systemClock{}
	}
	if cfg.Weight == nil {
		cfg.Weight = DefaultWeight
	}
	return &Limiter{cfg: cfg, slots: make(chan struct{}, cfg.Concurrency), turn: make(chan struct{}, 1)}, nil
}

var moscow = time.FixedZone("MSK", 3*60*60)

// DefaultWeight — вес по расписанию из _restrictions.md. Источник не указывает часовой пояс дат:
// SDK считает их полночью по Москве. Вес 5 — только у двух точных GET-endpoints отчёта об
// остатках, независимо от авторизации; источник вводит его «с февраля 2026» без дня, SDK
// применяет его всегда. Другое расписание задаётся LimiterConfig.Weight.
func DefaultWeight(profile Profile, method, path string, now time.Time) int {
	if method == http.MethodGet && (path == "report/stock/all" || path == "report/stock/bystore") {
		return 5
	}
	since := func(year int, month time.Month, day int) bool {
		return !now.Before(time.Date(year, month, day, 0, 0, 0, 0, moscow))
	}
	switch profile {
	case ProfileApp:
		return 1
	case ProfileDraftApp:
		if since(2026, time.April, 13) {
			return 4
		}
		return 1
	case ProfileEmployee:
		switch {
		case since(2026, time.December, 1):
			return 4
		case since(2026, time.September, 1):
			return 3
		case since(2026, time.May, 12):
			return 2
		}
		return 1
	}
	return 4
}

// acquire ждёт слот principal, бюджет и слот аккаунта для одной попытки запроса: ожидание
// бюджета не занимает слот аккаунта, нужный другим principal.
// release нужно вызвать после закрытия body ответа; потраченные единицы не возвращаются.
func (l *Limiter) acquire(ctx context.Context, method, path string, principal bool) (release func(), err error) {
	var held []chan struct{}
	release = func() {
		for _, slots := range held {
			<-slots
		}
	}
	take := func(slots chan struct{}) error {
		select {
		case slots <- struct{}{}:
			held = append(held, slots)
			return nil
		case <-ctx.Done():
			release()
			return fmt.Errorf("moysklad: wait for request slot: %w", ctx.Err())
		}
	}
	if principal {
		if err := take(l.slots); err != nil {
			return nil, err
		}
		if err := l.waitBudget(ctx, method, path); err != nil {
			release()
			return nil, err
		}
	}
	if err := take(l.cfg.Account.slots); err != nil {
		return nil, err
	}
	return release, nil
}

func (l *Limiter) waitBudget(ctx context.Context, method, path string) error {
	select {
	case l.turn <- struct{}{}:
	case <-ctx.Done():
		return fmt.Errorf("moysklad: wait for rate limit: %w", ctx.Err())
	}
	defer func() { <-l.turn }()

	for {
		now := l.cfg.Clock.Now()
		// Вес больше бюджета не поместился бы никогда; вес меньше 1 добавил бы бюджет.
		units := min(max(l.cfg.Weight(l.cfg.Profile, method, path, now), 1), l.cfg.Budget)
		wait := l.reserve(now, units)
		if wait <= 0 {
			return nil
		}
		select {
		case <-l.cfg.Clock.After(wait):
		case <-ctx.Done():
			return fmt.Errorf("moysklad: wait for rate limit: %w", ctx.Err())
		}
	}
}

// reserve списывает units, если они помещаются в окно, иначе возвращает время до освобождения.
func (l *Limiter) reserve(now time.Time, units int) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	start := now.Add(-l.cfg.Window)
	for len(l.used) > 0 && !l.used[0].at.After(start) {
		l.used = l.used[1:]
	}
	total := 0
	for _, a := range l.used {
		total += a.units
	}
	if total+units <= l.cfg.Budget {
		l.used = append(l.used, admission{at: now, units: units})
		return 0
	}
	// Ждём, пока из окна выйдет столько ранних списаний, чтобы units поместились.
	free := l.cfg.Budget - total
	for _, a := range l.used {
		free += a.units
		if free >= units {
			return a.at.Add(l.cfg.Window).Sub(now)
		}
	}
	return l.cfg.Window
}

type systemClock struct{}

func (systemClock) Now() time.Time                         { return time.Now() }
func (systemClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
