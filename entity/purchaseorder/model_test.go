package purchaseorder_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
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
				want       purchaseorder.PurchaseOrder
			}{
				{"upstream", "../../testdata/upstream/purchase_order.json", purchaseorder.PurchaseOrder{
					ID:           "a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6",
					AccountID:    "c229f794-26ae-11f1-0a83-1f2100000002",
					Name:         "00001",
					Description:  "test purchase order",
					Code:         "3RGKZCHtiIjQnB5IrFWdi0",
					ExternalCode: "3RGKZCHtiIjQnB5IrFWdi0",
					SyncID:       "a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6",
					Created:      documenttest.DateTime("2026-03-25 17:44:31.866"),
					Updated:      documenttest.DateTime("2026-03-25 17:44:31.829"),
					Deleted:      documenttest.DateTime("2026-03-25 17:44:31.866"),
					Moment:       documenttest.DateTime("2026-03-25 17:43:00.000"),
					Applicable:   true,
					Shared:       false,
					Printed:      false,
					Published:    false,
					Organization: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/organization/c340bf60-26ae-11f1-0a83-1f25000000a1",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/organization/metadata",
							Type:         model.MetaType("organization"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#mycompany/edit?id=c340bf60-26ae-11f1-0a83-1f25000000a1",
						},
					}),
					Owner: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/employee/c2f12a89-26ae-11f1-0a83-1f2500000055",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/employee/metadata",
							Type:         model.MetaType("employee"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#employee/edit?id=c2f12a89-26ae-11f1-0a83-1f2500000055",
						},
					}),
					Group: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/group/c22b8a8c-26ae-11f1-0a83-1f2100000003",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/group/metadata",
							Type:         model.MetaType("group"),
							MediaType:    "application/json",
						},
					}),
					Project: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/project/0a1b9d2e-2859-11f1-0a83-1f2500000241",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/project/metadata",
							Type:         model.MetaType("project"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#project/edit?id=0a1b9d2e-2859-11f1-0a83-1f2500000241",
						},
					}),
					State: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata/states/c3faae3e-26ae-11f1-0a83-1f25000000d7",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata",
							Type:         model.MetaType("state"),
							MediaType:    "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata",
						Type:         model.MetaType("purchaseorder"),
						MediaType:    "application/json",
						UUIDHref:     "https://online.moysklad.ru/app/#purchaseorder/edit?id=a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata/attributes/0d129ff5-2c8c-11e6-8a84-bae5000000f3",
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
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata/attributes/0d12a9a5-2c8c-11e6-8a84-bae5000000f4",
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
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/metadata/attributes/0d12b1e7-2c8c-11e6-8a84-bae5000000f5",
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
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/currency/c34f8b1f-26ae-11f1-0a83-1f25000000a9",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/currency/metadata",
								Type:         model.MetaType("currency"),
								MediaType:    "application/json",
								UUIDHref:     "https://online.moysklad.ru/app/#currency/edit?id=c34f8b1f-26ae-11f1-0a83-1f25000000a9",
							},
						}),
					}),
					Sum: float64(0.0),
					Agent: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/c340bf60-26ae-11f1-0a83-1f25000000a1",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/metadata",
							Type:         model.MetaType("counterparty"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#company/edit?id=c340bf60-26ae-11f1-0a83-1f25000000a1",
						},
					}),
					OrganizationAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/agentaccount/c340bf60-26ae-11f1-0a83-1f25000000a1",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/agentaccount/metadata",
							Type:         model.MetaType("agentaccount"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#mycompany/edit?id=c340bf60-26ae-11f1-0a83-1f25000000a1",
						},
					}),
					AgentAccount: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/agentaccount/c340bf60-26ae-11f1-0a83-1f25000000a1",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/agentaccount/metadata",
							Type:         model.MetaType("agentaccount"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#mycompany/edit?id=c340bf60-26ae-11f1-0a83-1f25000000a1",
						},
					}),
					Contract: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/contract/0696d9aa-2859-11f1-0a83-1f250000023d",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/contract/metadata",
							Type:         model.MetaType("contract"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#contract/edit?id=0696d9aa-2859-11f1-0a83-1f250000023d",
						},
					}),
					VATSum: float64(0.0),
					Positions: documenttest.Pointer(model.Collection[Position]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6/positions",
								Type:      model.MetaType("purchaseorderposition"),
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
											Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6/positions/b2c3d4e5-1234-1234-1234-b2c3d4e5f6a7",
											Type:      model.MetaType("purchaseorderposition"),
											MediaType: "application/json",
										},
										ID:        "b2c3d4e5-1234-1234-1234-b2c3d4e5f6a7",
										AccountID: "c229f794-26ae-11f1-0a83-1f2100000002",
									},
									Assortment: model.Reference{
										Meta: model.Meta{
											Href:         "https://api.moysklad.ru/api/remap/1.2/entity/product/d6877ac4-26ae-11f1-0a83-1f2500000149",
											MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/product/metadata",
											Type:         model.MetaType("product"),
											MediaType:    "application/json",
											UUIDHref:     "https://online.moysklad.ru/app/#good/edit?id=d686c1fc-26ae-11f1-0a83-1f2500000147",
										},
									},
									Quantity: float64(5.0),
								},
								Price:      float64(10000.0),
								Discount:   float64(0.0),
								VAT:        0,
								VATEnabled: false,
								Shipped:    float64(0.0),
								InTransit:  float64(0.0),
							},
						},
					}),
					VATEnabled:  true,
					VATIncluded: true,
					PayedSum:    float64(0.0),
					Store: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/store/c34a3943-26ae-11f1-0a83-1f25000000a4",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/store/metadata",
							Type:         model.MetaType("store"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#warehouse/edit?id=c34a3943-26ae-11f1-0a83-1f25000000a4",
						},
					}),
					Payments: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/paymentout/a5b6c7d8-1234-1234-1234-a5b6c7d8e9f0",
								Type:      model.MetaType("paymentout"),
								MediaType: "application/json",
							},
						},
						model.Reference{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/cashout/a5b6c7d8-1234-1234-1234-a5b6c7d8e9f0",
								Type:      model.MetaType("cashout"),
								MediaType: "application/json",
							},
						},
					},
					Files: documenttest.Pointer(model.Collection[model.Reference]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/09ad1032-2bf9-11f1-0a83-1f25000003a7/files",
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
									Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchaseorder/09ad1032-2bf9-11f1-0a83-1f25000003a7/files/c3616f62-c858-4f17-96fd-1a2e4dcb2fb6",
									Type:         model.MetaType("files"),
									MediaType:    "application/json",
									DownloadHref: "https://api.moysklad.ru/api/remap/1.2/download/c3616f62-c858-4f17-96fd-1a2e4dcb2fb6",
								},
							},
						},
					}),
					DeliveryPlannedMoment: documenttest.DateTime("2026-03-28 17:43:00.000"),
					InvoicedSum:           float64(0.0),
					ShippedSum:            float64(0.0),
					WaitSum:               float64(0.0),
					InternalOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/internalorder/d1e2f3a4-1234-1234-1234-d1e2f3a4b5c6",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/internalorder/metadata",
							Type:         model.MetaType("internalorder"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#internalorder/edit?id=d1e2f3a4-1234-1234-1234-d1e2f3a4b5c6",
						},
					}),
					CustomerOrders: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/22591b08-2859-11f1-0a83-1f2500000247",
								Type:      model.MetaType("customerorder"),
								MediaType: "application/json",
							},
						},
					},
					InvoicesIn: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/invoicein/e3f4a5b6-1234-1234-1234-e3f4a5b6c7d8",
								Type:      model.MetaType("invoicein"),
								MediaType: "application/json",
							},
						},
					},
					ProductionTasks: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/productiontask/b6c7d8e9-1234-1234-1234-b6c7d8e9f0a1",
								Type:      model.MetaType("productiontask"),
								MediaType: "application/json",
							},
						},
					},
					Supplies: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/supply/f4a5b6c7-1234-1234-1234-f4a5b6c7d8e9",
								Type:      model.MetaType("supply"),
								MediaType: "application/json",
							},
						},
					},
				}},
				{"unexpanded complete", "../../testdata/custom/purchaseorder/read.json", purchaseorder.PurchaseOrder{
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
						Href: "https://example.test/entity/purchaseorder/read-id",
						Type: model.MetaType("purchaseorder"),
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
								Type: model.MetaType("purchaseorderposition"),
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
					DeliveryPlannedMoment: documenttest.DateTime("2026-11-01 12:34:56.000"),
					InvoicedSum:           float64(4321.25),
					ShippedSum:            float64(4321.25),
					WaitSum:               float64(4321.25),
					InternalOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/internalOrder",
							Type: model.MetaType("internalOrder"),
						},
					}),
					CustomerOrders: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/customerOrders",
								Type: model.MetaType("customerOrders"),
							},
						},
					},
					InvoicesIn: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/invoicesIn",
								Type: model.MetaType("invoicesIn"),
							},
						},
					},
					ProductionTasks: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/productionTasks",
								Type: model.MetaType("productionTasks"),
							},
						},
					},
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
					var got purchaseorder.PurchaseOrder
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
				var got purchaseorder.PurchaseOrder
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
