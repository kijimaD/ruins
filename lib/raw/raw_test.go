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
		RecipeIndex:       map[string]int{},
		CommandTableIndex: map[string]int{},
		DropTableIndex:    map[string]int{},
		SpriteSheetIndex:  map[string]int{},
		TileIndex:         map[string]int{},
		PropIndex:         map[string]int{},
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
	loc := gc.ItemLocationInBackpack
	entitySpec, err := raw.NewItemSpec("リペア", &loc)
	assert.NoError(t, err)
	assert.NotNil(t, entitySpec.Name)
	assert.NotNil(t, entitySpec.Item)
	assert.NotNil(t, entitySpec.Description)
	assert.NotNil(t, entitySpec.SpriteRender)
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

	// 現在の実装ではスプライト情報なしでも生成される（デフォルト値が設定される）
	loc := gc.ItemLocationInBackpack
	entitySpec, err := raw.NewItemSpec("テストアイテム", &loc)
	assert.NoError(t, err)
	assert.NotNil(t, entitySpec.SpriteRender)
	assert.Equal(t, "field", entitySpec.SpriteRender.SpriteSheetName)
	assert.Equal(t, "field_item", entitySpec.SpriteRender.SpriteKey)
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
	entitySpec, err := raw.NewPlayerSpec("テストプレイヤー")
	assert.NoError(t, err)

	// 基本コンポーネントの確認
	assert.NotNil(t, entitySpec.Name)
	assert.NotNil(t, entitySpec.Player)

	// SpriteRenderコンポーネントの確認
	assert.NotNil(t, entitySpec.SpriteRender, "SpriteRenderコンポーネントが設定されていない")
	assert.Equal(t, "field", entitySpec.SpriteRender.SpriteSheetName, "SpriteSheetNameが正しくない")
	assert.Equal(t, "player", entitySpec.SpriteRender.SpriteKey, "SpriteKeyが正しくない")
	assert.Equal(t, gc.DepthNumPlayer, entitySpec.SpriteRender.Depth, "Depthが正しくない")
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
	entitySpec, err := raw.NewPlayerSpec("スプライトなしキャラ")
	assert.NoError(t, err)
	assert.NotNil(t, entitySpec.SpriteRender)
	assert.Equal(t, "", entitySpec.SpriteRender.SpriteSheetName)
	assert.Equal(t, "", entitySpec.SpriteRender.SpriteKey)
}

func TestGenerateMaterialWithSprite(t *testing.T) {
	t.Parallel()
	str := `
[[Items]]
Name = "テスト素材"
Description = "スプライト付き素材"
SpriteSheetName = "field"
SpriteKey = "field_item"
Stackable = true
`
	raw, err := Load(str)
	assert.NoError(t, err)
	loc := gc.ItemLocationInBackpack
	entitySpec, err := raw.NewItemSpec("テスト素材", &loc)
	assert.NoError(t, err)

	// 基本コンポーネントの確認
	assert.NotNil(t, entitySpec.Name)
	assert.NotNil(t, entitySpec.Description)
	// GenerateItemはStackableコンポーネントを付与しない
	assert.Nil(t, entitySpec.Stackable)

	// SpriteRenderコンポーネントの確認
	assert.NotNil(t, entitySpec.SpriteRender, "SpriteRenderコンポーネントが設定されていない")
	assert.Equal(t, "field", entitySpec.SpriteRender.SpriteSheetName, "SpriteSheetNameが正しくない")
	assert.Equal(t, "field_item", entitySpec.SpriteRender.SpriteKey, "SpriteKeyが正しくない")
	assert.Equal(t, gc.DepthNumRug, entitySpec.SpriteRender.Depth, "Depthが正しくない")
}

func TestGenerateMaterialWithoutSprite(t *testing.T) {
	t.Parallel()
	str := `
[[Items]]
Name = "スプライトなし素材"
Description = "スプライトなし素材"
`
	raw, err := Load(str)
	assert.NoError(t, err)

	// 現在の実装ではスプライト情報なしでも生成される（デフォルト値が設定される）
	loc := gc.ItemLocationInBackpack
	entitySpec, err := raw.NewItemSpec("スプライトなし素材", &loc)
	assert.NoError(t, err)
	assert.NotNil(t, entitySpec.SpriteRender)
	assert.Equal(t, "field", entitySpec.SpriteRender.SpriteSheetName)
	assert.Equal(t, "field_item", entitySpec.SpriteRender.SpriteKey)
	// Stackable=true が設定されていないので Stackable コンポーネントは付かない
	assert.Nil(t, entitySpec.Stackable)
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

	// 床タイルの取得をテスト
	floorTile, err := master.GetTile("GenerateTestFloor")
	require.NoError(t, err, "床タイルの取得に失敗")
	assert.True(t, floorTile.Walkable)

	// 壁タイルの取得をテスト
	wallTile, err := master.GetTile("GenerateTestWall")
	require.NoError(t, err, "壁タイルの取得に失敗")
	assert.False(t, wallTile.Walkable)

	// 存在しないタイルのテスト（エラーが発生する）
	_, err = master.GetTile("NonExistent")
	assert.Error(t, err, "存在しないタイルでエラーが発生すべき")
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

	// GetTile のテスト（存在するタイル）
	tileRaw, err := master.GetTile("Helper1")
	require.NoError(t, err, "タイル取得に失敗")
	assert.True(t, tileRaw.Walkable)

	// GetTile のテスト（存在しないタイル）
	_, err = master.GetTile("NonExistent")
	assert.Error(t, err, "存在しないタイルでエラーが発生すべき")
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
		tile, err := master.GetTile(tc.name)
		require.NoError(t, err, "タイル取得に失敗: %s", tc.name)
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

	// 実際のタイル取得テスト
	floorTile, err := master.GetTile("Floor")
	require.NoError(t, err, "床タイル取得に失敗")
	assert.True(t, floorTile.Walkable)

	// 壁タイルテスト
	wallTile, err := master.GetTile("Wall")
	require.NoError(t, err, "壁タイル取得に失敗")
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
