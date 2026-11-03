package trackingcode_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
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
		data, err := codec.Marshal([]trackingcode.Request{{ID: codeID, CIS: "012345678912345576", Type: "transportpack",
			TrackingCodes: []trackingcode.Request{{CIS: "c", Type: "trackingcode"}}}, {CIS: "010463003759026521LjJfNII5aXL-B", Type: "trackingcode"}})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`[{"id":"`+codeID+`","cis":"012345678912345576","type":"transportpack","trackingCodes":[{"cis":"c","type":"trackingcode"}]},`+
			`{"cis":"010463003759026521LjJfNII5aXL-B","type":"trackingcode"}]`))
	}
}

func TestHTTP(t *testing.T) {
	var got []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	if r, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{CIS: "c", Type: "trackingcode"}}); err != nil || r[0].Value.ID != codeID {
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

func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := trackingcode.NewService(c)
	ctx := t.Context()
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
			_, err := s.Save(ctx, "supply", docID, "x", []trackingcode.Request{{CIS: "c", Type: "t"}}, o...)
			return err
		},
		"save new cis": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{Type: "t"}}, o...)
			return err
		},
		"save bad id": func(o ...transport.RequestOption) error {
			_, err := s.Save(ctx, "supply", docID, posID, []trackingcode.Request{{ID: "x"}}, o...)
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
