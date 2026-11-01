package positions

import "github.com/eteeet/moysklad-sdk-go/model"

type Core struct {
	Meta      model.Meta `json:"meta"`
	ID        string     `json:"id"`
	AccountID string     `json:"accountId"`
}

type Item struct {
	Core
	Assortment model.Reference `json:"assortment"`
	Quantity   float64         `json:"quantity"`
}

// Price, Overhead и Cost — в копейках; Overhead доступен только для чтения.
// Общая форма строк товарных документов: поле, которого нет у документа, остаётся нулевым.
// Reserve и TaxSystem — у заказа покупателя; Reason — у оприходования и списания; Slot, Pack,
// Country, GTD, Things — по таблице документа; Declaration (fields=declaration), Stock
// (fields=stock), TrackingCodes1162 и Shipped — только для чтения. Things: таблица enter
// называет тип Object(String), остальные документы и примеры — Array(String); SDK читает массив.
type Position struct {
	Item
	Price      float64          `json:"price"`
	Discount   float64          `json:"discount"`
	VAT        int              `json:"vat"`
	VATEnabled bool             `json:"vatEnabled"`
	Cost       float64          `json:"cost"`
	Shipped    float64          `json:"shipped"`
	InTransit  float64          `json:"inTransit"`
	Overhead   float64          `json:"overhead"`
	SourceSlot *model.Reference `json:"sourceSlot,omitzero"`
	TargetSlot *model.Reference `json:"targetSlot,omitzero"`

	Reserve           float64            `json:"reserve,omitzero"`
	TaxSystem         string             `json:"taxSystem,omitzero"`
	Reason            string             `json:"reason,omitzero"`
	Slot              *model.Reference   `json:"slot,omitzero"`
	Pack              *Pack              `json:"pack,omitzero"`
	Country           *model.Reference   `json:"country,omitzero"`
	GTD               *GTD               `json:"gtd,omitzero"`
	Things            []string           `json:"things,omitzero"`
	TrackingCodes     []TrackingCode     `json:"trackingCodes,omitzero"`
	TrackingCodes1162 []TrackingCode1162 `json:"trackingCodes_1162,omitzero"`
	Declaration       []Declaration      `json:"declaration,omitzero"`
	Stock             *Stock             `json:"stock,omitzero"`
}

// Pack — упаковка товара позиции. Таблицы строк ссылаются на «Упаковки Товара» справочника
// product (id, quantity, uom, barcodes); примеры запросов передают только id.
type Pack struct {
	ID       string           `json:"id"`
	Quantity float64          `json:"quantity,omitzero"`
	UOM      *model.Reference `json:"uom,omitzero"`
	Barcodes []model.Barcode  `json:"barcodes,omitzero"`
}

// GTD — грузовая таможенная декларация строки: объект с номером name.
type GTD struct {
	Name string `json:"name"`
}

// TrackingCode — код маркировки; вложенные коды допустимы у consumerpack и transportpack.
// Type открыт: trackingcode, consumerpack, transportpack.
type TrackingCode struct {
	CIS           string         `json:"cis"`
	Type          string         `json:"type"`
	TrackingCodes []TrackingCode `json:"trackingCodes,omitzero"`
}

// TrackingCode1162 — код маркировки в формате тега 1162, только для чтения.
type TrackingCode1162 struct {
	CIS1162           string             `json:"cis_1162"`
	Type              string             `json:"type"`
	TrackingCodes1162 []TrackingCode1162 `json:"trackingCodes_1162,omitzero"`
}

// Declaration — прослеживаемость импортного товара (fields=declaration), только для чтения.
type Declaration struct {
	GTD      string           `json:"gtd"`
	RNPT     string           `json:"rnpt"`
	Country  *model.Reference `json:"country,omitzero"`
	Quantity float64          `json:"quantity"`
}

// Stock — остатки и себестоимость позиции (fields=stock), только для чтения. Cost — копейки;
// у возвратов без основания сервер его не возвращает (nil).
type Stock struct {
	Cost      *float64 `json:"cost,omitzero"`
	Quantity  float64  `json:"quantity"`
	Reserve   float64  `json:"reserve"`
	InTransit float64  `json:"intransit"`
	Available float64  `json:"available"`
}
