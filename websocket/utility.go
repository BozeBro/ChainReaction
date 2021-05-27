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
	"Black", "darkkhaki",
	"Tan", "SlateBlue",
	"Tomato", "Cyan",
	"Olive", "cornsilk",
	"mediumspringgreen", "darkslategray",
	"peachpuff", "Maroon",
	"RosyBrown", "Yellow",
	"Magenta", "Indigo",
	"mediumvioletred", "Moccasin",
}

// RandomColor grabs a random color from global COLORS
func RandomColor() string {
	clength := len(COLORS)
	ran := rand.Intn(clength)
	return COLORS[ran]
}
