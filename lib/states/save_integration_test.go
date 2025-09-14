package states

import (
	"os"
	"path/filepath"
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/save"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestSaveLoadIntegration(t *testing.T) {
	t.Parallel()
	// テスト用ディレクトリを準備
	testDir := "./test_save_integration"
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用エンティティを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(10)})
	player.AddComponent(world.Components.Operator, &gc.Operator{})

	npc := world.Manager.NewEntity()
	npc.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(15), Y: gc.Tile(20)})
	npc.AddComponent(world.Components.AIVision, &gc.AIVision{
		ViewDistance: gc.Pixel(160),
		TargetEntity: &player,
	})
	npc.AddComponent(world.Components.AIRoaming, &gc.AIRoaming{
		SubState:              gc.AIRoamingWaiting,
		StartSubStateTurn:     1,
		DurationSubStateTurns: 3,
	})

	// セーブマネージャーを作成
	saveManager := save.NewSerializationManager(testDir)

	// セーブテスト
	err = saveManager.SaveWorld(world, "test_slot")
	require.NoError(t, err)

	// セーブファイルの存在確認
	saveFile := filepath.Join(testDir, "test_slot.json")
	_, err = os.Stat(saveFile)
	assert.NoError(t, err, "Save file should exist")

	// 新しいワールドを作成
	newWorld, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// ロードテスト
	err = saveManager.LoadWorld(newWorld, "test_slot")
	require.NoError(t, err)

	// データの検証
	playerCount := 0
	npcCount := 0

	newWorld.Manager.Join(newWorld.Components.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity.HasComponent(newWorld.Components.Operator) {
			playerCount++
		}
		if entity.HasComponent(newWorld.Components.AIVision) {
			npcCount++
		}
	}))

	assert.Equal(t, 1, playerCount, "Should have 1 player")
	assert.Equal(t, 1, npcCount, "Should have 1 NPC")
}

func TestSaveSlotInfo(t *testing.T) {
	t.Parallel()
	// テスト用ディレクトリを準備
	testDir := "./test_save_slots"
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

	// セーブマネージャーを作成
	saveManager := save.NewSerializationManager(testDir)

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// SaveMenuStateを作成してテスト
	saveMenu := &SaveMenuState{}
	saveMenu.saveManager = saveManager

	// 初期状態（セーブファイルなし）
	slots := saveMenu.getSaveSlotInfo()
	assert.Len(t, slots, 3)
	for i, slot := range slots {
		assert.Equal(t, false, slot.Exists)
		assert.Contains(t, slot.Label, "[空]")
		assert.Equal(t, "空のスロット", slot.Description)
		t.Logf("Slot %d: %s - %s", i+1, slot.Label, slot.Description)
	}

	// 1つのセーブファイルを作成
	err = saveManager.SaveWorld(world, "slot1")
	require.NoError(t, err)

	// セーブファイル作成後の状態
	slots = saveMenu.getSaveSlotInfo()
	assert.True(t, slots[0].Exists, "Slot 1 should exist")
	assert.Contains(t, slots[0].Label, "1 [")
	assert.NotContains(t, slots[0].Label, "[空]")
	assert.Contains(t, slots[0].Description, "保存日時")
	assert.False(t, slots[1].Exists, "Slot 2 should not exist")
	assert.Contains(t, slots[1].Label, "[空]")
	assert.False(t, slots[2].Exists, "Slot 3 should not exist")
	assert.Contains(t, slots[2].Label, "[空]")

	t.Logf("After save - Slot 1: %s - %s", slots[0].Label, slots[0].Description)
}
