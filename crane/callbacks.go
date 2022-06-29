package crane

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

func hookCrateCollision(crane *Crane) cp.CollisionBeginFunc {
	return func(arb *cp.Arbiter, space *cp.Space, data interface{}) bool {
		if crane.hookJoint == nil && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			// Get pointers to the two bodies in the collision pair and define local variables for them.
			// Their order matches the order of the collision types passed
			// to the collision handler this function was defined for
			hook, crate := arb.Bodies()

			// additions and removals can't be done in a normal callback.
			// Schedule a post step callback to do it.
			// Use the hook as the key and pass along the arbiter.
			space.AddPostStepCallback(func(space *cp.Space, b1, b2 interface{}) {
				//hook := b1.(*cp.Body)
				//crate := b2.(*cp.Body)
				pos := VectorFrom(hook.Position(), crate.Position()).Normalize().Mult(11).Add(hook.Position())
				crane.hookJoint = space.AddConstraint(cp.NewPivotJoint(hook, crate, pos))
			}, hook, crate)
		}

		return true
	}
}
