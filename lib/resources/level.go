package resources

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
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
	// 探索済みタイルのマップ。座標をキーとして使用
	ExploredTiles map[gc.GridElement]bool
	// ミニマップの設定
	MinimapSettings MinimapSettings
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

// Level は現在の階層
// タイル計算メソッドを提供する
// TODO: 状態として持たないほうがいいかも
type Level struct {
	// 横のタイル数
	TileWidth gc.Tile
	// 縦のタイル数
	TileHeight gc.Tile
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

// Width はステージ幅。横の全体ピクセル数
func (l *Level) Width() gc.Pixel {
	return gc.Pixel(int(l.TileWidth) * int(consts.TileSize))
}

// Height はステージ縦。縦の全体ピクセル数
func (l *Level) Height() gc.Pixel {
	return gc.Pixel(int(l.TileHeight) * int(consts.TileSize))
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
