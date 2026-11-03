// Package discount — скидки: общий список всех видов (entity/discount) и округление копеек.
// Накопительная, персональная скидка и спеццена — пакеты accumulationdiscount, personaldiscount
// и specialpricediscount.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#2-skidki
package discount

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Discount — строка общего списка скидок. Вид — Meta.Type: discount (округление копеек),
// accumulationdiscount, personaldiscount, specialpricediscount; неизвестный вид не отвергается,
// его поля вне общего набора — в Raw (исходный JSON строки). Поля видов приходят только у своих
// скидок. Таблица называет коллекцию групп productfolders, примеры и фикстуры — productFolders:
// читаются оба ключа.
type Discount struct {
	Meta           model.Meta        `json:"meta"`
	ID             string            `json:"id"`
	AccountID      string            `json:"accountId"`
	Name           string            `json:"name"`
	Active         bool              `json:"active"`
	AllProducts    bool              `json:"allProducts"`
	AllAgents      bool              `json:"allAgents"`
	AgentTags      []string          `json:"agentTags"`
	Assortment     []model.Reference `json:"assortment"`
	ProductFolders []model.Reference `json:"productFolders"`
	Levels         []Level           `json:"levels"`
	Discount       float64           `json:"discount"`
	UsePriceType   bool              `json:"usePriceType"`
	SpecialPrice   *SpecialPrice     `json:"specialPrice,omitzero"`
	Raw            jsontext.Value    `json:"-"`
}

func (d *Discount) UnmarshalJSON(data []byte) error {
	type plain Discount
	var v struct {
		plain
		Lower []model.Reference `json:"productfolders"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*d = Discount(v.plain)
	if d.ProductFolders == nil {
		d.ProductFolders = v.Lower
	}
	d.Raw = bytes.Clone(data)
	return nil
}

// Level — уровень накопительной скидки: Amount — сумма накоплений в копейках, Discount — процент.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#4-levels
type Level struct {
	Amount   float64 `json:"amount"`
	Discount float64 `json:"discount"`
}

// SpecialPrice — спеццена: тип цены и значение в копейках (если выбрано фиксированное значение).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#4-specialprice
type SpecialPrice struct {
	PriceType *model.PriceType `json:"priceType,omitzero"`
	Value     float64          `json:"value"`
}

// SpecialPriceRequest — спеццена в запросе по таблице: priceType и value (копейки). Пример создания
// источника передаёт вместо этого meta типа цены — SDK следует таблице; сервером эта форма записи
// не проверена.
type SpecialPriceRequest struct {
	PriceType model.Field[model.Reference] `json:"priceType,omitzero"`
	Value     model.Field[float64]         `json:"value,omitzero"`
}

// RoundOffRequest — изменяемые поля округления копеек: name и active; agentTags и allAgents
// приходят в ответе, но не меняются.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-izmenit-okruglenie-kopeek
type RoundOffRequest struct {
	Name   model.Field[string] `json:"name,omitzero"`
	Active model.Field[bool]   `json:"active,omitzero"`
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: discount service requires a client")
	}
	return &Service{client: client}, nil
}

// List читает одну страницу всех скидок без expand и автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-poluchit-vse-skidki
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Discount], error) {
	return list.Entity[Discount](ctx, s.client, "discount", options, nil, opts...)
}

// UpdateRoundOff — PUT entity/discount/{id}: изменение округления копеек (создания и удаления
// источник не описывает).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-izmenit-okruglenie-kopeek
func (s *Service) UpdateRoundOff(ctx context.Context, id string, input RoundOffRequest, opts ...transport.RequestOption) (*Discount, error) {
	transport.ResetResponseInfo(opts)
	operation := "update discount " + id
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update discount: id must be a UUID")
	}
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	var result Discount
	resp, err := s.client.Send(ctx, http.MethodPut, "entity/discount/"+id, nil, input, &result, opts...)
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

// UpdateMany — «Массовое обновление округления копеек»: один POST entity/discount, только
// существующие объекты (meta по ID). Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-massovoe-obnovlenie-okrugleniya-kopeek
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[RoundOffRequest], opts ...transport.RequestOption) ([]batch.Result[Discount], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Discount](ctx, s.client, "discount", items, opts)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
