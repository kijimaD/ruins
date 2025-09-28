package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// EntityPlan はマップ生成計画を表す
// タイル配置とエンティティ配置を事前に計画し、後で一括実行する
type EntityPlan struct {
	Width        int          // マップ幅
	Height       int          // マップ高さ
	Tiles        []TileSpec   // タイル配置計画
	Entities     []EntitySpec // エンティティ配置計画
	PlayerStartX int          // プレイヤー開始X座標
	PlayerStartY int          // プレイヤー開始Y座標
	HasPlayerPos bool         // プレイヤー位置が設定されているか
}

// TileSpec はタイル配置仕様
type TileSpec struct {
	X        int  // X座標
	Y        int  // Y座標
	TileType Tile // タイルタイプ
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
		Tiles:        make([]TileSpec, 0),
		Entities:     make([]EntitySpec, 0),
		PlayerStartX: 0,
		PlayerStartY: 0,
		HasPlayerPos: false,
	}
}

// AddTile はタイル配置を計画に追加する
func (mp *EntityPlan) AddTile(x, y int, tileType Tile) {
	mp.Tiles = append(mp.Tiles, TileSpec{
		X:        x,
		Y:        y,
		TileType: tileType,
	})
}

// AddFloor は床エンティティを計画に追加する
func (mp *EntityPlan) AddFloor(x, y int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeFloor,
	})
}

// AddWall は壁エンティティを計画に追加する
func (mp *EntityPlan) AddWall(x, y int, spriteNumber int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeWall,
		WallSprite: &spriteNumber,
	})
}

// AddWallWithType は壁タイプを指定して壁エンティティを計画に追加する
// スプライト番号はmapspawnerで決定される
func (mp *EntityPlan) AddWallWithType(x, y int, wallType WallType) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeWall,
		WallType:   &wallType,
	})
}

// AddProp は置物エンティティを計画に追加する
func (mp *EntityPlan) AddProp(x, y int, propType gc.PropType) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeProp,
		PropType:   &propType,
	})
}

// AddWarpNext は進行ワープホールを計画に追加する
func (mp *EntityPlan) AddWarpNext(x, y int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeWarpNext,
	})
}

// AddWarpEscape は脱出ワープホールを計画に追加する
func (mp *EntityPlan) AddWarpEscape(x, y int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeWarpEscape,
	})
}

// AddNPC はNPCエンティティを計画に追加する
func (mp *EntityPlan) AddNPC(x, y int, npcType string) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeNPC,
		NPCType:    &npcType,
	})
}

// AddItem はアイテムエンティティを計画に追加する
func (mp *EntityPlan) AddItem(x, y int, itemType string) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeItem,
		ItemType:   &itemType,
	})
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
