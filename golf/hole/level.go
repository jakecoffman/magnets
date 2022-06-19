package hole

import "github.com/jakecoffman/cp"

type Line struct {
	A, B cp.Vector
}

type Level struct {
	Walls     []Line
	BallStart cp.Vector
	Hole      cp.Vector
}
