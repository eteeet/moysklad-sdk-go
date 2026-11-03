package counterparty

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Note — событие контрагента (meta.type note). Agent, Author и AuthorApplication только
// читаются; AuthorApplication есть у событий, созданных решением.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-zaprosy-sobytiya-kontragenta
type Note struct {
	Meta              model.Meta       `json:"meta"`
	ID                string           `json:"id"`
	AccountID         string           `json:"accountId"`
	Created           model.DateTime   `json:"created,omitzero"`
	Description       string           `json:"description"`
	Agent             *model.Reference `json:"agent,omitzero"`
	Author            *model.Reference `json:"author,omitzero"`
	AuthorApplication *model.Reference `json:"authorApplication,omitzero"`
}

// NoteRequest — единственное записываемое поле события: текст description (обязателен при
// создании).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-dobavit-sobytie
type NoteRequest struct {
	Description string `json:"description"`
}

func (s *Service) notes(counterpartyID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "counterparty", ParentID: counterpartyID, Path: "notes", Item: "note", MetaType: "note"}
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-spisok-sobytij
func (s *Service) ListNotes(ctx context.Context, counterpartyID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Note], error) {
	return nested.List[Note](ctx, s.notes(counterpartyID), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-poluchit-sobytie
func (s *Service) GetNote(ctx context.Context, counterpartyID, noteID string, opts ...transport.RequestOption) (*Note, error) {
	return nested.Get[Note](ctx, s.notes(counterpartyID), noteID, opts)
}

// CreateNote — пример ответа источника — массив из одного события; принимается и объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-dobavit-sobytie
func (s *Service) CreateNote(ctx context.Context, counterpartyID string, input NoteRequest, opts ...transport.RequestOption) (*Note, error) {
	if input.Description == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create counterparty %s note: description must not be empty", counterpartyID)
	}
	return nested.Create[Note](ctx, s.notes(counterpartyID), input, opts)
}

// UpdateNote заменяет текст события; пустой description — ошибка до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-sobytie
func (s *Service) UpdateNote(ctx context.Context, counterpartyID, noteID string, input NoteRequest, opts ...transport.RequestOption) (*Note, error) {
	return nested.Update[Note](ctx, s.notes(counterpartyID), noteID, input, opts)
}

// DeleteNote — DELETE события. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-udalit-sobytie
func (s *Service) DeleteNote(ctx context.Context, counterpartyID, noteID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.notes(counterpartyID), noteID, opts)
}
