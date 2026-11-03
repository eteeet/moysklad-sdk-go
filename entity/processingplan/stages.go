package processingplan

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Этапы — entity/processingplan/{id}/stages: только чтение и изменение. Создавать и удалять этапы
// источник не позволяет: они строго соответствуют этапам в позициях техпроцесса. Повторов нет.

func (s *Service) stages(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "processingplan", ParentID: id, Path: "stages", Item: "stage", MetaType: "processingplanstages"}
}

// GetStages читает одну страницу этапов: только limit и offset, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-etapy-tehkarty
func (s *Service) GetStages(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Stage], error) {
	return nested.List[Stage](ctx, s.stages(id), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-otdelnyj-etap-tehkarty
func (s *Service) GetStage(ctx context.Context, id, stageID string, opts ...transport.RequestOption) (*Stage, error) {
	return nested.Get[Stage](ctx, s.stages(id), stageID, opts)
}

// UpdateStage: обязательных полей нет; пустой input отклоняется до запроса. Одновременно cost
// этапа и cost техкарты сервер не принимает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-izmenit-otdelnyj-etap-tehkarty
func (s *Service) UpdateStage(ctx context.Context, id, stageID string, input StageRequest, opts ...transport.RequestOption) (*Stage, error) {
	return nested.Update[Stage](ctx, s.stages(id), stageID, input, opts)
}
