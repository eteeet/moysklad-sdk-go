package group_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/group"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "11111111-1111-1111-1111-111111111111"

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
	full   = group.CreateRequest{Name: "Отдел SDK", Index: model.Value(0)}
	update = group.UpdateRequest{Name: model.Value("Отдел 2"), Index: model.Value(3)}
)

func TestGroupFixture(t *testing.T) {
	data := read(t, "upstream/group.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got group.Group
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "Основной" || got.Index != 0 || got.AccountID == "" {
				t.Fatalf("group = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse group.Group
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
			{"create full", full, read(t, "custom/group/create-full.golden.json")},
			{"update", update, read(t, "custom/group/update.golden.json")},
			{"create minimal", group.CreateRequest{Name: "G"}, []byte(`{"name":"G"}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[group.CreateRequest](), reflect.TypeFor[group.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *group.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := group.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/group/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*group.Group, error)
	}{
		{"get", "GET", "entity/group/" + testID, nil, func(ctx context.Context, c *transport.Client) (*group.Group, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/group", read(t, "custom/group/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*group.Group, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/group/" + testID, read(t, "custom/group/update.golden.json"), func(ctx context.Context, c *transport.Client) (*group.Group, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *group.Group) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("group = %+v", got)
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
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
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
			_, err := s.Create(ctx, group.CreateRequest{Index: model.Value(1)}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, group.UpdateRequest{}, o...)
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
	listtest.Run(t, "group", read(t, "upstream/group.json"), "", "", func(c *transport.Client) listtest.ListFunc[group.Group] {
		return service(c).List
	}, func(t *testing.T, got group.Group) {
		if got.ID != testID {
			t.Fatalf("row = %+v", got)
		}
	})
}
