package moysklad

import (
	"fmt"
	"time"
)

// Error — ответ МойСклад со статусом 4xx или 5xx. Достаётся из цепочки через errors.As.
type Error struct {
	StatusCode int
	// Code — код первой ошибки МойСклад, например 1049 (превышена частота) или 1073 (параллельность).
	Code int
	// Message — текст первой ошибки.
	Message string
	// Errors — все ошибки из тела ответа.
	Errors []ErrorItem

	retryAfter time.Duration // из X-Lognex-Retry-After
}

func (e *Error) Error() string {
	return fmt.Sprintf("moysklad: HTTP %d: %s (код %d)", e.StatusCode, e.Message, e.Code)
}
