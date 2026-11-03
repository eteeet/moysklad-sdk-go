// Package trackingcode — коды маркировки отдельной позиции документа:
// entity/{type}/{id}/positions/{positionID}/trackingCodes. Отдельной сущности вне позиции нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/tracking-code
package trackingcode

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// MaxLimit — наибольший limit списка кодов (по умолчанию сервер отдаёт 100).
const MaxLimit = 100

// Форматы кодов в списке (codetype); по умолчанию сервер отдаёт gs1.
const (
	CodeTypeGS1  = "gs1"
	CodeType1162 = "tag_1162"
	CodeTypeAll  = "all"
)

var entityType = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// TrackingCode — код маркировки (meta.type trackingcode). Type открыт: trackingcode,
// consumerpack, transportpack …; вложенные коды — у упаковок. CIS1162 — формат тега 1162,
// только чтение.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/tracking-code#3-kod-markirovki
type TrackingCode struct {
	Meta          model.Meta     `json:"meta"`
	ID            string         `json:"id"`
	CIS           string         `json:"cis"`
	CIS1162       string         `json:"cis_1162"`
	Type          string         `json:"type"`
	TrackingCodes []TrackingCode `json:"trackingCodes"`
}

// Request — код в массовом запросе: без ID создаётся, с ID изменяется. ID только адресует код
// (метод массового изменения требует его), значение id не меняется — readonly-маркер таблицы
// этому не противоречит. cis и type необходимы только при создании (проверяется до запроса и у
// вложенных кодов); вложенные коды заменяют прежние целиком.
type Request struct {
	ID            string              `json:"id,omitempty"`
	CIS           model.Field[string] `json:"cis,omitzero"`
	Type          model.Field[string] `json:"type,omitzero"`
	TrackingCodes []Request           `json:"trackingCodes,omitzero"`
}

// Options — страница списка: Limit 0 — значение сервера (100), 1–100; CodeType — формат.
type Options struct {
	Limit    int
	Offset   int
	CodeType string
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: trackingcode service requires a client")
	}
	return &Service{client: client}, nil
}

// List читает одну страницу кодов позиции. Порядок кодов первого уровня фиксирован, вложенных — нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/tracking-code#3-poluchit-kody-markirovki-pozicii-dokumenta
func (s *Service) List(ctx context.Context, typ, id, positionID string, options Options, opts ...transport.RequestOption) (*model.Collection[TrackingCode], error) {
	transport.ResetResponseInfo(opts)
	operation := "list " + typ + " " + id + " position " + positionID + " trackingCodes"
	if err := owner(typ, id, positionID); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	if options.Limit < 0 || options.Limit > MaxLimit || options.Offset < 0 {
		return nil, fmt.Errorf("moysklad: %s: limit must be 0..%d and offset must not be negative", operation, MaxLimit)
	}
	params := url.Values{}
	if options.Limit > 0 {
		params.Set("limit", strconv.Itoa(options.Limit))
	}
	if options.Offset > 0 {
		params.Set("offset", strconv.Itoa(options.Offset))
	}
	if options.CodeType != "" {
		params.Set("codetype", options.CodeType)
	}
	return list.Page[TrackingCode](ctx, s.client, operation, endpoint(typ, id, positionID), params, opts...)
}

// Save — «Массовое создание и обновление Кодов маркировки»: один POST массивом, элементы с ID
// изменяются, без ID — создаются. Ответ — по элементу на каждый (правила batch); повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/tracking-code#3-massovoe-sozdanie-i-obnovlenie-kodov-markirovki
func (s *Service) Save(ctx context.Context, typ, id, positionID string, items []Request, opts ...transport.RequestOption) ([]batch.Result[TrackingCode], error) {
	transport.ResetResponseInfo(opts)
	operation := "save " + typ + " " + id + " position " + positionID + " trackingCodes"
	if err := owner(typ, id, positionID); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	for i, item := range items {
		if err := check(item); err != nil {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: %w", operation, i, err)
		}
	}
	return write.CreateManyNested[TrackingCode](ctx, s.client, operation, endpoint(typ, id, positionID), items, opts)
}

// DeleteMany — POST …/trackingCodes/delete с meta кодов; вложенные коды удаляются вместе с ними.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/tracking-code#3-massovoe-udalenie-kodov-markirovki
func (s *Service) DeleteMany(ctx context.Context, typ, id, positionID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	operation := "delete many " + typ + " " + id + " position " + positionID + " trackingCodes"
	if err := owner(typ, id, positionID); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return write.DeleteManyNested(ctx, s.client, operation, endpoint(typ, id, positionID), "trackingcode", ids, opts)
}

// check проверяет код и вложенные коды: UUID у изменяемых, непустые cis и type у новых.
func check(item Request) error {
	if item.ID != "" && !validID(item.ID) {
		return fmt.Errorf("id must be a UUID")
	}
	cis, _ := item.CIS.Get()
	typ, _ := item.Type.Get()
	if item.ID == "" && (cis == "" || typ == "") {
		return fmt.Errorf("cis and type are required for a new code")
	}
	for i, nested := range item.TrackingCodes {
		if err := check(nested); err != nil {
			return fmt.Errorf("trackingCodes[%d]: %w", i, err)
		}
	}
	return nil
}

func endpoint(typ, id, positionID string) string {
	return "entity/" + typ + "/" + id + "/positions/" + positionID + "/trackingCodes"
}

func owner(typ, id, positionID string) error {
	if !entityType.MatchString(typ) {
		return fmt.Errorf("entity type must match [a-z][a-z0-9]*")
	}
	if !validID(id) || !validID(positionID) {
		return fmt.Errorf("document and position ids must be UUIDs")
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
