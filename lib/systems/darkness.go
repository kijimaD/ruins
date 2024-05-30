package systems

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils/camera"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	visionImage *ebiten.Image // 視界を表現する黒背景
	blackImage  *ebiten.Image // 影生成時の、マスクのベースとして使う黒画像
)

const (
	visionNgon = 10
)

// 周囲を暗くする
func DarknessSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Game.(*resources.Game)
	visionImage = ebiten.NewImage(gameResources.Level.Width(), gameResources.Level.Height())
	blackImage = ebiten.NewImage(gameResources.Level.Width(), gameResources.Level.Height())

	visionImage.Fill(color.Black)
	blackImage.Fill(color.Black)

	gameComponents := world.Components.Game.(*gc.Components)

	var pos *gc.Position
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.Player,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos = gameComponents.Position.Get(entity).(*gc.Position)
	}))

	// 視界以外をグラデーションを入れながら塗りつぶし
	// TODO: 光源用のコンポーネントを追加したほうがよさそう
	{
		vs := visionVertices(visionNgon, pos.X, pos.Y, 500)
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Blend = ebiten.BlendSourceIn
		indices := []uint16{}
		for i := 0; i < visionNgon; i++ {
			indices = append(indices, uint16(i), uint16(i+1)%uint16(visionNgon), uint16(visionNgon))
		}
		visionImage.DrawTriangles(vs, indices, blackImage, opt)
	}

	// 壁の影。影をキャストする用のコンポーネントを追加したほうがよさそう
	// FIXME: 見えない部分もすべてのブロックに影がついて重いのでいったんオフにしておく
	// world.Manager.Join(
	// 	gameComponents.SpriteRender,
	// 	gameComponents.BlockView,
	// 	gameComponents.BlockPass,
	// ).Visit(ecs.Visit(func(entity ecs.Entity) {
	// 	switch {
	// 	case entity.HasComponent(gameComponents.Position):
	// 		pos := gameComponents.Position.Get(entity).(*gc.Position)
	// 		sprite := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)

	// 		spriteWidth := float32(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
	// 		spriteHeight := float32(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

	// 		vector.DrawFilledRect(visionImage, float32(pos.X)-16, float32(pos.Y)-16, spriteWidth, spriteHeight+4, color.RGBA{0, 0, 0, 140}, true)
	// 		vector.DrawFilledRect(visionImage, float32(pos.X)-16, float32(pos.Y)-16, spriteWidth, spriteHeight+16, color.RGBA{0, 0, 0, 80}, true)
	// 	case entity.HasComponent(gameComponents.GridElement):
	// 		grid := gameComponents.GridElement.Get(entity).(*gc.GridElement)
	// 		sprite := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)

	// 		spriteWidth := float32(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
	// 		spriteHeight := float32(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

	// 		vector.DrawFilledRect(visionImage, float32(grid.Row)*spriteWidth, float32(grid.Col)*spriteWidth, spriteWidth, spriteHeight+4, color.RGBA{0, 0, 0, 140}, true)
	// 		vector.DrawFilledRect(visionImage, float32(grid.Row)*spriteWidth, float32(grid.Col)*spriteWidth, spriteWidth, spriteHeight+16, color.RGBA{0, 0, 0, 80}, true)
	// 	}
	// }))

	// 光源の中心付近を明るくする
	{
		vs := visionVertices(visionNgon, pos.X, pos.Y, 64)
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Blend = ebiten.BlendClear
		indices := []uint16{}
		for i := 0; i < visionNgon; i++ {
			indices = append(indices, uint16(i), uint16(i+1)%uint16(visionNgon), uint16(visionNgon))
		}
		visionImage.DrawTriangles(vs, indices, blackImage, opt)
	}

	{
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(1)
		camera.SetTranslate(world, op)
		screen.DrawImage(visionImage, op)
	}
}

func visionVertices(num int, x int, y int, r int) []ebiten.Vertex {
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
