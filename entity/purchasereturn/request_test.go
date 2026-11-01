package purchasereturn_test

import (
	"encoding/json/jsontext"
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
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
				{"minimal Create", purchasereturn.CreateRequest{
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
				}, "../../testdata/custom/purchasereturn/create-minimal.golden.json", ""},
				{"full Create", fullCreateRequest(), "../../testdata/custom/purchasereturn/create-full.golden.json", ""},
				{"full Update", fullUpdateRequest(), "../../testdata/custom/purchasereturn/update-full.golden.json", ""},
				{"one field", purchasereturn.UpdateRequest{
					Applicable: model.Value(false),
				}, "", "{\"applicable\":false}"},
				{"SyncID only", purchasereturn.UpdateRequest{
					SyncID: model.Value("first-sync"),
				}, "", "{\"syncId\":\"first-sync\"}"},
				{"empty string", purchasereturn.UpdateRequest{
					Description: model.Value(""),
				}, "", "{\"description\":\"\"}"},
				{"rate null", purchasereturn.UpdateRequest{
					Rate: model.Null[model.DocumentRateRequest](),
				}, "", "{\"rate\":null}"},
				{"rate partial", purchasereturn.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{
						Value: model.Value(float64(0)),
					}),
				}, "", "{\"rate\":{\"value\":0}}"},
				{"rate empty", purchasereturn.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{}),
				}, "", "{\"rate\":{}}"},
				{"contract null", purchasereturn.UpdateRequest{
					Contract: model.Null[model.Reference](),
				}, "", "{\"contract\":null}"},
				{"positions empty", purchasereturn.UpdateRequest{
					Positions: model.Value([]PositionRequest{}),
				}, "", "{\"positions\":[]}"},
				{"positions null", purchasereturn.UpdateRequest{
					Positions: model.Null[[]PositionRequest](),
				}, "", "{\"positions\":null}"},
				{"positions nil", purchasereturn.UpdateRequest{Positions: model.Value([]purchasereturn.PositionRequest(nil))}, "", "{\"positions\":[]}"},
				{"attributes nil", purchasereturn.UpdateRequest{Attributes: model.Value([]model.Attribute(nil))}, "", "{\"attributes\":[]}"},
				{"unset", purchasereturn.UpdateRequest{}, "", "{}"},
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

func fullCreateRequest() purchasereturn.CreateRequest {
	return purchasereturn.CreateRequest{
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
		Supply: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/supply",
				Type: model.MetaType("supply"),
			},
		}),
	}
}

func fullUpdateRequest() purchasereturn.UpdateRequest {
	return purchasereturn.UpdateRequest{
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
		OrganizationAccount: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/organizationAccount",
				Type: model.MetaType("organizationAccount"),
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
	}
}

func TestUpdateExcludesReadonlyFields(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got purchasereturn.UpdateRequest
			if err := codec.Unmarshal([]byte(`{"id":"server-value","accountId":"server-value","created":"2026-11-01 12:34:56.000","updated":"2026-11-01 12:34:56.000","deleted":"2026-11-01 12:34:56.000","printed":true,"published":true,"meta":{"href":"https://example.test/readonly"},"sum":123.25,"agent":{"meta":{"href":"https://example.test/readonly"}},"agentAccount":{"meta":{"href":"https://example.test/readonly"}},"vatSum":123.25,"payedSum":123.25,"payments":[{"meta":{"href":"https://example.test/readonly"}}],"files":{"meta":{"href":"https://example.test/readonly"},"rows":[]},"supply":{"meta":{"href":"https://example.test/readonly"}},"factureIn":{"meta":{"href":"https://example.test/readonly"}},"factureOut":{"meta":{"href":"https://example.test/readonly"}}}`), &got); err != nil {
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
