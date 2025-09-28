package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// EntityPlan はマップ生成計画を表す
// タイル配置とエンティティ配置を事前に計画し、後で一括実行する
type EntityPlan struct {
	Width        int          // マップ幅
	Height       int          // マップ高さ
	Entities     []EntitySpec // エンティティ配置計画。計画は最終的にエンティティ列として表現される
	PlayerStartX int          // プレイヤー開始X座標
	PlayerStartY int          // プレイヤー開始Y座標
	HasPlayerPos bool         // プレイヤー位置が設定されているか
}

// EntitySpec はエンティティ配置仕様
type EntitySpec struct {
	X          int          // X座標
	Y          int          // Y座標
	EntityType EntityType   // エンティティタイプ
	PropType   *gc.PropType // 置物タイプ（置物の場合）
	NPCType    *string      // NPCタイプ（NPCの場合）
	ItemType   *string      // アイテムタイプ（アイテムの場合）
	WallSprite *int         // 壁スプライト番号（壁の場合、直接指定）
	WallType   *WallType    // 壁タイプ
}

// EntityType はエンティティの種類を表す
type EntityType string

const (
	// EntityTypeFloor は床エンティティ
	EntityTypeFloor EntityType = "Floor"
	// EntityTypeWall は壁エンティティ
	EntityTypeWall EntityType = "Wall"
	// EntityTypeWarpNext は進行ワープホール
	EntityTypeWarpNext EntityType = "WarpNext"
	// EntityTypeWarpEscape は脱出ワープホール
	EntityTypeWarpEscape EntityType = "WarpEscape"
	// EntityTypeProp は置物エンティティ
	EntityTypeProp EntityType = "Prop"
	// EntityTypeNPC はNPCエンティティ
	EntityTypeNPC EntityType = "NPC"
	// EntityTypeItem はアイテムエンティティ
	EntityTypeItem EntityType = "Item"
	// EntityTypePlayer はプレイヤーエンティティ
	EntityTypePlayer EntityType = "Player"
	// EntityTypeDoor はドアエンティティ
	EntityTypeDoor EntityType = "Door"
)

// NewEntityPlan は新しいEntityPlanを作成する
func NewEntityPlan(width, height int) *EntityPlan {
	return &EntityPlan{
		Width:        width,
		Height:       height,
		Entities:     make([]EntitySpec, 0),
		PlayerStartX: 0,
		PlayerStartY: 0,
		HasPlayerPos: false,
	}
}

// SetPlayerStartPosition はプレイヤーの開始位置を設定する
func (mp *EntityPlan) SetPlayerStartPosition(x, y int) {
	mp.PlayerStartX = x
	mp.PlayerStartY = y
	mp.HasPlayerPos = true
}

// GetPlayerStartPosition はプレイヤーの開始位置を取得する
// プレイヤー位置が設定されていない場合はfalseを返す
func (mp *EntityPlan) GetPlayerStartPosition() (int, int, bool) {
	return mp.PlayerStartX, mp.PlayerStartY, mp.HasPlayerPos
}
