package paymentout_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestModelJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name, path string
				want       paymentout.PaymentOut
			}{
				{"upstream", "../../testdata/upstream/payment_out.json", paymentout.PaymentOut{
					ID:           "64aa2e8c-3f95-11e6-8a84-bae50000010f",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "00001",
					Description:  "Платим контрагенту",
					Code:         "39393",
					ExternalCode: "1PPP5fDei13CwEHXYUpV62",
					SyncID:       "64aa2e8c-3f95-11e6-8a84-bae50000010f",
					Created:      documenttest.DateTime("2016-08-25 19:55:00"),
					Updated:      documenttest.DateTime("2016-07-01 17:09:15"),
					Deleted:      documenttest.DateTime("2016-08-25 19:55:00"),
					Moment:       documenttest.DateTime("2016-07-01 17:07:00"),
					Applicable:   true,
					Shared:       false,
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/paymentout/metadata/state/fb56c504-2e58-11e6-8a84-bae500000069",
							Type:      model.MetaType("state"),
							MediaType: "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/paymentout/64aa2e8c-3f95-11e6-8a84-bae50000010f",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/paymentout/metadata",
						Type:         model.MetaType("paymentout"),
						MediaType:    "application/json",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/paymentout/metadata/attributes/ebc69be9-3f91-11e6-8a84-bae500000104",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "ebc69be9-3f91-11e6-8a84-bae500000104",
							Name:  "Поступил",
							Type:  "boolean",
							Value: model.Value(jsontext.Value("false")),
						},
					},
					Rate: documenttest.Pointer(model.DocumentRate{
						Currency: documenttest.Pointer(model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/currency/faf45b9a-2e58-11e6-8a84-bae500000055",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/currency/metadata",
								Type:         model.MetaType("currency"),
								MediaType:    "application/json",
							},
						}),
						Value: documenttest.Pointer(float64(1)),
					}),
					Sum: float64(12525000),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/1496c09e-32ca-11e6-8a84-bae50000001b",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata",
							Type:         model.MetaType("counterparty"),
							MediaType:    "application/json",
						},
					}),
					OrganizationAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/organization/fae3561a-2e58-11e6-8a84-bae50000004e/accounts/fae39d66-2e58-11e6-8a84-bae50000004f",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
						},
					}),
					AgentAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/1496c09e-32ca-11e6-8a84-bae50000001b/accounts/1496cc36-32ca-11e6-8a84-bae50000001c",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
						},
					}),
					Contract: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/contract/4680d884-3f95-11e6-8a84-bae50000010c",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/contract/metadata",
							Type:         model.MetaType("contract"),
							MediaType:    "application/json",
						},
					}),
					VATSum:         float64(0),
					PaymentPurpose: "Оплата по счету за приемку",
					Operations: []model.FinanceOperation{
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/salesreturn/a8b8e1e3-3f85-11e6-8a84-bae50000008d",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/salesreturn/metadata",
								Type:         model.MetaType("salesreturn"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(50000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/supply/90ba347d-6b8b-11e6-8a84-bae5000000bd",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/supply/metadata",
								Type:         model.MetaType("supply"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(80000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/6bb4ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/metadata",
								Type:         model.MetaType("invoicein"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(120000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata",
								Type:         model.MetaType("purchaseorder"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(95000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/commissionreportout/7944ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/commissionreportout/metadata",
								Type:         model.MetaType("commissionreportout"),
								MediaType:    "application/json",
							},
							LinkedSum: float64(75000),
						},
						model.FinanceOperation{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/payroll/946f8595-17c3-11f1-4095-ec9f00000001",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/payroll/metadata",
								Type:         model.MetaType("payroll"),
								MediaType:    "application/json",
								UUIDHref:     "https://online.moysklad.ru/app/#payroll/edit?id=946f8595-17c3-11f1-4095-ec9f00000001",
							},
							LinkedSum: float64(30000.0),
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
					ExpenseItem: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/expenseitem/1be2350e-0479-11e5-b03a-448a5b426e7e",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/expenseitem/metadata",
							Type:         model.MetaType("expenseitem"),
							MediaType:    "application/json",
						},
					}),
					FactureIn: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/facturein/a1b2c3d4-aa6e-11e6-8a84-bc5200000999",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/facturein/metadata",
							Type:         model.MetaType("facturein"),
							MediaType:    "application/json",
						},
					}),
					NoClosingDocs: false,
					AccrualDate:   documenttest.DateTime("2026-09-23 00:29:14.514"),
				}},
				{"unexpanded complete", "../../testdata/custom/paymentout/read.json", paymentout.PaymentOut{
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
						Href: "https://example.test/entity/paymentout/read-id",
						Type: model.MetaType("paymentout"),
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
					OrganizationAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/organizationAccount",
							Type: model.MetaType("organizationAccount"),
						},
					}),
					AgentAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/agentAccount",
							Type: model.MetaType("agentAccount"),
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
					ExpenseItem: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/expenseItem",
							Type: model.MetaType("expenseItem"),
						},
					}),
					FactureIn: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/factureIn",
							Type: model.MetaType("factureIn"),
						},
					}),
					NoClosingDocs: false,
					AccrualDate:   documenttest.DateTime("2026-11-01 12:34:56.000"),
				}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := os.ReadFile(tc.path)
					if err != nil {
						t.Fatal(err)
					}
					var got paymentout.PaymentOut
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
				var got paymentout.PaymentOut
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

func TestPaymentOutAgentTypes(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, typ := range []model.MetaType{"counterparty", "employee", "organization"} {
				t.Run(string(typ), func(t *testing.T) {
					data := []byte(`{"agent":{"meta":{"href":"https://example.test/recipient","type":"` + string(typ) + `"}}}`)
					var got paymentout.PaymentOut
					if err := codec.Unmarshal(data, &got); err != nil {
						t.Fatal(err)
					}
					ref := model.Reference{Meta: model.Meta{Href: "https://example.test/recipient", Type: typ}}
					if got.Agent == nil || !reflect.DeepEqual(*got.Agent, ref) {
						t.Fatalf("agent = %+v", got.Agent)
					}
					encoded, err := codec.Marshal(paymentout.UpdateRequest{Agent: model.Value(ref)})
					if err != nil {
						t.Fatal(err)
					}
					documenttest.JSONEqual(t, encoded, data)
				})
			}
		})
	}
}
