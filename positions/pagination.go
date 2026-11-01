package positions

import "github.com/eteeet/moysklad-sdk-go/query"

// PageOptions задаёт одну страницу позиций; Limit=0 означает 1000.
// Допустимы Limit от 1 до 1000 и неотрицательный Offset.
type PageOptions = query.PageOptions
