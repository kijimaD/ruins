package loader

import (
	"image/color"
	"reflect"

	c "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/utils"
	w "github.com/kijimaD/ruins/lib/engine/world"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// EngineComponentList is the list of engine components
type EngineComponentList struct {
	SpriteRender *c.SpriteRender
	Transform    *c.Transform
}

// EntityComponentList is a list of preloaded entities with components
type EntityComponentList struct {
	Engine []EngineComponentList
	Game   []interface{}
}

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataContent []byte, world w.World, gameComponentList []interface{}) []ecs.Entity {
	entityComponentList := EntityComponentList{
		Engine: LoadEngineComponents(entityMetadataContent, world),
		Game:   gameComponentList,
	}
	return AddEntities(world, entityComponentList)
}

// AddEntities adds entities with engine and game components
func AddEntities(world w.World, entityComponentList EntityComponentList) []ecs.Entity {
	// Create new entities and add engine components
	entities := make([]ecs.Entity, len(entityComponentList.Engine))
	for iEntity := range entityComponentList.Engine {
		entities[iEntity] = world.Manager.NewEntity()
		AddEntityComponents(entities[iEntity], world.Components.Engine, entityComponentList.Engine[iEntity])
	}

	// Add game components
	if entityComponentList.Game != nil {
		if len(entityComponentList.Game) != len(entityComponentList.Engine) {
			utils.LogFatalf("incorrect size for game component list")
		}
		for iEntity := range entities {
			AddEntityComponents(entities[iEntity], world.Components.Game, entityComponentList.Game[iEntity])
		}
	}
	return entities
}

// AddEntityComponents adds loaded components to an entity
func AddEntityComponents(entity ecs.Entity, ecsComponentList interface{}, components interface{}) ecs.Entity {
	ecv := reflect.ValueOf(ecsComponentList).Elem()
	cv := reflect.ValueOf(components)
	for iField := 0; iField < cv.NumField(); iField++ {
		if !cv.Field(iField).IsNil() {
			component := cv.Field(iField).Elem()
			value := reflect.New(reflect.TypeOf(component.Interface()))
			value.Elem().Set(component)

			ecsComponent := ecv.FieldByName(component.Type().Name()).Interface().(ecs.DataComponent)
			entity.AddComponent(ecsComponent, value.Interface())
		}
	}
	return entity
}

type engineComponentListData struct {
	SpriteRender *spriteRenderData
}

type entity struct {
	Components engineComponentListData
}

type entityEngineMetadata struct {
	Entities []entity `toml:"entity"`
}

// LoadEngineComponents loads engine components from a TOML byte slice
func LoadEngineComponents(entityMetadataContent []byte, world w.World) []EngineComponentList {
	var entityEngineMetadata entityEngineMetadata
	utils.Try(toml.Decode(string(entityMetadataContent), &entityEngineMetadata))

	engineComponentList := make([]EngineComponentList, len(entityEngineMetadata.Entities))
	for iEntity, entity := range entityEngineMetadata.Entities {
		engineComponentList[iEntity] = processComponentsListData(world, entity.Components)
	}
	return engineComponentList
}

func processComponentsListData(world w.World, data engineComponentListData) EngineComponentList {
	return EngineComponentList{
		SpriteRender: processSpriteRenderData(world, data.SpriteRender),
	}
}

//
// SpriteRender
//

type fillData struct {
	Width  int
	Height int
	Color  [4]uint8
}

type spriteRenderData struct {
	Fill            *fillData
	SpriteSheetName string `toml:"sprite_sheet_name"`
	SpriteNumber    int    `toml:"sprite_number"`
}

func processSpriteRenderData(world w.World, spriteRenderData *spriteRenderData) *c.SpriteRender {
	if spriteRenderData == nil {
		return nil
	}
	if spriteRenderData.Fill != nil && spriteRenderData.SpriteSheetName != "" {
		utils.LogFatalf("fill and sprite_sheet_name fields are exclusive")
	}

	// Sprite is included in sprite sheet
	if spriteRenderData.SpriteSheetName != "" {
		// Add reference to sprite sheet
		spriteSheet, ok := (*world.Resources.SpriteSheets)[spriteRenderData.SpriteSheetName]
		if !ok {
			utils.LogFatalf("unable to find sprite sheet with name '%s'", spriteRenderData.SpriteSheetName)
		}
		return &c.SpriteRender{
			SpriteSheet:  &spriteSheet,
			SpriteNumber: spriteRenderData.SpriteNumber,
		}
	}

	// Sprite is a colored rectangle
	textureImage := ebiten.NewImage(spriteRenderData.Fill.Width, spriteRenderData.Fill.Height)
	textureImage.Fill(color.RGBA{
		R: spriteRenderData.Fill.Color[0],
		G: spriteRenderData.Fill.Color[1],
		B: spriteRenderData.Fill.Color[2],
		A: spriteRenderData.Fill.Color[3],
	})

	return &c.SpriteRender{
		SpriteSheet: &c.SpriteSheet{
			Texture: c.Texture{Image: textureImage},
			Sprites: []c.Sprite{{X: 0, Y: 0, Width: spriteRenderData.Fill.Width, Height: spriteRenderData.Fill.Height}},
		},
		SpriteNumber: 0,
	}
}
