package crane

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/magnets/crane/constant"
)

type Level struct {
	space *cp.Space

	junk []*cp.Shape
}

func NewLevel(space *cp.Space) *Level {
	w, h := constant.ScreenWidth, constant.ScreenHeight

	tl := V(0, 0)
	bl := V(0, h)
	br := V(w, h)
	tr := V(w, 0)

	addWallPoly(space, tl, tr, V(w, -1000), V(0, -1000))
	addWallPoly(space, tr, br, V(w+1000, 0), V(w+1000, h))
	addWallPoly(space, br, bl, V(0, 10000), V(w, 10000))
	addWallPoly(space, bl, tl, V(-1000, 0), V(-1000, h))

	const (
		boxMass = 10
		boxSize = 50
	)
	boxBody := space.AddBody(cp.NewBody(boxMass, cp.MomentForBox(boxMass, boxSize, 0)))
	boxBody.SetPosition(V(constant.ScreenWidth/2, constant.ScreenHeight/2))
	boxShape := space.AddShape(cp.NewBox(boxBody, boxSize, boxSize, 0))
	boxShape.SetFriction(0.7)
	boxShape.SetCollisionType(collisionCrate)

	return &Level{
		space: space,
		junk:  []*cp.Shape{boxShape},
	}
}

func (l *Level) Update() {

}

func (l *Level) Draw(screen *ebiten.Image) {

}

func addWall(space *cp.Space, a, b cp.Vector) {
	wall := space.AddShape(cp.NewSegment(space.StaticBody, a, b, 1))
	wall.SetElasticity(1)
	wall.SetFriction(1)
	wall.SetFilter(NotGrabbable)
}

func addWallPoly(space *cp.Space, a, b, c, d cp.Vector) {
	wall := space.AddShape(cp.NewPolyShape(space.StaticBody, 4, []cp.Vector{a, b, c, d}, cp.NewTransformIdentity(), 1))
	wall.SetElasticity(1)
	wall.SetFriction(1)
	wall.SetFilter(NotGrabbable)
}
