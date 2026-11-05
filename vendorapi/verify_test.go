package vendorapi_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/vendorapi"
)

const (
	incomingKey    = "incoming-secret"
	incomingHeader = `{"alg":"HS256","typ":"JWT"}`
	// Payload из примера документации: iat=1516239022, exp=iat+300.
	incomingPayload = `{"iat":1516239022,"exp":1516239322,"jti":"6S3BQLsaSRNdEnhPCoW9lplY2LozRUOq"}`
)

var (
	issuedAt  = time.Unix(1516239022, 0)
	expiresAt = time.Unix(1516239322, 0)
	validNow  = issuedAt.Add(time.Minute)
)

// incomingToken собирает входящий JWT МоегоСклада из сырых JSON header и payload и подписывает
// его HMAC-SHA256 независимо от пакета.
func incomingToken(header, payload, key string) string {
	enc := base64.RawURLEncoding
	return signParts(enc.EncodeToString([]byte(header)), enc.EncodeToString([]byte(payload)), key)
}

// signParts подписывает уже закодированные части: так битая часть получает верную подпись.
func signParts(header, payload, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(header + "." + payload))
	return header + "." + payload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func TestVerifyToken(t *testing.T) {
	claims, err := vendorapi.VerifyToken(incomingToken(incomingHeader, incomingPayload, incomingKey), incomingKey, validNow)
	if err != nil {
		t.Fatal(err)
	}
	if !claims.IssuedAt.Equal(issuedAt) || !claims.ExpiresAt.Equal(expiresAt) || claims.ID != "6S3BQLsaSRNdEnhPCoW9lplY2LozRUOq" {
		t.Fatalf("claims = %+v", claims)
	}
}

// Дробные секунды NumericDate допустимы (RFC 7519 §2), лишние claims игнорируются.
func TestVerifyTokenFractionalAndExtraClaims(t *testing.T) {
	payload := `{"iat":1516239022.5,"exp":1516239322.25,"jti":"x","sub":"ignored"}`
	claims, err := vendorapi.VerifyToken(incomingToken(incomingHeader, payload, incomingKey), incomingKey, validNow)
	if err != nil {
		t.Fatal(err)
	}
	if !claims.IssuedAt.Equal(issuedAt.Add(500*time.Millisecond)) || !claims.ExpiresAt.Equal(expiresAt.Add(250*time.Millisecond)) {
		t.Fatalf("claims = %+v", claims)
	}
}

// Токен действителен строго до exp: now == exp уже истёк (RFC 7519 §4.1.4).
func TestVerifyTokenExpired(t *testing.T) {
	token := incomingToken(incomingHeader, incomingPayload, incomingKey)
	for _, now := range []time.Time{expiresAt, expiresAt.Add(time.Nanosecond), expiresAt.Add(time.Hour)} {
		claims, err := vendorapi.VerifyToken(token, incomingKey, now)
		if claims != nil || !errors.Is(err, vendorapi.ErrTokenExpired) || errors.Is(err, vendorapi.ErrInvalidToken) {
			t.Errorf("now = %v: claims = %v, err = %v", now, claims, err)
		}
	}
	if _, err := vendorapi.VerifyToken(token, incomingKey, expiresAt.Add(-time.Nanosecond)); err != nil {
		t.Errorf("before exp: %v", err)
	}
}

// Пустой ключ — ошибка конфигурации, а не неверный токен.
func TestVerifyTokenEmptyKey(t *testing.T) {
	claims, err := vendorapi.VerifyToken(incomingToken(incomingHeader, incomingPayload, ""), "", validNow)
	if claims != nil || err == nil || errors.Is(err, vendorapi.ErrInvalidToken) || errors.Is(err, vendorapi.ErrTokenExpired) {
		t.Fatalf("claims = %v, err = %v", claims, err)
	}
}

func TestVerifyTokenRejects(t *testing.T) {
	valid := incomingToken(incomingHeader, incomingPayload, incomingKey)
	parts := strings.Split(valid, ".")
	tampered := incomingToken(incomingHeader, strings.Replace(incomingPayload, "1516239322", "1916239322", 1), incomingKey)
	tamperedParts := strings.Split(tampered, ".")
	payload := func(p string) string { return incomingToken(incomingHeader, p, incomingKey) }
	// 32 байта подписи — 43 символа: младшие 2 бита последнего символа обязаны быть нулями.
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	last := strings.IndexByte(alphabet, parts[2][len(parts[2])-1])
	nonCanonical := parts[0] + "." + parts[1] + "." + parts[2][:len(parts[2])-1] + string(alphabet[last|1])
	header := func(h string) string { return incomingToken(h, incomingPayload, incomingKey) }

	for name, token := range map[string]string{
		"empty":            "",
		"two parts":        parts[0] + "." + parts[1],
		"four parts":       valid + "." + parts[2],
		"padded base64":    parts[0] + "=." + parts[1] + "." + parts[2],
		"broken header":    signParts("!!!", parts[1], incomingKey),
		"broken payload":   signParts(parts[0], "!!!", incomingKey),
		"broken signature": parts[0] + "." + parts[1] + ".!!!",
		"empty signature":  parts[0] + "." + parts[1] + ".",
		"non-canonical":    nonCanonical,
		// encoding/base64 пропускает \r и \n: без проверки алфавита подпись и части были бы неоднозначны.
		"CRLF signature":  parts[0] + "." + parts[1] + "." + parts[2][:10] + "\r\n" + parts[2][10:],
		"LF signature":    parts[0] + "." + parts[1] + "." + parts[2] + "\n",
		"CRLF payload":    signParts(parts[0], parts[1][:10]+"\r\n"+parts[1][10:], incomingKey),
		"LF header":       signParts(parts[0]+"\n", parts[1], incomingKey),
		"space":           signParts(parts[0], parts[1]+" ", incomingKey),
		"foreign key":     incomingToken(incomingHeader, incomingPayload, "other-secret"),
		"changed payload": parts[0] + "." + tamperedParts[1] + "." + parts[2],

		"alg none":       incomingToken(`{"alg":"none","typ":"JWT"}`, incomingPayload, incomingKey),
		"alg none empty": strings.Join([]string{base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`)), parts[1], ""}, "."),
		"alg HS384":      header(`{"alg":"HS384","typ":"JWT"}`),
		"alg HS512":      header(`{"alg":"HS512","typ":"JWT"}`),
		"alg RS256":      header(`{"alg":"RS256","typ":"JWT"}`),
		"alg lowercase":  header(`{"alg":"hs256","typ":"JWT"}`),
		"alg missing":    header(`{"typ":"JWT"}`),
		"alg null":       header(`{"alg":null,"typ":"JWT"}`),
		"typ null":       header(`{"alg":"HS256","typ":null}`),
		"typ other":      header(`{"alg":"HS256","typ":"at+jwt"}`),
		"typ number":     header(`{"alg":"HS256","typ":1}`),
		"header array":   header(`["HS256"]`),
		"header broken":  header(`{"alg":`),
		// crit требует понимать перечисленные расширения; SDK их не поддерживает (RFC 7515 §4.1.11).
		"crit": header(`{"alg":"HS256","typ":"JWT","crit":["exp"]}`),

		"iat missing":    payload(`{"exp":1516239322,"jti":"x"}`),
		"exp missing":    payload(`{"iat":1516239022,"jti":"x"}`),
		"jti missing":    payload(`{"iat":1516239022,"exp":1516239322}`),
		"iat null":       payload(`{"iat":null,"exp":1516239322,"jti":"x"}`),
		"exp null":       payload(`{"iat":1516239022,"exp":null,"jti":"x"}`),
		"jti null":       payload(`{"iat":1516239022,"exp":1516239322,"jti":null}`),
		"iat string":     payload(`{"iat":"1516239022","exp":1516239322,"jti":"x"}`),
		"exp string":     payload(`{"iat":1516239022,"exp":"123","jti":"x"}`),
		"exp bool":       payload(`{"iat":1516239022,"exp":true,"jti":"x"}`),
		"exp negative":   payload(`{"iat":1516239022,"exp":-1,"jti":"x"}`),
		"exp huge":       payload(`{"iat":1516239022,"exp":1e300,"jti":"x"}`),
		"jti empty":      payload(`{"iat":1516239022,"exp":1516239322,"jti":""}`),
		"jti number":     payload(`{"iat":1516239022,"exp":1516239322,"jti":42}`),
		"duplicate exp":  payload(`{"iat":1516239022,"exp":1,"exp":1516239322,"jti":"x"}`),
		"payload array":  payload(`[1516239022]`),
		"payload null":   payload(`null`),
		"payload broken": payload(`{"iat":`),
	} {
		t.Run(name, func(t *testing.T) {
			claims, err := vendorapi.VerifyToken(token, incomingKey, validNow)
			if claims != nil || !errors.Is(err, vendorapi.ErrInvalidToken) {
				t.Fatalf("claims = %v, err = %v", claims, err)
			}
			if strings.Contains(err.Error(), incomingKey) {
				t.Errorf("secret in error: %v", err)
			}
		})
	}
}

// typ необязателен (RFC 7519 §5.1) и без учёта регистра (RFC 7515 §4.1.9): legacy проверял
// токены МоегоСклада без typ, и живые токены строгую проверку не прошли.
func TestVerifyTokenTyp(t *testing.T) {
	for name, header := range map[string]string{
		"typ missing":   `{"alg":"HS256"}`,
		"typ lowercase": `{"alg":"HS256","typ":"jwt"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := vendorapi.VerifyToken(incomingToken(header, incomingPayload, incomingKey), incomingKey, validNow); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// Отказ по header называет alg и typ: без них отказ живого токена не разобрать по логам.
func TestVerifyTokenHeaderErrorNamesAlgTyp(t *testing.T) {
	_, err := vendorapi.VerifyToken(incomingToken(`{"alg":"HS512","typ":"at+jwt"}`, incomingPayload, incomingKey), incomingKey, validNow)
	if !errors.Is(err, vendorapi.ErrInvalidToken) || !strings.Contains(err.Error(), `alg="HS512"`) || !strings.Contains(err.Error(), `typ="at+jwt"`) {
		t.Fatalf("err = %v", err)
	}
}

func TestVerifyRequest(t *testing.T) {
	token := incomingToken(incomingHeader, incomingPayload, incomingKey)
	req, err := http.NewRequest(http.MethodPut, "https://example.com/api/moysklad/vendor/1.0/apps/a/b", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	claims, err := vendorapi.VerifyRequest(req, incomingKey, validNow)
	if err != nil || claims.ID != "6S3BQLsaSRNdEnhPCoW9lplY2LozRUOq" {
		t.Fatalf("claims = %+v, err = %v", claims, err)
	}
	// Схема авторизации нечувствительна к регистру (RFC 7235 §2.1).
	req.Header.Set("Authorization", "bearer "+token)
	if _, err := vendorapi.VerifyRequest(req, incomingKey, validNow); err != nil {
		t.Fatalf("lowercase scheme: %v", err)
	}

	for name, values := range map[string][]string{
		"missing":      nil,
		"empty":        {""},
		"basic":        {"Basic " + token},
		"no token":     {"Bearer "},
		"double space": {"Bearer  " + token},
		"two headers":  {"Bearer " + token, "Bearer " + token},
	} {
		t.Run(name, func(t *testing.T) {
			req.Header["Authorization"] = values
			claims, err := vendorapi.VerifyRequest(req, incomingKey, validNow)
			if claims != nil || !errors.Is(err, vendorapi.ErrInvalidToken) {
				t.Fatalf("claims = %v, err = %v", claims, err)
			}
		})
	}
}
