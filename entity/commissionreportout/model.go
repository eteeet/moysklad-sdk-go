package commissionreportout

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// CommissionReportOut — выданный отчёт комиссионера. Модель чтения без expand: суммы — копейки
// и только для чтения; RewardPercent — целый процент, не деньги. Строки — GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout
type CommissionReportOut struct {
	Meta                  model.Meta                         `json:"meta"`
	ID                    string                             `json:"id"`
	AccountID             string                             `json:"accountId"`
	Name                  string                             `json:"name"`
	Description           string                             `json:"description"`
	Code                  string                             `json:"code"`
	ExternalCode          string                             `json:"externalCode"`
	SyncID                string                             `json:"syncId"`
	Created               model.DateTime                     `json:"created,omitzero"`
	Updated               model.DateTime                     `json:"updated,omitzero"`
	Deleted               model.DateTime                     `json:"deleted,omitzero"`
	Moment                model.DateTime                     `json:"moment,omitzero"`
	Applicable            bool                               `json:"applicable"`
	Shared                bool                               `json:"shared"`
	Printed               bool                               `json:"printed"`
	Published             bool                               `json:"published"`
	Organization          *model.Reference                   `json:"organization,omitzero"`
	OrganizationAccount   *model.Reference                   `json:"organizationAccount,omitzero"`
	Agent                 *model.Reference                   `json:"agent,omitzero"`
	AgentAccount          *model.Reference                   `json:"agentAccount,omitzero"`
	Contract              *model.Reference                   `json:"contract,omitzero"`
	Owner                 *model.Reference                   `json:"owner,omitzero"`
	Group                 *model.Reference                   `json:"group,omitzero"`
	Project               *model.Reference                   `json:"project,omitzero"`
	SalesChannel          *model.Reference                   `json:"salesChannel,omitzero"`
	State                 *model.Reference                   `json:"state,omitzero"`
	Attributes            []model.Attribute                  `json:"attributes"`
	Rate                  *model.DocumentRate                `json:"rate,omitzero"`
	VATEnabled            bool                               `json:"vatEnabled"`
	VATIncluded           bool                               `json:"vatIncluded"`
	Sum                   float64                            `json:"sum"`
	VATSum                float64                            `json:"vatSum"`
	PayedSum              float64                            `json:"payedSum"`
	CommitentSum          float64                            `json:"commitentSum"`
	CommissionPeriodStart model.DateTime                     `json:"commissionPeriodStart,omitzero"`
	CommissionPeriodEnd   model.DateTime                     `json:"commissionPeriodEnd,omitzero"`
	RewardType            string                             `json:"rewardType"` // PercentOfSales, None — открытый enum.
	RewardPercent         int                                `json:"rewardPercent"`
	Positions             *model.Collection[Position]        `json:"positions,omitzero"`
	Files                 *model.Collection[model.Reference] `json:"files,omitzero"`
	Payments              []model.Reference                  `json:"payments"`
}

// Position — строка выданного отчёта (meta.type commissionreportoutposition). Это собственная
// форма, не positions.Position: есть reward (вознаграждение). Price и Reward — копейки.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#4-pozicii-vydannogo-otcheta-komissionera
type Position struct {
	Meta       model.Meta       `json:"meta"`
	ID         string           `json:"id"`
	AccountID  string           `json:"accountId"`
	Assortment *model.Reference `json:"assortment,omitzero"`
	Quantity   float64          `json:"quantity"`
	Price      float64          `json:"price"`
	VAT        int              `json:"vat"`
	VATEnabled bool             `json:"vatEnabled"`
	Reward     float64          `json:"reward"`
	Pack       *positions.Pack  `json:"pack,omitzero"`
}
