package purchasereturn_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
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
				want       purchasereturn.PurchaseReturn
			}{
				{"upstream", "../../testdata/upstream/purchase_return.json", purchasereturn.PurchaseReturn{
					ID:           "0ea5ec6c-58f3-11f1-0a83-1c0f00000817",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "ATK-179-PR-00001",
					Description:  "Возврат поставщику ATK-179",
					Code:         "ATK179-PR-001",
					ExternalCode: "ATK179PurchaseReturnFull",
					Created:      documenttest.DateTime("2007-02-07 17:16:41"),
					Updated:      documenttest.DateTime("2016-07-04 11:35:26"),
					Moment:       documenttest.DateTime("2026-05-26 12:00:00.000"),
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/metadata/states/69e39d04-41be-11e6-8a84-bae5000000d2",
							Type:      model.MetaType("state"),
							MediaType: "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/metadata",
						Type:         model.MetaType("purchasereturn"),
						MediaType:    "application/json",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/metadata/attributes/4fc7e124-41be-11e6-8a84-bae5000000cb",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "4fc7e124-41be-11e6-8a84-bae5000000cb",
							Name:  "AttributeName1",
							Type:  "boolean",
							Value: model.Value(jsontext.Value("true")),
						},
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/metadata/attributes/4fc7e7e7-41be-11e6-8a84-bae5000000cc",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "4fc7e7e7-41be-11e6-8a84-bae5000000cc",
							Name:  "AttributeName2",
							Type:  "text",
							Value: model.Value(jsontext.Value("\"деталь\"")),
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
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/147c1f1b-32ca-11e6-8a84-bae500000004",
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/147c1f1b-32ca-11e6-8a84-bae500000004/accounts/147c3231-32ca-11e6-8a84-bae500000005",
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
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817/positions",
								Type:      model.MetaType("purchasereturnposition"),
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
											Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817/positions/9cfa8b56-41bf-11e6-8a84-bae500000014",
											Type:      model.MetaType("purchasereturnposition"),
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
											UUIDHref:     "https://online.moysklad.ru/app/#good/edit?id=e64d0a86-2a99-11e9-ac12-000c00000041",
										},
									},
									Quantity: float64(1),
								},
								Price:      float64(263000),
								Discount:   float64(0),
								VAT:        0,
								VATEnabled: false,
								Things:     []string{},
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
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/paymentin/7bb4ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/paymentin/metadata",
								Type:         model.MetaType("paymentin"),
								MediaType:    "application/json",
							},
						},
						model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/cashin/6aa4ef04-f831-11e5-7a69-971500188b19",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/cashin/metadata",
								Type:         model.MetaType("cashin"),
								MediaType:    "application/json",
							},
						},
					},
					Files: documenttest.Pointer(model.Collection[model.Reference]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817/files",
								Type:      model.MetaType("files"),
								MediaType: "application/json",
							},
							Size:   2,
							Limit:  1000,
							Offset: 0,
						},
						Rows: []model.Reference{
							model.Reference{
								Meta: model.Meta{
									Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817/files/f2728180-6afd-4d37-8a13-f3b48069bbb6",
									Type:         model.MetaType("files"),
									MediaType:    "application/json",
									DownloadHref: "https://api.moysklad.ru/api/remap/1.2/download/f2728180-6afd-4d37-8a13-f3b48069bbb6",
								},
							},
							model.Reference{
								Meta: model.Meta{
									Href:         "https://api.moysklad.ru/api/remap/1.2/entity/purchasereturn/0ea5ec6c-58f3-11f1-0a83-1c0f00000817/files/933e41ac-1946-4bf0-9b21-51f2051f3e9f",
									Type:         model.MetaType("files"),
									MediaType:    "application/json",
									DownloadHref: "https://api.moysklad.ru/api/remap/1.2/download/933e41ac-1946-4bf0-9b21-51f2051f3e9f",
								},
							},
						},
					}),
					Supply: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/supply/7585391b-41c0-11e6-8a84-bae5000000de",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/supply/metadata",
							Type:         model.MetaType("supply"),
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
					FactureOut: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/factureout/9dd4ef04-f831-11e5-7a69-971500188b19",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/factureout/metadata",
							Type:         model.MetaType("factureout"),
							MediaType:    "application/json",
						},
					}),
				}},
				{"unexpanded complete", "../../testdata/custom/purchasereturn/read.json", purchasereturn.PurchaseReturn{
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
						Href: "https://example.test/entity/purchasereturn/read-id",
						Type: model.MetaType("purchasereturn"),
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
								Type: model.MetaType("purchasereturnposition"),
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
					Supply: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/supply",
							Type: model.MetaType("supply"),
						},
					}),
					FactureIn: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/factureIn",
							Type: model.MetaType("factureIn"),
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
					var got purchasereturn.PurchaseReturn
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
				var got purchasereturn.PurchaseReturn
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
