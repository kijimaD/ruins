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
	Raws        Raws
	ItemIndex   map[string]int
	MemberIndex map[string]int
}

type Raws struct {
	Items   []Item   `toml:"item"`
	Members []Member `toml:"member"`
}

// tomlで入力として受け取る項目
type Item struct {
	Name        string
	Description string
	Consumable  bool
}

type Member struct {
	Name string
	HP   int
	SP   int
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
	rw.MemberIndex = map[string]int{}
	utils.Try(toml.Decode(string(entityMetadataContent), &rw.Raws))

	for i, item := range rw.Raws.Items {
		rw.ItemIndex[item.Name] = i
	}
	for i, member := range rw.Raws.Members {
		rw.MemberIndex[member.Name] = i
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
	cl.InBackpack = &gc.InBackpack{} // フィールドにスポーンするときもあるので、引数で変えられるようにする
	cl.Item = &gc.Item{}
	cl.Name = &gc.Name{Name: item.Name}
	cl.Description = &gc.Description{Description: item.Description}
	if item.Consumable {
		cl.Consumable = &gc.Consumable{}
	}

	return cl
}

func (rw *RawMaster) GenerateMember(name string, inParty bool) gloader.GameComponentList {
	memberIdx, ok := rw.MemberIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	member := rw.Raws.Members[memberIdx]
	cl := gloader.GameComponentList{}
	cl.Member = &gc.Member{}
	cl.Name = &gc.Name{Name: member.Name}
	if inParty {
		cl.InParty = &gc.InParty{}
	}
	cl.Pools = &gc.Pools{
		HP:    gc.Pool{Max: member.HP, Current: member.HP - 10},
		SP:    gc.Pool{Max: member.SP, Current: member.SP - 10},
		Level: 1,
	}

	return cl
}
