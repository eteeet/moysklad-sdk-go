package pricetype_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/pricetype"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "88888888-8888-8888-8888-888888888888"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func fixture(t *testing.T) model.PriceType {
	t.Helper()
	var p model.PriceType
	if err := documenttest.Codecs[1].Unmarshal(read(t, "upstream/price_type.json"), &p); err != nil {
		t.Fatal(err)
	}
	return p
}

func code(s string) *string { return &s }

func replaceItems(t *testing.T) []pricetype.ReplaceItem {
	t.Helper()
	wholesale := model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/a8967d6b-b026-11e7-9464-d04800000000",
		Type: "pricetype", MediaType: "application/json"}
	return []pricetype.ReplaceItem{
		pricetype.Existing(fixture(t)),
		// Переименование с пустым внешним кодом: сервер сгенерирует новый.
		{Meta: &wholesale, Name: "Оптовая цена", ExternalCode: code("")},
		{Name: "Цена для друзей", ExternalCode: code("friends")},
	}
}

func TestFixtureAndRequest(t *testing.T) {
	want := read(t, "custom/pricetype/replace.golden.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var p model.PriceType
			if err := codec.Unmarshal(read(t, "upstream/price_type.json"), &p); err != nil || p.ID != testID || p.ExternalCode != "default" {
				t.Fatalf("price type = %+v, err = %v", p, err)
			}
			got, err := codec.Marshal(replaceItems(t))
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, got, want)
		})
	}
}

const listBody = `[{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/` + testID + `"},"id":"` + testID +
	`","name":"Цена продажи","externalCode":"default"},{"meta":{"href":"h2"},"id":"b","name":"Цена для друзей","externalCode":"friends"}]`

const oneBody = `{"meta":{"href":"h"},"id":"` + testID + `","name":"Цена продажи","externalCode":"default"}`

type call struct {
	name, method, path string
	body               []byte
	success            string
	invoke             func(context.Context, *pricetype.Service) (int, error) // число типов в ответе
}

func calls(t *testing.T) []call {
	t.Helper()
	return []call{
		{"list", "GET", "/context/companysettings/pricetype", nil, listBody, func(ctx context.Context, s *pricetype.Service) (int, error) {
			got, err := s.List(ctx)
			return len(got), err
		}},
		{"get", "GET", "/context/companysettings/pricetype/" + testID, nil, oneBody, func(ctx context.Context, s *pricetype.Service) (int, error) {
			return one(s.Get(ctx, testID))
		}},
		{"default", "GET", "/context/companysettings/pricetype/default", nil, oneBody, func(ctx context.Context, s *pricetype.Service) (int, error) {
			return one(s.GetDefault(ctx))
		}},
		{"replace", "POST", "/context/companysettings/pricetype", read(t, "custom/pricetype/replace.golden.json"), listBody, func(ctx context.Context, s *pricetype.Service) (int, error) {
			got, err := s.ReplaceAll(ctx, replaceItems(t))
			return len(got), err
		}},
	}
}

// one — число типов цен в ответе одиночного метода: 1, если id на месте.
func one(p *model.PriceType, err error) (int, error) {
	if err != nil || p.ID == "" {
		return 0, err
	}
	return 1, nil
}

func newService(t *testing.T, url string) *pricetype.Service {
	t.Helper()
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: url, Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	s, err := pricetype.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// HTTP-контракт: точные method/path/headers/body, plain/gzip, ошибки API, неверные тела, один запрос.
func TestHTTP(t *testing.T) {
	for _, op := range calls(t) {
		for _, tc := range []struct {
			name, body string
			status     int
			gzip, ok   bool
			apiErr     bool
		}{
			{"plain", op.success, 200, false, true, false},
			{"gzip", op.success, 200, true, true, false},
			{"forbidden", `{"errors":[{"error":"admin only","code":1017}]}`, 403, false, false, true},
			{"too many", `{"errors":[{"error":"more than 100","code":3000}]}`, 412, false, false, true},
			{"wrong shape", `[]`, 200, false, false, false},
			{"null", `null`, 200, false, false, false},
			{"malformed", `[{"id":`, 200, false, false, false},
			{"no content", ``, 204, false, false, false},
		} {
			t.Run(op.name+"/"+tc.name, func(t *testing.T) {
				var n atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					n.Add(1)
					body := new(bytes.Buffer)
					body.ReadFrom(r.Body)
					if r.Method != op.method || r.URL.RequestURI() != op.path || r.Header.Get("Authorization") != "Bearer test-token" ||
						r.Header.Get("Accept-Encoding") != "gzip" {
						t.Errorf("request = %s %s %v", r.Method, r.URL.RequestURI(), r.Header)
					}
					if op.body != nil {
						documenttest.JSONEqual(t, body.Bytes(), op.body)
					} else if body.Len() != 0 {
						t.Errorf("GET body = %s", body)
					}
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
				count, err := op.invoke(t.Context(), newService(t, server.URL))
				if n.Load() != 1 {
					t.Fatalf("calls = %d", n.Load())
				}
				if tc.ok {
					if err != nil || count == 0 {
						t.Fatalf("count = %d, err = %v", count, err)
					}
					return
				}
				if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || tc.apiErr != (sdkerrors.StatusCode(err) == tc.status) {
					t.Fatalf("err = %v", err)
				}
			})
		}
		documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*int, error) {
			s, err := pricetype.NewService(c)
			if err != nil {
				return nil, err
			}
			n, err := op.invoke(ctx, s)
			if err != nil {
				return nil, err
			}
			return &n, nil
		})
	}
}

// Пустой и слишком длинный список, новый тип без имени — ошибки до запроса.
func TestReplaceAllValidation(t *testing.T) {
	var n atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { n.Add(1) }))
	defer server.Close()
	s := newService(t, server.URL)
	tooMany := make([]pricetype.ReplaceItem, pricetype.MaxPriceTypes+1)
	for i := range tooMany {
		tooMany[i] = pricetype.ReplaceItem{Name: "p"}
	}
	for name, items := range map[string][]pricetype.ReplaceItem{
		"empty":            nil,
		"too many":         tooMany,
		"new no name":      {{ExternalCode: code("x")}},
		"blank name":       {{Name: "   "}},
		"existing no name": {{Meta: &model.Meta{Href: "h"}}},
		"existing no href": {pricetype.Existing(model.PriceType{Name: "x"})},
	} {
		if _, err := s.ReplaceAll(t.Context(), items); err == nil {
			t.Errorf("%s accepted", name)
		}
	}
	if _, err := s.Get(t.Context(), "default"); err == nil {
		t.Error("non-UUID id accepted")
	}
	if n.Load() != 0 {
		t.Fatalf("calls = %d", n.Load())
	}
	if _, err := pricetype.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}
