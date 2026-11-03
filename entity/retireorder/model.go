package retireorder

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// RetireOrder — вывод кодов маркировки из оборота. Модель чтения без expand: Sum и VATSum —
// копейки. TrackingType, RetireOrderType (способ вывода), SupportingTransaction и DocumentState —
// открытые enum: неизвестное значение сохраняется как есть. Строки — GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder
type RetireOrder struct {
	Meta                        model.Meta                  `json:"meta"`
	ID                          string                      `json:"id"`
	AccountID                   string                      `json:"accountId"`
	SyncID                      string                      `json:"syncId"`
	Name                        string                      `json:"name"`
	Description                 string                      `json:"description"`
	ExternalCode                string                      `json:"externalCode"`
	Created                     model.DateTime              `json:"created,omitzero"`
	Updated                     model.DateTime              `json:"updated,omitzero"`
	Deleted                     model.DateTime              `json:"deleted,omitzero"`
	Moment                      model.DateTime              `json:"moment,omitzero"`
	Shared                      bool                        `json:"shared"`
	Printed                     bool                        `json:"printed"`
	Published                   bool                        `json:"published"`
	Organization                *model.Reference            `json:"organization,omitzero"`
	OrganizationBranch          *model.Reference            `json:"organizationBranch,omitzero"`
	Agent                       *model.Reference            `json:"agent,omitzero"`
	DestinationCountry          *model.Reference            `json:"destinationCountry,omitzero"`
	Owner                       *model.Reference            `json:"owner,omitzero"`
	Group                       *model.Reference            `json:"group,omitzero"`
	State                       *model.Reference            `json:"state,omitzero"`
	Rate                        *model.DocumentRate         `json:"rate,omitzero"`
	Sum                         float64                     `json:"sum"`
	VATEnabled                  bool                        `json:"vatEnabled"`
	VATIncluded                 bool                        `json:"vatIncluded"`
	VATSum                      float64                     `json:"vatSum"`
	TrackingType                string                      `json:"trackingType"`
	RetireOrderType             string                      `json:"retireOrderType"`
	DocumentState               string                      `json:"documentState"`
	ReasonDescription           string                      `json:"reasonDescription"`
	PrimaryDocumentName         string                      `json:"primaryDocumentName"`
	StateContractID             string                      `json:"stateContractId"`
	SupportingTransaction       string                      `json:"supportingTransaction"`
	SupportingTransactionDate   model.DateTime              `json:"supportingTransactionDate,omitzero"`
	SupportingTransactionNumber string                      `json:"supportingTransactionNumber"`
	Positions                   *model.Collection[Position] `json:"positions,omitzero"`
}

// Position — строка вывода (meta.type retireorderposition). Price — копейки; VAT — ставка в
// процентах; TrackingCodes — коды маркировки (type открыт: trackingcode, consumerpack, …).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#4-pozicii-vyvoda-iz-oborota
type Position struct {
	positions.Item
	Price         float64                  `json:"price"`
	VAT           int                      `json:"vat"`
	VATEnabled    bool                     `json:"vatEnabled"`
	TrackingCodes []positions.TrackingCode `json:"trackingCodes"`
}

// PositionRequest — строка в Create/Update документа, CreatePositions и UpdatePosition. Для новой
// строки источник требует assortment и trackingCodes — проверяет сервер; в UpdatePosition
// обязательных полей нет. ID сопоставляет существующую строку в Update документа.
type PositionRequest struct {
	positions.ItemRequest
	Price         model.Field[float64]                  `json:"price,omitzero"`
	VAT           model.Field[int]                      `json:"vat,omitzero"`
	VATEnabled    model.Field[bool]                     `json:"vatEnabled,omitzero"`
	TrackingCodes model.Field[[]positions.TrackingCode] `json:"trackingCodes,omitzero"`
}
