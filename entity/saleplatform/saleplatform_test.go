package saleplatform_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/saleplatform"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "39aaefe6-0304-4b85-ad3c-e8b93549f426"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestSalePlatformFixture(t *testing.T) {
	data := read(t, "upstream/sale_platform.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got saleplatform.SalePlatform
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "Avito" || got.SalePlatformGroup != "MARKETPLACES" {
				t.Fatalf("saleplatform = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse saleplatform.SalePlatform
			if err := codec.Unmarshal([]byte(`{"id":"x","futureField":{"a":1}}`), &sparse); err != nil || sparse.ID != "x" {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func service(c *transport.Client) *saleplatform.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := saleplatform.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/saleplatform/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*saleplatform.SalePlatform, error)
	}{
		{"get", "GET", "entity/saleplatform/" + testID, nil, func(ctx context.Context, c *transport.Client) (*saleplatform.SalePlatform, error) {
			return service(c).Get(ctx, testID)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *saleplatform.SalePlatform) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("saleplatform = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Ошибки входа — до HTTP; ResponseInfo прошлого вызова сброшен.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
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

func TestServiceList(t *testing.T) {
	listtest.Run(t, "saleplatform", read(t, "upstream/sale_platform.json"), "", "", func(c *transport.Client) listtest.ListFunc[saleplatform.SalePlatform] {
		return service(c).List
	}, func(t *testing.T, got saleplatform.SalePlatform) {
		if got.ID != testID || got.SalePlatformGroup != "MARKETPLACES" {
			t.Fatalf("row = %+v", got)
		}
	})
}
