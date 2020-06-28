package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var img *ebiten.Image

func init() {
	// Create a source image that consists of 4 parts: red, green, blue and yellow.
	img, _ = ebiten.NewImage(100, 100, ebiten.FilterDefault)
	ebitenutil.DrawRect(img, 0, 0, 50, 50, color.RGBA{0xff, 0, 0, 0xff})
	ebitenutil.DrawRect(img, 50, 0, 50, 50, color.RGBA{0, 0xff, 0, 0xff})
	ebitenutil.DrawRect(img, 0, 50, 50, 50, color.RGBA{0, 0, 0xff, 0xff})
	ebitenutil.DrawRect(img, 50, 50, 50, 50, color.RGBA{0xff, 0xff, 0, 0xff})
}

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render the original image at the upper-left corner.
	screen.DrawImage(img, nil)

	// Render the sub-image. Only the red part should be rendered.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Rotate(math.Pi / 8)
	op.GeoM.Translate(200, 100)
	screen.DrawImage(img.SubImage(image.Rect(0, 0, 50, 50)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowTitle("Bleeding Edges")
	ebiten.SetWindowSize(320, 240)
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
