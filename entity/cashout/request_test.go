package cashout_test

import (
	"encoding/json/jsontext"
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
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
				{"minimal Create", cashout.CreateRequest{
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
				}, "../../testdata/custom/cashout/create-minimal.golden.json", ""},
				{"full Create", fullCreateRequest(), "../../testdata/custom/cashout/create-full.golden.json", ""},
				{"full Update", fullUpdateRequest(), "../../testdata/custom/cashout/update-full.golden.json", ""},
				{"one field", cashout.UpdateRequest{
					Applicable: model.Value(false),
				}, "", "{\"applicable\":false}"},
				{"SyncID only", cashout.UpdateRequest{
					SyncID: model.Value("first-sync"),
				}, "", "{\"syncId\":\"first-sync\"}"},
				{"empty string", cashout.UpdateRequest{
					Description: model.Value(""),
				}, "", "{\"description\":\"\"}"},
				{"rate null", cashout.UpdateRequest{
					Rate: model.Null[model.DocumentRateRequest](),
				}, "", "{\"rate\":null}"},
				{"rate partial", cashout.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{
						Value: model.Value(float64(0)),
					}),
				}, "", "{\"rate\":{\"value\":0}}"},
				{"rate empty", cashout.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{}),
				}, "", "{\"rate\":{}}"},
				{"contract null", cashout.UpdateRequest{
					Contract: model.Null[model.Reference](),
				}, "", "{\"contract\":null}"},
				{"operations empty", cashout.UpdateRequest{
					Operations: model.Value([]model.FinanceOperationRequest{}),
				}, "", "{\"operations\":[]}"},
				{"operations nil", cashout.UpdateRequest{Operations: model.Value([]model.FinanceOperationRequest(nil))}, "", "{\"operations\":[]}"},
				{"zero sums", cashout.UpdateRequest{
					Sum:    model.Value(float64(0)),
					VATSum: model.Value(float64(0)),
				}, "", "{\"sum\":0,\"vatSum\":0}"},
				{"attributes nil", cashout.UpdateRequest{Attributes: model.Value([]model.Attribute(nil))}, "", "{\"attributes\":[]}"},
				{"unset", cashout.UpdateRequest{}, "", "{}"},
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

func fullCreateRequest() cashout.CreateRequest {
	return cashout.CreateRequest{
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

func fullUpdateRequest() cashout.UpdateRequest {
	return cashout.UpdateRequest{
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
			var got cashout.UpdateRequest
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
