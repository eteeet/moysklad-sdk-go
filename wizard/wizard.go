// Package wizard — автозаполнение цен, скидок, НДС и себестоимости шаблона документа:
// POST wizard/{тип документа}?action=…. Документ не создаётся и не сохраняется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/autofill
package wizard

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Значения action; набор открыт. Поддержку действия типом документа (скидки — не у inventory, move,
// internalorder, enter, loss, payroll; цены — не у inventory; себестоимость — только у возвратов
// без основания) проверяет сервер.
const (
	EvaluatePrice    = "evaluate_price"
	EvaluateVAT      = "evaluate_vat"
	EvaluateDiscount = "evaluate_discount"
	EvaluateCost     = "evaluate_cost"
)

var documentType = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// Template — шаблон автозаполнения: поля, влияющие на расчёт. Agent нужен для цен и скидок (без
// него — Organization), Organization — для НДС, Store — для себестоимости; проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/autofill#3-shablon-avtozapolneniya
type Template struct {
	Organization model.Field[model.Reference]           `json:"organization,omitzero"`
	Agent        model.Field[model.Reference]           `json:"agent,omitzero"`
	Store        model.Field[model.Reference]           `json:"store,omitzero"`
	VATEnabled   model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded  model.Field[bool]                      `json:"vatIncluded,omitzero"`
	Rate         model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Moment       model.Field[model.DateTime]            `json:"moment,omitzero"`
	Positions    []PositionRequest                      `json:"positions"`
}

// PositionRequest — позиция шаблона: Price — копейки, Discount — процент скидки или наценки.
type PositionRequest struct {
	Assortment model.Reference      `json:"assortment"`
	Quantity   model.Field[float64] `json:"quantity,omitzero"`
	Price      model.Field[float64] `json:"price,omitzero"`
	Discount   model.Field[float64] `json:"discount,omitzero"`
	VAT        model.Field[int]     `json:"vat,omitzero"`
	VATEnabled model.Field[bool]    `json:"vatEnabled,omitzero"`
}

// Result — заполненный шаблон. Raw — ответ целиком: поля, не влияющие на расчёт, сервер
// возвращает как переданы.
type Result struct {
	Organization *model.Reference    `json:"organization,omitzero"`
	Agent        *model.Reference    `json:"agent,omitzero"`
	Store        *model.Reference    `json:"store,omitzero"`
	VATEnabled   bool                `json:"vatEnabled"`
	VATIncluded  bool                `json:"vatIncluded"`
	Rate         *model.DocumentRate `json:"rate,omitzero"`
	Moment       model.DateTime      `json:"moment,omitzero"`
	Positions    []Position          `json:"positions"`
	Raw          jsontext.Value      `json:"-"`
}

func (r *Result) UnmarshalJSON(data []byte) error {
	type plain Result
	var v plain
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = Result(v)
	r.Raw = bytes.Clone(data)
	return nil
}

// Position — рассчитанная позиция: Price, DiscountedPrice (с учётом скидки и НДС), Sum (при
// переданном quantity) и Cost — копейки; Discount — процент.
type Position struct {
	Assortment      model.Reference `json:"assortment"`
	Quantity        float64         `json:"quantity"`
	Price           float64         `json:"price"`
	Discount        float64         `json:"discount"`
	VAT             int             `json:"vat"`
	VATEnabled      bool            `json:"vatEnabled"`
	DiscountedPrice float64         `json:"discountedPrice"`
	Sum             float64         `json:"sum"`
	Cost            float64         `json:"cost"`
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: wizard service requires a client")
	}
	return &Service{client: client}, nil
}

// Fill — POST wizard/{document}?action=a,b: один запрос с одним или несколькими действиями.
// Документ не создаётся; POST не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/autofill#3-zapros-avtozapolneniya
func (s *Service) Fill(ctx context.Context, document string, actions []string, template Template, opts ...transport.RequestOption) (*Result, error) {
	transport.ResetResponseInfo(opts)
	operation := "fill " + document + " template"
	if !documentType.MatchString(document) {
		return nil, fmt.Errorf("moysklad: %s: document type must match [a-z][a-z0-9]*", operation)
	}
	if len(actions) == 0 {
		return nil, fmt.Errorf("moysklad: %s: at least one action is required", operation)
	}
	for i, action := range actions {
		if action == "" || strings.Contains(action, ",") {
			return nil, fmt.Errorf("moysklad: %s: actions[%d] must be a non-empty value without commas", operation, i)
		}
	}
	if template.Positions == nil {
		template.Positions = []PositionRequest{}
	}
	var result *Result
	resp, err := s.client.Send(ctx, http.MethodPost, "wizard/"+document, url.Values{"action": {strings.Join(actions, ",")}}, template, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result == nil {
		return nil, fmt.Errorf("moysklad: %s: response is not an object", operation)
	}
	return result, nil
}
