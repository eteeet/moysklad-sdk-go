package model

// Набор значений открыт; допустимость распределения определяет сервер.
type OverheadDistribution string

const (
	OverheadDistributionWeight OverheadDistribution = "weight"
	OverheadDistributionVolume OverheadDistribution = "volume"
	OverheadDistributionPrice  OverheadDistribution = "price"
)

// Распределение выполняет сервер; SDK не пересчитывает и не округляет суммы.
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L47
type DocumentOverhead struct {
	Sum          float64              `json:"sum"` // В копейках с ограничениями двоичной точности float64.
	Distribution OverheadDistribution `json:"distribution"`
}

// Незаданные поля не отправляются; допустимость null и распределение проверяет сервер.
type DocumentOverheadRequest struct {
	Sum          Field[float64]              `json:"sum,omitzero"`
	Distribution Field[OverheadDistribution] `json:"distribution,omitzero"`
}
