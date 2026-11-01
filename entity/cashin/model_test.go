package cashin_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestModelJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name, path string
				want       cashin.CashIn
			}{
				{"upstream", "../../testdata/upstream/cash_in.json", cashin.CashIn{
					ID:           "b0ec3c98-3fac-11e6-8a84-bae500000131",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "00001",
					Description:  "Приходный Ордер.",
					Code:         "E7vH1TaKii9bkuw2pekD22",
					ExternalCode: "E7vH1TaKii9bkuw2pekD22",
					SyncID:       "b0ec3c98-3fac-11e6-8a84-bae500000131",
					Created:      documenttest.DateTime("2016-08-25 19:55:00"),
					Updated:      documenttest.DateTime("2016-07-04 08:55:58"),
					Deleted:      documenttest.DateTime("2026-03-25 17:44:31.866"),
					Moment:       documenttest.DateTime("2016-07-04 08:54:00"),
					Applicable:   true,
					Shared:       true,
					Printed:      true,
					Published:    true,
					Organization: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/organization/fae3561a-2e58-11e6-8a84-bae50000004e",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/organization/metadata",
							Type:         model.MetaType("organization"),
							MediaType:    "application/json",
						},
					}),
					Owner: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/employee/faba7f37-2e58-11e6-8a84-bae500000028",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/employee/metadata",
							Type:         model.MetaType("employee"),
							MediaType:    "application/json",
						},
					}),
					Group: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/group/f97aa1fb-2e58-11e6-8a84-bae500000002",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/group/metadata",
							Type:         model.MetaType("group"),
							MediaType:    "application/json",
						},
					}),
					Project: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/project/722e39f0-313e-11e6-8a84-bae500000008",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/project/metadata",
							Type:         model.MetaType("project"),
							MediaType:    "application/json",
						},
					}),
					State: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/cashin/metadata/states/aaa6c45a-3fac-11e6-8a84-bae50000012b",
							Type:      model.MetaType("state"),
							MediaType: "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/cashin/b0ec3c98-3fac-11e6-8a84-bae500000131",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/cashin/metadata",
						Type:         model.MetaType("cashin"),
						MediaType:    "application/json",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata/attributes/0d129ff5-2c8c-11e6-8a84-bae5000000f3",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "0d129ff5-2c8c-11e6-8a84-bae5000000f3",
							Name:  "AttributeName1",
							Type:  "double",
							Value: model.Value(jsontext.Value("0.75")),
						},
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata/attributes/0d12a9a5-2c8c-11e6-8a84-bae5000000f4",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "0d12a9a5-2c8c-11e6-8a84-bae5000000f4",
							Name:  "AttributeName2",
							Type:  "time",
							Value: model.Value(jsontext.Value("\"2016-06-07 12:52:33\"")),
						},
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata/attributes/0d12b1e7-2c8c-11e6-8a84-bae5000000f5",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "0d12b1e7-2c8c-11e6-8a84-bae5000000f5",
							Name:  "AttributeName3",
							Type:  "boolean",
							Value: model.Value(jsontext.Value("false")),
						},
					},
					Rate: documenttest.Pointer(model.DocumentRate{
						Currency: documenttest.Pointer(model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/currency/cdbc62de-3f68-11e6-8a84-bae500000050",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/currency/metadata",
								Type:         model.MetaType("currency"),
								MediaType:    "application/json",
							},
						}),
						Value: documenttest.Pointer(float64(63)),
					}),
					Sum: float64(3189),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/1489a08d-32ca-11e6-8a84-bae50000000d",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata",
							Type:         model.MetaType("counterparty"),
							MediaType:    "application/json",
						},
					}),
					Contract: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/contract/2b34d43f-3f52-11e6-8a84-bae500000066",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/contract/metadata",
							Type:         model.MetaType("contract"),
							MediaType:    "application/json",
						},
					}),
					VATSum:         float64(3189),
					PaymentPurpose: "Оплата заказа",
					Operations: []model.FinanceOperation{
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/22591b08-2859-11f1-0a83-1f2500000247",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata",
								Type:         model.MetaType("customerorder"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(200000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/metadata",
								Type:         model.MetaType("purchasereturn"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(45000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/demand/7944ef04-047e-11e9-ac12-000a00000002",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata",
								Type:         model.MetaType("demand"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(150000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/e3f4a5b6-1234-1234-1234-e3f4a5b6c7d8",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/metadata",
								Type:         model.MetaType("invoiceout"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(110000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/commissionreportin/7944ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/commissionreportin/metadata",
								Type:         model.MetaType("commissionreportin"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(65000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/retailshift/90ba347d-6b8b-11e6-8a84-bae5000000bd",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/retailshift/metadata",
								Type:         model.MetaType("retailshift"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(80000),
						},
					},
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/10b0654b-2859-11f1-0a83-1f2500000244",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/metadata",
							Type:         model.MetaType("saleschannel"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#saleschannel/edit?id=10b0654b-2859-11f1-0a83-1f2500000244",
						},
					}),
					FactureOut: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/factureout/02eb1458-aa6e-11e6-8a84-bc5200000076",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/factureout/metadata",
							Type:         model.MetaType("factureout"),
							MediaType:    "application/json",
						},
					}),
				}},
				{"unexpanded complete", "../../testdata/custom/cashin/read.json", cashin.CashIn{
					ID:           "read-id",
					AccountID:    "read-accountId",
					Name:         "request-name",
					Description:  "",
					Code:         "",
					ExternalCode: "request-externalCode",
					SyncID:       "request-syncId",
					Created:      documenttest.DateTime("2026-11-01 10:20:30.000"),
					Updated:      documenttest.DateTime("2026-11-01 10:20:30.000"),
					Deleted:      documenttest.DateTime("2026-11-01 10:20:30.000"),
					Moment:       documenttest.DateTime("2026-11-01 12:34:56.000"),
					Applicable:   false,
					Shared:       false,
					Printed:      true,
					Published:    true,
					Organization: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/organization",
							Type: model.MetaType("organization"),
						},
					}),
					Owner: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/owner",
							Type: model.MetaType("owner"),
						},
					}),
					Group: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/group",
							Type: model.MetaType("group"),
						},
					}),
					Project: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/project",
							Type: model.MetaType("project"),
						},
					}),
					State: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/state",
							Type: model.MetaType("state"),
						},
					}),
					Meta: model.Meta{
						Href: "https://example.test/entity/cashin/read-id",
						Type: model.MetaType("cashin"),
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href: "https://example.test/attributes/long",
							},
							Type:  "long",
							Value: model.Value(jsontext.Value("9007199254740993")),
						},
					},
					Rate: documenttest.Pointer(model.DocumentRate{
						Currency: documenttest.Pointer(model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/currency",
								Type: model.MetaType("currency"),
							},
						}),
					}),
					Sum: float64(123.25),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/agent",
							Type: model.MetaType("agent"),
						},
					}),
					Contract: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/contract",
							Type: model.MetaType("contract"),
						},
					}),
					VATSum:         float64(0),
					PaymentPurpose: "request-paymentPurpose",
					Operations: []model.FinanceOperation{
						model.FinanceOperation{
							Meta: model.Meta{
								Href: "https://example.test/operation",
								Type: model.MetaType("futureoperation"),
							},
							LinkedSum: float64(1.25),
						},
					},
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/salesChannel",
							Type: model.MetaType("salesChannel"),
						},
					}),
					FactureOut: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/factureOut",
							Type: model.MetaType("factureOut"),
						},
					}),
				}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := os.ReadFile(tc.path)
					if err != nil {
						t.Fatal(err)
					}
					var got cashin.CashIn
					if err := codec.Unmarshal(data, &got); err != nil {
						t.Fatal(err)
					}
					documenttest.NormalizeAttributeJSON(t, got.Attributes)
					documenttest.NormalizeAttributeJSON(t, tc.want.Attributes)
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("model = %+v, want %+v", got, tc.want)
					}
				})
			}
			for _, data := range []string{`{}`, `{"created":null,"updated":null,"deleted":null,"rate":null,"organization":null,"futureField":1}`} {
				var got cashin.CashIn
				if err := codec.Unmarshal([]byte(data), &got); err != nil {
					t.Fatal(err)
				}
				if got.ID != "" || got.Organization != nil || got.Rate != nil || !got.Created.IsZero() || !got.Updated.IsZero() || !got.Deleted.IsZero() {
					t.Fatalf("invented fields: %+v", got)
				}
			}
		})
	}
}
