package contract

import "github.com/eteeet/moysklad-sdk-go/model"

// Contract — договор. Name — номер договора; Sum — копейки; RewardPercent — целый процент (0–100), не деньги.
// ContractType (Commission, Sales) и RewardType (PercentOfSales, None) — открытые enum.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract
type Contract struct {
	Meta                model.Meta          `json:"meta"`
	ID                  string              `json:"id"`
	AccountID           string              `json:"accountId"`
	Updated             model.DateTime      `json:"updated,omitzero"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Code                string              `json:"code"`
	ExternalCode        string              `json:"externalCode"`
	Archived            bool                `json:"archived"`
	Moment              model.DateTime      `json:"moment,omitzero"`
	Sum                 float64             `json:"sum"`
	ContractType        string              `json:"contractType"`
	RewardType          string              `json:"rewardType"`
	RewardPercent       int                 `json:"rewardPercent"`
	OwnAgent            *model.Reference    `json:"ownAgent,omitzero"`
	Agent               *model.Reference    `json:"agent,omitzero"`
	OrganizationAccount *model.Reference    `json:"organizationAccount,omitzero"`
	AgentAccount        *model.Reference    `json:"agentAccount,omitzero"`
	Rate                *model.DocumentRate `json:"rate,omitzero"`
	State               *model.Reference    `json:"state,omitzero"`
	Attributes          []model.Attribute   `json:"attributes"`
	Printed             bool                `json:"printed"`
	Published           bool                `json:"published"`
	Shared              bool                `json:"shared"`
	Owner               *model.Reference    `json:"owner,omitzero"`
	Group               *model.Reference    `json:"group,omitzero"`
}
