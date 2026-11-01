package paymentout_test

import (
	"encoding/json/jsontext"
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name         string
				input        any
				golden, want string
			}{
				{"minimal Create", paymentout.CreateRequest{
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
					ExpenseItem: model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/expenseItem",
							Type: model.MetaType("expenseItem"),
						},
					},
				}, "../../testdata/custom/paymentout/create-minimal.golden.json", ""},
				{"full Create", fullCreateRequest(), "../../testdata/custom/paymentout/create-full.golden.json", ""},
				{"full Update", fullUpdateRequest(), "../../testdata/custom/paymentout/update-full.golden.json", ""},
				{"one field", paymentout.UpdateRequest{
					Applicable: model.Value(false),
				}, "", "{\"applicable\":false}"},
				{"SyncID only", paymentout.UpdateRequest{
					SyncID: model.Value("first-sync"),
				}, "", "{\"syncId\":\"first-sync\"}"},
				{"empty string", paymentout.UpdateRequest{
					Description: model.Value(""),
				}, "", "{\"description\":\"\"}"},
				{"rate null", paymentout.UpdateRequest{
					Rate: model.Null[model.DocumentRateRequest](),
				}, "", "{\"rate\":null}"},
				{"rate partial", paymentout.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{
						Value: model.Value(float64(0)),
					}),
				}, "", "{\"rate\":{\"value\":0}}"},
				{"rate empty", paymentout.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{}),
				}, "", "{\"rate\":{}}"},
				{"contract null", paymentout.UpdateRequest{
					Contract: model.Null[model.Reference](),
				}, "", "{\"contract\":null}"},
				{"operations empty", paymentout.UpdateRequest{
					Operations: model.Value([]model.FinanceOperationRequest{}),
				}, "", "{\"operations\":[]}"},
				{"operations nil", paymentout.UpdateRequest{Operations: model.Value([]model.FinanceOperationRequest(nil))}, "", "{\"operations\":[]}"},
				{"zero sums", paymentout.UpdateRequest{
					Sum:    model.Value(float64(0)),
					VATSum: model.Value(float64(0)),
				}, "", "{\"sum\":0,\"vatSum\":0}"},
				{"attributes nil", paymentout.UpdateRequest{Attributes: model.Value([]model.Attribute(nil))}, "", "{\"attributes\":[]}"},
				{"unset", paymentout.UpdateRequest{}, "", "{}"},
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

func fullCreateRequest() paymentout.CreateRequest {
	return paymentout.CreateRequest{
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
		Sum: model.Value(float64(123.25)),
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
		VATSum:         model.Value(float64(0)),
		PaymentPurpose: model.Value("request-paymentPurpose"),
		Operations: model.Value([]model.FinanceOperationRequest{
			model.FinanceOperationRequest{
				Meta: model.Meta{
					Href: "https://example.test/operation",
					Type: model.MetaType("futureoperation"),
				},
				LinkedSum: model.Value(float64(0)),
			},
			model.FinanceOperationRequest{
				Meta: model.Meta{
					Href: "https://example.test/second-operation",
					Type: model.MetaType("demand"),
				},
			},
		}),
		SalesChannel: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/salesChannel",
				Type: model.MetaType("salesChannel"),
			},
		}),
		ExpenseItem: model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/expenseItem",
				Type: model.MetaType("expenseItem"),
			},
		},
		FactureIn: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/factureIn",
				Type: model.MetaType("factureIn"),
			},
		}),
		NoClosingDocs: model.Value(false),
		AccrualDate:   model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
	}
}

func fullUpdateRequest() paymentout.UpdateRequest {
	return paymentout.UpdateRequest{
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
		Sum: model.Value(float64(123.25)),
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
		VATSum:         model.Value(float64(0)),
		PaymentPurpose: model.Value("request-paymentPurpose"),
		Operations: model.Value([]model.FinanceOperationRequest{
			model.FinanceOperationRequest{
				Meta: model.Meta{
					Href: "https://example.test/operation",
					Type: model.MetaType("futureoperation"),
				},
				LinkedSum: model.Value(float64(0)),
			},
			model.FinanceOperationRequest{
				Meta: model.Meta{
					Href: "https://example.test/second-operation",
					Type: model.MetaType("demand"),
				},
			},
		}),
		SalesChannel: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/salesChannel",
				Type: model.MetaType("salesChannel"),
			},
		}),
		ExpenseItem: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/expenseItem",
				Type: model.MetaType("expenseItem"),
			},
		}),
		FactureIn: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/factureIn",
				Type: model.MetaType("factureIn"),
			},
		}),
		NoClosingDocs: model.Value(false),
		AccrualDate:   model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
	}
}

func TestUpdateExcludesReadonlyFields(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got paymentout.UpdateRequest
			if err := codec.Unmarshal([]byte(`{"id":"server-value","accountId":"server-value","created":"2026-11-01 12:34:56.000","updated":"2026-11-01 12:34:56.000","deleted":"2026-11-01 12:34:56.000","printed":true,"published":true,"meta":{"href":"https://example.test/readonly"}}`), &got); err != nil {
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
