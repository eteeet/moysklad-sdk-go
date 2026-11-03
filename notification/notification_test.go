package notification_test

import (
	"compress/gzip"
	"context"
	"encoding/json/v2"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/notification"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const id = "bcd815b9-56ca-11e9-c0a8-100b00000001"

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestGetFixture(t *testing.T) {
	data := readFile(t, "../testdata/custom/notification/get.json")
	moscow := time.FixedZone("MSK", 3*3600)
	for _, codec := range documenttest.Codecs {
		var got notification.Notification
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.Meta.Type != "NotificationTaskAssigned" || got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/notification/"+id ||
			got.ID != id || got.AccountID != "c45c23d5-5640-11e9-c0a8-100a00000001" ||
			!time.Time(got.Created).Equal(time.Date(2019, 4, 4, 14, 13, 56, 0, moscow)) || got.Read ||
			got.Title != "Новая задача: Задача" || got.Description != "Администратор назначил вам задачу Задача" {
			t.Fatalf("%s: notification = %+v", codec.Name, got)
		}
		documenttest.JSONEqual(t, got.Raw, data)
	}
}

func TestListFixture(t *testing.T) {
	data := readFile(t, "../testdata/custom/notification/list.json")
	for _, codec := range documenttest.Codecs {
		var got model.Collection[notification.Notification]
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.Meta.Size != 3 || got.Meta.Limit != 1000 || len(got.Rows) != 4 {
			t.Fatalf("%s: meta = %+v, rows = %d", codec.Name, got.Meta, len(got.Rows))
		}
		var types []model.MetaType
		for _, row := range got.Rows {
			types = append(types, row.Meta.Type)
		}
		want := []model.MetaType{"NotificationGoodCountTooLow", "NotificationImportCompleted", "RetireOrderByDemandNotificationCompleted", "NotificationTaskChanged"}
		if !equal(types, want) {
			t.Fatalf("%s: types = %v", codec.Name, types)
		}
		// Поля конкретного типа — только в Raw; long не проходит через float64.
		var retire struct {
			TaskState        string         `json:"taskState"`
			StartMoment      model.DateTime `json:"startMoment"`
			CreatedDocsCount int64          `json:"createdDocsCount"`
		}
		if err := json.Unmarshal(got.Rows[2].Raw, &retire); err != nil {
			t.Fatal(err)
		}
		if retire.TaskState != "COMPLETED" || retire.CreatedDocsCount != 1 || retire.StartMoment.IsZero() || !got.Rows[2].Read {
			t.Fatalf("%s: retire = %+v", codec.Name, retire)
		}
		var good struct {
			Good           model.Reference `json:"good"`
			ActualBalance  float64         `json:"actualBalance"`
			MinimumBalance float64         `json:"minimumBalance"`
		}
		if err := json.Unmarshal(got.Rows[0].Raw, &good); err != nil {
			t.Fatal(err)
		}
		if good.Good.Meta.Type != "product" || good.ActualBalance != 0 || good.MinimumBalance != 1 {
			t.Fatalf("%s: good = %+v", codec.Name, good)
		}
	}
}

func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Детали типа из _notification_types_details читаются из Raw в структуру вызывающего.
func TestTaskChangedRaw(t *testing.T) {
	data := readFile(t, "../testdata/upstream/notification_task_changed.json")
	type change struct {
		OldValue *string `json:"oldValue"`
		NewValue *string `json:"newValue"`
	}
	for _, codec := range documenttest.Codecs {
		var got notification.Notification
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.Meta.Type != "NotificationTaskChanged" || got.ID != "a805beb5-5adf-11e9-9bea-3ff700000025" || !got.Read ||
			got.Title != "Задача изменена: Новое описание задачи" {
			t.Fatalf("%s: notification = %+v", codec.Name, got)
		}
		var details struct {
			PerformedBy model.Reference `json:"performedBy"`
			Task        struct {
				model.Reference
				Deadline model.DateTime `json:"deadline"`
			} `json:"task"`
			Diff map[string]change `json:"diff"`
		}
		if err := json.Unmarshal(got.Raw, &details); err != nil {
			t.Fatal(err)
		}
		if details.PerformedBy.Meta.Type != "employee" || details.Task.Meta.Type != "task" || details.Task.Deadline.IsZero() ||
			details.Diff["description"].OldValue == nil || *details.Diff["description"].OldValue != "Старое описание задачи" ||
			details.Diff["deadline"].OldValue != nil || *details.Diff["deadline"].NewValue != "2019-04-02 18:53:00.000" {
			t.Fatalf("%s: details = %+v", codec.Name, details)
		}
	}
}

// Неизвестный тип уведомления не отвергается; Raw — копия, а не ссылка на буфер.
func TestUnknownTypeAndRawCopy(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		data := []byte(`{"meta":{"type":"FutureNotification"},"id":"` + id + `","read":false,"extra":{"n":12345678901234567890}}`)
		var got notification.Notification
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		data[0] = 'X'
		if got.Meta.Type != "FutureNotification" || got.ID != id || got.Raw[0] != '{' ||
			!strings.Contains(string(got.Raw), `12345678901234567890`) {
			t.Fatalf("%s: got = %+v raw = %s", codec.Name, got, got.Raw)
		}
		out, err := codec.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(out), "extra") || strings.Contains(string(out), "Raw") {
			t.Fatalf("%s: marshal = %s", codec.Name, out)
		}
		if err := codec.Unmarshal([]byte(`{"read":"yes"}`), &got); err == nil {
			t.Fatalf("%s: string read accepted", codec.Name)
		}
	}
}

type response struct {
	name   string
	status int
	body   string
	gzip   bool
	api    bool // ошибка — *sdkerrors.APIError с этим статусом
	ok     bool
}

func errorResponses() []response {
	return []response{
		{"bad request", 400, `{"errors":[{"error":"constraint","code":22004}]}`, false, true, false},
		{"not found", 404, `{"errors":[{"error":"missing","code":1001}]}`, false, true, false},
		{"forbidden gzip", 403, `{"errors":[{"error":"denied","code":1016}]}`, true, true, false},
		{"rate limit", 429, `{"errors":[{"error":"slow","code":1073}]}`, false, true, false},
		{"server error", 500, `{"errors":[{"error":"server","code":1000}]}`, false, true, false},
		{"html error", 502, `<html>gateway</html>`, false, true, false},
		{"unexpected status", 201, `{}`, false, false, false},
	}
}

// serve отвечает одним ответом и проверяет точный запрос: method, path, query, заголовки и тело.
func serve(t *testing.T, method, path, query, body string, r response) (*notification.Service, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		calls.Add(1)
		if req.Method != method || req.URL.Path != "/api/remap/1.2/"+path || req.URL.RawQuery != query {
			t.Errorf("request = %s %s", req.Method, req.URL.RequestURI())
		}
		for key, want := range map[string]string{"Authorization": "Bearer test-token", "Accept": "application/json;charset=utf-8",
			"Accept-Encoding": "gzip", "User-Agent": "moysklad-sdk-go"} {
			if req.Header.Get(key) != want {
				t.Errorf("header %s = %q", key, req.Header.Get(key))
			}
		}
		data, err := io.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		if body == "" {
			if len(data) != 0 || req.Header.Get("Content-Type") != "" {
				t.Errorf("unexpected body %s, content-type %q", data, req.Header.Get("Content-Type"))
			}
		} else {
			documenttest.JSONEqual(t, data, []byte(body))
			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("content-type = %q", req.Header.Get("Content-Type"))
			}
		}
		if r.gzip {
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(r.status)
			zw := gzip.NewWriter(w)
			zw.Write([]byte(r.body))
			zw.Close()
			return
		}
		w.WriteHeader(r.status)
		io.WriteString(w, r.body)
	}))
	t.Cleanup(server.Close)
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	return newService(t, client), &calls
}

func newService(t *testing.T, client *transport.Client) *notification.Service {
	t.Helper()
	s, err := notification.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// checkResult проверяет итог одного вызова: одна отправка, ResponseInfo, ошибка с одним
// префиксом moysklad: и контекстом операции.
func checkResult(t *testing.T, r response, operation string, calls *atomic.Int32, info transport.ResponseInfo, err error) {
	t.Helper()
	if calls.Load() != 1 || info.StatusCode != r.status || info.Attempts != 1 {
		t.Fatalf("calls = %d, info = %+v, err = %v", calls.Load(), info, err)
	}
	if r.ok {
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		return
	}
	if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), operation) {
		t.Fatalf("err = %v", err)
	}
	var apiErr *sdkerrors.APIError
	if errors.As(err, &apiErr) != r.api || r.api && apiErr.StatusCode != r.status {
		t.Fatalf("API error = %v", err)
	}
	if r.status == 400 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 22004) {
		t.Errorf("API details = %+v", apiErr)
	}
}

func TestList(t *testing.T) {
	page := string(readFile(t, "../testdata/custom/notification/list.json"))
	cases := []struct {
		response
		options notification.ListOptions
		query   string
	}{
		{response{"defaults", 200, page, false, false, true}, notification.ListOptions{}, ""},
		{response{"limit offset gzip", 200, page, true, false, true}, notification.ListOptions{Limit: 1000, Offset: 40}, "limit=1000&offset=40"},
		{response{"empty page", 200, `{"meta":{"href":"h","size":0},"rows":[]}`, false, false, true}, notification.ListOptions{Limit: 1}, "limit=1"},
		{response{"null rows", 200, `{"meta":{"href":"h"},"rows":null}`, false, false, false}, notification.ListOptions{}, ""},
		{response{"no meta", 200, `{"rows":[]}`, false, false, false}, notification.ListOptions{}, ""},
		{response{"malformed", 200, `{"meta":`, false, false, false}, notification.ListOptions{}, ""},
	}
	for _, r := range errorResponses() {
		cases = append(cases, struct {
			response
			options notification.ListOptions
			query   string
		}{r, notification.ListOptions{}, ""})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, calls := serve(t, "GET", "notification", tc.query, "", tc.response)
			var info transport.ResponseInfo
			got, err := s.List(t.Context(), tc.options, transport.WithResponseInfo(&info))
			checkResult(t, tc.response, "list notifications", calls, info, err)
			if !tc.ok {
				if got != nil {
					t.Fatalf("got = %+v", got)
				}
				return
			}
			if tc.body == page && (len(got.Rows) != 4 || got.Rows[1].Meta.Type != "NotificationImportCompleted" ||
				got.Rows[1].ID != "0f423542-5b7a-11e9-9bea-3ff70000000f" || len(got.Rows[3].Raw) == 0) {
				t.Fatalf("rows = %+v", got.Rows)
			}
		})
	}
}

func TestGet(t *testing.T) {
	body := string(readFile(t, "../testdata/custom/notification/get.json"))
	cases := []response{
		{"plain", 200, body, false, false, true},
		{"gzip", 200, body, true, false, true},
		{"malformed", 200, `{"id":`, false, false, false},
		{"null", 200, `null`, false, false, false},
		{"empty body", 200, ``, false, false, false},
		{"missing id", 200, `{"meta":{"type":"NotificationTaskAssigned"}}`, false, false, false},
		{"no content", 204, ``, false, false, false},
	}
	for _, r := range append(cases, errorResponses()...) {
		t.Run(r.name, func(t *testing.T) {
			s, calls := serve(t, "GET", "notification/"+id, "", "", r)
			var info transport.ResponseInfo
			got, err := s.Get(t.Context(), id, transport.WithResponseInfo(&info))
			checkResult(t, r, "get notification "+id, calls, info, err)
			if !r.ok {
				if got != nil {
					t.Fatalf("got = %+v", got)
				}
				return
			}
			if got.ID != id || got.Meta.Type != "NotificationTaskAssigned" || got.Read || !strings.Contains(string(got.Raw), `"task"`) {
				t.Fatalf("got = %+v", got)
			}
		})
	}
}

// Операции без тела ответа: Delete, MarkAsRead, MarkAllAsRead.
func TestNoBodyOperations(t *testing.T) {
	ops := []struct {
		operation, method, path string
		invoke                  func(*notification.Service, ...transport.RequestOption) error
	}{
		{"delete notification " + id, "DELETE", "notification/" + id, func(s *notification.Service, o ...transport.RequestOption) error {
			return s.Delete(t.Context(), id, o...)
		}},
		{"mark notification " + id + " as read", "PUT", "notification/" + id + "/markasread", func(s *notification.Service, o ...transport.RequestOption) error {
			return s.MarkAsRead(t.Context(), id, o...)
		}},
		{"mark all notifications as read", "PUT", "notification/markasreadall", func(s *notification.Service, o ...transport.RequestOption) error {
			return s.MarkAllAsRead(t.Context(), o...)
		}},
	}
	cases := []response{
		{"empty 200", 200, ``, false, false, true},
		{"body ignored", 200, `not json`, false, false, true},
		{"gzip", 200, `{}`, true, false, true},
		{"no content", 204, ``, false, false, true},
	}
	for _, op := range ops {
		for _, r := range append(cases, errorResponses()...) {
			t.Run(op.method+" "+op.path+" "+r.name, func(t *testing.T) {
				s, calls := serve(t, op.method, op.path, "", "", r)
				var info transport.ResponseInfo
				err := op.invoke(s, transport.WithResponseInfo(&info))
				checkResult(t, r, op.operation, calls, info, err)
			})
		}
		t.Run(op.path+" local errors", func(t *testing.T) {
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
				s := newService(t, c)
				var err error
				switch op.method {
				case "DELETE":
					err = s.Delete(ctx, id)
				default:
					if strings.HasSuffix(op.path, "markasreadall") {
						err = s.MarkAllAsRead(ctx)
					} else {
						err = s.MarkAsRead(ctx, id)
					}
				}
				if err != nil {
					return nil, err
				}
				return &struct{}{}, nil
			})
		})
		// Запись не повторяется даже при RetryPolicy клиента.
		t.Run(op.path+" no retry", func(t *testing.T) {
			var calls atomic.Int32
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.WriteHeader(http.StatusServiceUnavailable)
			}))
			defer server.Close()
			client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{}})
			if err != nil {
				t.Fatal(err)
			}
			if err := op.invoke(newService(t, client)); sdkerrors.StatusCode(err) != 503 || calls.Load() != 1 {
				t.Fatalf("err = %v, calls = %d", err, calls.Load())
			}
		})
	}
}

func TestLocalValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := newService(t, client)
	for _, bad := range []string{"", "bcd815b9", "bcd815b9-56ca-11e9-c0a8-100b0000000z", "bcd815b9x56ca-11e9-c0a8-100b00000001", "../" + id[3:]} {
		for name, invoke := range map[string]func(...transport.RequestOption) error{
			"get": func(o ...transport.RequestOption) error {
				got, err := s.Get(t.Context(), bad, o...)
				if got != nil {
					t.Errorf("got = %+v", got)
				}
				return err
			},
			"delete":     func(o ...transport.RequestOption) error { return s.Delete(t.Context(), bad, o...) },
			"markasread": func(o ...transport.RequestOption) error { return s.MarkAsRead(t.Context(), bad, o...) },
		} {
			info := transport.ResponseInfo{StatusCode: 503, Attempts: 2}
			if err := invoke(transport.WithResponseInfo(&info)); err == nil || strings.Count(err.Error(), "moysklad:") != 1 ||
				!strings.Contains(err.Error(), "UUID") || info.StatusCode != 0 || info.Attempts != 0 {
				t.Errorf("%s %q: err = %v, info = %+v", name, bad, err, info)
			}
		}
	}
	for _, options := range []notification.ListOptions{{Limit: 1001}, {Limit: -1}, {Offset: -1}} {
		info := transport.ResponseInfo{StatusCode: 503}
		got, err := s.List(t.Context(), options, transport.WithResponseInfo(&info))
		if got != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "list notifications") || info.StatusCode != 0 {
			t.Errorf("%+v: got = %+v, err = %v, info = %+v", options, got, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	for name, invoke := range map[string]func(context.Context, *transport.Client) (*struct{}, error){
		"list": func(ctx context.Context, c *transport.Client) (*struct{}, error) {
			_, err := newService(t, c).List(ctx, notification.ListOptions{})
			return nil, err
		},
		"get": func(ctx context.Context, c *transport.Client) (*struct{}, error) {
			_, err := newService(t, c).Get(ctx, id)
			return nil, err
		},
	} {
		t.Run(name, func(t *testing.T) { documenttest.LocalErrors(t, invoke) })
	}
	if _, err := notification.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}
