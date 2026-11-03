package processingstage_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/processingstage"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "d2308bcc-8fd9-11ed-ac12-000b000000c1"
	stagePath = "entity/processingstage/" + testID
	performer = "https://api.moysklad.ru/api/remap/1.2/entity/employee/83db6e80-761b-4e3d-aee8-641299ed0898"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(href, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: href, Type: model.MetaType(typ), MediaType: "application/json"}}
}

func service(c *transport.Client) *processingstage.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := processingstage.NewService(c)
	return s
}

var (
	full = processingstage.CreateRequest{
		Name: "Этап 1", Description: model.Value(""), ExternalCode: model.Value("456"),
		Archived: model.Value(false), Shared: model.Value(false),
		AllPerformers: model.Value(false), DistributionRequired: model.Value(false),
		StandardHourCost: model.Value(15050.5),
		Performers:       model.Value([]model.Reference{ref(performer, "employee")}),
		Owner:            model.Value(ref("https://api.moysklad.ru/api/remap/1.2/entity/employee/d105a6bf-8fd9-11ed-ac12-000b0000004f", "employee")),
		Group:            model.Value(ref("https://api.moysklad.ru/api/remap/1.2/entity/group/d0668856-8fd9-11ed-ac12-000e00000001", "group")),
	}
	update = processingstage.UpdateRequest{
		Name: model.Value("Этап 1.1"), AllPerformers: model.Value(false), DistributionRequired: model.Value(true),
		StandardHourCost: model.Value(0.0), Performers: model.Value([]model.Reference{}),
	}
)

func TestProcessingStageFixture(t *testing.T) {
	data := read(t, "custom/processingstage/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got processingstage.ProcessingStage
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.AccountID != "d063f3f3-8fd9-11ed-ac12-000e00000000" || got.Meta.Type != "processingstage" ||
				got.Name != "Основной этап" || got.ExternalCode != "sTV9PL-HjZkNgDMUqvKKe3" || got.Description != "" ||
				!got.Shared || got.Archived || got.AllPerformers || got.DistributionRequired || got.StandardHourCost != 0 ||
				!time.Time(got.Updated).Equal(time.Time(documenttest.DateTime("2023-01-11 09:01:18.363"))) || got.MaterialStore != nil ||
				got.Owner == nil || got.Owner.Meta.Type != "employee" || got.Group == nil || got.Group.Meta.Type != "group" {
				t.Fatalf("processingstage = %+v", got)
			}
			// allPerformers=false и непустой performers — «только сотрудники из списка».
			if len(got.Performers) != 1 || got.Performers[0].Meta.Href != "http://localhost/api/remap/1.2/entity/employee/83db6e80-761b-4e3d-aee8-641299ed0898" || got.Performers[0].Meta.Type != "employee" {
				t.Fatalf("performers = %+v", got.Performers)
			}
			var sparse processingstage.ProcessingStage
			if err := codec.Unmarshal([]byte(`{"id":"x","allPerformers":true,"performers":[],"materialStore":{"meta":{"type":"store"}},"futureField":{"a":1}}`), &sparse); err != nil ||
				sparse.ID != "x" || !sparse.AllPerformers || sparse.Performers == nil || len(sparse.Performers) != 0 || sparse.MaterialStore == nil || sparse.MaterialStore.Meta.Type != "store" {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

// Стоимость нормо-часа — копейки float64: 0, 0.0, дробь и большое значение без потерь.
func TestMoneyJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got processingstage.ProcessingStage
			if err := codec.Unmarshal([]byte(`{"standardHourCost":`+raw+`}`), &got); err != nil || got.StandardHourCost != want {
				t.Fatalf("%s: standardHourCost=%s: got %v, err = %v", codec.Name, raw, got.StandardHourCost, err)
			}
		}
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/processingstage/create-full.golden.json")},
			{"update", update, read(t, "custom/processingstage/update.golden.json")},
			{"create minimal", processingstage.CreateRequest{Name: "Э"}, []byte(`{"name":"Э"}`)},
			{"create all performers", processingstage.CreateRequest{Name: "Э", AllPerformers: model.Value(true)}, []byte(`{"name":"Э","allPerformers":true}`)},
			// Флаг и список вместе уходят как заданы: итог определяет сервер.
			{"update flag and list", processingstage.UpdateRequest{AllPerformers: model.Value(true), Performers: model.Value([]model.Reference{ref(performer, "employee")})},
				[]byte(`{"allPerformers":true,"performers":[{"meta":{"href":"` + performer + `","type":"employee","mediaType":"application/json"}}]}`)},
			{"update nil performers", processingstage.UpdateRequest{Performers: model.Value[[]model.Reference](nil)}, []byte(`{"performers":[]}`)},
			{"update explicit false", processingstage.UpdateRequest{AllPerformers: model.Value(false)}, []byte(`{"allPerformers":false}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[processingstage.CreateRequest](), reflect.TypeFor[processingstage.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "MaterialStore"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestServiceHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: флаг и исполнители — только из ответа.
	const response = `{"id":"server-id","name":"server-name","allPerformers":true,"standardHourCost":4040,"performers":[{"meta":{"href":"` + performer + `","type":"employee"}}]}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*processingstage.ProcessingStage, error)
	}{
		{"get", "GET", stagePath, nil, func(ctx context.Context, c *transport.Client) (*processingstage.ProcessingStage, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/processingstage", read(t, "custom/processingstage/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*processingstage.ProcessingStage, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", stagePath, read(t, "custom/processingstage/update.golden.json"), func(ctx context.Context, c *transport.Client) (*processingstage.ProcessingStage, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *processingstage.ProcessingStage) {
						if got.ID != "server-id" || got.Name != "server-name" || !got.AllPerformers || got.StandardHourCost != 4040 ||
							len(got.Performers) != 1 || got.Performers[0].Meta.Href != performer {
							t.Fatalf("processingstage = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/processingstage"
	stageMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"processingstage","mediaType":"application/json"}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", stagePath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"create many", "POST", "entity/processingstage", `[{"name":"Этап 2"},{"name":"Этап 3","allPerformers":false}]`,
			`[{"id":"x"},{"errors":[{"error":"bad","code":3000}]}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []processingstage.CreateRequest{{Name: "Этап 2"}, {Name: "Этап 3", AllPerformers: model.Value(false)}})
				if err == nil && (len(r) != 2 || r[0].Value == nil || r[0].Value.ID != "x" || r[1].Value != nil || len(r[1].Errors) != 1 || r[1].Errors[0].Code != 3000) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/processingstage", `[{"meta":` + stageMeta + `,"name":"Этап 1","description":"Подготовка"}]`, `[{"id":"` + testID + `"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[processingstage.UpdateRequest]{{ID: testID, Request: processingstage.UpdateRequest{Name: model.Value("Этап 1"), Description: model.Value("Подготовка")}}})
				if err == nil && (len(r) != 1 || r[0].Value == nil || r[0].Value.ID != testID) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/processingstage/delete", `[{"meta":` + stageMeta + `}]`, `[{"info":"ok"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).DeleteMany(ctx, []string{testID})
				if err == nil && (len(r) != 1 || r[0].Info != "ok") {
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingstage") ||
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
			_, err := s.Create(ctx, processingstage.CreateRequest{AllPerformers: model.Value(true)}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, processingstage.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingstage.UpdateRequest]{{ID: "bad", Request: update}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingstage.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
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
	listtest.Run(t, "processingstage", read(t, "custom/processingstage/read.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[processingstage.ProcessingStage] {
		return service(c).List
	}, func(t *testing.T, got processingstage.ProcessingStage) {
		if got.ID != testID || got.AllPerformers || len(got.Performers) != 1 {
			t.Fatalf("row = %+v", got)
		}
	})
}
