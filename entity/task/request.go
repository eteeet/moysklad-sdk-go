package task

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: description и assignee необходимы при создании. Done — Field: Value(false)
// уходит явным false. Files и Notes не записываются (файлы — пакет files, комментарии —
// CreateComment); readonly поля отсутствуют.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-sozdat-zadachu
type CreateRequest struct {
	Description string                       `json:"description"`
	Assignee    model.Reference              `json:"assignee"`
	DueToDate   model.Field[model.DateTime]  `json:"dueToDate,omitzero"`
	Done        model.Field[bool]            `json:"done,omitzero"`
	Agent       model.Field[model.Reference] `json:"agent,omitzero"`
	Operation   model.Field[model.Reference] `json:"operation,omitzero"`
	State       model.Field[model.Reference] `json:"state,omitzero"`
}

// UpdateRequest — изменяемые поля задачи; незаданные не отправляются.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-izmenit-zadachu
type UpdateRequest struct {
	Description model.Field[string]          `json:"description,omitzero"`
	Assignee    model.Field[model.Reference] `json:"assignee,omitzero"`
	DueToDate   model.Field[model.DateTime]  `json:"dueToDate,omitzero"`
	Done        model.Field[bool]            `json:"done,omitzero"`
	Agent       model.Field[model.Reference] `json:"agent,omitzero"`
	Operation   model.Field[model.Reference] `json:"operation,omitzero"`
	State       model.Field[model.Reference] `json:"state,omitzero"`
}

// CommentRequest — текст комментария (text, String(4096)); необходим при создании.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-sozdat-kommentarij-zadachi
type CommentRequest struct {
	Text string `json:"text"`
}

// StateCreateRequest — новый тип задачи: name и color необходимы («необходимо и достаточно»).
// Color — ARGB целым числом из 4 байт (rgb(162, 198, 23) → 10667543); указатель отличает «не
// задан» от 0. StateType в таблице не помечен «только для чтения», но в описании создания не
// назван — Field, по умолчанию не отправляется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-sozdat-tip-zadachi
type StateCreateRequest struct {
	Name      string              `json:"name"`
	Color     *int64              `json:"color"`
	StateType model.Field[string] `json:"stateType,omitzero"`
}
