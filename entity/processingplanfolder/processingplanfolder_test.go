package processingplanfolder_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/processingplanfolder"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "3ea73e1a-2cad-11e6-8a84-bae50000001d"
	fixtureID = "8c329385-bee9-11ec-0a82-04c500000102"
	path      = "entity/processingplanfolder/" + testID
)

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

func service(c *transport.Client) *processingplanfolder.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := processingplanfolder.NewService(c)
	return s
}

var (
	owner = ref("entity/employee/faba7f37-2e58-11e6-8a84-bae500000028", "employee")
	group = ref("entity/group/f97aa1fb-2e58-11e6-8a84-bae500000002", "group")
	full  = processingplanfolder.CreateRequest{Name: "Группа 2", Description: model.Value("Техкарты мебели"), Code: model.Value("13321"),
		ExternalCode: model.Value("extGroup"), Archived: model.Value(false), Shared: model.Value(true), Owner: model.Value(owner), Group: model.Value(group)}
	update = processingplanfolder.UpdateRequest{Name: model.Value("Группа 1"), Description: model.Value(""), Code: model.Value("groupCode"),
		ExternalCode: model.Value("extGroupCode"), Archived: model.Value(true), Shared: model.Value(false), Owner: model.Null[model.Reference](), Group: model.Value(group)}
)

// read.json — пример «Получить Группу техкарт» источника без изменений.
func TestFixture(t *testing.T) {
	data := read(t, "custom/processingplanfolder/read.json")
	updated := time.Date(2022, 4, 19, 9, 48, 28, 563_000_000, time.FixedZone("MSK", 3*3600))
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got processingplanfolder.ProcessingPlanFolder
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.AccountID != "da7d9bbe-2c97-11e6-8a84-bae500000001" || got.Name != "Группа Мебель" ||
				got.ExternalCode != "furnitureCode" || got.Archived || got.Shared || got.PathName != "" ||
				got.Meta.Type != "processingplanfolder" || got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/entity/processingplanfolder/"+fixtureID ||
				got.Owner == nil || got.Owner.Meta.Href != owner.Meta.Href || got.Group == nil || got.Group.Meta.Href != group.Meta.Href ||
				!time.Time(got.Updated).Equal(updated) {
				t.Fatalf("processingplanfolder = %+v", got)
			}
			// Поля, которых нет в примере: code, description, непустые pathName и archived.
			var other processingplanfolder.ProcessingPlanFolder
			if err := codec.Unmarshal([]byte(`{"id":"x","code":"007","description":"d","pathName":"Корень/Мебель","archived":true,"shared":true,"futureField":{"a":1}}`), &other); err != nil ||
				other.Code != "007" || other.Description != "d" || other.PathName != "Корень/Мебель" || !other.Archived || !other.Shared {
				t.Fatalf("other = %+v, err = %v", other, err)
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
			{"create full", full, read(t, "custom/processingplanfolder/create-full.golden.json")},
			{"update", update, read(t, "custom/processingplanfolder/update.golden.json")},
			{"create minimal", processingplanfolder.CreateRequest{Name: "Группа 1"}, []byte(`{"name":"Группа 1"}`)},
			{"update archived false", processingplanfolder.UpdateRequest{Archived: model.Value(false)}, []byte(`{"archived":false}`)},
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
	// Readonly по таблице и родительская группа (её нет в таблице) не пишутся.
	for _, typ := range []reflect.Type{reflect.TypeFor[processingplanfolder.CreateRequest](), reflect.TypeFor[processingplanfolder.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "PathName", "ProcessingPlanFolder", "Parent"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains %s", typ.Name(), name)
			}
		}
	}
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/processingplanfolder/` + testID + `"},"id":"` + testID + `","name":"n","pathName":"Корень"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*processingplanfolder.ProcessingPlanFolder, error)
	}{
		{"get", "GET", path, nil, func(ctx context.Context, c *transport.Client) (*processingplanfolder.ProcessingPlanFolder, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/processingplanfolder", read(t, "custom/processingplanfolder/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*processingplanfolder.ProcessingPlanFolder, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", path, read(t, "custom/processingplanfolder/update.golden.json"), func(ctx context.Context, c *transport.Client) (*processingplanfolder.ProcessingPlanFolder, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *processingplanfolder.ProcessingPlanFolder) {
						if got.ID != testID || got.Name != "n" || got.PathName != "Корень" {
							t.Fatalf("processingplanfolder = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Delete и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/processingplanfolder"
	itemMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"processingplanfolder","mediaType":"application/json"}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", path, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"create many", "POST", "entity/processingplanfolder", `[{"name":"Группа Мебель","externalCode":"furnitureCode"}]`, `[{"id":"` + testID + `","name":"Группа Мебель"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []processingplanfolder.CreateRequest{{Name: "Группа Мебель", ExternalCode: model.Value("furnitureCode")}})
				if err == nil && (len(r) != 1 || r[0].Value == nil || r[0].Value.ID != testID || r[0].Value.Name != "Группа Мебель") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/processingplanfolder", `[{"meta":` + itemMeta + `,"name":"Группа Авто"}]`, `[{"id":"` + testID + `","name":"Группа Авто"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[processingplanfolder.UpdateRequest]{{ID: testID, Request: processingplanfolder.UpdateRequest{Name: model.Value("Группа Авто")}}})
				if err == nil && (len(r) != 1 || r[0].Value == nil || r[0].Value.Name != "Группа Авто") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/processingplanfolder/delete", `[{"meta":` + itemMeta + `}]`, `[{"info":"удалена"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).DeleteMany(ctx, []string{testID})
				if err == nil && (len(r) != 1 || r[0].Info != "удалена") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, st := range []struct {
				status int
				ok     bool
			}{{200, true}, {204, tc.empty}, {201, false}, {403, false}, {503, false}} {
				var calls atomic.Int32
				var server *httptest.Server
				server = memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					if r.Method != tc.method || r.URL.Path != "/api/remap/1.2/"+tc.path {
						t.Errorf("request = %s %s", r.Method, r.URL)
					}
					if tc.body == "" {
						if len(b) != 0 {
							t.Errorf("body = %s", b)
						}
					} else {
						documenttest.JSONEqual(t, b, []byte(strings.ReplaceAll(tc.body, "{base}", server.URL+"/api/remap/1.2")))
					}
					w.WriteHeader(st.status)
					switch st.status {
					case 200:
						io.WriteString(w, tc.response)
					case 403:
						io.WriteString(w, `{"errors":[{"error":"forbidden","code":1016}]}`)
					}
				}))
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				err = tc.invoke(t.Context(), client)
				server.Close()
				if (err == nil) != st.ok || calls.Load() != 1 {
					t.Fatalf("status %d: err = %v, calls = %d", st.status, err, calls.Load())
				}
				var apiErr *sdkerrors.APIError
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingplanfolder") ||
					st.status == 403 && !sdkerrors.HasCode(err, 1016) || errors.As(err, &apiErr) != (st.status >= 400)) {
					t.Errorf("status %d: error = %v", st.status, err)
				}
			}
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
				if err := tc.invoke(ctx, c); err != nil {
					return nil, err
				}
				return &struct{}{}, nil
			})
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
			_, err := s.Create(ctx, processingplanfolder.CreateRequest{Code: model.Value("x")}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, processingplanfolder.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingplanfolder.UpdateRequest]{{ID: "bad", Request: update}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingplanfolder.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingplanfolder") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := processingplanfolder.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "processingplanfolder", read(t, "custom/processingplanfolder/read.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[processingplanfolder.ProcessingPlanFolder] {
		return service(c).List
	}, func(t *testing.T, got processingplanfolder.ProcessingPlanFolder) {
		if got.ID != fixtureID || got.ExternalCode != "furnitureCode" {
			t.Fatalf("row = %+v", got)
		}
	})
}

// Метаданные читает общий пакет metadata; metadata.json — пример «Метаданные Групп техкарт» источника.
func TestMetadata(t *testing.T) {
	data := read(t, "custom/processingplanfolder/metadata.json")
	invoke := func(ctx context.Context, c *transport.Client) (*metadata.Entity, error) {
		s, err := metadata.NewService(c)
		if err != nil {
			return nil, err
		}
		return s.Get(ctx, "processingplanfolder")
	}
	for _, tc := range documenttest.ResponseCases(string(data)) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/processingplanfolder/metadata", nil, invoke, func(got *metadata.Entity) {
				if got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/entity/processingplanfolder/metadata" || got.Meta.MediaType != "application/json" ||
					got.Attributes != nil || len(got.States) != 0 || got.CreateShared {
					t.Fatalf("metadata = %+v", got)
				}
				documenttest.JSONEqual(t, got.Raw, data)
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
