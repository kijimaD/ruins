// Package mapplanner の文字列ベースマップビルダー
// 文字列で直接タイルとエンティティ配置を指定できるシステム
package mapplanner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
)

// StringMapPlanner は文字列ベースでマップを生成するプランナー
type StringMapPlanner struct {
	TileMap       []string                // タイル配置を表す文字列のスライス
	EntityMap     []string                // エンティティ配置を表す文字列のスライス
	TileMapping   map[rune]Tile           // 文字からTileへのマッピング
	EntityMapping map[rune]EntityTemplate // 文字からEntityTemplateへのマッピング
}

// EntityTemplate はエンティティ生成計画
type EntityTemplate struct {
	EntityType EntityType
	Data       interface{} // PropType等の追加データ
}

// BuildInitial は文字列定義からマップ構造を初期化する
func (b StringMapPlanner) BuildInitial(buildData *MetaPlan) {
	if len(b.TileMap) == 0 {
		return
	}

	// マップサイズを文字列から自動検出
	height := len(b.TileMap)
	width := 0
	for _, row := range b.TileMap {
		if len(row) > width {
			width = len(row)
		}
	}

	// 既存のBuildDataサイズを保持（引数で指定されたサイズを優先）
	existingWidth := int(buildData.Level.TileWidth)
	existingHeight := int(buildData.Level.TileHeight)

	// 文字列サイズと指定サイズが異なる場合は指定サイズを使用
	if existingWidth > 0 && existingHeight > 0 {
		width = existingWidth
		height = existingHeight
	} else {
		// BuildDataのサイズを文字列サイズで調整
		buildData.Level.TileWidth = gc.Tile(width)
		buildData.Level.TileHeight = gc.Tile(height)
	}

	buildData.Tiles = make([]Tile, width*height)

	// タイルマップを解析してタイルを配置
	for y, row := range b.TileMap {
		for x, char := range row {
			if x >= width || y >= height {
				continue
			}

			tileIdx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			if tileType, exists := b.TileMapping[char]; exists {
				buildData.Tiles[tileIdx] = tileType
			} else {
				// デフォルトは床タイル
				buildData.Tiles[tileIdx] = TileFloor
			}
		}
	}

	// エンティティマップが指定されている場合はエンティティを配置
	if len(b.EntityMap) > 0 {
		b.parseEntities(buildData, width, height)
	}
}

// parseEntities はエンティティマップを解析してエンティティを配置
func (b StringMapPlanner) parseEntities(plannerMap *MetaPlan, width, height int) {
	// エンティティマップから特殊なタイル（ワープホールなど）を配置
	for y := 0; y < len(b.EntityMap) && y < height; y++ {
		for x := 0; x < len(b.EntityMap[y]) && x < width; x++ {
			char := rune(b.EntityMap[y][x])

			// ワープホールなどのタイルとして扱うべきエンティティを処理
			switch char {
			case 'w':
				// ワープホール（次のレベル）をタイルとして配置
				tileIdx := plannerMap.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				plannerMap.Tiles[tileIdx] = TileWarpNext
			case 'e':
				// エスケープワープホールをタイルとして配置
				tileIdx := plannerMap.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				plannerMap.Tiles[tileIdx] = TileWarpEscape
			}

			// TODO: 他のエンティティ（NPC、アイテムなど）の処理は
			// MetaPlanにEntitiesフィールドを追加後に実装
		}
	}
}

// NewStringMapPlanner は文字列ベースマッププランナーを作成する
func NewStringMapPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	planner := &StringMapPlanner{
		TileMapping:   getDefaultTileMapping(),
		EntityMapping: getDefaultEntityMapping(),
	}

	chain := NewPlannerChain(width, height, seed)
	chain.StartWith(planner)
	return chain
}

// NewStringMapPlannerWithMaps は文字列マップを指定してプランナーを作成する
func NewStringMapPlannerWithMaps(tileMap, entityMap []string, seed uint64) *PlannerChain {
	planner := &StringMapPlanner{
		TileMap:       tileMap,
		EntityMap:     entityMap,
		TileMapping:   getDefaultTileMapping(),
		EntityMapping: getDefaultEntityMapping(),
	}

	// サイズを文字列から自動検出
	height := len(tileMap)
	width := 0
	for _, row := range tileMap {
		if len(row) > width {
			width = len(row)
		}
	}

	chain := NewPlannerChain(gc.Tile(width), gc.Tile(height), seed)
	chain.StartWith(planner)
	return chain
}

// getDefaultTileMapping はデフォルトのタイル文字マッピングを返す
func getDefaultTileMapping() map[rune]Tile {
	return map[rune]Tile{
		'#': TileWall,       // 壁
		'f': TileFloor,      // 床（建物内・住宅）
		'r': TileFloor,      // 道路（実際は床として処理、視覚的に区別される）
		'w': TileWarpNext,   // 進行ワープポータル
		'e': TileWarpEscape, // 帰還ワープポータル
		' ': TileFloor,      // 空白は床として扱う
		'.': TileFloor,      // 空地・庭
	}
}

// getDefaultEntityMapping はデフォルトのエンティティ文字マッピングを返す
func getDefaultEntityMapping() map[rune]EntityTemplate {
	return map[rune]EntityTemplate{
		'@': {EntityType: EntityTypePlayer},
		'D': {EntityType: EntityTypeDoor},
		'C': {EntityType: EntityTypeProp, Data: gc.PropTypeChair},     // 椅子
		'T': {EntityType: EntityTypeProp, Data: gc.PropTypeTable},     // テーブル
		'B': {EntityType: EntityTypeProp, Data: gc.PropTypeBed},       // ベッド
		'S': {EntityType: EntityTypeProp, Data: gc.PropTypeBookshelf}, // 本棚（店舗の棚）
		'A': {EntityType: EntityTypeProp, Data: gc.PropTypeAltar},     // 祭壇（農家・家具）
		'H': {EntityType: EntityTypeProp, Data: gc.PropTypeChest},     // 宝箱（住宅の壁）
		'O': {EntityType: EntityTypeProp, Data: gc.PropTypeTorch},     // 松明
		'M': {EntityType: EntityTypeProp, Data: gc.PropTypeBarrel},    // 樽
		'R': {EntityType: EntityTypeProp, Data: gc.PropTypeCrate},     // 木箱
		'&': {EntityType: EntityTypeNPC, Data: "街の住人"},                // NPC（街の住人）
		// 空のエンティティは明示的に含めない（何も配置されない）
	}
}

// SetTileMapping はカスタムタイルマッピングを設定する
func (b *StringMapPlanner) SetTileMapping(mapping map[rune]Tile) *StringMapPlanner {
	b.TileMapping = mapping
	return b
}

// SetEntityMapping はカスタムエンティティマッピングを設定する
func (b *StringMapPlanner) SetEntityMapping(mapping map[rune]EntityTemplate) *StringMapPlanner {
	b.EntityMapping = mapping
	return b
}

// AddTileMapping は追加のタイルマッピングを設定する
func (b *StringMapPlanner) AddTileMapping(char rune, tileType Tile) *StringMapPlanner {
	if b.TileMapping == nil {
		b.TileMapping = getDefaultTileMapping()
	}
	b.TileMapping[char] = tileType
	return b
}

// AddEntityMapping は追加のエンティティマッピングを設定する
func (b *StringMapPlanner) AddEntityMapping(char rune, entityType EntityType, data interface{}) *StringMapPlanner {
	if b.EntityMapping == nil {
		b.EntityMapping = getDefaultEntityMapping()
	}
	b.EntityMapping[char] = EntityTemplate{EntityType: entityType, Data: data}
	return b
}

// BuildEntityPlanFromStrings は文字列マップから直接EntityPlanを生成する
func BuildEntityPlanFromStrings(tileMap, entityMap []string) (*EntityPlan, error) {
	if err := ValidateStringMap(tileMap, entityMap); err != nil {
		return nil, err
	}

	// マップサイズを文字列から自動検出
	height := len(tileMap)
	width := 0
	for _, row := range tileMap {
		if len(row) > width {
			width = len(row)
		}
	}

	plan := NewEntityPlan(width, height)

	// デフォルトマッピングを取得
	tileMapping := getDefaultTileMapping()
	entityMapping := getDefaultEntityMapping()

	// タイルマップを解析してタイルエンティティを生成
	for y, row := range tileMap {
		for x, char := range row {
			if x >= width || y >= height {
				continue
			}

			var tileType Tile
			if mappedType, exists := tileMapping[char]; exists {
				tileType = mappedType
			} else {
				tileType = TileFloor // デフォルト
			}

			// タイルタイプに応じてエンティティを追加
			switch tileType {
			case TileFloor:
				plan.AddFloor(x, y)
			case TileWall:
				plan.AddWall(x, y, 0) // スプライト番号は0をデフォルト
			case TileWarpNext:
				plan.AddWarpNext(x, y)
			case TileWarpEscape:
				plan.AddWarpEscape(x, y)
			default:
				// その他のタイプは床として扱う
				plan.AddFloor(x, y)
			}
		}
	}

	// エンティティマップを解析してエンティティを追加
	if len(entityMap) > 0 {
		for y, row := range entityMap {
			if y >= height {
				break
			}
			for x, char := range row {
				if x >= width || char == ' ' || char == '.' {
					continue
				}

				// ワープホールを直接処理（entityMappingを使わない）
				if char == 'w' {
					plan.AddWarpNext(x, y)
					continue
				}
				if char == 'e' {
					plan.AddWarpEscape(x, y)
					continue
				}

				if entityPlan, exists := entityMapping[char]; exists {
					switch entityPlan.EntityType {
					case EntityTypeProp:
						if propType, ok := entityPlan.Data.(gc.PropType); ok {
							plan.AddProp(x, y, propType)
						}
					case EntityTypeNPC:
						if npcType, ok := entityPlan.Data.(string); ok {
							plan.AddNPC(x, y, npcType)
						} else {
							plan.AddNPC(x, y, "デフォルトNPC")
						}
					case EntityTypeItem:
						if itemName, ok := entityPlan.Data.(string); ok {
							plan.AddItem(x, y, itemName)
						} else {
							plan.AddItem(x, y, "デフォルトアイテム")
						}
					case EntityTypePlayer:
						// プレイヤーの開始位置を設定
						plan.SetPlayerStartPosition(x, y)
						// EntityTypeDoorは現在EntityTemplateに追加メソッドがないため
						// 将来的に追加予定
					}
				}
			}
		}
	}

	return plan, nil
}

// ValidateStringMap は文字列マップの妥当性をチェックする
func ValidateStringMap(tileMap, entityMap []string) error {
	if len(tileMap) == 0 {
		return fmt.Errorf("タイルマップが空です")
	}

	// 行の長さの一貫性をチェック
	expectedWidth := len(tileMap[0])
	for i, row := range tileMap {
		if len(row) != expectedWidth {
			return fmt.Errorf("タイルマップの行 %d の長さが不一致: 期待値 %d, 実際 %d", i, expectedWidth, len(row))
		}
	}

	// エンティティマップのサイズをチェック
	if len(entityMap) > 0 {
		if len(entityMap) != len(tileMap) {
			return fmt.Errorf("エンティティマップとタイルマップの行数が不一致: %d vs %d", len(entityMap), len(tileMap))
		}
		for i, row := range entityMap {
			if len(row) != expectedWidth {
				return fmt.Errorf("エンティティマップの行 %d の長さが不一致: 期待値 %d, 実際 %d", i, expectedWidth, len(row))
			}
		}
	}

	return nil
}
