package worldhelper

import (
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/raw"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// SpawnProp は置物を生成する
func SpawnProp(world w.World, propName string, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	// RawMasterから置物の設定を生成
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	entitySpec, err := rawMaster.NewPropSpec(propName)
	if err != nil {
		return ecs.Entity(0), err
	}

	// 床を下敷きとして配置（既に床がある場合は配置しない）
	if !hasFloorAt(world, x, y) {
		_, _ = SpawnTile(world, "Floor", x, y, nil)
	}

	// 位置情報を設定
	entitySpec.GridElement = &gc.GridElement{X: x, Y: y}

	// エンティティを生成
	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, entitySpec)
	entities, err := entities.AddEntities(world, componentList)
	if err != nil {
		return ecs.Entity(0), err
	}
	return entities[len(entities)-1], nil
}

// SpawnDoor はドアを生成する
func SpawnDoor(world w.World, x gc.Tile, y gc.Tile, orientation gc.DoorOrientation) (ecs.Entity, error) {
	// 床を下敷きとして配置（既に床がある場合は配置しない）
	if !hasFloorAt(world, x, y) {
		_, err := SpawnTile(world, "Floor", x, y, nil)
		if err != nil {
			return ecs.Entity(0), err
		}
	}

	// スプライトキーを決定（閉じたドア）
	var spriteKey string
	if orientation == gc.DoorOrientationHorizontal {
		spriteKey = "door_horizontal_closed"
	} else {
		spriteKey = "door_vertical_closed"
	}

	// EntitySpecを構築
	entitySpec := gc.EntitySpec{
		Name:        &gc.Name{Name: "ドア"},
		Description: &gc.Description{Description: "開閉できるドア"},
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheetName: "field",
			SpriteKey:       spriteKey,
			Depth:           gc.DepthNumTaller,
		},
		BlockPass: &gc.BlockPass{}, // 閉じているので通行不可
		BlockView: &gc.BlockView{}, // 閉じているので視線を遮る
		Door: &gc.Door{
			IsOpen:      false,
			Orientation: orientation,
		},
	}

	// エンティティを生成
	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, entitySpec)
	ents, err := entities.AddEntities(world, componentList)
	if err != nil {
		return ecs.Entity(0), err
	}
	return ents[len(ents)-1], nil
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
