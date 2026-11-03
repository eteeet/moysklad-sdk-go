package paymentout_test

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "11111111-1111-1111-1111-111111111111"

func TestServiceGet(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
		s, err := paymentout.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Get(ctx, testID)
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/paymentout/"+testID+"", nil, invoke, func(got *paymentout.PaymentOut) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceCreate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
		s, err := paymentout.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Create(ctx, fullCreateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/paymentout/create-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "POST", "entity/paymentout", body, invoke, func(got *paymentout.PaymentOut) {
				if got.ID != "server-id" || got.Name != "server-name" {
					t.Fatalf("response = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
func TestServiceUpdate(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
		s, err := paymentout.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Update(ctx, testID, fullUpdateRequest())
	}
	for _, tc := range documenttest.ResponseCases(`{"id":"server-id","name":"server-name","positions":{"meta":{"href":"https://example.test/positions","size":2}}}`) {
		t.Run(tc.Name, func(t *testing.T) {
			body, err := os.ReadFile("../../testdata/custom/paymentout/update-full.golden.json")
			if err != nil {
				t.Fatal(err)
			}
			documenttest.HTTP(t, tc, "PUT", "entity/paymentout/"+testID+"", body, invoke, func(got *paymentout.PaymentOut) {
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
		invoke func(context.Context, *transport.Client) (*paymentout.PaymentOut, error)
	}{
		{"minimal Create", "POST", "{\"organization\":{\"meta\":{\"href\":\"https://example.test/organization\",\"type\":\"organization\"}},\"agent\":{\"meta\":{\"href\":\"https://example.test/agent\",\"type\":\"agent\"}},\"expenseItem\":{\"meta\":{\"href\":\"https://example.test/expenseItem\",\"type\":\"expenseItem\"}}}", func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
			s, err := paymentout.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Create(ctx, paymentout.CreateRequest{
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
			})
		}},
		{"minimal Update", "PUT", "{\"applicable\":false}", func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
			s, err := paymentout.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, paymentout.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"SyncID Update", "PUT", "{\"syncId\":\"first-sync\"}", func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
			s, err := paymentout.NewService(client)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, paymentout.UpdateRequest{
				SyncID: model.Value("first-sync"),
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := "entity/paymentout"
			if tc.method == "PUT" {
				path += "/" + testID
			}
			documenttest.HTTP(t, documenttest.ResponseCase{Name: "plain", Success: true, Status: 200, Body: `{"id":"response"}`}, tc.method, path, []byte(tc.body), tc.invoke, func(got *paymentout.PaymentOut) {
				if got.ID != "response" {
					t.Fatal(got)
				}
			})
		})
	}
}
func TestServiceValidation(t *testing.T) {
	if s, err := paymentout.NewService(nil); err == nil || s != nil {
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
	service, err := paymentout.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		invoke func() (*paymentout.PaymentOut, error)
	}{
		{"bad Get id", func() (*paymentout.PaymentOut, error) { return service.Get(t.Context(), "../invalid") }},
		{"bad Update id", func() (*paymentout.PaymentOut, error) {
			return service.Update(t.Context(), "invalid", paymentout.UpdateRequest{
				Applicable: model.Value(false),
			})
		}},
		{"empty Update", func() (*paymentout.PaymentOut, error) {
			return service.Update(t.Context(), testID, paymentout.UpdateRequest{})
		}},
		{"missing organization", func() (*paymentout.PaymentOut, error) {
			return service.Create(t.Context(), paymentout.CreateRequest{
				Organization: model.Reference{
					Meta: model.Meta{},
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
			})
		}},
		{"missing agent", func() (*paymentout.PaymentOut, error) {
			return service.Create(t.Context(), paymentout.CreateRequest{
				Organization: model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/organization",
						Type: model.MetaType("organization"),
					},
				},
				Agent: model.Reference{
					Meta: model.Meta{},
				},
				ExpenseItem: model.Reference{
					Meta: model.Meta{
						Href: "https://example.test/expenseItem",
						Type: model.MetaType("expenseItem"),
					},
				},
			})
		}},
		{"missing expenseItem", func() (*paymentout.PaymentOut, error) {
			return service.Create(t.Context(), paymentout.CreateRequest{
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

// Синтетический отказ проверяет передачу APIError и отсутствие локального бизнес-валидатора.
func TestServiceNoClosingDocsConflict(t *testing.T) {
	operations := model.Value([]model.FinanceOperationRequest{{Meta: model.Meta{Href: "https://example.test/supply", Type: "supply"}}})
	createInput := paymentout.CreateRequest{Organization: model.Reference{Meta: model.Meta{Href: "https://example.test/organization", Type: "organization"}}, Agent: model.Reference{Meta: model.Meta{Href: "https://example.test/agent", Type: "agent"}}, ExpenseItem: model.Reference{Meta: model.Meta{Href: "https://example.test/expenseItem", Type: "expenseItem"}}, NoClosingDocs: model.Value(true), Operations: operations}
	updateInput := paymentout.UpdateRequest{NoClosingDocs: model.Value(true), Operations: operations}
	for _, tc := range []struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*paymentout.PaymentOut, error)
	}{
		{name: "Create", method: "POST", path: "entity/paymentout", body: `{"organization":{"meta":{"href":"https://example.test/organization","type":"organization"}},"agent":{"meta":{"href":"https://example.test/agent","type":"agent"}},"expenseItem":{"meta":{"href":"https://example.test/expenseItem","type":"expenseItem"}},"noClosingDocs":true,"operations":[{"meta":{"href":"https://example.test/supply","type":"supply"}}]}`, invoke: func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
			service, err := paymentout.NewService(client)
			if err != nil {
				return nil, err
			}
			return service.Create(ctx, createInput)
		}},
		{name: "Update", method: "PUT", path: "entity/paymentout" + "/" + testID, body: `{"noClosingDocs":true,"operations":[{"meta":{"href":"https://example.test/supply","type":"supply"}}]}`, invoke: func(ctx context.Context, client *transport.Client) (*paymentout.PaymentOut, error) {
			service, err := paymentout.NewService(client)
			if err != nil {
				return nil, err
			}
			return service.Update(ctx, testID, updateInput)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			documenttest.HTTP(t, documenttest.ResponseCase{Name: "noClosingDocs conflict", Status: 400, API: true, ErrorCode: 1234, Body: `{"errors":[{"error":"noClosingDocs incompatible with operations","code":1234}]}`}, tc.method, tc.path, []byte(tc.body), tc.invoke, func(*paymentout.PaymentOut) { t.Fatal("conflict accepted") })
		})
	}
}
