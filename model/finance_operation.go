package model

// Связь платежа с операцией без expanded-полей; SDK не пересчитывает суммы.
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L507
type FinanceOperation struct {
	Meta      Meta    `json:"meta"`
	LinkedSum float64 `json:"linkedSum"` // Оплаченная сумма в копейках с ограничениями двоичной точности float64.
}

// LinkedSum можно опустить; Null кодируется без гарантии принятия сервером.
// Внешний operations:[] передаётся точно, но не обещает очистку или merge на сервере.
type FinanceOperationRequest struct {
	Meta      Meta           `json:"meta"`
	LinkedSum Field[float64] `json:"linkedSum,omitzero"`
}
