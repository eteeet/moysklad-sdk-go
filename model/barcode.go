package model

// Barcode — один штрихкод: ключ — формат в нижнем регистре (ean13, ean8, code128, gtin, upc …),
// значение — сам штрихкод. Набор форматов открыт.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service#5-shtrihkody
type Barcode map[string]string
