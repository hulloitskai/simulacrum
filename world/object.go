package world

import (
	"image"

	"github.com/cockroachdb/errors"
	"github.com/hajimehoshi/ebiten"
)

// An object wraps an Entity with positional state information.
type object struct {
	ent Entity
	pos Position
}

func newObject(ent Entity, pos Position) *object {
	return &object{
		ent: ent,
		pos: pos,
	}
}

// Entity returns the object's underlying Entity.
func (obj *object) Entity() Entity { return obj.ent }

// Draw draws the object on a screen.
func (obj *object) Draw(screen *ebiten.Image) error {
	// Render entity.
	src, err := obj.ent.Render()
	if err != nil {
		return errors.Wrap(err, "render entity")
	}

	// Configure draw options.
	var opts ebiten.DrawImageOptions
	opts.GeoM.Translate(float64(obj.pos.X), float64(obj.pos.Y))

	// Draw entity, with conversion if necessary.
	img, ok := src.(*ebiten.Image)
	if !ok {
		img, err = ebiten.NewImageFromImage(src, ebiten.FilterDefault)
		if err != nil {
			return errors.Wrap(err, "convert image")
		}
	}
	return screen.DrawImage(img, &opts)
}

// Update is called once a tick to instruct the object to update its internal
// state.
func (obj *object) Update(w *World) error {
	ent := obj.ent
	if err := ent.Update(w, obj.pos); err != nil {
		return err
	}

	// Update the object's position using its entity's motion.
	{
		var (
			eb     = ent.Bounds()
			wb     = w.Bounds()
			pos    = obj.pos
			dx, dy = ent.Motion()
		)

		pos.X += dx
		pos.Y += dy

		// Prevent object from going out-of-bounds.
		switch {
		case int(pos.X)+1+eb.Dx() >= wb.Dx():
			pos.X = float64(wb.Dx() - eb.Dx() - 1)
		case pos.X < 0:
			pos.X = 0
		case int(pos.Y)+1+eb.Dy() >= wb.Dy():
			pos.Y = float64(wb.Dy() - eb.Dy() - 1)
		case pos.Y < 0:
			pos.Y = 0
		}

		obj.pos = pos
	}
	return nil
}

// Bounds returns the bounds of the object.
func (obj *object) Bounds() image.Rectangle {
	return obj.ent.Bounds().Add(obj.pos.Point())
}
