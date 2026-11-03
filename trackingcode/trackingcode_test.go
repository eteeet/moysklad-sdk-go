package trackingcode_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/trackingcode"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	docID  = "161d0246-1477-11ec-ac18-000b00000001"
	posID  = "161d25a8-1477-11ec-ac18-000b00000002"
	codeID = "f591e101-1bc0-11ec-ac18-000c00000004"
	path   = "/entity/supply/" + docID + "/positions/" + posID + "/trackingCodes"
)

func TestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var got trackingcode.TrackingCode
		if err := codec.Unmarshal([]byte(`{"id":"`+codeID+`","cis":"012345678912345576","cis_1162":"x","type":"transportpack",`+
			`"trackingCodes":[{"cis":"010463003759026521LjJfNII5aXL-B","type":"trackingcode"}]}`), &got); err != nil ||
			got.ID != codeID || got.CIS1162 != "x" || got.Type != "transportpack" || len(got.TrackingCodes) != 1 || got.TrackingCodes[0].Type != "trackingcode" {
			t.Fatalf("%s: code = %+v, err = %v", codec.Name, got, err)
		}
		// Пример источника; изменение по id без cis и type их не отправляет.
		data, err := codec.Marshal([]trackingcode.Request{
			{ID: codeID, CIS: model.Value("012345678912345576"), Type: model.Value("transportpack"),
				TrackingCodes: []trackingcode.Request{{CIS: model.Value("c"), Type: model.Value("trackingcode")}}},
			{CIS: model.Value("010463003759026521LjJfNII5aXL-B"), Type: model.Value("trackingcode")},
			{ID: codeID, TrackingCodes: []trackingcode.Request{}},
		})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`[{"id":"`+codeID+`","cis":"012345678912345576","type":"transportpack","trackingCodes":[{"cis":"c","type":"trackingcode"}]},`+
			`{"cis":"010463003759026521LjJfNII5aXL-B","type":"trackingcode"},{"id":"`+codeID+`","trackingCodes":[]}]`))
	}
}

func TestHTTP(t *testing.T) {
	var got []string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		switch {
		case r.Method == http.MethodGet:
			io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[{"id":"`+codeID+`","cis":"c","type":"trackingcode"}]}`)
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		default:
			io.WriteString(w, `[{"id":"`+codeID+`","cis":"c","type":"trackingcode"}]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := trackingcode.NewService(c)
	ctx := t.Context()
	if l, err := s.List(ctx, "supply", docID, posID, trackingcode.Options{Limit: 50, Offset: 50, CodeType: trackingcode.CodeTypeAll}); err != nil || len(l.Rows) != 1 {
		t.Fatalf("list = %+v, err = %v", l, err)
	}
	if _, err := s.List(ctx, "supply", docID, posID, trackingcode.Options{}); err != nil {
		t.Fatal(err)
	}
	if r, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{CIS: model.Value("c"), Type: model.Value("trackingcode")}}); err != nil || r[0].Value.ID != codeID {
		t.Fatalf("save = %+v, err = %v", r, err)
	}
	if r, err := s.DeleteMany(ctx, "supply", docID, posID, []string{codeID}); err != nil || r[0].Info != "ok" {
		t.Fatalf("delete = %+v, err = %v", r, err)
	}
	want := []string{
		"GET " + path + "?codetype=all&limit=50&offset=50 ",
		"GET " + path + " ",
		"POST " + path + ` [{"cis":"c","type":"trackingcode"}]`,
		"POST " + path + `/delete [{"meta":{"href":"` + server.URL + path + "/" + codeID + `","type":"trackingcode","mediaType":"application/json"}}]`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// Ответы сервера для каждого метода: plain и gzip, ошибки API, неожиданный статус и тело;
// одна отправка, errors.As и единый префикс moysklad:.
func TestResponses(t *testing.T) {
	type call func(context.Context, *trackingcode.Service) error
	for _, op := range []struct {
		name, ok string
		invoke   call
	}{
		{"list", `{"meta":{"href":"x"},"rows":[]}`, func(ctx context.Context, s *trackingcode.Service) error {
			_, err := s.List(ctx, "supply", docID, posID, trackingcode.Options{})
			return err
		}},
		{"save", `[{"id":"` + codeID + `"}]`, func(ctx context.Context, s *trackingcode.Service) error {
			_, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{ID: codeID, CIS: model.Value("c")}})
			return err
		}},
		{"delete many", `[{"info":"ok"}]`, func(ctx context.Context, s *trackingcode.Service) error {
			_, err := s.DeleteMany(ctx, "supply", docID, posID, []string{codeID})
			return err
		}},
	} {
		for _, tc := range []struct {
			name, body string
			status     int
			gzip, ok   bool
			api        bool
		}{
			{"plain", op.ok, 200, false, true, false},
			{"gzip", op.ok, 200, true, true, false},
			{"api error", `{"errors":[{"error":"bad","code":3000}]}`, 400, false, false, true},
			{"rate limit", `{"errors":[{"error":"slow","code":1073}]}`, 429, false, false, true},
			{"unexpected status", op.ok, 201, false, false, false},
			{"null", `null`, 200, false, false, false},
			{"malformed", `{"meta":`, 200, false, false, false},
		} {
			t.Run(op.name+"/"+tc.name, func(t *testing.T) {
				var calls atomic.Int32
				server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					body := []byte(tc.body)
					if tc.gzip {
						var buf bytes.Buffer
						zw := gzip.NewWriter(&buf)
						zw.Write(body)
						zw.Close()
						body = buf.Bytes()
						w.Header().Set("Content-Encoding", "gzip")
					}
					w.WriteHeader(tc.status)
					w.Write(body)
				}))
				defer server.Close()
				c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				s, _ := trackingcode.NewService(c)
				err = op.invoke(t.Context(), s)
				if (err == nil) != tc.ok {
					t.Fatalf("err = %v", err)
				}
				var apiErr *sdkerrors.APIError
				if err != nil && (errors.As(err, &apiErr) != tc.api || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "supply "+docID)) {
					t.Fatalf("error = %v", err)
				}
				// GET повторяет только 429 и 5xx по политике; остальные ответы — одна отправка.
				if want := map[bool]int32{true: 3, false: 1}[op.name == "list" && tc.status == 429]; calls.Load() != want {
					t.Fatalf("calls = %d, want %d", calls.Load(), want)
				}
			})
		}
		t.Run(op.name+"/cancelled", func(t *testing.T) {
			c, err := transport.NewClient(transport.Config{BaseURL: "http://127.0.0.1:1"})
			if err != nil {
				t.Fatal(err)
			}
			s, _ := trackingcode.NewService(c)
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			if err := op.invoke(ctx, s); !errors.Is(err, context.Canceled) {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := trackingcode.NewService(c)
	ctx := t.Context()
	newCode := trackingcode.Request{CIS: model.Value("c"), Type: model.Value("t")}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"list type": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, "../x", docID, posID, trackingcode.Options{}, o...)
			return err
		},
		"list limit": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, "supply", docID, posID, trackingcode.Options{Limit: 101}, o...)
			return err
		},
		"list offset": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, "supply", docID, posID, trackingcode.Options{Offset: -1}, o...)
			return err
		},
		"save pos": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, "x", []trackingcode.Request{newCode}, o...)
			return err
		},
		"save new cis": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{Type: model.Value("t")}}, o...)
			return err
		},
		"save empty cis": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{CIS: model.Value(""), Type: model.Value("t")}}, o...)
			return err
		},
		"save bad id": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{ID: "x"}}, o...)
			return err
		},
		"save nested": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{ID: codeID, TrackingCodes: []trackingcode.Request{{CIS: model.Value("c")}}}}, o...)
			return err
		},
		"save empty": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, posID, nil, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, "supply", "x", posID, []string{codeID}, o...)
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
