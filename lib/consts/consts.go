package consts

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// ========== ウィンドウサイズ ==========

const (
	// MinGameWidth はゲームウィンドウの最小幅
	MinGameWidth = 960
	// MinGameHeight はゲームウィンドウの最小高さ
	MinGameHeight = 720
)

// ========== ゲーム定数 ==========

const (
	// TileSize はタイルの寸法
	TileSize gc.Pixel = 32
	// MapTileWidth はマップの横タイル数
	MapTileWidth = 50
	// MapTileHeight はマップの縦タイル数
	MapTileHeight = 50
	// GameClearDepth はゲームクリアとなる深度
	GameClearDepth = 100
)
