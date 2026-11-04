package vendorapi_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/transport"
	"github.com/eteeet/moysklad-sdk-go/vendorapi"
)

const (
	appID     = "5f3c5489-6a17-48b7-9fe5-b2000eb807fe"
	accountID = "f088b0a7-9490-4a57-b804-393163e7680f"
	statusURI = "/api/vendor/1.0/apps/" + appID + "/" + accountID + "/status"
)

func TestAppStatusFixture(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got vendorapi.AppStatus
			if err := codec.Unmarshal(fixture(t, "status.json"), &got); err != nil {
				t.Fatal(err)
			}
			expiry := time.Date(2024, 1, 19, 18, 50, 12, 0, time.FixedZone("", 3*3600))
			sub := got.Subscription
			if got.Status != vendorapi.StatusActivated || got.Cause != vendorapi.CauseInstall || sub == nil ||
				sub.TariffID != "23ca69d4-2657-40c4-8ba1-6ce24ddeac2e" || sub.Trial || sub.TariffName != "Basic" ||
				sub.ExpiryMoment == nil || !sub.ExpiryMoment.Equal(expiry) || sub.NotForResale || sub.Partner {
				t.Fatalf("status = %+v, subscription = %+v", got, sub)
			}
			// В статусе токена нет; resource в примере — со слешем в конце.
			if len(got.Access) != 1 || got.Access[0].Resource != "https://api.moysklad.ru/api/remap/1.2/" ||
				len(got.Access[0].Scope) != 1 || got.Access[0].Scope[0] != "custom" || got.Access[0].AccessToken != "" {
				t.Fatalf("access = %+v", got.Access)
			}
			documenttest.JSONEqual(t, got.Access[0].Permissions, []byte(`{"viewDashboard":true,"viewAudit":true,"useAllWebhooks":true,
				"useAllAttributeMetadata":true,"supply":{"update":"ALL","view":"ALL"},
				"purchaseReturn":{"delete":"ALL","update":"ALL","approve":"ALL","view":"ALL","create":"ALL","print":"ALL"}}`))
		})
	}
}

// GET status повторяется RetryPolicy при 429 и 5xx; успех — 200 и непустой status.
func TestServiceStatus(t *testing.T) {
	cases := documenttest.ResponseCases(string(fixture(t, "status.json")))
	cases = append(cases, documenttest.ResponseCase{Name: "empty status", Status: 200, Body: `{"status":""}`},
		documenttest.ResponseCase{Name: "documented not installed", Status: 404, Body: `{"errors":[{"error":"not installed","code":2004}]}`, API: true, ErrorCode: 2004})
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var calls atomic.Int32
			s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet || r.URL.RequestURI() != statusURI || r.ContentLength != 0 || r.Header.Get("Content-Type") != "" {
					t.Errorf("request = %s %s, length = %d", r.Method, r.URL.RequestURI(), r.ContentLength)
				}
				outgoingClaims(t, r.Header.Get("Authorization"), "secret-key")
				respond(w, tc)
			})
			var info transport.ResponseInfo
			got, err := s.Status(t.Context(), appID, accountID, transport.WithResponseInfo(&info))
			want := int32(1)
			if tc.Status == 429 || tc.Status >= 500 {
				want = 3
			}
			if calls.Load() != want || info.StatusCode != tc.Status || info.Header.Get("X-Lognex-Auth") != "test" {
				t.Fatalf("calls = %d, info = %+v", calls.Load(), info)
			}
			if tc.Success {
				if err != nil || got == nil || got.Status != vendorapi.StatusActivated {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
				return
			}
			if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "vendor app status") {
				t.Fatalf("got = %+v, err = %v", got, err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.API || tc.ErrorCode != 0 && apiErr.Errors[0].Code != tc.ErrorCode {
				t.Fatalf("API error = %v", err)
			}
		})
	}
}

// Каждая попытка GET подписана заново: МойСклад отклонил бы повтор jti.
func TestServiceStatusRetryFreshJTI(t *testing.T) {
	var mu sync.Mutex
	var jtis []string
	s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, payload := outgoingClaims(t, r.Header.Get("Authorization"), "secret-key")
		mu.Lock()
		jtis = append(jtis, payload["jti"].(string))
		attempt := len(jtis)
		mu.Unlock()
		if attempt == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		io.WriteString(w, `{"status":"Activating","cause":"Install"}`)
	})
	got, err := s.Status(t.Context(), appID, accountID)
	if err != nil || got.Status != vendorapi.StatusActivating || got.Access != nil || got.Subscription != nil {
		t.Fatalf("got = %+v, err = %v", got, err)
	}
	if len(jtis) != 2 || jtis[0] == jtis[1] || !hexJTI.MatchString(jtis[0]) || !hexJTI.MatchString(jtis[1]) {
		t.Fatalf("jti = %v", jtis)
	}
}

// 36 символов, но не UUID: проверка длины их не отсечёт, а / и ? изменили бы путь.
var badIDs = [][2]string{{"", accountID}, {appID, ""}, {"not-a-uuid", accountID}, {appID, accountID + "/x"}, {appID + "x", accountID},
	{"5f3c5489-6a17-48b7-9fe5-b2000eb8?x=1", accountID}, {appID, "f088b0a7-9490-4a57-b804-393163e7680/"},
	{"zf3c5489-6a17-48b7-9fe5-b2000eb807fe", accountID}, {"5f3c5489x6a17-48b7-9fe5-b2000eb807fe", accountID}}

func TestServiceStatusRejectsIDs(t *testing.T) {
	var calls atomic.Int32
	s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	for _, ids := range badIDs {
		info := transport.ResponseInfo{StatusCode: 200, Header: http.Header{"X": {"stale"}}}
		got, err := s.Status(t.Context(), ids[0], ids[1], transport.WithResponseInfo(&info))
		if err == nil || got != nil || info.StatusCode != 0 || info.Header != nil {
			t.Errorf("%q: got = %v, err = %v, info = %+v", ids, got, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Errorf("calls = %d", calls.Load())
	}
}

func TestServiceStatusLocalErrors(t *testing.T) {
	documenttest.LocalErrors(t, func(ctx context.Context, client *transport.Client) (*vendorapi.AppStatus, error) {
		s, err := vendorapi.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Status(ctx, appID, accountID)
	})
}

// PUT {"status"} с пустым ответом 200; PUT не повторяется. 404 — решение не установлено,
// 409 — переход по жизненному циклу невозможен.
func TestServiceSetStatus(t *testing.T) {
	for _, tc := range []documenttest.ResponseCase{
		{Name: "ok", Status: 200, Success: true},
		{Name: "not installed", Status: 404, Body: `{"errors":[{"error":"not installed","code":2004}]}`, API: true, ErrorCode: 2004},
		{Name: "no transition", Status: 409, Body: `{"errors":[{"error":"conflict","code":2005}]}`, API: true},
		{Name: "server error", Status: 500, Body: `{"errors":[{"error":"server","code":1000}]}`, API: true},
		{Name: "unexpected status", Status: 204},
	} {
		for _, status := range []vendorapi.Status{vendorapi.StatusActivating, vendorapi.StatusSettingsRequired, vendorapi.StatusActivated} {
			var calls atomic.Int32
			s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				body, _ := io.ReadAll(r.Body)
				if r.Method != http.MethodPut || r.URL.RequestURI() != statusURI || r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("request = %s %s %q", r.Method, r.URL.RequestURI(), r.Header.Get("Content-Type"))
				}
				documenttest.JSONEqual(t, body, []byte(`{"status":"`+string(status)+`"}`))
				outgoingClaims(t, r.Header.Get("Authorization"), "secret-key")
				respond(w, tc)
			})
			var info transport.ResponseInfo
			err := s.SetStatus(t.Context(), appID, accountID, status, transport.WithResponseInfo(&info))
			if calls.Load() != 1 || info.StatusCode != tc.Status || info.Header.Get("X-Lognex-Auth") != "test" {
				t.Fatalf("%s %s: calls = %d, info = %+v", tc.Name, status, calls.Load(), info)
			}
			var apiErr *sdkerrors.APIError
			if tc.Success != (err == nil) || errors.As(err, &apiErr) != tc.API || tc.API && apiErr.StatusCode != tc.Status ||
				tc.ErrorCode != 0 && apiErr.Errors[0].Code != tc.ErrorCode ||
				err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "set vendor app status")) {
				t.Fatalf("%s %s: err = %v", tc.Name, status, err)
			}
		}
	}
}

// Остальные статусы выставляет только МойСклад; ошибка — до HTTP и обнуляет ResponseInfo.
func TestServiceSetStatusRejects(t *testing.T) {
	var calls atomic.Int32
	s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	check := func(name string, err error, info transport.ResponseInfo) {
		t.Helper()
		if err == nil || info.StatusCode != 0 || info.Header != nil {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	for _, status := range []vendorapi.Status{"", vendorapi.StatusActivationFailed, vendorapi.StatusDeactivating,
		vendorapi.StatusDeactivationFailed, vendorapi.StatusSuspended, "activated", "Unknown"} {
		info := transport.ResponseInfo{StatusCode: 200, Header: http.Header{"X": {"stale"}}}
		check(string(status), s.SetStatus(t.Context(), appID, accountID, status, transport.WithResponseInfo(&info)), info)
	}
	for _, ids := range badIDs {
		info := transport.ResponseInfo{StatusCode: 200, Header: http.Header{"X": {"stale"}}}
		check(ids[0]+" "+ids[1], s.SetStatus(t.Context(), ids[0], ids[1], vendorapi.StatusActivated, transport.WithResponseInfo(&info)), info)
	}
	if calls.Load() != 0 {
		t.Errorf("calls = %d", calls.Load())
	}
}

func TestServiceSetStatusLocalErrors(t *testing.T) {
	documenttest.LocalErrors(t, func(ctx context.Context, client *transport.Client) (*struct{}, error) {
		s, err := vendorapi.NewService(client)
		if err != nil {
			return nil, err
		}
		if err := s.SetStatus(ctx, appID, accountID, vendorapi.StatusActivated); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})
}
