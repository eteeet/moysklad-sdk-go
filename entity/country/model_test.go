package country_test

import (
	"os"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/country"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCountryUpstreamFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/country.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got country.Country
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			want := country.Country{
				Meta: model.Meta{
					Href:         "https://api.moysklad.ru/api/remap/1.2/entity/country/44444444-4444-4444-4444-444444444444",
					MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/country/metadata",
					Type:         "country", MediaType: "application/json",
				},
				ID: "44444444-4444-4444-4444-444444444444", AccountID: "abcdef12-1234-1234-1234-123456789012",
				Name: "Россия", Description: "Российская Федерация", Code: "643", ExternalCode: "RU", Shared: false,
				Updated: model.DateTime(time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC)),
				Owner: &model.Reference{Meta: model.Meta{
					Href:         "https://api.moysklad.ru/api/remap/1.2/entity/employee/b2c670fa-01b5-11f1-0a81-072700000055",
					MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/employee/metadata", Type: "employee", MediaType: "application/json",
					UUIDHref: "https://online.moysklad.ru/app/#employee/edit?id=b2c670fa-01b5-11f1-0a81-072700000055",
				}},
				Group: &model.Reference{Meta: model.Meta{
					Href:         "https://api.moysklad.ru/api/remap/1.2/entity/group/b205f9b3-01b5-11f1-0a81-181d00000003",
					MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/group/metadata", Type: "group", MediaType: "application/json",
				}},
			}
			if got.Meta != want.Meta || got.ID != want.ID || got.AccountID != want.AccountID || got.Name != want.Name ||
				got.Description != want.Description || got.Code != want.Code || got.ExternalCode != want.ExternalCode || got.Shared != want.Shared ||
				!time.Time(got.Updated).Equal(time.Time(want.Updated)) || got.Owner == nil || got.Group == nil ||
				*got.Owner != *want.Owner || *got.Group != *want.Group {
				t.Fatalf("country = %+v", got)
			}
		})
	}
}

func TestCountryPredefinedFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/custom/country/predefined.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got country.Country
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			wantMeta := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/country/000d77a9-3000-4f81-a995-6b9cffdee1d2",
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/country/metadata",
				Type:         "country", MediaType: "application/json",
				UUIDHref: "https://online.moysklad.ru/app/#country/edit?id=000d77a9-3000-4f81-a995-6b9cffdee1d2",
			}
			if got.Meta != wantMeta || got.ID != "000d77a9-3000-4f81-a995-6b9cffdee1d2" ||
				got.Name != "Марокко" || got.Description != "Королевство Марокко" || got.Code != "504" || got.ExternalCode != "504" ||
				!time.Time(got.Updated).Equal(time.Date(2012, 11, 2, 7, 4, 13, 0, time.UTC)) {
				t.Fatalf("country = %+v", got)
			}
			if got.AccountID != "" || got.Owner != nil || got.Group != nil || got.Shared {
				t.Errorf("absent fields were invented: %+v", got)
			}
		})
	}
}

func TestCountryMissingAndUnknownFields(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, data := range []string{
				`{"name":"SDK Country","code":"0007","shared":true,"future":{"enabled":true}}`,
				`{"name":"SDK Country","code":"0007","shared":true,"updated":null,"owner":null,"group":null}`,
			} {
				var got country.Country
				if err := codec.unmarshal([]byte(data), &got); err != nil {
					t.Fatal(err)
				}
				if got.Name != "SDK Country" || got.Code != "0007" || !got.Shared {
					t.Fatalf("country = %+v", got)
				}
				if got.Meta != (model.Meta{}) || got.ID != "" || got.AccountID != "" || got.Description != "" ||
					got.ExternalCode != "" || !got.Updated.IsZero() || got.Owner != nil || got.Group != nil {
					t.Errorf("absent fields were invented: %+v", got)
				}
			}
		})
	}
}
