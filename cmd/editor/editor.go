package main

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/magnets/golf/constant"
	"github.com/jakecoffman/magnets/golf/hole"
	"github.com/sqweek/dialog"
	"image/color"
	"log"
	"os"
)

func main() {
	ebiten.SetWindowSize(constant.ScreenWidth, constant.ScreenHeight)
	ebiten.SetWindowTitle("Magnets")
	if err := ebiten.RunGame(NewEditor()); err != nil {
		log.Fatal(err)
	}
}

func NewEditor() *Editor {
	return &Editor{
		Level: hole.Level{
			BallStart: V(100, 100),
			Hole:      V(300, 300),
		},
	}
}

type Editor struct {
	Level      hole.Level
	isDialogUp bool

	firstClick *cp.Vector
}

func (e *Editor) Update() error {
	if !e.isDialogUp && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		e.isDialogUp = true
		go func() {
			filename, err := dialog.File().Filter("JSON", "json").Save()
			if err != nil {
				return
			}
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err = json.NewEncoder(f).Encode(e.Level); err != nil {
				panic(err)
			}
			e.isDialogUp = false
		}()
	}
	if !e.isDialogUp && inpututil.IsKeyJustPressed(ebiten.KeyL) {
		go func() {
			filename, err := dialog.File().Filter("JSON", "json").Load()
			if err != nil {
				return
			}
			f, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err = json.NewDecoder(f).Decode(&e.Level); err != nil {
				panic(err)
			}
			e.isDialogUp = false
		}()
	}

	x, y := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if e.firstClick == nil {
			e.firstClick = Vptr(x, y)
		} else {
			e.Level.Walls = append(e.Level.Walls, hole.Line{*e.firstClick, V(x, y)})
			e.firstClick = Vptr(x, y)
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if e.firstClick != nil {
			e.firstClick = nil
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		e.Level.BallStart = V(x, y)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		e.Level.Hole = V(x, y)
	}

	return nil
}

func (e *Editor) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	if e.firstClick != nil {
		ebitenutil.DrawLine(screen, e.firstClick.X, e.firstClick.Y, float64(x), float64(y), color.RGBA{255, 0, 0, 255})
	}
	for _, line := range e.Level.Walls {
		ebitenutil.DrawLine(screen, line.A.X, line.A.Y, line.B.X, line.B.Y, color.White)
	}
	ebitenutil.DrawRect(screen, e.Level.BallStart.X-5, e.Level.BallStart.Y-5, 10, 10, color.RGBA{0, 255, 0, 255})
	ebitenutil.DrawRect(screen, e.Level.Hole.X-25, e.Level.Hole.Y-25, 50, 50, color.RGBA{0, 0, 255, 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("(S)ave\n(L)oad\nPlace (B)all\nPlace (H)ole"))
}

func (e *Editor) Layout(w, h int) (int, int) {
	return w, h
}

type number interface {
	int | float64 | float32
}

func V[T number](x, y T) cp.Vector {
	return cp.Vector{float64(x), float64(y)}
}

func Vptr[T number](x, y T) *cp.Vector {
	return &cp.Vector{float64(x), float64(y)}
}
