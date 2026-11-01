package supply_test

import (
	"encoding/json/jsontext"
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/supply"
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
				{"minimal Create", supply.CreateRequest{
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
					Store: model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/store",
							Type: model.MetaType("store"),
						},
					},
				}, "../../testdata/custom/supply/create-minimal.golden.json", ""},
				{"full Create", fullCreateRequest(), "../../testdata/custom/supply/create-full.golden.json", ""},
				{"full Update", fullUpdateRequest(), "../../testdata/custom/supply/update-full.golden.json", ""},
				{"one field", supply.UpdateRequest{
					Applicable: model.Value(false),
				}, "", "{\"applicable\":false}"},
				{"SyncID only", supply.UpdateRequest{
					SyncID: model.Value("first-sync"),
				}, "", "{\"syncId\":\"first-sync\"}"},
				{"empty string", supply.UpdateRequest{
					Description: model.Value(""),
				}, "", "{\"description\":\"\"}"},
				{"rate null", supply.UpdateRequest{
					Rate: model.Null[model.DocumentRateRequest](),
				}, "", "{\"rate\":null}"},
				{"rate partial", supply.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{
						Value: model.Value(float64(0)),
					}),
				}, "", "{\"rate\":{\"value\":0}}"},
				{"rate empty", supply.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{}),
				}, "", "{\"rate\":{}}"},
				{"contract null", supply.UpdateRequest{
					Contract: model.Null[model.Reference](),
				}, "", "{\"contract\":null}"},
				{"positions empty", supply.UpdateRequest{
					Positions: model.Value([]PositionRequest{}),
				}, "", "{\"positions\":[]}"},
				{"positions null", supply.UpdateRequest{
					Positions: model.Null[[]PositionRequest](),
				}, "", "{\"positions\":null}"},
				{"positions nil", supply.UpdateRequest{Positions: model.Value([]supply.PositionRequest(nil))}, "", "{\"positions\":[]}"},
				{"overhead {}", supply.UpdateRequest{
					Overhead: model.Value(model.DocumentOverheadRequest{}),
				}, "", "{\"overhead\":{}}"},
				{"overhead {'sum': 0}", supply.UpdateRequest{
					Overhead: model.Value(model.DocumentOverheadRequest{
						Sum: model.Value(float64(0)),
					}),
				}, "", "{\"overhead\":{\"sum\":0}}"},
				{"overhead None", supply.UpdateRequest{
					Overhead: model.Null[model.DocumentOverheadRequest](),
				}, "", "{\"overhead\":null}"},
				{"attributes nil", supply.UpdateRequest{Attributes: model.Value([]model.Attribute(nil))}, "", "{\"attributes\":[]}"},
				{"unset", supply.UpdateRequest{}, "", "{}"},
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

func fullCreateRequest() supply.CreateRequest {
	return supply.CreateRequest{
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
			},
		}),
		VATEnabled:  model.Value(false),
		VATIncluded: model.Value(false),
		Store: model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/store",
				Type: model.MetaType("store"),
			},
		},
		IncomingDate:   model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		IncomingNumber: model.Value("00007"),
		PurchaseOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/purchaseOrder",
				Type: model.MetaType("purchaseOrder"),
			},
		}),
		Overhead: model.Value(model.DocumentOverheadRequest{
			Sum:          model.Value(float64(12.25)),
			Distribution: model.Value(model.OverheadDistribution("weight")),
		}),
	}
}

func fullUpdateRequest() supply.UpdateRequest {
	return supply.UpdateRequest{
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
		IncomingDate:   model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		IncomingNumber: model.Value("00007"),
		PurchaseOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/purchaseOrder",
				Type: model.MetaType("purchaseOrder"),
			},
		}),
		Overhead: model.Value(model.DocumentOverheadRequest{
			Sum:          model.Value(float64(12.25)),
			Distribution: model.Value(model.OverheadDistribution("weight")),
		}),
	}
}

func TestUpdateExcludesReadonlyFields(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got supply.UpdateRequest
			if err := codec.Unmarshal([]byte(`{"id":"server-value","accountId":"server-value","created":"2026-11-01 12:34:56.000","updated":"2026-11-01 12:34:56.000","deleted":"2026-11-01 12:34:56.000","printed":true,"published":true,"meta":{"href":"https://example.test/readonly"},"sum":123.25,"vatSum":123.25,"payedSum":123.25,"payments":[{"meta":{"href":"https://example.test/readonly"}}],"files":{"meta":{"href":"https://example.test/readonly"},"rows":[]},"factureIn":{"meta":{"href":"https://example.test/readonly"}},"productionTask":{"meta":{"href":"https://example.test/readonly"}},"invoicesIn":[{"meta":{"href":"https://example.test/readonly"}}],"returns":[{"meta":{"href":"https://example.test/readonly"}}]}`), &got); err != nil {
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
