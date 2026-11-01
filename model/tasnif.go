package model

// TasnifUZ — данные справочника ТАСНИФ (поле mod__tasnif__uz), только для аккаунтов Узбекистана:
// читается и пишется с заголовком X-Lognex-Accept-Region: uz (transport.WithAcceptRegion).
// Объект передаётся целиком: поля с пустой строкой не отправляются, очистить весь объект — null.
type TasnifUZ struct {
	IKPU          string `json:"ikpu,omitempty"`
	PackCode      string `json:"packCode,omitempty"`
	BarcodeTasnif string `json:"barcodeTasnif,omitempty"`
}
