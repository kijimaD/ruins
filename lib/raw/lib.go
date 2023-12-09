package raw

import (
	"log"

	"github.com/BurntSushi/toml"
	gc "github.com/kijimaD/sokotwo/lib/components"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	gloader "github.com/kijimaD/sokotwo/lib/loader"
)

type RawMaster struct {
	Raws      Raws
	ItemIndex map[string]int // repair: 0 みたいな
}

type Raws struct {
	Items []Item `toml:"item"`
}

type Item struct {
	Name string
}

func (rw *RawMaster) Load(entityMetadataContent string) {
	rw.ItemIndex = map[string]int{}
	utils.Try(toml.Decode(string(entityMetadataContent), &rw.Raws))

	for i, item := range rw.Raws.Items {
		rw.ItemIndex[item.Name] = i
	}
}

func (rw *RawMaster) GenerateItem(name string) gloader.Entity {
	itemIdx, ok := rw.ItemIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	item := rw.Raws.Items[itemIdx]
	entity := gloader.Entity{}
	entity.Components.Item = &gc.Item{}
	entity.Components.Name = &gc.Name{Name: item.Name}

	return entity
}
