package loader

import (
	"log"
	"os"

	"github.com/kijimaD/sokotwo/lib/engine/loader"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
)

func PreloadEntities(entityMetadataPath string, world w.World) loader.EntityComponentList {
	b, err := os.ReadFile(entityMetadataPath)
	if err != nil {
		log.Fatal(b)
	}
	return loader.EntityComponentList{
		Engine: loader.LoadEngineComponents(b, world),
	}
}
