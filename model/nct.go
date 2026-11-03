package model

// NctKZ — данные Национального каталога товаров (поле mod__nct__kz), только для аккаунтов
// Казахстана; регион сервер проверяет сам. NTIN — идентификационный код товара.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-modul-dlya-kazahstana-s-informaciej-iz-nacionalnogo-kataloga-tovarov
type NctKZ struct {
	NTIN string `json:"ntin,omitempty"`
}
