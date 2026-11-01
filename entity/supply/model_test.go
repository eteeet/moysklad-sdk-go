package supply_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/supply"
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
				want       supply.Supply
			}{
				{"upstream", "../../testdata/upstream/supply.json", supply.Supply{
					ID:           "a3a404e4-2e5d-11e6-8a84-bae5000000fd",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "00001",
					Description:  "Приемка товаров ATK-179",
					Code:         "ATK179-SUP-001",
					ExternalCode: "0ULmxwN1jHJwT9nYPawwO2",
					Created:      documenttest.DateTime("2016-06-10 09:26:00"),
					Updated:      documenttest.DateTime("2016-06-27 10:48:47"),
					Moment:       documenttest.DateTime("2016-06-10 09:26:00"),
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/supply/metadata/states/918e5abd-3f66-11e6-8a84-bae500000083",
							Type:      model.MetaType("state"),
							MediaType: "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/supply/a3a404e4-2e5d-11e6-8a84-bae5000000fd",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/supply/metadata",
						Type:         model.MetaType("supply"),
						MediaType:    "application/json",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/supply/metadata/attributes/a31685ae-3f62-11e6-8a84-bae50000007b",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "a31685ae-3f62-11e6-8a84-bae50000007b",
							Name:  "AttributeName1",
							Type:  "time",
							Value: model.Value(jsontext.Value("\"2016-06-10 09:26:00\"")),
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
					Sum: float64(350000),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/faf44002-2e58-11e6-8a84-bae500000053",
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/faf44002-2e58-11e6-8a84-bae500000053/accounts/147c3231-32ca-11e6-8a84-bae500000005",
							Type:      model.MetaType("account"),
							MediaType: "application/json",
						},
					}),
					Contract: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/contract/91287d9a-41bc-11e6-8a84-bae500000082",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/contract/metadata",
							Type:         model.MetaType("contract"),
							MediaType:    "application/json",
						},
					}),
					VATSum: float64(58333),
					Positions: documenttest.Pointer(model.Collection[Position]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/supply/a3a404e4-2e5d-11e6-8a84-bae5000000fd/positions",
								Type:      model.MetaType("supplyposition"),
								MediaType: "application/json",
							},
							Size:   1,
							Limit:  1000,
							Offset: 0,
						},
						Rows: []Position{
							Position{
								Item: positions.Item{
									Core: positions.Core{
										Meta: model.Meta{
											Href:      "https://api.moysklad.ru/api/remap/1.2/entity/supply/a3a404e4-2e5d-11e6-8a84-bae5000000fd/positions/9cfa8b56-41bf-11e6-8a84-bae500000014",
											Type:      model.MetaType("supplyposition"),
											MediaType: "application/json",
										},
										ID:        "9cfa8b56-41bf-11e6-8a84-bae500000014",
										AccountID: "f976ed28-2e58-11e6-8a84-bae500000001",
									},
									Assortment: model.Reference{
										Meta: model.Meta{
											Href:         "https://api.moysklad.ru/api/remap/1.2/entity/product/20485cfd-2e62-11e6-8a84-bae500000112",
											MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/product/metadata",
											Type:         model.MetaType("product"),
											MediaType:    "application/json",
										},
									},
									Quantity: float64(10),
								},
								Price:         float64(35000),
								Discount:      float64(0),
								VAT:           18,
								VATEnabled:    true,
								Overhead:      float64(5000),
								Things:        []string{},
								TrackingCodes: []positions.TrackingCode{{CIS: "010460406000301021N4N57RTCBUZQZ", Type: "trackingcode", TrackingCodes: []positions.TrackingCode{}}},
							},
						},
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
					Files: documenttest.Pointer(model.Collection[model.Reference]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/supply/a3a404e4-2e5d-11e6-8a84-bae5000000fd/files",
								Type:      model.MetaType("files"),
								MediaType: "application/json",
							},
							Size:   1,
							Limit:  1000,
							Offset: 0,
						},
						Rows: []model.Reference{
							model.Reference{
								Meta: model.Meta{
									Href:         "https://api.moysklad.ru/api/remap/1.2/entity/supply/a3a404e4-2e5d-11e6-8a84-bae5000000fd/files/f2728180-6afd-4d37-8a13-f3b48069bbb6",
									Type:         model.MetaType("files"),
									MediaType:    "application/json",
									DownloadHref: "https://api.moysklad.ru/api/remap/1.2/download/f2728180-6afd-4d37-8a13-f3b48069bbb6",
								},
							},
						},
					}),
					IncomingDate:   documenttest.DateTime("2016-06-04 00:00:00"),
					IncomingNumber: "INV-2016-0604",
					PurchaseOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/22b4caaa-3f74-11e6-8a84-bae500000069",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata",
							Type:         model.MetaType("purchaseorder"),
							MediaType:    "application/json",
						},
					}),
					FactureIn: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/facturein/8cc4ef04-f831-11e5-7a69-971500188b19",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/facturein/metadata",
							Type:         model.MetaType("facturein"),
							MediaType:    "application/json",
						},
					}),
					ProductionTask: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/productiontask/5aa4ef04-f831-11e5-7a69-971500188b19",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/productiontask/metadata",
							Type:         model.MetaType("productiontask"),
							MediaType:    "application/json",
						},
					}),
					InvoicesIn: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/6bb4ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/metadata",
								Type:         model.MetaType("invoicein"),
								MediaType:    "application/json",
							},
						},
					},
					Returns: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/metadata",
								Type:         model.MetaType("purchasereturn"),
								MediaType:    "application/json",
							},
						},
					},
					Overhead: documenttest.Pointer(model.DocumentOverhead{
						Sum:          float64(50000),
						Distribution: model.OverheadDistribution("weight"),
					}),
				}},
				{"unexpanded complete", "../../testdata/custom/supply/read.json", supply.Supply{
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
						Href: "https://example.test/entity/supply/read-id",
						Type: model.MetaType("supply"),
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
								Type: model.MetaType("supplyposition"),
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
					Files: documenttest.Pointer(model.Collection[model.Reference]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href: "https://example.test/files",
								Type: model.MetaType("files"),
							},
							Size:   2,
							Limit:  1000,
							Offset: 0,
						},
					}),
					IncomingDate:   documenttest.DateTime("2026-11-01 12:34:56.000"),
					IncomingNumber: "00007",
					PurchaseOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/purchaseOrder",
							Type: model.MetaType("purchaseOrder"),
						},
					}),
					FactureIn: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/factureIn",
							Type: model.MetaType("factureIn"),
						},
					}),
					ProductionTask: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/productionTask",
							Type: model.MetaType("productionTask"),
						},
					}),
					InvoicesIn: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/invoicesIn",
								Type: model.MetaType("invoicesIn"),
							},
						},
					},
					Returns: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/returns",
								Type: model.MetaType("returns"),
							},
						},
					},
					Overhead: documenttest.Pointer(model.DocumentOverhead{
						Sum:          float64(12.25),
						Distribution: model.OverheadDistribution("volume"),
					}),
				}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := os.ReadFile(tc.path)
					if err != nil {
						t.Fatal(err)
					}
					var got supply.Supply
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
				var got supply.Supply
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
