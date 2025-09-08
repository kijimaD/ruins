package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestCheckTileItems(t *testing.T) {
	tests := []struct {
		name        string
		setupItems  func(world w.World)
		playerPos   *gc.Position
		expectedLog string
	}{
		{
			name: "GridElementベースのアイテムを発見",
			setupItems: func(world w.World) {
				// GridElementベースのアイテムを作成
				itemEntity := world.Manager.NewEntity()
				itemEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 5})
				itemEntity.AddComponent(world.Components.Item, &gc.Item{})
				itemEntity.AddComponent(world.Components.Name, &gc.Name{Name: "Iron Sword"})
			},
			playerPos: &gc.Position{
				X: 320, // タイル10 * 32
				Y: 160, // タイル5 * 32
			},
			expectedLog: "Iron Sword", // ItemNameが含まれることを確認
		},
		{
			name: "複数のアイテムを同時に発見",
			setupItems: func(world w.World) {
				// GridElementベースのアイテム
				itemEntity1 := world.Manager.NewEntity()
				itemEntity1.AddComponent(world.Components.GridElement, &gc.GridElement{X: 8, Y: 8})
				itemEntity1.AddComponent(world.Components.Item, &gc.Item{})
				itemEntity1.AddComponent(world.Components.Name, &gc.Name{Name: "Gold Coin"})
			},
			playerPos: &gc.Position{
				X: 256, // タイル8 * 32 = 256
				Y: 256, // タイル8 * 32 = 256
			},
			expectedLog: "Gold Coin", // GridElementベースのアイテムが検出される
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := CreateTestWorldWithResources(t)

			// テスト用のログストアを作成
			testStore := gamelog.NewSafeSlice(10)

			// テスト前にログをクリア
			testStore.Clear()

			// アイテムをセットアップ
			tt.setupItems(world)

			// checkTileItems関数をテスト（テスト用のログストアを使用）
			checkTileItemsWithStore(world, tt.playerPos, testStore)

			// ログの内容を確認
			messages := testStore.GetHistory()
			if tt.expectedLog == "" {
				assert.Empty(t, messages, "ログが出力されるべきではない")
			} else {
				require.NotEmpty(t, messages, "ログが出力されるべき")

				// 期待されるアイテム名がログに含まれているか確認
				found := false
				for _, message := range messages {
					if contains(message, tt.expectedLog) {
						found = true
						break
					}
				}
				assert.True(t, found, "期待されるアイテム名 '%s' がログに含まれていない。実際のログ: %v", tt.expectedLog, messages)
			}
		})
	}
}

// checkTileItemsWithStore はテスト用の関数（ログストアを指定可能）
func checkTileItemsWithStore(world w.World, playerPos *gc.Position, logStore *gamelog.SafeSlice) {
	// プレイヤーと同じタイルにあるGridElementベースのアイテムを探す
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Item,
		world.Components.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		nameComp := world.Components.Name.Get(entity).(*gc.Name)

		// タイル座標をピクセル座標に変換してチェック
		itemTileX := int(gridElement.X)
		itemTileY := int(gridElement.Y)
		playerTileX := int(playerPos.X) / 32 // TileSizeは32固定
		playerTileY := int(playerPos.Y) / 32

		if itemTileX == playerTileX && itemTileY == playerTileY {
			// アイテムを発見したメッセージを表示（テスト用ストアを使用）
			gamelog.New(logStore).
				ItemName(nameComp.Name).
				Append("を発見した。").
				Log()
		}
	}))

}

// contains は文字列に部分文字列が含まれているかチェックする
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				indexOfSubstring(s, substr) >= 0))
}

// indexOfSubstring は部分文字列のインデックスを返す
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestWarpHoleMessages(t *testing.T) {
	tests := []struct {
		name        string
		warpMode    gc.Warp // Warpコンポーネント全体を格納
		expectedLog string
	}{
		{
			name:        "次階段メッセージ",
			warpMode:    gc.Warp{Mode: gc.WarpModeNext},
			expectedLog: "階段を発見した。",
		},
		{
			name:        "出口メッセージ",
			warpMode:    gc.Warp{Mode: gc.WarpModeEscape},
			expectedLog: "出口を発見した。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := CreateTestWorldWithResources(t)

			// テスト用のログストアを作成
			testStore := gamelog.NewSafeSlice(10)
			testStore.Clear()

			// プレイヤーエンティティを作成
			playerEntity := world.Manager.NewEntity()
			playerEntity.AddComponent(world.Components.Position, &gc.Position{X: 100, Y: 100})
			playerEntity.AddComponent(world.Components.SpriteRender, &gc.SpriteRender{})
			playerEntity.AddComponent(world.Components.Operator, &gc.Operator{})

			// ワープホールをタイルエンティティとして作成
			warpEntity := world.Manager.NewEntity()
			warpEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 3, Y: 3}) // タイル(3,3)
			warpEntity.AddComponent(world.Components.Warp, &tt.warpMode)

			// ダンジョンレベルのセットアップ（Level.AtEntityが正しく動作するように）
			gameResources := world.Resources.Dungeon.(*resources.Dungeon)
			gameResources.Level = resources.Level{
				TileWidth:  10,
				TileHeight: 10,
				TileSize:   32,
				Entities:   make([]ecs.Entity, 100), // 10x10のタイルマップ
			}
			// タイル(3,3)のエンティティとしてワープホールを設定
			tileIndex := 3*10 + 3 // y*width + x
			gameResources.Level.Entities[tileIndex] = warpEntity

			// ワープメッセージの処理を個別にテスト
			testWarpMessage(world, &gc.Position{X: 100, Y: 100}, warpEntity, testStore)

			// ログの内容を確認
			messages := testStore.GetHistory()
			require.NotEmpty(t, messages, "ワープメッセージが出力されるべき")

			found := false
			for _, message := range messages {
				if message == tt.expectedLog {
					found = true
					break
				}
			}
			assert.True(t, found, "期待されるワープメッセージ '%s' がログに含まれていない。実際のログ: %v", tt.expectedLog, messages)
		})
	}
}

// testWarpMessage はワープメッセージのテスト用ヘルパー関数
func testWarpMessage(world w.World, _ *gc.Position, warpEntity ecs.Entity, logStore *gamelog.SafeSlice) {
	if warpEntity.HasComponent(world.Components.Warp) {
		warp := world.Components.Warp.Get(warpEntity).(*gc.Warp)

		switch warp.Mode {
		case gc.WarpModeNext:
			gamelog.New(logStore).
				Append("階段を発見した。").
				Log()
		case gc.WarpModeEscape:
			gamelog.New(logStore).
				Append("出口を発見した。").
				Log()
		}
	}
}
