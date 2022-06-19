package magnets

import (
	"github.com/jakecoffman/cp"
	"log"
)

const (
	collisionBall = iota + 1
	collisionHole
	collisionWall
	collisionMagnet
)

func ballHoleCallback(game *Game) cp.CollisionPreSolveFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		if game.lose {
			return arb.Ignore()
		}
		log.Println("HOLE IN 1")
		game.playing = false
		return arb.Ignore()
	}
}

func ballWallCallback(game *Game) cp.CollisionPreSolveFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		game.lose = true
		return true
	}
}

func ballMagnetCallback(game *Game) cp.CollisionPreSolveFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		return arb.Ignore()
	}
}
