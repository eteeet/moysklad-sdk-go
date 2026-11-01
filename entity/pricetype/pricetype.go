// Package pricetype читает и редактирует типы цен компании. Типы цен меняются только полным
// упорядоченным списком (ReplaceAll): одиночного создания и изменения API не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/pricetype
package pricetype

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// MaxPriceTypes — больше типов цен API создать не даёт.
const MaxPriceTypes = 100

// ReplaceItem — элемент списка для ReplaceAll. Существующий тип цены передаётся с Meta и
// наименованием (новым — для переименования), новый — без Meta. Наименование непустое у всех;
// уникальность проверяет сервер. ExternalCode — указатель, а не Field: удалить внешний код
// нельзя, null не отправляется; пустая строка — сервер сгенерирует случайный код.
type ReplaceItem struct {
	Meta         *model.Meta `json:"meta,omitzero"`
	Name         string      `json:"name"`
	ExternalCode *string     `json:"externalCode,omitzero"`
}

// Existing — элемент для типа цены, который остаётся в списке без изменений: meta,
// наименование и внешний код целиком, как в примерах источника.
func Existing(p model.PriceType) ReplaceItem {
	meta, code := p.Meta, p.ExternalCode
	item := ReplaceItem{Meta: &meta, Name: p.Name}
	if code != "" {
		item.ExternalCode = &code
	}
	return item
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: pricetype service requires a client")
	}
	return &Service{client: client}, nil
}

const endpoint = "context/companysettings/pricetype"

// List читает все типы цен: ответ — массив, а не коллекция, пагинации нет. Тип по умолчанию
// возвращает GetDefault.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/pricetype#3-poluchit-spisok-vseh-tipov-cen
func (s *Service) List(ctx context.Context, opts ...transport.RequestOption) ([]model.PriceType, error) {
	transport.ResetResponseInfo(opts)
	return s.list(ctx, "list price types", http.MethodGet, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/pricetype#3-poluchit-tip-ceny-po-id
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*model.PriceType, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get price type: id %q must be a UUID", id)
	}
	return s.one(ctx, "get price type "+id, endpoint+"/"+id, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/pricetype#3-poluchit-tip-ceny-po-umolchaniyu
func (s *Service) GetDefault(ctx context.Context, opts ...transport.RequestOption) (*model.PriceType, error) {
	transport.ResetResponseInfo(opts)
	return s.one(ctx, "get default price type", endpoint+"/default", opts)
}

// ReplaceAll сохраняет весь список типов цен в переданном порядке (1…MaxPriceTypes элементов —
// ограничение API):
// не вошедшие в items сервер удаляет, первый становится типом по умолчанию. SDK не читает
// список перед записью: изменения, сделанные другими клиентами после вашего List, будут
// перезаписаны. Не повторять автоматически после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/pricetype#3-redaktirovanie-spiska-tipov-cen
func (s *Service) ReplaceAll(ctx context.Context, items []ReplaceItem, opts ...transport.RequestOption) ([]model.PriceType, error) {
	transport.ResetResponseInfo(opts)
	const operation = "replace price types"
	if len(items) == 0 || len(items) > MaxPriceTypes {
		return nil, fmt.Errorf("moysklad: %s: list must contain 1..%d price types", operation, MaxPriceTypes)
	}
	for i, item := range items {
		if strings.TrimSpace(item.Name) == "" {
			return nil, fmt.Errorf("moysklad: %s: price type #%d must have a non-blank name", operation, i)
		}
		if item.Meta != nil && item.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: existing price type #%d has no meta.href", operation, i)
		}
	}
	return s.list(ctx, operation, http.MethodPost, items, opts)
}

func (s *Service) list(ctx context.Context, operation, method string, body any, opts []transport.RequestOption) ([]model.PriceType, error) {
	var result []model.PriceType
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	// Хотя бы один тип цены есть всегда: удалить все API не даёт.
	if len(result) == 0 {
		return nil, fmt.Errorf("moysklad: %s: response must be a non-empty array", operation)
	}
	return result, nil
}

func (s *Service) one(ctx context.Context, operation, path string, opts []transport.RequestOption) (*model.PriceType, error) {
	var result model.PriceType
	resp, err := s.client.Send(ctx, http.MethodGet, path, nil, nil, &result, opts...)
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
