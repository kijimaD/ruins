package loader

import (
	"log"
	"os"

	"github.com/x-hgg-x/goecsengine/loader"
	w "github.com/x-hgg-x/goecsengine/world"
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
