package expenseitem_test

import (
	"os"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestExpenseItemUpstreamFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/expense_item.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got expenseitem.ExpenseItem
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			wantMeta := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/expenseitem/fb0a4b75-2e58-11e6-8a84-bae500000058",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/expenseitem/metadata",
				Type:         "expenseitem", MediaType: "application/json",
			}
			if got.Meta != wantMeta || got.ID != "fb0a4b75-2e58-11e6-8a84-bae500000058" ||
				got.AccountID != "f976ed28-2e58-11e6-8a84-bae500000001" || got.Name != "Аренда" ||
				got.Description != "Аренда" || got.Code != "Аренда" || got.ExternalCode != "IVslr34uhCUuglxPD7Idm0" || !got.Archived {
				t.Fatalf("expenseitem = %+v", got)
			}
			if !time.Time(got.Updated).Equal(time.Date(2016, 6, 9, 15, 43, 58, 0, time.UTC)) {
				t.Errorf("updated = %s", time.Time(got.Updated))
			}
			if got.OperatingExpenses || got.Shared || got.Owner != nil || got.Group != nil {
				t.Errorf("absent fields were invented: %+v", got)
			}
		})
	}
}

func TestExpenseItemPredefinedFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/custom/expenseitem/predefined.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got expenseitem.ExpenseItem
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			wantMeta := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/expenseitem/1be2395a-0479-11e5-baee-448a5b426e7e",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/expenseitem/metadata",
				Type:         "expenseitem", MediaType: "application/json",
			}
			if got.Meta != wantMeta || got.ID != "1be2395a-0479-11e5-baee-448a5b426e7e" ||
				got.Name != "Возврат" || got.Description != "Расходы по возвратам не учитываются в отчете «Прибыли и убытки»" ||
				got.Code != "3" || got.ExternalCode != "3" || got.Archived || got.OperatingExpenses {
				t.Fatalf("expenseitem = %+v", got)
			}
			if got.AccountID != "" || got.Owner != nil || got.Group != nil || got.Shared {
				t.Errorf("absent fields were invented: %+v", got)
			}
			if !time.Time(got.Updated).Equal(time.Date(2015, 5, 27, 14, 3, 10, 0, time.UTC)) {
				t.Errorf("updated = %s", time.Time(got.Updated))
			}
		})
	}
}

func TestExpenseItemSupplementalFields(t *testing.T) {
	for _, codec := range codecs {
		for _, tc := range []struct {
			name                        string
			data                        string
			operating, archived, shared bool
		}{
			{"true", `{"code":"003","operatingExpenses":true,"archived":true,"shared":true}`, true, true, true},
			{"false", `{"code":"003","operatingExpenses":false,"archived":false,"shared":false}`, false, false, false},
			{"absent", `{"code":"003"}`, false, false, false},
		} {
			t.Run(codec.name+"/"+tc.name, func(t *testing.T) {
				var got expenseitem.ExpenseItem
				if err := codec.unmarshal([]byte(tc.data), &got); err != nil {
					t.Fatal(err)
				}
				if got.Code != "003" || got.OperatingExpenses != tc.operating || got.Archived != tc.archived || got.Shared != tc.shared {
					t.Fatalf("expenseitem = %+v", got)
				}
			})
		}
		t.Run(codec.name+"/references and unknown fields", func(t *testing.T) {
			var got expenseitem.ExpenseItem
			data := []byte(`{"owner":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111","type":"employee","mediaType":"application/json"}},"group":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/group/22222222-2222-2222-2222-222222222222","type":"group","mediaType":"application/json"}},"futureField":{"enabled":true}}`)
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			wantOwner := model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111", Type: "employee", MediaType: "application/json"}
			wantGroup := model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/group/22222222-2222-2222-2222-222222222222", Type: "group", MediaType: "application/json"}
			if got.Owner == nil || got.Owner.Meta != wantOwner || got.Group == nil || got.Group.Meta != wantGroup {
				t.Fatalf("references = %+v, %+v", got.Owner, got.Group)
			}
			if got.Meta != (model.Meta{}) || got.ID != "" || got.AccountID != "" || got.Name != "" || got.Description != "" ||
				got.Code != "" || got.ExternalCode != "" || !got.Updated.IsZero() || got.Shared || got.Archived || got.OperatingExpenses {
				t.Errorf("absent fields were invented: %+v", got)
			}
		})
	}
}
