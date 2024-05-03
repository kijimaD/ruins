package resources

import "github.com/kijimaD/ruins/lib/engine/loader"

type MenuPrefabs struct {
	InventoryMenu loader.EntityComponentList
	CraftMenu     loader.EntityComponentList
	EquipMenu     loader.EntityComponentList
}

type FieldPrefabs struct {
	LevelInfo   loader.EntityComponentList
	PackageInfo loader.EntityComponentList
}

type Prefabs struct {
	Menu  MenuPrefabs
	Intro loader.EntityComponentList
	Field FieldPrefabs
}
