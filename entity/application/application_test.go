package application_test

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/application"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const appID = "b58a6312-f958-11e9-ac12-000a00000020"

func checkApplication(t *testing.T, got *application.Application) {
	t.Helper()
	if got.ID != appID || got.AppUID != "test.moysklad@reqwy1" ||
		got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/entity/application/"+appID ||
		got.Meta.MetadataHref != "https://api.moysklad.ru/api/remap/1.2/entity/application/metadata" ||
		got.Meta.Type != "application" || got.Meta.MediaType != "application/json" {
		t.Fatalf("application = %+v", got)
	}
}

func TestFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/custom/application/application.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got application.Application
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			checkApplication(t, &got)
			// appUid — строка как есть, без разбора на alias/vendor/account.
			var raw application.Application
			if err := codec.Unmarshal([]byte(`{"id":"x","appUid":"a.b@c@d"}`), &raw); err != nil || raw.AppUID != "a.b@c@d" {
				t.Fatalf("appUid = %q, err = %v", raw.AppUID, err)
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	body, err := os.ReadFile("../../testdata/custom/application/application.json")
	if err != nil {
		t.Fatal(err)
	}
	invoke := func(ctx context.Context, c *transport.Client) (*application.Application, error) {
		s, err := application.NewService(c)
		if err != nil {
			return nil, err
		}
		return s.Get(ctx, appID)
	}
	for _, tc := range documenttest.ResponseCases(string(body)) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/application/"+appID, nil, invoke, func(got *application.Application) {
				checkApplication(t, got)
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

// Неверный id отклоняется до HTTP; ResponseInfo сбрасывается.
func TestLocalValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, err := application.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"", "b58a6312-f958-11e9-ac12-000a0000002", "../../../context/application-ac12-0", "test.moysklad@reqwy1"} {
		info := transport.ResponseInfo{StatusCode: 999, Attempts: 3}
		got, err := s.Get(t.Context(), id, transport.WithResponseInfo(&info))
		if got != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get application") {
			t.Errorf("%q: got = %+v, err = %v", id, got, err)
		}
		if info.StatusCode != 0 || info.Attempts != 0 {
			t.Errorf("%q: response info not reset: %+v", id, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid id sent: %d", calls.Load())
	}
	if _, err := application.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}
