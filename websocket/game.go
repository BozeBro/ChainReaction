package websocket

// Game interface that provides functions that give general flow of a board game
type Game interface {
	InitBoard(int, int)
	MovePiece(int, int, string) ([][][]int, [][][]int)
	UpdateColor(string, string) bool
	GetRows() int
	GetCols() int
	GetBoard() []*Squares
	IsLegalMove(int, int, string) bool
	Minimax(string, bool, []string, int, int, int, int, int) (int, [2]int)
}

// IsLegalMove Tells if a move is allowed
// Must be same color square or empty
// Cannot be off the board
func (c *Chain) IsLegalMove(x, y int, color string) bool {
	validRow := 0 <= x && x < c.Squares[0].Len
	validCol := 0 <= y && y < c.Len
	validColor := color == c.Squares[y].Color[x] || c.Squares[y].Color[x] == ""
	return validRow && validCol && validColor
}
