package resources

import (
	"github.com/kijimaD/ruins/lib/engine/loader"
	w "github.com/kijimaD/ruins/lib/engine/world"
	gloader "github.com/kijimaD/ruins/lib/loader"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
)

// フィールド上に表示される床を生成する
func SpawnFloor(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, gloader.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 2,
			Depth:        ec.DepthNumFloor,
		},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される壁を生成する
func SpawnFieldWall(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, gloader.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 1,
			Depth:        ec.DepthNumTaller,
		},
		BlockView: &gc.BlockView{},
		BlockPass: &gc.BlockPass{},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される進行ワープホールを生成する
func SpawnFieldWarpNext(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	SpawnFloor(world, x, y) // 下敷き描画

	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, gloader.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 4,
			Depth:        ec.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeNext},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される脱出ワープホールを生成する
func SpawnFieldWarpEscape(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	SpawnFloor(world, x, y) // 下敷き描画

	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, gloader.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 5,
			Depth:        ec.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeEscape},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示されるプレイヤーを生成する
func SpawnPlayer(world w.World, x int, y int) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	{
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, gloader.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Player:   &gc.Player{},
			SpriteRender: &ec.SpriteRender{
				SpriteSheet:  &fieldSpriteSheet,
				SpriteNumber: 3,
				Depth:        ec.DepthNumTaller,
			},
		})
		loader.AddEntities(world, componentList)
	}
	// カメラ
	{
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, gloader.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Camera:   &gc.Camera{Scale: 0.1, ScaleTo: 1},
		})
		loader.AddEntities(world, componentList)
	}
}
