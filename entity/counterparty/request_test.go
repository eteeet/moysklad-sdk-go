package counterparty_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCreateRequestMinimal(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(counterparty.CreateRequest{Name: "SDK Counterparty"})
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, []byte(`{"name":"SDK Counterparty"}`))
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
			assertJSONEqual(t, data, golden(t, "create"))
		})
	}
}

func fullCreateRequest(t *testing.T) counterparty.CreateRequest {
	return counterparty.CreateRequest{
		Name: "SDK Counterparty", Description: model.Value("Комментарий"), Code: model.Value("001"), ExternalCode: model.Value("sdk-cp"),
		Archived: model.Value(false), Shared: model.Value(true), CompanyType: model.Value(model.CompanyTypeEntrepreneur),
		Email: model.Value("sdk@example.com"), Phone: model.Value("+7 495 0000000"), Fax: model.Value("+7 495 0000001"),
		LegalTitle: model.Value("ИП Иванов"), LegalFirstName: model.Value("Иван"), LegalLastName: model.Value("Иванов"), LegalMiddleName: model.Value("Иванович"),
		INN: model.Value("001234567890"), KPP: model.Value("001234567"), OGRN: model.Value("0012345678901"), OGRNIP: model.Value("001234567890123"), OKPO: model.Value("00123456"),
		CertificateNumber: model.Value("0001"), CertificateDate: model.Value(date(t, "2024-01-01T21:00:00Z")),
		BirthDate: model.Value(date(t, "2000-01-01T21:00:00Z")), Sex: model.Value("future-sex"),
		ActualAddress: model.Value("Строковый фактический"), LegalAddress: model.Value("Строковый юридический"),
		ActualAddressFull: model.Value(fullAddressRequest()), LegalAddressFull: model.Value(model.AddressRequest{City: model.Value("Казань")}),
		Attributes: model.Value(requestAttributes()), Tags: model.Value([]string{"покупатель", "поставщик"}), DiscountCardNumber: model.Value("00042"),
		Owner: model.Value(requestRef("employee")), Group: model.Value(requestRef("group")), BonusProgram: model.Value(requestRef("bonusprogram")), State: model.Value(requestRef("state")),
		PriceType: model.Value(model.PriceType{Meta: model.Meta{Href: "https://example.com/pricetype"}, ID: "price-id", Name: "Оптовая", ExternalCode: "price-code"}),
		SyncID:    model.Value("11111111-1111-1111-1111-111111111111"),
	}
}

func fullAddressRequest() model.AddressRequest {
	return model.AddressRequest{
		AddInfo: model.Value("подъезд 1"), Apartment: model.Value("03"), City: model.Value("Москва"), Comment: model.Value("домофон"),
		FIASCodeRU: model.Value("000001"), Country: model.Value(requestRef("country")), House: model.Value("02"), PostalCode: model.Value("001234"),
		Region: model.Value(requestRef("region")), Street: model.Value("Лесная"),
	}
}

func requestRef(typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://example.com/" + typ}}
}

func requestAttributes() []model.Attribute {
	return []model.Attribute{
		{Meta: model.Meta{Href: "https://example.com/attribute"}, Value: model.Value(jsontext.Value(`9007199254740993`))},
		model.AttrNull(model.Meta{Href: "https://example.com/reset"}),
		model.AttrFileNull(model.Meta{Href: "https://example.com/file"}),
	}
}

func golden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/custom/counterparty/" + name + ".golden.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestUpdateRequestFull(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(fullUpdateRequest(t))
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, golden(t, "update"))
		})
	}
}

func fullUpdateRequest(t *testing.T) counterparty.UpdateRequest {
	return counterparty.UpdateRequest{
		Name: model.Value("SDK Counterparty updated"), Description: model.Value("Комментарий"), Code: model.Value("001"), ExternalCode: model.Value("sdk-cp"),
		Archived: model.Value(true), Shared: model.Value(false), CompanyType: model.Value(model.CompanyTypeEntrepreneur),
		Email: model.Value("sdk@example.com"), Phone: model.Value("+7 495 0000000"), Fax: model.Value("+7 495 0000001"),
		LegalTitle: model.Value("ИП Иванов"), LegalFirstName: model.Value("Иван"), LegalLastName: model.Value("Иванов"), LegalMiddleName: model.Value("Иванович"),
		INN: model.Value("001234567890"), KPP: model.Value("001234567"), OGRN: model.Value("0012345678901"), OGRNIP: model.Value("001234567890123"), OKPO: model.Value("00123456"),
		CertificateNumber: model.Value("0001"), CertificateDate: model.Value(date(t, "2024-01-01T21:00:00Z")),
		BirthDate: model.Value(date(t, "2000-01-01T21:00:00Z")), Sex: model.Value("future-sex"),
		ActualAddress: model.Value("Строковый фактический"), LegalAddress: model.Value("Строковый юридический"),
		ActualAddressFull: model.Value(fullAddressRequest()), LegalAddressFull: model.Value(model.AddressRequest{City: model.Value("Казань")}),
		Attributes: model.Value(requestAttributes()), Tags: model.Value([]string{"покупатель", "поставщик"}), DiscountCardNumber: model.Value("00042"),
		Owner: model.Value(requestRef("employee")), Group: model.Value(requestRef("group")), BonusProgram: model.Value(requestRef("bonusprogram")), State: model.Value(requestRef("state")),
		PriceType: model.Value(model.PriceType{Meta: model.Meta{Href: "https://example.com/pricetype"}, ID: "price-id", Name: "Оптовая", ExternalCode: "price-code"}),
		SyncID:    model.Value("11111111-1111-1111-1111-111111111111"),
	}
}

func TestRequestExcludedFields(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[counterparty.CreateRequest](), reflect.TypeFor[counterparty.UpdateRequest]()} {
		for _, name := range []string{"Meta", "ID", "AccountID", "Created", "Updated", "BonusPoints", "SalesAmount", "Accounts", "Contactpersons", "ContactPersons", "Notes", "Files", "Discounts"} {
			if _, exists := typ.FieldByName(name); exists {
				t.Errorf("%s contains excluded %s", typ.Name(), name)
			}
		}
	}
}

func TestRequestFieldStates(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input any
				want  string
			}{
				{"zero create", counterparty.CreateRequest{}, `{"name":""}`},
				{"zero update", counterparty.UpdateRequest{}, `{}`},
				{"update sync ID", counterparty.UpdateRequest{SyncID: model.Value("11111111-1111-1111-1111-111111111111")}, `{"syncId":"11111111-1111-1111-1111-111111111111"}`},
				{"update empty sync ID", counterparty.UpdateRequest{SyncID: model.Value("")}, `{"syncId":""}`},
				{"update null sync ID", counterparty.UpdateRequest{SyncID: model.Null[string]()}, `{"syncId":null}`},
				{"partial update", counterparty.UpdateRequest{Phone: model.Value("changed")}, `{"phone":"changed"}`},
				{"create zero attribute", counterparty.CreateRequest{Name: "cp", Attributes: model.Value([]model.Attribute{{Meta: model.Meta{Href: "https://example.com/zero"}, Value: model.Value(jsontext.Value(`0`))}})}, `{"name":"cp","attributes":[{"meta":{"href":"https://example.com/zero"},"value":0}]}`},
				{"update zero attribute", counterparty.UpdateRequest{Attributes: model.Value([]model.Attribute{{Meta: model.Meta{Href: "https://example.com/zero"}, Value: model.Value(jsontext.Value(`0`))}})}, `{"attributes":[{"meta":{"href":"https://example.com/zero"},"value":0}]}`},
				{"create false/empty", counterparty.CreateRequest{Name: "cp", Archived: model.Value(false), Shared: model.Value(false), Description: model.Value(""), Tags: model.Value([]string{}), Attributes: model.Value([]model.Attribute{}), SyncID: model.Value("")}, `{"name":"cp","archived":false,"shared":false,"description":"","tags":[],"attributes":[],"syncId":""}`},
				{"update false/empty", counterparty.UpdateRequest{Name: model.Value(""), Archived: model.Value(false), Shared: model.Value(false), Description: model.Value(""), Tags: model.Value([]string{}), Attributes: model.Value([]model.Attribute{})}, `{"name":"","archived":false,"shared":false,"description":"","tags":[],"attributes":[]}`},
				{"create nil slices", counterparty.CreateRequest{Name: "cp", Tags: model.Value[[]string](nil), Attributes: model.Value[[]model.Attribute](nil)}, `{"name":"cp","tags":[],"attributes":[]}`},
				{"update nil slices", counterparty.UpdateRequest{Tags: model.Value[[]string](nil), Attributes: model.Value[[]model.Attribute](nil)}, `{"tags":[],"attributes":[]}`},
				{"create nulls", counterparty.CreateRequest{Name: "cp", Description: model.Null[string](), CompanyType: model.Null[model.CompanyType](), Archived: model.Null[bool](), Shared: model.Null[bool](), Tags: model.Null[[]string](), Attributes: model.Null[[]model.Attribute](), ActualAddress: model.Null[string](), LegalAddress: model.Null[string](), ActualAddressFull: model.Null[model.AddressRequest](), LegalAddressFull: model.Null[model.AddressRequest](), Owner: model.Null[model.Reference](), Group: model.Null[model.Reference](), BonusProgram: model.Null[model.Reference](), State: model.Null[model.Reference](), PriceType: model.Null[model.PriceType](), CertificateDate: model.Null[model.DateTime](), BirthDate: model.Null[model.DateTime](), SyncID: model.Null[string]()}, `{"name":"cp","description":null,"companyType":null,"archived":null,"shared":null,"tags":null,"attributes":null,"actualAddress":null,"legalAddress":null,"actualAddressFull":null,"legalAddressFull":null,"owner":null,"group":null,"bonusProgram":null,"state":null,"priceType":null,"certificateDate":null,"birthDate":null,"syncId":null}`},
				{"update nulls", counterparty.UpdateRequest{Name: model.Null[string](), CompanyType: model.Null[model.CompanyType](), Archived: model.Null[bool](), Shared: model.Null[bool](), Tags: model.Null[[]string](), Attributes: model.Null[[]model.Attribute](), ActualAddress: model.Null[string](), LegalAddress: model.Null[string](), ActualAddressFull: model.Null[model.AddressRequest](), LegalAddressFull: model.Null[model.AddressRequest](), Owner: model.Null[model.Reference](), Group: model.Null[model.Reference](), BonusProgram: model.Null[model.Reference](), State: model.Null[model.Reference](), PriceType: model.Null[model.PriceType](), CertificateDate: model.Null[model.DateTime](), BirthDate: model.Null[model.DateTime]()}, `{"name":null,"companyType":null,"archived":null,"shared":null,"tags":null,"attributes":null,"actualAddress":null,"legalAddress":null,"actualAddressFull":null,"legalAddressFull":null,"owner":null,"group":null,"bonusProgram":null,"state":null,"priceType":null,"certificateDate":null,"birthDate":null}`},
				{"create address empty object", counterparty.CreateRequest{Name: "cp", ActualAddressFull: model.Value(model.AddressRequest{}), LegalAddressFull: model.Value(model.AddressRequest{})}, `{"name":"cp","actualAddressFull":{},"legalAddressFull":{}}`},
				{"update address empty object", counterparty.UpdateRequest{ActualAddressFull: model.Value(model.AddressRequest{}), LegalAddressFull: model.Value(model.AddressRequest{})}, `{"actualAddressFull":{},"legalAddressFull":{}}`},
				{"create empty companyType", counterparty.CreateRequest{Name: "cp", CompanyType: model.Value(model.CompanyType(""))}, `{"name":"cp","companyType":""}`},
				{"update empty companyType", counterparty.UpdateRequest{CompanyType: model.Value(model.CompanyType(""))}, `{"companyType":""}`},
				{"create unknown companyType", counterparty.CreateRequest{Name: "cp", CompanyType: model.Value(model.CompanyType("future"))}, `{"name":"cp","companyType":"future"}`},
				{"update unknown companyType", counterparty.UpdateRequest{CompanyType: model.Value(model.CompanyType("future"))}, `{"companyType":"future"}`},
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

func TestRequestNestedAddressFieldStates(t *testing.T) {
	emptyAddress := model.AddressRequest{
		AddInfo: model.Value(""), Apartment: model.Value(""), City: model.Value(""), Comment: model.Value(""), FIASCodeRU: model.Value(""),
		Country: model.Null[model.Reference](), House: model.Value(""), PostalCode: model.Value(""), Region: model.Null[model.Reference](), Street: model.Value(""),
	}
	nullAddress := model.AddressRequest{
		AddInfo: model.Null[string](), Apartment: model.Null[string](), City: model.Null[string](), Comment: model.Null[string](), FIASCodeRU: model.Null[string](),
		Country: model.Null[model.Reference](), House: model.Null[string](), PostalCode: model.Null[string](), Region: model.Null[model.Reference](), Street: model.Null[string](),
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				address model.AddressRequest
				want    string
			}{
				{"empty", emptyAddress, `{"addInfo":"","apartment":"","city":"","comment":"","fiasCode__ru":"","country":null,"house":"","postalCode":"","region":null,"street":""}`},
				{"null", nullAddress, `{"addInfo":null,"apartment":null,"city":null,"comment":null,"fiasCode__ru":null,"country":null,"house":null,"postalCode":null,"region":null,"street":null}`},
			} {
				for _, input := range []any{
					counterparty.CreateRequest{Name: "cp", ActualAddress: model.Value(""), ActualAddressFull: model.Value(tc.address), LegalAddressFull: model.Value(tc.address)},
					counterparty.UpdateRequest{ActualAddress: model.Value(""), ActualAddressFull: model.Value(tc.address), LegalAddressFull: model.Value(tc.address)},
				} {
					data, err := codec.marshal(input)
					if err != nil {
						t.Fatal(err)
					}
					prefix := ""
					if _, ok := input.(counterparty.CreateRequest); ok {
						prefix = `"name":"cp",`
					}
					assertJSONEqual(t, data, []byte(`{`+prefix+`"actualAddress":"","actualAddressFull":`+tc.want+`,"legalAddressFull":`+tc.want+`}`))
				}
			}
		})
	}
}

func TestRequestCompanyTypeConstants(t *testing.T) {
	for _, codec := range codecs {
		for _, tc := range []struct {
			value model.CompanyType
			want  string
		}{
			{model.CompanyTypeLegal, "legal"}, {model.CompanyTypeEntrepreneur, "entrepreneur"}, {model.CompanyTypeIndividual, "individual"},
			{model.CompanyTypeLegalUZ, "legalUZ"}, {model.CompanyTypeEntrepreneurUZ, "entrepreneurUZ"}, {model.CompanyTypeIndividualUZ, "individualUZ"},
			{model.CompanyTypeLegalKZ, "legalKZ"}, {model.CompanyTypeEntrepreneurKZ, "entrepreneurKZ"}, {model.CompanyTypeIndividualKZ, "individualKZ"},
		} {
			t.Run(codec.name+"/"+tc.want, func(t *testing.T) {
				created, err := codec.marshal(counterparty.CreateRequest{Name: "cp", CompanyType: model.Value(tc.value)})
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, created, []byte(`{"name":"cp","companyType":"`+tc.want+`"}`))
				updated, err := codec.marshal(counterparty.UpdateRequest{CompanyType: model.Value(tc.value)})
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, updated, []byte(`{"companyType":"`+tc.want+`"}`))
			})
		}
	}
}
