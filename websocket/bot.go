package websocket

import (
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
				myCircles += board[y].Cur[x]
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
				enemyCircles += board[y].Cur[x]
			}
		}
	}
	score += myCircles
	// A player can only win when they have more than one circle
	if enemyCircles == 0 && myCircles > 1 {
		return 10000
	} else if myCircles == 0 && enemyCircles > 1 {
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
	board := copyBoard(oldBoard)
	lengths := make([]int, 0)
	for y := 0; y < len(board); y++ {
		for x := 0; x < board[0].Len; x++ {
			if board[y].Cur[x] == board[y].Max[x]-1 && board[y].Color[x] == color {
				amount := 0
				visiting := [][]int{{x, y}}
				for len(visiting) > 0 {
					last := len(visiting) - 1
					nx, ny := visiting[last][0], visiting[last][1]
					visiting = visiting[:last]
					board[ny].Cur[nx] = 0
					amount++
					total, coords := findneighbors(nx, ny, board[0].Len, len(board))
					for _, coord := range coords {
						newX, newY := coord[0], coord[1]
						if board[newY].Cur[newX] == total-1 && board[newY].Color[newX] == color {
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
func (c *Chain) Max(color, nextColor string, depth, alpha, beta, movedx, movedy int) (int, [2]int) {
	sq := [2]int{movedx, movedy}
	val := -100000
	if depth == 0 {
		// use because we are evaluating previous player's move
		boardValue := BrillHeuristic(c.Squares, color)
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
				// revert side effects from MovePiece
				maxVal, _ := c.Min(nextColor, color, depth-1, alpha, beta, x, y)
				for client, squares := range newClients {
					c.Hub.Clients[client] = squares
				}
				length := len(players)
				c.Hub.Colors = make([]string, length, length)
				for index, color := range players {
					c.Hub.Colors[index] = color
				}
				replaceBoard(c.Squares, newBoard)
				if maxVal > val {
					sq = [2]int{x, y}
					val = maxVal
					if val > alpha {
						alpha = val
						if alpha >= beta {
							return alpha, sq

						}
					}
				}
			}
		}
	}
	return val, sq
}

// Minimizing player
// Return smallest number possible
// Look at Max() for more documentation
func (c *Chain) Min(color, nextColor string, depth, alpha, beta, movedx, movedy int) (int, [2]int) {
	sq := [2]int{movedx, movedy}
	val := 100000
	if depth == 0 {
		// nextColor is actually our original color
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
				minVal, _ := c.Max(nextColor, color, depth-1, alpha, beta, x, y)
				for client, squares := range newClients {
					c.Hub.Clients[client] = squares
				}
				length := len(players)
				c.Hub.Colors = make([]string, length, length)
				for index, color := range players {
					c.Hub.Colors[index] = color
				}
				replaceBoard(c.Squares, newBoard)
				if minVal < val {
					sq = [2]int{x, y}
					val = minVal
					if val < beta {
						beta = val
						if alpha >= beta {
							return beta, sq
						}
					}
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
