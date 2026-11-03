package eventfeed_test

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
	"github.com/eteeet/moysklad-sdk-go/eventfeed"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	docID   = "e4609c69-00bc-11ef-ac12-00120000001a"
	noteID  = "844a0ef9-19ac-11ef-ac12-000b00000000"
	note2ID = "c49d7517-5b69-40dc-9254-6721f72a860a"
	mention = "{{employee;861d34a9-f1b3-11ee-ac12-00110000004e}}, привет"
	created = "2024-05-24 12:03:40.849"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/custom/eventfeed/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *eventfeed.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := eventfeed.NewService(c)
	return s
}

func checkNote(t *testing.T, n eventfeed.Note, id, description string) {
	t.Helper()
	if n.ID != id || n.Meta.Type != "eventnote" || n.AccountID != "85284e22-f1b3-11ee-ac12-000f00000001" ||
		n.Description != description || n.Author == nil || n.Author.Meta.Type != "employee" || n.AuthorApplication != nil {
		t.Fatalf("note = %+v", n)
	}
}

func TestFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var page struct {
				Meta struct {
					Type string `json:"type"`
					Size int    `json:"size"`
				} `json:"meta"`
				Rows []eventfeed.Note `json:"rows"`
			}
			if err := codec.Unmarshal(read(t, "list.json"), &page); err != nil {
				t.Fatal(err)
			}
			if page.Meta.Type != "eventnote" || page.Meta.Size != 2 || len(page.Rows) != 2 {
				t.Fatalf("page = %+v", page)
			}
			checkNote(t, page.Rows[1], noteID, mention)
			// Упоминание — строка без разбора; момент с миллисекундами.
			want := time.Time(documenttest.DateTime(created))
			if got := time.Time(page.Rows[0].Created); !got.Equal(want) || got.Nanosecond() != 849000000 {
				t.Fatalf("created = %v, want %v", got, want)
			}
			var one eventfeed.Note
			if err := codec.Unmarshal(read(t, "get.json"), &one); err != nil {
				t.Fatal(err)
			}
			checkNote(t, one, noteID, mention)
			var many []eventfeed.Note
			if err := codec.Unmarshal(read(t, "batch.json"), &many); err != nil || len(many) != 2 {
				t.Fatalf("batch = %+v, err = %v", many, err)
			}
			checkNote(t, many[0], note2ID, "Привет, {{employee;730c1b3d-00ba-11ef-ac12-00120000000d}}! Это новое Событие")
			var app eventfeed.Note
			if err := codec.Unmarshal([]byte(`{"id":"x","authorApplication":{"meta":{"type":"application"}}}`), &app); err != nil ||
				app.AuthorApplication == nil || app.AuthorApplication.Meta.Type != "application" {
				t.Fatalf("application = %+v, err = %v", app, err)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	input := eventfeed.NoteRequest{Description: "Привет, {{employee;730c1b3d-00ba-11ef-ac12-00120000000d}}! Как у тебя дела?"}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			got, err := codec.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, got, read(t, "create.golden.json"))
		})
	}
	typ := reflect.TypeFor[eventfeed.NoteRequest]()
	if typ.NumField() != 1 {
		t.Fatalf("NoteRequest has %d fields, want only description", typ.NumField())
	}
}

// Get/Create/Update: путь с типом документа, тело, ошибки, одна отправка; Create и Update
// принимают и массив из одного события, как в примере Create источника.
func TestHTTP(t *testing.T) {
	get, created := read(t, "get.json"), read(t, "create.json")
	body := read(t, "create.golden.json")
	input := eventfeed.NoteRequest{Description: "Привет, {{employee;730c1b3d-00ba-11ef-ac12-00120000000d}}! Как у тебя дела?"}
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*eventfeed.Note, error)
	}{
		{"get", "GET", "entity/customerorder/" + docID + "/notes/" + noteID, nil, func(ctx context.Context, c *transport.Client) (*eventfeed.Note, error) {
			return service(c).Get(ctx, "customerorder", docID, noteID)
		}},
		{"create", "POST", "entity/invoicein/" + docID + "/notes", body, func(ctx context.Context, c *transport.Client) (*eventfeed.Note, error) {
			return service(c).Create(ctx, "invoicein", docID, input)
		}},
		{"update", "PUT", "entity/purchasereturn/" + docID + "/notes/" + noteID, body, func(ctx context.Context, c *transport.Client) (*eventfeed.Note, error) {
			return service(c).Update(ctx, "purchasereturn", docID, noteID, input)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			cases := documenttest.ResponseCases(string(get))
			if op.method == "POST" {
				// Источник отвечает на создание 201 (Collection.Created): это успех, а не неожиданный статус.
				for i := range cases {
					if cases[i].Status == 201 {
						cases[i].Name, cases[i].Success = "source 201", true
					}
				}
				cases = append(cases, documenttest.ResponseCase{Name: "source 201 array", Status: 201, Body: string(created), Success: true})
			}
			if op.method != "GET" {
				cases = append(cases, documenttest.ResponseCase{Name: "source array", Status: 200, Body: string(created), Success: true},
					documenttest.ResponseCase{Name: "array of two", Status: 200, Body: "[" + string(get) + "," + string(get) + "]"},
					documenttest.ResponseCase{Name: "empty array", Status: 200, Body: "[]"})
			}
			for _, tc := range cases {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *eventfeed.Note) {
						if got.ID != noteID || got.Author == nil || !strings.Contains(got.Description, "{{employee;") {
							t.Fatalf("note = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

type request struct{ method, uri, body string }

// List, Delete и массовые операции: метод, путь, query, тело и разбор ответа.
func TestCollections(t *testing.T) {
	var got []request
	var deleteStatus atomic.Int32
	deleteStatus.Store(http.StatusNoContent)
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, request{r.Method, r.URL.RequestURI(), string(body)})
		switch {
		case r.Method == http.MethodGet:
			w.Write(read(t, "list.json"))
		case r.Method == http.MethodDelete:
			w.WriteHeader(int(deleteStatus.Load()))
		case strings.HasSuffix(r.URL.Path, "/delete"):
			// Источник: 200 с пустым телом.
		default:
			// Массовое создание и обновление: источник отвечает 201.
			w.WriteHeader(http.StatusCreated)
			w.Write(read(t, "batch.json"))
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(c), t.Context()
	page, err := s.List(ctx, "supply", docID, query.PageOptions{Limit: 10, Offset: 40})
	if err != nil || len(page.Rows) != 2 || page.Meta.Size != 2 || page.Rows[0].Description != mention {
		t.Fatalf("list = %+v, err = %v", page, err)
	}
	if _, err := s.List(ctx, "customerorder", docID, query.PageOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(ctx, "customerorder", docID, noteID); err != nil {
		t.Fatal(err)
	}
	deleteStatus.Store(http.StatusOK)
	if err := s.Delete(ctx, "customerorder", docID, noteID); err != nil {
		t.Fatal(err)
	}
	deleteStatus.Store(http.StatusAccepted)
	if err := s.Delete(ctx, "customerorder", docID, noteID); err == nil || strings.Count(err.Error(), "moysklad:") != 1 {
		t.Fatalf("delete 202 err = %v", err)
	}
	if err := s.DeleteMany(ctx, "customerorder", docID, []string{noteID, note2ID}); err != nil {
		t.Fatal(err)
	}
	created, err := s.CreateMany(ctx, "customerorder", docID, []eventfeed.NoteRequest{{Description: "a"}, {Description: "b"}})
	if err != nil || len(created) != 2 || created[0].Value == nil || created[0].Value.ID != note2ID {
		t.Fatalf("create many = %+v, err = %v", created, err)
	}
	updated, err := s.UpdateMany(ctx, "customerorder", docID, []batch.Update[eventfeed.NoteRequest]{
		{ID: noteID, Request: eventfeed.NoteRequest{Description: "u"}}, {ID: note2ID, Request: eventfeed.NoteRequest{Description: "v"}}})
	if err != nil || len(updated) != 2 || updated[1].Value == nil || updated[1].Value.ID != noteID {
		t.Fatalf("update many = %+v, err = %v", updated, err)
	}

	base := "/entity/customerorder/" + docID + "/notes"
	meta := func(id string) string {
		return `{"href":"` + server.URL + base + "/" + id + `","type":"eventnote","mediaType":"application/json"}`
	}
	want := []request{
		{"GET", "/entity/supply/" + docID + "/notes?limit=10&offset=40", ""},
		{"GET", base + "?limit=1000&offset=0", ""},
		{"DELETE", base + "/" + noteID, ""},
		{"DELETE", base + "/" + noteID, ""},
		{"DELETE", base + "/" + noteID, ""},
		{"POST", base + "/delete", `[{"meta":` + meta(noteID) + `},{"meta":` + meta(note2ID) + `}]`},
		{"POST", base, `[{"description":"a"},{"description":"b"}]`},
		{"POST", base, `[{"meta":` + meta(noteID) + `,"description":"u"},{"meta":` + meta(note2ID) + `,"description":"v"}]`},
	}
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

// Ошибки входа — до HTTP; ResponseInfo сброшен, префикс moysklad: один.
func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(c), t.Context()
	note := eventfeed.NoteRequest{Description: "d"}
	cases := map[string]func(...transport.RequestOption) error{
		"get bad note": func(o ...transport.RequestOption) error {
			_, err := s.Get(ctx, "customerorder", docID, "bad", o...)
			return err
		},
		"get bad document id": func(o ...transport.RequestOption) error {
			_, err := s.Get(ctx, "customerorder", docID+"/x", noteID, o...)
			return err
		},
		"list fields": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, "customerorder", docID, query.PageOptions{Fields: "x"}, o...)
			return err
		},
		"list limit": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, "customerorder", docID, query.PageOptions{Limit: 1001}, o...)
			return err
		},
		"create empty": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, "customerorder", docID, eventfeed.NoteRequest{}, o...)
			return err
		},
		"create bad id": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, "customerorder", "bad", note, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, "customerorder", docID, noteID, eventfeed.NoteRequest{}, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error { return s.Delete(ctx, "customerorder", docID, "bad", o...) },
		"delete many empty": func(o ...transport.RequestOption) error {
			return s.DeleteMany(ctx, "customerorder", docID, nil, o...)
		},
		"delete many bad": func(o ...transport.RequestOption) error {
			return s.DeleteMany(ctx, "customerorder", docID, []string{"bad"}, o...)
		},
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, "customerorder", docID, nil, o...)
			return err
		},
		"update many no fields": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, "customerorder", docID, []batch.Update[eventfeed.NoteRequest]{{ID: noteID}}, o...)
			return err
		},
	}
	// Тип документа проверяется по форме в каждом методе.
	for _, typ := range []string{"", "CustomerOrder", "customer/order", "../x", "1supply", "supply?x=1"} {
		cases["list type "+typ] = func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, typ, docID, query.PageOptions{}, o...)
			return err
		}
		cases["get type "+typ] = func(o ...transport.RequestOption) error {
			_, err := s.Get(ctx, typ, docID, noteID, o...)
			return err
		}
		cases["create type "+typ] = func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, typ, docID, note, o...)
			return err
		}
		cases["update type "+typ] = func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, typ, docID, noteID, note, o...)
			return err
		}
		cases["delete type "+typ] = func(o ...transport.RequestOption) error { return s.Delete(ctx, typ, docID, noteID, o...) }
		cases["create many type "+typ] = func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, typ, docID, []eventfeed.NoteRequest{note}, o...)
			return err
		}
		cases["update many type "+typ] = func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, typ, docID, []batch.Update[eventfeed.NoteRequest]{{ID: noteID, Request: note}}, o...)
			return err
		}
		cases["delete many type "+typ] = func(o ...transport.RequestOption) error {
			return s.DeleteMany(ctx, typ, docID, []string{noteID}, o...)
		}
	}
	for name, invoke := range cases {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := eventfeed.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}
