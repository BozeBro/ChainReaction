package webserver

type Chain struct {
	Len     int
	Squares []*Squares
}
type Squares struct {
	Len   int      // Length of each array
	Cur   []int    // How many circles are in the square
	Max   []int    // Carrying Capacity of the square
	Color []string // The color that occupies. "" if none
}



/*
Animation data is data sent to Front end that will show animation
Moved / static data are the circles that remain from the explosion
*/
//explodeFunc used to clean syntax
type explodeFunc func([][]int, string) ([][]int, [][]int, [][]int)

func Chained(explode explodeFunc, exp [][]int, color string) ([][][]int, [][][]int) {
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
func (c *Chain) Explode(exp [][]int, color string) ([][]int, [][]int, [][]int) {
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
				animations = append(animations, []int{d[0], d[1]})
				sq := c.Squares[y]
				sq.Color[x] = color
				sq.Cur[x] += 1
				if sq.Cur[x] >= sq.Max[x] {
					sq.Cur[x] = 0
					sq.Color[x] = ""
					expN = append(expN, []int{x, y})
				}
				moved = append(moved, []int{x, y, sq.Cur[x]})
			}
		}
	}
	return expN, moved, animations
}
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
	return Chained(c.Explode, [][]int{[]int{x, y}}, color)
}

func (c *Chain) InitBoard(rows, cols int) {
	rows, cols = makeLegal(rows), makeLegal(cols)
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
	c.Squares =  sq
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

func (c *Chain) GetRows() int {
	return c.Squares[0].Len
}
func (c *Chain) GetCols() int {
	return c.Len
}
func makeLegal(dimension int) int {
	/*
		Disallow any dimension under 5 and over 30
	*/
	if dimension < 5 {
		dimension = 5
	} else if dimension > 30 {
		dimension = 30
	}
	return dimension
}
