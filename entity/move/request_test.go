package move_test

import (
	"encoding/json/jsontext"
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/move"
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
				{"minimal Create", move.CreateRequest{
					Organization: model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/organization",
							Type: model.MetaType("organization"),
						},
					},
					SourceStore: model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/sourceStore",
							Type: model.MetaType("sourceStore"),
						},
					},
					TargetStore: model.Reference{
						Meta: model.Meta{
							Href: "https://example.test/targetStore",
							Type: model.MetaType("targetStore"),
						},
					},
				}, "../../testdata/custom/move/create-minimal.golden.json", ""},
				{"full Create", fullCreateRequest(), "../../testdata/custom/move/create-full.golden.json", ""},
				{"full Update", fullUpdateRequest(), "../../testdata/custom/move/update-full.golden.json", ""},
				{"one field", move.UpdateRequest{
					Applicable: model.Value(false),
				}, "", "{\"applicable\":false}"},
				{"SyncID only", move.UpdateRequest{
					SyncID: model.Value("first-sync"),
				}, "", "{\"syncId\":\"first-sync\"}"},
				{"empty string", move.UpdateRequest{
					Description: model.Value(""),
				}, "", "{\"description\":\"\"}"},
				{"rate null", move.UpdateRequest{
					Rate: model.Null[model.DocumentRateRequest](),
				}, "", "{\"rate\":null}"},
				{"rate partial", move.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{
						Value: model.Value(float64(0)),
					}),
				}, "", "{\"rate\":{\"value\":0}}"},
				{"rate empty", move.UpdateRequest{
					Rate: model.Value(model.DocumentRateRequest{}),
				}, "", "{\"rate\":{}}"},
				{"positions empty", move.UpdateRequest{
					Positions: model.Value([]PositionRequest{}),
				}, "", "{\"positions\":[]}"},
				{"positions null", move.UpdateRequest{
					Positions: model.Null[[]PositionRequest](),
				}, "", "{\"positions\":null}"},
				{"positions nil", move.UpdateRequest{Positions: model.Value([]move.PositionRequest(nil))}, "", "{\"positions\":[]}"},
				{"overhead {}", move.UpdateRequest{
					Overhead: model.Value(model.DocumentOverheadRequest{}),
				}, "", "{\"overhead\":{}}"},
				{"overhead {'sum': 0}", move.UpdateRequest{
					Overhead: model.Value(model.DocumentOverheadRequest{
						Sum: model.Value(float64(0)),
					}),
				}, "", "{\"overhead\":{\"sum\":0}}"},
				{"overhead None", move.UpdateRequest{
					Overhead: model.Null[model.DocumentOverheadRequest](),
				}, "", "{\"overhead\":null}"},
				{"attributes nil", move.UpdateRequest{Attributes: model.Value([]model.Attribute(nil))}, "", "{\"attributes\":[]}"},
				{"unset", move.UpdateRequest{}, "", "{}"},
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

func fullCreateRequest() move.CreateRequest {
	return move.CreateRequest{
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
				SourceSlot: model.Value(model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/sourceSlot",
						Type: model.MetaType("sourceSlot"),
					},
				}),
				TargetSlot: model.Value(model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/targetSlot",
						Type: model.MetaType("targetSlot"),
					},
				}),
			},
		}),
		SourceStore: model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/sourceStore",
				Type: model.MetaType("sourceStore"),
			},
		},
		TargetStore: model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/targetStore",
				Type: model.MetaType("targetStore"),
			},
		},
		InternalOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/internalOrder",
				Type: model.MetaType("internalOrder"),
			},
		}),
		CustomerOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/customerOrder",
				Type: model.MetaType("customerOrder"),
			},
		}),
		Overhead: model.Value(model.DocumentOverheadRequest{
			Sum:          model.Value(float64(12.25)),
			Distribution: model.Value(model.OverheadDistribution("weight")),
		}),
	}
}

func fullUpdateRequest() move.UpdateRequest {
	return move.UpdateRequest{
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
				SourceSlot: model.Value(model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/sourceSlot",
						Type: model.MetaType("sourceSlot"),
					},
				}),
				TargetSlot: model.Value(model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/targetSlot",
						Type: model.MetaType("targetSlot"),
					},
				}),
			},
		}),
		SourceStore: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/sourceStore",
				Type: model.MetaType("sourceStore"),
			},
		}),
		TargetStore: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/targetStore",
				Type: model.MetaType("targetStore"),
			},
		}),
		InternalOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/internalOrder",
				Type: model.MetaType("internalOrder"),
			},
		}),
		CustomerOrder: model.Value(model.Reference{
			Meta: model.Meta{
				Href: "https://example.test/customerOrder",
				Type: model.MetaType("customerOrder"),
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
			var got move.UpdateRequest
			if err := codec.Unmarshal([]byte(`{"id":"server-value","accountId":"server-value","created":"2026-11-01 12:34:56.000","updated":"2026-11-01 12:34:56.000","deleted":"2026-11-01 12:34:56.000","printed":true,"published":true,"meta":{"href":"https://example.test/readonly"},"sum":123.25,"files":{"meta":{"href":"https://example.test/readonly"},"rows":[]},"demand":{"meta":{"href":"https://example.test/readonly"}},"supply":{"meta":{"href":"https://example.test/readonly"}},"productionTasks":[{"meta":{"href":"https://example.test/readonly"}}]}`), &got); err != nil {
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
