package bonusprogram

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязательные поля метода — обычные значения, поэтому false и 0 всегда
// отправляются. allProducts источник описывает как «всегда true». Диапазоны курсов и процента,
// приветственные баллы (не пустые при welcomeBonusesEnabled, не отрицательные) и тарифную опцию
// postponedBonusesDelayDays проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram#3-sozdat-bonusnuyu-programmu
type CreateRequest struct {
	Name                    string `json:"name"`
	Active                  bool   `json:"active"`
	AllProducts             bool   `json:"allProducts"`
	AllAgents               bool   `json:"allAgents"`
	EarnRateRoublesToPoint  int    `json:"earnRateRoublesToPoint"`
	SpendRatePointsToRouble int    `json:"spendRatePointsToRouble"`
	MaxPaidRatePercents     int    `json:"maxPaidRatePercents"`

	AgentTags                 model.Field[[]string] `json:"agentTags,omitzero"`
	EarnWhileRedeeming        model.Field[bool]     `json:"earnWhileRedeeming,omitzero"`
	WelcomeBonusesEnabled     model.Field[bool]     `json:"welcomeBonusesEnabled,omitzero"`
	WelcomeBonusesMode        model.Field[string]   `json:"welcomeBonusesMode,omitzero"`
	WelcomeBonusesValue       model.Field[int]      `json:"welcomeBonusesValue,omitzero"`
	PostponedBonusesDelayDays model.Field[int]      `json:"postponedBonusesDelayDays,omitzero"`
}

// UpdateRequest передаёт только изменяемые поля.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonusprogram#3-izmenit-bonusnuyu-programmu
type UpdateRequest struct {
	Name                      model.Field[string]   `json:"name,omitzero"`
	Active                    model.Field[bool]     `json:"active,omitzero"`
	AllProducts               model.Field[bool]     `json:"allProducts,omitzero"`
	AllAgents                 model.Field[bool]     `json:"allAgents,omitzero"`
	AgentTags                 model.Field[[]string] `json:"agentTags,omitzero"`
	EarnRateRoublesToPoint    model.Field[int]      `json:"earnRateRoublesToPoint,omitzero"`
	SpendRatePointsToRouble   model.Field[int]      `json:"spendRatePointsToRouble,omitzero"`
	MaxPaidRatePercents       model.Field[int]      `json:"maxPaidRatePercents,omitzero"`
	EarnWhileRedeeming        model.Field[bool]     `json:"earnWhileRedeeming,omitzero"`
	WelcomeBonusesEnabled     model.Field[bool]     `json:"welcomeBonusesEnabled,omitzero"`
	WelcomeBonusesMode        model.Field[string]   `json:"welcomeBonusesMode,omitzero"`
	WelcomeBonusesValue       model.Field[int]      `json:"welcomeBonusesValue,omitzero"`
	PostponedBonusesDelayDays model.Field[int]      `json:"postponedBonusesDelayDays,omitzero"`
}
