package systems

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/camera"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	visionNgon = 12
)

var (
	// 影生成時の、マスクのベースとして使う黒画像
	blackImage *ebiten.Image
	// AI視界可視化用の円形画像
	aiVisionCircleImage *ebiten.Image
)

// VisionSystem は探索範囲エリアを表示する
func VisionSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
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
		camera.SetTranslate(world, op)
		screen.DrawImage(gameResources.Level.VisionImage, op)
	}

	// デバッグモード時のみAI視界を可視化
	cfg := config.Get()
	if cfg.Debug {
		drawAIVisionRanges(world, screen)
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

// drawAIVisionRanges はデバッグ時にAIの視界範囲を描画する
func drawAIVisionRanges(world w.World, screen *ebiten.Image) {
	// AI視界可視化用の円形画像を初期化（1回だけ）
	if aiVisionCircleImage == nil {
		// 円の直径を最大視界距離の2倍（300ピクセル）に設定
		size := 300
		aiVisionCircleImage = ebiten.NewImage(size, size)
		
		// 半透明の緑色で円を描画
		radius := float64(size / 2)
		center := float64(size / 2)
		
		// 円形の頂点を作成
		vertices := []ebiten.Vertex{}
		indices := []uint16{}
		
		// 中心点
		vertices = append(vertices, ebiten.Vertex{
			DstX:   float32(center),
			DstY:   float32(center),
			SrcX:   0,
			SrcY:   0,
			ColorR: 0.0,
			ColorG: 1.0,
			ColorB: 0.0,
			ColorA: 0.3, // 半透明
		})
		
		// 円周上の点
		circlePoints := 32
		for i := 0; i < circlePoints; i++ {
			angle := 2 * math.Pi * float64(i) / float64(circlePoints)
			x := center + radius * math.Cos(angle)
			y := center + radius * math.Sin(angle)
			
			vertices = append(vertices, ebiten.Vertex{
				DstX:   float32(x),
				DstY:   float32(y),
				SrcX:   0,
				SrcY:   0,
				ColorR: 0.0,
				ColorG: 1.0,
				ColorB: 0.0,
				ColorA: 0.3,
			})
			
			// 三角形のインデックス
			if i < circlePoints {
				indices = append(indices, 0, uint16(i+1), uint16((i+1)%circlePoints+1))
			}
		}
		
		// 円を描画
		opt := &ebiten.DrawTrianglesOptions{}
		// 1x1ピクセルの白い画像を作成
		whiteImg := ebiten.NewImage(1, 1)
		whiteImg.Fill(color.White)
		aiVisionCircleImage.DrawTriangles(vertices, indices, whiteImg, opt)
	}
	
	// 各NPCの視界を描画
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIVision,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		vision := world.Components.AIVision.Get(entity).(*gc.AIVision)
		
		// 視界範囲に応じてスケールを調整
		scale := vision.ViewDistance / 150.0 // 150は基準の視界距離
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(
			float64(position.X) - (300.0 * scale / 2.0), // 中心に配置
			float64(position.Y) - (300.0 * scale / 2.0),
		)
		camera.SetTranslate(world, op)
		
		screen.DrawImage(aiVisionCircleImage, op)
	}))
}
