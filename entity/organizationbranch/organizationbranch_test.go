package organizationbranch_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/organizationbranch"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "b2e9eb47-f831-11e5-7a69-971500188b20"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ), MediaType: "application/json"}}
}

var (
	full   = organizationbranch.CreateRequest{Name: "Филиал", Organization: ref("entity/organization/fae3561a-2e58-11e6-8a84-bae50000004e", "organization"), Description: model.Value(""), Code: model.Value("b-1"), ExternalCode: model.Value("ext"), Archived: model.Value(false), Email: model.Value("b@example.com"), Phone: model.Value(""), Fax: model.Value(""), ActualAddress: model.Value(""), ActualAddressFull: model.Value(model.AddressRequest{City: model.Value("Москва")}), RequisitesRU: model.Value(organizationbranch.RequisitesRURequest{KPP: model.Value("773301001")}), SyncID: model.Value("6f6a3a5e-0c1f-4c3a-9d1c-3b2a1d4e5f60"), Shared: model.Value(true), Owner: model.Value(ref("entity/employee/faba7f37-2e58-11e6-8a84-bae500000028", "employee")), Group: model.Value(ref("entity/group/faba7f37-2e58-11e6-8a84-bae500000029", "group"))}
	update = organizationbranch.UpdateRequest{Name: model.Value("Филиал 2"), RequisitesRU: model.Value(organizationbranch.RequisitesRURequest{KPP: model.Value("")}), ActualAddress: model.Value("")}
)

func TestOrganizationBranchFixture(t *testing.T) {
	data := read(t, "custom/organizationbranch/branch.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got organizationbranch.OrganizationBranch
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "Московский филиал" || got.Organization == nil || got.Organization.Meta.Type != "organization" || got.RequisitesRU == nil || got.RequisitesRU.KPP != "123456789" || got.ActualAddressFull == nil || got.ActualAddressFull.City != "Москва" || got.Email != "admin@moysklad.ru" || got.Created.IsZero() || got.Updated.IsZero() {
				t.Fatalf("organizationbranch = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse organizationbranch.OrganizationBranch
			if err := codec.Unmarshal([]byte(`{"id":"x","futureField":{"a":1}}`), &sparse); err != nil || sparse.ID != "x" {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/organizationbranch/create-full.golden.json")},
			{"update", update, read(t, "custom/organizationbranch/update.golden.json")},
			{"create minimal", organizationbranch.CreateRequest{Name: "B", Organization: ref("entity/organization/1", "organization")}, []byte(`{"name":"B","organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/1","type":"organization","mediaType":"application/json"}}}`)},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, tc.want)
			})
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[organizationbranch.CreateRequest](), reflect.TypeFor[organizationbranch.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "Created"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	if _, ok := reflect.TypeFor[organizationbranch.UpdateRequest]().FieldByName("Organization"); ok {
		t.Error("UpdateRequest contains create-only Organization")
	}
}

func service(c *transport.Client) *organizationbranch.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := organizationbranch.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organizationbranch/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*organizationbranch.OrganizationBranch, error)
	}{
		{"get", "GET", "entity/organizationbranch/" + testID, nil, func(ctx context.Context, c *transport.Client) (*organizationbranch.OrganizationBranch, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/organizationbranch", read(t, "custom/organizationbranch/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*organizationbranch.OrganizationBranch, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/organizationbranch/" + testID, read(t, "custom/organizationbranch/update.golden.json"), func(ctx context.Context, c *transport.Client) (*organizationbranch.OrganizationBranch, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *organizationbranch.OrganizationBranch) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("organizationbranch = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Ошибки входа — до HTTP; ResponseInfo прошлого вызова сброшен.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no name": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, organizationbranch.CreateRequest{Organization: ref("entity/organization/1", "organization")}, o...)
			return err
		},
		"create no organization": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, organizationbranch.CreateRequest{Name: "B"}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, organizationbranch.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "organizationbranch", read(t, "custom/organizationbranch/branch.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[organizationbranch.OrganizationBranch] {
		return service(c).List
	}, func(t *testing.T, got organizationbranch.OrganizationBranch) {
		if got.ID != testID || got.Code != "MSC-001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
