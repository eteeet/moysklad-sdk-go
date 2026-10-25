package model_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Вторая страница вебхуков: https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-meta#3-metadannye-kollekcii
const webhookPage = `{
  "context": {"employee": {"meta": {"href": "https://api.moysklad.ru/api/remap/1.2/context/employee", "type": "employee", "mediaType": "application/json"}}},
  "meta": {
    "href": "https://api.moysklad.ru/api/remap/1.2/entity/webhook?limit=1&offset=1",
    "type": "webhook",
    "mediaType": "application/json",
    "size": 3,
    "limit": 1,
    "offset": 1,
    "nextHref": "https://api.moysklad.ru/api/remap/1.2/entity/webhook?offset=2&limit=1",
    "previousHref": "https://api.moysklad.ru/api/remap/1.2/entity/webhook?offset=0&limit=1"
  },
  "rows": [{"meta": {"href": "https://api.moysklad.ru/api/remap/1.2/entity/webhook/07598ccd-072c-11e8-7a6c-d2a90010c896", "type": "webhook", "mediaType": "application/json"}}]
}`

func TestCollectionUnmarshal(t *testing.T) {
	var got model.Collection[model.Reference]
	if err := json.Unmarshal([]byte(webhookPage), &got); err != nil {
		t.Fatal(err)
	}
	if got.Context == nil || got.Context.Employee == nil || got.Context.Employee.Meta.Type != "employee" {
		t.Fatalf("context = %+v", got.Context)
	}

	m := got.Meta
	if m.Type != "webhook" || m.Size != 3 || m.Limit != 1 || m.Offset != 1 {
		t.Fatalf("meta = %+v", m)
	}
	if m.NextHref != "https://api.moysklad.ru/api/remap/1.2/entity/webhook?offset=2&limit=1" {
		t.Fatalf("nextHref = %q", m.NextHref)
	}
	if m.PreviousHref != "https://api.moysklad.ru/api/remap/1.2/entity/webhook?offset=0&limit=1" {
		t.Fatalf("previousHref = %q", m.PreviousHref)
	}
	if len(got.Rows) != 1 || got.Rows[0].Meta.Type != "webhook" {
		t.Fatalf("rows = %+v", got.Rows)
	}
}

func TestMetaPermitted(t *testing.T) {
	var absent, present model.Meta
	if err := json.Unmarshal([]byte(`{"href":"h"}`), &absent); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(`{"href":"h","permitted":false}`), &present); err != nil {
		t.Fatal(err)
	}
	if absent.Permitted != nil {
		t.Fatalf("permitted = %v, want nil when field is absent", *absent.Permitted)
	}
	if present.Permitted == nil || *present.Permitted {
		t.Fatalf("permitted = %v, want false", present.Permitted)
	}
}
