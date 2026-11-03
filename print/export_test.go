package print_test

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
	"github.com/eteeet/moysklad-sdk-go/print"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const documentID = "a86708d2-f8d3-4e67-8f04-6101158da808"

func ref(href string, typ model.MetaType) model.Reference {
	return model.Reference{Meta: model.Meta{Href: href, Type: typ, MediaType: "application/json"}}
}

func exportRequest() print.ExportRequest {
	return print.ExportRequest{
		Template:  ref(api+"entity/demand/metadata/customtemplate/daca545a-1a35-11e7-8a84-bae500000001", "customtemplate"),
		Extension: "xls",
	}
}

func exportSetRequest() print.ExportSetRequest {
	return print.ExportSetRequest{Templates: []print.TemplateCopies{
		{Template: ref(api+"entity/demand/metadata/customtemplate/daca545a-1a35-11e7-8a84-bae500000001", "customtemplate"), Count: 2},
		{Template: ref(api+"entity/factureout/metadata/embeddedtemplate/3d2685b4-cf64-4fd1-87c8-e109966b364b", "embeddedtemplate"), Count: 3},
		{Template: ref(api+"entity/invoiceout/metadata/embeddedtemplate/6f3c9a47-6772-4944-9723-92d0d7be2a9c", "embeddedtemplate"), Count: 1},
	}}
}

func labelRequest() print.LabelRequest {
	org := ref(api+"entity/organization/107430bc-36e7-11e7-8a7f-40d000000090", "organization")
	org.Meta.MetadataHref = api + "entity/organization/metadata"
	org.Meta.UUIDHref = "https://online.moysklad.ru/app/#mycompany/edit?id=107430bc-36e7-11e7-8a7f-40d000000090"
	return print.LabelRequest{
		Organization: org,
		Count:        10,
		SalePrice:    print.LabelSalePrice{PriceType: ref(api+"context/companysettings/pricetype/672559f1-cbf3-11e1-9eb9-889ffa6f49fd", "pricetype")},
		Template:     ref(api+"entity/assortment/metadata/embeddedtemplate/f8e295eb-15c6-3184-b934-14fe90b3ea81", "embeddedtemplate"),
	}
}

func TestRequestGolden(t *testing.T) {
	for _, tc := range []struct {
		golden string
		input  any
	}{
		{"export.golden.json", exportRequest()},
		{"export-set.golden.json", exportSetRequest()},
		{"labels.golden.json", labelRequest()},
	} {
		for _, codec := range documenttest.Codecs {
			t.Run(tc.golden+"/"+codec.Name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, read(t, tc.golden))
			})
		}
	}
}

type exportCall struct {
	name, entity, golden, context string
	invoke                        func(context.Context, *print.Service, ...transport.RequestOption) (*print.Result, error)
}

func exportCalls() []exportCall {
	return []exportCall{
		{"export", "demand", "export.golden.json", "export demand " + documentID, func(ctx context.Context, s *print.Service, opts ...transport.RequestOption) (*print.Result, error) {
			return s.Export(ctx, "demand", documentID, exportRequest(), opts...)
		}},
		{"export set", "demand", "export-set.golden.json", "export set demand " + documentID, func(ctx context.Context, s *print.Service, opts ...transport.RequestOption) (*print.Result, error) {
			return s.ExportSet(ctx, "demand", documentID, exportSetRequest(), opts...)
		}},
		{"labels", "product", "labels.golden.json", "export labels product " + documentID, func(ctx context.Context, s *print.Service, opts ...transport.RequestOption) (*print.Result, error) {
			return s.ExportLabels(ctx, "product", documentID, labelRequest(), opts...)
		}},
	}
}

// Печать — один POST без следования redirect: 303 и 202 различимы, Location не запрашивается,
// в том числе на чужом origin; прочие ответы — ошибки с одним префиксом.
func TestExportHTTP(t *testing.T) {
	var foreignCalls atomic.Int32
	foreign := memnet.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { foreignCalls.Add(1) }))
	defer foreign.Close()

	for _, op := range exportCalls() {
		for _, tc := range []struct {
			name, location, body string
			status               int
			gzip                 bool
			want                 *print.Result // nil — ошибка
			api                  int           // код APIError или 0
		}{
			{name: "see other relative", status: 303, location: "/api/remap/1.2/download/f?sig=s", want: &print.Result{StatusCode: 303, Location: "{origin}/api/remap/1.2/download/f?sig=s"}},
			{name: "see other gzip", status: 303, gzip: true, location: "/api/remap/1.2/download/f", want: &print.Result{StatusCode: 303, Location: "{origin}/api/remap/1.2/download/f"}},
			{name: "see other foreign", status: 303, location: foreign.URL + "/tmp/f.pdf", want: &print.Result{StatusCode: 303, Location: foreign.URL + "/tmp/f.pdf"}},
			{name: "accepted", status: 202, location: "{origin}/api/remap/1.2/status/1", want: &print.Result{StatusCode: 202, Location: "{origin}/api/remap/1.2/status/1"}},
			{name: "see other without location", status: 303},
			{name: "bad location", status: 303, location: "http://[::1/f?sig=secret"},
			{name: "ok", status: 200, body: `{}`},
			{name: "found", status: 302, location: "/x"},
			{name: "bad request", status: 400, body: `{"errors":[{"error":"template","code":3000}]}`, api: 3000},
			{name: "server error", status: 500, body: `<html>`, api: -1},
		} {
			t.Run(op.name+"/"+tc.name, func(t *testing.T) {
				var calls atomic.Int32
				var origin string
				server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					body, _ := io.ReadAll(r.Body)
					if r.Method != http.MethodPost || r.URL.RequestURI() != "/api/remap/1.2/entity/"+op.entity+"/"+documentID+"/export" ||
						r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept-Encoding") != "gzip" ||
						r.Header.Get("Accept") != "application/json;charset=utf-8" || r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("request = %s %s %v", r.Method, r.URL.RequestURI(), r.Header)
					}
					documenttest.JSONEqual(t, body, read(t, op.golden))
					if tc.location != "" {
						w.Header().Set("Location", strings.ReplaceAll(tc.location, "{origin}", origin))
					}
					w.Header().Set("Content-Type", "application/json")
					payload := []byte(tc.body)
					if tc.gzip {
						var buf bytes.Buffer
						zw := gzip.NewWriter(&buf)
						zw.Write(payload)
						zw.Close()
						payload = buf.Bytes()
						w.Header().Set("Content-Encoding", "gzip")
					}
					w.WriteHeader(tc.status)
					w.Write(payload)
				}))
				defer server.Close()
				origin = server.URL
				auth, err := transport.BearerToken("test-token")
				if err != nil {
					t.Fatal(err)
				}
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Auth: auth})
				if err != nil {
					t.Fatal(err)
				}
				var info transport.ResponseInfo
				got, err := op.invoke(t.Context(), service(t, client), transport.WithResponseInfo(&info))
				if calls.Load() != 1 || foreignCalls.Load() != 0 {
					t.Fatalf("calls = %d, foreign = %d", calls.Load(), foreignCalls.Load())
				}
				// Неразбираемый Location net/http отвергает без ответа: StatusCode 0.
				if wantStatus := tc.status; (tc.name != "bad location" && info.StatusCode != wantStatus) || info.Attempts != 1 {
					t.Errorf("info = %+v", info)
				}
				if tc.want != nil {
					want := *tc.want
					want.Location = strings.ReplaceAll(want.Location, "{origin}", origin)
					if err != nil || got == nil || *got != want {
						t.Fatalf("result = %+v, err = %v, want %+v", got, err, want)
					}
					return
				}
				if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), op.context) {
					t.Fatalf("result = %+v, err = %v", got, err)
				}
				if strings.Contains(err.Error(), "secret") {
					t.Errorf("Location leaked into error: %v", err)
				}
				var apiErr *sdkerrors.APIError
				if errors.As(err, &apiErr) != (tc.api != 0) {
					t.Fatalf("API error = %v", err)
				}
				if tc.api > 0 && (apiErr.StatusCode != tc.status || len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != tc.api) {
					t.Errorf("API details = %+v", apiErr)
				}
			})
		}
		t.Run(op.name+"/local", func(t *testing.T) {
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*print.Result, error) {
				return op.invoke(ctx, service(t, c))
			})
		})
	}
}

// Обязательные поля и форма type/id проверяются до HTTP, ResponseInfo сброшен.
func TestExportValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(t, client)
	noTemplate := exportRequest()
	noTemplate.Template = model.Reference{}
	noExtension := exportRequest()
	noExtension.Extension = ""
	zeroCount := exportSetRequest()
	zeroCount.Templates[1].Count = 0
	setNoTemplate := exportSetRequest()
	setNoTemplate.Templates[2].Template = model.Reference{}
	labels := func(edit func(*print.LabelRequest)) print.LabelRequest {
		r := labelRequest()
		edit(&r)
		return r
	}
	ctx := t.Context()
	for name, invoke := range map[string]func(transport.RequestOption) (*print.Result, error){
		"export type": func(o transport.RequestOption) (*print.Result, error) {
			return s.Export(ctx, "Demand", documentID, exportRequest(), o)
		},
		"export id": func(o transport.RequestOption) (*print.Result, error) {
			return s.Export(ctx, "demand", "x", exportRequest(), o)
		},
		"export template": func(o transport.RequestOption) (*print.Result, error) {
			return s.Export(ctx, "demand", documentID, noTemplate, o)
		},
		"export extension": func(o transport.RequestOption) (*print.Result, error) {
			return s.Export(ctx, "demand", documentID, noExtension, o)
		},
		"set type": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportSet(ctx, "de-mand", documentID, exportSetRequest(), o)
		},
		"set empty": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportSet(ctx, "demand", documentID, print.ExportSetRequest{}, o)
		},
		"set count": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportSet(ctx, "demand", documentID, zeroCount, o)
		},
		"set template": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportSet(ctx, "demand", documentID, setNoTemplate, o)
		},
		"labels id": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportLabels(ctx, "product", "", labelRequest(), o)
		},
		"labels organization": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportLabels(ctx, "product", documentID, labels(func(r *print.LabelRequest) { r.Organization = model.Reference{} }), o)
		},
		"labels count": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportLabels(ctx, "product", documentID, labels(func(r *print.LabelRequest) { r.Count = 0 }), o)
		},
		"labels price type": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportLabels(ctx, "product", documentID, labels(func(r *print.LabelRequest) { r.SalePrice = print.LabelSalePrice{} }), o)
		},
		"labels template": func(o transport.RequestOption) (*print.Result, error) {
			return s.ExportLabels(ctx, "product", documentID, labels(func(r *print.LabelRequest) { r.Template = model.Reference{} }), o)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 303, Attempts: 1}
		got, err := invoke(transport.WithResponseInfo(&info))
		if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || info.StatusCode != 0 || info.Attempts != 0 {
			t.Errorf("%s: result = %+v, err = %v, info = %+v", name, got, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Errorf("calls = %d", calls.Load())
	}
}

// Download — отдельный GET по Location готовой формы через same-origin transport.Download.
func TestDownload(t *testing.T) {
	var foreignCalls atomic.Int32
	foreign := memnet.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { foreignCalls.Add(1) }))
	defer foreign.Close()
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodGet || r.URL.RequestURI() != "/api/remap/1.2/download/f?sig=s" || r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("request = %s %s %v", r.Method, r.URL.RequestURI(), r.Header)
		}
		io.WriteString(w, "%PDF-1.4")
	}))
	defer server.Close()
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	s := service(t, client)

	var buf bytes.Buffer
	var info transport.ResponseInfo
	n, err := s.Download(t.Context(), print.Result{StatusCode: 303, Location: server.URL + "/api/remap/1.2/download/f?sig=s"}, &buf, transport.WithResponseInfo(&info))
	if err != nil || n != 8 || buf.String() != "%PDF-1.4" || info.StatusCode != 200 || calls.Load() != 1 {
		t.Fatalf("n = %d, body = %q, err = %v, info = %+v, calls = %d", n, buf.String(), err, info, calls.Load())
	}

	for name, result := range map[string]print.Result{
		"accepted":       {StatusCode: 202, Location: server.URL + "/api/remap/1.2/status/1"},
		"no location":    {StatusCode: 303},
		"foreign origin": {StatusCode: 303, Location: foreign.URL + "/tmp/f.pdf"},
	} {
		buf.Reset()
		info := transport.ResponseInfo{StatusCode: 200}
		n, err := s.Download(t.Context(), result, &buf, transport.WithResponseInfo(&info))
		if err == nil || n != 0 || buf.Len() != 0 || strings.Count(err.Error(), "moysklad:") != 1 || info.StatusCode != 0 {
			t.Errorf("%s: n = %d, err = %v, info = %+v", name, n, err, info)
		}
	}
	if calls.Load() != 1 || foreignCalls.Load() != 0 {
		t.Errorf("calls = %d, foreign = %d", calls.Load(), foreignCalls.Load())
	}
}
