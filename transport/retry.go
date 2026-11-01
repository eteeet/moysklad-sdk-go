package transport

import (
	"context"
	"errors"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
)

const (
	DefaultMaxAttempts = 3
	DefaultBaseDelay   = 500 * time.Millisecond
	DefaultMaxDelay    = 30 * time.Second
	DefaultMaxWait     = time.Minute
)

// RetryPolicy повторяет только чтения GET/HEAD при статусах retryStatuses и временном сбое
// соединения (обрыв, отказ, таймаут). Запись, выпуск токена и async-запросы не повторяются:
// повтор мог бы выполнить действие дважды. 429 с кодом 1073 (параллельные запросы) повторяется
// как любой 429 и не меняет бюджет Limiter. Числа по умолчанию — выбор SDK, а не лимиты API.
type RetryPolicy struct {
	MaxAttempts int           // всего попыток, включая первую; 0 — DefaultMaxAttempts
	BaseDelay   time.Duration // первая пауза экспоненты; 0 — DefaultBaseDelay
	MaxDelay    time.Duration // потолок экспоненты; подсказку сервера не обрезает; 0 — DefaultMaxDelay
	// Если нужная пауза (с учётом подсказки сервера) больше MaxWait, повтора нет; 0 — DefaultMaxWait.
	MaxWait time.Duration
	Clock   Clock
	// Rand возвращает [0, 1) для full jitter; nil — math/rand/v2.
	Rand func() float64
}

func (p *RetryPolicy) withDefaults() RetryPolicy {
	r := *p
	if r.MaxAttempts <= 0 {
		r.MaxAttempts = DefaultMaxAttempts
	}
	if r.BaseDelay <= 0 {
		r.BaseDelay = DefaultBaseDelay
	}
	if r.MaxDelay <= 0 {
		r.MaxDelay = DefaultMaxDelay
	}
	if r.MaxWait <= 0 {
		r.MaxWait = DefaultMaxWait
	}
	if r.Clock == nil {
		r.Clock = systemClock{}
	}
	if r.Rand == nil {
		r.Rand = rand.Float64
	}
	return r
}

var retryStatuses = map[int]bool{
	http.StatusTooManyRequests:     true,
	http.StatusInternalServerError: true,
	http.StatusBadGateway:          true,
	http.StatusServiceUnavailable:  true,
	http.StatusGatewayTimeout:      true,
}

// async=true создаёт задачу даже у GET; он может прийти и в query, и в строке endpoint.
func retryableCall(method, endpoint string, query url.Values, options RequestOptions) bool {
	if method != http.MethodGet && method != http.MethodHead || options.NoRetry || query.Get("async") == "true" {
		return false
	}
	u, err := url.Parse(endpoint)
	return err == nil && !strings.Contains(strings.ToLower(u.RawQuery), "async=true")
}

// retryDelay решает, повторять ли попытку attempt (с 1), и возвращает паузу.
func (p RetryPolicy) retryDelay(attempt int, resp *http.Response, err error) (time.Duration, bool) {
	if attempt >= p.MaxAttempts || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return 0, false
	}
	var hint time.Duration
	switch {
	case resp == nil:
		if !temporaryNetError(err) {
			return 0, false
		}
	case retryStatuses[resp.StatusCode]:
		// Только чистая APIError: ошибка чтения или неизвестная кодировка ответа не повторяются.
		var apiErr *sdkerrors.APIError
		if !errors.As(err, &apiErr) || error(apiErr) != err {
			return 0, false
		}
		info := ResponseInfo{Header: resp.Header}
		for _, h := range []func() (time.Duration, bool){info.LognexRetryAfter, info.ResetAfter,
			func() (time.Duration, bool) { return info.RetryAfter(p.Clock.Now()) }} {
			if d, ok := h(); ok {
				hint = max(hint, d)
			}
		}
	default:
		// 2xx с ошибкой decode, 4xx и прочие ответы не повторяются.
		return 0, false
	}
	backoff := min(p.BaseDelay<<(attempt-1), p.MaxDelay)
	if backoff <= 0 { // переполнение сдвига
		backoff = p.MaxDelay
	}
	delay := max(hint, time.Duration(p.Rand()*float64(backoff)))
	return delay, delay <= p.MaxWait
}

// Временные сбои соединения; TLS, неверный URL и прочие постоянные ошибки не повторяются.
func temporaryNetError(err error) bool {
	var netErr net.Error
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.EPIPE) ||
		errors.As(err, &netErr) && netErr.Timeout()
}

// sleep ждёт паузу; false — пауза не уложится в дедлайн ctx или ctx отменён (тогда ctxErr не nil).
func (p RetryPolicy) sleep(ctx context.Context, d time.Duration) (waited bool, ctxErr error) {
	// Дедлайн контекста — по системному времени, независимо от Clock.
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < d {
		return false, nil
	}
	select {
	case <-p.Clock.After(d):
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}
