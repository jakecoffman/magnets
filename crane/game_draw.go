package crane

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/cpebiten"
)

func (g *Game) Draw(screen *ebiten.Image) {
	opts := cpebiten.NewDrawOptions(screen)
	cp.DrawSpace(g.space, opts)
	opts.Flush()

	g.crane.Draw(screen)
	g.level.Draw(screen)

	ebitenutil.DebugPrint(screen, "Control the crane by moving the mouse. Right click to release.")
}
