package save

import (
	"encoding/json"
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
	t.Parallel()
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
	stableID1Again := manager.GetStableID(entity1)
	assert.Equal(t, stableID1, stableID1Again)

	// 逆引きが正しく動作することを確認
	retrievedEntity1, exists1 := manager.GetEntity(stableID1)
	assert.True(t, exists1)
	assert.Equal(t, entity1, retrievedEntity1)

	retrievedEntity2, exists2 := manager.GetEntity(stableID2)
	assert.True(t, exists2)
	assert.Equal(t, entity2, retrievedEntity2)
}

func TestStableIDGeneration(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	world := createTestWorld()
	registry := NewComponentRegistry()

	// ワールドから自動初期化
	err := registry.InitializeFromWorld(world)
	require.NoError(t, err)

	// 型情報が正しく登録されていることを確認
	visionInfo, exists := registry.GetTypeInfoByName("AIVision")
	assert.True(t, exists)
	assert.Equal(t, "AIVision", visionInfo.Name)

	// 存在しない型
	_, exists = registry.GetTypeInfoByName("NonExistent")
	assert.False(t, exists)
}

func TestSerializationManager_SaveAndLoad(t *testing.T) {
	t.Parallel()
	// テストディレクトリを準備
	testDir := "./test_saves"
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

	// シリアライゼーションマネージャーを作成
	manager := NewSerializationManager(testDir)

	// テスト用ワールドを作成
	world := createTestWorld()

	// プレイヤーエンティティを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(10)})
	player.AddComponent(world.Components.Player, &gc.Player{})

	// NPCエンティティを作成
	npc := world.Manager.NewEntity()
	npc.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(15), Y: gc.Tile(20)})
	npc.AddComponent(world.Components.AIVision, &gc.AIVision{
		ViewDistance: gc.Pixel(160),
		TargetEntity: &player, // プレイヤーを参照
	})
	npc.AddComponent(world.Components.AIRoaming, &gc.AIRoaming{
		SubState:              gc.AIRoamingWaiting,
		StartSubStateTurn:     1,
		DurationSubStateTurns: 3,
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

	// まずプレイヤーを見つける
	newWorld.Manager.Join(newWorld.Components.Player).Visit(ecs.Visit(func(entity ecs.Entity) {
		restoredPlayer = entity
		playerFound = true

		// GridElementをチェック
		grid := newWorld.Components.GridElement.Get(entity).(*gc.GridElement)
		assert.Equal(t, gc.Tile(5), grid.X)
		assert.Equal(t, gc.Tile(10), grid.Y)

	}))

	// エンティティの検証
	entityCount := 0
	newWorld.Manager.Join(newWorld.Components.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		entityCount++

		// NPCを特定
		if entity.HasComponent(newWorld.Components.AIVision) {
			npcFound = true

			// GridElementをチェック
			grid := newWorld.Components.GridElement.Get(entity).(*gc.GridElement)
			assert.Equal(t, gc.Tile(15), grid.X)
			assert.Equal(t, gc.Tile(20), grid.Y)

			// AIVisionをチェック
			vision := newWorld.Components.AIVision.Get(entity).(*gc.AIVision)
			assert.Equal(t, gc.Pixel(160), vision.ViewDistance)

			// エンティティ参照が正しく復元されているかチェック
			if vision.TargetEntity != nil && playerFound {
				assert.Equal(t, restoredPlayer, *vision.TargetEntity)
			}

			// AIRoamingをチェック
			if entity.HasComponent(newWorld.Components.AIRoaming) {
				roaming := newWorld.Components.AIRoaming.Get(entity).(*gc.AIRoaming)
				assert.Equal(t, gc.AIRoamingWaiting, roaming.SubState)
				assert.Equal(t, 3, roaming.DurationSubStateTurns)
			}
		}
	}))

	t.Logf("Total entities restored: %d", entityCount)
	t.Logf("Player found: %v, NPC found: %v", playerFound, npcFound)

	assert.True(t, playerFound, "Player entity should be restored")
	assert.True(t, npcFound, "NPC entity should be restored")
}

func TestSerializationManager_EmptyWorld(t *testing.T) {
	t.Parallel()
	// テストディレクトリを準備
	testDir := "./test_saves_empty"
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

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
	newWorld.Manager.Join().Visit(ecs.Visit(func(_ ecs.Entity) {
		entityCount++
	}))

	assert.Equal(t, 0, entityCount)
}

func TestSerializationManager_InvalidFile(t *testing.T) {
	t.Parallel()
	// テストディレクトリを準備
	testDir := "./test_saves_invalid"
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

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

func TestValidJSONButNoChecksum(t *testing.T) {
	t.Parallel()
	// テストディレクトリを準備
	testDir := "./test_saves_valid_no_checksum"
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// 有効なJSONだがチェックサムがないセーブデータを作成
	validJSONNoChecksum := `{
		"version": "1.0.0",
		"timestamp": "2024-01-01T00:00:00Z",
		"world": {
			"entities": []
		}
	}`
	err = os.WriteFile(testDir+"/valid_no_checksum.json", []byte(validJSONNoChecksum), 0644)
	require.NoError(t, err)

	// シリアライゼーションマネージャーを作成
	manager := NewSerializationManager(testDir)
	world := createTestWorld()

	// 有効なJSONだがチェックサムなしのファイルの読み込み（失敗するはず）
	err = manager.LoadWorld(world, "valid_no_checksum")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save data validation failed")
	assert.Contains(t, err.Error(), "checksum field is missing")
}

func TestChecksumValidation(t *testing.T) {
	t.Parallel()
	world := createTestWorld()

	// テスト用のセーブディレクトリ
	tempDir := t.TempDir()
	manager := NewSerializationManager(tempDir)

	// テストエンティティを作成
	entity := world.Manager.NewEntity()
	world.Components.Name.Set(entity, &gc.Name{Name: "TestEntity"})

	// セーブ実行
	err := manager.SaveWorld(world, "test_checksum")
	require.NoError(t, err)

	// セーブファイルを読み込み
	data, err := manager.loadDataImpl("test_checksum")
	require.NoError(t, err)

	// JSONをパース
	var saveData Data
	err = json.Unmarshal(data, &saveData)
	require.NoError(t, err)

	// 正常なチェックサム検証
	err = manager.validateChecksum(&saveData)
	assert.NoError(t, err)

	// チェックサムを改ざん
	originalChecksum := saveData.Checksum
	saveData.Checksum = "invalid_checksum"

	// 改ざんされたチェックサムでの検証（失敗するはず）
	err = manager.validateChecksum(&saveData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum mismatch")

	// データを改ざん（チェックサムは元に戻す）
	saveData.Checksum = originalChecksum
	saveData.Version = "tampered_version"

	// データ改ざんでの検証（失敗するはず）
	err = manager.validateChecksum(&saveData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum mismatch")
}

func TestTamperedSaveDataLoad(t *testing.T) {
	t.Parallel()
	world := createTestWorld()

	// テスト用のセーブディレクトリ
	tempDir := t.TempDir()
	manager := NewSerializationManager(tempDir)

	// テストエンティティを作成
	entity := world.Manager.NewEntity()
	world.Components.Name.Set(entity, &gc.Name{Name: "TestEntity"})

	// セーブ実行
	err := manager.SaveWorld(world, "test_tampered")
	require.NoError(t, err)

	// セーブファイルを直接読み込み・改ざん
	data, err := manager.loadDataImpl("test_tampered")
	require.NoError(t, err)

	var saveData Data
	err = json.Unmarshal(data, &saveData)
	require.NoError(t, err)

	// データを改ざん
	saveData.Version = "hacked_version"

	// 改ざんされたデータを書き戻し
	tamperedData, err := json.MarshalIndent(saveData, "", "  ")
	require.NoError(t, err)

	err = manager.saveDataImpl("test_tampered", tamperedData)
	require.NoError(t, err)

	// 改ざんされたデータのロードを試行（失敗するはず）
	err = manager.LoadWorld(world, "test_tampered")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save data validation failed")
}

func TestDeterministicHashCalculation(t *testing.T) {
	t.Parallel()
	world := createTestWorld()

	// テスト用のセーブディレクトリ
	tempDir := t.TempDir()
	manager := NewSerializationManager(tempDir)

	// 同じテストエンティティを作成
	entity1 := world.Manager.NewEntity()
	world.Components.Name.Set(entity1, &gc.Name{Name: "TestEntity1"})
	entity2 := world.Manager.NewEntity()
	world.Components.Name.Set(entity2, &gc.Name{Name: "TestEntity2"})

	// セーブを複数回実行
	err := manager.SaveWorld(world, "test_deterministic_1")
	require.NoError(t, err)

	err = manager.SaveWorld(world, "test_deterministic_2")
	require.NoError(t, err)

	// 両方のセーブファイルを読み込み
	data1, err := manager.loadDataImpl("test_deterministic_1")
	require.NoError(t, err)

	data2, err := manager.loadDataImpl("test_deterministic_2")
	require.NoError(t, err)

	var saveData1, saveData2 Data
	err = json.Unmarshal(data1, &saveData1)
	require.NoError(t, err)
	err = json.Unmarshal(data2, &saveData2)
	require.NoError(t, err)

	// タイムスタンプを除いてチェックサムを再計算
	// 同じワールド状態であれば、同じチェックサムが生成されるはず
	saveData1.Timestamp = saveData2.Timestamp // タイムスタンプを同一にする

	checksum1 := manager.calculateChecksum(&saveData1)
	checksum2 := manager.calculateChecksum(&saveData2)

	// 決定的ハッシュなので、同じデータからは同じハッシュが生成されるはず
	assert.Equal(t, checksum1, checksum2, "同じワールド状態からは同じチェックサムが生成されるべき")
}

func TestHashConsistencyAcrossRuns(t *testing.T) {
	t.Parallel()
	world := createTestWorld()

	tempDir := t.TempDir()
	manager := NewSerializationManager(tempDir)

	// テストエンティティを作成（複数のコンポーネント付き）
	entity := world.Manager.NewEntity()
	world.Components.Name.Set(entity, &gc.Name{Name: "ConsistencyTest"})
	world.Components.GridElement.Set(entity, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(10)})

	// ワールドデータを抽出
	worldData := manager.extractWorldData(world)

	// 同じデータから複数回ハッシュを計算
	data := Data{
		Version: "1.0.0",
		World:   worldData,
	}

	hash1 := manager.calculateChecksum(&data)
	hash2 := manager.calculateChecksum(&data)
	hash3 := manager.calculateChecksum(&data)

	// 同じデータからは必ず同じハッシュが生成されるはず
	assert.Equal(t, hash1, hash2, "同一データから生成されるハッシュは一致するべき")
	assert.Equal(t, hash2, hash3, "同一データから生成されるハッシュは一致するべき")
	assert.NotEmpty(t, hash1, "ハッシュは空でないべき")
}

func TestMissingChecksumValidation(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	manager := NewSerializationManager(tempDir)

	// チェックサムなしのセーブデータを作成
	saveDataWithoutChecksum := Data{
		Version:   "1.0.0",
		Timestamp: time.Now(),
		World: WorldSaveData{
			Entities: []EntitySaveData{},
		},
		// Checksumフィールドは空
	}

	// チェックサム検証（失敗するはず）
	err := manager.validateChecksum(&saveDataWithoutChecksum)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum field is missing")
}

func TestOldSaveDataWithoutChecksum(t *testing.T) {
	t.Parallel()
	world := createTestWorld()
	tempDir := t.TempDir()
	manager := NewSerializationManager(tempDir)

	// テストエンティティを作成
	entity := world.Manager.NewEntity()
	world.Components.Name.Set(entity, &gc.Name{Name: "TestEntity"})

	// チェックサムなしの古いフォーマットのセーブデータを手動作成
	oldFormatData := map[string]interface{}{
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"world": map[string]interface{}{
			"entities": []interface{}{},
		},
		// checksumフィールドなし
	}

	// JSONにシリアライズ
	oldFormatJSON, err := json.MarshalIndent(oldFormatData, "", "  ")
	require.NoError(t, err)

	// セーブファイルとして書き込み
	err = manager.saveDataImpl("old_format_test", oldFormatJSON)
	require.NoError(t, err)

	// 古いフォーマットのロードを試行（失敗するはず）
	err = manager.LoadWorld(world, "old_format_test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save data validation failed")
	assert.Contains(t, err.Error(), "checksum field is missing")
}
