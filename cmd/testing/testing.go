package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Testing")
	if err := ebiten.RunGame(&TestGame{}); err != nil {
		log.Fatal(err)
	}
}

type TestGame struct {
}

func (t *TestGame) Update() error {
	return nil
}

func (t *TestGame) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, 0, 0, screenWidth, screenHeight, color.White)
}

func (t *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
