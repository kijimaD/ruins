package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// カメラを考慮した画像配置オプションをセットする
// TODO: ズーム率を追加する
func SetTranslate(world w.World, op *ebiten.DrawImageOptions) {
	gameComponents := world.Components.Game.(*gc.Components)
	// var camera gc.Camera
	var cPos *gc.Position
	world.Manager.Join(
		gameComponents.Camera,
		gameComponents.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// camera = gameComponents.Camera.Get(entity).(*gc.Camera)
		cPos = gameComponents.Position.Get(entity).(*gc.Position)
	}))

	cx, cy := float64(world.Resources.ScreenDimensions.Width/2), float64(world.Resources.ScreenDimensions.Height/2)

	// カメラ位置
	op.GeoM.Translate(float64(-cPos.X), float64(-cPos.Y))
	// 画面の中央
	op.GeoM.Translate(float64(cx), float64(cy))
	op.GeoM.Scale(1, 1)
}
