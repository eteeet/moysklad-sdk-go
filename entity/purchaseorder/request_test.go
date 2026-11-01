package purchaseorder_test

import (
	"encoding/json/jsontext"
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

type PositionRequest = positions.PositionRequest

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name         string
				input        any
				golden, want string
			}{
				{"minimal Create", purchaseorder.CreateRequest{
					Organization: model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/organization",
							Type: model.MetaType("organization"),
						},
					},
					Agent: model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/agent",
							Type: model.MetaType("agent"),
						},
					},
				}, "../../testdata/custom/purchaseorder/create-minimal.golden.json", ""},
				{"full Create", fullCreateRequest(), "../../testdata/custom/purchaseorder/create-full.golden.json", ""},
				{"full Update", fullUpdateRequest(), "../../testdata/custom/purchaseorder/update-full.golden.json", ""},
				{"one field", purchaseorder.UpdateRequest{
					Applicable: model.Value(false),
				}, "", "{\"applicable\":false}"},
				{"SyncID only", purchaseorder.UpdateRequest{
					SyncID: model.Value("first-sync"),
				}, "", "{\"syncId\":\"first-sync\"}"},
				{"empty string", purchaseorder.UpdateRequest{
					Description: model.Value(""),
				}, "", "{\"description\":\"\"}"},
				{"rate null", purchaseorder.UpdateRequest{
					Rate: model.Null[model.DocumentRateRequest](),
				}, "", "{\"rate\":null}"},
				{"rate partial", purchaseorder.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{
						Value: model.Value(float64(0)),
					}),
				}, "", "{\"rate\":{\"value\":0}}"},
				{"rate empty", purchaseorder.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{}),
				}, "", "{\"rate\":{}}"},
				{"contract null", purchaseorder.UpdateRequest{
					Contract: model.Null[model.Reference](),
				}, "", "{\"contract\":null}"},
				{"positions empty", purchaseorder.UpdateRequest{
					Positions: model.Value([]PositionRequest{}),
				}, "", "{\"positions\":[]}"},
				{"positions null", purchaseorder.UpdateRequest{
					Positions: model.Null[[]PositionRequest](),
				}, "", "{\"positions\":null}"},
				{"positions nil", purchaseorder.UpdateRequest{Positions: model.Value([]purchaseorder.PositionRequest(nil))}, "", "{\"positions\":[]}"},
				{"inTransit 0", purchaseorder.UpdateRequest{
					Positions: model.Value([]PositionRequest{
						PositionRequest{
							ItemRequest: positions.ItemRequest{
								CoreRequest: positions.CoreRequest{
									ID: model.Value("row"),
								},
							},
							InTransit: model.Value(float64(0)),
						},
					}),
				}, "", "{\"positions\":[{\"id\":\"row\",\"inTransit\":0}]}"},
				{"inTransit 0.25", purchaseorder.UpdateRequest{
					Positions: model.Value([]PositionRequest{
						PositionRequest{
							ItemRequest: positions.ItemRequest{
								CoreRequest: positions.CoreRequest{
									ID: model.Value("row"),
								},
							},
							InTransit: model.Value(float64(0.25)),
						},
					}),
				}, "", "{\"positions\":[{\"id\":\"row\",\"inTransit\":0.25}]}"},
				{"inTransit None", purchaseorder.UpdateRequest{
					Positions: model.Value([]PositionRequest{
						PositionRequest{
							ItemRequest: positions.ItemRequest{
								CoreRequest: positions.CoreRequest{
									ID: model.Value("row"),
								},
							},
							InTransit: model.Null[float64](),
						},
					}),
				}, "", "{\"positions\":[{\"id\":\"row\",\"inTransit\":null}]}"},
				{"attributes nil", purchaseorder.UpdateRequest{Attributes: model.Value([]model.Attribute(nil))}, "", "{\"attributes\":[]}"},
				{"unset", purchaseorder.UpdateRequest{}, "", "{}"},
			} {
				t.Run(tc.name, func(t *testing.T) {
					want := []byte(tc.want)
					if tc.golden != "" {
						var err error
						want, err = os.ReadFile(tc.golden)
						if err != nil {
							t.Fatal(err)
						}
					}
					got, err := codec.Marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					documenttest.JSONEqual(t, got, want)
				})
			}
		})
	}
}

func fullCreateRequest() purchaseorder.CreateRequest {
	return purchaseorder.CreateRequest{
		Name:         model.Value("request-name"),
		Description:  model.Value(""),
		Code:         model.Value(""),
		ExternalCode: model.Value("request-externalCode"),
		SyncID:       model.Value("request-syncId"),
		Moment:       model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable:   model.Value(false),
		Shared:       model.Value(false),
		Organization: model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/organization",
				Type: model.MetaType("organization"),
			},
		},
		Owner: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/owner",
				Type: model.MetaType("owner"),
			},
		}),
		Group: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/group",
				Type: model.MetaType("group"),
			},
		}),
		Project: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/project",
				Type: model.MetaType("project"),
			},
		}),
		State: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/state",
				Type: model.MetaType("state"),
			},
		}),
		Attributes: model.Value([]model.Attribute{
			model.Attribute{
				Meta: model.Meta{
					Href: "https://example.test/attributes/long",
				},
				Type:  "long",
				Value: model.Value(jsontext.Value("9007199254740993")),
			},
		}),
		Rate: model.Value(model.DocumentRateRequest{
			Currency: model.Value(model.Reference{
				Meta: model.Meta{
					Href: "https://example.test/currency",
					Type: model.MetaType("currency"),
				},
			}),
		}),
		Agent: model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/agent",
				Type: model.MetaType("agent"),
			},
		},
		OrganizationAccount: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/organizationAccount",
				Type: model.MetaType("organizationAccount"),
			},
		}),
		AgentAccount: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/agentAccount",
				Type: model.MetaType("agentAccount"),
			},
		}),
		Contract: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/contract",
				Type: model.MetaType("contract"),
			},
		}),
		Positions: model.Value([]PositionRequest{
			PositionRequest{
				ItemRequest: positions.ItemRequest{
					CoreRequest: positions.CoreRequest{},
					Assortment: model.Value(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/assortment",
							Type: model.MetaType("assortment"),
						},
					}),
					Quantity: model.Value(float64(1.25)),
				},
				Price:      model.Value(float64(12345.67)),
				Discount:   model.Value(float64(-5.5)),
				VAT:        model.Value(0),
				VATEnabled: model.Value(false),
				InTransit:  model.Value(float64(0.25)),
			},
		}),
		VATEnabled:  model.Value(false),
		VATIncluded: model.Value(false),
		Store: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/store",
				Type: model.MetaType("store"),
			},
		}),
		DeliveryPlannedMoment: model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		InternalOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/internalOrder",
				Type: model.MetaType("internalOrder"),
			},
		}),
	}
}

func fullUpdateRequest() purchaseorder.UpdateRequest {
	return purchaseorder.UpdateRequest{
		Name:         model.Value("request-name"),
		Description:  model.Value(""),
		Code:         model.Value(""),
		ExternalCode: model.Value("request-externalCode"),
		SyncID:       model.Value("request-syncId"),
		Moment:       model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable:   model.Value(false),
		Shared:       model.Value(false),
		Organization: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/organization",
				Type: model.MetaType("organization"),
			},
		}),
		Owner: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/owner",
				Type: model.MetaType("owner"),
			},
		}),
		Group: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/group",
				Type: model.MetaType("group"),
			},
		}),
		Project: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/project",
				Type: model.MetaType("project"),
			},
		}),
		State: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/state",
				Type: model.MetaType("state"),
			},
		}),
		Attributes: model.Value([]model.Attribute{
			model.Attribute{
				Meta: model.Meta{
					Href: "https://example.test/attributes/long",
				},
				Type:  "long",
				Value: model.Value(jsontext.Value("9007199254740993")),
			},
		}),
		Rate: model.Value(model.DocumentRateRequest{
			Currency: model.Value(model.Reference{
				Meta: model.Meta{
					Href: "https://example.test/currency",
					Type: model.MetaType("currency"),
				},
			}),
		}),
		Agent: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/agent",
				Type: model.MetaType("agent"),
			},
		}),
		OrganizationAccount: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/organizationAccount",
				Type: model.MetaType("organizationAccount"),
			},
		}),
		AgentAccount: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/agentAccount",
				Type: model.MetaType("agentAccount"),
			},
		}),
		Contract: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/contract",
				Type: model.MetaType("contract"),
			},
		}),
		Positions: model.Value([]PositionRequest{
			PositionRequest{
				ItemRequest: positions.ItemRequest{
					CoreRequest: positions.CoreRequest{
						ID: model.Value("existing-row"),
					},
					Assortment: model.Value(model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/assortment",
							Type: model.MetaType("assortment"),
						},
					}),
					Quantity: model.Value(float64(1.25)),
				},
				Price:      model.Value(float64(12345.67)),
				Discount:   model.Value(float64(-5.5)),
				VAT:        model.Value(0),
				VATEnabled: model.Value(false),
				InTransit:  model.Value(float64(0.25)),
			},
		}),
		VATEnabled:  model.Value(false),
		VATIncluded: model.Value(false),
		Store: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/store",
				Type: model.MetaType("store"),
			},
		}),
		DeliveryPlannedMoment: model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		InternalOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/internalOrder",
				Type: model.MetaType("internalOrder"),
			},
		}),
	}
}

func TestUpdateExcludesReadonlyFields(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got purchaseorder.UpdateRequest
			if err := codec.Unmarshal([]byte(`{"id":"server-value","accountId":"server-value","created":"2026-11-01 12:34:56.000","updated":"2026-11-01 12:34:56.000","deleted":"2026-11-01 12:34:56.000","printed":true,"published":true,"meta":{"href":"https://example.test/readonly"},"sum":123.25,"vatSum":123.25,"payedSum":123.25,"payments":[{"meta":{"href":"https://example.test/readonly"}}],"files":{"meta":{"href":"https://example.test/readonly"},"rows":[]},"invoicedSum":123.25,"shippedSum":123.25,"waitSum":123.25,"customerOrders":[{"meta":{"href":"https://example.test/readonly"}}],"invoicesIn":[{"meta":{"href":"https://example.test/readonly"}}],"productionTasks":[{"meta":{"href":"https://example.test/readonly"}}],"supplies":[{"meta":{"href":"https://example.test/readonly"}}]}`), &got); err != nil {
				t.Fatal(err)
			}
			data, err := codec.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, data, []byte(`{}`))
		})
	}
}
