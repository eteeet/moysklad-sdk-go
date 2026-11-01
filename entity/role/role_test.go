package role_test

import (
	"context"
	"encoding/json/jsontext"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/role"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "736da682-ad8b-11eb-0a80-17ef000000d4"

func flag(v bool) model.Permission { return model.Permission{Flag: &v} }

func actions(kv ...string) model.Permission {
	p := model.Permission{Actions: map[string]model.PermissionScope{}}
	for i := 0; i < len(kv); i += 2 {
		p.Actions[kv[i]] = model.PermissionScope(kv[i+1])
	}
	return p
}

func TestRoleFixtures(t *testing.T) {
	custom, err := os.ReadFile("../../testdata/upstream/custom_role.json")
	if err != nil {
		t.Fatal(err)
	}
	preset, err := os.ReadFile("../../testdata/upstream/employee_role.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got role.Role
			if err := codec.Unmarshal(custom, &got); err != nil {
				t.Fatal(err)
			}
			flags, scopes := 0, 0
			for _, p := range got.Permissions {
				if p.Flag != nil {
					flags++
				} else if p.Actions != nil {
					scopes++
				}
			}
			if got.ID != testID || got.Name != "CustomRole" || got.Meta.Type != "customrole" || flags == 0 || scopes == 0 ||
				flags+scopes != len(got.Permissions) {
				t.Fatalf("role = %s %s, flags = %d, scopes = %d of %d", got.ID, got.Name, flags, scopes, len(got.Permissions))
			}
			// Права записываются обратно без потерь.
			out, err := codec.Marshal(got.Permissions)
			if err != nil {
				t.Fatal(err)
			}
			var raw struct {
				Permissions jsontext.Value `json:"permissions"`
			}
			if err := codec.Unmarshal(custom, &raw); err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, out, raw.Permissions)

			var sys role.Role
			if err := codec.Unmarshal(preset, &sys); err != nil {
				t.Fatal(err)
			}
			if sys.Meta.Type != "systemrole" || sys.ID != "" || sys.Permissions != nil {
				t.Fatalf("preset = %+v", sys)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name   string
			input  any
			golden string
			want   string
		}{
			{"create", role.CreateRequest{Name: "Кладовщик", Permissions: model.Value(model.Permissions{
				"importData": flag(true), "viewAudit": flag(false),
				"supply": actions("view", "ALL", "create", "OWN", "approve", "NO"), "script": actions("view", "OWN", "done", "ASSIGNEE"),
			})}, "../../testdata/custom/role/create.golden.json", ""},
			{"update", role.UpdateRequest{Permissions: model.Value(model.Permissions{
				"exportData": flag(false), "demand": actions("view", "OWN_GROUP"),
			})}, "../../testdata/custom/role/update.golden.json", ""},
			{"create name only", role.CreateRequest{Name: "R"}, "", `{"name":"R"}`},
			{"update name", role.UpdateRequest{Name: model.Value("R2")}, "", `{"name":"R2"}`},
			{"empty permissions", role.UpdateRequest{Permissions: model.Value(model.Permissions{})}, "", `{"permissions":{}}`},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				want := []byte(tc.want)
				if tc.golden != "" {
					var err error
					if want, err = os.ReadFile(tc.golden); err != nil {
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
	}
}

const serverRole = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/role/` + testID + `"},"id":"` + testID + `","name":"R","permissions":{"importData":true}}`

func checkRole(t *testing.T) func(*role.Role) {
	return func(got *role.Role) {
		if got.ID != testID || got.Permissions["importData"].Flag == nil {
			t.Fatalf("role = %+v", got)
		}
	}
}

func TestServiceGetCreateUpdate(t *testing.T) {
	create, err := os.ReadFile("../../testdata/custom/role/create.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	update, err := os.ReadFile("../../testdata/custom/role/update.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	// Запросы читаются из тех же golden, которые проверяет TestRequestJSON.
	var createInput role.CreateRequest
	var updateInput role.UpdateRequest
	if err := documenttest.Codecs[1].Unmarshal(create, &createInput); err != nil {
		t.Fatal(err)
	}
	if err := documenttest.Codecs[1].Unmarshal(update, &updateInput); err != nil {
		t.Fatal(err)
	}
	for _, op := range []struct {
		method, path string
		body         []byte
		invoke       func(context.Context, *transport.Client) (*role.Role, error)
	}{
		{"GET", "entity/role/" + testID, nil, func(ctx context.Context, c *transport.Client) (*role.Role, error) {
			s, err := role.NewService(c)
			if err != nil {
				return nil, err
			}
			return s.Get(ctx, testID)
		}},
		{"POST", "entity/role", create, func(ctx context.Context, c *transport.Client) (*role.Role, error) {
			s, err := role.NewService(c)
			if err != nil {
				return nil, err
			}
			return s.Create(ctx, createInput)
		}},
		{"PUT", "entity/role/" + testID, update, func(ctx context.Context, c *transport.Client) (*role.Role, error) {
			s, err := role.NewService(c)
			if err != nil {
				return nil, err
			}
			return s.Update(ctx, testID, updateInput)
		}},
	} {
		t.Run(op.method, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(serverRole) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, checkRole(t))
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Предустановленные роли: точный путь по имени, ответ без id, критерий — meta.href.
func TestServicePreset(t *testing.T) {
	preset, err := os.ReadFile("../../testdata/upstream/employee_role.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []role.Preset{role.PresetOwner, role.PresetAdmin, role.PresetIndividual, role.PresetCashier, role.PresetWorker, "future"} {
		invoke := func(ctx context.Context, c *transport.Client) (*role.Role, error) {
			s, err := role.NewService(c)
			if err != nil {
				return nil, err
			}
			return s.Preset(ctx, p)
		}
		t.Run(string(p), func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(string(preset)) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, "GET", "entity/role/"+string(p), nil, invoke, func(got *role.Role) {
						if got.Meta.Type != "systemrole" || got.ID != "" {
							t.Fatalf("preset = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, invoke)
		})
	}
	// Небезопасный сегмент пути не отправляется.
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, err := role.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []role.Preset{"../employee", "Admin", "", "a/b"} {
		if _, err := s.Preset(t.Context(), p); err == nil {
			t.Errorf("%q accepted", p)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceValidation(t *testing.T) {
	client, err := transport.NewClient(transport.Config{BaseURL: "http://example.test"})
	if err != nil {
		t.Fatal(err)
	}
	s, err := role.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(t.Context(), "admin"); err == nil {
		t.Error("preset name accepted as id")
	}
	if _, err := s.Create(t.Context(), role.CreateRequest{}); err == nil {
		t.Error("Create without name accepted")
	}
	if _, err := s.Update(t.Context(), testID, role.UpdateRequest{}); err == nil {
		t.Error("empty Update accepted")
	}
	if _, err := role.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/custom_role.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "role", row, "", "", func(c *transport.Client) listtest.ListFunc[role.Role] {
		s, err := role.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got role.Role) {
		if got.ID != testID || got.Name != "CustomRole" {
			t.Fatalf("row = %s %s", got.ID, got.Name)
		}
	})
}
