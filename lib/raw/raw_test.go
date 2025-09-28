package raw

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	str := `
[[item]]
name = "リペア"
description = "半分程度回復する"

[[item]]
name = "回復薬"
description = "半分程度回復する"
`
	raw, err := Load(str)
	assert.NoError(t, err)

	expect := Master{
		Raws: Raws{
			Items: []Item{
				Item{Name: "リペア", Description: "半分程度回復する"},
				Item{Name: "回復薬", Description: "半分程度回復する"},
			},
		},
		ItemIndex: map[string]int{
			"リペア": 0,
			"回復薬": 1,
		},
		MemberIndex:       map[string]int{},
		MaterialIndex:     map[string]int{},
		RecipeIndex:       map[string]int{},
		CommandTableIndex: map[string]int{},
		DropTableIndex:    map[string]int{},
		SpriteSheetIndex:  map[string]int{},
		TileIndex:         map[string]int{},
	}
	assert.Equal(t, expect, raw)
}

func TestGenerateItem(t *testing.T) {
	t.Parallel()
	str := `
[[item]]
name = "リペア"
`
	raw, err := Load(str)
	assert.NoError(t, err)
	entity, err := raw.GenerateItem("リペア", gc.ItemLocationInBackpack)
	assert.NoError(t, err)
	assert.NotNil(t, entity.Name)
	assert.NotNil(t, entity.Item)
	assert.NotNil(t, entity.Description)
}

func TestLoadTilesFromRaw(t *testing.T) {
	t.Parallel()

	// テスト用のTOMLデータ（タイル定義を含む）
	tomlData := `
[[tile]]
Name = "TestFloor"
Description = "テスト用床タイル"
Walkable = true

[[tile]]
Name = "TestWall"
Description = "テスト用壁タイル"
Walkable = false

[[item]]
name = "テストアイテム"
description = "テスト用"
`

	master, err := Load(tomlData)
	require.NoError(t, err, "raw.goからの読み込みに失敗")

	// 基本的なタイルが定義されていることを確認
	expectedTiles := []string{"TestFloor", "TestWall"}
	for _, tileName := range expectedTiles {
		_, ok := master.TileIndex[tileName]
		assert.True(t, ok, "タイル '%s' がインデックスに見つかりません", tileName)
	}
}

func TestGenerateTileFromRaw(t *testing.T) {
	t.Parallel()

	tomlData := `
[[tile]]
Name = "GenerateTestFloor"
Description = "生成テスト用床タイル"
Walkable = true

[[tile]]
Name = "GenerateTestWall"
Description = "生成テスト用壁タイル"
Walkable = false
`

	master, err := Load(tomlData)
	require.NoError(t, err, "テストTOMLの読み込みに失敗")

	// 床タイルの生成をテスト
	floorTile := master.GenerateTile("GenerateTestFloor")
	assert.True(t, floorTile.Walkable)

	// 壁タイルの生成をテスト
	wallTile := master.GenerateTile("GenerateTestWall")
	assert.False(t, wallTile.Walkable)

	// 存在しないタイルのテスト（panicが発生する）
	assert.Panics(t, func() {
		master.GenerateTile("NonExistent")
	}, "存在しないタイルでpanicが発生すべき")
}

// TestGenerateTileSpecFromRaw - TileSpecは削除されたためこのテストは不要

func TestTileHelperFunctionsFromRaw(t *testing.T) {
	t.Parallel()

	tomlData := `
[[tile]]
Name = "Helper1"
Description = "ヘルパー関数テスト1"
Walkable = true

[[tile]]
Name = "Helper2"
Description = "ヘルパー関数テスト2"
Walkable = false
`

	master, err := Load(tomlData)
	require.NoError(t, err, "テストTOMLの読み込みに失敗")

	// GenerateTile のテスト（存在するタイル）
	tileRaw := master.GenerateTile("Helper1")
	assert.True(t, tileRaw.Walkable)

	// GenerateTile のテスト（存在しないタイル）
	assert.Panics(t, func() {
		master.GenerateTile("NonExistent")
	}, "存在しないタイルでpanicが発生すべき")
}

func TestTilePropertiesFromRaw(t *testing.T) {
	t.Parallel()

	// Walkableフィールドのテスト
	tomlData := `
[[tile]]
Name = "EmptyLike"
Description = "空のような性質"
Walkable = false

[[tile]]
Name = "FloorLike"
Description = "床のような性質"
Walkable = true

[[tile]]
Name = "WallLike"
Description = "壁のような性質"
Walkable = false
`

	master, err := Load(tomlData)
	require.NoError(t, err, "テストTOMLの読み込みに失敗")

	testCases := []struct {
		name         string
		expectedWalk bool
	}{
		{"EmptyLike", false},
		{"FloorLike", true},
		{"WallLike", false},
	}

	for _, tc := range testCases {
		tile := master.GenerateTile(tc.name)
		assert.Equal(t, tc.expectedWalk, tile.Walkable, "Walkableが期待値と一致しない: %s", tc.name)
	}
}

func TestLoadFromRealTileFile(t *testing.T) {
	t.Parallel()

	// 実際のraw.tomlファイルからタイル定義を読み込み
	master, err := LoadFromFile("metadata/entities/raw/raw.toml")
	require.NoError(t, err, "実際のraw.tomlファイルの読み込みに失敗")

	// 基本的なタイルが定義されていることを確認
	expectedTiles := []string{"Empty", "Floor", "Wall"}
	for _, tileName := range expectedTiles {
		_, ok := master.TileIndex[tileName]
		assert.True(t, ok, "タイル '%s' が実際のファイルで見つかりません", tileName)
	}

	// 実際のタイル生成テスト
	floorTile := master.GenerateTile("Floor")
	assert.True(t, floorTile.Walkable)

	// 壁タイルテスト
	wallTile := master.GenerateTile("Wall")
	assert.False(t, wallTile.Walkable)
}
