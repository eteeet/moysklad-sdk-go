package model_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestPriceValueJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				input string
				want  float64
			}{
				{`{"value":0}`, 0},
				{`{"value":500.0}`, 500},
				{`{"value":23553000.0}`, 23553000},
				{`{"value":346347237000.0}`, 346347237000},
				{`{"value":9007199254740991}`, 9007199254740991},
				{`{"value":500.25}`, 500.25},
				{`{"value":0.1}`, 0.1},
			} {
				t.Run(tc.input, func(t *testing.T) {
					var got model.Price
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if got.Value != tc.want {
						t.Fatalf("value = %v, want %v", got.Value, tc.want)
					}
					encoded, err := codec.marshal(got)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.input)
					var roundTrip model.Price
					if err := codec.unmarshal(encoded, &roundTrip); err != nil {
						t.Fatal(err)
					}
					if roundTrip.Value != tc.want {
						t.Fatalf("value after roundtrip = %v, want %v", roundTrip.Value, tc.want)
					}
				})
			}
		})
	}
}

func TestPriceCurrencyJSON(t *testing.T) {
	const input = `{"value":0,"currency":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/currency/22222222-2222-2222-2222-222222222222","type":"currency","mediaType":"application/json"}}}`
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got model.Price
			if err := codec.unmarshal([]byte(input), &got); err != nil {
				t.Fatal(err)
			}
			want := model.Reference{Meta: model.Meta{
				Href:      "https://api.moysklad.ru/api/remap/1.2/entity/currency/22222222-2222-2222-2222-222222222222",
				Type:      "currency",
				MediaType: "application/json",
			}}
			if got.Currency != want {
				t.Fatalf("currency = %+v, want %+v", got.Currency, want)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, encoded, input)
		})
	}
}

func TestSalePriceJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				price model.SalePrice
				want  string
			}{
				{"zero", model.SalePrice{}, `{"value":0}`},
				{
					"price type meta only",
					model.SalePrice{Value: 500, PriceType: model.PriceType{Meta: model.Meta{
						Href: "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/33333333-3333-3333-3333-333333333333",
					}}},
					`{"value":500,"priceType":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/33333333-3333-3333-3333-333333333333"}}}`,
				},
				{
					"full",
					model.SalePrice{
						Value: 346347237000,
						Currency: model.Reference{Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/currency/22222222-2222-2222-2222-222222222222",
							Type:      "currency",
							MediaType: "application/json",
						}},
						PriceType: model.PriceType{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/33333333-3333-3333-3333-333333333333",
								Type:      "pricetype",
								MediaType: "application/json",
							},
							ID:           "33333333-3333-3333-3333-333333333333",
							Name:         "Цена продажи",
							ExternalCode: "default",
						},
					},
					`{
					  "value":346347237000.0,
					  "currency":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/currency/22222222-2222-2222-2222-222222222222","type":"currency","mediaType":"application/json"}},
					  "priceType":{
					    "meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/33333333-3333-3333-3333-333333333333","type":"pricetype","mediaType":"application/json"},
					    "id":"33333333-3333-3333-3333-333333333333","name":"Цена продажи","externalCode":"default"
					  }
					}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.price)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got model.SalePrice
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.price) {
						t.Fatalf("price = %+v, want %+v", got, tc.price)
					}
				})
			}
		})
	}
}

func TestPricesProductFixtures(t *testing.T) {
	type productPrices struct {
		BuyPrice   *model.Price      `json:"buyPrice"`
		MinPrice   *model.Price      `json:"minPrice"`
		SalePrices []model.SalePrice `json:"salePrices"`
	}
	currency := model.Reference{Meta: model.Meta{
		Href:         "https://api.moysklad.ru/api/remap/1.2/entity/currency/22222222-2222-2222-2222-222222222222",
		MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/currency/metadata",
		Type:         "currency",
		MediaType:    "application/json",
	}}
	full := productPrices{
		BuyPrice: &model.Price{Value: 50000, Currency: currency},
		MinPrice: &model.Price{Value: 1000, Currency: model.Reference{Meta: model.Meta{
			Href:         "https://api.moysklad.ru/api/remap/1.2/entity/currency/b31787dd-01b5-11f1-0a81-0727000000a9",
			MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/currency/metadata",
			Type:         "currency",
			MediaType:    "application/json",
			UUIDHref:     "https://online.moysklad.ru/app/#currency/edit?id=b31787dd-01b5-11f1-0a81-0727000000a9",
		}}},
		SalePrices: []model.SalePrice{{
			Value:    100000,
			Currency: currency,
			PriceType: model.PriceType{
				Meta: model.Meta{
					Href:      "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/33333333-3333-3333-3333-333333333333",
					Type:      "pricetype",
					MediaType: "application/json",
				},
				ID:           "33333333-3333-3333-3333-333333333333",
				Name:         "Цена продажи",
				ExternalCode: "default",
			},
		}},
	}
	for _, tc := range []struct {
		file string
		want productPrices
	}{
		{"product.json", productPrices{SalePrices: []model.SalePrice{}}},
		{"product_full.json", full},
	} {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile("../testdata/upstream/" + tc.file)
			if err != nil {
				t.Fatal(err)
			}
			for _, codec := range codecs {
				t.Run(codec.name, func(t *testing.T) {
					var got productPrices
					if err := codec.unmarshal(data, &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got.BuyPrice, tc.want.BuyPrice) {
						t.Fatalf("buyPrice = %+v, want %+v", got.BuyPrice, tc.want.BuyPrice)
					}
					if !reflect.DeepEqual(got.MinPrice, tc.want.MinPrice) {
						t.Fatalf("minPrice = %+v, want %+v", got.MinPrice, tc.want.MinPrice)
					}
					if !reflect.DeepEqual(got.SalePrices, tc.want.SalePrices) {
						t.Fatalf("salePrices = %+v, want %+v", got.SalePrices, tc.want.SalePrices)
					}
				})
			}
		})
	}
}

func TestPriceTypeFixture(t *testing.T) {
	data, err := os.ReadFile("../testdata/upstream/price_type.json")
	if err != nil {
		t.Fatal(err)
	}
	want := model.PriceType{
		Meta: model.Meta{
			Href:      "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/88888888-8888-8888-8888-888888888888",
			Type:      "pricetype",
			MediaType: "application/json",
		},
		ID:           "88888888-8888-8888-8888-888888888888",
		Name:         "Цена продажи",
		ExternalCode: "default",
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got model.PriceType
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("priceType = %+v, want %+v", got, want)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, encoded, string(data))
		})
	}
}

func TestPriceRequestFieldsJSON(t *testing.T) {
	type request struct {
		BuyPrice   model.Field[model.Price]       `json:"buyPrice,omitzero"`
		MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
		SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value request
				want  string
			}{
				{"unset", request{}, `{}`},
				{"zero buy price", request{BuyPrice: model.Value(model.Price{})}, `{"buyPrice":{"value":0}}`},
				{"zero min price", request{MinPrice: model.Value(model.Price{})}, `{"minPrice":{"value":0}}`},
				{"empty sale prices", request{SalePrices: model.Value([]model.SalePrice{})}, `{"salePrices":[]}`},
				{"nil sale prices", request{SalePrices: model.Value[[]model.SalePrice](nil)}, `{"salePrices":[]}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
				})
			}
		})
	}
}
