package vendorapi_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/vendorapi"
)

var hexJTI = regexp.MustCompile(`^[0-9a-f]{64}$`)

// outgoingClaims разбирает исходящий JWT независимо от пакета: проверяет HMAC-SHA256 и
// возвращает header и payload как JSON-объекты.
func outgoingClaims(t *testing.T, authorization, secretKey string) (header, payload map[string]any) {
	t.Helper()
	token, ok := strings.CutPrefix(authorization, "Bearer ")
	if !ok {
		t.Fatalf("Authorization = %q", authorization)
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("token parts = %d", len(parts))
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	if got := base64.RawURLEncoding.EncodeToString(mac.Sum(nil)); got != parts[2] {
		t.Fatalf("signature = %s, want %s", parts[2], got)
	}
	for i, dst := range []*map[string]any{&header, &payload} {
		data, err := base64.RawURLEncoding.DecodeString(parts[i])
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, dst); err != nil {
			t.Fatal(err)
		}
	}
	return header, payload
}

func authorize(t *testing.T, auth func(*http.Request) error) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "https://apps-api.moysklad.ru/api/vendor/1.0/context/x", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := auth(req); err != nil {
		t.Fatal(err)
	}
	return req.Header.Get("Authorization")
}

func TestAppAuthToken(t *testing.T) {
	auth, err := vendorapi.AppAuth("dummyapp.example-vendor", "secret-key")
	if err != nil {
		t.Fatal(err)
	}
	before := time.Now().Unix()
	header, payload := outgoingClaims(t, authorize(t, auth), "secret-key")
	after := time.Now().Unix()

	if len(header) != 2 || header["alg"] != "HS256" || header["typ"] != "JWT" {
		t.Errorf("header = %v", header)
	}
	iat, _ := payload["iat"].(float64)
	jti, _ := payload["jti"].(string)
	if len(payload) != 3 || payload["sub"] != "dummyapp.example-vendor" || int64(iat) < before || int64(iat) > after || iat != float64(int64(iat)) || !hexJTI.MatchString(jti) {
		t.Errorf("payload = %v", payload)
	}
	if _, ok := payload["exp"]; ok {
		t.Error("exp must be absent")
	}
}

// Каждый вызов Auth выпускает новый одноразовый токен.
func TestAppAuthFreshJTI(t *testing.T) {
	auth, err := vendorapi.AppAuth("app", "secret-key")
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for range 3 {
		_, payload := outgoingClaims(t, authorize(t, auth), "secret-key")
		jti := payload["jti"].(string)
		if seen[jti] {
			t.Fatalf("jti %s repeated", jti)
		}
		seen[jti] = true
	}
}

func TestAppAuthRejectsEmpty(t *testing.T) {
	for _, tc := range [][2]string{{"", "secret"}, {"app", ""}, {"", ""}} {
		if auth, err := vendorapi.AppAuth(tc[0], tc[1]); err == nil || auth != nil {
			t.Errorf("AppAuth(%q, %q) = %v, %v", tc[0], tc[1], auth != nil, err)
		}
	}
}
