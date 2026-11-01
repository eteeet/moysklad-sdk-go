package invoiceout_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
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
				want       invoiceout.InvoiceOut
			}{
				{"upstream", "../../testdata/upstream/invoice_out.json", invoiceout.InvoiceOut{
					ID:           "6afee657-0888-11e6-9464-e4de00000037",
					AccountID:    "84e60e93-f504-11e5-8a84-bae500000008",
					Name:         "change name",
					Description:  "Комментарий к счету",
					Code:         "change_code",
					ExternalCode: "666777666",
					SyncID:       "734a9e26-45a2-4ead-849c-e144daeb854d",
					Created:      documenttest.DateTime("2016-08-25 19:55:00"),
					Updated:      documenttest.DateTime("2016-04-22 16:21:19"),
					Moment:       documenttest.DateTime("2016-04-19 13:50:24"),
					Applicable:   false,
					Shared:       false,
					Printed:      true,
					Published:    true,
					Organization: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/organization/850c8195-f504-11e5-8a84-bae50000015e",
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
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/6afee657-0888-11e6-9464-e4de00000037",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/metadata",
						Type:         model.MetaType("invoiceout"),
						MediaType:    "application/json",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/metadata/attributes/0cd74e1e-2e59-11e6-8a84-bae50000008a",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "51d191f4-0887-11e6-9464-e4de00000079",
							Name:  "AttributeName1",
							Type:  "string",
							Value: model.Value(jsontext.Value("\"new_val\"")),
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
					Sum: float64(5840),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/850f1667-f504-11e5-8a84-bae500000163",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata",
							Type:         model.MetaType("counterparty"),
							MediaType:    "application/json",
						},
					}),
					OrganizationAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/account/850c8af2-f504-11e5-8a84-bae50000015f",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
						},
					}),
					AgentAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/account/850f1eab-f504-11e5-8a84-bae500000164",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
						},
					}),
					Positions: documenttest.Pointer(model.Collection[Position]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/invoiceout/6afee657-0888-11e6-9464-e4de00000037/positions",
								Type:      model.MetaType("invoiceposition"),
								MediaType: "application/json",
							},
							Size:   2,
							Limit:  1000,
							Offset: 0,
						},
						Rows: []Position{},
					}),
					PayedSum: float64(0),
					Store: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/store/850ee995-f504-11e5-8a84-bae500000160",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/store/metadata",
							Type:         model.MetaType("store"),
							MediaType:    "application/json",
						},
					}),
					ShippedSum: float64(0),
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/56446e7f-3633-11ec-ac13-000d00000000",
							Type:      model.MetaType("saleschannel"),
							MediaType: "application/json",
						},
					}),
				}},
				{"unexpanded complete", "../../testdata/custom/invoiceout/read.json", invoiceout.InvoiceOut{
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
						Href: "https://example.test/entity/invoiceout/read-id",
						Type: model.MetaType("invoiceout"),
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
								Type: model.MetaType("invoiceoutposition"),
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
					CustomerOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/customerOrder",
							Type: model.MetaType("customerOrder"),
						},
					}),
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/salesChannel",
							Type: model.MetaType("salesChannel"),
						},
					}),
					Demands: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/demands",
								Type: model.MetaType("demands"),
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
					var got invoiceout.InvoiceOut
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
				var got invoiceout.InvoiceOut
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
