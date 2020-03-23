package life

import (
	"image/color"
	"math/rand"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

var colors = []color.Color{
	colorHex("#FFD54F"),
	colorHex("#FFCA28"),
	colorHex("#FFC107"),
	colorHex("#FFB300"),
	colorHex("#FFA000"),
	colorHex("#FFA000"),
	colorHex("#FF6F00"),
}

func colorHex(hex string) color.Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		panic(err)
	}
	return c
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func randColor() color.Color { return colors[rng.Intn(len(colors))] }
