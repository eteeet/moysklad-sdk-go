package task_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/task"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	listTaskID = "47bc1d9b-8b87-11e8-d9ce-84d900000017"
	commentID  = "2050793f-b8c6-11f0-0a80-204500000006"
	stateID    = "6262b270-60c3-11e7-6adb-ede50000000d"
	mention    = "Привет, {{employee;861d34a9-f1b3-11ee-ac12-00110000004e}}! Задача актуальна?"
)

func TestCommentExamples(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var page model.Collection[task.Comment]
			if err := codec.Unmarshal(read(t, "custom/task/comments.json"), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Rows) != 2 || page.Meta.Type != "tasknote" || page.Rows[1].Text != "текст комментария 2" ||
				!time.Time(page.Rows[0].Moment).Equal(at("2018-07-19 22:10:01")) || page.Rows[0].Author == nil {
				t.Fatalf("comments = %+v", page)
			}
			var created []task.Comment
			if err := codec.Unmarshal(read(t, "custom/task/comment-create.json"), &created); err != nil || len(created) != 1 {
				t.Fatalf("created = %+v, err = %v", created, err)
			}
			c := created[0]
			if c.ID != commentID || c.AccountID != "5ae35472-b8c5-11f0-0a82-042d00000002" || c.Text != mention ||
				c.Files == nil || c.Files.Meta.Type != "files" || time.Time(c.Moment).Nanosecond() != 338000000 {
				t.Fatalf("created = %+v", c)
			}
			// Пример Get источника — без id (см. GetComment).
			var one task.Comment
			if err := codec.Unmarshal(read(t, "custom/task/comment.json"), &one); err != nil || one.ID != "" || one.Text != "текст комментария 1" {
				t.Fatalf("comment = %+v, err = %v", one, err)
			}
			var meta metadata.Entity
			if err := codec.Unmarshal(read(t, "custom/task/metadata.json"), &meta); err != nil || len(meta.States) != 2 {
				t.Fatalf("metadata = %+v, err = %v", meta, err)
			}
			if s := meta.States[1]; s.Name != "Связаться" || s.Color != 10667543 || s.StateType != "Regular" || s.EntityType != "task" {
				t.Fatalf("state = %+v", s)
			}
		})
	}
}

// Get/Create/Update комментариев и запись типов задач: общий HTTP-контракт.
func TestCommentHTTP(t *testing.T) {
	created := read(t, "custom/task/comment-create.json")
	item := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(string(created)), "["), "]")
	commentBody := read(t, "custom/task/comment.golden.json")
	input := task.CommentRequest{Text: mention}
	base := "entity/task/" + listTaskID + "/notes"
	comment := func(f func(context.Context, *task.Service) (*task.Comment, error)) func(context.Context, *transport.Client) (*task.Comment, error) {
		return func(ctx context.Context, c *transport.Client) (*task.Comment, error) { return f(ctx, service(c)) }
	}
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*task.Comment, error)
	}{
		{"get", "GET", base + "/" + commentID, nil, comment(func(ctx context.Context, s *task.Service) (*task.Comment, error) {
			return s.GetComment(ctx, listTaskID, commentID)
		})},
		{"create", "POST", base, commentBody, comment(func(ctx context.Context, s *task.Service) (*task.Comment, error) {
			return s.CreateComment(ctx, listTaskID, input)
		})},
		{"update", "PUT", base + "/" + commentID, commentBody, comment(func(ctx context.Context, s *task.Service) (*task.Comment, error) {
			return s.UpdateComment(ctx, listTaskID, commentID, input)
		})},
	} {
		t.Run(op.name, func(t *testing.T) {
			cases := documenttest.ResponseCases(item)
			if op.method != "GET" {
				cases = append(cases, documenttest.ResponseCase{Name: "source array", Status: 200, Body: string(created), Success: true},
					documenttest.ResponseCase{Name: "empty array", Status: 200, Body: "[]"})
			}
			for _, tc := range cases {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *task.Comment) {
						if got.ID != commentID || got.Text != mention {
							t.Fatalf("comment = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}

	stateResponse := string(read(t, "custom/task/state.json"))
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*metadata.State, error)
	}{
		{"create state", "POST", "entity/task/metadata/states", read(t, "custom/task/state-create.golden.json"), func(ctx context.Context, c *transport.Client) (*metadata.State, error) {
			return service(c).CreateState(ctx, task.StateCreateRequest{Name: "Встреча", Color: documenttest.Pointer[int64](69446)})
		}},
		{"update state", "PUT", "entity/task/metadata/states/" + stateID, []byte(`{"color":255}`), func(ctx context.Context, c *transport.Client) (*metadata.State, error) {
			return service(c).UpdateState(ctx, stateID, metadata.StateUpdateRequest{Color: model.Value[int64](255)})
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(stateResponse) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *metadata.State) {
						if got.ID != stateID || got.Name != "Встреча" || got.Color != 69446 || got.EntityType != "task" {
							t.Fatalf("state = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Списки, удаления и массовые операции комментариев и типов задач.
func TestCommentCollections(t *testing.T) {
	c, got := recorder(t, func(r *http.Request) (int, []byte) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/metadata"):
			return 200, read(t, "custom/task/metadata.json")
		case r.Method == http.MethodGet:
			return 200, read(t, "custom/task/comments.json")
		case r.Method == http.MethodDelete, strings.HasSuffix(r.URL.Path, "/delete"):
			return 200, nil
		case strings.HasSuffix(r.URL.Path, "/states"):
			return 200, read(t, "custom/task/states-batch.json")
		}
		return 500, nil
	})
	s, ctx := service(c), t.Context()
	filter, err := query.Filter(query.GreaterOrEq("updated", "2016-04-15 15:48:46"), query.Eq("updatedBy", "admin@admin"))
	if err != nil {
		t.Fatal(err)
	}
	page, err := s.ListComments(ctx, listTaskID, task.CommentsOptions{Limit: 100, Offset: 40, Filter: filter})
	if err != nil || len(page.Rows) != 2 || page.Rows[0].Text != "текст комментария 1" {
		t.Fatalf("comments = %+v, err = %v", page, err)
	}
	if _, err := s.ListComments(ctx, listTaskID, task.CommentsOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteComment(ctx, listTaskID, commentID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteManyComments(ctx, listTaskID, []string{commentID}); err != nil {
		t.Fatal(err)
	}
	states, err := s.ListStates(ctx)
	if err != nil || len(states) != 2 || states[0].ID != "4f70c518-60a1-11e7-6adb-ede500000003" || states[0].Color != 15106326 {
		t.Fatalf("states = %+v, err = %v", states, err)
	}
	if err := s.DeleteState(ctx, stateID); err != nil {
		t.Fatal(err)
	}
	newStates, err := s.CreateManyStates(ctx, []task.StateCreateRequest{{Name: "Встреча", Color: documenttest.Pointer[int64](8767198)}, {Name: "Связаться", Color: documenttest.Pointer[int64](34617)}})
	if err != nil || len(newStates) != 2 || newStates[1].Value == nil || newStates[1].Value.Color != 34617 {
		t.Fatalf("create many states = %+v, err = %v", newStates, err)
	}
	changed, err := s.UpdateManyStates(ctx, []batch.Update[metadata.StateUpdateRequest]{
		{ID: "b55d2ddf-60c3-11e7-6adb-ede500000010", Request: metadata.StateUpdateRequest{Name: model.Value("Встреча")}},
		{ID: "b56215dc-60c3-11e7-6adb-ede500000013", Request: metadata.StateUpdateRequest{Color: model.Value[int64](34617)}}})
	if err != nil || len(changed) != 2 || changed[0].Value == nil || changed[0].Value.Name != "Встреча" {
		t.Fatalf("update many states = %+v, err = %v", changed, err)
	}

	base := "/entity/task/" + listTaskID + "/notes"
	stateMeta := func(id string) string {
		return `{"href":"` + c.Href("entity/task/metadata/states/"+id) + `","type":"state","mediaType":"application/json"}`
	}
	want := []request{
		{"GET", base + "?filter=updated%3E%3D2016-04-15+15%3A48%3A46%3BupdatedBy%3Dadmin%40admin&limit=100&offset=40", ""},
		{"GET", base, ""},
		{"DELETE", base + "/" + commentID, ""},
		{"POST", base + "/delete", `[{"meta":{"href":"` + c.Href(base[1:]+"/"+commentID) + `","type":"tasknote","mediaType":"application/json"}}]`},
		{"GET", "/entity/task/metadata", ""},
		{"DELETE", "/entity/task/metadata/states/" + stateID, ""},
		{"POST", "/entity/task/metadata/states", `[{"name":"Встреча","color":8767198},{"name":"Связаться","color":34617}]`},
		{"POST", "/entity/task/metadata/states", `[{"meta":` + stateMeta("b55d2ddf-60c3-11e7-6adb-ede500000010") + `,"name":"Встреча"},{"meta":` + stateMeta("b56215dc-60c3-11e7-6adb-ede500000013") + `,"color":34617}]`},
	}
	checkRequests(t, *got, want)

	// Массовое создание комментариев: тело — пример источника, ответ — массив комментариев.
	c2, got2 := recorder(t, func(r *http.Request) (int, []byte) {
		return 200, []byte(`[{"id":"a","text":"Привет"},{"id":"b","text":"Могу сделать сегодня!"}]`)
	})
	res, err := service(c2).CreateManyComments(ctx, listTaskID, []task.CommentRequest{{Text: mention}, {Text: "Могу сделать сегодня!"}})
	if err != nil || len(res) != 2 || res[1].Value == nil || res[1].Value.Text != "Могу сделать сегодня!" {
		t.Fatalf("create many comments = %+v, err = %v", res, err)
	}
	checkRequests(t, *got2, []request{{"POST", base, `[{"text":"` + mention + `"},{"text":"Могу сделать сегодня!"}]`}})
}
