// Package mapplanner の文字列ベースマップビルダー
// 文字列で直接タイルとエンティティ配置を指定できるシステム
package mapplanner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
)

// StringMapPlanner は文字列ベースでマップを生成するプランナー
type StringMapPlanner struct {
	TileMap       []string                // タイル配置を表す文字列のスライス
	EntityMap     []string                // エンティティ配置を表す文字列のスライス
	TileMapping   map[rune]string         // 文字からタイル名へのマッピング
	EntityMapping map[rune]EntityTemplate // 文字からEntityTemplateへのマッピング
}

// EntityTemplate はエンティティ生成計画
type EntityTemplate struct {
	EntityType EntityType
	Data       interface{} // PropType等の追加データ
}

// PlanInitial は文字列定義からマップ構造を初期化する
func (b StringMapPlanner) PlanInitial(planData *MetaPlan) error {
	if len(b.TileMap) == 0 {
		return nil
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
	existingWidth := int(planData.Level.TileWidth)
	existingHeight := int(planData.Level.TileHeight)

	// 文字列サイズと指定サイズが異なる場合は指定サイズを使用
	if existingWidth > 0 && existingHeight > 0 {
		width = existingWidth
		height = existingHeight
	} else {
		// BuildDataのサイズを文字列サイズで調整
		planData.Level.TileWidth = gc.Tile(width)
		planData.Level.TileHeight = gc.Tile(height)
	}

	planData.Tiles = make([]raw.TileRaw, width*height)

	// タイルマップを解析してタイルを配置
	for y, row := range b.TileMap {
		for x, char := range row {
			if x >= width || y >= height {
				continue
			}

			tileIdx := planData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			if tileName, exists := b.TileMapping[char]; exists {
				planData.Tiles[tileIdx] = planData.GenerateTile(tileName)
			} else {
				return fmt.Errorf("未知のタイル文字 '%c' が位置 (%d, %d) で見つかりました", char, x, y)
			}
		}
	}

	// エンティティマップが指定されている場合はエンティティを配置
	if len(b.EntityMap) > 0 {
		if err := b.parseEntities(planData, width, height); err != nil {
			return err
		}
	}

	return nil
}

// parseEntities はエンティティマップを解析して計画に追加する
func (b StringMapPlanner) parseEntities(plannerMap *MetaPlan, width, height int) error {
	// エンティティマップからワープポータルを解析してMetaPlanに追加
	for y := 0; y < len(b.EntityMap) && y < height; y++ {
		for x := 0; x < len(b.EntityMap[y]) && x < width; x++ {
			char := rune(b.EntityMap[y][x])

			// 空文字やドット（空エンティティ）はスキップ
			if char == '.' {
				continue
			}

			// ワープポータル文字を処理（床タイル + MetaPlanへのワープポータル情報追加）
			switch char {
			case 'w':
				// 進行ワープポータル位置に床タイルを配置
				tileIdx := plannerMap.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				plannerMap.Tiles[tileIdx] = plannerMap.GenerateTile("Floor")
				// MetaPlanにワープポータル情報を追加
				plannerMap.WarpPortals = append(plannerMap.WarpPortals, WarpPortal{
					X:    x,
					Y:    y,
					Type: WarpPortalNext,
				})
			case 'e':
				// 帰還ワープポータル位置に床タイルを配置
				tileIdx := plannerMap.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				plannerMap.Tiles[tileIdx] = plannerMap.GenerateTile("Floor")
				// MetaPlanにワープポータル情報を追加
				plannerMap.WarpPortals = append(plannerMap.WarpPortals, WarpPortal{
					X:    x,
					Y:    y,
					Type: WarpPortalEscape,
				})
			default:
				// その他のエンティティ文字をチェック
				if _, exists := b.EntityMapping[char]; !exists {
					return fmt.Errorf("未知のエンティティ文字 '%c' が位置 (%d, %d) で見つかりました", char, x, y)
				}
			}
		}
	}

	return nil
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
func getDefaultTileMapping() map[rune]string {
	return map[rune]string{
		'#': "Wall",  // 壁
		'f': "Floor", // 床（建物内・住宅）
		'r': "Floor", // 道路 (TODO: スプライト未対応)
		'.': "Floor", // 空地・庭
	}
}

// getDefaultEntityMapping はデフォルトのエンティティ文字マッピングを返す
func getDefaultEntityMapping() map[rune]EntityTemplate {
	return map[rune]EntityTemplate{
		'@': {EntityType: EntityTypePlayer},
		'&': {EntityType: EntityTypeNPC},
		'C': {EntityType: EntityTypeProp, Data: gc.PropTypeChair},
		'T': {EntityType: EntityTypeProp, Data: gc.PropTypeTable},
		'S': {EntityType: EntityTypeProp, Data: gc.PropTypeBookshelf},
		'M': {EntityType: EntityTypeProp, Data: gc.PropTypeBarrel},
		'R': {EntityType: EntityTypeProp, Data: gc.PropTypeCrate},
	}
}

// SetTileMapping はカスタムタイルマッピングを設定する
func (b *StringMapPlanner) SetTileMapping(mapping map[rune]string) *StringMapPlanner {
	b.TileMapping = mapping
	return b
}

// SetEntityMapping はカスタムエンティティマッピングを設定する
func (b *StringMapPlanner) SetEntityMapping(mapping map[rune]EntityTemplate) *StringMapPlanner {
	b.EntityMapping = mapping
	return b
}

// AddTileMapping は追加のタイルマッピングを設定する
func (b *StringMapPlanner) AddTileMapping(char rune, tileName string) *StringMapPlanner {
	if b.TileMapping == nil {
		b.TileMapping = getDefaultTileMapping()
	}
	b.TileMapping[char] = tileName
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

// BuildEntityPlanFromStrings は文字列マップからEntityPlanを生成する
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

			var tileName string
			if mappedName, exists := tileMapping[char]; exists {
				tileName = mappedName
			} else {
				return nil, fmt.Errorf("未知のタイル文字 '%c' が位置 (%d, %d) で見つかりました", char, x, y)
			}

			// タイル名に応じてエンティティを追加
			if tileName == "Floor" {
				plan.AddFloor(x, y)
			} else {
				plan.AddWall(x, y, 0) // スプライト番号は0をデフォルト
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

				// ワープポータルを処理
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
				} else {
					return nil, fmt.Errorf("未知のエンティティ文字 '%c' が位置 (%d, %d) で見つかりました", char, x, y)
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
