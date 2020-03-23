package world

import (
	"bytes"
	"image"
	"sort"

	"github.com/cockroachdb/errors"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
)

// The World is the setting that the simulacrum takes place in.
type World struct {
	// World dimensions, in terms of number of tiles.
	cols, rows int
	size       int // width and height of each tile

	image *ebiten.Image // tile image
	tiles []int         // tile numbers
	count int           // number of tiles in row of image

	// A set of objects.
	objects []*object
}

// New creates a World with the given size.
func New(width, height int) *World {
	src, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		panic(errors.Wrap(err, "decode image"))
	}
	img, err := ebiten.NewImageFromImage(src, ebiten.FilterDefault)
	if err != nil {
		panic(errors.Wrap(err, "convert image"))
	}

	// Generate tile data.
	var (
		size  = 16
		cols  = 25
		dx    = width / size
		dy    = height / size
		tiles = make([]int, dx*dy)
	)
	for i := range tiles {
		n := rng.Intn(32)
		switch n {
		case 0:
			tiles[i] = 218 // leaf 1 tile
		case 1:
			tiles[i] = 219 // leaf 2 tile
		case 2:
			tiles[i] = 244 // leaf 3 tile
		default:
			tiles[i] = 243 // plain tile
		}
	}

	return &World{
		cols:  dx,
		rows:  dy,
		size:  size,
		count: cols,
		image: img,
		tiles: tiles,
	}
}

// NewForScreen creates a World with a size matching the target screen.
func NewForScreen(screen *ebiten.Image) *World {
	b := screen.Bounds()
	return New(b.Dx(), b.Dy())
}

// Bounds returns the bounds of the World.
func (w *World) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.cols*w.size, w.rows*w.size)
}

// Draw draws the world on the screen.
func (w *World) Draw(screen *ebiten.Image) error {
	for x := 0; x < w.cols; x++ {
		for y := 0; y < w.rows; y++ {
			var (
				num = w.tiles[y*w.size+x]
				nx  = (num % w.count) * w.size
				ny  = (num / w.count) * w.size
				img = w.image.SubImage(image.Rect(nx, ny, nx+w.size, ny+w.size))
			)

			var opts ebiten.DrawImageOptions
			opts.GeoM.Translate(float64(x*w.size), float64(y*w.size))

			if err := screen.DrawImage(img.(*ebiten.Image), &opts); err != nil {
				return errors.Wrapf(err, "(%d, %d)", x, y)
			}
		}
	}

	for _, obj := range w.objects {
		if err := obj.Draw(screen); err != nil {
			return errors.Wrap(err, "draw object")
		}
	}
	return nil
}

// Update instructs the world to update its internal state.
func (w *World) Update(screen *ebiten.Image) error {
	// Update each object.
	for i, obj := range w.objects {
		if err := obj.Update(w); err != nil {
			return errors.Wrapf(err, "object %d", i)
		}
	}

	// Sort w.objects such that any object that is "below" another in the world
	// will be rendered above the other one.
	sort.Slice(w.objects, func(i, j int) bool {
		var (
			ib = w.objects[i].Bounds()
			jb = w.objects[j].Bounds()
		)
		return ib.Min.Y < jb.Min.Y
	})

	// Discard unused variables.
	_ = screen
	return nil
}

// Spawn spawns a new Entity somewhere in the World.
func (w *World) Spawn(ent Entity) error {
	var (
		b = w.Bounds()
		x = rng.Intn(b.Dx())
		y = rng.Intn(b.Dy())
	)
	return w.SpawnAt(ent, image.Pt(x, y))
}

// SpawnAt spawns an Entity at the given position.
func (w *World) SpawnAt(ent Entity, pt image.Point) error {
	if !pt.In(w.Bounds()) {
		return errors.New("world: out of bounds")
	}
	obj := newObject(ent, PointPos(pt))
	w.objects = append(w.objects, obj)
	return nil
}

// Count returns the number of entities in the world.
func (w *World) Count() int { return len(w.objects) }

// EntityNear returns an Entity within a certain pixel distance of pos.
//
// Returns nil if no such Entity is found.
func (w *World) EntityNear(pos Position, within int) Entity {
	var (
		pt = pos.Point()
		dt = image.Pt(within, within)
	)
	area := image.Rectangle{Min: pt.Sub(dt), Max: pt.Add(dt)}
	for _, obj := range w.objects {
		if obj.Bounds().Overlaps(area) {
			return obj.Entity()
		}
	}
	return nil
}
