package model

import "encoding/json/v2"

// Типы неснижаемого остатка; набор открыт.
const (
	MinimumStockAllWarehouseSum  = "ALL_WAREHOUSE_SUM"
	MinimumStockAllWarehouseSame = "ALL_WAREHOUSE_SAME"
	MinimumStockWarehouseVaried  = "WAREHOUSE_VARIED"
)

// MinimumStock — неснижаемый остаток товара или модификации; выводится с fields=minimumStock
// (в SDK — только List: Get параметра fields не принимает) и в ответе Create/Update, где поле
// передано. Quantity — у ALL_WAREHOUSE_SUM и ALL_WAREHOUSE_SAME, StoreBalances — у
// WAREHOUSE_VARIED: без expand только метаданные коллекции, строки читает ListStoreBalances.
// Источник пишет коллекцию то storebalances (таблица, большинство примеров), то storeBalances
// (часть примеров, фикстура OpenAPI): читаются оба ключа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-nesnizhaemyj-ostatok
type MinimumStock struct {
	Type          string                    `json:"type"`
	Quantity      float64                   `json:"quantity"`
	StoreBalances *Collection[StoreBalance] `json:"storebalances,omitzero"`
}

func (m *MinimumStock) UnmarshalJSON(data []byte) error {
	type plain MinimumStock
	var v struct {
		plain
		Alt *Collection[StoreBalance] `json:"storeBalances,omitzero"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*m = MinimumStock(v.plain)
	if m.StoreBalances == nil {
		m.StoreBalances = v.Alt
	}
	return nil
}

// StoreBalance — неснижаемый остаток на одном складе (meta.type minimumstock).
type StoreBalance struct {
	Meta      Meta       `json:"meta"`
	ID        string     `json:"id"`
	AccountID string     `json:"accountId"`
	Store     *Reference `json:"store,omitzero"`
	Quantity  float64    `json:"quantity"`
}

// MinimumStockRequest — неснижаемый остаток в Create/Update товара или модификации. Тип
// ALL_WAREHOUSE_SAME только читается (проверяет сервер). StoreBalances — для WAREHOUSE_VARIED:
// элемент с Meta изменяет существующую позицию, без Meta — создаёт.
type MinimumStockRequest struct {
	Type          string                       `json:"type"`
	Quantity      Field[float64]               `json:"quantity,omitzero"`
	StoreBalances Field[[]StoreBalanceRequest] `json:"storebalances,omitzero"`
}

// StoreBalanceRequest — неснижаемый остаток склада: store и quantity. Meta таблица помечает «только
// для чтения», но пример изменения товара определяет ею существующую строку storebalances — SDK её
// передаёт; в CreateStoreBalance и UpdateStoreBalance Meta не нужна.
type StoreBalanceRequest struct {
	Meta     Field[Meta]      `json:"meta,omitzero"`
	Store    Field[Reference] `json:"store,omitzero"`
	Quantity Field[float64]   `json:"quantity,omitzero"`
}
