package resources

import "github.com/kijimaD/ruins/lib/engine/loader"

type FieldPrefabs struct {
	LevelInfo   loader.EntityComponentList
	PackageInfo loader.EntityComponentList
}

type Prefabs struct {
	Intro loader.EntityComponentList
	Field FieldPrefabs
}
