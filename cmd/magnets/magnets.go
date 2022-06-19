package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/magnets/magnets"
	"github.com/jakecoffman/magnets/magnets/constant"
	"log"
)

func main() {
	ebiten.SetWindowSize(constant.ScreenWidth, constant.ScreenHeight)
	ebiten.SetWindowTitle("Magnets")
	if err := ebiten.RunGame(magnets.NewGame()); err != nil {
		log.Fatal(err)
	}
}
