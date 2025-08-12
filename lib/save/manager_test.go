package save

import (
	"os"
	"testing"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// テスト用のワールドを作成
func createTestWorld() w.World {
	world, err := game.InitWorld(960, 720)
	if err != nil {
		panic(err)
	}
	return world
}

func TestStableIDManager(t *testing.T) {
	world := createTestWorld()
	manager := NewStableIDManager()

	// エンティティを作成
	entity1 := world.Manager.NewEntity()
	entity2 := world.Manager.NewEntity()

	// 安定IDを取得
	stableID1 := manager.GetStableID(entity1)
	stableID2 := manager.GetStableID(entity2)

	// IDが異なることを確認
	assert.NotEqual(t, stableID1, stableID2)

	// IDが安定していることを確認
	stableID1_again := manager.GetStableID(entity1)
	assert.Equal(t, stableID1, stableID1_again)

	// 逆引きが正しく動作することを確認
	retrievedEntity1, exists1 := manager.GetEntity(stableID1)
	assert.True(t, exists1)
	assert.Equal(t, entity1, retrievedEntity1)

	retrievedEntity2, exists2 := manager.GetEntity(stableID2)
	assert.True(t, exists2)
	assert.Equal(t, entity2, retrievedEntity2)
}

func TestStableIDGeneration(t *testing.T) {
	world := createTestWorld()
	manager := NewStableIDManager()

	// エンティティを作成
	entity1 := world.Manager.NewEntity()
	stableID1 := manager.GetStableID(entity1)

	// エンティティを削除
	manager.UnregisterEntity(entity1)

	// 新しいエンティティを作成
	entity2 := world.Manager.NewEntity()
	stableID2 := manager.GetStableID(entity2)

	// インデックスは再利用されるが、世代が異なることを確認
	if stableID1.Index == stableID2.Index {
		assert.NotEqual(t, stableID1.Generation, stableID2.Generation)
	}

	// 古いIDは無効になることを確認
	assert.False(t, manager.IsValid(stableID1))
	assert.True(t, manager.IsValid(stableID2))
}

func TestComponentRegistry(t *testing.T) {
	world := createTestWorld()
	registry := NewComponentRegistry()

	// ワールドから自動初期化
	err := registry.InitializeFromWorld(world)
	require.NoError(t, err)

	// 型情報が正しく登録されていることを確認
	positionInfo, exists := registry.GetTypeInfoByName("Position")
	assert.True(t, exists)
	assert.Equal(t, "Position", positionInfo.Name)

	velocityInfo, exists := registry.GetTypeInfoByName("Velocity")
	assert.True(t, exists)
	assert.Equal(t, "Velocity", velocityInfo.Name)

	// 存在しない型
	_, exists = registry.GetTypeInfoByName("NonExistent")
	assert.False(t, exists)
}

func TestSerializationManager_SaveAndLoad(t *testing.T) {
	// テストディレクトリを準備
	testDir := "./test_saves"
	defer os.RemoveAll(testDir)

	// シリアライゼーションマネージャーを作成
	manager := NewSerializationManager(testDir)

	// テスト用ワールドを作成
	world := createTestWorld()

	// プレイヤーエンティティを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Position, &gc.Position{X: gc.Pixel(100), Y: gc.Pixel(200)})
	player.AddComponent(world.Components.Velocity, &gc.Velocity{
		Angle:        45.0,
		Speed:        2.5,
		MaxSpeed:     5.0,
		ThrottleMode: gc.ThrottleModeFront,
	})
	player.AddComponent(world.Components.Operator, &gc.Operator{})

	// NPCエンティティを作成
	npc := world.Manager.NewEntity()
	npc.AddComponent(world.Components.Position, &gc.Position{X: gc.Pixel(300), Y: gc.Pixel(400)})
	npc.AddComponent(world.Components.AIVision, &gc.AIVision{
		ViewDistance: gc.Pixel(160),
		TargetEntity: &player, // プレイヤーを参照
	})
	npc.AddComponent(world.Components.AIRoaming, &gc.AIRoaming{
		SubState:         gc.AIRoamingWaiting,
		StartSubState:    time.Now(),
		DurationSubState: time.Second * 3,
	})

	// 保存
	err := manager.SaveWorld(world, "test_slot")
	require.NoError(t, err)

	// 新しいワールドを作成
	newWorld := createTestWorld()

	// 読み込み
	err = manager.LoadWorld(newWorld, "test_slot")
	require.NoError(t, err)

	// データが正しく復元されているかチェック
	var restoredPlayer ecs.Entity
	var playerFound, npcFound bool

	// エンティティの検証
	entityCount := 0
	newWorld.Manager.Join(newWorld.Components.Position).Visit(ecs.Visit(func(entity ecs.Entity) {
		entityCount++

		// プレイヤーを特定
		if entity.HasComponent(newWorld.Components.Operator) {
			restoredPlayer = entity
			playerFound = true

			// Positionをチェック
			pos := newWorld.Components.Position.Get(entity).(*gc.Position)
			assert.Equal(t, gc.Pixel(100), pos.X)
			assert.Equal(t, gc.Pixel(200), pos.Y)

			// Velocityをチェック
			if entity.HasComponent(newWorld.Components.Velocity) {
				vel := newWorld.Components.Velocity.Get(entity).(*gc.Velocity)
				assert.Equal(t, 45.0, vel.Angle)
				assert.Equal(t, 2.5, vel.Speed)
				assert.Equal(t, 5.0, vel.MaxSpeed)
				assert.Equal(t, gc.ThrottleModeFront, vel.ThrottleMode)
			}
		}

		// NPCを特定
		if entity.HasComponent(newWorld.Components.AIVision) {
			npcFound = true

			// Positionをチェック
			pos := newWorld.Components.Position.Get(entity).(*gc.Position)
			assert.Equal(t, gc.Pixel(300), pos.X)
			assert.Equal(t, gc.Pixel(400), pos.Y)

			// AIVisionをチェック
			vision := newWorld.Components.AIVision.Get(entity).(*gc.AIVision)
			assert.Equal(t, gc.Pixel(160), vision.ViewDistance)

			// エンティティ参照が正しく復元されているかチェック
			if vision.TargetEntity != nil {
				assert.Equal(t, restoredPlayer, *vision.TargetEntity)
			}

			// AIRoamingをチェック
			if entity.HasComponent(newWorld.Components.AIRoaming) {
				roaming := newWorld.Components.AIRoaming.Get(entity).(*gc.AIRoaming)
				assert.Equal(t, gc.AIRoamingWaiting, roaming.SubState)
				assert.Equal(t, time.Second*3, roaming.DurationSubState)
			}
		}
	}))

	t.Logf("Total entities restored: %d", entityCount)
	t.Logf("Player found: %v, NPC found: %v", playerFound, npcFound)

	assert.True(t, playerFound, "Player entity should be restored")
	assert.True(t, npcFound, "NPC entity should be restored")
}

func TestSerializationManager_EmptyWorld(t *testing.T) {
	// テストディレクトリを準備
	testDir := "./test_saves_empty"
	defer os.RemoveAll(testDir)

	// シリアライゼーションマネージャーを作成
	manager := NewSerializationManager(testDir)

	// 空のワールドを作成
	world := createTestWorld()

	// 保存
	err := manager.SaveWorld(world, "empty_slot")
	require.NoError(t, err)

	// 読み込み
	newWorld := createTestWorld()
	err = manager.LoadWorld(newWorld, "empty_slot")
	require.NoError(t, err)

	// エンティティが存在しないことを確認
	entityCount := 0
	newWorld.Manager.Join().Visit(ecs.Visit(func(entity ecs.Entity) {
		entityCount++
	}))

	assert.Equal(t, 0, entityCount)
}

func TestSerializationManager_InvalidFile(t *testing.T) {
	// テストディレクトリを準備
	testDir := "./test_saves_invalid"
	defer os.RemoveAll(testDir)

	// 無効なJSONファイルを作成
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	invalidJSON := `{"invalid": json}`
	err = os.WriteFile(testDir+"/invalid_slot.json", []byte(invalidJSON), 0644)
	require.NoError(t, err)

	// シリアライゼーションマネージャーを作成
	manager := NewSerializationManager(testDir)
	world := createTestWorld()

	// 無効なファイルの読み込みでエラーが発生することを確認
	err = manager.LoadWorld(world, "invalid_slot")
	assert.Error(t, err)
}
