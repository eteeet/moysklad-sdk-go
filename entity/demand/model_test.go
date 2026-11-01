package demand_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/demand"
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
				want       demand.Demand
			}{
				{"upstream", "../../testdata/upstream/demand.json", demand.Demand{
					ID:           "c7d8e9f0-2859-11f1-0a83-1f2500000301",
					AccountID:    "c229f794-26ae-11f1-0a83-1f2100000002",
					Name:         "00001",
					Description:  "golden demand fixture",
					Code:         "DEM-001",
					ExternalCode: "ext-demand-001",
					SyncID:       "c7d8e9f0-2859-11f1-0a83-1f2500000301",
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
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata/states/c3faae3e-26ae-11f1-0a83-1f25000000d7",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata",
							Type:         model.MetaType("state"),
							MediaType:    "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/demand/c7d8e9f0-2859-11f1-0a83-1f2500000301",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata",
						Type:         model.MetaType("demand"),
						MediaType:    "application/json",
						UUIDHref:     "https://online.moysklad.ru/app/#demand/edit?id=c7d8e9f0-2859-11f1-0a83-1f2500000301",
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
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/demand/c7d8e9f0-2859-11f1-0a83-1f2500000301/positions",
								Type:      model.MetaType("demandposition"),
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
											Href:      "https://api.moysklad.ru/api/remap/1.2/entity/demand/c7d8e9f0-2859-11f1-0a83-1f2500000301/positions/d8e9f0a1-2859-11f1-0a83-1f2500000302",
											Type:      model.MetaType("demandposition"),
											MediaType: "application/json",
										},
										ID:        "d8e9f0a1-2859-11f1-0a83-1f2500000302",
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
									Quantity: float64(1.0),
								},
								Price:             float64(10000.0),
								Discount:          float64(0.0),
								VAT:               0,
								VATEnabled:        false,
								Things:            []string{},
								TrackingCodes:     []positions.TrackingCode{},
								TrackingCodes1162: []positions.TrackingCode1162{},
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
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/10b0654b-2859-11f1-0a83-1f2500000244",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/metadata",
							Type:         model.MetaType("saleschannel"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#saleschannel/edit?id=10b0654b-2859-11f1-0a83-1f2500000244",
						},
					}),
					CustomerOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/22591b08-2859-11f1-0a83-1f2500000247",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata",
							Type:         model.MetaType("customerorder"),
							MediaType:    "application/json",
							UUIDHref:     "https://online.moysklad.ru/app/#customerorder/edit?id=22591b08-2859-11f1-0a83-1f2500000247",
						},
					}),
					ShipmentAddress: "101000, г Москва, ул Тверская, д 1",
					ShipmentAddressFull: documenttest.Pointer(model.Address{
						City:       "Москва",
						House:      "1",
						PostalCode: "101000",
						Street:     "ул Тверская",
					}),
				}},
				{"unexpanded complete", "../../testdata/custom/demand/read.json", demand.Demand{
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
						Href: "https://example.test/entity/demand/read-id",
						Type: model.MetaType("demand"),
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
								Type: model.MetaType("demandposition"),
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
					SalesChannel: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/salesChannel",
							Type: model.MetaType("salesChannel"),
						},
					}),
					CustomerOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/customerOrder",
							Type: model.MetaType("customerOrder"),
						},
					}),
					FactureOut: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/factureOut",
							Type: model.MetaType("factureOut"),
						},
					}),
					ShipmentAddress: "request-shipmentAddress",
					ShipmentAddressFull: documenttest.Pointer(model.Address{
						AddInfo:    "вход 2",
						Apartment:  "7",
						City:       "Москва",
						Comment:    "позвонить",
						House:      "3",
						PostalCode: "101000",
						Street:     "Мясницкая",
						Country: documenttest.Pointer(model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/country",
								Type: model.MetaType("country"),
							},
						}),
						Region: documenttest.Pointer(model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/region",
								Type: model.MetaType("region"),
							},
						}),
					}),
					Returns: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/returns",
								Type: model.MetaType("returns"),
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
					InvoicesOut: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/invoicesOut",
								Type: model.MetaType("invoicesOut"),
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
					var got demand.Demand
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
				var got demand.Demand
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
