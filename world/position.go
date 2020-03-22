package world

import "image"

// A Position describes a point in the World.
type Position struct {
	X, Y float64
}

// Pos creates a Position.
func Pos(x, y float64) Position {
	return Position{
		X: x,
		Y: y,
	}
}

// Point converts a Position into an image.Point.
func (pos Position) Point() image.Point {
	return image.Pt(int(pos.X), int(pos.Y))
}

// PointPos creates a Position from an image.Point.
func PointPos(pt image.Point) Position {
	return Position{
		X: float64(pt.X),
		Y: float64(pt.Y),
	}
}
