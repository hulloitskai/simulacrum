package world

import (
	"image"
)

// An Entity can be placed, and can interact with the world.
type Entity interface {
	// Bounds returns the bounds of the Entity without rendering it.
	Bounds() image.Rectangle

	// Render renders an image of the Entity.
	Render() (image.Image, error)

	// Motion describes the desired velocity of an Entity.
	Motion() (dx, dy float64)

	// Update is called once a tick to instruct the Entity to update its internal
	// state.
	//
	// The pos of an Entity is the coordinate of its left-bottommost pixel.
	Update(w *World, pos Position)
}
