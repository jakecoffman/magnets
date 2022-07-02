package crane

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/magnets/crane/constant"
	"math/rand"
	"time"
)

type Level struct {
	space *cp.Space

	cars []*Car
}

func NewLevel(space *cp.Space) *Level {
	w, h := constant.ScreenWidth, constant.ScreenHeight

	tl := V(0, 0)
	bl := V(0, h)
	br := V(w, h)
	tr := V(w, 0)

	addWallPoly(space, tl, tr, V(w, -1000), V(0, -1000))
	addWallPoly(space, tr, br, V(w+1000, 0), V(w+1000, h))
	addWallPoly(space, br, bl.Add(V(100, 0)), V(0, 10000), V(w, 10000))
	addWallPoly(space, bl, tl, V(-1000, 0), V(-1000, h))

	addWall(space, V(100, constant.ScreenHeight), V(100, constant.ScreenHeight-200))

	level := &Level{
		space: space,
	}

	for i := 0; i < 10; i++ {
		level.addJunk(space)
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			space.AddPostStepCallback(func(space *cp.Space, key interface{}, data interface{}) {
				level.addJunk(space)
			}, nil, nil)
		}
	}()

	return level
}

func (l *Level) Update() {

}

func (l *Level) Draw(screen *ebiten.Image) {
	for _, car := range l.cars {
		car.Draw(screen)
	}
}

func addWall(space *cp.Space, a, b cp.Vector) {
	wall := space.AddShape(cp.NewSegment(space.StaticBody, a, b, 5))
	wall.SetElasticity(1)
	wall.SetFriction(1)
	wall.SetFilter(NotGrabbable)
}

func addWallPoly(space *cp.Space, a, b, c, d cp.Vector) {
	wall := space.AddShape(cp.NewPolyShape(space.StaticBody, 4, []cp.Vector{a, b, c, d}, cp.NewTransformIdentity(), 1))
	wall.SetElasticity(0)
	wall.SetFriction(10)
	wall.SetFilter(NotGrabbable)
}

func (l *Level) addJunk(space *cp.Space) {
	if rand.Intn(2)%2 == 0 {
		addBigBox(space)
	} else {
		body, shape := AddMarchedJunk(space, "car.png")
		car := &Car{
			body:  body,
			shape: shape,
			img:   Extract("car.png"),
		}
		l.cars = append(l.cars, car)
	}
}

func addBigBox(space *cp.Space) (*cp.Body, *cp.Shape) {
	const (
		boxMass = 5
		boxSize = 50
	)
	boxBody := space.AddBody(cp.NewBody(boxMass, cp.MomentForBox(boxMass, boxSize, 0)))
	x := rand.Intn(constant.ScreenWidth-100) + 100
	boxBody.SetPosition(V(x, constant.ScreenHeight/2))
	boxShape := space.AddShape(cp.NewBox(boxBody, boxSize, boxSize, 0))
	boxShape.SetFriction(1)
	boxShape.SetCollisionType(collisionCrate)
	return boxBody, boxShape
}
