package crane

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	"github.com/jakecoffman/magnets/crane/constant"
	"image/png"
	_ "image/png"
	"math/rand"
)

//go:embed assets
var assets embed.FS

func Extract(imgName string) *ebiten.Image {
	f, err := assets.Open("assets/" + imgName)
	if err != nil {
		panic(err)
	}
	img, _, err := ebitenutil.NewImageFromReader(f)
	if err != nil {
		panic(err)
	}
	return img
}

func AddMarchedJunk(space *cp.Space, imageName string) (*cp.Body, *cp.Shape) {
	f, err := assets.Open("assets/" + imageName)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	b := img.Bounds()
	bb := cp.BB{float64(b.Min.X), float64(b.Min.Y), float64(b.Max.X), float64(b.Max.Y)}

	sampleFunc := func(point cp.Vector) float64 {
		x := point.X
		y := point.Y
		rect := img.Bounds()

		if x < float64(rect.Min.X) || x > float64(rect.Max.X) || y < float64(rect.Min.Y) || y > float64(rect.Max.Y) {
			return 0.0
		}
		_, _, _, a := img.At(int(x), int(y)).RGBA()
		return float64(a) / 0xffff
	}

	//lineSet := MarchHard(bb, 100, 100, 0.2, PolyLineCollectSegment, sampleFunc)
	lineSet := cp.MarchSoft(bb, 300, 300, 0.5, cp.PolyLineCollectSegment, sampleFunc)

	line := lineSet.Lines[0].SimplifyCurves(.1)
	offset := V((b.Max.X-b.Min.X)/2., (b.Max.Y-b.Min.Y)/2.)
	// center the verts on origin
	for i, l := range line.Verts {
		line.Verts[i] = l.Sub(offset)
	}

	x := rand.Intn(constant.ScreenWidth-100) + 100
	body := space.AddBody(cp.NewBody(10, cp.MomentForPoly(10, len(line.Verts), line.Verts, Z(), 1)))
	body.SetPosition(V(x, constant.ScreenHeight/2))
	//body.SetPosition(V(rand.Intn(640)-320, rand.Intn(480)-240))
	shape := space.AddShape(cp.NewPolyShape(body, len(line.Verts), line.Verts, cp.NewTransformIdentity(), 0))
	shape.SetElasticity(0)
	shape.SetFriction(1)
	shape.SetCollisionType(collisionCrate)
	return body, shape
	// or use the outline of the shape with lines if you don't want a polygon
	//for i := 0; i < len(line.Verts)-1; i++ {
	//	a := line.Verts[i]
	//	b := line.Verts[i+1]
	//	AddSegment(space, body, a, b, 0)
	//}
}

func AddSegment(space *cp.Space, body *cp.Body, a, b cp.Vector, radius float64) *cp.Shape {
	// swap so we always draw the same direction horizontally
	if a.X < b.X {
		a, b = b, a
	}

	seg := cp.NewSegment(body, a, b, radius).Class.(*cp.Segment)
	shape := space.AddShape(seg.Shape)
	shape.SetElasticity(0)
	shape.SetFriction(0.7)
	shape.SetCollisionType(collisionCrate)

	return shape
}
