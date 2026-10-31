package positions

// PageOptions задаёт одну страницу позиций; Limit=0 означает 1000.
// Допустимы Limit от 1 до 1000 и неотрицательный Offset.
type PageOptions struct {
	Limit  int
	Offset int
}
