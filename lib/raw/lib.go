package raw

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/sokotwo/assets"
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

func LoadFromFile(path string) RawMaster {
	bs, err := assets.FS.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	rw := Load(string(bs))
	return rw
}

func Load(entityMetadataContent string) RawMaster {
	rw := RawMaster{}
	rw.ItemIndex = map[string]int{}
	utils.Try(toml.Decode(string(entityMetadataContent), &rw.Raws))

	for i, item := range rw.Raws.Items {
		rw.ItemIndex[item.Name] = i
	}
	return rw
}

func (rw *RawMaster) GenerateItem(name string) gloader.GameComponentList {
	itemIdx, ok := rw.ItemIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	item := rw.Raws.Items[itemIdx]
	cl := gloader.GameComponentList{}
	cl.Item = &gc.Item{}
	cl.Name = &gc.Name{Name: item.Name}

	return cl
}
