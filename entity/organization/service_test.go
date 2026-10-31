package organization_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const fixtureOrganizationID = "4b9d5bec-0575-11e6-9464-e4de00000008"

func TestServiceGet(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/organization.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, compressed := range []bool{false, true} {
		name := "plain"
		if compressed {
			name = "gzip"
		}
		t.Run(name, func(t *testing.T) {
			var calls atomic.Int32
			service := newOrganizationService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assertRequest(t, r, http.MethodGet, "/api/remap/1.2/entity/organization/"+fixtureOrganizationID, nil)
				writeOrganizationResponse(t, w, fixture, compressed)
			}))
			got, err := service.Get(t.Context(), fixtureOrganizationID)
			if err != nil {
				t.Fatal(err)
			}
			assertServerOrganization(t, got)
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want 1", calls.Load())
			}
		})
	}
}

func TestServiceCreate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/organization.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name       string
		input      organization.CreateRequest
		want       []byte
		compressed bool
	}{
		{"full/plain", fullCreateRequest(t), organizationGolden(t, "create"), false},
		{"full/gzip", fullCreateRequest(t), organizationGolden(t, "create"), true},
		{"minimal", organization.CreateRequest{Name: "SDK Organization"}, []byte(`{"name":"SDK Organization"}`), false},
		{"whitespace name", organization.CreateRequest{Name: " \t "}, []byte(`{"name":" \t "}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newOrganizationService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assertRequest(t, r, http.MethodPost, "/api/remap/1.2/entity/organization", tc.want)
				writeOrganizationResponse(t, w, fixture, tc.compressed)
			}))
			got, err := service.Create(t.Context(), tc.input)
			if err != nil {
				t.Fatal(err)
			}
			assertServerOrganization(t, got)
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one POST without preliminary GET", calls.Load())
			}
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/organization.json")
	if err != nil {
		t.Fatal(err)
	}
	const routeID = "11111111-1111-1111-1111-111111111111"
	for _, tc := range []struct {
		name       string
		input      organization.UpdateRequest
		want       []byte
		compressed bool
	}{
		{"full/plain", fullUpdateRequest(t), organizationGolden(t, "update"), false},
		{"full/gzip", fullUpdateRequest(t), organizationGolden(t, "update"), true},
		{"one field", organization.UpdateRequest{Name: model.Value("renamed")}, []byte(`{"name":"renamed"}`), false},
		{"sync ID", organization.UpdateRequest{SyncID: model.Value("11111111-1111-1111-1111-111111111111")}, []byte(`{"syncId":"11111111-1111-1111-1111-111111111111"}`), false},
		{"empty sync ID", organization.UpdateRequest{SyncID: model.Value("")}, []byte(`{"syncId":""}`), false},
		{"null sync ID", organization.UpdateRequest{SyncID: model.Null[string]()}, []byte(`{"syncId":null}`), false},
		{"empty name", organization.UpdateRequest{Name: model.Value("")}, []byte(`{"name":""}`), false},
		{"payerVat false", organization.UpdateRequest{PayerVAT: model.Value(false)}, []byte(`{"payerVat":false}`), false},
		{"empty attributes", organization.UpdateRequest{Attributes: model.Value([]model.Attribute{})}, []byte(`{"attributes":[]}`), false},
		{"null", organization.UpdateRequest{BonusProgram: model.Null[model.Reference]()}, []byte(`{"bonusProgram":null}`), false},
		{"empty structured address", organization.UpdateRequest{ActualAddressFull: model.Value(model.AddressRequest{})}, []byte(`{"actualAddressFull":{}}`), false},
		{"partial structured address", organization.UpdateRequest{LegalAddressFull: model.Value(model.AddressRequest{City: model.Value(""), Region: model.Null[model.Reference]()})}, []byte(`{"legalAddressFull":{"city":"","region":null}}`), false},
		{"unknown company type", organization.UpdateRequest{CompanyType: model.Value(model.CompanyType("futureKind"))}, []byte(`{"companyType":"futureKind"}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newOrganizationService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assertRequest(t, r, http.MethodPut, "/api/remap/1.2/entity/organization/"+routeID, tc.want)
				writeOrganizationResponse(t, w, fixture, tc.compressed)
			}))
			got, err := service.Update(t.Context(), routeID, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			assertServerOrganization(t, got)
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one PUT without preliminary GET", calls.Load())
			}
		})
	}
}

func TestServiceResponseErrors(t *testing.T) {
	for _, operation := range serviceOperations() {
		t.Run(operation.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				status  int
				body    string
				wantAPI bool
			}{
				{"bad request", 400, `{"errors":[{"error":"invalid requisites","code":3000}]}`, true},
				{"forbidden", 403, `{"errors":[{"error":"forbidden","code":1000}]}`, true},
				{"not found", 404, `{"errors":[{"error":"not found","code":1001}]}`, true},
				{"rate limit", 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
				{"server error", 500, `{"errors":[{"error":"server failed","code":5000}]}`, true},
				{"HTML error", 502, "<html>bad gateway</html>", true},
				{"unexpected status", 201, `{"id":"unexpected"}`, false},
				{"no content", 204, "", false},
				{"missing id", 200, `{"name":"invalid"}`, false},
				{"empty id", 200, `{"id":""}`, false},
				{"null", 200, "null", false},
				{"malformed JSON", 200, `{"id":"partial",`, false},
				{"invalid date", 200, `{"id":"partial","created":"invalid"}`, false},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var calls atomic.Int32
					service := newOrganizationService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						calls.Add(1)
						w.WriteHeader(tc.status)
						if _, err := io.WriteString(w, tc.body); err != nil {
							t.Error(err)
						}
					}))
					got, err := operation.call(service, t.Context())
					if got != nil || err == nil {
						t.Fatalf("organization = %+v, error = %v", got, err)
					}
					assertErrorContext(t, err, operation.context)
					var apiErr *sdkerrors.APIError
					if errors.As(err, &apiErr) != tc.wantAPI {
						t.Fatalf("APIError = %v, error = %v", apiErr, err)
					}
					if tc.wantAPI && apiErr.StatusCode != tc.status {
						t.Errorf("status = %d, want %d", apiErr.StatusCode, tc.status)
					}
					if tc.status == 400 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 3000 || apiErr.Errors[0].Message != "invalid requisites") {
						t.Errorf("API details = %+v", apiErr.Errors)
					}
					if calls.Load() != 1 {
						t.Fatalf("calls = %d, want 1 without retry", calls.Load())
					}
				})
			}
		})
	}
}

func TestServiceErrorWrapping(t *testing.T) {
	for _, operation := range serviceOperations() {
		t.Run(operation.name, func(t *testing.T) {
			for _, authError := range []bool{false, true} {
				name := "transport"
				if authError {
					name = "authorization"
				}
				t.Run(name, func(t *testing.T) {
					wantErr := errors.New("boundary failed")
					var calls atomic.Int32
					cfg := transport.Config{HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
						calls.Add(1)
						return nil, wantErr
					})}}
					if authError {
						cfg.Auth = func(*http.Request) error { return wantErr }
					}
					client, err := transport.NewClient(cfg)
					if err != nil {
						t.Fatal(err)
					}
					service, err := organization.NewService(client)
					if err != nil {
						t.Fatal(err)
					}
					got, err := operation.call(service, t.Context())
					if got != nil || !errors.Is(err, wantErr) {
						t.Fatalf("organization = %+v, error = %v, want wrapped boundary error", got, err)
					}
					assertErrorContext(t, err, operation.context)
					wantCalls := int32(1)
					if authError {
						wantCalls = 0
					}
					if calls.Load() != wantCalls {
						t.Errorf("calls = %d, want %d", calls.Load(), wantCalls)
					}
				})
			}
		})
	}
}

func TestServiceContext(t *testing.T) {
	for _, operation := range serviceOperations() {
		t.Run(operation.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newOrganizationService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.WriteHeader(http.StatusInternalServerError)
			}))
			if got, err := operation.call(service, nil); err == nil || got != nil {
				t.Fatalf("nil context: organization = %+v, error = %v", got, err)
			} else {
				assertErrorContext(t, err, operation.context)
			}
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			if got, err := operation.call(service, ctx); got != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("canceled context: organization = %+v, error = %v", got, err)
			} else {
				assertErrorContext(t, err, operation.context)
			}
			if calls.Load() != 0 {
				t.Fatalf("calls = %d, want 0", calls.Load())
			}
		})
	}
}

func TestServiceValidation(t *testing.T) {
	if service, err := organization.NewService(nil); service != nil || err == nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	service := newOrganizationService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	for _, id := range []string{"", "metadata", ".", "..", "../metadata", "x/y", fixtureOrganizationID + "/metadata", "4b9d5bec-0575-11e6-9464-e4de0000000g", "4b9d5bec057511e69464e4de00000008", "4b9d5bec_0575-11e6-9464-e4de00000008"} {
		if got, err := service.Get(t.Context(), id); got != nil || err == nil {
			t.Errorf("Get accepted invalid id %q", id)
		}
		if got, err := service.Update(t.Context(), id, organization.UpdateRequest{PayerVAT: model.Value(false)}); got != nil || err == nil {
			t.Errorf("Update accepted invalid id %q", id)
		}
	}
	if got, err := service.Create(t.Context(), organization.CreateRequest{}); got != nil || err == nil {
		t.Fatal("Create accepted empty name")
	}
	if got, err := service.Update(t.Context(), fixtureOrganizationID, organization.UpdateRequest{}); got != nil || err == nil {
		t.Fatal("Update accepted unset request")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d, want 0", calls.Load())
	}
}

func TestServiceUUIDVersion(t *testing.T) {
	for _, id := range []string{"00000000-0000-f000-0000-000000000000", "4B9D5BEC-0575-11E6-9464-E4DE00000008"} {
		t.Run(id, func(t *testing.T) {
			var calls atomic.Int32
			service := newOrganizationService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.URL.Path != "/api/remap/1.2/entity/organization/"+id {
					t.Errorf("path = %s", r.URL.Path)
				}
				_, _ = io.WriteString(w, `{"id":"server-id"}`)
			}))
			if got, err := service.Get(t.Context(), id); err != nil || got == nil || got.ID != "server-id" {
				t.Fatalf("Get: organization = %+v, error = %v", got, err)
			}
			if got, err := service.Update(t.Context(), id, organization.UpdateRequest{PayerVAT: model.Value(false)}); err != nil || got == nil || got.ID != "server-id" {
				t.Fatalf("Update: organization = %+v, error = %v", got, err)
			}
			if calls.Load() != 2 {
				t.Fatalf("calls = %d, want 2", calls.Load())
			}
		})
	}
}

func serviceOperations() []struct {
	name    string
	context string
	call    func(*organization.Service, context.Context) (*organization.Organization, error)
} {
	return []struct {
		name    string
		context string
		call    func(*organization.Service, context.Context) (*organization.Organization, error)
	}{
		{"get", "get organization " + fixtureOrganizationID, func(s *organization.Service, ctx context.Context) (*organization.Organization, error) {
			return s.Get(ctx, fixtureOrganizationID)
		}},
		{"create", "create organization", func(s *organization.Service, ctx context.Context) (*organization.Organization, error) {
			return s.Create(ctx, organization.CreateRequest{Name: "SDK Organization"})
		}},
		{"update", "update organization " + fixtureOrganizationID, func(s *organization.Service, ctx context.Context) (*organization.Organization, error) {
			return s.Update(ctx, fixtureOrganizationID, organization.UpdateRequest{PayerVAT: model.Value(false)})
		}},
		{"update sync ID", "update organization " + fixtureOrganizationID, func(s *organization.Service, ctx context.Context) (*organization.Organization, error) {
			return s.Update(ctx, fixtureOrganizationID, organization.UpdateRequest{SyncID: model.Value("11111111-1111-1111-1111-111111111111")})
		}},
	}
}

func assertErrorContext(t *testing.T, err error, operation string) {
	t.Helper()
	if strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), operation) {
		t.Errorf("error context = %v, want %q and one moysklad:", err, operation)
	}
}

func assertServerOrganization(t *testing.T, got *organization.Organization) {
	t.Helper()
	if got == nil || got.ID != fixtureOrganizationID || got.Meta.Type != "organization" || got.Name != "ОАО СветПром" || got.Code != "666" ||
		got.AccountID != "84e60e93-f504-11e5-8a84-bae500000008" || got.CompanyType != "legal" || got.Created.IsZero() || got.Updated.IsZero() ||
		got.Owner == nil || got.Group == nil || got.BonusProgram == nil || got.ActualAddressFull == nil || got.LegalAddressFull == nil ||
		!got.PayerVAT || got.Director != "Администратор" || got.ChiefAccountant != "Администратор" || got.TrackingContractNumber != "12345678" ||
		got.TrackingContractDate.IsZero() || len(got.Attributes) != 1 {
		t.Fatalf("organization = %+v", got)
	}
}

func assertRequest(t *testing.T, r *http.Request, method, path string, wantBody []byte) {
	t.Helper()
	if r.Method != method || r.URL.Path != path || r.URL.RawQuery != "" {
		t.Errorf("request = %s %s, want %s %s without query", r.Method, r.URL, method, path)
	}
	if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept") != "application/json;charset=utf-8" ||
		r.Header.Get("Accept-Encoding") != "gzip" || r.Header.Get("User-Agent") != "moysklad-sdk-go" {
		t.Error("required request headers missing")
	}
	contentType := "application/json"
	if wantBody == nil {
		contentType = ""
	}
	if r.Header.Get("Content-Type") != contentType {
		t.Errorf("Content-Type = %q, want %q", r.Header.Get("Content-Type"), contentType)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		return
	}
	if wantBody == nil {
		if len(body) != 0 {
			t.Errorf("GET body = %q", body)
		}
		return
	}
	assertJSONEqual(t, body, wantBody)
}

func writeOrganizationResponse(t *testing.T, w http.ResponseWriter, body []byte, compressed bool) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if compressed {
		w.Header().Set("Content-Encoding", "gzip")
		writer := gzip.NewWriter(w)
		if _, err := writer.Write(body); err != nil {
			t.Error(err)
		}
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
		return
	}
	if _, err := w.Write(body); err != nil {
		t.Error(err)
	}
}

func newOrganizationService(t *testing.T, handler http.Handler) *organization.Service {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Timeout: 5 * time.Second, Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	service, err := organization.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
