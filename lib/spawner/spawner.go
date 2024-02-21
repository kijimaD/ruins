package spawner

import (
	"github.com/kijimaD/ruins/lib/engine/loader"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raw"
)

func SpawnItem(world w.World, name string) {
	componentList := loader.EntityComponentList{}
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	componentList.Game = append(componentList.Game, rawMaster.GenerateItem(name))
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
	loader.AddEntities(world, componentList)
}

func SpawnMember(world w.World, name string, inParty bool) {
	componentList := loader.EntityComponentList{}
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	componentList.Game = append(componentList.Game, rawMaster.GenerateMember(name, inParty))
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
	loader.AddEntities(world, componentList)
}

func SpawnAllMaterials(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.MaterialIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateMaterial(k))
		componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
		loader.AddEntities(world, componentList)
	}
}

func SpawnAllRecipes(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.RecipeIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateRecipe(k))
		componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
		loader.AddEntities(world, componentList)
	}
}
