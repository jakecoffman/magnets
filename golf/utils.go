package golf

import "github.com/jakecoffman/cp"

type number interface {
	int | float64 | float32
}

func V[T number](x, y T) cp.Vector {
	return cp.Vector{float64(x), float64(y)}
}

func VectorFrom(a, b cp.Vector) cp.Vector {
	return b.Sub(a)
}
