package contexts

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/internal/single"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// CompanySettings — настройки компании. Currency, PriceTypes и AccountCountry только читаются:
// типы цен изменяет pricetype.Service.ReplaceAll. DiscountStrategy — открытый enum (bySum,
// byPriority).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/companysettings#2-nastrojki-kompanii
type CompanySettings struct {
	Meta                     model.Meta        `json:"meta"`
	Currency                 *model.Reference  `json:"currency,omitzero"`
	PriceTypes               []model.PriceType `json:"priceTypes"`
	DiscountStrategy         string            `json:"discountStrategy"`
	GlobalOperationNumbering bool              `json:"globalOperationNumbering"`
	CheckShippingStock       bool              `json:"checkShippingStock"`
	CheckMinPrice            bool              `json:"checkMinPrice"`
	UseRecycleBin            bool              `json:"useRecycleBin"`
	UseCompanyAddress        bool              `json:"useCompanyAddress"`
	CompanyAddress           string            `json:"companyAddress"`
	AccountCountry           string            `json:"accountCountry"`
}

// CompanySettingsRequest — поля из списка редактируемых настроек источника; изменение
// частичное.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/companysettings#3-izmenit-nastrojki-kompanii
type CompanySettingsRequest struct {
	DiscountStrategy         model.Field[string] `json:"discountStrategy,omitzero"`
	GlobalOperationNumbering model.Field[bool]   `json:"globalOperationNumbering,omitzero"`
	CheckShippingStock       model.Field[bool]   `json:"checkShippingStock,omitzero"`
	CheckMinPrice            model.Field[bool]   `json:"checkMinPrice,omitzero"`
	UseRecycleBin            model.Field[bool]   `json:"useRecycleBin,omitzero"`
	UseCompanyAddress        model.Field[bool]   `json:"useCompanyAddress,omitzero"`
	CompanyAddress           model.Field[string] `json:"companyAddress,omitzero"`
}

// CompanySettings — GET context/companysettings.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/companysettings#3-poluchit-nastrojki-kompanii
func (s *Service) CompanySettings(ctx context.Context, opts ...transport.RequestOption) (*CompanySettings, error) {
	return single.Get[CompanySettings](ctx, s.client, "get companysettings", "context/companysettings", opts)
}

// UpdateCompanySettings — PUT context/companysettings; ответ — настройки целиком.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/companysettings#3-izmenit-nastrojki-kompanii
func (s *Service) UpdateCompanySettings(ctx context.Context, input CompanySettingsRequest, opts ...transport.RequestOption) (*CompanySettings, error) {
	return single.Put[CompanySettings](ctx, s.client, "update companysettings", "context/companysettings", input, opts)
}

// UserSettings — настройки пользователя. Ссылки Default* — значения по умолчанию в документах;
// без href — «не задано». Locale (ru_RU, en_US), DefaultScreen и PrintFormat —
// открытые enum. MailFooter таблица помечает Boolean, но пример и фикстура источника передают
// текст подписи: SDK следует данным — string.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/usersettings#2-nastrojki-polzovatelya
type UserSettings struct {
	Meta                        model.Meta   `json:"meta"`
	DefaultCompany              *SettingsRef `json:"defaultCompany,omitzero"`
	DefaultCustomerCounterparty *SettingsRef `json:"defaultCustomerCounterparty,omitzero"`
	DefaultPurchaseCounterparty *SettingsRef `json:"defaultPurchaseCounterparty,omitzero"`
	DefaultProject              *SettingsRef `json:"defaultProject,omitzero"`
	DefaultPlace                *SettingsRef `json:"defaultPlace,omitzero"`
	Locale                      string       `json:"locale"`
	MailFooter                  string       `json:"mailFooter"`
	FieldsPerRow                int          `json:"fieldsPerRow"`
	DefaultScreen               string       `json:"defaultScreen"`
	PrintFormat                 string       `json:"printFormat"`
	AutoShowReports             bool         `json:"autoShowReports"`
}

// UserSettingsRequest: Default* — Meta без обёртки {"meta": …}, как в примере изменения источника.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/usersettings#3-izmenit-nastrojki-polzovatelya
type UserSettingsRequest struct {
	DefaultCompany              model.Field[model.Meta] `json:"defaultCompany,omitzero"`
	DefaultCustomerCounterparty model.Field[model.Meta] `json:"defaultCustomerCounterparty,omitzero"`
	DefaultPurchaseCounterparty model.Field[model.Meta] `json:"defaultPurchaseCounterparty,omitzero"`
	DefaultProject              model.Field[model.Meta] `json:"defaultProject,omitzero"`
	DefaultPlace                model.Field[model.Meta] `json:"defaultPlace,omitzero"`
	Locale                      model.Field[string]     `json:"locale,omitzero"`
	MailFooter                  model.Field[string]     `json:"mailFooter,omitzero"`
	FieldsPerRow                model.Field[int]        `json:"fieldsPerRow,omitzero"`
	DefaultScreen               model.Field[string]     `json:"defaultScreen,omitzero"`
	PrintFormat                 model.Field[string]     `json:"printFormat,omitzero"`
	AutoShowReports             model.Field[bool]       `json:"autoShowReports,omitzero"`
}

// SettingsRef — ссылка настроек пользователя. Примеры источника передают Meta без обёртки,
// фикстура OpenAPI — {"meta": …}: чтение принимает обе формы, запись — форма источника.
type SettingsRef struct {
	model.Meta
}

func (r *SettingsRef) UnmarshalJSON(data []byte) error {
	var wrapped struct {
		Meta *model.Meta `json:"meta"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}
	if wrapped.Meta != nil {
		r.Meta = *wrapped.Meta
		return nil
	}
	return json.Unmarshal(data, &r.Meta)
}

// UserSettings — GET context/usersettings (настройки пользователя токена).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/usersettings#3-poluchit-nastrojki-polzovatelya
func (s *Service) UserSettings(ctx context.Context, opts ...transport.RequestOption) (*UserSettings, error) {
	return single.Get[UserSettings](ctx, s.client, "get usersettings", "context/usersettings", opts)
}

// UpdateUserSettings — PUT context/usersettings; ответ — настройки целиком.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/usersettings#3-izmenit-nastrojki-polzovatelya
func (s *Service) UpdateUserSettings(ctx context.Context, input UserSettingsRequest, opts ...transport.RequestOption) (*UserSettings, error) {
	return single.Put[UserSettings](ctx, s.client, "update usersettings", "context/usersettings", input, opts)
}

// Subscription — подписка компании; только чтение, без meta. Role (USER/ADMIN) и Tariff —
// открытые строки. SubscriptionEndDate — Unix-время в миллисекундах; у бессрочных тарифов
// источник значение не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/companysubscription#2-podpiska-kompanii
type Subscription struct {
	Role                          string `json:"role"`
	Tariff                        string `json:"tariff"`
	IsSubscriptionChangeAvailable bool   `json:"isSubscriptionChangeAvailable"`
	SubscriptionEndDate           int64  `json:"subscriptionEndDate"`
}

// Subscription — GET accountSettings/subscription. Ответ без meta и id: ошибка только если
// тело не JSON-объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/companysubscription#3-poluchit-podpisku-kompanii
func (s *Service) Subscription(ctx context.Context, opts ...transport.RequestOption) (*Subscription, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get subscription"
	var result *Subscription
	resp, err := s.client.Send(ctx, http.MethodGet, "accountSettings/subscription", nil, nil, &result, opts...)
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
