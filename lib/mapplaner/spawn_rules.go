// Package mapbuilder のスポーンルールシステム
// インターフェースベースの戦略パターンにより、
// 各種エンティティ（Props、NPCs、アイテムなど）の
// 柔軟なスポーン制御を提供する
package mapplaner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// SpawnRule はスポーンルールを表すインターフェース
// 各種エンティティのスポーン条件と処理を定義する
type SpawnRule interface {
	// ShouldSpawn はこのルールを適用すべきかどうかを判定する
	ShouldSpawn(builderType BuilderType, buildData *BuilderMap) bool
	// Execute はスポーン処理を実行する
	Execute(world w.World, buildData *BuilderMap) error
	// GetName はルール名を返す（ログやデバッグ用）
	GetName() string
}

// DungeonSpawnRule はダンジョン用のスポーンルール
// 探索重視の環境で、戦闘・アイテム収集に適したエンティティを配置
type DungeonSpawnRule struct {
	name        string
	builderType BuilderType // 実行時のBuilderTypeを記録
}

// NewDungeonSpawnRule はダンジョン用スポーンルールを作成する
func NewDungeonSpawnRule() SpawnRule {
	return &DungeonSpawnRule{
		name: "DungeonSpawnRule",
	}
}

// ShouldSpawn はダンジョン用スポーンルールを適用すべきかを判定する
func (r *DungeonSpawnRule) ShouldSpawn(builderType BuilderType, _ *BuilderMap) bool {
	// ダンジョン系のBuilderTypeで適用
	switch builderType.Name {
	case BuilderTypeSmallRoom.Name,
		BuilderTypeBigRoom.Name:
		// builderTypeを記録して、Execute時に使用
		r.builderType = builderType
		return true
	default:
		return false
	}
}

// Execute はダンジョン用エンティティの配置を実行する
func (r *DungeonSpawnRule) Execute(_ w.World, _ *BuilderMap) error {
	// ダンジョンではPropsを配置しない（シンプルなダンジョン）

	// TODO: NPCs、フィールドアイテム、ポータルなどの配置
	// - 敵性NPCの配置
	// - 宝箱、武器、回復アイテムの配置
	// - 隠し通路、トラップの配置

	return nil
}

// GetName はダンジョン用スポーンルール名を返す
func (r *DungeonSpawnRule) GetName() string {
	return r.name
}

// TownSpawnRule は街用のスポーンルール
// 生活重視の環境で、NPCとの交流や休息に適したエンティティを配置
type TownSpawnRule struct {
	name string
}

// NewTownSpawnRule は街用スポーンルールを作成する
func NewTownSpawnRule() SpawnRule {
	return &TownSpawnRule{
		name: "TownSpawnRule",
	}
}

// ShouldSpawn は街用スポーンルールを適用すべきかを判定する
func (r *TownSpawnRule) ShouldSpawn(builderType BuilderType, _ *BuilderMap) bool {
	// 街系のBuilderTypeで適用
	switch builderType.Name {
	case BuilderTypeTown.Name:
		return true
	default:
		return false
	}
}

// Execute は街用エンティティの配置を実行する
func (r *TownSpawnRule) Execute(world w.World, buildData *BuilderMap) error {
	// Props（家具）の配置
	if err := r.spawnTownProps(world, buildData); err != nil {
		return err
	}

	// TODO: 街用エンティティの配置
	// - 友好NPCs（商人、住民など）
	// - ショップアイテム
	// - 休息施設（宿屋など）

	return nil
}

// GetName は街用スポーンルール名を返す
func (r *TownSpawnRule) GetName() string {
	return r.name
}

// FixedPropPlacement は固定配置する家具の情報
type FixedPropPlacement struct {
	PropType gc.PropType
	X        gc.Tile
	Y        gc.Tile
}

// spawnTownProps は市街地用の家具を配置する
func (r *TownSpawnRule) spawnTownProps(world w.World, buildData *BuilderMap) error {
	// 各建物に固定配置する家具を定義
	fixedPlacements := r.getFixedPropPlacements(buildData)

	// 固定位置に家具を配置
	for _, placement := range fixedPlacements {
		if err := r.placePropAtPosition(world, buildData, placement); err != nil {
			// エラーをログに記録するが、処理は継続
			continue
		}
	}

	return nil
}

// getFixedPropPlacements は各建物の固定家具配置を返す
func (r *TownSpawnRule) getFixedPropPlacements(buildData *BuilderMap) []FixedPropPlacement {
	centerX := int(buildData.Level.TileWidth) / 2
	centerY := int(buildData.Level.TileHeight) / 2

	placements := []FixedPropPlacement{}

	// === 北の文教区域 ===

	// 図書館（書籍と閲覧席を配置）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeBookshelf, gc.Tile(centerX - 8), gc.Tile(centerY - 19)}, // 北壁沿い
		{gc.PropTypeBookshelf, gc.Tile(centerX - 6), gc.Tile(centerY - 19)}, // 北壁沿い
		{gc.PropTypeTable, gc.Tile(centerX - 7), gc.Tile(centerY - 17)},     // 閲覧机
		{gc.PropTypeChair, gc.Tile(centerX - 7), gc.Tile(centerY - 16)},     // 閲覧用椅子
		{gc.PropTypeTable, gc.Tile(centerX - 4), gc.Tile(centerY - 15)},     // 学習机
	}...)

	// 学校（教育設備を配置）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeBookshelf, gc.Tile(centerX + 6), gc.Tile(centerY - 19)},  // 北壁
		{gc.PropTypeBookshelf, gc.Tile(centerX + 8), gc.Tile(centerY - 19)},  // 北壁
		{gc.PropTypeBookshelf, gc.Tile(centerX + 10), gc.Tile(centerY - 19)}, // 北壁
		{gc.PropTypeBookshelf, gc.Tile(centerX + 5), gc.Tile(centerY - 16)},  // 西壁
		{gc.PropTypeBookshelf, gc.Tile(centerX + 11), gc.Tile(centerY - 16)}, // 東壁
		{gc.PropTypeTable, gc.Tile(centerX + 8), gc.Tile(centerY - 15)},      // 教卓
		{gc.PropTypeChair, gc.Tile(centerX + 8), gc.Tile(centerY - 14)},      // 教師用椅子
	}...)

	// === 東の居住区域 ===

	// 住民の家1（居住用家具）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeBed, gc.Tile(centerX + 13), gc.Tile(centerY - 7)},   // 寝室
		{gc.PropTypeTable, gc.Tile(centerX + 15), gc.Tile(centerY - 5)}, // 食事台
		{gc.PropTypeChair, gc.Tile(centerX + 15), gc.Tile(centerY - 4)}, // 食事用椅子
		{gc.PropTypeChair, gc.Tile(centerX + 16), gc.Tile(centerY - 5)}, // 食事用椅子
	}...)

	// 住民の家2（居住用家具）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeBed, gc.Tile(centerX + 14), gc.Tile(centerY + 2)},   // 寝室
		{gc.PropTypeTable, gc.Tile(centerX + 16), gc.Tile(centerY + 4)}, // 食事台
		{gc.PropTypeChair, gc.Tile(centerX + 16), gc.Tile(centerY + 5)}, // 食事用椅子
	}...)

	// === 中央街区・南の公共区域 ===

	// 公民館（集会用の座席配置）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeChair, gc.Tile(centerX - 6), gc.Tile(centerY + 12)}, // 集会用座席
		{gc.PropTypeChair, gc.Tile(centerX - 4), gc.Tile(centerY + 12)}, // 集会用座席
		{gc.PropTypeChair, gc.Tile(centerX + 4), gc.Tile(centerY + 12)}, // 集会用座席
		{gc.PropTypeChair, gc.Tile(centerX + 6), gc.Tile(centerY + 12)}, // 集会用座席
	}...)

	// 事務所（管理業務用設備）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeBed, gc.Tile(centerX + 12), gc.Tile(centerY + 13)},       // 休憩用ベッド
		{gc.PropTypeTable, gc.Tile(centerX + 14), gc.Tile(centerY + 15)},     // 事務机
		{gc.PropTypeChair, gc.Tile(centerX + 14), gc.Tile(centerY + 16)},     // 事務用椅子
		{gc.PropTypeBookshelf, gc.Tile(centerX + 18), gc.Tile(centerY + 14)}, // 書類棚
	}...)

	// === 西の商業区域 ===

	// 商店（販売設備）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeTable, gc.Tile(centerX - 18), gc.Tile(centerY - 4)}, // 販売台
		{gc.PropTypeChair, gc.Tile(centerX - 17), gc.Tile(centerY - 4)}, // 客用椅子
		{gc.PropTypeChair, gc.Tile(centerX - 19), gc.Tile(centerY - 4)}, // 客用椅子
		{gc.PropTypeTable, gc.Tile(centerX - 15), gc.Tile(centerY - 1)}, // 商品台
		{gc.PropTypeTable, gc.Tile(centerX - 13), gc.Tile(centerY + 1)}, // 商品台
	}...)

	// 倉庫（保管設備）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeCrate, gc.Tile(centerX - 16), gc.Tile(centerY + 8)},   // 木箱
		{gc.PropTypeCrate, gc.Tile(centerX - 14), gc.Tile(centerY + 8)},   // 木箱
		{gc.PropTypeBarrel, gc.Tile(centerX - 12), gc.Tile(centerY + 10)}, // 樽
		{gc.PropTypeBarrel, gc.Tile(centerX - 10), gc.Tile(centerY + 12)}, // 樽
		{gc.PropTypeTable, gc.Tile(centerX - 15), gc.Tile(centerY + 13)},  // 作業台
	}...)

	// === 郊外区域 ===

	// 小さな住宅（シンプルな生活空間）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeBed, gc.Tile(centerX - 14), gc.Tile(centerY - 13)},       // 寝床
		{gc.PropTypeTable, gc.Tile(centerX - 12), gc.Tile(centerY - 11)},     // 食事机
		{gc.PropTypeChair, gc.Tile(centerX - 12), gc.Tile(centerY - 10)},     // 椅子
		{gc.PropTypeBookshelf, gc.Tile(centerX - 10), gc.Tile(centerY - 12)}, // 本棚
	}...)

	// 公園（休憩設備）
	placements = append(placements, []FixedPropPlacement{
		{gc.PropTypeChair, gc.Tile(centerX + 14), gc.Tile(centerY + 12)}, // 休憩用座席
		{gc.PropTypeChair, gc.Tile(centerX + 16), gc.Tile(centerY + 14)}, // 休憩用座席
	}...)

	return placements
}

// placePropAtPosition は指定位置に家具を配置する
func (r *TownSpawnRule) placePropAtPosition(world w.World, buildData *BuilderMap, placement FixedPropPlacement) error {
	// 配置位置が有効かチェック
	if !buildData.IsSpawnableTile(world, placement.X, placement.Y) {
		return fmt.Errorf("位置 (%d,%d) は配置不可能", placement.X, placement.Y)
	}

	// 家具を配置
	return worldhelper.PlacePropAt(world, placement.PropType, placement.X, placement.Y)
}

// SpawnRuleEngine はスポーンルールエンジン
// 複数のSpawnRuleを管理し、適切なルールを選択・実行する
type SpawnRuleEngine struct {
	rules []SpawnRule
}

// NewSpawnRuleEngine は新しいスポーンルールエンジンを作成する
func NewSpawnRuleEngine() *SpawnRuleEngine {
	engine := &SpawnRuleEngine{
		rules: []SpawnRule{},
	}

	// デフォルトルールを登録
	engine.AddRule(NewDungeonSpawnRule())
	engine.AddRule(NewTownSpawnRule())

	return engine
}

// AddRule はスポーンルールを追加する
func (e *SpawnRuleEngine) AddRule(rule SpawnRule) {
	e.rules = append(e.rules, rule)
}

// ExecuteRules は適用可能なルールを実行する
func (e *SpawnRuleEngine) ExecuteRules(builderType BuilderType, world w.World, buildData *BuilderMap) error {
	for _, rule := range e.rules {
		if rule.ShouldSpawn(builderType, buildData) {
			if err := rule.Execute(world, buildData); err != nil {
				return err
			}
		}
	}
	return nil
}
