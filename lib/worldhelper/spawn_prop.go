package worldhelper

import (
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/raw"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// SpawnProp は置物を生成する統一関数
func SpawnProp(world w.World, propName string, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	// RawMasterから置物の設定を生成
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	entitySpec, err := rawMaster.NewPropSpec(propName)
	if err != nil {
		return ecs.Entity(0), err
	}

	// 床を下敷きとして配置（既に床がある場合は配置しない）
	if !hasFloorAt(world, x, y) {
		_, _ = SpawnFloor(world, x, y, "field", "floor")
	}

	// 位置情報を設定
	entitySpec.GridElement = &gc.GridElement{X: x, Y: y}

	// エンティティを生成
	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, entitySpec)
	entities := entities.AddEntities(world, componentList)
	return entities[len(entities)-1], nil
}

// hasFloorAt は指定位置に床が存在するかチェックする
func hasFloorAt(world w.World, x, y gc.Tile) bool {
	floorExists := false

	world.Manager.Join(
		world.Components.GridElement,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)

		if gridElement.X == x && gridElement.Y == y && spriteRender.Depth == gc.DepthNumFloor {
			floorExists = true
			return
		}
	}))

	return floorExists
}
