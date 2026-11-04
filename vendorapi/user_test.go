package vendorapi_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/transport"
	"github.com/eteeet/moysklad-sdk-go/vendorapi"
)

const userToken = "a1b2c3d4e5f6478901234567890abcdef1234567"

func TestUserContextFixture(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got vendorapi.UserContext
			if err := codec.Unmarshal(fixture(t, "user-context.json"), &got); err != nil {
				t.Fatal(err)
			}
			if got != (vendorapi.UserContext{AccountID: "b0b309ee-13e3-11e9-9109-f8fc0001f188", UserID: "b0a02321-13e3-11e9-912f-f3d4002516e3",
				UserUID: "admin@elonmask", Role: vendorapi.RoleAdmin}) {
				t.Fatalf("%+v", got)
			}
			// Значения ролей — дословно из таблицы источника; новая роль не теряется.
			var future vendorapi.UserContext
			if err := codec.Unmarshal([]byte(`{"role":"auditor"}`), &future); err != nil || future.Role != "auditor" {
				t.Fatalf("%+v, %v", future, err)
			}
			for got, want := range map[vendorapi.Role]string{vendorapi.RoleAdmin: "admin", vendorapi.RoleCashier: "cashier",
				vendorapi.RoleWorker: "worker", vendorapi.RoleIndividual: "individual"} {
				if string(got) != want {
					t.Errorf("role %q, want %q", got, want)
				}
			}
		})
	}
}

// Токен — только в теле POST; POST одноразового токена не повторяется даже при RetryPolicy и 5xx.
func TestServiceUserContext(t *testing.T) {
	cases := documenttest.ResponseCases(string(fixture(t, "user-context.json")))
	cases = append(cases,
		documenttest.ResponseCase{Name: "documented empty token", Status: 400, Body: `{"errors":[{"error":"empty","code":3008}]}`, API: true, ErrorCode: 3008},
		documenttest.ResponseCase{Name: "documented jwt", Status: 401, Body: `{"errors":[{"error":"auth","code":1056}]}`, API: true},
		documenttest.ResponseCase{Name: "documented used token", Status: 404, Body: `{"errors":[{"error":"used","code":3007}]}`, API: true, ErrorCode: 3007},
		documenttest.ResponseCase{Name: "no account", Status: 200, Body: `{"userId":"u","userUid":"x","role":"admin"}`},
		documenttest.ResponseCase{Name: "no user", Status: 200, Body: `{"accountId":"a","userUid":"x","role":"admin"}`},
		documenttest.ResponseCase{Name: "no uid", Status: 200, Body: `{"accountId":"a","userId":"u","role":"admin"}`},
		documenttest.ResponseCase{Name: "no role", Status: 200, Body: `{"accountId":"a","userId":"u","userUid":"x"}`},
	)
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var calls atomic.Int32
			s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				body, _ := io.ReadAll(r.Body)
				if r.Method != http.MethodPost || r.URL.RequestURI() != "/api/vendor/1.0/context/user" || r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("request = %s %s", r.Method, r.URL.RequestURI())
				}
				for key, values := range r.Header {
					if strings.Contains(strings.Join(values, " "), userToken) {
						t.Errorf("token in header %s", key)
					}
				}
				documenttest.JSONEqual(t, body, []byte(`{"token":"`+userToken+`"}`))
				outgoingClaims(t, r.Header.Get("Authorization"), "secret-key")
				respond(w, tc)
			})
			var info transport.ResponseInfo
			got, err := s.UserContext(t.Context(), userToken, transport.WithResponseInfo(&info))
			if calls.Load() != 1 || info.StatusCode != tc.Status || info.Header.Get("X-Lognex-Auth") != "test" {
				t.Fatalf("calls = %d, info = %+v", calls.Load(), info)
			}
			if tc.Success {
				if err != nil || got == nil || got.UserUID != "admin@elonmask" {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
				return
			}
			if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "vendor user context") ||
				strings.Contains(err.Error(), userToken) {
				t.Fatalf("got = %+v, err = %v", got, err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.API || tc.ErrorCode != 0 && apiErr.Errors[0].Code != tc.ErrorCode {
				t.Fatalf("API error = %v", err)
			}
		})
	}
}

func TestServiceUserContextRejectsEmpty(t *testing.T) {
	var calls atomic.Int32
	s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	info := transport.ResponseInfo{StatusCode: 200, Header: http.Header{"X": {"stale"}}}
	got, err := s.UserContext(t.Context(), "", transport.WithResponseInfo(&info))
	if err == nil || got != nil || info.StatusCode != 0 || info.Header != nil || calls.Load() != 0 {
		t.Fatalf("got = %v, err = %v, info = %+v, calls = %d", got, err, info, calls.Load())
	}
}

func TestServiceUserContextLocalErrors(t *testing.T) {
	documenttest.LocalErrors(t, func(ctx context.Context, client *transport.Client) (*vendorapi.UserContext, error) {
		s, err := vendorapi.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.UserContext(ctx, userToken)
	})
}
