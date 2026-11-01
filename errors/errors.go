package errors

// Документация: https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-obrabotka-oshibok
import (
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/model"
)

type Detail struct {
	Message      string       `json:"error"`
	Code         int          `json:"code,omitzero"`
	Parameter    string       `json:"parameter,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
	MoreInfo     string       `json:"moreInfo,omitempty"`
	Line         int          `json:"line,omitzero"`
	Column       int          `json:"column,omitzero"`
	Dependencies []model.Meta `json:"dependencies,omitempty"`
	Meta         *model.Meta  `json:"meta,omitzero"`
}

// StatusCode заполняет transport из HTTP-ответа. Errors пуст, если ответ пришёл
// без JSON-тела (например, 415 без gzip).
type APIError struct {
	StatusCode int      `json:"-"`
	Errors     []Detail `json:"errors"`
}

func (e *APIError) Error() string {
	prefix := fmt.Sprintf("moysklad: HTTP %d", e.StatusCode)
	messages := make([]string, 0, len(e.Errors))

	for _, detail := range e.Errors {
		message := detail.Message
		if message == "" {
			message = detail.ErrorMessage
		}
		if detail.Code != 0 {
			message = fmt.Sprintf("[%d] %s", detail.Code, message)
		}
		if message = strings.TrimSpace(message); message != "" {
			messages = append(messages, message)
		}
	}

	if len(messages) == 0 {
		return prefix
	}
	return prefix + ": " + strings.Join(messages, "; ")
}

// HasCode проверяет детали этого APIError; Code 0 (detail без кода) не совпадает никогда.
func (e *APIError) HasCode(code int) bool {
	if code == 0 {
		return false
	}
	for _, detail := range e.Errors {
		if detail.Code == code {
			return true
		}
	}
	return false
}

// HasCode проверяет все APIError в цепочке err, включая каждую ветку errors.Join.
// Классифицировать ошибки нужно по коду, а не по локализованному тексту.
func HasCode(err error, code int) bool {
	switch e := err.(type) {
	case nil:
		return false
	case *APIError:
		if e.HasCode(code) {
			return true
		}
	case interface{ Unwrap() []error }:
		for _, branch := range e.Unwrap() {
			if HasCode(branch, code) {
				return true
			}
		}
		return false
	}
	return HasCode(stderrors.Unwrap(err), code)
}

// StatusCode возвращает HTTP-статус первого APIError в цепочке err или 0, если APIError нет.
// HTTP-статус и код ошибки МойСклад — независимые признаки.
func StatusCode(err error) int {
	var apiErr *APIError
	if stderrors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return 0
}
