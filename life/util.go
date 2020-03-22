package life

import (
	"image/color"
	"math/rand"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

var colors = []color.Color{
	// colorHex("#E1500F"),
	colorHex("#F28A28"),
	colorHex("#FEAA38"),
	colorHex("#FECA4B"),
	colorHex("#7567F1"),
	colorHex("#5C44D9"),
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
