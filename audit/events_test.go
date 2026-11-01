package audit_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/audit"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const contextID = "7944ef04-f831-11e5-7a69-971500188b19"

func TestEventFixtures(t *testing.T) {
	single := readFixture(t, "audit_event.json")
	page := readFixture(t, "audit_event_list.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got audit.Event
			if err := codec.Unmarshal(single, &got); err != nil {
				t.Fatal(err)
			}
			if got.EventType != "update" || got.EntityType != "product" || got.Name != "some product" || got.Audit == nil ||
				got.Entity == nil || got.UID != "admin@1" {
				t.Fatalf("event = %+v", got)
			}
			changes, err := got.FieldChanges()
			if err != nil || string(changes["weighed"].OldValue) != "false" || string(changes["weighed"].NewValue) != "true" {
				t.Fatalf("changes = %+v, err = %v", changes, err)
			}
			var list model.Collection[audit.Event]
			if err := codec.Unmarshal(page, &list); err != nil || len(list.Rows) == 0 || list.Rows[0].Diff == nil {
				t.Fatalf("list = %+v, err = %v", list, err)
			}
		})
	}
}

// Diff без потерь: большие целые, вложенные объекты, null, удалённые и неизвестные поля.
func TestEventDiffLossless(t *testing.T) {
	const diff = `{"sum":{"oldValue":12345678901234567890,"newValue":null},"positions":{"oldValue":[{"a":1}],"newValue":{"x":{"y":2}}},"removed":{"oldValue":"x"},"unknownMarker":{"oldValue":1,"newValue":2,"extra":true}}`
	for _, codec := range documenttest.Codecs {
		var got audit.Event
		if err := codec.Unmarshal([]byte(`{"eventType":"update","diff":`+diff+`}`), &got); err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, got.Diff, []byte(diff))
		// Неизвестное поле события сохраняется в Raw.
		var withExtra audit.Event
		if err := codec.Unmarshal([]byte(`{"eventType":"update","futureField":{"x":1}}`), &withExtra); err != nil ||
			!strings.Contains(string(withExtra.Raw), "futureField") {
			t.Fatalf("raw = %s, err = %v", withExtra.Raw, err)
		}
		changes, err := got.FieldChanges()
		if err != nil {
			t.Fatal(err)
		}
		if string(changes["sum"].OldValue) != "12345678901234567890" || string(changes["sum"].NewValue) != "null" ||
			string(changes["positions"].NewValue) != `{"x":{"y":2}}` || changes["removed"].NewValue != nil {
			t.Fatalf("%s: changes = %+v", codec.Name, changes)
		}
	}
}

func TestEventSpecificDiff(t *testing.T) {
	for _, tc := range []struct {
		name, diff string
		dst        any
		check      func(any) bool
	}{
		{"registration", `{"account":"shop","country":"RU"}`, &audit.RegistrationDiff{}, func(v any) bool {
			d := v.(*audit.RegistrationDiff)
			return d.Account == "shop" && d.Country == "RU"
		}},
		{"publication", `{"templateName":"Счёт","publicationHref":"https://p.example/1"}`, &audit.PublicationDiff{}, func(v any) bool {
			d := v.(*audit.PublicationDiff)
			return d.TemplateName == "Счёт" && d.PublicationHref == "https://p.example/1"
		}},
		{"email", `{"senderEmail":"a@x","targetEmail":"b@x","subjectEmail":"s","text":"t"}`, &audit.EmailDiff{}, func(v any) bool {
			d := v.(*audit.EmailDiff)
			return d.SenderEmail == "a@x" && d.TargetEmail == "b@x" && d.SubjectEmail == "s" && d.Text == "t"
		}},
	} {
		e := audit.Event{Diff: []byte(tc.diff)}
		if err := e.DecodeDiff(tc.dst); err != nil || !tc.check(tc.dst) {
			t.Errorf("%s: %+v, err = %v", tc.name, tc.dst, err)
		}
	}
	// Событие создания без diff; diff не той формы — ошибка без изменения Diff.
	if changes, err := (audit.Event{}).FieldChanges(); err != nil || len(changes) != 0 {
		t.Fatalf("create event: %v, %v", changes, err)
	}
	e := audit.Event{Diff: []byte(`{"account":"shop"}`)}
	if _, err := e.FieldChanges(); err == nil || !strings.HasPrefix(err.Error(), "moysklad: ") || string(e.Diff) != `{"account":"shop"}` {
		t.Fatalf("err = %v", err)
	}
}

func TestEventsEndpoints(t *testing.T) {
	page := readFixture(t, "audit_event_list.json")
	for _, op := range []struct {
		path   string
		invoke func(*audit.Service, ...transport.RequestOption) (*model.Collection[audit.Event], error)
	}{
		{"/audit/" + contextID + "/events", func(s *audit.Service, opts ...transport.RequestOption) (*model.Collection[audit.Event], error) {
			return s.EventsByContext(context.Background(), contextID, opts...)
		}},
		{"/entity/product/" + contextID + "/audit", func(s *audit.Service, opts ...transport.RequestOption) (*model.Collection[audit.Event], error) {
			return s.EventsByEntity(context.Background(), "product", contextID, opts...)
		}},
	} {
		for _, tc := range []struct {
			status int
			gzip   bool
		}{{200, false}, {200, true}, {403, false}} {
			status := tc.status
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet || r.URL.RequestURI() != op.path {
					t.Errorf("request = %s %s", r.Method, r.URL.RequestURI())
				}
				body := page
				if status != 200 {
					body = []byte(`{"errors":[{"error":"no access","code":1016}]}`)
				}
				if tc.gzip {
					var buf bytes.Buffer
					zw := gzip.NewWriter(&buf)
					zw.Write(body)
					zw.Close()
					body = buf.Bytes()
					w.Header().Set("Content-Encoding", "gzip")
				}
				w.WriteHeader(status)
				w.Write(body)
			}))
			var info transport.ResponseInfo
			got, err := op.invoke(newService(t, server.URL), transport.WithResponseInfo(&info))
			server.Close()
			if calls.Load() != 1 || info.StatusCode != status {
				t.Fatalf("%s: calls = %d, info = %+v", op.path, calls.Load(), info)
			}
			if status == 200 && (err != nil || len(got.Rows) == 0) ||
				status != 200 && (got != nil || sdkerrors.StatusCode(err) != status || strings.Count(err.Error(), "moysklad:") != 1) {
				t.Fatalf("%s %d: got = %+v, err = %v", op.path, status, got, err)
			}
		}
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*model.Collection[audit.Event], error) {
		s, err := audit.NewService(c)
		if err != nil {
			return nil, err
		}
		return s.EventsByEntity(ctx, "demand", contextID)
	})
}

// Небезопасный тип и не-UUID не отправляются; метаданные прошлого вызова сбрасываются.
func TestEventsValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	s := newService(t, server.URL)
	info := transport.ResponseInfo{StatusCode: 503}
	if _, err := s.EventsByContext(t.Context(), "../x", transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
		t.Errorf("bad context id: err = %v, info = %+v", err, info)
	}
	for _, typ := range []string{"../entity", "Product", "", "a/b", "customentity"} {
		if _, err := s.EventsByEntity(t.Context(), typ, contextID); err == nil {
			t.Errorf("type %q accepted", typ)
		}
	}
	if _, err := s.EventsByEntity(t.Context(), "product", "bad"); err == nil {
		t.Error("bad entity id accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
