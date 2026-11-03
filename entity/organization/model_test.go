package organization_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestOrganizationUpstream(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/organization.json")
	if err != nil {
		t.Fatal(err)
	}
	const root = "https://api.moysklad.ru/api/remap/1.2/entity/"
	want := organization.Organization{
		Meta: model.Meta{Href: root + "organization/4b9d5bec-0575-11e6-9464-e4de00000008", MetadataHref: root + "organization/metadata", Type: "organization", MediaType: "application/json"},
		ID:   "4b9d5bec-0575-11e6-9464-e4de00000008", AccountID: "84e60e93-f504-11e5-8a84-bae500000008",
		Created: date(t, "2007-02-07T17:16:41+03:00"), Updated: date(t, "2016-04-18T17:57:27+03:00"),
		Name: "ОАО СветПром", Description: "юрлицо, делающее маленькую прибыль", Code: "666", ExternalCode: "666АААА666",
		CompanyType: model.CompanyTypeLegal, TrackingContractNumber: "12345678", TrackingContractDate: date(t, "2007-02-07T00:00:00+03:00"),
		LegalTitle: "ООО Великое Свет Пром", INN: "87654321", KPP: "15312532", OGRN: "12345", OKPO: "12345", OGRNIP: "58632598y21jk",
		CertificateNumber: "klghvew983412", CertificateDate: date(t, "2016-04-30T00:00:00+03:00"),
		Email: "svetprom@mail.svet", Phone: "22222222", Fax: "bello123", PayerVAT: true, Director: "Администратор", ChiefAccountant: "Администратор",
		ActualAddress: "125009, Россия, г Москва, Москва, ул Тверская, 1, 111, addInfo",
		LegalAddress:  "125009, Россия, г Москва, Москва, ул Тверская, 1, 123, addInfo",
		Owner:         fixtureReference(root, "employee", "faba7f37-2e58-11e6-8a84-bae500000028"),
		Group:         fixtureReference(root, "group", "f97aa1fb-2e58-11e6-8a84-bae500000002"),
		BonusProgram:  fixtureReference(root, "bonusprogram", "129626ee-ac91-11e9-ac12-000d00000009"),
		Attributes:    []model.Attribute{{Meta: model.Meta{Href: root + "organization/metadata/attributes/0d129ff5-2c8c-11e6-8a84-bae5000000f3", Type: "attributemetadata", MediaType: "application/json"}, ID: "0d129ff5-2c8c-11e6-8a84-bae5000000f3", Name: "AttributeName1", Type: "double", Value: model.Value(jsontext.Value("0.75"))}},
	}
	want.ActualAddressFull = &model.Address{
		PostalCode: "125009", City: "Москва", Street: "ул Тверская", House: "1", Apartment: "111", AddInfo: "addinfo", Comment: "some words about address",
		Country: fixtureReference(root, "country", "9df7c2c3-7782-4c5c-a8ed-1102af611608"), Region: fixtureReference(root, "region", "00000000-0000-0000-0000-000000000077"),
	}
	want.LegalAddressFull = &model.Address{
		PostalCode: "125009", City: "Москва", Street: "ул Тверская", House: "1", Apartment: "123", AddInfo: "addinfo", Comment: "some words about address",
		Country: fixtureReference(root, "country", "9df7c2c3-7782-4c5c-a8ed-1102af611608"), Region: fixtureReference(root, "region", "00000000-0000-0000-0000-000000000077"),
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got organization.Organization
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("organization = %+v, want %+v", got, want)
			}
			if value, err := got.Attributes[0].AsDouble(); err != nil || value != 0.75 {
				t.Errorf("attribute = %v, error = %v", value, err)
			}
		})
	}
}

func TestOrganizationOptionalFields(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name string
				data string
				want organization.Organization
			}{
				{"absent", `{}`, organization.Organization{}},
				{"null", `{"created":null,"updated":null,"certificateDate":null,"trackingContractDate":null,"actualAddressFull":null,"legalAddressFull":null,"owner":null,"group":null,"bonusProgram":null,"attributes":null}`, organization.Organization{}},
				{"additional fields", `{"name":"ИП","archived":true,"shared":true,"companyType":"futureKind","legalFirstName":"Иван","legalLastName":"Петров","legalMiddleName":"Иванович","inn":"001234567890","kpp":null,"ogrn":null,"syncId":"001-sync","bonusPoints":17,"payerVat":false,"director":"Иван Петров","directorPosition":"Директор","chiefAccountant":"Мария Петрова","trackingContractNumber":"00042","trackingContractDate":"2025-01-02 03:04:05.000","certificateDate":"2020-02-03 00:00:00","isEgaisEnable":true,"actualAddressFull":{"fiasCode__ru":"0001","country":null,"futurePart":{}},"legalAddressFull":{"fiasCode__ru":"0002","region":null},"accounts":{"rows":[{}]},"mod__requisites__uz":{},"mod__requisites__kz":{},"companyVat__ru":{},"advancePaymentVat":20,"chiefAccountSign":{},"directorSign":{},"stamp":{},"fsrarId":"ignored","utmUrl":"ignored","futureField":true}`,
					organization.Organization{Name: "ИП", Archived: true, Shared: true, CompanyType: model.CompanyType("futureKind"), LegalFirstName: "Иван", LegalLastName: "Петров", LegalMiddleName: "Иванович", INN: "001234567890", SyncID: "001-sync", BonusPoints: 17, Director: "Иван Петров", DirectorPosition: "Директор", ChiefAccountant: "Мария Петрова", TrackingContractNumber: "00042", TrackingContractDate: date(t, "2025-01-02T03:04:05+03:00"), CertificateDate: date(t, "2020-02-03T00:00:00+03:00"), IsEGAISEnable: true, ActualAddressFull: &model.Address{FIASCodeRU: "0001"}, LegalAddressFull: &model.Address{FIASCodeRU: "0002"},
						CompanyVat: &organization.CompanyVat{}, AdvancePaymentVat: 20, ChiefAccountSign: &organization.Image{}, DirectorSign: &organization.Image{}, Stamp: &organization.Image{},
						RequisitesUZ: &model.RequisitesUZ{}, RequisitesKZ: &model.RequisitesKZ{}}},
				{"egais false", `{"isEgaisEnable":false,"payerVat":false}`, organization.Organization{}},
				{"regional type", `{"companyType":"individualKZ"}`, organization.Organization{CompanyType: model.CompanyTypeIndividualKZ}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got organization.Organization
					if err := codec.unmarshal([]byte(tc.data), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Errorf("organization = %+v, want %+v", got, tc.want)
					}
				})
			}
		})
	}
}

func fixtureReference(root, kind, id string) *model.Reference {
	return &model.Reference{Meta: model.Meta{Href: root + kind + "/" + id, MetadataHref: root + kind + "/metadata", Type: model.MetaType(kind), MediaType: "application/json"}}
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
