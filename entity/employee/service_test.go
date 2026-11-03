package employee_test

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "33333333-3333-3333-3333-333333333333"

const serverEmployee = `{"id":"server-id","lastName":"Иванов","name":"Иванов И."}`

func checkEmployee(t *testing.T) func(*employee.Employee) {
	return func(got *employee.Employee) {
		if got.ID != "server-id" || got.LastName != "Иванов" {
			t.Fatalf("response = %+v", got)
		}
	}
}

func TestServiceGet(t *testing.T) {
	invoke := func(ctx context.Context, client *transport.Client) (*employee.Employee, error) {
		s, err := employee.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Get(ctx, testID)
	}
	for _, tc := range documenttest.ResponseCases(serverEmployee) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/employee/"+testID, nil, invoke, checkEmployee(t))
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestServiceCreate(t *testing.T) {
	body, err := os.ReadFile("../../testdata/custom/employee/create-full.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	invoke := func(ctx context.Context, client *transport.Client) (*employee.Employee, error) {
		s, err := employee.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Create(ctx, fullCreateRequest())
	}
	for _, tc := range documenttest.ResponseCases(serverEmployee) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "POST", "entity/employee", body, invoke, checkEmployee(t))
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestServiceUpdate(t *testing.T) {
	body, err := os.ReadFile("../../testdata/custom/employee/update-full.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	invoke := func(ctx context.Context, client *transport.Client) (*employee.Employee, error) {
		s, err := employee.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Update(ctx, testID, fullUpdateRequest())
	}
	for _, tc := range documenttest.ResponseCases(serverEmployee) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "PUT", "entity/employee/"+testID, body, invoke, checkEmployee(t))
		})
	}
	documenttest.LocalErrors(t, invoke)
}

// Локальная проверка аргументов не отправляет запросов.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, err := employee.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(t.Context(), "bad"); err == nil {
		t.Error("invalid id accepted")
	}
	if _, err := s.Create(t.Context(), employee.CreateRequest{FirstName: model.Value("Петр")}); err == nil {
		t.Error("Create without lastName accepted")
	}
	if _, err := s.Update(t.Context(), testID, employee.UpdateRequest{FirstName: model.Value("Петр")}); err == nil {
		t.Error("Update without lastName accepted")
	}
	if _, err := s.Update(t.Context(), "bad", employee.UpdateRequest{LastName: "Л"}); err == nil {
		t.Error("Update with invalid id accepted")
	}
	if _, err := employee.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/employee.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "employee", row, "", "archived", func(c *transport.Client) listtest.ListFunc[employee.Employee] {
		s, err := employee.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got employee.Employee) {
		if got.ID != testID || got.UID != "admin@company" {
			t.Fatalf("row = %+v", got)
		}
	})
}
