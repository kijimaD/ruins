package resources

import "github.com/kijimaD/sokotwo/lib/engine/loader"

type MenuPrefabs struct {
	MainMenu loader.EntityComponentList
}

type Prefabs struct {
	Menu  MenuPrefabs
	Intro loader.EntityComponentList
}
