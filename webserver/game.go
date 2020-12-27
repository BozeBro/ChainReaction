package webserver

type Chain struct {
	Rows, Cols int
	squares [][]int
}

type Game interface {
	InitBoard(rows, cols int) [][]int
}