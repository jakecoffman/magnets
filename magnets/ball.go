package magnets

import (
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/cpebiten"
	"github.com/jakecoffman/magnets/magnets/constant"
)

type Ball struct {
	shape *cp.Shape
	body  *cp.Body
}

func NewBall(game *Game, position cp.Vector) *Ball {
	const ballMass, ballRadius = 1.0, 5.0
	shape := cpebiten.AddCircle(game.Space, position, ballMass, ballRadius)
	shape.SetFilter(Grabbable)
	shape.SetElasticity(1)
	shape.SetFriction(0)
	shape.SetCollisionType(collisionBall)

	body := shape.Body()
	body.SetVelocityUpdateFunc(func(body *cp.Body, gravity cp.Vector, damping float64, dt float64) {
		if !game.playing {
			body.SetVelocity(0, 0)
			return
		}
		if body.Velocity().LengthSq() < 1 {
			body.SetVelocity(0, 0)
		}
		// get pulled towards magnets
		for _, magnet := range game.magnets {
			dir := VectorFrom(body.Position(), magnet.shape.Body().Position()).Normalize()
			//dist := body.Position().DistanceSq(magnet.shape.Body().Position())
			gravity = gravity.Add(dir.Mult(constant.GravityStrength)) //.Mult(1000 / dist))
		}
		body.UpdateVelocity(gravity, damping, dt)
	})

	return &Ball{
		body:  body,
		shape: shape,
	}
}

func (b *Ball) Update() {

}
