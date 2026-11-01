package bonusprogram

import "github.com/eteeet/moysklad-sdk-go/model"

// BonusProgram — бонусная программа. Курсы (рублей за балл, баллов за рубль), проценты, дни и
// приветственные баллы — целые числа, не деньги. WelcomeBonusesMode открыт: REGISTRATION,
// FIRST_PURCHASE. Баланс и правила начисления SDK не рассчитывает. Поля без маркера «обязательное
// при ответе» (курсы, проценты, приветственные баллы) при отсутствии читаются как 0 или "".
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram
type BonusProgram struct {
	Meta                      model.Meta `json:"meta"`
	ID                        string     `json:"id"`
	AccountID                 string     `json:"accountId"`
	Name                      string     `json:"name"`
	Active                    bool       `json:"active"`
	AllProducts               bool       `json:"allProducts"`
	AllAgents                 bool       `json:"allAgents"`
	AgentTags                 []string   `json:"agentTags"`
	EarnRateRoublesToPoint    int        `json:"earnRateRoublesToPoint"`
	SpendRatePointsToRouble   int        `json:"spendRatePointsToRouble"`
	MaxPaidRatePercents       int        `json:"maxPaidRatePercents"`
	EarnWhileRedeeming        bool       `json:"earnWhileRedeeming"`
	WelcomeBonusesEnabled     bool       `json:"welcomeBonusesEnabled"`
	WelcomeBonusesMode        string     `json:"welcomeBonusesMode"`
	WelcomeBonusesValue       int        `json:"welcomeBonusesValue"`
	PostponedBonusesDelayDays int        `json:"postponedBonusesDelayDays"`
}
