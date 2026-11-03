package processingstage

import "github.com/eteeet/moysklad-sdk-go/model"

// ProcessingStage — этап производства. Исполнители: AllPerformers=true и пустой Performers — любой
// активный сотрудник; AllPerformers=false и пустой Performers — никто; AllPerformers=false и
// непустой Performers — только перечисленные сотрудники.
// StandardHourCost — стоимость нормо-часа, float64 в копейках, как остальные суммы API.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingstage#2-etap-proizvodstva
type ProcessingStage struct {
	Meta                 model.Meta        `json:"meta"`
	ID                   string            `json:"id"`
	AccountID            string            `json:"accountId"`
	Updated              model.DateTime    `json:"updated,omitzero"`
	Name                 string            `json:"name"`
	Description          string            `json:"description"`
	ExternalCode         string            `json:"externalCode"`
	Archived             bool              `json:"archived"`
	Shared               bool              `json:"shared"`
	AllPerformers        bool              `json:"allPerformers"`
	DistributionRequired bool              `json:"distributionRequired"`
	StandardHourCost     float64           `json:"standardHourCost"`
	Performers           []model.Reference `json:"performers"`
	MaterialStore        *model.Reference  `json:"materialStore,omitzero"`
	Owner                *model.Reference  `json:"owner,omitzero"`
	Group                *model.Reference  `json:"group,omitzero"`
}
