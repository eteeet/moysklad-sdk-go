// Package positions содержит общие JSON-контракты позиций документов.
//
// Для embedded Create/Update используй model.Field[[]T] с json:",omitzero".
// Unset сохраняет позиции при Update; Value с nil или пустым slice отправляет [].
//
// Массив задаёт full replacement: matching ID обновляет строку, новая идёт без ID,
// отсутствующие старые строки удаляются. Запись выполняет документ без GET или merge.
//
// API ограничивает embedded массив 1000 строками; лимит проверяет writer документа.
// JSON-типы не выполняют business validation; Null не обещает допустимость очистки.
package positions

// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1037
