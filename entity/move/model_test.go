package move_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/move"
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
				want       move.Move
			}{
				{"upstream", "../../testdata/upstream/move.json", move.Move{
					ID:           "5ac52c9a-4f67-11e6-8a84-bae500000084",
					AccountID:    "f976ed28-2e58-11e6-8a84-bae500000001",
					Name:         "705507",
					Description:  "Переезд 2",
					Code:         "12345",
					ExternalCode: "12345678",
					SyncID:       "5ac52c9a-4f67-11e6-8a84-bae500000084",
					Created:      documenttest.DateTime("2016-08-25 19:55:00"),
					Updated:      documenttest.DateTime("2016-07-21 20:44:20"),
					Deleted:      documenttest.DateTime("2016-08-26 19:55:00"),
					Moment:       documenttest.DateTime("2016-07-21 19:44:33"),
					Applicable:   false,
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
							Href:      "https://api.moysklad.ru/api/remap/1.2/entity/move/metadata/states/64c4689e-4f67-11e6-8a84-bae5000000b5",
							Type:      model.MetaType("state"),
							MediaType: "application/json",
						},
					}),
					Meta: model.Meta{
						Href:         "https://api.moysklad.ru/api/remap/1.2/entity/move/5ac52c9a-4f67-11e6-8a84-bae500000084",
						MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/move/metadata",
						Type:         model.MetaType("move"),
						MediaType:    "application/json",
						UUIDHref:     "https://online.moysklad.ru/app/#move/edit?id=5ac52c9a-4f67-11e6-8a84-bae500000084",
					},
					Attributes: []model.Attribute{
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/move/metadata/attributes/954d62cc-4f67-11e6-8a84-bae5000000bb",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "954d62cc-4f67-11e6-8a84-bae5000000bb",
							Name:  "AttributeName1",
							Type:  "time",
							Value: model.Value(jsontext.Value("\"2016-07-21 20:28:53\"")),
						},
						model.Attribute{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/move/metadata/attributes/954d6ab9-4f67-11e6-8a84-bae5000000bc",
								Type:      model.MetaType("attributemetadata"),
								MediaType: "application/json",
							},
							ID:    "954d6ab9-4f67-11e6-8a84-bae5000000bc",
							Name:  "AttributeName2",
							Type:  "long",
							Value: model.Value(jsontext.Value("43")),
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
						Value: documenttest.Pointer(float64(1.0)),
					}),
					Sum: float64(670.0),
					Positions: documenttest.Pointer(model.Collection[Position]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/move/5ac52c9a-4f67-11e6-8a84-bae500000084/positions",
								Type:      model.MetaType("moveposition"),
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
											Href:      "https://api.moysklad.ru/api/remap/1.2/entity/move/5ac52c9a-4f67-11e6-8a84-bae500000084/positions/5ac5407b-4f67-11e6-8a84-bae500000085",
											Type:      model.MetaType("moveposition"),
											MediaType: "application/json",
										},
										ID:        "5ac5407b-4f67-11e6-8a84-bae500000085",
										AccountID: "f976ed28-2e58-11e6-8a84-bae500000001",
									},
									Assortment: model.Reference{
										Meta: model.Meta{
											Href:         "https://api.moysklad.ru/api/remap/1.2/entity/product/4f2a0659-3304-11e6-8a84-bae50001c6a1",
											MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/product/metadata",
											Type:         model.MetaType("product"),
											MediaType:    "application/json",
											UUIDHref:     "https://online.moysklad.ru/app/#good/edit?id=e64d0a86-2a99-11e9-ac12-000c00000041",
										},
									},
									Quantity: float64(2.0),
								},
								Price:    float64(670.0),
								Overhead: float64(0),
								SourceSlot: documenttest.Pointer(model.Reference{
									Meta: model.Meta{
										Href:      "https://api.moysklad.ru/api/remap/1.2/entity/store/e94a6e65-4f64-11e6-8a84-bae500000066/slots/7fce2da5-684d-11ee-ac12-000c000000a2",
										Type:      model.MetaType("slot"),
										MediaType: "application/json",
									},
								}),
								TargetSlot: documenttest.Pointer(model.Reference{
									Meta: model.Meta{
										Href:      "https://api.moysklad.ru/api/remap/1.2/entity/store/faf3ff5b-2e58-11e6-8a84-bae500000050/slots/7fce37a5-684d-11ee-ac12-000c000000a3",
										Type:      model.MetaType("slot"),
										MediaType: "application/json",
									},
								}),
							},
						},
					}),
					Files: documenttest.Pointer(model.Collection[model.Reference]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href:      "https://api.moysklad.ru/api/remap/1.2/entity/move/5ac52c9a-4f67-11e6-8a84-bae500000084/files",
								Type:      model.MetaType("file"),
								MediaType: "application/json",
							},
							Size:   0,
							Limit:  1000,
							Offset: 0,
						},
						Rows: []model.Reference{},
					}),
					SourceStore: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/store/e94a6e65-4f64-11e6-8a84-bae500000066",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/store/metadata",
							Type:         model.MetaType("store"),
							MediaType:    "application/json",
						},
					}),
					TargetStore: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/store/faf3ff5b-2e58-11e6-8a84-bae500000050",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/store/metadata",
							Type:         model.MetaType("store"),
							MediaType:    "application/json",
						},
					}),
					InternalOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/internalorder/64e426af-b0d8-11e6-8a84-bae500000064",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/internalorder/metadata",
							Type:         model.MetaType("internalorder"),
							MediaType:    "application/json",
						},
					}),
					CustomerOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/64e426af-b0d8-11e6-8a84-bae500000064",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata",
							Type:         model.MetaType("customerorder"),
							MediaType:    "application/json",
						},
					}),
					Demand: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/demand/544236cb-4f6a-11e6-8a84-bae500000003",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata",
							Type:         model.MetaType("demand"),
							MediaType:    "application/json",
						},
					}),
					Supply: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href:         "https://api.moysklad.ru/api/remap/1.2/entity/supply/544236cb-4f6a-11e6-8a84-bae500000004",
							MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/supply/metadata",
							Type:         model.MetaType("supply"),
							MediaType:    "application/json",
						},
					}),
					ProductionTasks: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href:         "https://api.moysklad.ru/api/remap/1.2/entity/productiontask/544236cb-4f6a-11e6-8a84-bae500000005",
								MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/productiontask/metadata",
								Type:         model.MetaType("productiontask"),
								MediaType:    "application/json",
							},
						},
					},
					Overhead: documenttest.Pointer(model.DocumentOverhead{
						Sum:          float64(5555.0),
						Distribution: model.OverheadDistribution("weight"),
					}),
				}},
				{"unexpanded complete", "../../testdata/custom/move/read.json", move.Move{
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
						Href: "https://example.test/entity/move/read-id",
						Type: model.MetaType("move"),
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
					Positions: documenttest.Pointer(model.Collection[Position]{
						Meta: model.CollectionMeta{
							Meta: model.Meta{
								Href: "https://example.test/positions",
								Type: model.MetaType("moveposition"),
							},
							Size:   2,
							Limit:  1000,
							Offset: 0,
						},
					}),
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
					SourceStore: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/sourceStore",
							Type: model.MetaType("sourceStore"),
						},
					}),
					TargetStore: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/targetStore",
							Type: model.MetaType("targetStore"),
						},
					}),
					InternalOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/internalOrder",
							Type: model.MetaType("internalOrder"),
						},
					}),
					CustomerOrder: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/customerOrder",
							Type: model.MetaType("customerOrder"),
						},
					}),
					Demand: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/demand",
							Type: model.MetaType("demand"),
						},
					}),
					Supply: documenttest.Pointer(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/supply",
							Type: model.MetaType("supply"),
						},
					}),
					ProductionTasks: []model.Reference{
						model.Reference{
							Meta: model.Meta{
								Href: "https://example.test/productionTasks",
								Type: model.MetaType("productionTasks"),
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
					var got move.Move
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
				var got move.Move
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
