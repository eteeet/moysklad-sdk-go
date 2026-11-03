package assortment

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/single"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Settings — настройки справочника товаров (meta.type assortmentsettings). Поле общего доступа
// называется createdShared — в отличие от createShared справочника контрагентов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#3-nastrojki-spravochnika
type Settings struct {
	Meta            model.Meta            `json:"meta"`
	UniqueCodeRules model.UniqueCodeRules `json:"uniqueCodeRules"`
	BarcodeRules    BarcodeRules          `json:"barcodeRules"`
	CreatedShared   bool                  `json:"createdShared"`
}

// BarcodeRules — правила штрихкодов: автоматический EAN13 и префикс весовых товаров
// (WeightBarcodePrefix — число формата X или XX, допустимость проверяет сервер).
type BarcodeRules struct {
	FillEAN13Barcode    bool `json:"fillEAN13Barcode"`
	WeightBarcode       bool `json:"weightBarcode"`
	WeightBarcodePrefix int  `json:"weightBarcodePrefix"`
}

type BarcodeRulesRequest struct {
	FillEAN13Barcode    model.Field[bool] `json:"fillEAN13Barcode,omitzero"`
	WeightBarcode       model.Field[bool] `json:"weightBarcode,omitzero"`
	WeightBarcodePrefix model.Field[int]  `json:"weightBarcodePrefix,omitzero"`
}

// SettingsRequest — частичное изменение. Пример источника пишет checkUniqueCodeBoolean;
// SDK следует таблице атрибутов: checkUniqueCode.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#3-izmenit-nastrojki-spravochnika-tovarov
type SettingsRequest struct {
	UniqueCodeRules model.Field[model.UniqueCodeRulesRequest] `json:"uniqueCodeRules,omitzero"`
	BarcodeRules    model.Field[BarcodeRulesRequest]          `json:"barcodeRules,omitzero"`
	CreatedShared   model.Field[bool]                         `json:"createdShared,omitzero"`
}

// GetSettings — GET entity/assortment/settings.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#3-poluchit-nastrojki-spravochnika-tovarov
func (s *Service) GetSettings(ctx context.Context, opts ...transport.RequestOption) (*Settings, error) {
	return single.Get[Settings](ctx, s.client, "get assortment settings", "entity/assortment/settings", opts)
}

// UpdateSettings — PUT entity/assortment/settings; ответ — настройки целиком.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#3-izmenit-nastrojki-spravochnika-tovarov
func (s *Service) UpdateSettings(ctx context.Context, input SettingsRequest, opts ...transport.RequestOption) (*Settings, error) {
	return single.Put[Settings](ctx, s.client, "update assortment settings", "entity/assortment/settings", input, opts)
}
