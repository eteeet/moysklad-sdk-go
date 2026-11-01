package metadata_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const attrID = "c92a0e25-120c-11f1-ac12-001100000002"

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestEntityFixtures(t *testing.T) {
	for _, tc := range []struct {
		file        string
		attrMeta    bool // коллекция {meta[, rows]}
		attrRows    int
		hasStates   bool
		extraRawKey string
	}{
		{file: "product_metadata.json", attrMeta: true},
		{file: "product_metadata_expanded.json", attrMeta: true, attrRows: 2},
		{file: "document_metadata.json", attrMeta: true, hasStates: true},
		{file: "document_metadata_expanded.json", attrMeta: true, attrRows: -1, hasStates: true},
		{file: "counterparty_metadata.json", attrMeta: true, hasStates: true, extraRawKey: `"tags"`},
		{file: "counterparty_metadata_expanded.json", attrMeta: true, attrRows: -1, hasStates: true},
		{file: "counterparty_metadata_minimum.json"},
	} {
		data := readFile(t, "upstream/"+tc.file)
		for _, codec := range documenttest.Codecs {
			t.Run(tc.file+"/"+codec.Name, func(t *testing.T) {
				var got metadata.Entity
				if err := codec.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				if got.Meta.Href == "" || (got.Attributes != nil) != tc.attrMeta {
					t.Fatalf("metadata = %+v", got)
				}
				if tc.attrMeta && (got.Attributes.Meta == nil || tc.attrRows >= 0 && len(got.Attributes.Rows) != tc.attrRows ||
					tc.attrRows < 0 && len(got.Attributes.Rows) == 0) {
					t.Fatalf("attributes = %+v", got.Attributes)
				}
				if tc.hasStates != (len(got.States) > 0) || tc.hasStates && (got.States[0].Color == 0 || got.States[0].StateType == "") {
					t.Fatalf("states = %+v", got.States)
				}
				documenttest.JSONEqual(t, got.Raw, data)
				if tc.extraRawKey != "" && !strings.Contains(string(got.Raw), tc.extraRawKey) {
					t.Fatalf("raw lost %s", tc.extraRawKey)
				}
			})
		}
	}
}

// Большинство примеров источника отдаёт attributes массивом, без id у описаний.
func TestEntityAttributesArray(t *testing.T) {
	const data = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata"},"attributes":[` +
		`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/customerorder/metadata/attributes/1"},"name":"Ref","type":"customentity",` +
		`"required":false,"customEntityMeta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/companysettings/metadata/customEntities/2"}}],"createShared":true}`
	for _, codec := range documenttest.Codecs {
		var got metadata.Entity
		if err := codec.Unmarshal([]byte(data), &got); err != nil {
			t.Fatalf("%s: %v", codec.Name, err)
		}
		if got.Attributes == nil || got.Attributes.Meta != nil || len(got.Attributes.Rows) != 1 ||
			got.Attributes.Rows[0].CustomEntityMeta == nil || got.Attributes.Rows[0].ID != "" || !got.CreateShared {
			t.Fatalf("%s: metadata = %+v", codec.Name, got)
		}
	}
}

func TestCompanySettingsFixture(t *testing.T) {
	data := readFile(t, "upstream/company_settings_metadata.json")
	for _, codec := range documenttest.Codecs {
		var got metadata.CompanySettings
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if len(got.CustomEntities) == 0 || got.CustomEntities[0].ID == "" || got.CustomEntities[0].EntityMeta.Href == "" ||
			got.CustomEntities[0].Attributes == nil {
			t.Fatalf("%s: settings = %+v", codec.Name, got)
		}
		documenttest.JSONEqual(t, got.Raw, data)
	}
}

func service(c *transport.Client) *metadata.Service {
	s, err := metadata.NewService(c)
	if err != nil {
		panic(err)
	}
	return s
}

// Общий HTTP-контракт: plain/gzip, ошибки API, неверные тела, один запрос, локальные ошибки.
func TestHTTPContract(t *testing.T) {
	entityBody := string(readFile(t, "upstream/document_metadata.json"))
	attribute := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/metadata/attributes/` + attrID + `"},"name":"test2","type":"long"}`
	attributes := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/metadata/attributes"},"rows":[` + attribute + `]}`
	catalog := `{"demand":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/demand/metadata"}}}`
	for _, op := range []struct {
		name, path, body string
		invoke           func(context.Context, *transport.Client) (any, error)
	}{
		{"catalog", "entity/metadata", catalog, func(ctx context.Context, c *transport.Client) (any, error) {
			got, err := service(c).Catalog(ctx, nil)
			if err != nil {
				return nil, err
			}
			return got, nil
		}},
		{"get", "entity/demand/metadata", entityBody, func(ctx context.Context, c *transport.Client) (any, error) {
			return nilIfErr(service(c).Get(ctx, "demand"))
		}},
		{"customentity", "entity/customentity/" + attrID + "/metadata", entityBody, func(ctx context.Context, c *transport.Client) (any, error) {
			return nilIfErr(service(c).GetCustomEntity(ctx, attrID))
		}},
		{"companysettings", "context/companysettings/metadata", string(readFile(t, "upstream/company_settings_metadata.json")), func(ctx context.Context, c *transport.Client) (any, error) {
			return nilIfErr(service(c).GetCompanySettings(ctx))
		}},
		{"attributes", "entity/product/metadata/attributes?limit=10", attributes, func(ctx context.Context, c *transport.Client) (any, error) {
			return nilIfErr(service(c).ListAttributes(ctx, "product", metadata.AttributesPage{Limit: 10}))
		}},
		{"attribute", "entity/product/metadata/attributes/" + attrID, attribute, func(ctx context.Context, c *transport.Client) (any, error) {
			return nilIfErr(service(c).GetAttribute(ctx, "product", attrID))
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			invoke := func(ctx context.Context, c *transport.Client) (*any, error) {
				got, err := op.invoke(ctx, c)
				if err != nil {
					return nil, err
				}
				return &got, nil
			}
			for _, tc := range documenttest.ResponseCases(op.body) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, "GET", op.path, nil, invoke, func(*any) {})
				})
			}
			documenttest.LocalErrors(t, invoke)
		})
	}
}

func nilIfErr[T any](v *T, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Фильтр каталога — type=a;type=b одним параметром filter.
func TestCatalogQuery(t *testing.T) {
	var uri string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri = r.URL.RequestURI()
		w.Write([]byte(`{"service":{"meta":{"href":"h1"}},"demand":{"meta":{"href":"h2"}}}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	got, err := service(client).Catalog(t.Context(), []string{"service", "demand"})
	if err != nil || uri != "/entity/metadata?filter=type%3Dservice%3Btype%3Ddemand" || len(got) != 2 || got["demand"].Meta.Href != "h2" {
		t.Fatalf("got = %+v, err = %v, uri = %s", got, err, uri)
	}
}

// Отказ в доступе и отсутствие метаданных у типа — APIError с кодом и ResponseInfo.
func TestServerErrors(t *testing.T) {
	for _, tc := range []struct {
		status, code int
	}{{401, 1056}, {403, 1016}, {404, 1003}} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
			w.Write([]byte(`{"errors":[{"error":"x","code":` + strconv.Itoa(tc.code) + `}]}`))
		}))
		client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
		if err != nil {
			t.Fatal(err)
		}
		var info transport.ResponseInfo
		got, err := service(client).Get(t.Context(), "retailshift", transport.WithResponseInfo(&info))
		server.Close()
		if got != nil || !sdkerrors.HasCode(err, tc.code) || info.StatusCode != tc.status {
			t.Fatalf("%d: got = %+v, err = %v", tc.status, got, err)
		}
	}
}

func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(client)
	for name, call := range map[string]func() error{
		"catalog type":    func() error { _, err := s.Catalog(t.Context(), []string{"product", "../x"}); return err },
		"get type":        func() error { _, err := s.Get(t.Context(), "../x"); return err },
		"customentity id": func() error { _, err := s.GetCustomEntity(t.Context(), "bad"); return err },
		"list type": func() error {
			_, err := s.ListAttributes(t.Context(), "Product", metadata.AttributesPage{})
			return err
		},
		"list limit": func() error {
			_, err := s.ListAttributes(t.Context(), "product", metadata.AttributesPage{Limit: 1001})
			return err
		},
		"attribute id":   func() error { _, err := s.GetAttribute(t.Context(), "product", "bad"); return err },
		"attribute type": func() error { _, err := s.GetAttribute(t.Context(), "a/b", attrID); return err },
	} {
		if err := call(); err == nil || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatal("invalid arguments sent")
	}
	if _, err := metadata.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}

// Регрессия ревью 3306518: AttributeList записывается в той же форме, в какой читается,
// и строки не теряются при сохранении и повторном чтении.
func TestAttributeListRoundTrip(t *testing.T) {
	const row = `{"meta":{"href":"https://x/attributes/1"},"id":"1","name":"Цвет","type":"string","required":false,"show":true,"description":""}`
	for _, tc := range []struct {
		name, data string
		rows       int
	}{
		{"array", `[` + row + `]`, 1},
		{"empty array", `[]`, 0},
		{"collection meta only", `{"meta":{"href":"https://x/attributes","size":3,"limit":0,"offset":0}}`, 0},
		{"collection with rows", `{"meta":{"href":"https://x/attributes","size":1,"limit":0,"offset":0},"rows":[` + row + `]}`, 1},
	} {
		for _, codec := range documenttest.Codecs {
			t.Run(tc.name+"/"+codec.Name, func(t *testing.T) {
				var list metadata.AttributeList
				if err := codec.Unmarshal([]byte(tc.data), &list); err != nil {
					t.Fatal(err)
				}
				out, err := codec.Marshal(list)
				if err != nil {
					t.Fatal(err)
				}
				var again metadata.AttributeList
				if err := codec.Unmarshal(out, &again); err != nil {
					t.Fatal(err)
				}
				if len(again.Rows) != tc.rows || (again.Meta == nil) != (list.Meta == nil) {
					t.Fatalf("round trip = %s, rows = %d", out, len(again.Rows))
				}
				if tc.rows > 0 {
					documenttest.JSONEqual(t, out, []byte(tc.data))
				}
			})
		}
	}
}

// Сохранённые Entity и CompanySettings читаются обратно с теми же описаниями доп. полей.
func TestMetadataRoundTrip(t *testing.T) {
	entityRows := func(e metadata.Entity) int { return len(e.Attributes.Rows) }
	settingsRows := func(c metadata.CompanySettings) int {
		n := 0
		for _, ce := range c.CustomEntities {
			if ce.Attributes != nil {
				n += 1 + len(ce.Attributes.Rows)
			}
		}
		return n
	}
	roundTrip(t, readFile(t, "upstream/product_metadata_expanded.json"), entityRows)
	roundTrip(t, readFile(t, "upstream/company_settings_metadata.json"), settingsRows)
}

func roundTrip[T any](t *testing.T, data []byte, rows func(T) int) {
	t.Helper()
	for _, codec := range documenttest.Codecs {
		var first, second T
		if err := codec.Unmarshal(data, &first); err != nil {
			t.Fatal(err)
		}
		out, err := codec.Marshal(first)
		if err != nil {
			t.Fatal(err)
		}
		if err := codec.Unmarshal(out, &second); err != nil {
			t.Fatal(err)
		}
		if rows(first) == 0 || rows(second) != rows(first) {
			t.Fatalf("%s: rows before = %d, after = %d: %s", codec.Name, rows(first), rows(second), out)
		}
	}
}
