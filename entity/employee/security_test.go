package employee_test

import (
	"context"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestSecurityFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/employee_security.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		var got employee.Security
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if !got.IsActive || got.Login != "example@lognex" || got.Group == nil || got.Group.Name != "Основной" ||
			len(got.AuthorizedHosts) != 1 || got.AuthorizedIPNetwork != "80.8.8.8" || got.Role == nil ||
			got.Role.Meta.Type != "individualrole" {
			t.Fatalf("%s: security = %+v", codec.Name, got)
		}
	}
}

func securityRequest() employee.UpdateSecurityRequest {
	yes := true
	return employee.UpdateSecurityRequest{
		Role: model.Value(employee.SecurityRoleRequest{
			Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/role/individual", Type: "individualrole", MediaType: "application/json"},
			Permissions: model.Value(model.Permissions{
				"importData": {Flag: &yes},
				"supply":     {Actions: map[string]model.PermissionScope{"view": model.ScopeAll}},
			}),
		}),
		Group:               model.Value(reference("https://api.moysklad.ru/api/remap/1.2/entity/group/f4b74c5e-443a-11eb-ac12-001000000002", "group")),
		AuthorizedHosts:     model.Value([]string{}),
		AuthorizedIPNetwork: model.Value("80.8.8.8"),
		AuthorizedIPNetmask: model.Value("255.255.255.0"),
	}
}

func TestSecurityRequestJSON(t *testing.T) {
	want, err := os.ReadFile("../../testdata/custom/employee/security-update.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		got, err := codec.Marshal(securityRequest())
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, got, want)
		// Роль без прав: сервер оставит прежние права индивидуальной роли.
		got, err = codec.Marshal(employee.UpdateSecurityRequest{Role: model.Value(employee.SecurityRoleRequest{Meta: model.Meta{Href: "h"}})})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, got, []byte(`{"role":{"meta":{"href":"h"}}}`))
	}
}

// Критерий успеха прав — group, обязательный в ответе; id у прав нет.
func TestServiceSecurity(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/employee_security.json")
	if err != nil {
		t.Fatal(err)
	}
	golden, err := os.ReadFile("../../testdata/custom/employee/security-update.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, op := range []struct {
		method string
		body   []byte
		invoke func(context.Context, *transport.Client) (*employee.Security, error)
	}{
		{"GET", nil, func(ctx context.Context, c *transport.Client) (*employee.Security, error) {
			s, err := employee.NewService(c)
			if err != nil {
				return nil, err
			}
			return s.GetSecurity(ctx, testID)
		}},
		{"PUT", golden, func(ctx context.Context, c *transport.Client) (*employee.Security, error) {
			s, err := employee.NewService(c)
			if err != nil {
				return nil, err
			}
			return s.UpdateSecurity(ctx, testID, securityRequest())
		}},
	} {
		t.Run(op.method, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(string(fixture)) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, "entity/employee/"+testID+"/security", op.body, op.invoke, func(got *employee.Security) {
						if got.Login != "example@lognex" || got.Role == nil || got.Role.Meta.Type != "individualrole" {
							t.Fatalf("security = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestServiceSecurityValidation(t *testing.T) {
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
	var info transport.ResponseInfo
	info.StatusCode = 503
	if _, err := s.GetSecurity(t.Context(), "bad", transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
		t.Errorf("invalid id: err = %v, info = %+v", err, info)
	}
	if _, err := s.UpdateSecurity(t.Context(), testID, employee.UpdateSecurityRequest{}); err == nil {
		t.Error("empty security update accepted")
	}
	if _, err := s.UpdateSecurity(t.Context(), "bad", securityRequest()); err == nil {
		t.Error("invalid id accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
