package life

import (
	"image"
	"image/color"

	"github.com/cockroachdb/errors"
	"github.com/hajimehoshi/ebiten"

	"go.stevenxie.me/simulacrum/life/motion"
	"go.stevenxie.me/simulacrum/world"
)

// A Blob is a basic living entity.
type Blob struct {
	color color.Color

	age       int
	horniness int

	intent motion.Intent
	motion *motion.Controller
}

// NewBlob creates a Blob with random properties.
func NewBlob() *Blob {
	return &Blob{
		color:  randColor(), // random color in colors
		motion: motion.NewController(),
	}
}

var _ world.Entity = (*Blob)(nil)

// Bounds returns the bounds of the Blob without actually drawing it.
func (b *Blob) Bounds() image.Rectangle {
	size := b.Size()
	return image.Rect(0, 0, size, size*2)
}

// Size returns the size factor of the Blob.
func (b *Blob) Size() int {
	s := b.age / 600
	switch {
	case s < 6:
		s = 6
	case s > 20:
		s = 20
	}
	return s
}

// Render returns an image representation of the Blob.
func (b *Blob) Render() (image.Image, error) {
	// Blob sizing.
	var (
		sz = b.Size()
		bw = sz     // blob width
		bh = sz * 2 // blob height
	)

	// Draw canvas.
	canvas, err := ebiten.NewImage(sz*2, sz*2, ebiten.FilterDefault)
	if err != nil {
		return nil, errors.Wrap(err, "create canvas")
	}

	// Draw shadow on canvas.
	if err := func() error {
		img, err := ebiten.NewImage(bw, bh, ebiten.FilterDefault)
		if err != nil {
			return err
		}
		if err = img.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 100}); err != nil {
			return err
		}

		var (
			opts ebiten.DrawImageOptions
			b    = canvas.Bounds()
		)
		opts.GeoM.Translate(float64(b.Dx())*0.84, float64(b.Dy())/2)
		opts.GeoM.Skew(-0.7, 0)
		return canvas.DrawImage(img, &opts)
	}(); err != nil {
		return nil, errors.Wrap(err, "draw shadow")
	}

	// Draw blob on canvas.
	if err := func() error {
		img, err := ebiten.NewImage(bw, bh, ebiten.FilterDefault)
		if err != nil {
			return err
		}

		// Use Blob's color, but flash white if really horny.
		c := b.color
		if b.horniness >= reallyHorny {
			if b.age%60 < 30 {
				c = color.White
			}
		}

		if err = img.Fill(c); err != nil {
			return err
		}
		return canvas.DrawImage(img, nil)
	}(); err != nil {
		return nil, err
	}

	return canvas, nil
}

// Motion returns the velocity of the Blob.
func (b *Blob) Motion() (x, y float64) { return b.motion.Motion() }

const reallyHorny = 3

// Update updates the Blob's internal state according to its surroundings.
func (b *Blob) Update(w *world.World, pos world.Position) error {
	// Update age and horniness.
	b.age++
	if b.age == 1800 {
		// Make the Blob really horny the second it is "of age".
		b.horniness = reallyHorny
	} else if b.age > 1800 { // approx. 30 seconds
		if b.age%600 == 0 { // horniness goes up at 10-second intervals
			b.horniness++
		}
	}

	// If there is no next movement intent, choose one at random (but tend
	// to stay still, 1 / 24 times).
	//
	// Then, update the motion.
	if !b.motion.Next() {
		n := rng.Intn(104)
		if n > 4 {
			n = 0
		}
		b.intent = motion.Intent(n)
	}
	b.motion.Update(b.intent)

	if b.horniness >= reallyHorny {
		if other, ok := w.EntityNear(pos, b.Size()).(Creature); ok {
			if offspring := b.Fuck(other); offspring != nil {
				center := image.Pt(int(pos.X)+(b.Size()/2), int(pos.Y))
				if err := w.SpawnAt(offspring, center); err != nil {
					return errors.Wrap(err, "spawn offspring")
				}
			}
		}
	}

	// Discard unused variables.
	_, _ = w, pos
	return nil
}

// Fuck implements Creature.Fuck.
func (b *Blob) Fuck(other Creature) (offspring Creature) {
	o, ok := other.(*Blob)
	if !ok {
		return nil
	}
	if o == b {
		return nil
	}
	if (b.age < 1800) || (o.age < 1800) {
		return nil
	}
	b.horniness = 0
	o.horniness = 0
	return NewBlob()
}
