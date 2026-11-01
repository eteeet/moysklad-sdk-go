package invoicein_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

type Position = positions.Position

func TestModelJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name, path string
				want       invoicein.InvoiceIn
			}{
				{"upstream", "../../testdata/upstream/invoice_in.json", invoicein.InvoiceIn{
					ID:           "36c000ad-3f52-11e6-8a84-bae500000069",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "00002",
					Description:  "Счет поставщика по заказу",
					Code:         "INV-IN-001",
					ExternalCode: "oS8jXsK8hHmxYzn8MzT1d0",
					SyncID:       "734a9e26-45a2-4ead-849c-e144daeb854d",
					Created:      documenttest.DateTime("2016-06-10 09:26:00"),
					Updated:      documenttest.DateTime("2016-07-01 09:08:21"),
					Moment:       documenttest.DateTime("2016-07-01 09:07:00"),
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/metadata/states/918e5abd-3f66-11e6-8a84-bae500000083",
							Type:      model.MetaType("state"),
							MediaType: "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/36c000ad-3f52-11e6-8a84-bae500000069",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/metadata",
						Type:         model.MetaType("invoicein"),
						MediaType:    "application/json",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/metadata/attributes/e97c5ea8-3f52-11e6-8a84-bae50000006f",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "e97c5ea8-3f52-11e6-8a84-bae50000006f",
							Name:  "AttributeName1",
							Type:  "long",
							Value: model.Value(jsontext.Value("7")),
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
					Sum: float64(100000),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/1489a08d-32ca-11e6-8a84-bae50000000d",
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/1489a08d-32ca-11e6-8a84-bae50000000d/accounts/1489ad8d-32ca-11e6-8a84-bae50000000e",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
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
					VATSum: float64(16667),
					Positions: documenttest.Pointer(model.Collection[Position]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/36c000ad-3f52-11e6-8a84-bae500000069/positions",
								Type:      model.MetaType("invoiceposition"),
								MediaType: "application/json",
							},
							Size:   2,
							Limit:  1000,
							Offset: 0,
						},
						Rows: []Position{},
					}),
					VATEnabled:  true,
					VATIncluded: true,
					PayedSum:    float64(0),
					Store: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/store/faf3ff5b-2e58-11e6-8a84-bae500000050",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/store/metadata",
							Type:         model.MetaType("store"),
							MediaType:    "application/json",
						},
					}),
					Payments: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/paymentout/7bb4ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/paymentout/metadata",
								Type:         model.MetaType("paymentout"),
								MediaType:    "application/json",
							},
						},
					},
					PaymentPlannedMoment: documenttest.DateTime("2016-07-07 09:07:00"),
					ShippedSum:           float64(0),
					IncomingDate:         documenttest.DateTime("2016-07-28 00:00:00"),
					IncomingNumber:       "124124",
					PurchaseOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/22b4caaa-3f74-11e6-8a84-bae500000069",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata",
							Type:         model.MetaType("purchaseorder"),
							MediaType:    "application/json",
						},
					}),
					Supplies: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/supply/60c7896d-cf7c-11e7-6a80-332a00000004",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/supply/metadata",
								Type:         model.MetaType("supply"),
								MediaType:    "application/json",
							},
						},
					},
				}},
				{"unexpanded complete", "../../testdata/custom/invoicein/read.json", invoicein.InvoiceIn{
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
						Href: "https://example.test/entity/invoicein/read-id",
						Type: model.MetaType("invoicein"),
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
					Sum: float64(4321.25),
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
					VATSum: float64(4321.25),
					Positions: documenttest.Pointer(model.Collection[Position]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href: "https://example.test/positions",
								Type: model.MetaType("invoiceinposition"),
							},
							Size:   2,
							Limit:  1000,
							Offset: 0,
						},
					}),
					VATEnabled:  false,
					VATIncluded: false,
					PayedSum:    float64(4321.25),
					Store: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/store",
							Type: model.MetaType("store"),
						},
					}),
					Payments: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/payments",
								Type: model.MetaType("payments"),
							},
						},
					},
					PaymentPlannedMoment: documenttest.DateTime("2026-11-01 12:34:56.000"),
					ShippedSum:           float64(4321.25),
					IncomingDate:         documenttest.DateTime("2026-11-01 12:34:56.000"),
					IncomingNumber:       "00007",
					PurchaseOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/purchaseOrder",
							Type: model.MetaType("purchaseOrder"),
						},
					}),
					Supplies: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/supplies",
								Type: model.MetaType("supplies"),
							},
						},
					},
				}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := os.ReadFile(tc.path)
					if err != nil {
						t.Fatal(err)
					}
					var got invoicein.InvoiceIn
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
				var got invoicein.InvoiceIn
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
