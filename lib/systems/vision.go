package systems

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	visionNgon = 12
)

var (
	// 影生成時の、マスクのベースとして使う黒画像
	blackImage *ebiten.Image
)

// VisionSystem は探索範囲エリアを表示する
func VisionSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Game.(*resources.Game)
	if gameResources.Level.VisionImage == nil {
		img := ebiten.NewImage(int(gameResources.Level.Width()), int(gameResources.Level.Height()))
		img.Fill(color.Black)
		gameResources.Level.VisionImage = img
	}

	if blackImage == nil {
		img := ebiten.NewImage(int(gameResources.Level.Width()), int(gameResources.Level.Height()))
		img.Fill(color.Black)
		blackImage = img
	}

	var pos *gc.Position
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	// 視界以外をグラデーションを入れながら塗りつぶし
	// TODO: 光源用のコンポーネントを追加したほうがよさそう
	{
		vs := visionVertices(visionNgon, pos.X, pos.Y, 160)
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Blend = ebiten.BlendSourceIn
		indices := []uint16{}
		for i := 0; i < visionNgon; i++ {
			if i < 65536 && visionNgon < 65536 { // uint16範囲チェック
				indices = append(indices, uint16(i), uint16(i+1)%uint16(visionNgon), uint16(visionNgon))
			}
		}
		gameResources.Level.VisionImage.DrawTriangles(vs, indices, blackImage, opt)
	}
	{
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(1)
		utils.SetTranslate(world, op)
		screen.DrawImage(gameResources.Level.VisionImage, op)
	}
}

func visionVertices(num int, x gc.Pixel, y gc.Pixel, r gc.Pixel) []ebiten.Vertex {
	vs := []ebiten.Vertex{}
	for i := 0; i < num; i++ {
		rate := float64(i) / float64(num)
		cr := 0.0
		cg := 0.0
		cb := 0.0
		vs = append(vs, ebiten.Vertex{
			DstX:   float32(float64(r)*math.Cos(2*math.Pi*rate)) + float32(x),
			DstY:   float32(float64(r)*math.Sin(2*math.Pi*rate)) + float32(y),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(cr),
			ColorG: float32(cg),
			ColorB: float32(cb),
			ColorA: 1,
		})
	}

	vs = append(vs, ebiten.Vertex{
		DstX:   float32(x),
		DstY:   float32(y),
		SrcX:   0,
		SrcY:   0,
		ColorR: 0,
		ColorG: 0,
		ColorB: 0,
		ColorA: 0,
	})

	return vs
}
