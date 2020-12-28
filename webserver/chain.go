package webserver

type Chain struct {
	Len     int
	Squares []*Squares
	Hub     *Hub
}
type Squares struct {
	Len   int      // Length of each array
	Cur   []int    // How many circles are in the square
	Max   []int    // Carrying Capacity of the square
	Color []string // The color that occupies. "" if none
}

func (c *Chain) InitBoard(rows, cols int) {
	rows, cols = makeLegal(5, rows, 30), makeLegal(5, cols, 30)
	sq := make([]*Squares, cols)
	c.Len = cols
	for y := 0; y < cols; y++ {
		sq[y] = &Squares{
			Len:   rows,
			Cur:   make([]int, rows),
			Max:   make([]int, rows),
			Color: make([]string, rows),
		}
		for x := 0; x < rows; x++ {
			sq[y].Max[x], _ = c.findneighbors(x, y, rows, cols)
		}
	}
	c.Squares = sq
}
func (c *Chain) findneighbors(x, y, rows, cols int) (int, [][]int) {
	// Returns maximum neighbors and their coords
	totalNeighbros := 0
	coords := make([][]int, 0, 4)
	for _, v := range [][]int{
		[]int{1, 0}, []int{-1, 0}, []int{0, 1}, []int{0, -1}} {
		nx := x + v[0]
		ny := y + v[1]
		if IsLegalMove(c, nx, ny) {
			totalNeighbros += 1
			coords = append(coords, []int{nx, ny})
		}
	}
	return totalNeighbros, coords
}

/*
Animation data is data sent to Front end that will show animation
Moved / static data are the circles that remain from the explosion
*/

func (c *Chain) MovePiece(x, y int, color string) ([][][]int, [][][]int) {
	/*
		x : x coordinate of the user clicked square
		y : y coordinate of the user clicked square
		color : color of the user
	*/
	c.Squares[y].Cur[x] += 1
	if c.Squares[y].Cur[x] < c.Squares[y].Max[x] {
		return [][][]int{[][]int{[]int{x, y}}}, make([][][]int, 0)
	}
	c.Squares[y].Cur[x] = 0
	return chained(c.explode, [][]int{[]int{x, y}}, color)
}

//explodeFunc used to clean syntax
type explodeFunc func([][]int, string) ([][]int, [][]int, [][]int)

func chained(explode explodeFunc, exp [][]int, color string) ([][][]int, [][][]int) {
	/*
		explode - function to execute to receive animation data
		exp -  nested array that contains coords of exploding squares
		color - Color of the person that is moving
		-----
		Loops through the explode function until no more until come out.
		Receives animation data and static data
	*/
	animations := make([][][]int, 0)
	moved := make([][][]int, 0)
	for len(exp) != 0 {
		newExp, newAni, newMove := explode(exp, color)
		animations = append(animations, newAni)
		moved = append(moved, newMove)
		exp = newExp
	}
	return animations, moved
}
func (c *Chain) explode(exp [][]int, color string) ([][]int, [][]int, [][]int) {
	/*
		exp - Current exploding squares
		color - color of the user
		Function to handle exploding squares.
		Returns a frame of animation and static data
	*/
	expN := make([][]int, 0)
	moved := make([][]int, 0)
	animations := make([][]int, 0)
	for _, coords := range exp {
		// d is all the possible neighbors of the coords
		for _, d := range [][]int{
			[]int{1, 0},
			[]int{-1, 0},
			[]int{0, 1},
			[]int{0, -1},
		} {
			x, y := coords[0]+d[0], coords[1]+d[1]
			if IsLegalMove(c, x, y) {
				sq := c.Squares[y]
				isdead := c.UpdateColor(color, sq.Color[x])
				deletedColor := sq.Color[x]
				animations = append(animations, []int{d[0], d[1]})
				sq.Color[x] = color
				sq.Cur[x] += 1
				if sq.Cur[x] >= sq.Max[x] {
					isdead = c.UpdateColor(color, sq.Color[x])
					sq.Cur[x] = 0
					sq.Color[x] = ""
					expN = append(expN, []int{x, y})
				}
				moved = append(moved, []int{x, y, sq.Cur[x]})
				if isdead {
					for index := 0; index < len(c.Hub.Colors); index++ {
						if c.Hub.Colors[index]  == deletedColor {
							c.Hub.Colors = append(c.Hub.Colors[:index], c.Hub.Colors[index+1:]...)
						}
					}
					if len(c.Hub.Colors) == 1 {
						return nil, animations, moved
					}
				}
			}
		}
	}
	return expN, moved, animations
}
func (c *Chain) UpdateColor(newColor, oldColor string) bool {
	/*
		Update amount of squares each player controls.
		true if oldColor is dead / out of squares
	*/
	if newColor == oldColor {
		return false
	}
	dead := false
	for client := range c.Hub.Clients {
		if client.Color == newColor {
			c.Hub.Clients[client] += -1
			if c.Hub.Clients[client] == 0 {
				dead = true
			}
		} else if client.Color == oldColor {
			c.Hub.Clients[client] += 1
		}
	}
	return dead
}

func (c *Chain) GetRows() int {
	return c.Squares[0].Len
}
func (c *Chain) GetCols() int {
	return c.Len
}
func makeLegal(lower, dimension, upper int) int {
	/*
		Disallow any dimension under 5 and over 30
	*/
	if dimension < lower {
		dimension = lower
	} else if dimension > upper {
		dimension = upper
	}
	return dimension
}
