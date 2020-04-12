package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xff, 0, 0, 0xff})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Fill")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
