package golf

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	_ "image/png"
)

type Magnet struct {
	shape *cp.Shape
	img   *ebiten.Image

	pos cp.Vector
}

func NewMagnet(space *cp.Space, position cp.Vector) *Magnet {
	const magnetSize = 50
	body := cp.NewKinematicBody()
	body.SetPosition(position)

	circle := cp.NewCircle(body, magnetSize, cp.Vector{}).Class.(*cp.Circle)

	f, err := assets.Open("assets/magnet.png")
	if err != nil {
		panic(err)
	}
	img, _, err := ebitenutil.NewImageFromReader(f)
	if err != nil {
		panic(err)
	}

	shape := space.AddShape(circle.Shape)
	shape.SetElasticity(0)
	shape.SetFriction(0.7)
	shape.SetCollisionType(collisionMagnet)
	return &Magnet{
		shape: shape,
		img:   img,
		pos:   position,
	}
}

func (m *Magnet) Update() {}

func (m *Magnet) Draw(screen *ebiten.Image, alpha float64) {
	current := m.shape.Body().Position()
	state := current.Mult(alpha).Add(m.pos.Mult(1. - alpha))
	opt := &ebiten.DrawImageOptions{}
	w, h := m.img.Size()
	opt.GeoM.Translate(float64(-w/2), float64(-h/2))
	opt.GeoM.Translate(state.X, state.Y)
	screen.DrawImage(m.img, opt)
}
