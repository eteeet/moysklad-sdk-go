package move_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "11111111-1111-1111-1111-111111111111"

func TestServiceGet(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*move.Move, error) {
		s, err := move.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Get(ctx, testID)
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/move/"+testID+"", nil, invoke, func(got *move.Move) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
				if got.Positions == nil || got.Positions.Rows != nil || got.Positions.Meta.Size != 2 {
					t.Fatalf("metadata/rows = %+v", got.Positions)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceCreate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*move.Move, error) {
		s, err := move.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Create(ctx, fullCreateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/move/create-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "POST", "entity/move", body, invoke, func(got *move.Move) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
				if got.Positions == nil || got.Positions.Rows != nil || got.Positions.Meta.Size != 2 {
					t.Fatalf("metadata/rows = %+v", got.Positions)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceUpdate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*move.Move, error) {
		s, err := move.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Update(ctx, testID, fullUpdateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/move/update-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "PUT", "entity/move/"+testID+"", body, invoke, func(got *move.Move) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
				if got.Positions == nil || got.Positions.Rows != nil || got.Positions.Meta.Size != 2 {
					t.Fatalf("metadata/rows = %+v", got.Positions)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceMinimalBodies(t *testing.T) {
	for _, tc := range []struct {
		name   string
		method string
		body   string
		invoke func(context.Context, *transport.Client) (*move.Move, error)
	}{
		{"minimal Create", "POST", "{\"organization\":{\"meta\":{\"href\":\"https://example.test/organization\",\"type\":\"organization\"}},\"sourceStore\":{\"meta\":{\"href\":\"https://example.test/sourceStore\",\"type\":\"sourceStore\"}},\"targetStore\":{\"meta\":{\"href\":\"https://example.test/targetStore\",\"type\":\"targetStore\"}}}", func(ctx context.Context, client *transport.Client) (*move.Move, error) {
			s, err := move.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Create(ctx, move.CreateRequest{
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
			})
		}},
		{"minimal Update", "PUT", "{\"applicable\":false}", func(ctx context.Context, client *transport.Client) (*move.Move, error) {
			s, err := move.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, move.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"SyncID Update", "PUT", "{\"syncId\":\"first-sync\"}", func(ctx context.Context, client *transport.Client) (*move.Move, error) {
			s, err := move.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, move.UpdateRequest{
				SyncID: model.Value("first-sync"),
			})
		}},
		{"positions omitted", "PUT", "{\"description\":\"\"}", func(ctx context.Context, client *transport.Client) (*move.Move, error) {
			s, err := move.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, move.UpdateRequest{
				Description: model.Value(""),
			})
		}},
		{"positions empty", "PUT", "{\"positions\":[]}", func(ctx context.Context, client *transport.Client) (*move.Move, error) {
			s, err := move.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, move.UpdateRequest{
				Positions: model.Value([]PositionRequest{}),
			})
		}},
		{"positions null", "PUT", "{\"positions\":null}", func(ctx context.Context, client *transport.Client) (*move.Move, error) {
			s, err := move.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, move.UpdateRequest{
				Positions: model.Null[[]PositionRequest](),
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := "entity/move"
			if tc.method == "PUT" {
				path += "/" + testID
			}
			documenttest.HTTP(t, documenttest.ResponseCase{Name: "plain", Success: true, Status: 200, Body: `{"id":"response"}`}, tc.method, path, []byte(tc.body), tc.invoke, func(got *move.Move) {
				if got.ID != "response" {
					t.Fatal(got)
				}
			})
		})
	}
}
func TestServiceValidation(t *testing.T) {
	if s, err := move.NewService(nil); err == nil || s != nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_, _ = w.Write([]byte(`{"id":"response"}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	service, err := move.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		invoke func() (*move.Move, error)
	}{
		{"bad Get id", func() (*move.Move, error) { return service.Get(t.Context(), "../invalid") }},
		{"bad Update id", func() (*move.Move, error) {
			return service.Update(t.Context(), "invalid", move.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"empty Update", func() (*move.Move, error) { return service.Update(t.Context(), testID, move.UpdateRequest{}) }},
		{"missing organization", func() (*move.Move, error) {
			return service.Create(t.Context(), move.CreateRequest{
				Organization: model.Reference{
					Meta: model.Meta{},
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
			})
		}},
		{"missing sourceStore", func() (*move.Move, error) {
			return service.Create(t.Context(), move.CreateRequest{
				Organization: model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/organization",
						Type: model.MetaType("organization"),
					},
				},
				SourceStore: model.Reference{
					Meta: model.Meta{},
				},
				TargetStore: model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/targetStore",
						Type: model.MetaType("targetStore"),
					},
				},
			})
		}},
		{"missing targetStore", func() (*move.Move, error) {
			return service.Create(t.Context(), move.CreateRequest{
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
					Meta: model.Meta{},
				},
			})
		}},
		{"too many Create positions", func() (*move.Move, error) {
			return service.Create(t.Context(), move.CreateRequest{
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
				Positions: model.Value(make([]move.PositionRequest, 1001))})
		}},
		{"too many Update positions", func() (*move.Move, error) {
			return service.Update(t.Context(), testID, move.UpdateRequest{Positions: model.Value(make([]move.PositionRequest, 1001))})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.invoke()
			if err == nil || got != nil {
				t.Fatalf("result=%+v error=%v", got, err)
			}
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("validation sent %d requests", calls.Load())
	}
	for _, size := range []int{0, 1000} {
		input := move.CreateRequest{
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
		}
		input.Positions = model.Value(make([]move.PositionRequest, size))
		if got, err := service.Create(t.Context(), input); err != nil || got == nil {
			t.Fatalf("Create size=%d err=%v", size, err)
		}
		if got, err := service.Update(t.Context(), testID, move.UpdateRequest{Positions: model.Value(make([]move.PositionRequest, size))}); err != nil || got == nil {
			t.Fatalf("Update size=%d err=%v", size, err)
		}
	}
}
func TestServiceGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7}} {
		limit := page.Limit
		if limit == 0 {
			limit = 1000
		}
		query := fmt.Sprintf("?limit=%d&offset=%d", limit, page.Offset)
		invoke := func(ctx context.Context, client *transport.Client) (*model.Collection[move.Position], error) {
			s, err := move.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.GetPositions(ctx, testID, page)
		}
		for _, tc := range documenttest.ResponseCases(`{"meta":{"href":"https://example.test/positions"},"rows":[{"id":"row","quantity":1.25,"price":12.75,"overhead":0.0,"cost":7.25,"shipped":0.5,"inTransit":0.75,"sourceSlot":{"meta":{"href":"https://example.test/source"}},"targetSlot":{"meta":{"href":"https://example.test/target"}}}]}`) {
			t.Run(tc.Name+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/move/"+testID+"/positions"+query, nil, invoke, func(got *model.Collection[move.Position]) {
					if len(got.Rows) != 1 || got.Rows[0].ID != "row" || got.Rows[0].Quantity != 1.25 || got.Rows[0].Price != 12.75 || got.Rows[0].Cost != 7.25 || got.Rows[0].Shipped != 0.5 || got.Rows[0].InTransit != 0.75 || got.Rows[0].SourceSlot == nil || got.Rows[0].TargetSlot == nil {
						t.Fatalf("rows = %+v", got.Rows)
					}
				})
			})
		}
		for _, tc := range []documenttest.ResponseCase{{Name: "empty page", Success: true, Status: 200, Body: `{"meta":{"href":"https://example.test/positions"},"rows":[]}`}, {Name: "missing meta", Status: 200, Body: `{"rows":[]}`}, {Name: "missing rows", Status: 200, Body: `{"meta":{"href":"https://example.test/positions"}}`}, {Name: "null rows", Status: 200, Body: `{"meta":{"href":"https://example.test/positions"},"rows":null}`}} {
			t.Run(tc.Name+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/move/"+testID+"/positions"+query, nil, invoke, func(got *model.Collection[move.Position]) {
					if got.Rows == nil || len(got.Rows) != 0 {
						t.Fatal(got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
	client, err := transport.NewClient(transport.Config{HTTPClient: &http.Client{Transport: rejectTransport{t}}})
	if err != nil {
		t.Fatal(err)
	}
	s, err := move.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		id   string
		page positions.PageOptions
	}{{"invalid", positions.PageOptions{}}, {testID, positions.PageOptions{Limit: -1}}, {testID, positions.PageOptions{Limit: 1001}}, {testID, positions.PageOptions{Offset: -1}}} {
		got, err := s.GetPositions(t.Context(), tc.id, tc.page)
		if err == nil || got != nil {
			t.Fatalf("result=%+v err=%v", got, err)
		}
	}
}

type rejectTransport struct{ t *testing.T }

func (r rejectTransport) RoundTrip(*http.Request) (*http.Response, error) {
	r.t.Error("validation sent HTTP")
	return nil, fmt.Errorf("unexpected request")
}
