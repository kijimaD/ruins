package mapplaner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

func TestPropPlacer_BuildMeta(t *testing.T) {
	t.Parallel()
	// テスト用のマップを作成
	width, height := gc.Tile(20), gc.Tile(20)
	seed := uint64(12345)
	chain := NewBuilderChain(width, height, seed)

	// 部屋を作成してからPropPlacerをテスト
	chain.StartWith(RectRoomBuilder{})
	chain.With(NewFillAll(TileWall))
	chain.With(RoomDraw{})

	// マップを生成
	chain.Build()

	// PropPlacerをテスト
	propTypes := []gc.PropType{
		gc.PropTypeTable,
		gc.PropTypeChair,
	}
	placer := NewPropPlacer(0.5, propTypes) // 50%密度

	// 置物配置前の部屋数を確認
	roomCount := len(chain.BuildData.Rooms)
	if roomCount == 0 {
		t.Fatal("部屋が生成されていません")
	}

	// 置物を配置（実際のワールドがないので、ログ出力のみ）
	placer.BuildMeta(&chain.BuildData)

	t.Logf("生成された部屋数: %d", roomCount)
	t.Logf("テスト完了: PropPlacerが正常に実行されました")
}

func TestPropPlacer_NewPropPlacer(t *testing.T) {
	t.Parallel()
	// デフォルトの置物タイプでテスト
	placer1 := NewPropPlacer(0.3, []gc.PropType{})
	if len(placer1.PropTypes) == 0 {
		t.Error("デフォルトの置物タイプが設定されていません")
	}
	if placer1.PropDensity != 0.3 {
		t.Errorf("密度が期待値と異なります: %f", placer1.PropDensity)
	}

	// カスタムの置物タイプでテスト
	customTypes := []gc.PropType{gc.PropTypeTable}
	placer2 := NewPropPlacer(0.8, customTypes)
	if len(placer2.PropTypes) != 1 {
		t.Errorf("カスタム置物タイプ数が期待値と異なります: %d", len(placer2.PropTypes))
	}
	if placer2.PropTypes[0] != gc.PropTypeTable {
		t.Errorf("カスタム置物タイプが期待値と異なります: %s", placer2.PropTypes[0])
	}
}

func TestPropPlacer_GetValidPositions(t *testing.T) {
	t.Parallel()
	// テスト用のシンプルなマップを作成
	width, height := gc.Tile(10), gc.Tile(10)
	seed := uint64(12345)
	chain := NewBuilderChain(width, height, seed)

	// 単一の部屋を作成
	chain.StartWith(RectRoomBuilder{})
	chain.With(NewFillAll(TileWall))
	chain.With(RoomDraw{})
	chain.Build()

	if len(chain.BuildData.Rooms) == 0 {
		t.Fatal("部屋が生成されていません")
	}

	placer := NewPropPlacer(0.5, []gc.PropType{gc.PropTypeTable})
	room := chain.BuildData.Rooms[0]

	// 有効な位置を取得
	validPositions := placer.getValidPositions(&chain.BuildData, room)

	// 有効な位置が1つ以上あることを確認
	if len(validPositions) == 0 {
		t.Error("有効な位置が見つかりませんでした")
	}

	t.Logf("部屋サイズ: %dx%d", room.X2-room.X1, room.Y2-room.Y1)
	t.Logf("有効な位置数: %d", len(validPositions))

	// 各位置が実際に床タイルであることを確認
	for i, pos := range validPositions {
		idx := chain.BuildData.Level.XYTileIndex(pos.X, pos.Y)
		tile := chain.BuildData.Tiles[idx]
		if tile != TileFloor {
			t.Errorf("位置%d (%d, %d)が床タイルではありません: %v", i, pos.X, pos.Y, tile)
		}

		// 部屋の境界内であることを確認
		if pos.X <= room.X1 || pos.X >= room.X2-1 || pos.Y <= room.Y1 || pos.Y >= room.Y2-1 {
			t.Errorf("位置%d (%d, %d)が部屋の境界外です", i, pos.X, pos.Y)
		}
	}
}

func TestNewTownBuilder(t *testing.T) {
	t.Parallel()
	width, height := gc.Tile(30), gc.Tile(30)
	seed := uint64(98765)

	// 街ビルダーを作成してテスト
	chain := NewTownBuilder(width, height, seed)
	chain.Build()

	// 部屋が生成されていることを確認
	if len(chain.BuildData.Rooms) == 0 {
		t.Fatal("部屋が生成されていません")
	}

	// マップサイズが正しいことを確認
	if chain.BuildData.Level.TileWidth != width || chain.BuildData.Level.TileHeight != height {
		t.Errorf("マップサイズが期待値と異なります: %dx%d (期待値: %dx%d)",
			chain.BuildData.Level.TileWidth, chain.BuildData.Level.TileHeight, width, height)
	}

	t.Logf("生成された部屋数: %d", len(chain.BuildData.Rooms))
	t.Logf("置物付きマップビルダーのテスト完了")
}
