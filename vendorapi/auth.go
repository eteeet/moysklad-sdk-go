package vendorapi

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"time"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Header исходящего токена неизменен; кодируется один раз.
var outgoingHeader = base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

// AppAuth подписывает каждую попытку запроса новым одноразовым JWT решения: sub=appUID, iat,
// случайный jti. exp не передаётся: МойСклад ограничивает срок iat+300 с.
// secretKey служит только Vendor API; свои токены приложения подписывайте другим ключом.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/authentication#3-vhodyashie-zaprosy-razrabotchik-mojsklad
func AppAuth(appUID, secretKey string) (transport.Auth, error) {
	if appUID == "" || secretKey == "" {
		return nil, fmt.Errorf("moysklad: vendor app UID and secret key are required")
	}
	key := []byte(secretKey)
	return func(req *http.Request) error {
		var nonce [32]byte
		if _, err := rand.Read(nonce[:]); err != nil {
			return fmt.Errorf("moysklad: generate vendor token id: %w", err)
		}
		payload, err := json.Marshal(struct {
			Sub string `json:"sub"`
			Iat int64  `json:"iat"`
			Jti string `json:"jti"`
		}{appUID, time.Now().Unix(), hex.EncodeToString(nonce[:])})
		if err != nil {
			return fmt.Errorf("moysklad: encode vendor token: %w", err)
		}
		signed := outgoingHeader + "." + base64.RawURLEncoding.EncodeToString(payload)
		mac := hmac.New(sha256.New, key)
		mac.Write([]byte(signed))
		req.Header.Set("Authorization", "Bearer "+signed+"."+base64.RawURLEncoding.EncodeToString(mac.Sum(nil)))
		return nil
	}, nil
}
