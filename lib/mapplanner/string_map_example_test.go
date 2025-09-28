package mapplanner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
)

// 文字列ベースマップビルダーの使用例
func ExampleBuildEntityPlanFromStrings() {
	// タイル配置を文字列で定義
	tileMap := []string{
		"######......",
		"#ffff#......",
		"#fffffffff..",
		"#ffff#......",
		"######......",
	}

	// エンティティ配置を文字列で定義
	entityMap := []string{
		"............",
		"............",
		"....@.......",
		"..CT........",
		"............",
	}

	// EntityPlanを生成
	plan, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err != nil {
		panic(err)
	}

	fmt.Printf("マップサイズ: %dx%d\n", plan.Width, plan.Height)
	fmt.Printf("エンティティ数: %d\n", len(plan.Entities))

	// Props数をカウント
	propCount := 0
	for _, entity := range plan.Entities {
		if entity.EntityType == EntityTypeProp {
			propCount++
		}
	}
	fmt.Printf("Props数: %d\n", propCount)

	// Output:
	// マップサイズ: 12x5
	// エンティティ数: 62
	// Props数: 2
}

// EntityPlan統合使用例（mapspawnerは別パッケージで使用）
func Example_stringEntityPlanGeneration() {
	// 小さな街のレイアウト
	tileMap := []string{
		"##########",
		"#ffffffff#",
		"#ffffffff#",
		"#ffffffff#",
		"#ffffffff#",
		"##########",
	}

	entityMap := []string{
		"..........",
		"..........",
		"..CT......",
		"..........",
		"..........",
		"..........",
	}

	// EntityPlanを生成
	plan, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err != nil {
		panic(err)
	}

	fmt.Printf("統合例:\n")
	fmt.Printf("マップサイズ: %dx%d\n", plan.Width, plan.Height)
	fmt.Printf("エンティティ数: %d\n", len(plan.Entities))

	// Props数をカウント
	propCount := 0
	for _, entity := range plan.Entities {
		if entity.EntityType == EntityTypeProp {
			propCount++
		}
	}
	fmt.Printf("Props数: %d\n", propCount)

	// Output:
	// 統合例:
	// マップサイズ: 10x6
	// エンティティ数: 62
	// Props数: 2
}

// タイル文字のマッピング一覧を表示する例
func Example_tileMappings() {
	mapping := getDefaultTileMapping()

	fmt.Println("タイル文字マッピング:")
	fmt.Printf("'#' -> 壁 (TileWall)\n")
	fmt.Printf("'f' -> 床 (TileFloor)\n")
	fmt.Printf("'r' -> 道路 (TileFloor)\n")
	fmt.Printf("' ' -> 床 (TileFloor)\n")
	fmt.Printf("'.' -> 空地・庭 (TileFloor)\n")

	fmt.Printf("\n実際のマッピング数: %d\n", len(mapping))
	fmt.Printf("\nワープポータル ('w', 'e') はエンティティとして別途処理されます\n")

	// Output:
	// タイル文字マッピング:
	// '#' -> 壁 (TileWall)
	// 'f' -> 床 (TileFloor)
	// 'r' -> 道路 (TileFloor)
	// ' ' -> 床 (TileFloor)
	// '.' -> 空地・庭 (TileFloor)
	//
	// 実際のマッピング数: 4
	//
	// ワープポータル ('w', 'e') はエンティティとして別途処理されます
}

// エンティティ文字のマッピング一覧を表示する例
func Example_entityMappings() {
	mapping := getDefaultEntityMapping()

	fmt.Println("エンティティ文字マッピング:")
	fmt.Printf("'@' -> プレイヤー (EntityTypePlayer)\n")
	fmt.Printf("'&' -> NPC (EntityTypeNPC)\n")
	fmt.Printf("'C' -> 椅子 (PropTypeChair)\n")
	fmt.Printf("'T' -> テーブル (PropTypeTable)\n")
	fmt.Printf("'S' -> 本棚 (PropTypeBookshelf)\n")
	fmt.Printf("'M' -> 樽 (PropTypeBarrel)\n")
	fmt.Printf("'R' -> 木箱 (PropTypeCrate)\n")

	fmt.Printf("\n実際のマッピング数: %d\n", len(mapping))

	// Output:
	// エンティティ文字マッピング:
	// '@' -> プレイヤー (EntityTypePlayer)
	// '&' -> NPC (EntityTypeNPC)
	// 'C' -> 椅子 (PropTypeChair)
	// 'T' -> テーブル (PropTypeTable)
	// 'S' -> 本棚 (PropTypeBookshelf)
	// 'M' -> 樽 (PropTypeBarrel)
	// 'R' -> 木箱 (PropTypeCrate)
	//
	// 実際のマッピング数: 7
}

// 複雑なマップレイアウトの例
func Example_complexLayout() {
	// より複雑な街のレイアウト
	tileMap := []string{
		"##################",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"#ffffffffffffffff#",
		"##################",
	}

	entityMap := []string{
		"..................",
		"..................",
		"..CT.......BS.....",
		"..................",
		"..................",
		"..................",
		"..................",
		"..................",
		"..................",
		"..................",
		"..................",
	}

	plan, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err != nil {
		panic(err)
	}

	// 各種統計を表示
	tileTypeCount := make(map[string]int)
	propTypeCount := make(map[gc.PropType]int)

	for _, entity := range plan.Entities {
		switch entity.EntityType {
		case EntityTypeFloor:
			tileTypeCount["床"]++
		case EntityTypeWall:
			tileTypeCount["壁"]++
		case EntityTypeProp:
			if entity.PropType != nil {
				propTypeCount[*entity.PropType]++
			}
		}
	}

	fmt.Printf("複雑なレイアウト統計:\n")
	fmt.Printf("マップサイズ: %dx%d\n", plan.Width, plan.Height)
	fmt.Printf("床タイル数: %d\n", tileTypeCount["床"])
	fmt.Printf("壁タイル数: %d\n", tileTypeCount["壁"])
	fmt.Printf("総Props数: %d\n", len(propTypeCount))

	for propType, count := range propTypeCount {
		fmt.Printf("  %s: %d個\n", propType, count)
	}
}
