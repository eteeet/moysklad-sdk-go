package cashin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "11111111-1111-1111-1111-111111111111"

func TestServiceGet(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*cashin.CashIn, error) {
		s, err := cashin.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Get(ctx, testID)
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/cashin/"+testID+"", nil, invoke, func(got *cashin.CashIn) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceCreate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*cashin.CashIn, error) {
		s, err := cashin.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Create(ctx, fullCreateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/cashin/create-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "POST", "entity/cashin", body, invoke, func(got *cashin.CashIn) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceUpdate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*cashin.CashIn, error) {
		s, err := cashin.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Update(ctx, testID, fullUpdateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/cashin/update-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "PUT", "entity/cashin/"+testID+"", body, invoke, func(got *cashin.CashIn) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
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
		invoke func(context.Context, *transport.Client) (*cashin.CashIn, error)
	}{
		{"minimal Create", "POST", "{\"organization\":{\"meta\":{\"href\":\"https://example.test/organization\",\"type\":\"organization\"}},\"agent\":{\"meta\":{\"href\":\"https://example.test/agent\",\"type\":\"agent\"}}}", func(ctx context.Context, client *transport.Client) (*cashin.CashIn, error) {
			s, err := cashin.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Create(ctx, cashin.CreateRequest{
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
			})
		}},
		{"minimal Update", "PUT", "{\"applicable\":false}", func(ctx context.Context, client *transport.Client) (*cashin.CashIn, error) {
			s, err := cashin.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, cashin.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"SyncID Update", "PUT", "{\"syncId\":\"first-sync\"}", func(ctx context.Context, client *transport.Client) (*cashin.CashIn, error) {
			s, err := cashin.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, cashin.UpdateRequest{
				SyncID: model.Value("first-sync"),
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := "entity/cashin"
			if tc.method == "PUT" {
				path += "/" + testID
			}
			documenttest.HTTP(t, documenttest.ResponseCase{Name: "plain", Success: true, Status: 200, Body: `{"id":"response"}`}, tc.method, path, []byte(tc.body), tc.invoke, func(got *cashin.CashIn) {
				if got.ID != "response" {
					t.Fatal(got)
				}
			})
		})
	}
}
func TestServiceValidation(t *testing.T) {
	if s, err := cashin.NewService(nil); err == nil || s != nil {
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
	service, err := cashin.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		invoke func() (*cashin.CashIn, error)
	}{
		{"bad Get id", func() (*cashin.CashIn, error) { return service.Get(t.Context(), "../invalid") }},
		{"bad Update id", func() (*cashin.CashIn, error) {
			return service.Update(t.Context(), "invalid", cashin.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"empty Update", func() (*cashin.CashIn, error) { return service.Update(t.Context(), testID, cashin.UpdateRequest{}) }},
		{"missing organization", func() (*cashin.CashIn, error) {
			return service.Create(t.Context(), cashin.CreateRequest{
				Organization: model.Reference{
					Meta: model.Meta{},
				},
				Agent: model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/agent",
						Type: model.MetaType("agent"),
					},
				},
			})
		}},
		{"missing agent", func() (*cashin.CashIn, error) {
			return service.Create(t.Context(), cashin.CreateRequest{
				Organization: model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/organization",
						Type: model.MetaType("organization"),
					},
				},
				Agent: model.Reference{
					Meta: model.Meta{},
				},
			})
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
}
