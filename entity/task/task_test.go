package task_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/task"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	taskID     = "cfe5a6ae-8b87-11e8-d9ce-84d900000000"
	task2ID    = "9971ba00-8b88-11e8-d9ce-84d900000009"
	employeeID = "98fa7086-8aa1-11e8-7210-075e0000002c"
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

func ref(path, typ string, metadata bool) model.Reference {
	m := model.Meta{Href: api + path, Type: model.MetaType(typ), MediaType: "application/json"}
	if metadata {
		m.MetadataHref = api + "entity/" + typ + "/metadata"
	}
	return model.Reference{Meta: m}
}

func service(c *transport.Client) *task.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := task.NewService(c)
	return s
}

func at(v string) time.Time { return time.Time(documenttest.DateTime(v)) }

var (
	full = task.CreateRequest{
		Description: "Исправить реквизиты юрлица",
		Assignee:    ref("entity/employee/"+employeeID, "employee", true),
		DueToDate:   model.Value(documenttest.DateTime("2017-05-12 17:12:00")),
		Done:        model.Value(false),
		Agent:       model.Value(ref("entity/organization/9927bf0d-8aa1-11e8-7210-075e00000054", "organization", true)),
		State:       model.Value(ref("entity/task/metadata/states/4f70c518-60a1-11e7-6adb-ede500000003", "state", false)),
	}
	update = task.UpdateRequest{
		Description: model.Value("Уточнить контактные лица"),
		Done:        model.Value(true),
		DueToDate:   model.Null[model.DateTime](),
		Agent:       model.Null[model.Reference](),
		Operation:   model.Value(ref("entity/supply/7eb5b552-8aa3-11e8-7210-075e000000f1", "supply", true)),
	}
)

// Фикстура официального SDK: все поля модели, включая раскрытые комментарии и тип задачи.
func TestTaskFixture(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got task.Task
			if err := codec.Unmarshal(read(t, "upstream/task.json"), &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != taskID || got.Meta.Type != "task" || got.AccountID != "98182760-8aa1-11e8-7210-075e00000001" ||
				got.Description != "Уточнить контактные лица" || !got.Done {
				t.Fatalf("task = %+v", got)
			}
			for name, pair := range map[string][2]time.Time{
				"created":   {time.Time(got.Created), at("2018-07-19 22:13:25")},
				"updated":   {time.Time(got.Updated), at("2018-07-19 22:24:50")},
				"dueToDate": {time.Time(got.DueToDate), at("2018-07-20 22:09:00")},
				"completed": {time.Time(got.Completed), at("2018-07-19 22:24:24")},
			} {
				if !pair[0].Equal(pair[1]) {
					t.Errorf("%s = %v, want %v", name, pair[0], pair[1])
				}
			}
			for name, pair := range map[string]struct {
				ref *model.Reference
				typ model.MetaType
			}{
				"assignee": {got.Assignee, "employee"}, "author": {got.Author, "employee"}, "implementer": {got.Implementer, "employee"},
				"authorApplication": {got.AuthorApplication, "application"}, "agent": {got.Agent, "counterparty"},
				"operation": {got.Operation, "supply"}, "state": {got.State, "state"},
			} {
				if pair.ref == nil || pair.ref.Meta.Type != pair.typ || pair.ref.Meta.Href == "" {
					t.Errorf("%s = %+v", name, pair.ref)
				}
			}
			if got.Files == nil || got.Files.Meta.Type != "files" || got.Files.Meta.Size != 0 {
				t.Errorf("files = %+v", got.Files)
			}
			if got.Notes == nil || got.Notes.Meta.Type != "tasknote" || got.Notes.Meta.Size != 2 || len(got.Notes.Rows) != 2 {
				t.Fatalf("notes = %+v", got.Notes)
			}
			c := got.Notes.Rows[0]
			if c.ID != "5eb5f69d-8b88-11e8-d9ce-84d900000001" || c.Text != "Связался с клиентом, ожидаем подтверждение." ||
				!time.Time(c.Moment).Equal(at("2018-07-19 22:20:10")) || c.Author == nil || c.Files == nil || c.Files.Meta.Size != 1 {
				t.Fatalf("comment = %+v", c)
			}
		})
	}
}

// Примеры источника: список (привязка к контрагенту или документу) и массовый ответ.
func TestSourceExamples(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var page model.Collection[task.Task]
			if err := codec.Unmarshal(read(t, "custom/task/list.json"), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Rows) != 2 || page.Meta.Size != 2 {
				t.Fatalf("page = %+v", page)
			}
			a, b := page.Rows[0], page.Rows[1]
			if a.Agent == nil || a.Agent.Meta.Type != "counterparty" || a.Operation != nil || a.Done || !a.Completed.IsZero() || a.Implementer != nil ||
				a.Notes == nil || a.Notes.Meta.Size != 2 || a.Notes.Rows != nil {
				t.Fatalf("row 0 = %+v", a)
			}
			if b.Operation == nil || b.Operation.Meta.Type != "supply" || b.Agent != nil || b.Description != "task1" {
				t.Fatalf("row 1 = %+v", b)
			}
			var many []task.Task
			if err := codec.Unmarshal(read(t, "custom/task/batch.json"), &many); err != nil || len(many) != 2 {
				t.Fatalf("batch = %+v, err = %v", many, err)
			}
			if many[0].ID != task2ID || many[0].Done || many[1].ID != taskID || !many[1].Done || many[1].Implementer == nil ||
				!time.Time(many[1].Completed).Equal(at("2018-07-19 22:19:04")) {
				t.Fatalf("batch = %+v", many)
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
			{"create full", full, read(t, "custom/task/create.golden.json")},
			{"update", update, read(t, "custom/task/update.golden.json")},
			{"create minimal", task.CreateRequest{Description: "d", Assignee: ref("entity/employee/"+employeeID, "employee", false)},
				[]byte(`{"description":"d","assignee":{"meta":{"href":"` + api + `entity/employee/` + employeeID + `","type":"employee","mediaType":"application/json"}}}`)},
			{"update done false", task.UpdateRequest{Done: model.Value(false)}, []byte(`{"done":false}`)},
			{"comment", task.CommentRequest{Text: "Привет, {{employee;861d34a9-f1b3-11ee-ac12-00110000004e}}! Задача актуальна?"}, read(t, "custom/task/comment.golden.json")},
			{"state create", task.StateCreateRequest{Name: "Встреча", Color: documenttest.Pointer[int64](69446)}, read(t, "custom/task/state-create.golden.json")},
			{"state create black", task.StateCreateRequest{Name: "n", Color: documenttest.Pointer[int64](0), StateType: model.Value("Successful")},
				[]byte(`{"name":"n","color":0,"stateType":"Successful"}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[task.CreateRequest](), reflect.TypeFor[task.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Author", "AuthorApplication", "Completed", "Created", "Updated", "Implementer", "Files", "Notes"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestServiceHTTP(t *testing.T) {
	response := string(read(t, "upstream/task.json"))
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*task.Task, error)
	}{
		{"get", "GET", "entity/task/" + taskID, nil, func(ctx context.Context, c *transport.Client) (*task.Task, error) {
			return service(c).Get(ctx, taskID)
		}},
		{"create", "POST", "entity/task", read(t, "custom/task/create.golden.json"), func(ctx context.Context, c *transport.Client) (*task.Task, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/task/" + taskID, read(t, "custom/task/update.golden.json"), func(ctx context.Context, c *transport.Client) (*task.Task, error) {
			return service(c).Update(ctx, taskID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *task.Task) {
						if got.ID != taskID || got.Assignee == nil || got.Operation == nil || len(got.Notes.Rows) != 2 {
							t.Fatalf("task = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "task", read(t, "upstream/task.json"), "", "", func(c *transport.Client) listtest.ListFunc[task.Task] {
		return service(c).List
	}, func(t *testing.T, got task.Task) {
		if got.ID != taskID || got.Description != "Уточнить контактные лица" {
			t.Fatalf("row = %+v", got)
		}
	})
}

type request struct{ method, uri, body string }

// recorder отвечает по маршруту и записывает запросы.
func recorder(t *testing.T, respond func(r *http.Request) (int, []byte)) (*transport.Client, *[]request) {
	t.Helper()
	got := &[]request{}
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*got = append(*got, request{r.Method, r.URL.RequestURI(), string(body)})
		status, data := respond(r)
		w.WriteHeader(status)
		w.Write(data)
	}))
	t.Cleanup(server.Close)
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	return c, got
}

func checkRequests(t *testing.T, got, want []request) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("requests = %+v", got)
	}
	for i := range want {
		if got[i].method != want[i].method || got[i].uri != want[i].uri {
			t.Errorf("request %d = %s %s, want %s %s", i, got[i].method, got[i].uri, want[i].method, want[i].uri)
		}
		if want[i].body == "" {
			if got[i].body != "" {
				t.Errorf("request %d body = %s", i, got[i].body)
			}
		} else {
			documenttest.JSONEqual(t, []byte(got[i].body), []byte(want[i].body))
		}
	}
}

// Delete и массовые операции задач: тела с meta (metadataHref, как в примерах) и ответы.
func TestServiceWrite(t *testing.T) {
	c, got := recorder(t, func(r *http.Request) (int, []byte) {
		switch {
		case r.Method == http.MethodDelete:
			return 200, nil
		case strings.HasSuffix(r.URL.Path, "/delete"):
			return 200, []byte(`[{"info":"Сущность 'task' с UUID: ` + taskID + ` успешно удалена"},{"info":"ok"}]`)
		}
		return 200, read(t, "custom/task/batch.json")
	})
	s, ctx := service(c), t.Context()
	if err := s.Delete(ctx, taskID); err != nil {
		t.Fatal(err)
	}
	deleted, err := s.DeleteMany(ctx, []string{taskID, task2ID})
	if err != nil || len(deleted) != 2 || !strings.Contains(deleted[0].Info, taskID) {
		t.Fatalf("delete many = %+v, err = %v", deleted, err)
	}
	created, err := s.CreateMany(ctx, []task.CreateRequest{full, {Description: "d", Assignee: full.Assignee}})
	if err != nil || len(created) != 2 || created[0].Value == nil || created[0].Value.ID != task2ID {
		t.Fatalf("create many = %+v, err = %v", created, err)
	}
	updated, err := s.UpdateMany(ctx, []batch.Update[task.UpdateRequest]{{ID: task2ID, Request: task.UpdateRequest{Done: model.Value(false)}}, {ID: taskID, Request: update}})
	if err != nil || len(updated) != 2 || updated[1].Value == nil || !updated[1].Value.Done {
		t.Fatalf("update many = %+v, err = %v", updated, err)
	}
	meta := func(id string) string {
		return `{"href":"` + c.Href("entity/task/"+id) + `","metadataHref":"` + c.Href("entity/task/metadata") + `","type":"task","mediaType":"application/json"}`
	}
	checkRequests(t, *got, []request{
		{"DELETE", "/entity/task/" + taskID, ""},
		{"POST", "/entity/task/delete", `[{"meta":` + meta(taskID) + `},{"meta":` + meta(task2ID) + `}]`},
		{"POST", "/entity/task", "[" + string(read(t, "custom/task/create.golden.json")) + `,{"description":"d","assignee":` + mustJSON(t, full.Assignee) + `}]`},
		{"POST", "/entity/task", `[{"meta":` + meta(task2ID) + `,"done":false},{"meta":` + meta(taskID) + `,` + strings.TrimPrefix(string(read(t, "custom/task/update.golden.json")), "{") + `]`},
	})
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	data, err := documenttest.Codecs[1].Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// Ошибки входа — до HTTP; ResponseInfo сброшен, префикс moysklad: один.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(c), t.Context()
	comment := task.CommentRequest{Text: "t"}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, taskID+"/x", o...); return err },
		"create no description": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, task.CreateRequest{Assignee: full.Assignee}, o...)
			return err
		},
		"create no assignee": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, task.CreateRequest{Description: "d"}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, taskID, task.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"list fields": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, query.ListOptions{Fields: "stock"}, o...)
			return err
		},
		"create many empty": func(o ...transport.RequestOption) error { _, err := s.CreateMany(ctx, nil, o...); return err },
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[task.UpdateRequest]{{ID: "bad", Request: update}}, o...)
			return err
		},
		"delete many bad": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"comments bad task": func(o ...transport.RequestOption) error {
			_, err := s.ListComments(ctx, "bad", task.CommentsOptions{}, o...)
			return err
		},
		"comments limit": func(o ...transport.RequestOption) error {
			_, err := s.ListComments(ctx, taskID, task.CommentsOptions{Limit: 101}, o...)
			return err
		},
		"comments offset": func(o ...transport.RequestOption) error {
			_, err := s.ListComments(ctx, taskID, task.CommentsOptions{Offset: -1}, o...)
			return err
		},
		"get comment bad": func(o ...transport.RequestOption) error { _, err := s.GetComment(ctx, taskID, "bad", o...); return err },
		"create comment empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateComment(ctx, taskID, task.CommentRequest{}, o...)
			return err
		},
		"create comment bad task": func(o ...transport.RequestOption) error {
			_, err := s.CreateComment(ctx, "../x", comment, o...)
			return err
		},
		"create many comments empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyComments(ctx, taskID, nil, o...)
			return err
		},
		"update comment empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateComment(ctx, taskID, taskID, task.CommentRequest{}, o...)
			return err
		},
		"delete comment bad": func(o ...transport.RequestOption) error { return s.DeleteComment(ctx, "bad", taskID, o...) },
		"delete many comments bad task": func(o ...transport.RequestOption) error {
			return s.DeleteManyComments(ctx, "bad", []string{taskID}, o...)
		},
		"delete many comments empty": func(o ...transport.RequestOption) error { return s.DeleteManyComments(ctx, taskID, nil, o...) },
		"create state no color": func(o ...transport.RequestOption) error {
			_, err := s.CreateState(ctx, task.StateCreateRequest{Name: "n"}, o...)
			return err
		},
		"create state no name": func(o ...transport.RequestOption) error {
			_, err := s.CreateState(ctx, task.StateCreateRequest{Color: documenttest.Pointer[int64](1)}, o...)
			return err
		},
		"update state empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateState(ctx, taskID, metadata.StateUpdateRequest{}, o...)
			return err
		},
		"update state bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateState(ctx, "bad", metadata.StateUpdateRequest{Color: model.Value[int64](1)}, o...)
			return err
		},
		"delete state bad id": func(o ...transport.RequestOption) error { return s.DeleteState(ctx, "bad", o...) },
		"create many states empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyStates(ctx, nil, o...)
			return err
		},
		"update many states bad": func(o ...transport.RequestOption) error {
			_, err := s.UpdateManyStates(ctx, []batch.Update[metadata.StateUpdateRequest]{{ID: "bad"}}, o...)
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
	if _, err := task.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}
