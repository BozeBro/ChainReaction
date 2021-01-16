package websocket

import (
	"math/rand"
)

// Slice of colors that a player can be.
var COLORS = []string{
	"Brown",
	"BlueViolet", "Red",
	"Aquamarine", "Green",
	"Brown", "DarkOrange",
	"DeepPink",
}

// RandomColor grabs a random color from global COLORS
func RandomColor() string {
	clength := len(COLORS)
	rand := rand.Intn(clength)
	return COLORS[rand]
}
