package loader

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/ruins/assets"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/utils"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

type GameComponentList struct {
	GridElement      *gc.GridElement
	Player           *gc.Player
	Wall             *gc.Wall
	Warp             *gc.Warp
	Item             *gc.Item
	Name             *gc.Name
	Description      *gc.Description
	InBackpack       *gc.InBackpack
	Equipped         *gc.Equipped
	Consumable       *gc.Consumable
	InParty          *gc.InParty
	Member           *gc.Member
	Pools            *gc.Pools
	ProvidesHealing  *gc.ProvidesHealing
	InflictsDamage   *gc.InflictsDamage
	Attack           *gc.Attack
	Material         *gc.Material
	Recipe           *gc.Recipe
	Wearable         *gc.Wearable
	Attributes       *gc.Attributes
	EquipmentChanged *gc.EquipmentChanged
	Card             *gc.Card

	Position     *gc.Position
	SpriteRender *ec.SpriteRender
	BlockView    *gc.BlockView
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
	for iEntity, _ := range entityGameMetadata.Entities {
		gameComponentList[iEntity] = GameComponentList{}
	}

	interfaceSlice := make([]interface{}, len(gameComponentList))
	for i, v := range gameComponentList {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
