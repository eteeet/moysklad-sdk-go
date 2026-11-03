package saleschannel_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/saleschannel"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "53d1ff92-202b-11ec-0a82-050800000007"

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
	full   = saleschannel.CreateRequest{Name: "Канал SDK", Type: "MARKETPLACE", Description: model.Value(""), Code: model.Value("c-1"), ExternalCode: model.Value("ext"), Archived: model.Value(false), Shared: model.Value(false), Owner: model.Value(ref("entity/employee/faba7f37-2e58-11e6-8a84-bae500000028", "employee")), Group: model.Value(ref("entity/group/faba7f37-2e58-11e6-8a84-bae500000029", "group"))}
	update = saleschannel.UpdateRequest{Type: model.Value("FUTURE_TYPE"), Archived: model.Value(true)}
)

func TestSalesChannelFixture(t *testing.T) {
	data := read(t, "upstream/sales_channel.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got saleschannel.SalesChannel
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "Телефонный звонок" || got.Type != "OTHER" || got.Code != "PHONE_CALL" || got.Archived || got.Owner == nil || got.Updated.IsZero() {
				t.Fatalf("saleschannel = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse saleschannel.SalesChannel
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
			{"create full", full, read(t, "custom/saleschannel/create-full.golden.json")},
			{"update", update, read(t, "custom/saleschannel/update.golden.json")},
			{"create minimal", saleschannel.CreateRequest{Name: "C", Type: "OTHER"}, []byte(`{"name":"C","type":"OTHER"}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[saleschannel.CreateRequest](), reflect.TypeFor[saleschannel.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *saleschannel.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := saleschannel.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*saleschannel.SalesChannel, error)
	}{
		{"get", "GET", "entity/saleschannel/" + testID, nil, func(ctx context.Context, c *transport.Client) (*saleschannel.SalesChannel, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/saleschannel", read(t, "custom/saleschannel/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*saleschannel.SalesChannel, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/saleschannel/" + testID, read(t, "custom/saleschannel/update.golden.json"), func(ctx context.Context, c *transport.Client) (*saleschannel.SalesChannel, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *saleschannel.SalesChannel) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("saleschannel = %+v", got)
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
			_, err := s.Create(ctx, saleschannel.CreateRequest{Type: "OTHER"}, o...)
			return err
		},
		"create no type": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, saleschannel.CreateRequest{Name: "C"}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, saleschannel.UpdateRequest{}, o...)
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
	listtest.Run(t, "saleschannel", read(t, "upstream/sales_channel.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[saleschannel.SalesChannel] {
		return service(c).List
	}, func(t *testing.T, got saleschannel.SalesChannel) {
		if got.ID != testID || got.Type != "OTHER" {
			t.Fatalf("row = %+v", got)
		}
	})
}
