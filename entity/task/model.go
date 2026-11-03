// Package task — задачи (entity/task), их комментарии (entity/task/{id}/notes) и типы задач
// (entity/task/metadata/states). Для создания, изменения и удаления задач нужна тарифная опция
// CRM; права и опцию проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task
package task

import "github.com/eteeet/moysklad-sdk-go/model"

// Task — задача. Agent (контрагент или юрлицо) и Operation (документ) — ссылки без expand;
// источник допускает привязку только к одному из них, проверяет сервер. Completed и Implementer
// приходят у выполненной задачи. Notes — коллекция комментариев: без expand только meta,
// не больше 100 комментариев хранится вместе с задачей (остальные — ListComments).
// Files — только metadata коллекции файлов.
type Task struct {
	Meta              model.Meta                         `json:"meta"`
	ID                string                             `json:"id"`
	AccountID         string                             `json:"accountId"`
	Description       string                             `json:"description"`
	Done              bool                               `json:"done"`
	DueToDate         model.DateTime                     `json:"dueToDate,omitzero"`
	Created           model.DateTime                     `json:"created,omitzero"`
	Updated           model.DateTime                     `json:"updated,omitzero"`
	Completed         model.DateTime                     `json:"completed,omitzero"`
	Assignee          *model.Reference                   `json:"assignee,omitzero"`
	Author            *model.Reference                   `json:"author,omitzero"`
	AuthorApplication *model.Reference                   `json:"authorApplication,omitzero"`
	Implementer       *model.Reference                   `json:"implementer,omitzero"`
	Agent             *model.Reference                   `json:"agent,omitzero"`
	Operation         *model.Reference                   `json:"operation,omitzero"`
	State             *model.Reference                   `json:"state,omitzero"`
	Notes             *model.Collection[Comment]         `json:"notes,omitzero"`
	Files             *model.Collection[model.Reference] `json:"files,omitzero"`
}

// Comment — комментарий задачи (meta.type tasknote). Текст — поле text, как в запросах и
// примерах источника (таблица атрибутов называет его description). Text может содержать
// упоминания {{employee;uuid}}: SDK не разбирает текст. Author, AuthorApplication и Moment
// только читаются; ID, AccountID и Files есть в примерах ответа, но не в таблице.
type Comment struct {
	Meta              model.Meta                         `json:"meta"`
	ID                string                             `json:"id"`
	AccountID         string                             `json:"accountId"`
	Text              string                             `json:"text"`
	Moment            model.DateTime                     `json:"moment,omitzero"`
	Author            *model.Reference                   `json:"author,omitzero"`
	AuthorApplication *model.Reference                   `json:"authorApplication,omitzero"`
	Files             *model.Collection[model.Reference] `json:"files,omitzero"`
}
