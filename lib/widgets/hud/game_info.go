package hud

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	w "github.com/kijimaD/ruins/lib/world"
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

// Update はゲーム情報エリアを更新する
func (info *GameInfo) Update(_ w.World) {
	// 現在は更新処理なし
}

// Draw はゲーム情報エリアを描画する
func (info *GameInfo) Draw(screen *ebiten.Image, data GameInfoData) {
	if !info.enabled {
		return
	}

	// フロア情報を描画
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("floor: B%d", data.FloorNumber), 0, 200)
}
