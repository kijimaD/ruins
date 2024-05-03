package spawner

import (
	gc "github.com/kijimaD/ruins/lib/components"
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
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
	entities := loader.AddEntities(world, componentList)

	return entities[len(entities)-1]
}

// プレイアブルキャラを生成する
func SpawnMember(world w.World, name string, inParty bool) ecs.Entity {
	componentList := loader.EntityComponentList{}
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	componentList.Game = append(componentList.Game, rawMaster.GenerateMember(name, inParty))
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
	entities := loader.AddEntities(world, componentList)

	return entities[len(entities)-1]
}

// 所持素材の個数を0で初期化する
func SpawnAllMaterials(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.MaterialIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateMaterial(k, 0, raw.SpawnInBackpack))
		componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
		loader.AddEntities(world, componentList)
	}
}

// レシピ初期化
func SpawnAllRecipes(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.RecipeIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateRecipe(k))
		componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
		loader.AddEntities(world, componentList)
	}
}

// フィールド上に表示されるプレイヤーを生成する
func SpawnPlayer(world w.World, x int, y int) {
	gcl := gloader.GameComponentList{}
	gcl.Position = &gc.Position{x, y}
	gcl.Player = &gc.Player{}

	// fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, gloader.GameComponentList{
		Position: &gc.Position{x, y},
		Player:   &gc.Player{},
		// Render:   &gc.Render{SpriteSheet: &fieldSpriteSheet, SpriteNumber: 3},
	})
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
	loader.AddEntities(world, componentList)
}
