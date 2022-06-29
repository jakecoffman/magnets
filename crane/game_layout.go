package crane

import "github.com/jakecoffman/magnets/crane/constant"

func (g *Game) Layout(w, h int) (int, int) {
	return constant.ScreenWidth, constant.ScreenHeight
}
