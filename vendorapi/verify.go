package vendorapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json/v2"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

var (
	// ErrInvalidToken — подпись, формат, алгоритм или claims входящего JWT неверны.
	ErrInvalidToken = errors.New("moysklad: invalid vendor token")
	// ErrTokenExpired — подпись верна, но now >= exp.
	ErrTokenExpired = errors.New("moysklad: vendor token expired")
)

// maxNumericDate — 9999-12-31T23:59:59Z: больше значение не переводится в time.Time без потерь.
const maxNumericDate = 253402300799

const base64URLAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var b64 = base64.RawURLEncoding.Strict()

// TokenClaims — проверенные claims входящего JWT. Повтор ID (jti) отклоняет решение: SDK состояния не хранит.
type TokenClaims struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
	ID        string
}

// VerifyToken проверяет JWT запроса МоегоСклада к решению; iat с now не сравнивается, допуска часов нет.
// Тот же secretKey не должен подписывать токены самого решения, иначе они пройдут проверку.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/authentication#3-ishodyashie-zaprosy-mojsklad-razrabotchik
func VerifyToken(token, secretKey string, now time.Time) (*TokenClaims, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("moysklad: vendor secret key is required")
	}
	// encoding/base64 пропускает \r и \n: без проверки алфавита одна подпись имела бы много записей.
	if strings.ContainsFunc(token, func(r rune) bool { return r != '.' && !strings.ContainsRune(base64URLAlphabet, r) }) {
		return nil, fmt.Errorf("%w: characters outside base64url", ErrInvalidToken)
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: want 3 parts", ErrInvalidToken)
	}

	var header map[string]any
	if err := decodePart(parts[0], &header); err != nil {
		return nil, fmt.Errorf("%w: header: %v", ErrInvalidToken, err)
	}
	// crit требует понимать расширения, которых SDK не знает (RFC 7515 §4.1.11). typ
	// необязателен (RFC 7519 §5.1) и без учёта регистра (RFC 7515 §4.1.9).
	typ, hasTyp := header["typ"]
	typString, typIsString := typ.(string)
	if _, crit := header["crit"]; crit || header["alg"] != "HS256" || hasTyp && !(typIsString && strings.EqualFold(typString, "JWT")) {
		return nil, fmt.Errorf("%w: header alg=%.20q typ=%.20q: want alg HS256, typ JWT or none, no crit",
			ErrInvalidToken, fmt.Sprint(header["alg"]), fmt.Sprint(typ))
	}

	signature, err := b64.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("%w: signature encoding", ErrInvalidToken)
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return nil, fmt.Errorf("%w: signature mismatch", ErrInvalidToken)
	}

	var payload map[string]any
	if err := decodePart(parts[1], &payload); err != nil {
		return nil, fmt.Errorf("%w: payload: %v", ErrInvalidToken, err)
	}
	iat, okIat := numericDate(payload["iat"])
	exp, okExp := numericDate(payload["exp"])
	jti, okJti := payload["jti"].(string)
	if !okIat || !okExp || !okJti || jti == "" {
		return nil, fmt.Errorf("%w: iat and exp must be numbers, jti a non-empty string", ErrInvalidToken)
	}
	if !now.Before(exp) {
		return nil, ErrTokenExpired
	}
	return &TokenClaims{IssuedAt: iat, ExpiresAt: exp, ID: jti}, nil
}

// VerifyRequest проверяет токен единственного заголовка Authorization: Bearer.
func VerifyRequest(r *http.Request, secretKey string, now time.Time) (*TokenClaims, error) {
	values := r.Header.Values("Authorization")
	if len(values) != 1 {
		return nil, fmt.Errorf("%w: want one Authorization header", ErrInvalidToken)
	}
	scheme, token, ok := strings.Cut(values[0], " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.ContainsAny(token, " \t") {
		return nil, fmt.Errorf("%w: want Bearer token", ErrInvalidToken)
	}
	return VerifyToken(token, secretKey, now)
}

// decodePart разбирает base64url-часть JWT в JSON-объект; json/v2 отклоняет повтор имён.
func decodePart(part string, dst *map[string]any) error {
	data, err := b64.DecodeString(part)
	if err != nil {
		return errors.New("encoding")
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return errors.New("not a JSON object")
	}
	if *dst == nil {
		return errors.New("not a JSON object")
	}
	return nil
}

// numericDate переводит NumericDate (секунды, дробная часть допустима) в time.Time.
func numericDate(v any) (time.Time, bool) {
	f, ok := v.(float64)
	if !ok || f < 0 || f > maxNumericDate {
		return time.Time{}, false
	}
	sec, frac := math.Modf(f)
	return time.Unix(int64(sec), int64(math.Round(frac*1e9))), true
}
