package transport

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu      sync.Mutex
	now     time.Time
	waiters []fakeWaiter
	added   chan time.Duration // каждый After сообщает о себе: тесты ждут без sleep
	nows    chan time.Time     // если задан, каждый Now сообщает о себе
}

type fakeWaiter struct {
	at time.Time
	ch chan time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	now, nows := c.now, c.nows
	c.mu.Unlock()
	if nows != nil {
		nows <- now
	}
	return now
}

func (c *fakeClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := make(chan time.Time, 1)
	c.waiters = append(c.waiters, fakeWaiter{c.now.Add(d), ch})
	c.added <- d
	return ch
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
	kept := c.waiters[:0]
	for _, w := range c.waiters {
		if w.at.After(c.now) {
			kept = append(kept, w)
		} else {
			w.ch <- c.now
		}
	}
	c.waiters = kept
}

// waitAfter ждёт следующий вызов After и возвращает его длительность.
func (c *fakeClock) waitAfter(t *testing.T) time.Duration {
	t.Helper()
	select {
	case d := <-c.added:
		return d
	case <-t.Context().Done():
		t.Fatal("no clock waiter")
		return 0
	}
}

var testStart = time.Date(2026, 11, 1, 12, 0, 0, 0, moscow)

func newTestLimiter(t *testing.T, cfg LimiterConfig) (*Limiter, *fakeClock) {
	t.Helper()
	clock := &fakeClock{now: testStart, added: make(chan time.Duration, 100)}
	cfg.Clock = clock
	l, err := NewLimiter(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return l, clock
}

func TestDefaultWeight(t *testing.T) {
	at := func(y int, m time.Month, d int, sec int) time.Time {
		return time.Date(y, m, d, 0, 0, 0, 0, moscow).Add(time.Duration(sec) * time.Second)
	}
	for _, tc := range []struct {
		profile Profile
		method  string
		path    string
		now     time.Time
		want    int
	}{
		{ProfileApp, "GET", "entity/product", testStart, 1},
		{ProfileDraftApp, "GET", "entity/product", at(2026, 4, 13, -1), 1},
		{ProfileDraftApp, "GET", "entity/product", at(2026, 4, 13, 0), 4},
		{ProfileEmployee, "GET", "entity/product", at(2026, 5, 12, -1), 1},
		{ProfileEmployee, "GET", "entity/product", at(2026, 5, 12, 0), 2},
		{ProfileEmployee, "GET", "entity/product", at(2026, 9, 1, 0), 3},
		{ProfileEmployee, "POST", "entity/product", at(2026, 11, 30, 86399), 3},
		{ProfileEmployee, "GET", "entity/product", at(2026, 12, 1, 0), 4},
		// Граница в московском времени, а не UTC.
		{ProfileEmployee, "GET", "entity/product", time.Date(2026, 11, 30, 21, 0, 0, 0, time.UTC), 4},
		{ProfileUnknown, "GET", "entity/product", testStart, 4},
		{ProfileApp, "GET", "report/stock/all", testStart, 5},
		{ProfileEmployee, "GET", "report/stock/bystore", testStart, 5},
		{ProfileApp, "GET", "report/stock/all/current", testStart, 1},
		{ProfileApp, "GET", "report/stock/bystore/current", testStart, 1},
		{ProfileApp, "POST", "report/stock/all", testStart, 1},
	} {
		if got := DefaultWeight(tc.profile, tc.method, tc.path, tc.now); got != tc.want {
			t.Errorf("DefaultWeight(%d, %s %s, %v) = %d, want %d", tc.profile, tc.method, tc.path, tc.now, got, tc.want)
		}
	}
}

func acquireAsync(ctx context.Context, l *Limiter, path string) chan error {
	done := make(chan error, 1)
	go func() {
		release, err := l.acquire(ctx, http.MethodGet, path, true)
		if err == nil {
			release()
		}
		done <- err
	}()
	return done
}

func mustAcquire(t *testing.T, l *Limiter, principal bool) func() {
	t.Helper()
	release, err := l.acquire(t.Context(), http.MethodGet, "x", principal)
	if err != nil {
		t.Fatal(err)
	}
	return release
}

func TestLimiterBudgetWindow(t *testing.T) {
	for _, tc := range []struct {
		profile Profile
		fits    int
	}{
		{ProfileApp, 45},
		{ProfileEmployee, 15}, // вес 3 на 2026-11-01
		{ProfileUnknown, 11},
	} {
		l, clock := newTestLimiter(t, LimiterConfig{Profile: tc.profile, Concurrency: 100})
		for range tc.fits {
			mustAcquire(t, l, true)() // слот освобождён, единицы окна — нет
		}
		done := acquireAsync(t.Context(), l, "entity/product")
		if wait := clock.waitAfter(t); wait != 3*time.Second {
			t.Fatalf("profile %d: wait = %v", tc.profile, wait)
		}
		clock.Advance(3*time.Second - time.Nanosecond)
		select {
		case err := <-done:
			t.Fatalf("admitted before window end: %v", err)
		default:
		}
		clock.Advance(time.Nanosecond)
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
}

// Отменённый ожидающий не тратит бюджет, следующий в очереди получает его.
func TestLimiterQueueCancel(t *testing.T) {
	l, clock := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Budget: 2, Concurrency: 10})
	mustAcquire(t, l, true)()
	mustAcquire(t, l, true)()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	canceled := acquireAsync(ctx, l, "a")
	clock.waitAfter(t) // canceled держит очередь и ждёт окно
	second := acquireAsync(t.Context(), l, "b")
	cancel()
	if err := <-canceled; !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
	clock.waitAfter(t) // second занял очередь после отменённого
	clock.Advance(3 * time.Second)
	if err := <-second; err != nil {
		t.Fatal(err)
	}
	l.mu.Lock()
	used := len(l.used)
	l.mu.Unlock()
	if used != 1 {
		t.Fatalf("admissions in window = %d, want 1 (canceled waiter spent nothing)", used)
	}
}

// Ожидание бюджета не занимает слот аккаунта.
func TestLimiterBudgetWaitKeepsAccountSlot(t *testing.T) {
	account, err := NewAccount(1)
	if err != nil {
		t.Fatal(err)
	}
	l, clock := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Budget: 1, Account: account})
	mustAcquire(t, l, true)()
	waiting := acquireAsync(t.Context(), l, "x")
	clock.waitAfter(t)
	other, _ := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Account: account})
	mustAcquire(t, other, true)()
	clock.Advance(3 * time.Second)
	if err := <-waiting; err != nil {
		t.Fatal(err)
	}
}

func TestLimiterConcurrency(t *testing.T) {
	account, err := NewAccount(3)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Concurrency: 2, Account: account})
	b, _ := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Concurrency: 2, Account: account})
	// Отменённый контекст: занятый слот не дождаться, свободный выбирается без ожидания.
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	blocked := func(l *Limiter, principal bool) bool {
		release, err := l.acquire(canceled, "GET", "x", principal)
		if err == nil {
			release()
		}
		return errors.Is(err, context.Canceled)
	}

	r1 := mustAcquire(t, a, true)
	r2 := mustAcquire(t, a, true)
	if len(a.slots) != 2 || !blocked(a, true) {
		t.Fatal("principal limit not enforced")
	}
	r3 := mustAcquire(t, b, true)
	if len(account.slots) != 3 || !blocked(b, false) {
		t.Fatal("shared account limit not enforced")
	}
	r1()
	r4 := mustAcquire(t, b, false) // без principal: только слот аккаунта
	if len(b.slots) != 1 {
		t.Fatalf("principal slot taken by account-only call: %d", len(b.slots))
	}
	for _, r := range []func(){r2, r3, r4} {
		r()
	}
	if len(a.slots)+len(b.slots)+len(account.slots) != 0 {
		t.Fatal("slots leaked")
	}
}

func TestLimiterWeightClamp(t *testing.T) {
	for _, weight := range []int{-5, 0, 100} {
		l, _ := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Budget: 3,
			Weight: func(Profile, string, string, time.Time) int { return weight }})
		mustAcquire(t, l, true)()
		l.mu.Lock()
		got := l.used[0].units
		l.mu.Unlock()
		if want := min(max(weight, 1), 3); got != want {
			t.Errorf("weight %d: units = %d, want %d", weight, got, want)
		}
	}
}

func TestNewLimiterRejects(t *testing.T) {
	for _, cfg := range []LimiterConfig{{}, {Profile: ProfileApp, Budget: -1}, {Profile: ProfileApp, Window: -1},
		{Profile: ProfileApp, Concurrency: -1}, {Profile: Profile(99)}} {
		if _, err := NewLimiter(cfg); err == nil {
			t.Errorf("%+v accepted", cfg)
		}
	}
	if _, err := NewAccount(0); err == nil {
		t.Error("zero account concurrency accepted")
	}
}

// Регрессия ревью 46e9d06: бюджет списывается по времени фактической отправки. Запрос, долго
// ждавший занятый другим principal слот аккаунта, не открывает окно для следующего сразу.
func TestLimiterBudgetChargedAtSend(t *testing.T) {
	account, err := NewAccount(1)
	if err != nil {
		t.Fatal(err)
	}
	l, clock := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Budget: 1, Account: account})
	other, _ := newTestLimiter(t, LimiterConfig{Profile: ProfileApp, Account: account})
	busy := mustAcquire(t, other, true) // аккаунт занят другим principal

	clock.mu.Lock()
	clock.nows = make(chan time.Time, 100)
	clock.mu.Unlock()
	first := acquireAsync(t.Context(), l, "x")
	<-clock.nows // окно проверено в 12:00:00, дальше — ожидание слота аккаунта
	clock.Advance(4 * time.Second)
	busy()
	if err := <-first; err != nil {
		t.Fatal(err)
	}

	second := acquireAsync(t.Context(), l, "x")
	select {
	case err := <-second:
		t.Fatalf("second request admitted without waiting (err = %v): the first was sent at 12:00:04", err)
	case wait := <-clock.added:
		if wait != 3*time.Second {
			t.Fatalf("second request waits %v, want 3s", wait)
		}
	}
	clock.Advance(3 * time.Second)
	if err := <-second; err != nil {
		t.Fatal(err)
	}
}
