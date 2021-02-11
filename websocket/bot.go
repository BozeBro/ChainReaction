package websocket

import (
	"math"
	"math/rand"
)

func (c *Chain) RandMove(playerColor string) (int, int) {
	validSquares := make([][2]int, 0)
	for y, v := range c.Squares {
		for x, color := range v.Color {
			if color == "" || color == playerColor {
				validSquares = append(validSquares, [2]int{x, y})
			}
		}
	}
	sq := validSquares[rand.Intn(len(validSquares))]
	return sq[0], sq[1]
}

// Heuristic from https://brilliant.org/wiki/chain-reaction-game/
func BrillHeuristic(board []*Squares, color string) int {
	score := 0
	myCircles, enemyCircles := 0, 0
	for y := 0; y < len(board); y++ {
		for x := 0; x < board[0].Len; x++ {
			if board[y].Color[x] == color {
				myCircles++
				safe := true
				total, coords := findneighbors(x, y, board[0].Len, len(board))
				for _, coord := range coords {
					nx, ny := coord[0], coord[1]
					if board[ny].Color[nx] != color && board[ny].Cur[nx] == board[ny].Max[nx]-1 {
						adjTotal, _ := findneighbors(nx, ny, board[0].Len, len(board))
						score -= 5 - adjTotal
						safe = false
					}
				}
				if safe {
					if total == 3 {
						score += 2
					} else if total == 2 {
						score += 3
					}
					if board[y].Cur[x] == total-1 {
						score += 2
					}
				}
			} else {
				enemyCircles++
			}
		}
	}
	score += myCircles
	// A player can only win when they have more than one circle
	if enemyCircles == 0 && myCircles > 1 {
		return 10000
	} else if myCircles == 1 && enemyCircles > 0 {
		return -10000
	}
	for _, length := range findChains(board, color) {
		if length > 1 {
			score += 2 * length
		}
	}
	return score
}

func findChains(oldBoard []*Squares, color string) []int {
	board := oldBoard
	lengths := make([]int, 0)
	for y := 0; y < len(board); y++ {
		for x := 0; x < board[0].Len; x++ {
			if board[y].Cur[x] == board[y].Max[x]-1 && board[y].Color[x] == color {
				amount := 0
				visiting := [][]int{{x, y}}
				for len(visiting) > 0 {
					last := len(visiting) - 1
					x, y := visiting[last][0], visiting[last][1]
					visiting = visiting[:last]
					board[y].Cur[x] = 0
					total, coords := findneighbors(x, y, board[0].Len, len(board))
					amount++
					for _, coord := range coords {
						x, y := coord[0], coord[1]
						if board[y].Cur[x] == total-1 && board[y].Color[x] == color {
							visiting = append(visiting, coord)
						}
					}
				}
				lengths = append(lengths, amount)
			}
		}
	}
	return lengths
}

// Maximizing player
// Return greatest number possible
func (c *Chain) Max(color string, nextColor string, depth, alpha, beta, movedx, movedy int) (int, [2]int) {
	sq := [2]int{movedx, movedy}
	val := int(math.Inf(-1))
	if depth == 0 {
		// use nextColor because we are evaluating previous player's move
		boardValue := BrillHeuristic(c.Squares, nextColor)
		return boardValue, sq
	}
	newBoard := copyBoard(c.Squares)
	newClients := copyClients(c.Hub.Clients)
	players := make([]string, len(c.Hub.Colors), len(c.Hub.Colors))
	for ind, player := range c.Hub.Colors {
		players[ind] = player
	}
	for y := 0; y < c.Len; y++ {
		for x := 0; x < c.Squares[0].Len; x++ {
			if c.Squares[y].Color[x] == "" || c.Squares[y].Color[x] == color {
				c.MovePiece(x, y, color)
				if iswinner(c.Squares, color) {
					c.Hub.Clients = newClients
					c.Hub.Colors = players
					return 10000, [2]int{movedx, movedy}
				}
				maxVal, _ := c.Min(nextColor, color, depth-1, alpha, beta, x, y)
				if maxVal > val {
					sq = [2]int{x, y}
					val = maxVal
				}
				if val > alpha {
					alpha = val
				}
				c.Hub.Clients = newClients
				replaceBoard(c.Squares, newBoard)
				c.Hub.Colors = players
				if alpha >= beta {
					return val, sq
				}
			}
		}
	}
	return val, sq
}

// Minimizing player
// Return smallest number possible
// Look at Max() for more documentation
func (c *Chain) Min(color string, nextColor string, depth, alpha, beta, movedx, movedy int) (int, [2]int) {
	sq := [2]int{movedx, movedy}
	val := int(math.Inf(1))
	if depth == 0 {
		boardValue := BrillHeuristic(c.Squares, nextColor)
		return boardValue, sq
	}
	players := make([]string, len(c.Hub.Colors), len(c.Hub.Colors))
	for ind, player := range c.Hub.Colors {
		players[ind] = player
	}
	newBoard := copyBoard(c.Squares)
	newClients := copyClients(c.Hub.Clients)
	for y := 0; y < c.Len; y++ {
		for x := 0; x < c.Squares[0].Len; x++ {
			if c.Squares[y].Color[x] == "" || c.Squares[y].Color[x] == color {
				c.MovePiece(x, y, color)
				if iswinner(c.Squares, color) {
					// revert side effects
					c.Hub.Clients = newClients
					c.Hub.Colors = players
					replaceBoard(c.Squares, newBoard)
					return -10000, [2]int{movedx, movedy}
				}
				minVal, _ := c.Max(nextColor, color, depth-1, alpha, beta, x, y)
				if minVal < val {
					sq = [2]int{x, y}
					val = minVal
				}
				if val < beta {
					beta = val
				}
				c.Hub.Clients = newClients
				c.Hub.Colors = players
				replaceBoard(c.Squares, newBoard)
				if alpha >= beta {
					return val, sq
				}
			}
		}
	}
	return val, sq
}

func copyBoard(oldBoard []*Squares) []*Squares {
	newBoard := make([]*Squares, len(oldBoard))
	for y, valy := range oldBoard {
		newBoard[y] = &Squares{
			Len:   valy.Len,
			Color: make([]string, valy.Len),
			Cur:   make([]int, valy.Len),
			Max:   make([]int, valy.Len),
		}
		for x := 0; x < valy.Len; x++ {
			newBoard[y].Color[x] = valy.Color[x]
			newBoard[y].Cur[x] = valy.Cur[x]
			newBoard[y].Max[x] = valy.Max[x]
		}
	}
	return newBoard
}

func copyClients(oldClients map[*Client]int) map[*Client]int {
	newClients := make(map[*Client]int)
	for client, val := range oldClients {
		newClients[client] = val
	}
	return newClients
}

func iswinner(board []*Squares, color string) bool {
	isdead := true
	for _, squares := range board {
		for _, c := range squares.Color {
			if c != color || c != "" {
				isdead = false
				break
			}
		}
	}
	return isdead
}

// Put contents of newboard into oldboard in-place
func replaceBoard(oldboard []*Squares, newboard []*Squares) {
	for y, squares := range newboard {
		for x := 0; x < squares.Len; x++ {
			oldboard[y].Color[x] = newboard[y].Color[x]
			oldboard[y].Cur[x] = newboard[y].Cur[x]
			oldboard[y].Max[x] = newboard[y].Max[x]
		}
	}
}
