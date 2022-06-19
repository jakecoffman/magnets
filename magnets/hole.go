package magnets

import (
	"github.com/jakecoffman/cp"
)

type Hole struct {
	game  *Game
	body  *cp.Body
	shape *cp.Shape
}

func NewHole(game *Game, pos cp.Vector) *Hole {
	body := game.Space.AddBody(cp.NewStaticBody())
	body.SetPosition(pos)

	const holeRadius = 50
	circle := cp.NewCircle(body, holeRadius, cp.Vector{}).Class.(*cp.Circle)
	shape := game.Space.AddShape(circle.Shape)
	shape.SetFilter(NotGrabbable)
	shape.SetCollisionType(collisionHole)

	return &Hole{
		game:  game,
		shape: shape,
		body:  body,
	}
}

func (h *Hole) Update() {

}
