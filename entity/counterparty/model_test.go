package counterparty_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/model"
)

const apiURL = "https://api.moysklad.ru/api/remap/1.2/"

func TestCounterpartyFixtures(t *testing.T) {
	for _, tc := range []struct {
		file           string
		want           counterparty.Counterparty
		attributeCount int
	}{
		{"counterparty.json", counterparty.Counterparty{
			Meta: responseMeta("counterparty", "12a8b923-692c-11e6-8a84-bae500000056"),
			ID:   "12a8b923-692c-11e6-8a84-bae500000056", SyncID: "12a8b923-692c-11e6-8a84-bae500000056",
			AccountID: "abcdef12-1234-1234-1234-123456789012",
			Created:   date(t, "2024-01-10T06:00:00Z"), Updated: date(t, "2024-01-15T07:30:00Z"),
			Name: "ООО Тестовый контрагент", Description: "testDescirption", Code: "CP-001", ExternalCode: "EXT-CP-001",
			Shared: true, CompanyType: model.CompanyTypeLegal,
			LegalTitle:     "Индивидуальный предприниматель Чекамеев Алексей Петрович",
			LegalFirstName: "Алексей", LegalLastName: "Чекамеев", LegalMiddleName: "Петрович",
			INN: "132808730606", OKPO: "122121", OGRNIP: "313132804400022",
			CertificateNumber: "12121212", CertificateDate: date(t, "2026-02-03T21:00:00Z"),
			Email: "888888888@mail.ru", Phone: "888888888", Fax: "8888888",
			ActualAddress:     "430008, Россия, Респ Мордовия, г Саранск, рп Луховка",
			LegalAddress:      "430008, Россия, Респ Мордовия, г Саранск, рп Луховка",
			ActualAddressFull: fixtureAddress(), LegalAddressFull: fixtureAddress(),
			Tags: []string{"uyiuuiy"}, DiscountCardNumber: "1111", BonusPoints: 5,
			Owner: &model.Reference{Meta: model.Meta{
				Href:         apiURL + "entity/employee/b2c670fa-01b5-11f1-0a81-072700000055",
				MetadataHref: apiURL + "entity/employee/metadata", Type: "employee", MediaType: "application/json",
				UUIDHref: "https://online.moysklad.ru/app/#employee/edit?id=b2c670fa-01b5-11f1-0a81-072700000055",
			}},
			Group: &model.Reference{Meta: responseMeta("group", "b205f9b3-01b5-11f1-0a81-181d00000003")},
			BonusProgram: &model.Reference{Meta: model.Meta{
				Href:         apiURL + "entity/bonusprogram/1223d051-ba76-11e8-3353-995e0000005a",
				MetadataHref: apiURL + "entity/bonusprogram/metadata", Type: "bonusprogram", MediaType: "application/json",
				UUIDHref: "https://online.moysklad.ru/app/#discount/edit?id=1223d051-ba76-11e8-3353-995e0000005a",
			}},
			State: &model.Reference{Meta: model.Meta{
				Href:         apiURL + "entity/counterparty/metadata/states/b3b37a9b-01b5-11f1-0a81-0727000000d2",
				MetadataHref: apiURL + "entity/counterparty/metadata", Type: "state", MediaType: "application/json",
			}},
			PriceType: &model.PriceType{Meta: model.Meta{
				Href: apiURL + "context/companysettings/pricetype/b35e9728-01b5-11f1-0a81-0727000000aa",
				Type: "pricetype", MediaType: "application/json",
			}, ID: "b35e9728-01b5-11f1-0a81-0727000000aa", Name: "Цена продажи", ExternalCode: "cbcf493b-55bc-11d9-848a-00112f43529a"},
		}, 3},
		{"counterparty_legal.json", counterparty.Counterparty{
			Meta: responseMeta("counterparty", "12a8b923-692c-11e6-8a84-bae500000057"),
			ID:   "12a8b923-692c-11e6-8a84-bae500000057", AccountID: "abcdef12-1234-1234-1234-123456789012",
			Created: date(t, "2024-01-10T06:00:00Z"), Updated: date(t, "2024-01-15T07:30:00Z"),
			Name: "ООО Юридическое лицо", Code: "CP-002", ExternalCode: "EXT-CP-002", Shared: true,
			CompanyType: model.CompanyTypeLegal, LegalTitle: "Общество с ограниченной ответственностью \"Юридическое лицо\"",
			LegalAddress: "г. Москва, ул. Примерная, д. 1", ActualAddress: "г. Москва, ул. Примерная, д. 1, офис 100",
			INN: "7701234567", KPP: "770101001", OGRN: "1027700000001", OKPO: "12345678",
			Email: "legal@example.com", Phone: "+7 (495) 123-45-67", Fax: "+7 (495) 123-45-68",
			Tags: []string{}, Group: &model.Reference{Meta: responseMeta("group", "11111111-1111-1111-1111-111111111111")},
		}, 1},
	} {
		data, err := os.ReadFile("../../testdata/upstream/" + tc.file)
		if err != nil {
			t.Fatal(err)
		}
		for _, codec := range codecs {
			t.Run(tc.file+"/"+codec.name, func(t *testing.T) {
				var got counterparty.Counterparty
				if err := codec.unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				attributes := got.Attributes
				got.Attributes = nil
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("counterparty = %+v, want %+v", got, tc.want)
				}
				if len(attributes) != tc.attributeCount {
					t.Fatalf("attributes = %+v", attributes)
				}
				first := attributes[0]
				attributeResource := "counterparty"
				if tc.file == "counterparty_legal.json" {
					attributeResource = "couterparty"
				}
				if first.ID != "0d129ff5-2c8c-11e6-8a84-bae5000000f3" || first.Name != "AttributeName1" || first.Type != "double" ||
					first.Meta != (model.Meta{Href: apiURL + "entity/" + attributeResource + "/metadata/attributes/0d129ff5-2c8c-11e6-8a84-bae5000000f3", Type: "attributemetadata", MediaType: "application/json"}) {
					t.Errorf("attribute = %+v", first)
				}
				if value, err := first.AsDouble(); err != nil || value != 0.75 {
					t.Errorf("double = %v, %v", value, err)
				}
				if len(attributes) == 3 {
					for i, want := range []struct{ id, name, typ string }{
						{"0d12a9a5-2c8c-11e6-8a84-bae5000000f4", "AttributeName2", "time"},
						{"0d12b1e7-2c8c-11e6-8a84-bae5000000f5", "AttributeName3", "boolean"},
					} {
						a := attributes[i+1]
						if a.ID != want.id || a.Name != want.name || a.Type != want.typ || a.Meta != (model.Meta{Href: apiURL + "entity/counterparty/metadata/attributes/" + want.id, Type: "attributemetadata", MediaType: "application/json"}) {
							t.Errorf("attribute = %+v", a)
						}
					}
					if value, err := attributes[1].AsTime(); err != nil || !time.Time(value).Equal(time.Time(date(t, "2016-06-07T09:52:33Z"))) {
						t.Errorf("time = %v, %v", value, err)
					}
					if value, err := attributes[2].AsBool(); err != nil || value {
						t.Errorf("bool = %v, %v", value, err)
					}
				}
			})
		}
	}
}

func responseMeta(typ, id string) model.Meta {
	return model.Meta{Href: apiURL + "entity/" + typ + "/" + id, MetadataHref: apiURL + "entity/" + typ + "/metadata", Type: model.MetaType(typ), MediaType: "application/json"}
}

func fixtureAddress() *model.Address {
	return &model.Address{PostalCode: "430008", City: "г Саранск, рп Луховка", Comment: "выыфвы",
		Country: &model.Reference{Meta: model.Meta{
			Href: apiURL + "entity/country/9df7c2c3-7782-4c5c-a8ed-1102af611608", MetadataHref: apiURL + "entity/country/metadata", Type: "country", MediaType: "application/json",
			UUIDHref: "https://online.moysklad.ru/app/#country/edit?id=9df7c2c3-7782-4c5c-a8ed-1102af611608",
		}},
		Region: &model.Reference{Meta: responseMeta("region", "00000000-0000-0000-0000-000000000013")},
	}
}

func date(t *testing.T, value string) model.DateTime {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatal(err)
	}
	return model.DateTime(parsed.In(location))
}

func TestCounterpartyOptionalAndNewValues(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, number := range []string{"0", "0.0", "123456.75"} {
				var got counterparty.Counterparty
				data := []byte(`{"birthDate":"1990-05-06 00:00:00.000","certificateDate":null,"sex":"future-sex","companyType":"future-region","archived":true,"salesAmount":` + number + `,"inn":"00123","kpp":"00456","ogrn":"00789","ogrnip":"00012","okpo":"00034","actualAddressFull":{"addInfo":"entry","apartment":"03","house":"02","street":"Лесная","fiasCode__ru":"00001","postalCode":"001234","country":null,"futurePart":1},"legalAddressFull":null,"futureField":{},"mod__requisites__uz":{"inn":"123"},"mod__requisites__kz":{}}`)
				if err := codec.unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				wantAmount := 0.0
				if number == "123456.75" {
					wantAmount = 123456.75
				}
				want := counterparty.Counterparty{
					BirthDate: date(t, "1990-05-05T20:00:00Z"), Sex: "future-sex", CompanyType: "future-region", Archived: true,
					SalesAmount: wantAmount, INN: "00123", KPP: "00456", OGRN: "00789", OGRNIP: "00012", OKPO: "00034",
					ActualAddressFull: &model.Address{AddInfo: "entry", Apartment: "03", House: "02", Street: "Лесная", FIASCodeRU: "00001", PostalCode: "001234"},
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("counterparty = %+v, want %+v", got, want)
				}
			}
			for _, input := range []string{`{}`, `{"created":null,"updated":null,"birthDate":null,"certificateDate":null,"owner":null,"group":null,"bonusProgram":null,"state":null,"priceType":null,"actualAddressFull":null,"legalAddressFull":null,"attributes":null,"tags":null}`} {
				var got counterparty.Counterparty
				if err := codec.unmarshal([]byte(input), &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, counterparty.Counterparty{}) {
					t.Errorf("absent/null fields invented: %+v", got)
				}
			}
		})
	}
}
