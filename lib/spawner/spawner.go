package spawner

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	w "github.com/kijimaD/ruins/lib/engine/world"
	gloader "github.com/kijimaD/ruins/lib/loader"
	"github.com/kijimaD/ruins/lib/raw"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// アイテムを生成する
func SpawnItem(world w.World, name string, spawnType raw.SpawnType) ecs.Entity {
	componentList := loader.EntityComponentList{}
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	componentList.Game = append(componentList.Game, rawMaster.GenerateItem(name, spawnType))
	entities := loader.AddEntities(world, componentList)

	return entities[len(entities)-1]
}

// パーティに追加可能なキャラを生成する
func SpawnMember(world w.World, name string, inParty bool) ecs.Entity {
	componentList := loader.EntityComponentList{}
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	componentList.Game = append(componentList.Game, rawMaster.GenerateMember(name, inParty))
	entities := loader.AddEntities(world, componentList)

	return entities[len(entities)-1]
}

// 所持素材の個数を0で初期化する
func SpawnAllMaterials(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.MaterialIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateMaterial(k, 0, raw.SpawnInBackpack))
		loader.AddEntities(world, componentList)
	}
}

// レシピ初期化
func SpawnAllRecipes(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.RecipeIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateRecipe(k))
		loader.AddEntities(world, componentList)
	}
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
			Camera:   &gc.Camera{Scale: 01, ScaleTo: 1},
		})
		loader.AddEntities(world, componentList)
	}
}

// ================
// TODO: フィールド系は、ステージ初期化でしか使わないのでloaderに移動させる

// フィールド上に表示される床を生成する
func SpawnFloor(world w.World, x gc.Row, y gc.Col) {
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
	loader.AddEntities(world, componentList)
}

// TODO: これもGridElementにする
// フィールド上に表示される壁を生成する
func SpawnFieldWall(world w.World, x int, y int) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, gloader.GameComponentList{
		Position:     &gc.Position{X: x, Y: y, Depth: gc.DepthNumTaller},
		SpriteRender: &ec.SpriteRender{SpriteSheet: &fieldSpriteSheet, SpriteNumber: 1},
		BlockView:    &gc.BlockView{},
		BlockPass:    &gc.BlockPass{},
	})
	loader.AddEntities(world, componentList)
}

// フィールド上に表示される階段を生成する
func SpawnFieldWarpNext(world w.World, x gc.Row, y gc.Col) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, gloader.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 4,
			Depth:        ec.DepthNumRug,
		},
	})
	loader.AddEntities(world, componentList)
}
