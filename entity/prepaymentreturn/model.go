package prepaymentreturn

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки товарных документов; у возврата предоплаты — assortment, quantity, price (копейки),
// discount, vat, vatEnabled и pack (meta.type prepaymentreturnposition). Остальные поля общей формы
// остаются нулевыми.
type Position = positions.Position

// PrepaymentReturn — возврат предоплаты. Только чтение: Create/Update источник не описывает. Модель без expand:
// sum, cashSum, noCashSum, qrSum и vatSum — копейки; retailShift — ссылка на розничную смену.
// version и documents из примера ответа в таблице атрибутов нет — SDK их не читает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment-return
type PrepaymentReturn struct {
	Meta               model.Meta       `json:"meta"`
	ID                 string           `json:"id"`
	AccountID          string           `json:"accountId"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	Code               string           `json:"code"`
	ExternalCode       string           `json:"externalCode"`
	SyncID             string           `json:"syncId"`
	Created            model.DateTime   `json:"created,omitzero"`
	Updated            model.DateTime   `json:"updated,omitzero"`
	Deleted            model.DateTime   `json:"deleted,omitzero"`
	Moment             model.DateTime   `json:"moment,omitzero"`
	Applicable         bool             `json:"applicable"`
	Shared             bool             `json:"shared"`
	Printed            bool             `json:"printed"`
	Published          bool             `json:"published"`
	Organization       *model.Reference `json:"organization,omitzero"`
	OrganizationBranch *model.Reference `json:"organizationBranch,omitzero"`
	Agent              *model.Reference `json:"agent,omitzero"`
	Owner              *model.Reference `json:"owner,omitzero"`
	Group              *model.Reference `json:"group,omitzero"`
	State              *model.Reference `json:"state,omitzero"`
	RetailShift        *model.Reference `json:"retailShift,omitzero"`
	RetailStore        *model.Reference `json:"retailStore,omitzero"`
	// Prepayment — исходная предоплата (ссылка).
	Prepayment  *model.Reference    `json:"prepayment,omitzero"`
	Attributes  []model.Attribute   `json:"attributes"`
	Rate        *model.DocumentRate `json:"rate,omitzero"`
	Sum         float64             `json:"sum"`
	CashSum     float64             `json:"cashSum"`
	NoCashSum   float64             `json:"noCashSum"`
	QRSum       float64             `json:"qrSum"`
	VATSum      float64             `json:"vatSum"`
	VATEnabled  bool                `json:"vatEnabled"`
	VATIncluded bool                `json:"vatIncluded"`
	// TaxSystem открыт: GENERAL_TAX_SYSTEM, SIMPLIFIED_TAX_SYSTEM_INCOME, … PATENT_BASED.
	TaxSystem string                             `json:"taxSystem"`
	Positions *model.Collection[Position]        `json:"positions,omitzero"`
	Files     *model.Collection[model.Reference] `json:"files,omitzero"`
}
