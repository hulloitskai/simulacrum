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
	size  int
	color color.Color

	intent motion.Intent
	motion *motion.Controller
}

// NewBlob creates a Blob with random properties.
func NewBlob() *Blob {
	return &Blob{
		size:   6 + rng.Intn(6), // random between 6 and 12
		color:  randColor(),     // random color in colors
		motion: motion.NewController(),
	}
}

var _ world.Entity = (*Blob)(nil)

// Bounds returns the bounds of the Blob without actually drawing it.
func (b *Blob) Bounds() image.Rectangle {
	return image.Rect(0, 0, b.size, b.size*2)
}

// Render returns an image representation of the Blob.
func (b *Blob) Render() (image.Image, error) {
	// Blob sizing.
	var (
		bw = b.size     // blob width
		bh = b.size * 2 // blob height
	)

	// Draw canvas.
	canvas, err := ebiten.NewImage(
		b.size*2, // canvas width
		b.size*2, // canvas height
		ebiten.FilterDefault,
	)
	if err != nil {
		return nil, errors.Wrap(err, "blob: create canvas")
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
		return nil, errors.Wrap(err, "blob: draw shadow")
	}

	// Draw blob on canvas.
	if err := func() error {
		img, err := ebiten.NewImage(bw, bh, ebiten.FilterDefault)
		if err != nil {
			return err
		}
		if err = img.Fill(b.color); err != nil {
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

// Update updates the Blob's internal state according to its surroundings.
func (b *Blob) Update(w *world.World, pos world.Position) {
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

	// Discard unused variables.
	_, _ = w, pos
}
