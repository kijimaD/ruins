package loader

import (
	"log"

	"github.com/kijimaD/sokotwo/assets"
	gc "github.com/kijimaD/sokotwo/lib/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
)

type gameComponentList struct {
	GridElement *gc.GridElement
	Player      *gc.Player
	Wall        *gc.Wall
	Warp        *gc.Warp
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
