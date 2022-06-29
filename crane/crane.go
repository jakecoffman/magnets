package crane

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/magnets/crane/constant"
	"math"
)

type Crane struct {
	space *cp.Space

	dollyBody *cp.Body

	// Constraint used as a servo motor to move the dolly back and forth.
	dollyServo *cp.PivotJoint

	// Constraint used as a winch motor to lift the load.
	winchServo *cp.SlideJoint

	// Temporary joint used to hold the hook to the load.
	hookJoint *cp.Constraint

	magnetImg *ebiten.Image
	hookBody  *cp.Body
}

func NewCrane(space *cp.Space) *Crane {
	dollyBody := space.AddBody(cp.NewBody(10, cp.INFINITY))
	dollyBody.SetPosition(cp.Vector{100, 100})
	space.AddShape(cp.NewBox(dollyBody, 30, 30, 0))
	// Add a groove joint for it to move back and forth on.
	space.AddConstraint(cp.NewGrooveJoint(space.StaticBody, dollyBody, V(100, 100), V(constant.ScreenWidth-100, 100), Z()))

	// Add a pivot joint to act as a servo motor controlling its position
	// By updating the anchor points of the pivot joint, you can move the dolly.
	dollyServo := space.AddConstraint(cp.NewPivotJoint(space.StaticBody, dollyBody, dollyBody.Position())).Class.(*cp.PivotJoint)
	// Max force the dolly servo can generate.
	dollyServo.SetMaxForce(10_000)
	// Max speed of the dolly servo
	dollyServo.SetMaxBias(1000)
	// You can also change the error bias to control how it slows down.
	//dollyServo.SetErrorBias(0.2)

	const (
		hookMass   = 1
		hookRadius = 10
	)
	hookBody := space.AddBody(cp.NewBody(hookMass, cp.MomentForCircle(hookMass, hookRadius, 0, Z())))
	hookBody.SetPosition(cp.Vector{200, 200})

	// This will be used to figure out when the hook touches a box.
	sensor := space.AddShape(cp.NewCircle(hookBody, hookRadius, Z()))
	//sensor.SetSensor(true)
	sensor.SetCollisionType(collisionHook)

	// By updating the max length of the joint you can make it pull up the load.
	winchServo := space.AddConstraint(cp.NewSlideJoint(dollyBody, hookBody, Z(), V(0, -10), 0, cp.INFINITY)).Class.(*cp.SlideJoint)
	winchServo.SetMaxForce(300_000)
	winchServo.SetMaxBias(600)

	return &Crane{
		space:      space,
		dollyBody:  dollyBody,
		dollyServo: dollyServo,
		winchServo: winchServo,
		hookBody:   hookBody,
		magnetImg:  Extract("magnet.png"),
	}
}

func (c *Crane) Update() {
	x, y := ebiten.CursorPosition()
	mouse := V(x, y)

	// Set the first anchor point (the one attached to the static body) of the dolly servo to the mouse's x position.
	c.dollyServo.AnchorA = V(x, 100)

	// Set the max length of the winch servo to match the mouse's height.
	c.winchServo.Max = math.Max(mouse.Y-100, 50)

	if c.hookJoint != nil && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		c.space.RemoveConstraint(c.hookJoint)
		c.hookJoint = nil
	}
}

func (c *Crane) Draw(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	w, h := c.magnetImg.Size()
	opt.GeoM.Translate(float64(-w/2), float64(-h/2))
	opt.GeoM.Scale(.4, .4)
	opt.GeoM.Rotate(c.hookBody.Angle())
	opt.GeoM.Translate(c.hookBody.Position().X, c.hookBody.Position().Y)
	screen.DrawImage(c.magnetImg, opt)
}
