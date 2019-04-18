package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var img *ebiten.Image

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, 50)
	op.GeoM.Scale(1.5, 1)
	screen.DrawImage(img, op)
	return nil
}

func main() {
	if err := ebiten.Run(update, 640, 480, 1, "Geometry Matrix"); err != nil {
		log.Fatal(err)
	}
}
