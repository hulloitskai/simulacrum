package main

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"go.stevenxie.me/simulacrum/life"
	"go.stevenxie.me/simulacrum/world"
)

const (
	width  = 480
	height = 320
)

func main() {
	var (
		w      = world.New(width, height)
		update = func(screen *ebiten.Image) error {
			// Check for input.
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				blob := life.NewBlob()
				if err := w.Spawn(blob); err != nil {
					return errors.Wrap(err, "spawn")
				}
			}

			// Update the world.
			w.Update(screen)

			// If drawing is not skipped, draw the world.
			if ebiten.IsDrawingSkipped() {
				return nil
			}
			if err := w.Draw(screen); err != nil {
				return err
			}

			// Print the TPS (ticks-per-second).
			ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

			// If no entities, draw spawn instructions.
			if w.Count() == 0 {
				ebitenutil.DebugPrintAt(
					screen,
					"spawn with [space]",
					(width/2)-50,
					height-30,
				)
			}
			return nil
		}
	)
	if err := ebiten.Run(update, width, height, 2, "simulacrum"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
}
