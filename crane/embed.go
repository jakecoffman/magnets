package crane

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
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
