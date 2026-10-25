package errors

// Документация: https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-obrabotka-oshibok
import (
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
