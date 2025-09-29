package worldhelper

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/props"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// SpawnProp は置物を生成する統一関数
func SpawnProp(world w.World, propType gc.PropType, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	// PropManagerを取得（後でResourcesに追加する必要がある）
	propManager := props.NewPropManager()

	config, exists := propManager.GetConfig(propType)
	if !exists {
		return ecs.Entity(0), fmt.Errorf("未定義の置物タイプ: %s", propType)
	}

	// 床を下敷きとして配置（既に床がある場合は配置しない）
	if !hasFloorAt(world, x, y) {
		_, _ = SpawnFloor(world, x, y, "field", "floor")
	}

	// 置物エンティティを構築
	componentList := entities.ComponentList[gc.EntitySpec]{}
	entitySpec := gc.EntitySpec{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheetName: "field",
			SpriteKey:       config.SpriteKey,
			Depth:           gc.DepthNumRug, // アイテムと同じ深度
		},
		PropType: &propType,
		Name: &gc.Name{
			Name: config.Name,
		},
		Description: &gc.Description{
			Description: config.Description,
		},
	}

	// 設定に応じてコンポーネントを追加
	if config.BlocksMovement {
		entitySpec.BlockPass = &gc.BlockPass{}
	}
	if config.BlocksVisibility {
		entitySpec.BlockView = &gc.BlockView{}
	}

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
