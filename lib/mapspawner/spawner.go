package mapspawner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Spawn はEntityPlanに基づいてレベルを生成する
func Spawn(world w.World, plan *mapplanner.EntityPlan) (resources.Level, error) {
	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// Levelオブジェクトを作成
	totalTiles := plan.Width * plan.Height
	level := resources.Level{
		TileWidth:  gc.Tile(plan.Width),
		TileHeight: gc.Tile(plan.Height),
		Entities:   make([]ecs.Entity, totalTiles),
	}

	// エンティティを一括生成
	for _, entitySpec := range plan.Entities {
		entity, err := spawnEntityFromSpec(world, entitySpec)
		if err != nil {
			return resources.Level{}, fmt.Errorf("エンティティ生成エラー (%d, %d): %w", entitySpec.X, entitySpec.Y, err)
		}

		// エンティティをLevelスライスに登録
		if entity != ecs.Entity(0) {
			tileIdx := level.XYTileIndex(gc.Tile(entitySpec.X), gc.Tile(entitySpec.Y))
			level.Entities[tileIdx] = entity
		}
	}

	return level, nil
}

// spawnEntityFromSpec は EntitySpec に基づいてエンティティを生成する
func spawnEntityFromSpec(world w.World, spec mapplanner.EntitySpec) (ecs.Entity, error) {
	x := gc.Tile(spec.X)
	y := gc.Tile(spec.Y)

	switch spec.EntityType {
	case mapplanner.EntityTypeFloor:
		return worldhelper.SpawnFloor(world, x, y)

	case mapplanner.EntityTypeWall:
		if spec.WallSprite == nil {
			return ecs.Entity(0), fmt.Errorf("壁エンティティにスプライト番号が指定されていません")
		}
		return worldhelper.SpawnWall(world, x, y, *spec.WallSprite)

	case mapplanner.EntityTypeWarpNext:
		return worldhelper.SpawnFieldWarpNext(world, x, y)

	case mapplanner.EntityTypeWarpEscape:
		return worldhelper.SpawnFieldWarpEscape(world, x, y)

	case mapplanner.EntityTypeProp:
		if spec.PropType == nil {
			return ecs.Entity(0), fmt.Errorf("置物エンティティにPropTypeが指定されていません")
		}
		return worldhelper.SpawnProp(world, *spec.PropType, x, y)

	case mapplanner.EntityTypeNPC:
		if spec.NPCType == nil {
			return ecs.Entity(0), fmt.Errorf("NPCエンティティにNPCTypeが指定されていません")
		}
		return worldhelper.SpawnEnemy(world, spec.X, spec.Y, *spec.NPCType)

	case mapplanner.EntityTypeItem:
		if spec.ItemType == nil {
			return ecs.Entity(0), fmt.Errorf("アイテムエンティティにItemTypeが指定されていません")
		}
		return worldhelper.SpawnFieldItem(world, *spec.ItemType, x, y)

	default:
		return ecs.Entity(0), fmt.Errorf("未知のエンティティタイプ: %s", spec.EntityType)
	}
}
