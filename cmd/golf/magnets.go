package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/magnets/golf"
	"github.com/jakecoffman/magnets/golf/constant"
	"log"
)

func main() {
	ebiten.SetWindowSize(constant.ScreenWidth, constant.ScreenHeight)
	ebiten.SetWindowTitle("Magnets")
	if err := ebiten.RunGame(golf.NewGame()); err != nil {
		log.Fatal(err)
	}
}
