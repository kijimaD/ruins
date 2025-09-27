package mapplaner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
)

// MapPlan はマップ生成計画を表す
// タイル配置とエンティティ配置を事前に計画し、後で一括実行する
type MapPlan struct {
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
	WallSprite *int         // 壁スプライト番号（壁の場合）
}

// EntityType はエンティティの種類を表す
type EntityType int

const (
	// EntityTypeFloor は床エンティティ
	EntityTypeFloor EntityType = iota
	// EntityTypeWall は壁エンティティ
	EntityTypeWall
	// EntityTypeWarpNext は進行ワープホール
	EntityTypeWarpNext
	// EntityTypeWarpEscape は脱出ワープホール
	EntityTypeWarpEscape
	// EntityTypeProp は置物エンティティ
	EntityTypeProp
	// EntityTypeNPC はNPCエンティティ
	EntityTypeNPC
	// EntityTypeItem はアイテムエンティティ
	EntityTypeItem
	// EntityTypePlayer はプレイヤーエンティティ
	EntityTypePlayer
	// EntityTypeDoor はドアエンティティ
	EntityTypeDoor
)

// String はEntityTypeの文字列表現を返す
func (et EntityType) String() string {
	switch et {
	case EntityTypeFloor:
		return "Floor"
	case EntityTypeWall:
		return "Wall"
	case EntityTypeWarpNext:
		return "WarpNext"
	case EntityTypeWarpEscape:
		return "WarpEscape"
	case EntityTypeProp:
		return "Prop"
	case EntityTypeNPC:
		return "NPC"
	case EntityTypeItem:
		return "Item"
	case EntityTypePlayer:
		return "Player"
	case EntityTypeDoor:
		return "Door"
	default:
		return "Unknown"
	}
}

// NewMapPlan は新しいMapPlanを作成する
func NewMapPlan(width, height int) *MapPlan {
	return &MapPlan{
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
func (mp *MapPlan) AddTile(x, y int, tileType Tile) {
	mp.Tiles = append(mp.Tiles, TileSpec{
		X:        x,
		Y:        y,
		TileType: tileType,
	})
}

// AddFloor は床エンティティを計画に追加する
func (mp *MapPlan) AddFloor(x, y int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeFloor,
	})
}

// AddWall は壁エンティティを計画に追加する
func (mp *MapPlan) AddWall(x, y int, spriteNumber int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeWall,
		WallSprite: &spriteNumber,
	})
}

// AddProp は置物エンティティを計画に追加する
func (mp *MapPlan) AddProp(x, y int, propType gc.PropType) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeProp,
		PropType:   &propType,
	})
}

// AddWarpNext は進行ワープホールを計画に追加する
func (mp *MapPlan) AddWarpNext(x, y int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeWarpNext,
	})
}

// AddWarpEscape は脱出ワープホールを計画に追加する
func (mp *MapPlan) AddWarpEscape(x, y int) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeWarpEscape,
	})
}

// AddNPC はNPCエンティティを計画に追加する
func (mp *MapPlan) AddNPC(x, y int, npcType string) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeNPC,
		NPCType:    &npcType,
	})
}

// AddItem はアイテムエンティティを計画に追加する
func (mp *MapPlan) AddItem(x, y int, itemType string) {
	mp.Entities = append(mp.Entities, EntitySpec{
		X:          x,
		Y:          y,
		EntityType: EntityTypeItem,
		ItemType:   &itemType,
	})
}

// ValidatePlan は計画の妥当性をチェックする
func (mp *MapPlan) ValidatePlan() error {
	// 座標範囲チェック
	for _, tile := range mp.Tiles {
		if tile.X < 0 || tile.X >= mp.Width || tile.Y < 0 || tile.Y >= mp.Height {
			return NewValidationError("タイル座標が範囲外", tile.X, tile.Y)
		}
	}

	for _, entity := range mp.Entities {
		if entity.X < 0 || entity.X >= mp.Width || entity.Y < 0 || entity.Y >= mp.Height {
			return NewValidationError("エンティティ座標が範囲外", entity.X, entity.Y)
		}
	}

	return nil
}

// ValidationError は計画検証エラー
type ValidationError struct {
	Message string
	X, Y    int
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: (%d, %d)", e.Message, e.X, e.Y)
}

// NewValidationError は新しいValidationErrorを作成する
func NewValidationError(message string, x, y int) ValidationError {
	return ValidationError{
		Message: message,
		X:       x,
		Y:       y,
	}
}

// SetPlayerStartPosition はプレイヤーの開始位置を設定する
func (mp *MapPlan) SetPlayerStartPosition(x, y int) {
	mp.PlayerStartX = x
	mp.PlayerStartY = y
	mp.HasPlayerPos = true
}

// GetPlayerStartPosition はプレイヤーの開始位置を取得する
// プレイヤー位置が設定されていない場合はfalseを返す
func (mp *MapPlan) GetPlayerStartPosition() (int, int, bool) {
	return mp.PlayerStartX, mp.PlayerStartY, mp.HasPlayerPos
}
