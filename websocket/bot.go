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

// getChains tries to find optimizations to ignore reprocessing redundant squares
func getChains(oldBoard []*Squares, color string, x, y int) [][]int {
	board := copyBoard(oldBoard)
	visiting := [][]int{{x, y}}
	if board[y].Cur[x] == board[y].Max[x]-1 && board[y].Color[x] == color {
		for len(visiting) > 0 {
			last := len(visiting) - 1
			nx, ny := visiting[last][0], visiting[last][1]
			board[ny].Cur[nx] = 0
			total, coords := findneighbors(nx, ny, board[0].Len, len(board))
			for _, coord := range coords {
				newX, newY := coord[0], coord[1]
				if board[newY].Cur[newX] == total-1 {
					visiting = append(visiting, coord)
				}
			}
		}
	}
	return visiting
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

func (c *Chain) signiMoves(color string) [][]int {
	validMoves := make([][]int, 0, c.Squares[0].Len*c.Len)
	redundant := make(map[[2]int]bool, 0)
	// heuristic for ignoring moves that will lead to the same outcome
	// redundant acting as a set to account for already seen squares
	for y := 0; y < c.Len; y++ {
		for x := 0; x < c.Squares[0].Len; x++ {
			if (c.Squares[y].Color[x] == "" || c.Squares[y].Color[x] == color) && redundant[[2]int{x, y}] == false {
				chained := getChains(c.Squares, color, x, y)
				redundant[[2]int{x, y}] = true
				validMoves = append(validMoves, chained[0])
				for _, sq := range chained {
					nx, ny := sq[0], sq[1]
					// objects can get chained into other people's circles.
					redundant[[2]int{nx, ny}] = true
				}
			}
		}
	}
	return validMoves
}

// defaultMoves grabs all legal moves that are possible
func (c *Chain) defaultMoves(color string, rows, cols int) [][]int {
	moves := make([][]int, 0, rows*cols)
	for y := 0; y < cols; y++ {
		for x := 0; x < rows; x++ {
			if c.Squares[y].Color[x] == "" || c.Squares[y].Color[x] == color {
				moves = append(moves, []int{x, y})
			}
		}
	}
	return moves
}

// Maximizing player
// Return greatest number possible
func (c *Chain) Max(color, nextColor string, depth, alpha, beta, movedx, movedy int) (int, [2]int) {
	//counter++
	sq := [2]int{movedx, movedy}
	//println(depth, " In depth and Max ", counter)
	if depth == 0 {
		// use because we are evaluating previous player's move
		boardValue := BrillHeuristic(c.Squares, color)
		return boardValue, sq
	}
	newBoard := copyBoard(c.Squares)
	newClients := copyClients(c.Hub.Clients)
	players := make([]string, len(c.Hub.Colors))
	for ind, player := range c.Hub.Colors {
		players[ind] = player
	}
	//validMoves := c.signiMoves(color)
	validMoves := c.defaultMoves(color, c.Squares[0].Len, c.Len)
	for _, pos := range validMoves {
		x, y := pos[0], pos[1]
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
		if maxVal > alpha {
			sq = [2]int{x, y}
			alpha = maxVal
			if alpha >= beta {
				return alpha, sq
			}
		}
	}
	return alpha, sq
}

// Minimizing player
// Return smallest number possible
// Look at Max() for more documentation
func (c *Chain) Min(color, nextColor string, depth, alpha, beta, movedx, movedy int) (int, [2]int) {
	//counter++
	//println(depth, " In depth and Min ", counter)
	sq := [2]int{movedx, movedy}
	if depth == 0 {
		// nextColor is actually our original color
		boardValue := BrillHeuristic(c.Squares, nextColor)
		return boardValue, sq
	}
	players := make([]string, len(c.Hub.Colors))
	for ind, player := range c.Hub.Colors {
		players[ind] = player
	}
	newBoard := copyBoard(c.Squares)
	newClients := copyClients(c.Hub.Clients)
	validMoves := c.defaultMoves(color, c.Squares[0].Len, c.Len)
	for _, pos := range validMoves {
		x, y := pos[0], pos[1]
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
		if minVal < beta {
			sq = [2]int{x, y}
			beta = minVal
			if alpha >= beta {
				return beta, sq
			}
		}
	}
	return beta, sq
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
