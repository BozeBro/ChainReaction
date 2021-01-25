package websocket

import (
	"math/rand"
)

// A Player can only be a color within this const
var COLORS = []string{
	"Brown",
	"BlueViolet", "Red",
	"Aquamarine", "Green",
	"Brown", "DarkOrange",
	"DeepPink", "Gray", 
	"Black", "Granola",
	"Tan", "SlateBlue",
	"Tomato", "Cyan",
	"Olive", "Jellyfish",
	"PineGreen", "Cantaloupe",
	"DeepPeach", "Maroon",
	"RosyBrown", "CadillacPink",
	"Magenta", "WisteriaPurple",
	"Pearl", "Moccasin",
}

// RandomColor grabs a random color from global COLORS
func RandomColor() string {
	clength := len(COLORS)
	rand := rand.Intn(clength)
	return COLORS[rand]
}
