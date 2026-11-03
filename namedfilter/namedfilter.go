// Package namedfilter читает сохранённые фильтры списков: entity/{entityType}/namedfilter.
// Фильтр отдаётся только как id и наименование, параметры фильтрации API не возвращает.
// Создание, изменение и удаление источник не описывает — пакет только читает. Применяется
// фильтр в List сущности или документа того же типа: query.ListOptions{NamedFilter: f.Meta.Href}.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/named-filter
package namedfilter

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// NamedFilter — сохранённый фильтр. Все поля только для чтения.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/named-filter#3-sohranennyj-filtr
type NamedFilter struct {
	Meta      model.Meta       `json:"meta"`
	ID        string           `json:"id"`
	AccountID string           `json:"accountId"`
	Name      string           `json:"name"`
	Owner     *model.Reference `json:"owner,omitzero"` // сотрудник, без expand
}

// ListOptions — параметры списка. limit, offset и другие параметры источник для этого
// endpoint не описывает, поэтому их нет.
type ListOptions struct {
	// Owner — абсолютный href сотрудника, чьи фильтры нужны; пусто — фильтры текущего
	// пользователя. Чужие фильтры доступны администратору или решению — права проверяет сервер.
	Owner string
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: namedfilter service requires a client")
	}
	return &Service{client: client}, nil
}

// Тип сущности уходит в путь запроса: только латиница в нижнем регистре и цифры, как в metadata.
var entityTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// List читает одну страницу сохранённых фильтров типа entityType (например "product",
// "supply"); общего списка по всем типам API не даёт. Один GET, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/named-filter#3-poluchit-spisok-filtrov
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/named-filter#3-poluchit-spisok-filtrov-drugogo-polzovatelya
func (s *Service) List(ctx context.Context, entityType string, options ListOptions, opts ...transport.RequestOption) (*model.Collection[NamedFilter], error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: get namedfilters: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	var params url.Values
	if options.Owner != "" {
		if u, err := url.Parse(options.Owner); err != nil || !u.IsAbs() {
			return nil, fmt.Errorf("moysklad: get %s namedfilters: owner must be an absolute employee href", entityType)
		}
		params = url.Values{query.ParamOwner: {options.Owner}}
	}
	return list.Page[NamedFilter](ctx, s.client, "get "+entityType+" namedfilters", "entity/"+entityType+"/namedfilter", params, opts...)
}

// Get читает один сохранённый фильтр типа entityType по id.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/named-filter#3-poluchit-filtr-po-id
func (s *Service) Get(ctx context.Context, entityType, id string, opts ...transport.RequestOption) (*NamedFilter, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: get namedfilter: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get namedfilter: id %q must be a UUID", id)
	}
	operation := "get " + entityType + " namedfilter " + id
	var result NamedFilter
	resp, err := s.client.Send(ctx, http.MethodGet, "entity/"+entityType+"/namedfilter/"+id, nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
