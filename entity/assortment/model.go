package assortment

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/eteeet/moysklad-sdk-go/entity/bundle"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/entity/service"
	"github.com/eteeet/moysklad-sdk-go/entity/variant"
	"github.com/eteeet/moysklad-sdk-go/model"
)

// Row — строка ассортимента. По meta.type заполнено одно из Product, Service, Bundle, Variant,
// Consignment; при неизвестном типе — только Type и Raw. Raw — вся строка как пришла.
// Row — модель чтения: запись возвращает Raw байт в байт, изменения типизированных полей
// не записываются; Row без Raw пишется как null.
type Row struct {
	Type model.MetaType

	Product     *product.Product
	Service     *service.Service
	Bundle      *bundle.Bundle
	Variant     *variant.Variant
	Consignment *Consignment

	Balances

	Raw jsontext.Value `json:"-"`
}

// Balances — остатки строки. Приходят только из ListLegacy (для услуг и комплектов
// не выводятся); List без остатков их не возвращает.
type Balances struct {
	Stock     *float64 `json:"stock"`
	Reserve   *float64 `json:"reserve"`
	InTransit *float64 `json:"inTransit"`
	Quantity  *float64 `json:"quantity"`
}

func (r *Row) UnmarshalJSON(data []byte) error {
	var head struct {
		Meta model.Meta `json:"meta"`
		Balances
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return err
	}
	*r = Row{Type: head.Meta.Type, Balances: head.Balances, Raw: jsontext.Value(append([]byte(nil), data...))}
	var dst any
	switch r.Type {
	case "product":
		r.Product = new(product.Product)
		dst = r.Product
	case "service":
		r.Service = new(service.Service)
		dst = r.Service
	case "bundle":
		r.Bundle = new(bundle.Bundle)
		dst = r.Bundle
	case "variant":
		r.Variant = new(variant.Variant)
		dst = r.Variant
	case "consignment":
		r.Consignment = new(Consignment)
		dst = r.Consignment
	default:
		return nil
	}
	return json.Unmarshal(data, dst)
}

func (r Row) MarshalJSON() ([]byte, error) {
	if len(r.Raw) == 0 {
		return []byte("null"), nil
	}
	return r.Raw, nil
}

// Consignment — партия в ассортименте (groupBy=consignment), общая модель чтения без expand.
// Assortment — ссылка на товар, услугу или модификацию партии.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment
type Consignment struct {
	Meta         model.Meta        `json:"meta"`
	ID           string            `json:"id"`
	AccountID    string            `json:"accountId"`
	Updated      model.DateTime    `json:"updated,omitzero"`
	Name         string            `json:"name"`
	Code         string            `json:"code"`
	ExternalCode string            `json:"externalCode"`
	Description  string            `json:"description"`
	Label        string            `json:"label"`
	Archived     bool              `json:"archived"`
	ExpiryDate   model.DateTime    `json:"expiryDate,omitzero"`
	Assortment   model.Reference   `json:"assortment"`
	Barcodes     []model.Barcode   `json:"barcodes"`
	Attributes   []model.Attribute `json:"attributes"`
	// Images — только метаданные коллекции изображений товара партии.
	Images *model.Collection[jsontext.Value] `json:"images,omitzero"`
}
