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
		_, _ = SpawnFloor(world, x, y)
	}

	// 置物エンティティを構築
	componentList := entities.ComponentList[gc.EntitySpec]{}
	entitySpec := gc.EntitySpec{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			Name:         "field",
			SpriteNumber: config.SpriteNumber,
			Depth:        gc.DepthNumRug, // アイテムと同じ深度
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

// PlacePropAt は指定位置に置物を配置する
// 位置の安全性をチェックしてから配置を行う
func PlacePropAt(world w.World, propType gc.PropType, x, y gc.Tile) error {
	// 位置が安全かチェック
	if !isPositionSafeForProp(world, x, y) {
		return fmt.Errorf("位置 (%d,%d) は配置不可能: 他のエンティティと重複", x, y)
	}

	// 置物を配置
	_, err := SpawnProp(world, propType, x, y)
	return err
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

// isPositionSafeForProp は指定位置が家具配置に安全かチェックする
// プレイヤーや他のエンティティと重複しない位置かを確認（床は除外）
func isPositionSafeForProp(world w.World, x, y gc.Tile) bool {
	// 指定位置に床以外のエンティティが存在するかチェック
	entityExists := false

	world.Manager.Join(
		world.Components.GridElement,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)

		if gridElement.X == x && gridElement.Y == y {
			// 床は除外（床の上には家具を配置可能）
			if spriteRender.Depth == gc.DepthNumFloor {
				return
			}
			// 床以外のエンティティがある場合は配置不可
			entityExists = true
			return
		}
	}))

	// 床以外のエンティティが存在する場合は安全ではない
	return !entityExists
}
