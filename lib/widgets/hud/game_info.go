package hud

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// GameInfo はHUDの基本ゲーム情報エリア
type GameInfo struct {
	enabled bool
}

// NewGameInfo は新しいHUDGameInfoを作成する
func NewGameInfo() *GameInfo {
	return &GameInfo{
		enabled: true,
	}
}

// SetEnabled は有効/無効を設定する
func (info *GameInfo) SetEnabled(enabled bool) {
	info.enabled = enabled
}

// Update はゲーム情報エリアを更新する
func (info *GameInfo) Update(_ w.World) {
	// 現在は更新処理なし
}

// Draw はゲーム情報エリアを描画する
func (info *GameInfo) Draw(world w.World, screen *ebiten.Image) {
	if !info.enabled {
		return
	}

	// フロア情報を描画
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("floor: B%d", gameResources.Depth), 0, 200)

	// プレイヤーの速度情報を描画
	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.Operator,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("speed: %.2f", velocity.Speed), 0, 220)
	}))
}
