package spawner

import (
	"github.com/kijimaD/ruins/lib/engine/loader"
	w "github.com/kijimaD/ruins/lib/engine/world"
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
