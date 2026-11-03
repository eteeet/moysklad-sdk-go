package print

import "github.com/eteeet/moysklad-sdk-go/model"

// Template — шаблон печатной формы, только чтение: изменить или удалить шаблон через JSON API
// нельзя. Type открыт: entity — шаблон документа, mxtemplate — ценников и этикеток; неизвестное
// значение не отвергается. Content — ссылка на скачивание самого шаблона.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/template
type Template struct {
	Meta    model.Meta `json:"meta"`
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Type    string     `json:"type"`
	Content string     `json:"content"`
}

// Result — ответ на запрос печати, без следования redirect.
//
//   - StatusCode 303 (http.StatusSeeOther): форма готова, Location — временный адрес файла,
//     доступный 5 минут; скачать — Download отдельным вызовом.
//   - StatusCode 202 (http.StatusAccepted): Location — адрес статуса печати. Источник называет
//     этот вариант будущим и протокол опроса не описывает: SDK не опрашивает и не ждёт.
//
// Location разрешён относительно адреса запроса; origin ссылки на файл источник не фиксирует.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/print
type Result struct {
	StatusCode int
	Location   string
}
