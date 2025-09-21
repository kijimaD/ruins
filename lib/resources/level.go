package resources

import (
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Dungeon は冒険出発から帰還までを1セットとした情報を保持する。
// 冒険出発から帰還までは複数階層が存在し、複数階層を通しての情報を保持する必要がある。
type Dungeon struct {
	// ステート遷移発生イベント。各stateで処理する
	stateEvent StateEvent
	// 現在階のフィールド情報
	Level Level
	// 階層数
	Depth int
	// 探索済みタイルのマップ（キー: "x,y", 値: true）
	ExploredTiles map[string]bool
	// ミニマップの設定
	Minimap MinimapSettings
	// プレイヤーのタイル移動状態
	PlayerTileState PlayerTileState
}

// PlayerTileState はプレイヤーのタイル移動に関する状態を管理する
type PlayerTileState struct {
	// プレイヤーの前回のタイル位置（重複メッセージ防止用）
	LastTileX int
	LastTileY int
}

// ResetPlayerTileState はプレイヤーのタイル状態をリセットする（階層移動時に使用）
func (d *Dungeon) ResetPlayerTileState() {
	d.PlayerTileState = PlayerTileState{
		LastTileX: -1,
		LastTileY: -1,
	}
}

// Level は現在の階層
type Level struct {
	// 横のタイル数
	TileWidth gc.Tile
	// 縦のタイル数
	TileHeight gc.Tile
	// タイルエンティティ群
	Entities []ecs.Entity
	// 視界を表現する黒背景
	// 階層移動でリセットされる
	VisionImage *ebiten.Image
}

// XYTileIndex はタイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYTileIndex(tx gc.Tile, ty gc.Tile) TileIdx {
	return TileIdx(int(ty)*int(l.TileWidth) + int(tx))
}

// XYTileCoord はタイルスライスのインデックスからタイル座標を求める
func (l *Level) XYTileCoord(idx TileIdx) (gc.Pixel, gc.Pixel) {
	x := int(idx) % int(l.TileWidth)
	y := int(idx) / int(l.TileWidth)

	return gc.Pixel(x), gc.Pixel(y)
}

// AtEntity はxy座標から、該当するエンティティを求める
func (l *Level) AtEntity(x gc.Pixel, y gc.Pixel) ecs.Entity {
	tx := gc.Tile(int(x) / int(consts.TileSize))
	ty := gc.Tile(int(y) / int(consts.TileSize))
	idx := l.XYTileIndex(tx, ty)

	return l.Entities[idx]
}

// Width はステージ幅。横の全体ピクセル数
func (l *Level) Width() gc.Pixel {
	return gc.Pixel(int(l.TileWidth) * int(consts.TileSize))
}

// Height はステージ縦。縦の全体ピクセル数
func (l *Level) Height() gc.Pixel {
	return gc.Pixel(int(l.TileHeight) * int(consts.TileSize))
}

// GetStateEvent はStateEventを読み取り専用で取得する（クリアしない）
func (d *Dungeon) GetStateEvent() StateEvent {
	return d.stateEvent
}

// SetStateEvent はStateEventを設定する
func (d *Dungeon) SetStateEvent(event StateEvent) {
	d.stateEvent = event
}

// ConsumeStateEvent はStateEventを一度だけ読み取り、読み取り後にStateEventNoneで自動クリアする
func (d *Dungeon) ConsumeStateEvent() StateEvent {
	event := d.stateEvent
	d.stateEvent = StateEventNone
	return event
}

// MinimapSettings はミニマップの設定を管理する
type MinimapSettings struct {
	// ミニマップのサイズ（ピクセル単位）
	Width  int
	Height int
	// ミニマップの表示位置（画面右上に配置）
	OffsetX int
	OffsetY int
	// ミニマップのスケール（何ピクセルで1タイルを表すか）
	Scale int
}
