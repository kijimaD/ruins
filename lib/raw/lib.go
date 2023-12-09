package raw

import (
	"github.com/BurntSushi/toml"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
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
