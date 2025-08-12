package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kijimaD/ruins/lib/camera"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	// AI視界可視化用の円形画像
	aiVisionCircleImage *ebiten.Image
)

// HUDSystem はゲームの HUD 情報を描画する
func HUDSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("floor: B%d", gameResources.Depth), 0, 200)

	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.Operator,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("speed: %.2f", velocity.Speed), 0, 220)
	}))

	// デバッグモード時のみAI情報表示
	cfg := config.Get()
	if cfg.Debug {
		drawAIVisionRanges(world, screen)
		drawAIStates(world, screen)
	}
}

// drawAIStates はAIエンティティのステートをスプライトの近くに表示する
func drawAIStates(world w.World, screen *ebiten.Image) {
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIMoveFSM,
		world.Components.AIRoaming,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		roaming := world.Components.AIRoaming.Get(entity).(*gc.AIRoaming)

		// AIの現在の状態を判定
		var stateText string
		if entity.HasComponent(world.Components.AIChasing) {
			stateText = "CHASING"
		} else {
			switch roaming.SubState {
			case gc.AIRoamingWaiting:
				stateText = "WAITING"
			case gc.AIRoamingDriving:
				stateText = "ROAMING"
			case gc.AIRoamingChasing:
				stateText = "CHASING"
			default:
				stateText = "UNKNOWN"
			}
		}

		// スプライトの上にテキストを表示
		// カメラのオフセットを考慮して座標を計算
		var cameraPos gc.Position
		world.Manager.Join(
			world.Components.Camera,
			world.Components.Position,
		).Visit(ecs.Visit(func(camEntity ecs.Entity) {
			cameraPos = *world.Components.Position.Get(camEntity).(*gc.Position)
		}))

		// 画面座標に変換
		screenX := float64(position.X-cameraPos.X) + float64(world.Resources.ScreenDimensions.Width)/2
		screenY := float64(position.Y-cameraPos.Y) + float64(world.Resources.ScreenDimensions.Height)/2

		// スプライトの上20ピクセルの位置にテキスト表示
		ebitenutil.DebugPrintAt(screen, stateText, int(screenX)-20, int(screenY)-30)
	}))
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
			x := center + radius*math.Cos(angle)
			y := center + radius*math.Sin(angle)

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
		scale := vision.ViewDistance / 300.0 // 300は基準の視界距離

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(
			float64(position.X)-(300.0*scale/2.0), // 中心に配置
			float64(position.Y)-(300.0*scale/2.0),
		)
		camera.SetTranslate(world, op)

		screen.DrawImage(aiVisionCircleImage, op)
	}))
}
