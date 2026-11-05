package vendorapi

import (
	"cmp"
	"encoding/json/jsontext"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Secret — токен из тела запроса МоегоСклада. fmt, slog и json его не раскрывают;
// значение достаётся явно: string(s).
type Secret string

func (Secret) String() string { return "<redacted>" }

// Format перекрывает любой глагол fmt: ошибочный (%d) иначе печатает значение в %!d(...).
func (s Secret) Format(f fmt.State, _ rune)   { io.WriteString(f, s.String()) }
func (s Secret) LogValue() slog.Value         { return slog.StringValue(s.String()) }
func (s Secret) MarshalJSON() ([]byte, error) { return []byte(`"<redacted>"`), nil }

// ActivationPattern — путь от {endpointBase} для http.ServeMux: PUT (активация и смена подписки), DELETE
// (деактивация). Путь из endpointBase снимите http.StripPrefix; МойСклад повторяет запросы.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/vendor-endpoints#2-rest-endpointy-na-storone-razrabotchika-reshenij
const ActivationPattern = "/api/moysklad/vendor/1.0/apps/{appId}/{accountId}"

// RequestID — id запроса МоегоСклада для логов; при повторе запроса не меняется. Живой
// МойСклад (2026-11-05) шлёт X-Lognex-RequestId через дефис, документация пишет
// X_Lognex_RequestId; читаются оба, дефис первым. Пусто, если заголовка нет.
func RequestID(r *http.Request) string {
	return cmp.Or(r.Header.Get("X-Lognex-RequestId"), r.Header.Get("X_Lognex_RequestId"))
}

// Cause — причина активации, деактивации или события. Неизвестные значения сохраняются.
type Cause string

const (
	CauseInstall          Cause = "Install"
	CauseResume           Cause = "Resume"
	CauseTariffChanged    Cause = "TariffChanged"
	CauseAutoprolongation Cause = "Autoprolongation"
	CauseUninstall        Cause = "Uninstall"
	CauseSuspend          Cause = "Suspend"
	// CausePermissionsChanged приходит в PUT …/event; путь события в источнике противоречив.
	CausePermissionsChanged Cause = "PermissionsChanged"
)

// Status — статус решения на аккаунте. Решение может вернуть или выставить только
// Activating, SettingsRequired и Activated.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/moysklad-endpoints#4-tablica-vozmozhnyh-statusov-resheniya-na-akkaunte
type Status string

const (
	StatusActivating         Status = "Activating"
	StatusActivationFailed   Status = "ActivationFailed"
	StatusSettingsRequired   Status = "SettingsRequired"
	StatusActivated          Status = "Activated"
	StatusDeactivating       Status = "Deactivating"
	StatusDeactivationFailed Status = "DeactivationFailed"
	StatusSuspended          Status = "Suspended"
)

// ActivationRequest — тело PUT; имя — по разделу «Активация» источника, хотя TariffChanged и
// Autoprolongation по glossary не активация. Access при них и без доступа к API не приходит.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/vendor-endpoints#3-aktivaciya-resheniya-na-akkaunte
type ActivationRequest struct {
	AppUID       string        `json:"appUid"`
	AccountName  string        `json:"accountName"`
	Cause        Cause         `json:"cause"`
	Access       []Access      `json:"access,omitempty"`
	Subscription *Subscription `json:"subscription,omitempty"`
	Additional   *Additional   `json:"additional,omitempty"`
}

// Access — доступ к ресурсу из дескриптора. Permissions только при scope custom; источник
// называет его массивом, пример — объект, поэтому JSON хранится как есть.
type Access struct {
	Resource    string         `json:"resource"`
	Scope       []string       `json:"scope"`
	Permissions jsontext.Value `json:"permissions,omitempty"`
	AccessToken Secret         `json:"access_token"`
}

type Subscription struct {
	TariffID     string     `json:"tariffId"`
	Trial        bool       `json:"trial"`
	TariffName   string     `json:"tariffName,omitempty"`
	ExpiryMoment *time.Time `json:"expiryMoment,omitempty"`
	NotForResale bool       `json:"notForResale"`
	Partner      bool       `json:"partner"`
}

type Additional struct {
	FiscalAPI *FiscalAPI `json:"fiscalApi,omitempty"`
}

// FiscalAPI — регистрация установки; Token подписывает запросы МоегоСклада к решению по Fiscal API.
type FiscalAPI struct {
	ID    string `json:"id"`
	Token Secret `json:"token"`
}

// DeactivationRequest — тело DELETE; appUid и accountName есть только в примерах источника.
// Ответ — пустое тело: 200 или 204, если решение уже деактивировано.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/vendor-endpoints#3-deaktivaciya-resheniya-na-akkaunte
type DeactivationRequest struct {
	Cause       Cause  `json:"cause"`
	AppUID      string `json:"appUid,omitempty"`
	AccountName string `json:"accountName,omitempty"`
}

// ActivationResponse — тело ответа 200 на PUT; Content-Type: application/json обязателен.
type ActivationResponse struct {
	Status Status `json:"status"`
}
