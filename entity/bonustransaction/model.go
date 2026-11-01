package bonustransaction

import "github.com/eteeet/moysklad-sdk-go/model"

// BonusTransaction — бонусная операция: начисление (EARNING) или списание (SPENDING) баллов.
// BonusValue — баллы, целое число, не деньги. Перечисления открыты: TransactionStatus
// (WAIT_PROCESSING, COMPLETED, CANCELED) и CategoryType (REGULAR, WELCOME) только для чтения.
// UpdatedBy источник называет атрибутом только для фильтрации, но фикстура его возвращает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction
type BonusTransaction struct {
	Meta              model.Meta       `json:"meta"`
	ID                string           `json:"id"`
	AccountID         string           `json:"accountId"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Code              string           `json:"code"`
	ExternalCode      string           `json:"externalCode"`
	Created           model.DateTime   `json:"created,omitzero"`
	Updated           model.DateTime   `json:"updated,omitzero"`
	UpdatedBy         string           `json:"updatedBy"`
	Moment            model.DateTime   `json:"moment,omitzero"`
	ExecutionDate     model.DateTime   `json:"executionDate,omitzero"`
	Applicable        bool             `json:"applicable"`
	Shared            bool             `json:"shared"`
	Agent             *model.Reference `json:"agent,omitzero"`
	BonusProgram      *model.Reference `json:"bonusProgram,omitzero"`
	Organization      *model.Reference `json:"organization,omitzero"`
	Owner             *model.Reference `json:"owner,omitzero"`
	Group             *model.Reference `json:"group,omitzero"`
	ParentDocument    *model.Reference `json:"parentDocument,omitzero"`
	BonusValue        int              `json:"bonusValue"`
	TransactionType   string           `json:"transactionType"`
	TransactionStatus string           `json:"transactionStatus"`
	CategoryType      string           `json:"categoryType"`
}
