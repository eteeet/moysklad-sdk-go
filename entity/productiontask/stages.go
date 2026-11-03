package productiontask

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Stage — производственный этап (meta.type productionstage). Этапы создаёт сервер по позициям
// задания; создать или удалить этап нельзя. ProcessingUnitCost, LabourUnitCost и StandardHourCost —
// денежные float64 без пересчёта (источник единицу не называет, у денежных полей API — копейки);
// StandardHourUnit (нормо-часы) и *Quantity — не деньги. Поля статуса у этапа источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-proizvodstvennye-etapy
type Stage struct {
	Meta                 model.Meta                         `json:"meta"`
	ID                   string                             `json:"id"`
	AccountID            string                             `json:"accountId"`
	Stage                *model.Reference                   `json:"stage,omitzero"`
	ProductionRow        *model.Reference                   `json:"productionRow,omitzero"`
	MaterialStore        *model.Reference                   `json:"materialStore,omitzero"`
	Materials            *model.Collection[Material]        `json:"materials,omitzero"`
	Files                *model.Collection[model.Reference] `json:"files,omitzero"`
	Instruction          string                             `json:"instruction"`
	OrderingPosition     int                                `json:"orderingPosition"`
	TotalQuantity        float64                            `json:"totalQuantity"`
	CompletedQuantity    float64                            `json:"completedQuantity"`
	AvailableQuantity    float64                            `json:"availableQuantity"`
	BlockedQuantity      float64                            `json:"blockedQuantity"`
	SkippedQuantity      float64                            `json:"skippedQuantity"`
	ProcessingUnitCost   float64                            `json:"processingUnitCost"`
	LabourUnitCost       float64                            `json:"labourUnitCost"`
	StandardHourCost     float64                            `json:"standardHourCost"`
	StandardHourUnit     float64                            `json:"standardHourUnit"`
	EnableHourAccounting bool                               `json:"enableHourAccounting"`
	PlannedEndDate       model.DateTime                     `json:"plannedEndDate,omitzero"`
}

// StageUpdateRequest — поля этапа без маркера «Только для чтения». EnableHourAccounting=true
// включает расчёт labourUnitCost = standardHourCost * standardHourUnit; при включённом флаге сервер
// не принимает labourUnitCost — для фиксированного значения передайте Value(false). Materials
// заменяет все материалы этапа: материал с ID — существующий, без — новый.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-izmenit-proizvodstvennyj-etap
type StageUpdateRequest struct {
	EnableHourAccounting model.Field[bool]              `json:"enableHourAccounting,omitzero"`
	ProcessingUnitCost   model.Field[float64]           `json:"processingUnitCost,omitzero"`
	LabourUnitCost       model.Field[float64]           `json:"labourUnitCost,omitzero"`
	StandardHourCost     model.Field[float64]           `json:"standardHourCost,omitzero"`
	StandardHourUnit     model.Field[float64]           `json:"standardHourUnit,omitzero"`
	PlannedEndDate       model.Field[model.DateTime]    `json:"plannedEndDate,omitzero"`
	Materials            model.Field[[]MaterialRequest] `json:"materials,omitzero"`
}

// GetStages читает одну страницу этапов задания id: GET entity/productionstage с обязательным
// фильтром productionTask=<href задания>; только limit и offset, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-spisok-proizvodstvennyh-etapov-proizvodstvennogo-zadaniya
func (s *Service) GetStages(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Stage], error) {
	filter := "productionTask=" + s.client.Href("entity/productiontask/"+id)
	return getPage[Stage](ctx, s.client, "get productionstage of productiontask "+id, id, "entity/productionstage", page, filter, opts)
}

// UpdateStage: PUT entity/productionstage/{stageID}; пустой input отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-izmenit-proizvodstvennyj-etap
func (s *Service) UpdateStage(ctx context.Context, stageID string, input StageUpdateRequest, opts ...transport.RequestOption) (*Stage, error) {
	transport.ResetResponseInfo(opts)
	operation := "update productionstage " + stageID
	if !validID(stageID) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	var result Stage
	resp, err := s.client.Send(ctx, http.MethodPut, "entity/productionstage/"+stageID, nil, input, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}
