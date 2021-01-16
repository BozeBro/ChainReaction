package websocket

// Chain contains data relevant for Chain Reaction Game
// Satisfies the Game interface
type Chain struct {
	Len     int
	Squares []*Squares
	Hub     *Hub
}

// Squares contains data about each square
type Squares struct {
	Len   int      // Length of each array
	Cur   []int    // How many circles are in the square
	Max   []int    // Carrying Capacity of the square
	Color []string // The color that occupies. "" if none
}

// InitBoard Creates a board with dimensions rows x cols
// Each y level has a Squares struct which act as the whole x row
func (c *Chain) InitBoard(rows, cols int) {
	rows, cols = makeLegal(5, rows, 30), makeLegal(5, cols, 30)
	c.Squares = make([]*Squares, cols)
	c.Len = cols
	for y := 0; y < cols; y++ {
		c.Squares[y] = &Squares{
			Len:   rows,
			Cur:   make([]int, rows),
			Max:   make([]int, rows),
			Color: make([]string, rows),
		}
		for x := 0; x < rows; x++ {
			c.Squares[y].Max[x], _ = c.findneighbors(x, y, rows, cols)
		}
	}
}

// Checks if a square position exists on a board horizontally and vertically
func (c *Chain) findneighbors(x, y, rows, cols int) (int, [][]int) {
	totalNeighbros := 0
	coords := make([][]int, 0, 4)
	for _, v := range [][]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		nx := x + v[0]
		ny := y + v[1]
		if IsBounded(c, nx, ny) {
			totalNeighbros++
			coords = append(coords, []int{nx, ny})
		}
	}
	return totalNeighbros, coords
}

/*
Animation data is data sent to Front end that will show animation
Moved / static data are the circles that remain from the explosion
*/

// MovePiece is requirement for Game Interface
// MovePiece Moves the piece on the chain board
// It will call the chained(explode) function to handle explosion
func (c *Chain) MovePiece(x, y int, color string) ([][][]int, [][][]int) {
	/*
		x : x coordinate of the user clicked square
		y : y coordinate of the user clicked square
		color : color of the user
	*/
	c.Squares[y].Cur[x]++
	if c.Squares[y].Cur[x] < c.Squares[y].Max[x] {
		c.UpdateColor(color, c.Squares[y].Color[x])
		c.Squares[y].Color[x] = color
		return make([][][]int, 0), [][][]int{{{x, y, c.Squares[y].Cur[x]}}}
	}
	c.Squares[y].Cur[x] = 0
	c.Squares[y].Color[x] = ""
	c.UpdateColor("", color)
	return chained(c.explode, [][]int{{x, y}}, color)
}

//explodeFunc used to clean syntax
type explodeFunc func([][]int, string) ([][]int, [][]int, [][]int)

// chained is a helper function that will continually call c.explode
// appends animation pieces from c.explode to an array, which will be returned
func chained(explode explodeFunc, exp [][]int, color string) ([][][]int, [][][]int) {
	/*
		explode - function to execute to receive animation data
		exp -  nested array that contains coords of exploding squares
		color - Color of the person that is moving
		-----
		Loops through the explode function until no more until come out.
		Receives animation data and static data
	*/
	x, y := exp[0][0], exp[0][1]
	animations := make([][][]int, 0)
	moved := [][][]int{{{x, y, 0}}} // static is going to be zero
	for len(exp) != 0 {
		newExp, newAni, newMoves := explode(exp, color)
		animations = append(animations, newAni)
		moved = append(moved, newMoves)
		exp = newExp
	}
	return animations, moved
}

// explode simulates the actual game logic of Chain Logic
// Add animation data for each neighbor
// explode iterates each exploding square and check if neigboring squares will also explode
// Else it will just add a square and add to static animation data
func (c *Chain) explode(exp [][]int, color string) ([][]int, [][]int, [][]int) {
	/*
		exp - Current exploding squares
		color - color of the user that is making the move
		Function to handle exploding squares.
		Returns a frame of animation and static data
	*/
	expN := make([][]int, 0) // Neighbors that are going to explode next iteration
	moved := make([][]int, 0)
	animations := make([][]int, 0)
	for _, coords := range exp {
		// d is all the possible neighbors of the coords
		for _, d := range [][]int{
			{1, 0},
			{-1, 0},
			{0, 1},
			{0, -1},
		} {
			x, y := coords[0]+d[0], coords[1]+d[1]
			if !IsBounded(c, x, y) {
				continue
			}
			// (coords of explosion site), direction they are going
			animations = append(animations, []int{coords[0], coords[1], d[0], d[1]})
			sq := c.Squares[y]
			deletedColor := sq.Color[x]
			if c.UpdateColor(color, deletedColor) {
				// Remove the players from the list of alive people
				for index := 0; index < len(c.Hub.Colors); index++ {
					if c.Hub.Colors[index] == deletedColor {
						c.Hub.Colors = append(c.Hub.Colors[:index], c.Hub.Colors[index+1:]...)
						if index <= c.Hub.i {
							c.Hub.i--
						}
						break

					}
				}
			}
			sq.Color[x] = color
			sq.Cur[x]++
			if sq.Cur[x] == sq.Max[x] {
				sq.Cur[x] = 0
				sq.Color[x] = ""
				_ = c.UpdateColor("", color)
				expN = append(expN, []int{x, y})
			}
			moved = append(moved, []int{x, y, sq.Cur[x]})

		}
	}
	if len(c.Hub.Colors) == 1 {
		return make([][]int, 0), animations, moved
	}
	return expN, animations, moved
}

// UpdateColor updates amount of squares each player controls.
// results true if oldColor is dead / out of squares
func (c *Chain) UpdateColor(newColor, oldColor string) bool {
	/*
		Update amount of squares each player controls.
		true if oldColor is dead / out of squares
	*/
	dead := false
	if newColor == oldColor {
		return dead
	}
	for client := range c.Hub.Clients {
		if client.Color == oldColor {
			c.Hub.Clients[client]--
			if c.Hub.Clients[client] == 0 {
				dead = true
			}
		} else if client.Color == newColor {
			c.Hub.Clients[client]++
		}
	}
	return dead
}

// GetRows is a requirement for Game interface
// GetRows Gets the rows in the Chain Board
func (c *Chain) GetRows() int {
	return c.Squares[0].Len
}

// GetCols is a requirement for Game interface
// GetCols Gets the cols in the Chain Board
func (c *Chain) GetCols() int {
	return c.Len
}

// makeLegal makes sure that dimensions are legal
// Must be (lower, upper]
func makeLegal(lower, dimension, upper int) int {
	if dimension < lower {
		dimension = lower
	} else if dimension > upper {
		dimension = upper
	}
	return dimension
}
