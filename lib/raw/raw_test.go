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
[[Items]]
Name = "リペア"
Description = "半分程度回復する"

[[Items]]
Name = "回復薬"
Description = "半分程度回復する"
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
[[Items]]
Name = "リペア"
SpriteSheetName = "field"
SpriteKey = "repair_item"
`
	raw, err := Load(str)
	assert.NoError(t, err)
	entity, err := raw.GenerateItem("リペア", gc.ItemLocationInBackpack)
	assert.NoError(t, err)
	assert.NotNil(t, entity.Name)
	assert.NotNil(t, entity.Item)
	assert.NotNil(t, entity.Description)
	assert.NotNil(t, entity.SpriteRender)
}

func TestGenerateItemWithoutSprite(t *testing.T) {
	t.Parallel()
	str := `
[[Items]]
Name = "テストアイテム"
Description = "スプライトなしアイテム"
`
	raw, err := Load(str)
	assert.NoError(t, err)

	// 現在の実装ではスプライト情報なしでも生成される（空文字列が設定される）
	entity, err := raw.GenerateItem("テストアイテム", gc.ItemLocationInBackpack)
	assert.NoError(t, err)
	assert.NotNil(t, entity.SpriteRender)
	assert.Equal(t, "", entity.SpriteRender.SpriteSheetName)
	assert.Equal(t, "", entity.SpriteRender.SpriteKey)
}

func TestGenerateMemberWithSprite(t *testing.T) {
	t.Parallel()
	str := `
[[Members]]
Name = "テストプレイヤー"
Player = true
SpriteSheetName = "field"
SpriteKey = "player"
[Members.Attributes]
Vitality = 50
Strength = 50
Sensation = 5
Dexterity = 6
Agility = 5
Defense = 0
`
	raw, err := Load(str)
	assert.NoError(t, err)
	entity, err := raw.GeneratePlayer("テストプレイヤー")
	assert.NoError(t, err)

	// 基本コンポーネントの確認
	assert.NotNil(t, entity.Name)
	assert.NotNil(t, entity.Player)

	// SpriteRenderコンポーネントの確認
	assert.NotNil(t, entity.SpriteRender, "SpriteRenderコンポーネントが設定されていない")
	assert.Equal(t, "field", entity.SpriteRender.SpriteSheetName, "SpriteSheetNameが正しくない")
	assert.Equal(t, "player", entity.SpriteRender.SpriteKey, "SpriteKeyが正しくない")
	assert.Equal(t, gc.DepthNumPlayer, entity.SpriteRender.Depth, "Depthが正しくない")
}

func TestGenerateMemberWithoutSprite(t *testing.T) {
	t.Parallel()
	str := `
[[Members]]
Name = "スプライトなしキャラ"
Player = true
[Members.Attributes]
Vitality = 50
Strength = 50
Sensation = 5
Dexterity = 6
Agility = 5
Defense = 0
`
	raw, err := Load(str)
	assert.NoError(t, err)

	// 現在の実装ではスプライト情報なしでも生成される（空文字列が設定される）
	entity, err := raw.GeneratePlayer("スプライトなしキャラ")
	assert.NoError(t, err)
	assert.NotNil(t, entity.SpriteRender)
	assert.Equal(t, "", entity.SpriteRender.SpriteSheetName)
	assert.Equal(t, "", entity.SpriteRender.SpriteKey)
}

func TestGenerateMaterialWithSprite(t *testing.T) {
	t.Parallel()
	str := `
[[Materials]]
Name = "テスト素材"
Description = "スプライト付き素材"
SpriteSheetName = "field"
SpriteKey = "field_item"
`
	raw, err := Load(str)
	assert.NoError(t, err)
	entity, err := raw.GenerateMaterial("テスト素材", 5, gc.ItemLocationInBackpack)
	assert.NoError(t, err)

	// 基本コンポーネントの確認
	assert.NotNil(t, entity.Name)
	assert.NotNil(t, entity.Material)
	assert.NotNil(t, entity.Description)

	// SpriteRenderコンポーネントの確認
	assert.NotNil(t, entity.SpriteRender, "SpriteRenderコンポーネントが設定されていない")
	assert.Equal(t, "field", entity.SpriteRender.SpriteSheetName, "SpriteSheetNameが正しくない")
	assert.Equal(t, "field_item", entity.SpriteRender.SpriteKey, "SpriteKeyが正しくない")
	assert.Equal(t, gc.DepthNumRug, entity.SpriteRender.Depth, "Depthが正しくない")
}

func TestGenerateMaterialWithoutSprite(t *testing.T) {
	t.Parallel()
	str := `
[[Materials]]
Name = "スプライトなし素材"
Description = "スプライトなし素材"
`
	raw, err := Load(str)
	assert.NoError(t, err)

	// 現在の実装ではスプライト情報なしでも生成される（空文字列が設定される）
	entity, err := raw.GenerateMaterial("スプライトなし素材", 5, gc.ItemLocationInBackpack)
	assert.NoError(t, err)
	assert.NotNil(t, entity.SpriteRender)
	assert.Equal(t, "", entity.SpriteRender.SpriteSheetName)
	assert.Equal(t, "", entity.SpriteRender.SpriteKey)
}

func TestLoadTilesFromRaw(t *testing.T) {
	t.Parallel()

	// テスト用のTOMLデータ（タイル定義を含む）
	tomlData := `
[[Tiles]]
Name = "TestFloor"
Description = "テスト用床タイル"
Walkable = true

[[Tiles]]
Name = "TestWall"
Description = "テスト用壁タイル"
Walkable = false

[[Items]]
Name = "テストアイテム"
Description = "テスト用"
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
[[Tiles]]
Name = "GenerateTestFloor"
Description = "生成テスト用床タイル"
Walkable = true

[[Tiles]]
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
[[Tiles]]
Name = "Helper1"
Description = "ヘルパー関数テスト1"
Walkable = true

[[Tiles]]
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
[[Tiles]]
Name = "EmptyLike"
Description = "空のような性質"
Walkable = false

[[Tiles]]
Name = "FloorLike"
Description = "床のような性質"
Walkable = true

[[Tiles]]
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

func TestLoadWithUnknownFields(t *testing.T) {
	t.Parallel()

	// 未知のフィールドを含むTOMLデータ
	invalidToml := `
[[Items]]
Name = "テストアイテム"
Description = "正常なアイテム"
UnknownField = "これは未知のフィールド"

[[UnknownSection]]
SomeField = "これは未知のセクション"
`

	_, err := Load(invalidToml)
	assert.Error(t, err, "未知のフィールドがあるTOMLでエラーが発生すべき")
	assert.Contains(t, err.Error(), "unknown keys found in TOML", "エラーメッセージに未知のキーについての情報が含まれるべき")
}

func TestLoadWithValidFields(t *testing.T) {
	t.Parallel()

	// 正常なTOMLデータ（既知のフィールドのみ）
	validToml := `
[[Items]]
Name = "テストアイテム"
Description = "正常なアイテム"
SpriteSheetName = "test_sheet"
SpriteKey = "test_key"

[[Tiles]]
Name = "テストタイル"
Description = "正常なタイル"
Walkable = true
`

	master, err := Load(validToml)
	assert.NoError(t, err, "正常なTOMLでエラーが発生してはいけない")
	assert.Equal(t, 1, len(master.Raws.Items), "アイテムが1つ読み込まれるべき")
	assert.Equal(t, 1, len(master.Raws.Tiles), "タイルが1つ読み込まれるべき")
}
