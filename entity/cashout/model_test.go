package cashout_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestModelJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name, path string
				want       cashout.CashOut
			}{
				{"upstream", "../../testdata/upstream/cash_out.json", cashout.CashOut{
					ID:           "e446a227-41b1-11e6-8a84-bae500000005",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "0721",
					Description:  "E7vH1TaKii9bkuw2pekD22",
					Code:         "E7vH1TaKii9bkuw2pekD22",
					ExternalCode: "E7vH1TaKii9bkuw2pekD22",
					SyncID:       "e446a227-41b1-11e6-8a84-bae500000005",
					Created:      documenttest.DateTime("2016-08-25 19:55:00"),
					Updated:      documenttest.DateTime("2026-03-25 17:44:31.866"),
					Deleted:      documenttest.DateTime("2026-03-25 17:44:31.866"),
					Moment:       documenttest.DateTime("2012-06-27 16:52:24"),
					Applicable:   false,
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
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata/states/c3faae3e-26ae-11f1-0a83-1f25000000d7",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata",
							Type:         model.MetaType("state"),
							MediaType:    "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/cashout/e446a227-41b1-11e6-8a84-bae500000005",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/cashout/metadata",
						Type:         model.MetaType("cashout"),
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
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/currency/baac25f0-50ac-11e5-300d-c79b00000055",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/currency/metadata",
								Type:         model.MetaType("currency"),
								MediaType:    "application/json",
							},
						}),
					}),
					Sum: float64(80000.4),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/147c1f1b-32ca-11e6-8a84-bae500000004",
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
					VATSum:         float64(80000.4),
					PaymentPurpose: "test",
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
					NoClosingDocs: true,
					AccrualDate:   documenttest.DateTime("2026-09-23 00:29:14.514"),
				}},
				{"unexpanded complete", "../../testdata/custom/cashout/read.json", cashout.CashOut{
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
						Href: "https://example.test/entity/cashout/read-id",
						Type: model.MetaType("cashout"),
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
					var got cashout.CashOut
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
				var got cashout.CashOut
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
