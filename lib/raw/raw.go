package raw

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/ruins/assets"
	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/utils"
)

type RawMaster struct {
	Raws              Raws
	ItemIndex         map[string]int
	MaterialIndex     map[string]int
	RecipeIndex       map[string]int
	MemberIndex       map[string]int
	CommandTableIndex map[string]int
}

type Raws struct {
	Items         []Item         `toml:"item"`
	Materials     []Material     `toml:"material"`
	Recipes       []Recipe       `toml:"recipe"`
	Members       []Member       `toml:"member"`
	CommandTables []CommandTable `toml:"command_table"`
}

type Item struct {
	Name            string
	Description     string
	InflictsDamage  int
	Consumable      *Consumable      `toml:"consumable"`
	ProvidesHealing *ProvidesHealing `toml:"provides_healing"`
	Wearable        *Wearable        `toml:"wearable"`
	EquipBonus      *EquipBonus      `toml:"equip_bonus"`
	Card            *Card            `toml:"card"`
	Attack          *Attack          `toml:"attack"`
}

type ProvidesHealing struct {
	ValueType ValueType
	Amount    int
	Ratio     float64
}

type Consumable struct {
	UsableScene string
	TargetGroup string
	TargetNum   string
}

type Card struct {
	Cost        int
	TargetGroup string
	TargetNum   string
}

type Attack struct {
	Accuracy       int    // 命中率
	Damage         int    // 攻撃力
	AttackCount    int    // 攻撃回数
	Element        string // 攻撃属性
	AttackCategory string // 攻撃種別
}

type Wearable struct {
	Defense           int
	EquipmentCategory string
}

type EquipBonus struct {
	Vitality  int
	Strength  int
	Sensation int
	Dexterity int
	Agility   int
}

type Material struct {
	Name        string
	Description string
}

type Member struct {
	Name       string
	Attributes Attributes `toml:"attributes"`
}

type CommandTable struct {
	Name  string
	Table []CommandTableEntry `toml:"table"`
}

type CommandTableEntry struct {
	Card   string
	Weight float64
}

type Attributes struct {
	Vitality  int
	Strength  int
	Sensation int
	Dexterity int
	Agility   int
	Defense   int
}

// レシピ
type Recipe struct {
	Name   string
	Inputs []RecipeInput `toml:"inputs"`
}

// 合成の元になる素材
type RecipeInput struct {
	Name   string
	Amount int
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
	rw.MaterialIndex = map[string]int{}
	rw.RecipeIndex = map[string]int{}
	rw.MemberIndex = map[string]int{}
	rw.CommandTableIndex = map[string]int{}
	utils.Try(toml.Decode(string(entityMetadataContent), &rw.Raws))

	for i, item := range rw.Raws.Items {
		rw.ItemIndex[item.Name] = i
	}
	for i, material := range rw.Raws.Materials {
		rw.MaterialIndex[material.Name] = i
	}
	for i, recipe := range rw.Raws.Recipes {
		rw.RecipeIndex[recipe.Name] = i
	}
	for i, member := range rw.Raws.Members {
		rw.MemberIndex[member.Name] = i
	}
	for i, commandTable := range rw.Raws.CommandTables {
		rw.CommandTableIndex[commandTable.Name] = i
	}

	return rw
}

func (rw *RawMaster) GenerateItem(name string, locationType gc.ItemLocationType) components.GameComponentList {
	itemIdx, ok := rw.ItemIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	item := rw.Raws.Items[itemIdx]
	cl := components.GameComponentList{}
	cl.ItemLocationType = &locationType
	cl.Item = &gc.Item{}
	cl.Name = &gc.Name{Name: item.Name}
	cl.Description = &gc.Description{Description: item.Description}

	if item.Consumable != nil {
		if err := gc.TargetGroupType(item.Consumable.TargetGroup).Valid(); err != nil {
			log.Fatal(err)
		}
		if err := gc.TargetNumType(item.Consumable.TargetNum).Valid(); err != nil {
			log.Fatal(err)
		}
		targetType := gc.TargetType{
			TargetGroup: gc.TargetGroupType(item.Consumable.TargetGroup),
			TargetNum:   gc.TargetNumType(item.Consumable.TargetNum),
		}

		if err := gc.UsableSceneType(item.Consumable.UsableScene).Valid(); err != nil {
			log.Fatal(err)
		}
		cl.Consumable = &gc.Consumable{
			UsableScene: gc.UsableSceneType(item.Consumable.UsableScene),
			TargetType:  targetType,
		}
	}

	if item.ProvidesHealing != nil {
		if err := ValueType(item.ProvidesHealing.ValueType).Valid(); err != nil {
			log.Fatal(err)
		}
		switch item.ProvidesHealing.ValueType {
		case PercentageType:
			cl.ProvidesHealing = &gc.ProvidesHealing{Amount: gc.RatioAmount{Ratio: item.ProvidesHealing.Ratio}}
		case NumeralType:
			cl.ProvidesHealing = &gc.ProvidesHealing{Amount: gc.NumeralAmount{Numeral: item.ProvidesHealing.Amount}}
		}
	}
	if item.InflictsDamage != 0 {
		cl.InflictsDamage = &gc.InflictsDamage{Amount: item.InflictsDamage}
	}

	if item.Card != nil {
		if err := gc.TargetGroupType(item.Card.TargetGroup).Valid(); err != nil {
			log.Fatal(err)
		}
		if err := gc.TargetNumType(item.Card.TargetNum).Valid(); err != nil {
			log.Fatal(err)
		}

		cl.Card = &gc.Card{
			TargetType: gc.TargetType{
				TargetGroup: gc.TargetGroupType(item.Card.TargetGroup),
				TargetNum:   gc.TargetNumType(item.Card.TargetNum),
			},
			Cost: item.Card.Cost,
		}
	}

	if item.Attack != nil {
		if err := gc.ElementType(item.Attack.Element).Valid(); err != nil {
			log.Fatal(err)
		}
		if err := gc.AttackType(item.Attack.AttackCategory).Valid(); err != nil {
			log.Fatal(err)
		}

		cl.Attack = &gc.Attack{
			Accuracy:       item.Attack.Accuracy,
			Damage:         item.Attack.Damage,
			AttackCount:    item.Attack.AttackCount,
			Element:        gc.ElementType(item.Attack.Element),
			AttackCategory: gc.AttackType(item.Attack.AttackCategory),
		}
	}

	var bonus gc.EquipBonus
	if item.EquipBonus != nil {
		bonus = gc.EquipBonus{
			Vitality:  item.EquipBonus.Vitality,
			Strength:  item.EquipBonus.Strength,
			Sensation: item.EquipBonus.Sensation,
			Dexterity: item.EquipBonus.Dexterity,
			Agility:   item.EquipBonus.Agility,
		}
	}

	if item.Wearable != nil {
		if err := components.EquipmentType(item.Wearable.EquipmentCategory).Valid(); err != nil {
			log.Fatal(err)
		}
		cl.Wearable = &gc.Wearable{
			Defense:           item.Wearable.Defense,
			EquipmentCategory: components.EquipmentType(item.Wearable.EquipmentCategory),
			EquipBonus:        bonus,
		}
	}

	return cl
}

func (rw *RawMaster) GenerateMaterial(name string, amount int, locationType gc.ItemLocationType) components.GameComponentList {
	materialIdx, ok := rw.MaterialIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	cl := components.GameComponentList{}
	cl.Material = &gc.Material{Amount: amount}
	material := rw.Raws.Materials[materialIdx]
	cl.Name = &gc.Name{Name: material.Name}
	cl.Description = &gc.Description{Description: material.Description}
	cl.ItemLocationType = &locationType

	return cl
}

func (rw *RawMaster) GenerateRecipe(name string) components.GameComponentList {
	recipeIdx, ok := rw.RecipeIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	recipe := rw.Raws.Recipes[recipeIdx]
	cl := components.GameComponentList{}
	cl.Name = &gc.Name{Name: recipe.Name}
	cl.Recipe = &gc.Recipe{}
	for _, input := range recipe.Inputs {
		cl.Recipe.Inputs = append(cl.Recipe.Inputs, gc.RecipeInput{Name: input.Name, Amount: input.Amount})
	}

	// 説明文などのため、マッチしたitemの定義から持ってくる
	item := rw.GenerateItem(recipe.Name, gc.ItemLocationInBackpack)
	cl.Description = &gc.Description{Description: item.Description.Description}
	if item.Card != nil {
		cl.Card = item.Card
	}
	if item.Attack != nil {
		cl.Attack = item.Attack
	}
	if item.Wearable != nil {
		cl.Wearable = item.Wearable
	}
	if item.Consumable != nil {
		cl.Consumable = item.Consumable
	}

	return cl
}

func (rw *RawMaster) GenerateFighter(name string) components.GameComponentList {
	memberIdx, ok := rw.MemberIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	member := rw.Raws.Members[memberIdx]

	cl := components.GameComponentList{}
	cl.Name = &gc.Name{Name: member.Name}
	cl.Attributes = &gc.Attributes{
		Vitality:  gc.Attribute{Base: member.Attributes.Vitality},
		Strength:  gc.Attribute{Base: member.Attributes.Strength},
		Sensation: gc.Attribute{Base: member.Attributes.Sensation},
		Dexterity: gc.Attribute{Base: member.Attributes.Dexterity},
		Agility:   gc.Attribute{Base: member.Attributes.Agility},
		Defense:   gc.Attribute{Base: member.Attributes.Defense},
	}
	cl.Pools = &gc.Pools{
		Level: 1,
	}
	cl.EquipmentChanged = &gc.EquipmentChanged{}

	commandTableIdx, ok := rw.CommandTableIndex[name]
	if ok {
		commandTable := rw.Raws.CommandTables[commandTableIdx]
		cl.CommandTable = &gc.CommandTable{Name: commandTable.Name}
	}

	return cl
}

func (rw *RawMaster) GenerateMember(name string, inParty bool) components.GameComponentList {
	cl := rw.GenerateFighter(name)
	cl.FactionType = &gc.FactionAlly
	if inParty {
		cl.InParty = &gc.InParty{}
	}

	return cl
}

func (rw *RawMaster) GenerateEnemy(name string) components.GameComponentList {
	cl := rw.GenerateFighter(name)
	cl.FactionType = &gc.FactionEnemy

	return cl
}

func (rw *RawMaster) GetCommandTable(name string) CommandTable {
	ctIdx, ok := rw.CommandTableIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	commandTable := rw.Raws.CommandTables[ctIdx]

	return commandTable
}
