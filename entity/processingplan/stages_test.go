package processingplan_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/processingplan"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Пример ответа «Получить отдельный этап Техкарты» источника.
func TestStageFixture(t *testing.T) {
	data := read(t, "custom/processingplan/stage.json")
	for _, codec := range documenttest.Codecs {
		var got processingplan.Stage
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.ID != rowID || got.AccountID != "2f09c84d-bbd0-11ed-ac1c-000f00000001" || got.Meta.Type != "processingplanstages" || got.Cost != 1100 || got.LabourCost != 100 ||
			got.StandardHour != 200 || got.EnableHourAccounting || got.ProcessingProcessPosition == nil || got.ProcessingProcessPosition.Meta.Type != "processingprocessposition" {
			t.Fatalf("%s: stage = %+v", codec.Name, got)
		}
	}
}

func TestStagesHTTP(t *testing.T) {
	const row = `{"meta":{"type":"processingplanstages"},"id":"` + rowID + `","cost":5.5,"labourCost":10.0,"standardHour":15.0,"standardHourCost":300,"enableHourAccounting":true}`
	check := func(got *processingplan.Stage) {
		if got.ID != rowID || got.Cost != 5.5 || got.LabourCost != 10 || got.StandardHour != 15 || got.StandardHourCost != 300 || !got.EnableHourAccounting {
			t.Fatalf("stage = %+v", got)
		}
	}
	t.Run("get stages", func(t *testing.T) {
		pageCases(t, planPath+"/stages", "limit=2&offset=7", `{"meta":{"href":"h","type":"processingplanstages","size":1},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[processingplan.Stage], error) {
				return service(c).GetStages(ctx, testID, query.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[processingplan.Stage]) { check(&got.Rows[0]) })
	})
	t.Run("get stage", func(t *testing.T) {
		httpCases(t, "GET", planPath+"/stages/"+rowID, "", row, func(ctx context.Context, c *transport.Client) (*processingplan.Stage, error) {
			return service(c).GetStage(ctx, testID, rowID)
		}, check)
	})
	// Пример источника: {"cost": 5.5}.
	t.Run("update stage", func(t *testing.T) {
		httpCases(t, "PUT", planPath+"/stages/"+rowID, `{"cost":5.5}`, row, func(ctx context.Context, c *transport.Client) (*processingplan.Stage, error) {
			return service(c).UpdateStage(ctx, testID, rowID, processingplan.StageRequest{Cost: model.Value(5.5)})
		}, check)
	})
}

// Создавать и удалять этапы источник не позволяет — таких методов нет.
func TestStagesReadUpdateOnly(t *testing.T) {
	typ := reflect.TypeFor[*processingplan.Service]()
	for _, name := range []string{"CreateStages", "CreateStage", "DeleteStage", "DeleteManyStages"} {
		if _, ok := typ.MethodByName(name); ok {
			t.Errorf("unexpected method %s", name)
		}
	}
}

func TestStagesValidation(t *testing.T) {
	noHTTP(t, func(s *processingplan.Service, ctx context.Context) map[string]func(...transport.RequestOption) error {
		return map[string]func(...transport.RequestOption) error{
			"list bad id": func(o ...transport.RequestOption) error {
				_, err := s.GetStages(ctx, "bad", query.PageOptions{}, o...)
				return err
			},
			"list fields": func(o ...transport.RequestOption) error {
				_, err := s.GetStages(ctx, testID, query.PageOptions{Fields: "stock"}, o...)
				return err
			},
			"list limit": func(o ...transport.RequestOption) error {
				_, err := s.GetStages(ctx, testID, query.PageOptions{Limit: 1001}, o...)
				return err
			},
			"get bad stage": func(o ...transport.RequestOption) error {
				_, err := s.GetStage(ctx, testID, "bad", o...)
				return err
			},
			"update empty": func(o ...transport.RequestOption) error {
				_, err := s.UpdateStage(ctx, testID, rowID, processingplan.StageRequest{}, o...)
				return err
			},
			"update bad plan": func(o ...transport.RequestOption) error {
				_, err := s.UpdateStage(ctx, "bad", rowID, processingplan.StageRequest{Cost: model.Value(1.0)}, o...)
				return err
			},
		}
	})
}
