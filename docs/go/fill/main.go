package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.RGBA{0xff, 0, 0, 0xff})
	return nil
}

func main() {
	if err := ebiten.Run(update, 640, 480, 1, "Fill"); err != nil {
		log.Fatal(err)
	}
}
