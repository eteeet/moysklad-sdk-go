package processingprocess_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/processingprocess"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "117cae13-a612-11ed-ac12-000900000022"
	fixtureID  = "d5174779-862b-11eb-ac14-000900000007"
	positionID = "23a62e19-a6bb-11ed-ac12-000900000043"
	docPath    = "entity/processingprocess/" + testID
	api        = "https://api.moysklad.ru/api/remap/1.2/"
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
	return model.Reference{Meta: model.Meta{Href: api + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *processingprocess.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := processingprocess.NewService(c)
	return s
}

var (
	stageA  = ref("entity/processingstage/364b1107-9bd3-11ed-ac12-000c0000006a", "processingstage")
	stageB  = ref("entity/processingstage/c18373c0-9aea-11ed-ac12-000e000000c2", "processingstage")
	nextRef = ref(docPath+"/positions/"+positionID, "processingprocessposition")
	rowMeta = func(id string) model.Meta {
		return model.Meta{Href: api + docPath + "/positions/" + id, Type: "processingprocessposition", MediaType: "application/json"}
	}
	minimal    = processingprocess.CreateRequest{Name: "Изготавливаем двигатель", Positions: []processingprocess.PositionRequest{{ProcessingStage: stageA}}}
	fullCreate = processingprocess.CreateRequest{
		Name: "Изготавливаем двигатель",
		// Порядок этапов — порядок среза: A, затем B.
		Positions:   []processingprocess.PositionRequest{{ProcessingStage: stageA}, {ProcessingStage: stageB, NextPositions: model.Value([]model.Reference{})}},
		Description: model.Value(""), ExternalCode: model.Value("ext-1"), Archived: model.Value(false), Shared: model.Value(false),
		Owner: model.Value(ref("entity/employee/c0f98b1e-9aea-11ed-ac12-000e00000050", "employee")),
		Group: model.Value(ref("entity/group/c0b230f9-9aea-11ed-ac12-000b00000012", "group")),
	}
	fullUpdate = processingprocess.UpdateRequest{
		Name: model.Value("Изготавливаем двигатель №2"), ExternalCode: model.Value("dfsafsfsd1231231"), Archived: model.Value(false),
		Description: model.Null[string](), Owner: model.Null[model.Reference](),
		Positions: model.Value([]processingprocess.PositionUpdateRequest{
			{Meta: model.Value(rowMeta("117cb64b-a612-11ed-ac12-000900000024")), ProcessingStage: model.Value(stageB), NextPositions: model.Value([]model.Reference{nextRef})},
			{ProcessingStage: model.Value(stageA)},
			{Meta: model.Value(rowMeta(positionID)), NextPositions: model.Value([]model.Reference{})},
		}),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "custom/processingprocess/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got processingprocess.ProcessingProcess
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.AccountID != "dbb8cfc1-cbfa-11e1-6dfb-889ffa6f49fd" || got.Name != "Основной техпроцесс" || got.ExternalCode != "F1l43a3ojXZShfnzJCKsG3" ||
				got.Archived || !got.Shared || got.Meta.Type != "processingprocess" || got.Description != "" {
				t.Fatalf("process = %+v", got)
			}
			if !time.Time(got.Updated).Equal(time.Time(documenttest.DateTime("2021-03-16 10:47:18.160"))) {
				t.Errorf("updated = %v", got.Updated)
			}
			if got.Owner == nil || !strings.HasSuffix(got.Owner.Meta.Href, "/employee/872559f1-cbf3-11e1-9eb9-889ffa6f49fd") || got.Group == nil || got.Group.Meta.Type != "group" {
				t.Fatalf("references = %+v / %+v", got.Owner, got.Group)
			}
			if got.Positions == nil || got.Positions.Meta.Type != "processingprocessposition" || got.Positions.Meta.Size != 1 || got.Positions.Rows != nil {
				t.Fatalf("positions = %+v", got.Positions)
			}
			// Поля таблицы, которых нет в примере, и строки позиций при expand: порядок как в ответе.
			var sparse processingprocess.ProcessingProcess
			if err := codec.Unmarshal([]byte(`{"id":"x","description":"d","positions":{"meta":{"href":"h","size":2},"rows":[`+
				`{"id":"p1","processingstage":{"meta":{"type":"processingstage","href":"s1"}},"nextPositions":[{"meta":{"href":"p2","type":"processingprocessposition"}}]},`+
				`{"id":"p2","processingstage":{"meta":{"href":"s2"}},"nextPositions":[]}]},"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			rows := sparse.Positions.Rows
			if sparse.Description != "d" || len(rows) != 2 || rows[0].ID != "p1" || rows[1].ID != "p2" || rows[0].ProcessingStage.Meta.Href != "s1" ||
				len(rows[0].NextPositions) != 1 || rows[0].NextPositions[0].Meta.Href != "p2" || len(rows[1].NextPositions) != 0 || sparse.Owner != nil || !sparse.Updated.IsZero() {
				t.Fatalf("sparse = %+v", sparse)
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
			{"create full", fullCreate, read(t, "custom/processingprocess/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"name":"Изготавливаем двигатель","positions":[{"processingstage":{"meta":{"href":"` + stageA.Meta.Href + `","type":"processingstage"}}}]}`)},
			{"update", fullUpdate, read(t, "custom/processingprocess/update.golden.json")},
			{"position create", processingprocess.PositionRequest{ProcessingStage: stageB, NextPositions: model.Value([]model.Reference{nextRef})},
				[]byte(`{"processingstage":{"meta":{"href":"` + stageB.Meta.Href + `","type":"processingstage"}},"nextPositions":[{"meta":{"href":"` + nextRef.Meta.Href + `","type":"processingprocessposition"}}]}`)},
			{"position update", processingprocess.PositionUpdateRequest{NextPositions: model.Null[[]model.Reference]()}, []byte(`{"nextPositions":null}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[processingprocess.CreateRequest](), reflect.TypeFor[processingprocess.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	// meta строки — только адрес существующей позиции в Update техпроцесса.
	for _, typ := range []reflect.Type{reflect.TypeFor[processingprocess.PositionRequest](), reflect.TypeFor[processingprocess.PositionUpdateRequest]()} {
		for _, name := range []string{"ID", "AccountID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	if _, ok := reflect.TypeFor[processingprocess.PositionRequest]().FieldByName("Meta"); ok {
		t.Error("new position must not carry meta")
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и размер набора — только из ответа.
	const response = `{"id":"server-id","name":"server-name","archived":true,"positions":{"meta":{"href":"https://example.test/positions","type":"processingprocessposition","size":3}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*processingprocess.ProcessingProcess, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*processingprocess.ProcessingProcess, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/processingprocess", read(t, "custom/processingprocess/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*processingprocess.ProcessingProcess, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/processingprocess/update.golden.json"), func(ctx context.Context, c *transport.Client) (*processingprocess.ProcessingProcess, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *processingprocess.ProcessingProcess) {
						if got.ID != "server-id" || got.Name != "server-name" || !got.Archived || got.Positions.Meta.Size != 3 || got.Positions.Rows != nil {
							t.Fatalf("process = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func httpCases[T any](t *testing.T, method, path, body, response string, invoke func(context.Context, *transport.Client) (*T, error), check func(*T)) {
	t.Helper()
	var sent []byte
	if body != "" {
		sent = []byte(body)
	}
	for _, tc := range documenttest.ResponseCases(response) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, method, path, sent, invoke, check)
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestPositionsHTTP(t *testing.T) {
	// Пример ответа «Получить отдельную позицию» плюс nextPositions из таблицы.
	const row = `{"meta":{"href":"` + api + docPath + `/positions/` + positionID + `","type":"processingprocessposition","mediaType":"application/json"},` +
		`"id":"` + positionID + `","accountId":"c0b1ef18-9aea-11ed-ac12-000b00000011",` +
		`"processingstage":{"meta":{"href":"` + api + `entity/processingstage/4b5662f4-9bd3-11ed-ac12-000c00000070","type":"processingstage"}},` +
		`"nextPositions":[{"meta":{"href":"` + api + docPath + `/positions/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","type":"processingprocessposition"}},{"meta":{"href":"` + api + docPath + `/positions/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb","type":"processingprocessposition"}}]}`
	checkRow := func(got *processingprocess.Position) {
		if got.ID != positionID || got.AccountID != "c0b1ef18-9aea-11ed-ac12-000b00000011" || got.Meta.Type != "processingprocessposition" || got.ProcessingStage == nil ||
			!strings.HasSuffix(got.ProcessingStage.Meta.Href, "/4b5662f4-9bd3-11ed-ac12-000c00000070") || len(got.NextPositions) != 2 ||
			!strings.HasSuffix(got.NextPositions[0].Meta.Href, "/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa") || !strings.HasSuffix(got.NextPositions[1].Meta.Href, "/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb") {
			t.Fatalf("row = %+v", got)
		}
	}
	t.Run("get positions", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions?limit=2&offset=7", "", `{"meta":{"href":"h","size":1},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[processingprocess.Position], error) {
				return service(c).GetPositions(ctx, testID, query.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[processingprocess.Position]) { checkRow(&got.Rows[0]) })
	})
	t.Run("get position", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions/"+positionID, "", row, func(ctx context.Context, c *transport.Client) (*processingprocess.Position, error) {
			return service(c).GetPosition(ctx, testID, positionID)
		}, checkRow)
	})
	t.Run("update position", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/positions/"+positionID, `{"processingstage":{"meta":{"href":"`+stageA.Meta.Href+`","type":"processingstage"}},"nextPositions":[]}`, row,
			func(ctx context.Context, c *transport.Client) (*processingprocess.Position, error) {
				return service(c).UpdatePosition(ctx, testID, positionID, processingprocess.PositionUpdateRequest{ProcessingStage: model.Value(stageA), NextPositions: model.Value([]model.Reference{})})
			}, checkRow)
	})
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/processingprocess"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"processingprocess","mediaType":"application/json"}`
	stageJSON := `{"meta":{"href":"` + stageA.Meta.Href + `","type":"processingstage"}}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", docPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"create many", "POST", "entity/processingprocess", `[{"name":"Изготавливаем двигатель","positions":[{"processingstage":` + stageJSON + `}]}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []processingprocess.CreateRequest{minimal})
				if err == nil && r[0].Value.ID != "x" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/processingprocess", `[{"meta":` + docMeta + `,"name":"Штамповка"}]`, `[{"id":"x","name":"Штамповка"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[processingprocess.UpdateRequest]{{ID: testID, Request: processingprocess.UpdateRequest{Name: model.Value("Штамповка")}}})
				if err == nil && r[0].Value.Name != "Штамповка" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/processingprocess/delete", `[{"meta":` + docMeta + `}]`, `[{"info":"ok"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).DeleteMany(ctx, []string{testID})
				if err == nil && r[0].Info != "ok" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"create positions", "POST", docPath + "/positions", `[{"processingstage":` + stageJSON + `},{"processingstage":{"meta":{"href":"` + stageB.Meta.Href + `","type":"processingstage"}},"nextPositions":[]}]`,
			`[{"id":"` + positionID + `"},{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreatePositions(ctx, testID, []processingprocess.PositionRequest{{ProcessingStage: stageA}, {ProcessingStage: stageB, NextPositions: model.Value([]model.Reference{})}})
				if err == nil && (r[0].Value.ID != positionID || r[1].Value.ID != "x") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete position", "DELETE", docPath + "/positions/" + positionID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, positionID)
		}},
		{"delete many positions", "POST", docPath + "/positions/delete", `[{"meta":{"href":"` + base + `/` + testID + `/positions/` + positionID + `","type":"processingprocessposition","mediaType":"application/json"}}]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyPositions(ctx, testID, []string{positionID})
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingprocess") ||
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

// Ошибки входа — до HTTP, ResponseInfo сброшен; ровно 100 позиций уходят.
func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		io.WriteString(w, `{"id":"server-id"}`)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	newRows := func(n int) []processingprocess.PositionRequest {
		rows := make([]processingprocess.PositionRequest, n)
		for i := range rows {
			rows[i].ProcessingStage = stageA
		}
		return rows
	}
	updateRows := func(n int) model.Field[[]processingprocess.PositionUpdateRequest] {
		rows := make([]processingprocess.PositionUpdateRequest, n)
		for i := range rows {
			rows[i].Meta = model.Value(rowMeta(positionID))
		}
		return model.Value(rows)
	}
	limit := minimal
	limit.Positions = newRows(100)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, processingprocess.UpdateRequest{Positions: updateRows(100)}); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	create := func(change func(*processingprocess.CreateRequest)) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			input := minimal
			change(&input)
			_, err := s.Create(ctx, input, o...)
			return err
		}
	}
	update := func(input processingprocess.UpdateRequest) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, input, o...)
			return err
		}
	}
	calls.Store(0)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":      func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":        func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no name":  create(func(r *processingprocess.CreateRequest) { r.Name = "" }),
		"create no rows":  create(func(r *processingprocess.CreateRequest) { r.Positions = nil }),
		"create 101 rows": create(func(r *processingprocess.CreateRequest) { r.Positions = newRows(101) }),
		"create no stage": create(func(r *processingprocess.CreateRequest) {
			r.Positions = []processingprocess.PositionRequest{{ProcessingStage: stageA}, {}}
		}),
		"update empty":        update(processingprocess.UpdateRequest{}),
		"update bad id":       func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"update null rows":    update(processingprocess.UpdateRequest{Positions: model.Null[[]processingprocess.PositionUpdateRequest]()}),
		"update zero rows":    update(processingprocess.UpdateRequest{Positions: updateRows(0)}),
		"update 101 rows":     update(processingprocess.UpdateRequest{Positions: updateRows(101)}),
		"update new no stage": update(processingprocess.UpdateRequest{Positions: model.Value([]processingprocess.PositionUpdateRequest{{NextPositions: model.Value([]model.Reference{})}})}),
		"delete bad id":       func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"create many empty":   func(o ...transport.RequestOption) error { _, err := s.CreateMany(ctx, nil, o...); return err },
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingprocess.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingprocess.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"positions bad id": func(o ...transport.RequestOption) error {
			_, err := s.GetPositions(ctx, "bad", query.PageOptions{}, o...)
			return err
		},
		"positions fields": func(o ...transport.RequestOption) error {
			_, err := s.GetPositions(ctx, testID, query.PageOptions{Fields: "stock"}, o...)
			return err
		},
		"positions limit": func(o ...transport.RequestOption) error {
			_, err := s.GetPositions(ctx, testID, query.PageOptions{Limit: 1001}, o...)
			return err
		},
		"get position bad": func(o ...transport.RequestOption) error {
			_, err := s.GetPosition(ctx, testID, "bad", o...)
			return err
		},
		"create positions empty": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, testID, nil, o...)
			return err
		},
		"create positions bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, "bad", newRows(1), o...)
			return err
		},
		"create positions no stage": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, testID, []processingprocess.PositionRequest{{}}, o...)
			return err
		},
		"update position empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, processingprocess.PositionUpdateRequest{}, o...)
			return err
		},
		"update position meta": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, processingprocess.PositionUpdateRequest{Meta: model.Value(rowMeta(positionID)), ProcessingStage: model.Value(stageA)}, o...)
			return err
		},
		"update position bad": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, "bad", positionID, processingprocess.PositionUpdateRequest{ProcessingStage: model.Value(stageA)}, o...)
			return err
		},
		"delete position bad":         func(o ...transport.RequestOption) error { return s.DeletePosition(ctx, testID, "bad", o...) },
		"delete many positions empty": func(o ...transport.RequestOption) error { return s.DeleteManyPositions(ctx, testID, nil, o...) },
		"delete many positions 1001": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, make([]string, 1001), o...)
		},
		"delete many positions bad parent": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, "bad", []string{positionID}, o...)
		},
		"delete many positions bad row": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, []string{positionID, "bad"}, o...)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingprocess") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := processingprocess.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "processingprocess", read(t, "custom/processingprocess/read.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[processingprocess.ProcessingProcess] {
		return service(c).List
	}, func(t *testing.T, got processingprocess.ProcessingProcess) {
		if got.ID != fixtureID || got.Name != "Основной техпроцесс" || got.Positions == nil || got.Positions.Meta.Size != 1 {
			t.Fatalf("row = %+v", got)
		}
	})
}
