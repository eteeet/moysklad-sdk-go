package employee_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const accessID = "7944ef04-f831-11e5-7a69-971500188b19"

// Пример активации источника: login, group и роль кассира; ответ mailActivationRequired.
func TestAccessHTTP(t *testing.T) {
	activate := employee.ActivateRequest{
		Login: model.Value("newcashier@lognex"),
		Group: model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/group/f4b74c5e-443a-11eb-ac12-001000000003", Type: "group"}}),
		Role:  model.Value(employee.SecurityRoleRequest{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/role/cashier", Type: "systemrole"}}),
	}
	wantActivate := `{"login":"newcashier@lognex","group":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/group/f4b74c5e-443a-11eb-ac12-001000000003","type":"group"}},` +
		`"role":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/role/cashier","type":"systemrole"}}}`
	for _, tc := range []struct {
		name, path, body, response string
		status                     int
		ok                         bool
		invoke                     func(*employee.Service) (bool, error)
	}{
		{"activate new", "activate", wantActivate, `{"mailActivationRequired":true}`, 200, true, func(s *employee.Service) (bool, error) {
			a, err := s.Activate(t.Context(), accessID, activate)
			return a != nil && a.MailActivationRequired, err
		}},
		{"activate again", "activate", `{}`, `{"mailActivationRequired":false}`, 200, true, func(s *employee.Service) (bool, error) {
			a, err := s.Activate(t.Context(), accessID, employee.ActivateRequest{})
			return a != nil && !a.MailActivationRequired, err
		}},
		{"activate null", "activate", `{}`, `null`, 200, false, func(s *employee.Service) (bool, error) {
			a, err := s.Activate(t.Context(), accessID, employee.ActivateRequest{})
			return a == nil, err
		}},
		{"activate 204", "activate", `{}`, ``, 204, false, func(s *employee.Service) (bool, error) {
			a, err := s.Activate(t.Context(), accessID, employee.ActivateRequest{})
			return a == nil, err
		}},
		{"activate api error", "activate", `{}`, `{"errors":[{"error":"tariff","code":1024}]}`, 412, false, func(s *employee.Service) (bool, error) {
			a, err := s.Activate(t.Context(), accessID, employee.ActivateRequest{})
			return a == nil, err
		}},
		{"deactivate", "deactivate", ``, ``, 204, true, func(s *employee.Service) (bool, error) { return true, s.Deactivate(t.Context(), accessID) }},
		{"deactivate 200", "deactivate", ``, `{}`, 200, true, func(s *employee.Service) (bool, error) { return true, s.Deactivate(t.Context(), accessID) }},
		{"deactivate 202", "deactivate", ``, ``, 202, false, func(s *employee.Service) (bool, error) { return true, s.Deactivate(t.Context(), accessID) }},
		{"reset password", "resetpassword", ``, ``, 204, true, func(s *employee.Service) (bool, error) { return true, s.ResetPassword(t.Context(), accessID) }},
		{"reset password error", "resetpassword", ``, `{"errors":[{"error":"no email","code":3000}]}`, 400, false, func(s *employee.Service) (bool, error) {
			return true, s.ResetPassword(t.Context(), accessID)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				body, _ := io.ReadAll(r.Body)
				if r.Method != http.MethodPut || r.URL.Path != "/entity/employee/"+accessID+"/access/"+tc.path {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				if tc.body == "" && len(body) != 0 {
					t.Errorf("body = %s", body)
				} else if tc.body != "" {
					documenttest.JSONEqual(t, body, []byte(tc.body))
				}
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.response)
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			s, _ := employee.NewService(c)
			ok, err := tc.invoke(s)
			if tc.ok != (err == nil) || !ok || calls.Load() != 1 {
				t.Fatalf("ok = %v, err = %v, calls = %d", ok, err, calls.Load())
			}
			if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "employee "+accessID)) {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestAccessJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		got, err := codec.Marshal(employee.ActivateRequest{Login: model.Value(""), Group: model.Null[model.Reference]()})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, got, []byte(`{"login":"","group":null}`))
		var a employee.Activation
		if err := codec.Unmarshal([]byte(`{"mailActivationRequired":true,"future":1}`), &a); err != nil || !a.MailActivationRequired {
			t.Fatalf("%s: activation = %+v, err = %v", codec.Name, a, err)
		}
	}
}

func TestAccessLocalErrors(t *testing.T) {
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*employee.Activation, error) {
		s, _ := employee.NewService(c)
		return s.Activate(ctx, accessID, employee.ActivateRequest{})
	})
}

func TestAccessValidation(t *testing.T) {
	c, err := transport.NewClient(transport.Config{BaseURL: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := employee.NewService(c)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"activate": func(o ...transport.RequestOption) error {
			_, err := s.Activate(t.Context(), "bad", employee.ActivateRequest{}, o...)
			return err
		},
		"deactivate": func(o ...transport.RequestOption) error { return s.Deactivate(t.Context(), accessID+"/x", o...) },
		"reset":      func(o ...transport.RequestOption) error { return s.ResetPassword(t.Context(), "", o...) },
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
}
