package magnets

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/cpebiten"
	"github.com/jakecoffman/magnets/magnets/constant"
	"github.com/jakecoffman/magnets/magnets/hole"
	"image/color"
	"math/rand"
	"os"
	"time"
)

// Game is provided as a convenience for the examples since they all share similar logic.
type Game struct {
	// Space holds all shapes and bodies in the simulation
	Space *cp.Space

	// TicksPerSecond is the fixed physics tick rate. Set it higher if objects are going
	// through each other at the cost of higher CPU usage.
	TicksPerSecond float64

	// Accumulator shows the remaining time from the physics tick.
	Accumulator float64
	lastTime    float64

	mouseBody  *cp.Body
	mouseJoint *cp.Constraint
	touches    map[ebiten.TouchID]*touchInfo

	ball    *Ball
	hole    *Hole
	magnets []*Magnet

	levelNum int
	level    hole.Level

	playing bool

	initialClick   *cp.Vector
	lose           bool
	draggingMagnet *Magnet
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGame() *Game {
	space := cp.NewSpace()
	// slow the ball down over time
	//space.SetDamping(1)

	game := &Game{
		Space:          space,
		TicksPerSecond: 180.,
		mouseBody:      cp.NewKinematicBody(),
		touches:        map[ebiten.TouchID]*touchInfo{},

		levelNum: 0,
	}
	space.NewCollisionHandler(collisionBall, collisionHole).PreSolveFunc = ballHoleCallback(game)
	space.NewCollisionHandler(collisionBall, collisionWall).PreSolveFunc = ballWallCallback(game)
	space.NewCollisionHandler(collisionBall, collisionMagnet).PreSolveFunc = ballMagnetCallback(game)
	game.loadLevel()
	return game
}

func (g *Game) loadLevel() {
	f, err := assets.Open(fmt.Sprintf("assets/levels/level%v.json", g.levelNum+1))
	if err != nil {
		panic(err)
	}
	if err = json.NewDecoder(f).Decode(&g.level); err != nil {
		panic(err)
	}
	g.Space.EachShape(func(shape *cp.Shape) {
		g.Space.RemoveShape(shape)
	})
	g.Space.EachBody(func(body *cp.Body) {
		g.Space.RemoveBody(body)
	})

	g.playing = false
	g.Space.StaticBody = cp.NewStaticBody()
	g.ball = NewBall(g, g.level.BallStart)
	g.hole = NewHole(g, g.level.Hole)
	buildWallsAroundScreen(g.Space)
	for _, wall := range g.level.Walls {
		shape := cpebiten.AddWall(g.Space, g.Space.StaticBody, wall.A, wall.B, 1)
		shape.SetCollisionType(collisionWall)
	}
}

func buildWallsAroundScreen(space *cp.Space) {
	//container := space.AddBody(cp.NewKinematicBody())
	//container.SetPosition(cp.Vector{screenWidth / 2, screenHeight / 2})

	w, h := constant.ScreenWidth, constant.ScreenHeight

	tl := V(0-5, 0-5)
	bl := V(0-5, h+5)
	br := V(w+5, h+5)
	tr := V(w+5, 0-5)

	cpebiten.AddWall(space, space.StaticBody, tl, bl, 1)
	cpebiten.AddWall(space, space.StaticBody, bl, br, 1)
	cpebiten.AddWall(space, space.StaticBody, br, tr, 1)
	cpebiten.AddWall(space, space.StaticBody, tr, tl, 1)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.playing = false
		g.lose = false
		g.Space.RemoveBody(g.ball.body)
		g.Space.RemoveShape(g.ball.shape)
		g.ball = NewBall(g, g.level.BallStart)
		return nil
	}

	g.PhysicsTick()

	x, y := ebiten.CursorPosition()
	mouse := V(x, y)

	if inpututil.IsKeyJustPressed(ebiten.KeyF10) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		ebiten.SetVsyncEnabled(vsync)
		vsync = !vsync
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		g.magnets = append(g.magnets, NewMagnet(g.Space, V(x, y)))
	}

	g.ball.Update()
	for i := range g.magnets {
		g.magnets[i].Update()
	}
	if g.draggingMagnet != nil {
		g.draggingMagnet.shape.Body().SetPosition(mouse)
	}

	//// web stuff
	//for _, id := range inpututil.JustPressedTouchIDs() {
	//	x, y := ebiten.TouchPosition(id)
	//	touchPos := cp.Vector{float64(x), float64(y)}
	//
	//	body := cp.NewKinematicBody()
	//	body.SetPosition(touchPos)
	//	touch := &touchInfo{
	//		id:    id,
	//		body:  body,
	//		joint: handleGrab(g.Space, touchPos, body),
	//	}
	//	g.touches[id] = touch
	//}
	//for id, touch := range g.touches {
	//	if touch.joint != nil && inpututil.IsTouchJustReleased(id) {
	//		g.Space.RemoveConstraint(touch.joint)
	//		touch.joint = nil
	//		delete(g.touches, id)
	//	} else {
	//		x, y := ebiten.TouchPosition(id)
	//		touchPos := cp.Vector{float64(x), float64(y)}
	//		// calculate velocity so the object goes as fast as the touch moved
	//		newPoint := touch.body.Position().Lerp(touchPos, 0.25)
	//		touch.body.SetVelocityVector(newPoint.Sub(touch.body.Position()).Mult(60.0))
	//		touch.body.SetPosition(newPoint)
	//	}
	//}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		shape := handleGrab(g.Space, mouse)
		if shape != nil {
			if shape == g.ball.shape {
				g.initialClick = &mouse
			} else {
				for _, m := range g.magnets {
					if shape == m.shape {
						g.draggingMagnet = m
						break
					}
				}
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.initialClick != nil {
			g.playing = true
			g.lastTime = float64(time.Now().UnixNano()) / 1.e9
			dir := VectorFrom(mouse, *g.initialClick)
			g.ball.body.SetForce(dir.Mult(constant.HitStrength).Clamp(2e5))
			g.initialClick = nil
			g.ball.shape.SetFilter(NotGrabbable)
		}
		g.draggingMagnet = nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		shape := handleGrab(g.Space, mouse)
		for i, magnet := range g.magnets {
			if magnet.shape == shape {
				g.Space.RemoveBody(magnet.shape.Body())
				g.Space.RemoveShape(magnet.shape)
				g.magnets = append(g.magnets[:i], g.magnets[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (g *Game) PhysicsTick() {
	newTime := float64(time.Now().UnixNano()) / 1.e9
	frameTime := newTime - g.lastTime
	const maxUpdate = .25
	if frameTime > maxUpdate {
		frameTime = maxUpdate
	}
	g.lastTime = newTime
	g.Accumulator += frameTime

	dt := 1. / g.TicksPerSecond
	for g.Accumulator >= dt {
		for _, m := range g.magnets {
			m.pos = m.shape.Body().Position()
		}
		g.Space.Step(dt)
		g.Accumulator -= dt
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	mouse := V(x, y)

	if g.lose {
		screen.Fill(color.RGBA{235, 64, 52, 0xff})
	}

	opts := cpebiten.NewDrawOptions(screen)
	cp.DrawSpace(g.Space, opts)
	opts.Flush()

	alpha := g.Accumulator / (1. / 180.)
	for _, m := range g.magnets {
		m.Draw(screen, alpha)
	}

	if g.initialClick != nil {
		// draw a line representing the strength of the shot
		v := VectorFrom(mouse, *g.initialClick)
		v = v.Add(*g.initialClick)
		ebitenutil.DrawLine(screen, g.initialClick.X, g.initialClick.Y, v.X, v.Y, color.White)
	}

	out := fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS())
	out += fmt.Sprintf("\nClick and drag the ball to take a shot")
	out += fmt.Sprintf("\nPlace (O)mni magnet")
	out += fmt.Sprintf("\nPlace (D)irectional magnet")
	out += fmt.Sprintf("\nRight click to remove")
	out += fmt.Sprintf("\n(R)estart")
	ebitenutil.DebugPrint(screen, out)
}

func (g *Game) Layout(w int, h int) (int, int) {
	return w, h
}

var GrabbableMaskBit uint = 1 << 31

var Grabbable = cp.ShapeFilter{
	cp.NO_GROUP, GrabbableMaskBit, GrabbableMaskBit,
}
var NotGrabbable = cp.ShapeFilter{
	cp.NO_GROUP, ^GrabbableMaskBit, ^GrabbableMaskBit,
}

func handleGrab(space *cp.Space, pos cp.Vector) *cp.Shape {
	const radius = 5.0 // make it easier to grab stuff
	info := space.PointQueryNearest(pos, radius, Grabbable)
	return info.Shape
}

type touchInfo struct {
	id    ebiten.TouchID
	body  *cp.Body
	joint *cp.Constraint
}

var vsync bool
