package mapplanner

import (
	"strings"
	"testing"

	"github.com/kijimaD/ruins/lib/resources"
)

func TestStringMapPlanner_UnknownTileCharacterError(t *testing.T) {
	t.Parallel()
	// 未知のタイル文字を含むマップ
	tileMap := []string{
		"####",
		"#Xf#", // 'X'は未知のタイル文字
		"####",
	}

	entityMap := []string{
		"....",
		"....",
		"....",
	}

	_, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err == nil {
		t.Fatal("未知のタイル文字に対するエラーが期待されましたが、エラーが返されませんでした")
	}

	expectedErrorSubstring := "未知のタイル文字 'X'"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("期待されるエラーメッセージ '%s' が含まれていません。実際のエラー: %v", expectedErrorSubstring, err)
	}
}

func TestStringMapPlanner_UnknownEntityCharacterError(t *testing.T) {
	t.Parallel()
	// 未知のエンティティ文字を含むマップ
	tileMap := []string{
		"fff",
		"fff",
		"fff",
	}

	entityMap := []string{
		"...",
		".Z.", // 'Z'は未知のエンティティ文字
		"...",
	}

	_, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err == nil {
		t.Fatal("未知のエンティティ文字に対するエラーが期待されましたが、エラーが返されませんでした")
	}

	expectedErrorSubstring := "未知のエンティティ文字 'Z'"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("期待されるエラーメッセージ '%s' が含まれていません。実際のエラー: %v", expectedErrorSubstring, err)
	}
}

func TestStringMapPlanner_ValidCharactersNoError(t *testing.T) {
	t.Parallel()
	// 有効な文字のみを含むマップ
	tileMap := []string{
		"####",
		"#ff#",
		"####",
	}

	entityMap := []string{
		"....",
		".@..",
		"....",
	}

	plan, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err != nil {
		t.Fatalf("有効な文字に対してエラーが発生しました: %v", err)
	}

	if plan == nil {
		t.Fatal("EntityPlanがnilです")
	}

	// プレイヤーの開始位置が設定されているかチェック
	playerX, playerY, hasPlayer := plan.GetPlayerStartPosition()
	if !hasPlayer {
		t.Error("プレイヤーの開始位置が設定されていません")
	} else {
		if playerX != 1 || playerY != 1 {
			t.Errorf("プレイヤーの開始位置が間違っています: 期待値 (1, 1), 実際 (%d, %d)", playerX, playerY)
		}
	}
}

func TestStringMapBuilder_BuildInitialError(t *testing.T) {
	t.Parallel()
	// BuildInitialでエラーが発生することをテスト
	builder := &StringMapPlanner{
		TileMap: []string{
			"#Y#", // 'Y'は未知のタイル文字
			"#f#",
			"###",
		},
		EntityMap: []string{
			"...",
			"...",
			"...",
		},
		TileMapping:   getDefaultTileMapping(),
		EntityMapping: getDefaultEntityMapping(),
	}

	// テスト用のMetaPlanを作成
	metaPlan := &MetaPlan{
		Level:        resources.Level{TileWidth: 3, TileHeight: 3},
		RandomSource: NewRandomSource(42),
		Tiles:        make([]Tile, 9),
	}

	err := builder.BuildInitial(metaPlan)
	if err == nil {
		t.Fatal("BuildInitialでエラーが期待されましたが、エラーが返されませんでした")
	}

	expectedErrorSubstring := "未知のタイル文字 'Y'"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("期待されるエラーメッセージ '%s' が含まれていません。実際のエラー: %v", expectedErrorSubstring, err)
	}
}
