package webserver

import (
	"math/rand"
)
// COLORS provides available colors to choose from
// That is until I can find a color library.
var COLORS = []string{
	"Brown",
	"BlueViolet", "Red",
	"Aquamarine", "Green",
	"Brown", "DarkOrange",
	"DeepPink",
}
// RandomColor grabs a random color from COLORS
func RandomColor() string {
	clength := len(COLORS)
	rand := rand.Intn(clength)
	return COLORS[rand]
}
