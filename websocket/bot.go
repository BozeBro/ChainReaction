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
				visiting := make([][]int, 0)
				visiting = append(visiting, []int{x, y})
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
func (c *Chain) Minimax(color string, maximizing bool, players []string, depth, alpha, beta, movedx, movedy int) (int, [2]int) {
	square := [2]int{movedx, movedy}
	nextColor := players[0]
	curNums := c.Hub.Clients
	if color == players[0] {
		nextColor = players[1]
	}
	if depth == 0 {
		return BrillHeuristic(c.Squares, color), [2]int{movedx, movedy}
	} else if maximizing {
		val := int(math.Inf(-1))
		for y := 0; y < c.Len; y++ {
			for x := 0; x < c.Squares[0].Len; x++ {
				if c.Squares[y].Color[x] == color || c.Squares[y].Color[x] == "" {
					c.MovePiece(x, y, color)
					for client, squares := range c.Hub.Clients {
						if client.Color == color {
							if squares == 0 {
								c.Hub.Clients = curNums
								return 10000, [2]int{x, y}
							}
							break
						}
					}
					maxVal, _ := c.Minimax(nextColor, !maximizing, players, depth-1, alpha, beta, x, y)

					if maxVal > val {
						square = [2]int{x, y}
						val = maxVal
					}
					if val > alpha {
						alpha = val
					}
					c.Hub.Clients = curNums
					if alpha >= beta {
						return val, square
					}
				}
			}
		}
		return val, square
	}
	val := int(math.Inf(1))
	for y := 0; y < c.Len; y++ {
		for x := 0; x < c.Squares[0].Len; x++ {
			if c.Squares[y].Color[x] == color || c.Squares[y].Color[x] == "" {
				c.MovePiece(x, y, color)
				for client, squares := range c.Hub.Clients {
					if client.Color == color {
						if squares == 0 {
							c.Hub.Clients = curNums
							return 10000, [2]int{x, y}
						}
						break
					}
				}
				minVal, _ := c.Minimax(nextColor, !maximizing, players, depth-1, alpha, beta, x, y)

				if minVal < val {
					square = [2]int{x, y}
					val = minVal
				}
				if val > beta {
					beta = val
				}
				c.Hub.Clients = curNums
				if alpha >= beta {
					return val, square
				}
			}
		}
	}
	return val, square
}
