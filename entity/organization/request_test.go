package organization_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCreateRequestMinimal(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(organization.CreateRequest{Name: "SDK Organization"})
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, []byte(`{"name":"SDK Organization"}`))
		})
	}
}

func TestCreateRequestFull(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(fullCreateRequest(t))
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, organizationGolden(t, "create"))
		})
	}
}

func TestUpdateRequestFull(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(fullUpdateRequest(t))
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, organizationGolden(t, "update"))
		})
	}
}

func TestRequestExplicitValues(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input any
				want  string
			}{
				{"zero create", organization.CreateRequest{}, `{"name":""}`},
				{"zero update", organization.UpdateRequest{}, `{}`},
				{"update sync ID", organization.UpdateRequest{SyncID: model.Value("11111111-1111-1111-1111-111111111111")}, `{"syncId":"11111111-1111-1111-1111-111111111111"}`},
				{"update empty sync ID", organization.UpdateRequest{SyncID: model.Value("")}, `{"syncId":""}`},
				{"update null sync ID", organization.UpdateRequest{SyncID: model.Null[string]()}, `{"syncId":null}`},
				{"one field", organization.UpdateRequest{DirectorPosition: model.Value("Руководитель")}, `{"directorPosition":"Руководитель"}`},
				{"create empty and false", organization.CreateRequest{Name: "SDK", Description: model.Value(""), CompanyType: model.Value(model.CompanyType("")), Archived: model.Value(false), Shared: model.Value(false), PayerVAT: model.Value(false), SyncID: model.Value("")}, `{"name":"SDK","description":"","companyType":"","archived":false,"shared":false,"payerVat":false,"syncId":""}`},
				{"update empty and false", organization.UpdateRequest{Name: model.Value(""), Email: model.Value(""), INN: model.Value(""), TrackingContractNumber: model.Value(""), CompanyType: model.Value(model.CompanyType("")), Archived: model.Value(false), Shared: model.Value(false), PayerVAT: model.Value(false)}, `{"name":"","email":"","inn":"","trackingContractNumber":"","companyType":"","archived":false,"shared":false,"payerVat":false}`},
				{"create true", organization.CreateRequest{Name: "SDK", Shared: model.Value(true), PayerVAT: model.Value(true)}, `{"name":"SDK","shared":true,"payerVat":true}`},
				{"update true", organization.UpdateRequest{Archived: model.Value(true), Shared: model.Value(true), PayerVAT: model.Value(true)}, `{"archived":true,"shared":true,"payerVat":true}`},
				{"create nulls", organization.CreateRequest{Name: "SDK", Email: model.Null[string](), CompanyType: model.Null[model.CompanyType](), PayerVAT: model.Null[bool](), CertificateDate: model.Null[model.DateTime](), TrackingContractDate: model.Null[model.DateTime](), ActualAddress: model.Null[string](), LegalAddress: model.Null[string](), ActualAddressFull: model.Null[model.AddressRequest](), LegalAddressFull: model.Null[model.AddressRequest](), Attributes: model.Null[[]model.Attribute](), Owner: model.Null[model.Reference](), Group: model.Null[model.Reference](), BonusProgram: model.Null[model.Reference](), SyncID: model.Null[string]()}, `{"name":"SDK","email":null,"companyType":null,"payerVat":null,"certificateDate":null,"trackingContractDate":null,"actualAddress":null,"legalAddress":null,"actualAddressFull":null,"legalAddressFull":null,"attributes":null,"owner":null,"group":null,"bonusProgram":null,"syncId":null}`},
				{"update nulls", organization.UpdateRequest{Name: model.Null[string](), CompanyType: model.Null[model.CompanyType](), PayerVAT: model.Null[bool](), CertificateDate: model.Null[model.DateTime](), TrackingContractDate: model.Null[model.DateTime](), ActualAddress: model.Null[string](), LegalAddress: model.Null[string](), ActualAddressFull: model.Null[model.AddressRequest](), LegalAddressFull: model.Null[model.AddressRequest](), Attributes: model.Null[[]model.Attribute](), Owner: model.Null[model.Reference](), Group: model.Null[model.Reference](), BonusProgram: model.Null[model.Reference]()}, `{"name":null,"companyType":null,"payerVat":null,"certificateDate":null,"trackingContractDate":null,"actualAddress":null,"legalAddress":null,"actualAddressFull":null,"legalAddressFull":null,"attributes":null,"owner":null,"group":null,"bonusProgram":null}`},
				{"create empty attributes", organization.CreateRequest{Name: "SDK", Attributes: model.Value([]model.Attribute{})}, `{"name":"SDK","attributes":[]}`},
				{"update empty attributes", organization.UpdateRequest{Attributes: model.Value([]model.Attribute{})}, `{"attributes":[]}`},
				{"create nil attributes", organization.CreateRequest{Name: "SDK", Attributes: model.Value[[]model.Attribute](nil)}, `{"name":"SDK","attributes":[]}`},
				{"update nil attributes", organization.UpdateRequest{Attributes: model.Value[[]model.Attribute](nil)}, `{"attributes":[]}`},
				{"create empty addresses", organization.CreateRequest{Name: "SDK", ActualAddressFull: model.Value(model.AddressRequest{}), LegalAddressFull: model.Value(model.AddressRequest{})}, `{"name":"SDK","actualAddressFull":{},"legalAddressFull":{}}`},
				{"update empty addresses", organization.UpdateRequest{ActualAddressFull: model.Value(model.AddressRequest{}), LegalAddressFull: model.Value(model.AddressRequest{})}, `{"actualAddressFull":{},"legalAddressFull":{}}`},
				{"create partial address", organization.CreateRequest{Name: "SDK", ActualAddress: model.Value(""), ActualAddressFull: model.Value(model.AddressRequest{City: model.Value(""), Country: model.Null[model.Reference](), Comment: model.Null[string]()})}, `{"name":"SDK","actualAddress":"","actualAddressFull":{"city":"","country":null,"comment":null}}`},
				{"update partial address", organization.UpdateRequest{LegalAddress: model.Value(""), LegalAddressFull: model.Value(model.AddressRequest{City: model.Value(""), Region: model.Null[model.Reference](), Street: model.Null[string]()})}, `{"legalAddress":"","legalAddressFull":{"city":"","region":null,"street":null}}`},
				{"unknown company type", organization.UpdateRequest{CompanyType: model.Value(model.CompanyType("futureType"))}, `{"companyType":"futureType"}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := codec.marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, data, []byte(tc.want))
				})
			}
		})
	}
}

func TestRequestCompanyTypeConstants(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				value      model.CompanyType
				wantCreate string
				wantUpdate string
			}{
				{model.CompanyTypeLegal, `{"name":"SDK","companyType":"legal"}`, `{"companyType":"legal"}`},
				{model.CompanyTypeEntrepreneur, `{"name":"SDK","companyType":"entrepreneur"}`, `{"companyType":"entrepreneur"}`},
				{model.CompanyTypeIndividual, `{"name":"SDK","companyType":"individual"}`, `{"companyType":"individual"}`},
				{model.CompanyTypeLegalUZ, `{"name":"SDK","companyType":"legalUZ"}`, `{"companyType":"legalUZ"}`},
				{model.CompanyTypeEntrepreneurUZ, `{"name":"SDK","companyType":"entrepreneurUZ"}`, `{"companyType":"entrepreneurUZ"}`},
				{model.CompanyTypeIndividualUZ, `{"name":"SDK","companyType":"individualUZ"}`, `{"companyType":"individualUZ"}`},
				{model.CompanyTypeLegalKZ, `{"name":"SDK","companyType":"legalKZ"}`, `{"companyType":"legalKZ"}`},
				{model.CompanyTypeEntrepreneurKZ, `{"name":"SDK","companyType":"entrepreneurKZ"}`, `{"companyType":"entrepreneurKZ"}`},
				{model.CompanyTypeIndividualKZ, `{"name":"SDK","companyType":"individualKZ"}`, `{"companyType":"individualKZ"}`},
			} {
				t.Run(string(tc.value), func(t *testing.T) {
					for _, request := range []struct {
						name  string
						input any
						want  string
					}{
						{"create", organization.CreateRequest{Name: "SDK", CompanyType: model.Value(tc.value)}, tc.wantCreate},
						{"update", organization.UpdateRequest{CompanyType: model.Value(tc.value)}, tc.wantUpdate},
					} {
						t.Run(request.name, func(t *testing.T) {
							data, err := codec.marshal(request.input)
							if err != nil {
								t.Fatal(err)
							}
							assertJSONEqual(t, data, []byte(request.want))
						})
					}
				})
			}
		})
	}
}

func TestRequestZeroAttributeValue(t *testing.T) {
	attribute, err := model.AttrValue(model.Meta{Href: "https://example.test/attributes/zero"}, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input any
				want  string
			}{
				{"create", organization.CreateRequest{Name: "SDK", Attributes: model.Value([]model.Attribute{attribute})}, `{"name":"SDK","attributes":[{"meta":{"href":"https://example.test/attributes/zero"},"value":0}]}`},
				{"update", organization.UpdateRequest{Attributes: model.Value([]model.Attribute{attribute})}, `{"attributes":[{"meta":{"href":"https://example.test/attributes/zero"},"value":0}]}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := codec.marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, data, []byte(tc.want))
				})
			}
		})
	}
}

// Пример создания юрлица источника: подписи и печать — filename и content (Base64), НДС для
// авансов и настройки НДС.
func TestRequestVatAndImages(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G'}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(organization.UpdateRequest{
				CompanyVat:        model.Value(organization.CompanyVatRequest{UseCompanyVat: model.Value(true), DefaultCompanyVat: model.Value(0), DefaultCompanyVatEnabled: model.Value(false)}),
				AdvancePaymentVat: model.Value(20.0),
				DirectorSign:      model.Value(organization.ImageRequest{Filename: "directorSignTest.png", Content: png}),
				ChiefAccountSign:  model.Null[organization.ImageRequest](),
				Stamp:             model.Value(organization.ImageRequest{Filename: "stampTest.png", Content: []byte{}}),
			})
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, []byte(`{"companyVat__ru":{"useCompanyVat":true,"defaultCompanyVat":0,"defaultCompanyVatEnabled":false},"advancePaymentVat":20,"directorSign":{"filename":"directorSignTest.png","content":"iVBORw=="},"chiefAccountSign":null,"stamp":{"filename":"stampTest.png","content":""}}`))
		})
	}
}

// Ответ источника с подписями: meta файла, размер, миниатюра; «без НДС» — без defaultCompanyVat.
func TestVatAndImagesResponse(t *testing.T) {
	data := []byte(`{"payerVat":true,"companyVat__ru":{"useCompanyVat":true,"defaultCompanyVatEnabled":false},"advancePaymentVat":10.5,` +
		`"directorSign":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/download/7da0feb5-e110-4021-a49f-35db4ae75f13","mediaType":"application/octet-stream"},` +
		`"title":"directorSignTest","filename":"directorSignTest.png","size":70,"updated":"2020-09-23 07:37:26.417",` +
		`"miniature":{"href":"https://api.moysklad.ru/api/remap/1.2/download/7da0feb5-e110-4021-a49f-35db4ae75f13?miniature=true","mediaType":"image/png","downloadHref":"https://miniature-prod.moysklad.ru/x"}}}`)
	for _, codec := range codecs {
		var got organization.Organization
		if err := codec.unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		vat, sign := got.CompanyVat, got.DirectorSign
		if vat == nil || !vat.UseCompanyVat || vat.DefaultCompanyVat != nil || vat.DefaultCompanyVatEnabled == nil || *vat.DefaultCompanyVatEnabled ||
			got.AdvancePaymentVat != 10.5 || sign == nil || sign.Size != 70 || sign.Filename != "directorSignTest.png" || sign.Title != "directorSignTest" ||
			sign.Meta.Href == "" || sign.Miniature == nil || sign.Miniature.DownloadHref == "" || sign.Updated.IsZero() || got.Stamp != nil {
			t.Fatalf("%s: organization = %+v", codec.name, got)
		}
	}
}

func TestRequestReadOnlyAndDeferredFields(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[organization.CreateRequest](), reflect.TypeFor[organization.UpdateRequest]()} {
		for _, field := range reflect.VisibleFields(typ) {
			tag := strings.Split(field.Tag.Get("json"), ",")[0]
			for _, forbidden := range []string{"meta", "id", "accountId", "created", "updated", "bonusPoints", "isEgaisEnable", "accounts", "mod__requisites__uz", "mod__requisites__kz", "fsrarId", "utmUrl"} {
				if tag == forbidden {
					t.Errorf("%s contains %s", typ.Name(), forbidden)
				}
			}
		}
	}
}

func fullUpdateRequest(t *testing.T) organization.UpdateRequest {
	t.Helper()
	return organization.UpdateRequest{
		Name: model.Value("SDK Organization updated"), Description: model.Value(""), Code: model.Value("0007"), ExternalCode: model.Value("sdk-org"),
		Archived: model.Value(false), Shared: model.Value(true), CompanyType: model.Value(model.CompanyTypeEntrepreneur),
		Email: model.Value("org@example.test"), Phone: model.Value("+7 123"), Fax: model.Value("456"),
		LegalTitle: model.Value("ИП Петров"), LegalFirstName: model.Value("Иван"), LegalLastName: model.Value("Петров"), LegalMiddleName: model.Value("Иванович"),
		INN: model.Value("001234567890"), KPP: model.Value("001234567"), OGRN: model.Value("0001"), OGRNIP: model.Value("0002"), OKPO: model.Value("0003"),
		CertificateNumber: model.Value("007"), CertificateDate: model.Value(date(t, "2020-02-03T00:00:00+03:00")),
		ActualAddress: model.Value("Строка фактического адреса"), LegalAddress: model.Value("Строка юридического адреса"),
		ActualAddressFull: model.Value(fullAddressRequest("Москва")), LegalAddressFull: model.Value(fullAddressRequest("Тула")),
		Attributes:   model.Value(requestAttributes()),
		Owner:        model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/employee/owner"}}),
		Group:        model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/group/group"}}),
		BonusProgram: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/bonusprogram/bonus"}}),
		SyncID:       model.Value("11111111-1111-1111-1111-111111111111"),
		PayerVAT:     model.Value(true), Director: model.Value("Иван Петров"), DirectorPosition: model.Value("Директор"), ChiefAccountant: model.Value("Мария Петрова"),
		TrackingContractNumber: model.Value("00042"), TrackingContractDate: model.Value(date(t, "2025-01-02T03:04:05+03:00")),
	}
}

func fullCreateRequest(t *testing.T) organization.CreateRequest {
	t.Helper()
	return organization.CreateRequest{
		Name: "SDK Organization", Description: model.Value("Описание"), Code: model.Value("0007"), ExternalCode: model.Value("sdk-org"),
		Archived: model.Value(true), Shared: model.Value(false), CompanyType: model.Value(model.CompanyType("futureKind")),
		Email: model.Value("org@example.test"), Phone: model.Value("+7 123"), Fax: model.Value("456"),
		LegalTitle: model.Value("ИП Петров"), LegalFirstName: model.Value("Иван"), LegalLastName: model.Value("Петров"), LegalMiddleName: model.Value("Иванович"),
		INN: model.Value("001234567890"), KPP: model.Value("001234567"), OGRN: model.Value("0001"), OGRNIP: model.Value("0002"), OKPO: model.Value("0003"),
		CertificateNumber: model.Value("007"), CertificateDate: model.Value(date(t, "2020-02-03T00:00:00+03:00")),
		ActualAddress: model.Value("Строка фактического адреса"), LegalAddress: model.Value("Строка юридического адреса"),
		ActualAddressFull: model.Value(fullAddressRequest("Москва")), LegalAddressFull: model.Value(fullAddressRequest("Тула")),
		Attributes:   model.Value(requestAttributes()),
		Owner:        model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/employee/owner"}}),
		Group:        model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/group/group"}}),
		BonusProgram: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/bonusprogram/bonus"}}),
		SyncID:       model.Value("00000000-0000-4000-8000-000000000007"),
		PayerVAT:     model.Value(false), Director: model.Value("Иван Петров"), DirectorPosition: model.Value("Директор"), ChiefAccountant: model.Value("Мария Петрова"),
		TrackingContractNumber: model.Value("00042"), TrackingContractDate: model.Value(date(t, "2025-01-02T03:04:05+03:00")),
	}
}

func fullAddressRequest(city string) model.AddressRequest {
	return model.AddressRequest{
		AddInfo: model.Value("Корпус"), Apartment: model.Value("007"), City: model.Value(city), Comment: model.Value("Комментарий"), FIASCodeRU: model.Value("000001"),
		Country: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/country/ru"}}),
		House:   model.Value("01"), PostalCode: model.Value("001234"), Region: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/region/77"}}), Street: model.Value("Тверская"),
	}
}

func requestAttributes() []model.Attribute {
	return []model.Attribute{
		{Meta: model.Meta{Href: "https://example.test/attributes/long"}, Value: model.Value(jsontext.Value("9007199254740993"))},
		model.AttrNull(model.Meta{Href: "https://example.test/attributes/reset"}),
		model.AttrFileNull(model.Meta{Href: "https://example.test/attributes/file"}),
	}
}

func organizationGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/custom/organization/" + name + ".golden.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}
