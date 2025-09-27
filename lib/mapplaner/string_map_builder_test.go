package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

func TestStringEntityPlanner_BasicExample(t *testing.T) {
	t.Parallel()
	// あなたの例を使ったテスト
	tileMap := []string{
		"######......",
		"#~fff#......",
		"#fffffrrrrrr",
		"#ffff#......",
		"######......",
	}

	entityMap := []string{
		"............",
		"......D.....",
		"....@.......",
		"..CT........",
		"............",
	}

	// バリデーション
	err := ValidateStringMap(tileMap, entityMap)
	if err != nil {
		t.Fatalf("マップバリデーションに失敗: %v", err)
	}

	// EntityPlanを直接生成
	plan, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err != nil {
		t.Fatalf("EntityPlan生成に失敗: %v", err)
	}

	// サイズをチェック
	expectedWidth := 12
	expectedHeight := 5
	if plan.Width != expectedWidth {
		t.Errorf("幅が期待値と違います: 期待値 %d, 実際 %d", expectedWidth, plan.Width)
	}
	if plan.Height != expectedHeight {
		t.Errorf("高さが期待値と違います: 期待値 %d, 実際 %d", expectedHeight, plan.Height)
	}

	// エンティティの総数をチェック（タイルエンティティ + その他）
	// 5行12列 = 60タイル + Props(2) = 62エンティティ
	expectedMinEntities := 60 // 最低限タイルエンティティの数
	if len(plan.Entities) < expectedMinEntities {
		t.Errorf("エンティティ数が少なすぎます: 期待値最低 %d, 実際 %d", expectedMinEntities, len(plan.Entities))
	}

	// Props配置をチェック（椅子とテーブル）
	propCount := 0
	for _, entity := range plan.Entities {
		if entity.EntityType == EntityTypeProp {
			propCount++
			switch entity.X {
			case 2:
				if entity.Y != 3 {
					t.Errorf("椅子のY座標が違います: 期待値 3, 実際 %d", entity.Y)
				}
				if entity.PropType == nil || *entity.PropType != gc.PropTypeChair {
					t.Errorf("椅子のタイプが違います: 期待値 %v, 実際 %v", gc.PropTypeChair, entity.PropType)
				}
			case 3:
				if entity.Y != 3 {
					t.Errorf("テーブルのY座標が違います: 期待値 3, 実際 %d", entity.Y)
				}
				if entity.PropType == nil || *entity.PropType != gc.PropTypeTable {
					t.Errorf("テーブルのタイプが違います: 期待値 %v, 実際 %v", gc.PropTypeTable, entity.PropType)
				}
			}
		}
	}

	expectedPropCount := 2
	if propCount != expectedPropCount {
		t.Errorf("Props数が期待値と違います: 期待値 %d, 実際 %d", expectedPropCount, propCount)
	}
}

func TestStringEntityPlanner_CustomMapping(t *testing.T) {
	t.Parallel()
	// カスタムマッピングのテストは現在スキップ
	// TODO: 将来的にカスタムマッピング機能を実装予定
	t.Skip("カスタムマッピング機能は未実装")
}

func TestStringEntityPlanner_EmptyEntityMap(t *testing.T) {
	t.Parallel()
	tileMap := []string{
		"###",
		"#f#",
		"###",
	}

	// エンティティマップなし
	plan, err := BuildEntityPlanFromStrings(tileMap, nil)
	if err != nil {
		t.Fatalf("EntityPlan生成に失敗: %v", err)
	}

	// エンティティは最低限タイルエンティティが生成される
	// 3x3 = 9タイル
	expectedMinEntities := 9
	if len(plan.Entities) < expectedMinEntities {
		t.Errorf("エンティティ数が少なすぎます: 期待値最低 %d, 実際 %d", expectedMinEntities, len(plan.Entities))
	}

	// Props以外のエンティティのみが存在することをチェック
	propCount := 0
	for _, entity := range plan.Entities {
		if entity.EntityType == EntityTypeProp {
			propCount++
		}
	}
	if propCount != 0 {
		t.Errorf("Propsが生成されるべきではありません: 実際 %d", propCount)
	}
}

func TestStringEntityPlanner_AllPropTypes(t *testing.T) {
	t.Parallel()
	tileMap := []string{
		"ffffffffff",
		"ffffffffff",
		"ffffffffff",
	}

	entityMap := []string{
		"CTBSAHMOR.",
		"..........",
		"..........",
	}

	plan, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err != nil {
		t.Fatalf("EntityPlan生成に失敗: %v", err)
	}

	// Propエンティティをカウント
	propCount := 0
	propTypes := make(map[gc.PropType]int)

	for _, entity := range plan.Entities {
		if entity.EntityType == EntityTypeProp && entity.PropType != nil {
			propCount++
			propTypes[*entity.PropType]++
		}
	}

	expectedPropCount := 9
	if propCount != expectedPropCount {
		t.Errorf("Props数が期待値と違います: 期待値 %d, 実際 %d", expectedPropCount, propCount)
	}

	// 各PropTypeが1つずつ存在することをチェック
	expectedTypes := []gc.PropType{
		gc.PropTypeChair, gc.PropTypeTable, gc.PropTypeBed,
		gc.PropTypeBookshelf, gc.PropTypeAltar, gc.PropTypeChest,
		gc.PropTypeBarrel, gc.PropTypeCrate, gc.PropTypeTorch,
	}

	for _, expectedType := range expectedTypes {
		if count, exists := propTypes[expectedType]; !exists || count != 1 {
			t.Errorf("PropType %v の数が期待値と違います: 期待値 1, 実際 %d", expectedType, count)
		}
	}
}

func TestValidateStringMap_InvalidCases(t *testing.T) {
	t.Parallel()
	// 空のタイルマップ
	err := ValidateStringMap([]string{}, []string{})
	if err == nil {
		t.Error("空のタイルマップでエラーが発生しませんでした")
	}

	// 行の長さが不一致
	tileMap := []string{
		"###",
		"##", // 短い行
		"###",
	}
	err = ValidateStringMap(tileMap, []string{})
	if err == nil {
		t.Error("行の長さが不一致でエラーが発生しませんでした")
	}

	// エンティティマップとタイルマップの行数が不一致
	tileMap = []string{
		"###",
		"###",
		"###",
	}
	entityMap := []string{
		"...",
		"...", // 1行少ない
	}
	err = ValidateStringMap(tileMap, entityMap)
	if err == nil {
		t.Error("行数不一致でエラーが発生しませんでした")
	}
}
