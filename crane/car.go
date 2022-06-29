package crane

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

type Car struct {
	body  *cp.Body
	shape *cp.Shape
	img   *ebiten.Image
}

func (c *Car) Draw(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	w, h := c.img.Size()
	opt.GeoM.Translate(float64(-w/2), float64(-h/2))
	opt.GeoM.Scale(1.2, 1.2)
	opt.GeoM.Rotate(c.body.Angle())
	opt.GeoM.Translate(c.body.Position().X, c.body.Position().Y)
	screen.DrawImage(c.img, opt)
}
