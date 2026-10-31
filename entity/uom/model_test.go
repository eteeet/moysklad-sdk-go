package uom_test

import (
	"os"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestUOMFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/uom.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got uom.UOM
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			wantMeta := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/uom/77777777-7777-7777-7777-777777777777",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/uom/metadata",
				Type:         "uom",
				MediaType:    "application/json",
			}
			if got.Meta != wantMeta || got.ID != "77777777-7777-7777-7777-777777777777" ||
				got.AccountID != "abcdef12-1234-1234-1234-123456789012" ||
				got.Name != "шт" || got.Description != "Штука" ||
				got.Code != "796" || got.ExternalCode != "PCE" || got.Shared {
				t.Fatalf("uom = %+v", got)
			}
			if !time.Time(got.Updated).Equal(time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC)) {
				t.Errorf("updated = %s", time.Time(got.Updated))
			}
			wantOwner := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/employee/b2c670fa-01b5-11f1-0a81-072700000055",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/employee/metadata",
				Type:         "employee",
				MediaType:    "application/json",
				UUIDHref:     "https://online.moysklad.ru/app/#employee/edit?id=b2c670fa-01b5-11f1-0a81-072700000055",
			}
			if got.Owner == nil || got.Owner.Meta != wantOwner {
				t.Errorf("owner = %+v", got.Owner)
			}
			wantGroup := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/group/11111111-1111-1111-1111-111111111111",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/group/metadata",
				Type:         "group",
				MediaType:    "application/json",
			}
			if got.Group == nil || got.Group.Meta != wantGroup {
				t.Errorf("group = %+v", got.Group)
			}
		})
	}
}

func TestUOMPredefinedFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/custom/uom/predefined.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got uom.UOM
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			wantMeta := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/uom/061721df-9197-49a5-b637-7f5b4d3be969",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/uom/metadata",
				Type:         "uom",
				MediaType:    "application/json",
			}
			if got.Meta != wantMeta || got.ID != "061721df-9197-49a5-b637-7f5b4d3be969" ||
				got.Name != "дюйм" || got.Description != "Дюйм (25,4 мм)" ||
				got.Code != "039" || got.ExternalCode != "039" {
				t.Fatalf("uom = %+v", got)
			}
			if got.AccountID != "" || got.Owner != nil || got.Group != nil || got.Shared {
				t.Errorf("absent fields were invented: %+v", got)
			}
			if !time.Time(got.Updated).Equal(time.Date(2012, 11, 2, 7, 7, 44, 0, time.UTC)) {
				t.Errorf("updated = %s", time.Time(got.Updated))
			}
		})
	}
}

func TestUOMUnknownFields(t *testing.T) {
	data := []byte(`{"name":"Попугаи","code":"papagei","shared":true,"futureField":{"enabled":true}}`)
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got uom.UOM
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.Name != "Попугаи" || got.Code != "papagei" || !got.Shared {
				t.Fatalf("uom = %+v", got)
			}
			if got.Meta != (model.Meta{}) || got.ID != "" || got.AccountID != "" ||
				got.Description != "" || got.ExternalCode != "" || !got.Updated.IsZero() ||
				got.Owner != nil || got.Group != nil {
				t.Errorf("absent fields were invented: %+v", got)
			}
		})
	}
}
