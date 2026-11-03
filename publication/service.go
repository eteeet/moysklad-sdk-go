// Package publication публикует печатные формы документов: entity/{type}/{id}/publication.
// Нужны права на чтение и печать типа; список типов с публикациями (заказы, счета, отгрузки,
// платежи, договор, розничная смена и др.) задаёт сервер — SDK проверяет только форму типа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/publication
package publication

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: publication service requires a client")
	}
	return &Service{client: client}, nil
}

var entityTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// List читает публикации документа: GET entity/{type}/{id}/publication. Параметров списка
// источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/publication
func (s *Service) List(ctx context.Context, entityType, id string, opts ...transport.RequestOption) (*model.Collection[Publication], error) {
	transport.ResetResponseInfo(opts)
	if err := checkDocument(entityType, id); err != nil {
		return nil, fmt.Errorf("moysklad: get publications: %w", err)
	}
	return list.Page[Publication](ctx, s.client, "get "+entityType+" "+id+" publications", endpoint(entityType, id), nil, opts...)
}

// Get читает публикацию: GET entity/{type}/{id}/publication/{publicationId}.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/publication
func (s *Service) Get(ctx context.Context, entityType, id, publicationID string, opts ...transport.RequestOption) (*Publication, error) {
	transport.ResetResponseInfo(opts)
	if err := checkPublication(entityType, id, publicationID); err != nil {
		return nil, fmt.Errorf("moysklad: get publication: %w", err)
	}
	operation := "get " + entityType + " " + id + " publication " + publicationID
	return s.do(ctx, operation, http.MethodGet, endpoint(entityType, id)+"/"+publicationID, nil, opts)
}

// Create публикует документ: POST entity/{type}/{id}/publication. Новая публикация — 201, уже
// существующая по этому шаблону — 200 с ней же; статус виден в transport.ResponseInfo. Запись не
// повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/publication
func (s *Service) Create(ctx context.Context, entityType, id string, input CreateRequest, opts ...transport.RequestOption) (*Publication, error) {
	transport.ResetResponseInfo(opts)
	if err := checkDocument(entityType, id); err != nil {
		return nil, fmt.Errorf("moysklad: create publication: %w", err)
	}
	operation := "create " + entityType + " " + id + " publication"
	if input.Template.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: template.meta.href is required", operation)
	}
	return s.do(ctx, operation, http.MethodPost, endpoint(entityType, id), input, opts)
}

// Delete удаляет публикацию: DELETE entity/{type}/{id}/publication/{publicationId}, успех — 204
// (или 200); повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/publication
func (s *Service) Delete(ctx context.Context, entityType, id, publicationID string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if err := checkPublication(entityType, id, publicationID); err != nil {
		return fmt.Errorf("moysklad: delete publication: %w", err)
	}
	operation := "delete " + entityType + " " + id + " publication " + publicationID
	return write.Delete(ctx, s.client, operation, endpoint(entityType, id)+"/"+publicationID, opts)
}

// do — один запрос; успех — 200 (и 201 у POST) с meta.href: id у публикации нет.
func (s *Service) do(ctx context.Context, operation, method, path string, body any, opts []transport.RequestOption) (*Publication, error) {
	var result Publication
	resp, err := s.client.Send(ctx, method, path, nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && (method != http.MethodPost || resp.StatusCode != http.StatusCreated) {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	return &result, nil
}

func endpoint(entityType, id string) string {
	return "entity/" + entityType + "/" + id + "/publication"
}

func checkDocument(entityType, id string) error {
	if !entityTypePattern.MatchString(entityType) {
		return fmt.Errorf("entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return fmt.Errorf("id %q must be a UUID", id)
	}
	return nil
}

func checkPublication(entityType, id, publicationID string) error {
	if err := checkDocument(entityType, id); err != nil {
		return err
	}
	if !validID(publicationID) {
		return fmt.Errorf("publication id %q must be a UUID", publicationID)
	}
	return nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
