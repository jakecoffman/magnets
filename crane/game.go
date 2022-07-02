package crane

import (
	"github.com/jakecoffman/cp"
)

type Game struct {
	space *cp.Space

	level *Level
	crane *Crane
}

func NewGame() *Game {
	space := cp.NewSpace()
	space.Iterations = 30
	space.SetGravity(cp.Vector{0, 1000})
	space.SetDamping(0.4)

	level := NewLevel(space)
	crane := NewCrane(space)
	handler := space.NewCollisionHandler(collisionHook, collisionCrate)
	handler.BeginFunc = hookCrateCollision(crane)

	return &Game{
		space: space,
		crane: crane,
		level: level,
	}
}

const (
	collisionHook  = 1
	collisionCrate = 2
)

var GrabbableMaskBit uint = 1 << 31

var Grabbable = cp.ShapeFilter{
	cp.NO_GROUP, GrabbableMaskBit, GrabbableMaskBit,
}
var NotGrabbable = cp.ShapeFilter{
	cp.NO_GROUP, ^GrabbableMaskBit, ^GrabbableMaskBit,
}
