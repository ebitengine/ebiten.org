type Game interface {
	Update(screen *ebiten.Image) error
	Layout(outsideWidth, outsideHeight int)
		(screenWidth, screenHeight int)
}
