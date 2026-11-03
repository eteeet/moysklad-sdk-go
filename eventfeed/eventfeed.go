// Package eventfeed — лента событий документа: entity/{type}/{id}/notes (meta.type eventnote).
// Источник перечисляет purchasereturn, commissionreportout, customerorder, purchaseorder, supply,
// invoiceout и invoicein; тип документа SDK проверяет только по форме, допустимость — сервер.
// События контрагента — counterparty.Service (ListNotes и др.), комментарии задачи — task.Service.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed
package eventfeed

import (
	"context"
	"fmt"
	"regexp"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

var documentType = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// Note — событие документа. Description может содержать упоминания {{employee;uuid}}: SDK
// передаёт и возвращает текст как есть, без разбора. Author и AuthorApplication только
// читаются; AuthorApplication есть у событий, созданных решением.
type Note struct {
	Meta              model.Meta       `json:"meta"`
	ID                string           `json:"id"`
	AccountID         string           `json:"accountId"`
	Created           model.DateTime   `json:"created,omitzero"`
	Description       string           `json:"description"`
	Author            *model.Reference `json:"author,omitzero"`
	AuthorApplication *model.Reference `json:"authorApplication,omitzero"`
}

// NoteRequest — единственное записываемое поле события: description (String(4096), необходим
// при создании; длину проверяет сервер).
type NoteRequest struct {
	Description string `json:"description"`
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: eventfeed service requires a client")
	}
	return &Service{client: client}, nil
}

// notes проверяет тип документа до HTTP; при ошибке сбрасывает ResponseInfo сам, иначе это
// делает вызываемая функция nested.
func (s *Service) notes(operation, typ, id string, opts []transport.RequestOption) (nested.Collection, error) {
	if !documentType.MatchString(typ) {
		transport.ResetResponseInfo(opts)
		return nested.Collection{}, fmt.Errorf("moysklad: %s notes: document type %q must match [a-z][a-z0-9]*", operation, typ)
	}
	return nested.Collection{Client: s.client, Parent: typ, ParentID: id, Path: "notes", Item: "note", MetaType: "eventnote", Created: true}, nil
}

// List читает одну страницу событий документа: только limit (1–1000) и offset.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-poluchit-spisok-sobytij-dokumenta
func (s *Service) List(ctx context.Context, documentType, documentID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Note], error) {
	c, err := s.notes("list", documentType, documentID, opts)
	if err != nil {
		return nil, err
	}
	return nested.List[Note](ctx, c, page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-poluchit-sobytie
func (s *Service) Get(ctx context.Context, documentType, documentID, noteID string, opts ...transport.RequestOption) (*Note, error) {
	c, err := s.notes("get", documentType, documentID, opts)
	if err != nil {
		return nil, err
	}
	return nested.Get[Note](ctx, c, noteID, opts)
}

// Create добавляет событие; пустой description — ошибка до запроса. Пример ответа источника —
// массив из одного события; принимается и объект. Источник отвечает HTTP 201 — SDK принимает 201 и 200
// (nested.Collection.Created). Не повторять POST автоматически после
// неоднозначного сбоя: возможен дубль события.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-dobavit-sobytie
func (s *Service) Create(ctx context.Context, documentType, documentID string, input NoteRequest, opts ...transport.RequestOption) (*Note, error) {
	c, err := s.notes("create", documentType, documentID, opts)
	if err != nil {
		return nil, err
	}
	if input.Description == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create %s %s note: description must not be empty", documentType, documentID)
	}
	return nested.Create[Note](ctx, c, input, opts)
}

// Update заменяет текст события; пустой description — ошибка до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-obnovit-sobytie
func (s *Service) Update(ctx context.Context, documentType, documentID, noteID string, input NoteRequest, opts ...transport.RequestOption) (*Note, error) {
	c, err := s.notes("update", documentType, documentID, opts)
	if err != nil {
		return nil, err
	}
	return nested.Update[Note](ctx, c, noteID, input, opts)
}

// Delete — DELETE события (источник: 204). Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-udalit-sobytie
func (s *Service) Delete(ctx context.Context, documentType, documentID, noteID string, opts ...transport.RequestOption) error {
	c, err := s.notes("delete", documentType, documentID, opts)
	if err != nil {
		return err
	}
	return nested.Delete(ctx, c, noteID, opts)
}

// CreateMany и UpdateMany — один POST коллекции с массивом: элемент без meta создаёт событие,
// с meta (тип eventnote) — изменяет. Правила ответа — пакет batch; HTTP 201 по источнику (принимается и 200). Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-massovoe-sozdanie-i-obnovlenie-sobytij
func (s *Service) CreateMany(ctx context.Context, documentType, documentID string, items []NoteRequest, opts ...transport.RequestOption) ([]batch.Result[Note], error) {
	c, err := s.notes("create many", documentType, documentID, opts)
	if err != nil {
		return nil, err
	}
	return nested.CreateMany[Note](ctx, c, items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-massovoe-sozdanie-i-obnovlenie-sobytij
func (s *Service) UpdateMany(ctx context.Context, documentType, documentID string, items []batch.Update[NoteRequest], opts ...transport.RequestOption) ([]batch.Result[Note], error) {
	c, err := s.notes("update many", documentType, documentID, opts)
	if err != nil {
		return nil, err
	}
	return nested.UpdateMany[Note](ctx, c, items, opts)
}

// DeleteMany — один POST notes/delete с meta событий по ids. Источник отвечает 200 с пустым
// телом: результата по элементам нет. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed#3-massovoe-udalenie-sobytij
func (s *Service) DeleteMany(ctx context.Context, documentType, documentID string, ids []string, opts ...transport.RequestOption) error {
	c, err := s.notes("delete many", documentType, documentID, opts)
	if err != nil {
		return err
	}
	return nested.DeleteRows(ctx, c, "delete many "+documentType+" "+documentID+" notes", ids, opts)
}
