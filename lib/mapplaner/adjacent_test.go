package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestPlanData_AdjacentAnyFloor(t *testing.T) {
	t.Parallel()
	// テスト用のマップを作成
	width, height := gc.Tile(5), gc.Tile(5)
	buildData := &MetaPlan{
		Level: resources.Level{
			TileWidth:  width,
			TileHeight: height,
			Entities:   make([]ecs.Entity, int(width)*int(height)),
		},
		Tiles:     make([]Tile, int(width)*int(height)),
		Rooms:     []gc.Rect{},
		Corridors: [][]resources.TileIdx{},
	}

	// 全体を壁で埋める
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileWall
	}

	// 中央(2,2)を床にする
	centerIdx := buildData.Level.XYTileIndex(2, 2)
	buildData.Tiles[centerIdx] = TileFloor

	// テストケース1: 直交する隣接タイルは床を検出する
	upIdx := buildData.Level.XYTileIndex(1, 2)    // 上
	downIdx := buildData.Level.XYTileIndex(3, 2)  // 下
	leftIdx := buildData.Level.XYTileIndex(2, 1)  // 左
	rightIdx := buildData.Level.XYTileIndex(2, 3) // 右

	if !buildData.AdjacentAnyFloor(upIdx) {
		t.Error("上の隣接タイルで床を検出できていない")
	}
	if !buildData.AdjacentAnyFloor(downIdx) {
		t.Error("下の隣接タイルで床を検出できていない")
	}
	if !buildData.AdjacentAnyFloor(leftIdx) {
		t.Error("左の隣接タイルで床を検出できていない")
	}
	if !buildData.AdjacentAnyFloor(rightIdx) {
		t.Error("右の隣接タイルで床を検出できていない")
	}

	// テストケース2: 斜めの隣接タイルも床を検出する
	diagUpLeftIdx := buildData.Level.XYTileIndex(1, 1)    // 左上
	diagUpRightIdx := buildData.Level.XYTileIndex(1, 3)   // 右上
	diagDownLeftIdx := buildData.Level.XYTileIndex(3, 1)  // 左下
	diagDownRightIdx := buildData.Level.XYTileIndex(3, 3) // 右下

	if !buildData.AdjacentAnyFloor(diagUpLeftIdx) {
		t.Error("斜め左上の隣接タイルで床を検出できていない")
	}
	if !buildData.AdjacentAnyFloor(diagUpRightIdx) {
		t.Error("斜め右上の隣接タイルで床を検出できていない")
	}
	if !buildData.AdjacentAnyFloor(diagDownLeftIdx) {
		t.Error("斜め左下の隣接タイルで床を検出できていない")
	}
	if !buildData.AdjacentAnyFloor(diagDownRightIdx) {
		t.Error("斜め右下の隣接タイルで床を検出できていない")
	}

	// テストケース3: 離れたタイルは床を検出しない
	farIdx := buildData.Level.XYTileIndex(0, 0) // 離れた位置
	if buildData.AdjacentAnyFloor(farIdx) {
		t.Error("離れたタイルで床を誤検出している")
	}
}

func TestPlanData_AdjacentAnyFloor_WithWarpTiles(t *testing.T) {
	t.Parallel()
	// テスト用のマップを作成
	width, height := gc.Tile(5), gc.Tile(5)
	buildData := &MetaPlan{
		Level: resources.Level{
			TileWidth:  width,
			TileHeight: height,
			Entities:   make([]ecs.Entity, int(width)*int(height)),
		},
		Tiles:     make([]Tile, int(width)*int(height)),
		Rooms:     []gc.Rect{},
		Corridors: [][]resources.TileIdx{},
	}

	// 全体を壁で埋める
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileWall
	}

	// ワープポータルを配置（床 + エンティティ）
	warpNextIdx := buildData.Level.XYTileIndex(2, 2)
	warpEscapeIdx := buildData.Level.XYTileIndex(2, 3)
	buildData.Tiles[warpNextIdx] = TileFloor
	buildData.Tiles[warpEscapeIdx] = TileFloor

	// ワープポータルエンティティを追加
	buildData.WarpPortals = append(buildData.WarpPortals, WarpPortal{
		X:    2,
		Y:    2,
		Type: WarpPortalNext,
	})
	buildData.WarpPortals = append(buildData.WarpPortals, WarpPortal{
		X:    2,
		Y:    3,
		Type: WarpPortalEscape,
	})

	// 床タイルに隣接する場所から床の検出をテスト
	adjacentIdx := buildData.Level.XYTileIndex(1, 2) // (2,2)の床タイルの左隣
	if !buildData.AdjacentAnyFloor(adjacentIdx) {
		t.Error("床タイルに隣接する位置で隣接床検出が失敗")
	}

	adjacentEscapeIdx := buildData.Level.XYTileIndex(1, 3) // (2,3)の床タイルの左隣
	if !buildData.AdjacentAnyFloor(adjacentEscapeIdx) {
		t.Error("床タイルに隣接する位置で隣接床検出が失敗")
	}
}
