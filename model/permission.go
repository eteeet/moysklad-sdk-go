package model

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// Permissions — права по имени: сущности (supply, task …) или пользовательского права
// (importData, admin …). Набор имён открыт.
// Противоречия источника: таблицы называют поле Array(Object), но все примеры роли, прав
// сотрудника и контекста — объект; текст контекста пишет действие VIEW, примеры — view;
// таблица имён — customerOrder, пример — customerorder. SDK хранит ключи как пришли, без
// нормализации регистра, и следует примерам.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/custom-role
type Permissions map[string]Permission

// Permission — значение одного права: флаг (у пользовательских прав роли) или области по
// действиям (view, create, update, delete, print, approve, у задач done). Значение другой
// формы сохраняется в Raw и записывается обратно без изменений. Заполняется одно поле; если
// заполнено несколько, при записи используется первое из Raw, Flag, Actions.
type Permission struct {
	Flag    *bool
	Actions map[string]PermissionScope
	Raw     jsontext.Value
}

// PermissionScope — на какие объекты распространяется действие; набор значений открыт.
type PermissionScope string

const (
	ScopeNo             PermissionScope = "NO"
	ScopeOwn            PermissionScope = "OWN"
	ScopeOwnShared      PermissionScope = "OWN_SHARED"
	ScopeOwnGroup       PermissionScope = "OWN_GROUP"
	ScopeOwnGroupShared PermissionScope = "OWN_GROUP_SHARED"
	ScopeAll            PermissionScope = "ALL"
	// Только для задач.
	ScopeAuthor           PermissionScope = "AUTHOR"
	ScopeAssignee         PermissionScope = "ASSIGNEE"
	ScopeAuthorOrAssignee PermissionScope = "AUTHOR_OR_ASSIGNEE"
)

func (p Permission) MarshalJSON() ([]byte, error) {
	switch {
	case p.Raw != nil:
		return bytes.Clone(p.Raw), nil
	case p.Flag != nil:
		return json.Marshal(*p.Flag)
	case p.Actions != nil:
		return json.Marshal(p.Actions)
	}
	return []byte("null"), nil
}

func (p *Permission) UnmarshalJSON(data []byte) error {
	*p = Permission{}
	value := jsontext.Value(bytes.Clone(data))
	switch value.Kind() {
	case 't', 'f':
		var flag bool
		if err := json.Unmarshal(value, &flag); err != nil {
			return err
		}
		p.Flag = &flag
		return nil
	case '{':
		var actions map[string]PermissionScope
		if err := json.Unmarshal(value, &actions); err == nil {
			p.Actions = actions
			return nil
		}
	case 'n':
		return nil
	}
	// Декодер передаёт сюда только корректный JSON: неизвестную форму сохраняем как есть.
	p.Raw = value
	return nil
}
