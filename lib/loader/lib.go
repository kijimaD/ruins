package loader

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/sokotwo/assets"
	gc "github.com/kijimaD/sokotwo/lib/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
)

type GameComponentList struct {
	GridElement *gc.GridElement
	Player      *gc.Player
	Wall        *gc.Wall
	Warp        *gc.Warp
	Item        *gc.Item
	Name        *gc.Name
	Description *gc.Description
	InBackpack  *gc.InBackpack
	Consumable  *gc.Consumable
	InParty     *gc.InParty
	Member      *gc.Member
	Pools       *gc.Pools
}

type Entity struct {
	Components GameComponentList
}

func PreloadEntities(entityMetadataPath string, world w.World) loader.EntityComponentList {
	b, err := assets.FS.ReadFile(entityMetadataPath)
	if err != nil {
		log.Fatal(b)
	}
	return loader.EntityComponentList{
		Engine: loader.LoadEngineComponents(b, world),
	}
}

func PreloadGameEntities(entityMetadataPath string, world w.World) loader.EntityComponentList {
	b, err := assets.FS.ReadFile(entityMetadataPath)
	if err != nil {
		log.Fatal(b)
	}
	return loader.EntityComponentList{
		Game: LoadGameComponent(b, world),
	}
}

type entityGameMetadata struct {
	Entities []Entity `toml:"entity"`
}

func LoadGameComponent(entityMetadataContent []byte, world w.World) []interface{} {
	var entityGameMetadata entityGameMetadata
	utils.Try(toml.Decode(string(entityMetadataContent), &entityGameMetadata))

	gameComponentList := make([]GameComponentList, len(entityGameMetadata.Entities))
	for iEntity, entity := range entityGameMetadata.Entities {
		gameComponentList[iEntity] = processComponentsListData(world, entity.Components)
	}

	interfaceSlice := make([]interface{}, len(gameComponentList))
	for i, v := range gameComponentList {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

func processComponentsListData(world w.World, data GameComponentList) GameComponentList {
	return GameComponentList{
		Item:        data.Item,
		Name:        data.Name,
		Description: data.Description,
		InBackpack:  data.InBackpack,
		Consumable:  data.Consumable,
		InParty:     data.InParty,
		Member:      data.Member,
		Pools:       data.Pools,
	}
}
