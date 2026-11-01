package model

// MarkingUZ — модуль маркировки (поле mod__marking__uz), только для аккаунтов Узбекистана:
// читается и пишется с заголовком X-Lognex-Accept-Region: uz (transport.WithAcceptRegion).
// TrackingType открыт: ALCOHOL, BEER, HOME_APPLIANCE, NOT_TRACKED, TOBACCO, WATER …
type MarkingUZ struct {
	TrackingType string `json:"trackingType,omitempty"`
}
