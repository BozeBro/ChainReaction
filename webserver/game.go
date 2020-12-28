package webserver


type Game interface {
	InitBoard(int, int)
	MovePiece(int, int, string) ([][][]int, [][][]int)
	UpdateColor(string, string) bool
	GetRows() int
	GetCols() int
}
func (c *Chain) IsLegalMove(x, y int, color string) bool {
	validRow := 0 <= x && x < c.Squares[0].Len
	validCol := 0 <= y && y < c.Len
	validColor := color == c.Squares[y].Color[x] || c.Squares[y].Color[x] == ""
	return validRow && validCol && validColor
}
func IsLegalMove(g Game, x, y int) bool {
	validRow := 0 <= x && x < g.GetRows()
	validCol := 0 <= y && y < g.GetCols()
	return validRow && validCol
}
