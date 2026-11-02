package async_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/async"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const taskID = "498b8673-0308-11e6-9464-e4de00000089"

var moscow = time.FixedZone("MSK", 3*60*60)

func read(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/custom/async/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func newService(t *testing.T, h http.HandlerFunc, cfg transport.Config) *async.Service {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	cfg.BaseURL = server.URL
	c, err := transport.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	s, err := async.NewService(c)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// Примеры источника: DONE с resultUrl и deletionDate, API_ERROR с errors массивом, список
// статусов; обоими кодеками.
func TestTaskFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var done, apiErr async.Task
			if err := codec.Unmarshal(read(t, "task_done.json"), &done); err != nil {
				t.Fatal(err)
			}
			if done.ID != taskID || done.State != async.StateDone || !done.State.Finished() || done.ResultURL == "" ||
				done.DeletionDate.IsZero() || done.Owner == nil || done.Owner.Meta.Type != "employee" || done.Errors != nil {
				t.Fatalf("done = %+v", done)
			}
			if err := codec.Unmarshal(read(t, "task_api_error.json"), &apiErr); err != nil {
				t.Fatal(err)
			}
			if apiErr.State != async.StateAPIError || len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 39000 || apiErr.ResultURL != "" {
				t.Fatalf("api error = %+v", apiErr)
			}
			var list model.Collection[async.Task]
			if err := codec.Unmarshal(read(t, "list.json"), &list); err != nil {
				t.Fatal(err)
			}
			deletion := time.Time(list.Rows[1].DeletionDate).In(moscow).Format("2006-01-02 15:04:05.000")
			if len(list.Rows) != 2 || list.Rows[0].State != async.StateProcessing || list.Rows[0].State.Finished() ||
				list.Rows[1].State != async.StateDone || list.Rows[0].Owner.Meta.Type != "application" || deletion != "2021-04-21 15:07:05.996" {
				t.Fatalf("list = %+v", list)
			}
		})
	}
}

func TestStateFinished(t *testing.T) {
	for state, want := range map[async.State]bool{async.StatePending: false, async.StateProcessing: false, async.StateDone: true,
		async.StateError: true, async.StateCancel: true, async.StateAPIError: true, "FUTURE": false} {
		if state.Finished() != want {
			t.Errorf("%s: %v", state, !want)
		}
	}
}

// Get и List: путь, query, ResponseInfo, APIError.
func TestGetList(t *testing.T) {
	s := newService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/async/" + taskID:
			w.Write(read(t, "task_done.json"))
		case "/async?filter=state%3DDONE&limit=10":
			w.Write(read(t, "list.json"))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"errors":[{"error":"not found","code":1021}]}`))
		}
	}, transport.Config{})
	var info transport.ResponseInfo
	task, err := s.Get(t.Context(), taskID, transport.WithResponseInfo(&info))
	if err != nil || task.State != async.StateDone || info.StatusCode != http.StatusOK {
		t.Fatalf("task = %+v, err = %v, info = %+v", task, err, info)
	}
	list, err := s.List(t.Context(), query.ListOptions{Limit: 10, Filter: "state=DONE"})
	if err != nil || len(list.Rows) != 2 {
		t.Fatalf("list = %+v, err = %v", list, err)
	}
	if _, err := s.Get(t.Context(), "11111111-1111-1111-1111-111111111111"); !sdkerrors.HasCode(err, 1021) {
		t.Fatalf("err = %v", err)
	}
}

// Wait: первый Get сразу, затем каждые interval до завершения; срок задаёт ctx; ошибка Get
// возвращается.
func TestWait(t *testing.T) {
	var calls atomic.Int32
	states := []string{"PENDING", "PROCESSING", "API_ERROR"}
	s := newService(t, func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if r.URL.Path == "/async/"+taskID {
			w.Write([]byte(`{"id":"` + taskID + `","state":"` + states[min(int(n), len(states))-1] + `"}`))
			return
		}
		w.Write([]byte(`{"id":"` + taskID + `","state":"PROCESSING"}`))
	}, transport.Config{})
	task, err := s.Wait(t.Context(), taskID, time.Millisecond)
	if err != nil || task.State != async.StateAPIError || calls.Load() != 3 {
		t.Fatalf("task = %+v, err = %v, calls = %d", task, err, calls.Load())
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
	defer cancel()
	if task, err := s.Wait(ctx, "11111111-1111-1111-1111-111111111111", time.Millisecond); task != nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline: task = %+v, err = %v", task, err)
	}

	failing := newService(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusForbidden) }, transport.Config{})
	if _, err := failing.Wait(t.Context(), taskID, time.Millisecond); sdkerrors.StatusCode(err) != http.StatusForbidden {
		t.Fatalf("err = %v", err)
	}
}

// Result: GET async/{id}/result, redirect на ссылку другого origin без Authorization; ошибка
// API_ERROR — APIError, w не тронут.
func TestResult(t *testing.T) {
	file := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Errorf("authorization on result link")
		}
		w.Write([]byte(`{"rows":[]}`))
	}))
	defer file.Close()
	auth, _ := transport.BearerToken("secret")
	s := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/async/" + taskID + "/result":
			http.Redirect(w, r, file.URL+"/result?sig=1", http.StatusFound)
		case "/async/22222222-2222-2222-2222-222222222222/result":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"errors":[{"error":"результат еще не готов","code":61006}]}`))
		default:
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"errors":[{"error":"Доступ запрещен","code":1016}]}`))
		}
	}, transport.Config{Auth: auth})
	var out bytes.Buffer
	if n, err := s.Result(t.Context(), taskID, &out); err != nil || n != int64(out.Len()) || out.String() != `{"rows":[]}` {
		t.Fatalf("n = %d, err = %v, out = %q", n, err, out.String())
	}
	out.Reset()
	if _, err := s.Result(t.Context(), "11111111-1111-1111-1111-111111111111", &out); !sdkerrors.HasCode(err, 1016) || out.Len() != 0 {
		t.Fatalf("err = %v, out = %q", err, out.String())
	}
	if _, err := s.Result(t.Context(), "22222222-2222-2222-2222-222222222222", &out); !sdkerrors.HasCode(err, 61006) || out.Len() != 0 {
		t.Fatalf("not ready: err = %v, out = %q", err, out.String())
	}
}

// Cancel: PUT без тела, 204 — успех, другой статус — ошибка; 503 при включённых повторах не
// повторяется.
func TestCancel(t *testing.T) {
	for status, ok := range map[int]bool{http.StatusNoContent: true, http.StatusOK: false, http.StatusAccepted: false, http.StatusServiceUnavailable: false} {
		var calls atomic.Int32
		s := newService(t, func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			if r.Method != http.MethodPut || r.URL.Path != "/async/"+taskID+"/cancel" || r.ContentLength > 0 {
				t.Errorf("request = %s %s %d", r.Method, r.URL, r.ContentLength)
			}
			w.WriteHeader(status)
		}, transport.Config{Retry: &transport.RetryPolicy{MaxAttempts: 3}})
		if err := s.Cancel(t.Context(), taskID); (err == nil) != ok || calls.Load() != 1 {
			t.Fatalf("status %d: err = %v, calls = %d", status, err, calls.Load())
		}
	}
}

// Неверный id и interval — ошибка до запроса; ResponseInfo сброшен.
func TestRejectsBeforeRequest(t *testing.T) {
	var calls atomic.Int32
	s := newService(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }, transport.Config{})
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get": func(o ...transport.RequestOption) error { _, err := s.Get(t.Context(), "bad", o...); return err },
		"result": func(o ...transport.RequestOption) error {
			_, err := s.Result(t.Context(), "../x", &bytes.Buffer{}, o...)
			return err
		},
		"cancel": func(o ...transport.RequestOption) error { return s.Cancel(t.Context(), taskID+"/x", o...) },
		"wait id": func(o ...transport.RequestOption) error {
			_, err := s.Wait(t.Context(), "bad", time.Millisecond, o...)
			return err
		},
		"wait interval": func(o ...transport.RequestOption) error { _, err := s.Wait(t.Context(), taskID, 0, o...); return err },
		"list limit": func(o ...transport.RequestOption) error {
			_, err := s.List(t.Context(), query.ListOptions{Limit: 1001}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Fatalf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestNewServiceNilClient(t *testing.T) {
	if _, err := async.NewService(nil); err == nil || !strings.Contains(err.Error(), "client") {
		t.Fatalf("err = %v", err)
	}
}
