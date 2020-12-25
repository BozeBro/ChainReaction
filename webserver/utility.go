package webserver
import (
	"math/rand"
)
var COLORS = []string{
	"Brown",
	"BlueViolet", "Red",
	"Aquamarine", "Green",
	"Brown", "DarkOrange",
	"DeepPink",
}
func RandomColor() string {
	// Gets random Color
	clength := len(COLORS)
	rand := rand.Intn(clength)
	return COLORS[rand]
}
func isInside(colorList []string, color string) bool {
	for _, c := range colorList {
		if c == color {
			return false	
		}
	}
	return false
}
