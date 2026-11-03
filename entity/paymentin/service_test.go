package paymentin_test

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "11111111-1111-1111-1111-111111111111"

func TestServiceGet(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*paymentin.PaymentIn, error) {
		s, err := paymentin.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Get(ctx, testID)
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/paymentin/"+testID+"", nil, invoke, func(got *paymentin.PaymentIn) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceCreate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*paymentin.PaymentIn, error) {
		s, err := paymentin.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Create(ctx, fullCreateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/paymentin/create-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "POST", "entity/paymentin", body, invoke, func(got *paymentin.PaymentIn) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceUpdate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*paymentin.PaymentIn, error) {
		s, err := paymentin.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Update(ctx, testID, fullUpdateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/paymentin/update-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "PUT", "entity/paymentin/"+testID+"", body, invoke, func(got *paymentin.PaymentIn) {
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
		invoke func(context.Context, *transport.Client) (*paymentin.PaymentIn, error)
	}{
		{"minimal Create", "POST", "{\"organization\":{\"meta\":{\"href\":\"https://example.test/organization\",\"type\":\"organization\"}},\"agent\":{\"meta\":{\"href\":\"https://example.test/agent\",\"type\":\"agent\"}}}", func(ctx context.Context, client *transport.Client) (*paymentin.PaymentIn, error) {
			s, err := paymentin.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Create(ctx, paymentin.CreateRequest{
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
		{"minimal Update", "PUT", "{\"applicable\":false}", func(ctx context.Context, client *transport.Client) (*paymentin.PaymentIn, error) {
			s, err := paymentin.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, paymentin.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"SyncID Update", "PUT", "{\"syncId\":\"first-sync\"}", func(ctx context.Context, client *transport.Client) (*paymentin.PaymentIn, error) {
			s, err := paymentin.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, paymentin.UpdateRequest{
				SyncID: model.Value("first-sync"),
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := "entity/paymentin"
			if tc.method == "PUT" {
				path += "/" + testID
			}
			documenttest.HTTP(t, documenttest.ResponseCase{Name: "plain", Success: true, Status: 200, Body: `{"id":"response"}`}, tc.method, path, []byte(tc.body), tc.invoke, func(got *paymentin.PaymentIn) {
				if got.ID != "response" {
					t.Fatal(got)
				}
			})
		})
	}
}
func TestServiceValidation(t *testing.T) {
	if s, err := paymentin.NewService(nil); err == nil || s != nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_, _ = w.Write([]byte(`{"id":"response"}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	service, err := paymentin.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		invoke func() (*paymentin.PaymentIn, error)
	}{
		{"bad Get id", func() (*paymentin.PaymentIn, error) { return service.Get(t.Context(), "../invalid") }},
		{"bad Update id", func() (*paymentin.PaymentIn, error) {
			return service.Update(t.Context(), "invalid", paymentin.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"empty Update", func() (*paymentin.PaymentIn, error) {
			return service.Update(t.Context(), testID, paymentin.UpdateRequest{})
		}},
		{"missing organization", func() (*paymentin.PaymentIn, error) {
			return service.Create(t.Context(), paymentin.CreateRequest{
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
		{"missing agent", func() (*paymentin.PaymentIn, error) {
			return service.Create(t.Context(), paymentin.CreateRequest{
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
