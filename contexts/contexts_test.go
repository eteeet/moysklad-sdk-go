package contexts_test

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/contexts"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestEmployeeContextFixture(t *testing.T) {
	data, err := os.ReadFile("../testdata/custom/contexts/employee.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got contexts.EmployeeContext
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != "5bf6430e-204c-11ec-c0a8-300d00000042" || got.UID != "admin@f1" || got.Group == nil || len(got.Permissions) != 73 {
				t.Fatalf("context = %+v", got)
			}
			// Задачи имеют действие done; пользовательское право — только view.
			if got.Permissions["task"].Actions["done"] != model.ScopeAll || got.Permissions["admin"].Actions["view"] != model.ScopeAll ||
				len(got.Permissions["admin"].Actions) != 1 {
				t.Fatalf("permissions = %+v, %+v", got.Permissions["task"], got.Permissions["admin"])
			}
		})
	}
}

// Неизвестные права, действия и области сохраняются, отсутствующие поля не выдумываются.
func TestEmployeeContextOpenSets(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got contexts.EmployeeContext
			// Нестроковое значение права не ломает разбор контекста: оно сохраняется в Raw.
			data := `{"id":"x","code":"c1","permissions":{"newEntity":{"view":"OWN","archive":"SOME_NEW_SCOPE"},"flag":true,"odd":[1]}}`
			if err := codec.Unmarshal([]byte(data), &got); err != nil {
				t.Fatal(err)
			}
			p := got.Permissions
			if got.Code != "c1" || p["newEntity"].Actions["archive"] != "SOME_NEW_SCOPE" || p["newEntity"].Actions["view"] != model.ScopeOwn ||
				p["flag"].Flag == nil || string(p["odd"].Raw) != "[1]" || got.Salary != nil || got.Cashiers != nil {
				t.Fatalf("context = %+v", got)
			}
		})
	}
}

func TestServiceEmployee(t *testing.T) {
	body, err := os.ReadFile("../testdata/custom/contexts/employee.json")
	if err != nil {
		t.Fatal(err)
	}
	invoke := func(ctx context.Context, client *transport.Client) (*contexts.EmployeeContext, error) {
		s, err := contexts.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Employee(ctx)
	}
	for _, tc := range documenttest.ResponseCases(string(body)) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "context/employee", nil, invoke, func(got *contexts.EmployeeContext) {
				if got.ID != "5bf6430e-204c-11ec-c0a8-300d00000042" || len(got.Permissions) != 73 {
					t.Fatalf("context = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestServiceApplication(t *testing.T) {
	body, err := os.ReadFile("../testdata/custom/contexts/application.json")
	if err != nil {
		t.Fatal(err)
	}
	invoke := func(ctx context.Context, client *transport.Client) (*contexts.ApplicationContext, error) {
		s, err := contexts.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Application(ctx)
	}
	for _, tc := range documenttest.ResponseCases(string(body)) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "context/application", nil, invoke, func(got *contexts.ApplicationContext) {
				if got.ID != "b58a6312-f958-11e9-ac12-000a00000020" || got.AppUID != "test.moysklad@reqwy1" || got.Meta.Type != "application" {
					t.Fatalf("context = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestNewServiceNilClient(t *testing.T) {
	if _, err := contexts.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}

// Basic и Bearer авторизация, отказ в доступе 401/403 — APIError, один запрос.
func TestServiceAuth(t *testing.T) {
	body, err := os.ReadFile("../testdata/custom/contexts/employee.json")
	if err != nil {
		t.Fatal(err)
	}
	basic, err := transport.BasicAuth("admin@shop", "password")
	if err != nil {
		t.Fatal(err)
	}
	bearer, err := transport.BearerToken("app-token")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		auth   transport.Auth
		want   string
		status int
	}{
		{"basic", basic, "Basic YWRtaW5Ac2hvcDpwYXNzd29yZA==", 200},
		{"bearer", bearer, "Bearer app-token", 200},
		{"unauthorized", bearer, "Bearer app-token", 401},
		{"forbidden", basic, "Basic YWRtaW5Ac2hvcDpwYXNzd29yZA==", 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Header.Get("Authorization") != tc.want || r.URL.RequestURI() != "/context/employee" {
					t.Errorf("request = %s %v", r.URL.RequestURI(), r.Header)
				}
				w.WriteHeader(tc.status)
				if tc.status == 200 {
					w.Write(body)
					return
				}
				w.Write([]byte(`{"errors":[{"error":"denied","code":1056}]}`))
			}))
			defer server.Close()
			client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Auth: tc.auth})
			if err != nil {
				t.Fatal(err)
			}
			s, err := contexts.NewService(client)
			if err != nil {
				t.Fatal(err)
			}
			got, err := s.Employee(t.Context())
			if calls.Load() != 1 || (tc.status == 200) != (err == nil) || tc.status != 200 && (got != nil || sdkerrors.StatusCode(err) != tc.status) {
				t.Fatalf("got = %v, err = %v, calls = %d", got, err, calls.Load())
			}
		})
	}
}
