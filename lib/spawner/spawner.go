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
