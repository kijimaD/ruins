package raw

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/ruins/assets"
	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/utils"
	gloader "github.com/kijimaD/ruins/lib/loader"
)

type RawMaster struct {
	Raws          Raws
	ItemIndex     map[string]int
	MaterialIndex map[string]int
	RecipeIndex   map[string]int
	MemberIndex   map[string]int
}

type Raws struct {
	Items     []Item     `toml:"item"`
	Materials []Material `toml:"material"`
	Recipes   []Recipe   `toml:"recipe"`
	Members   []Member   `toml:"member"`
}

// tomlで入力として受け取る項目
type Item struct {
	Name            string
	Description     string
	ProvidesHealing int
	InflictsDamage  int
	Consumable      *Consumable `toml:"consumable"`
	Weapon          *Weapon     `toml:"weapon"`
}

type Consumable struct {
	UsableScene   string
	TargetFaction string
	TargetNum     string
}

type Weapon struct {
	Accuracy          int // 命中率。0~100%
	BaseDamage        int // ベース攻撃力
	AttackCount       int // 攻撃回数
	EnergyConsumption int // 攻撃で消費するエネルギー
	DamageAttr        string
}

type Material struct {
	Name        string
	Description string
}

type Member struct {
	Name string
	HP   int
	SP   int
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

	return rw
}

func (rw *RawMaster) GenerateItem(name string, spawnType SpawnType) gloader.GameComponentList {
	itemIdx, ok := rw.ItemIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	item := rw.Raws.Items[itemIdx]
	cl := gloader.GameComponentList{}
	if spawnType == SpawnInBackpack {
		cl.InBackpack = &gc.InBackpack{}
	}
	cl.Item = &gc.Item{}
	cl.Name = &gc.Name{Name: item.Name}
	cl.Description = &gc.Description{Description: item.Description}

	if item.Consumable != nil {
		var faction gc.TargetFactionType
		switch gc.TargetFactionType(item.Consumable.TargetFaction) {
		case gc.TargetFactionAlly:
			faction = gc.TargetFactionAlly
		case gc.TargetFactionEnemy:
			faction = gc.TargetFactionEnemy
		case gc.TargetFactionNone:
			faction = gc.TargetFactionNone
		default:
			log.Fatalf("invalid TargetFaction: %s", item.Consumable.TargetFaction)
		}

		var usableContext gc.UsableSceneType
		switch gc.UsableSceneType(item.Consumable.UsableScene) {
		case gc.UsableSceneAny:
			usableContext = gc.UsableSceneAny
		case gc.UsableSceneBattle:
			usableContext = gc.UsableSceneBattle
		case gc.UsableSceneField:
			usableContext = gc.UsableSceneField
		default:
			log.Fatalf("invalid UsableScene: %s", item.Consumable.UsableScene)
		}

		var targetCount gc.TargetNum
		switch gc.TargetNum(item.Consumable.TargetNum) {
		case gc.TargetSingle:
			targetCount = gc.TargetSingle
		case gc.TargetAll:
			targetCount = gc.TargetAll
		default:
			log.Fatalf("invalid TargetNum: %s", item.Consumable.TargetNum)
		}

		targetType := gc.TargetType{
			TargetFaction: faction,
			TargetNum:     targetCount,
		}
		cl.Consumable = &gc.Consumable{
			UsableScene: usableContext,
			TargetType:  targetType,
		}
	}
	if item.ProvidesHealing != 0 {
		cl.ProvidesHealing = &gc.ProvidesHealing{Amount: item.ProvidesHealing}
	}
	if item.InflictsDamage != 0 {
		cl.InflictsDamage = &gc.InflictsDamage{Amount: item.InflictsDamage}
	}
	if item.Weapon != nil {
		cl.Weapon = &gc.Weapon{
			Accuracy:          item.Weapon.Accuracy,
			BaseDamage:        item.Weapon.BaseDamage,
			AttackCount:       item.Weapon.AttackCount,
			EnergyConsumption: item.Weapon.EnergyConsumption,
			DamageAttr:        components.StringToDamangeAttrType(item.Weapon.DamageAttr),
		}
	}
	return cl
}

func (rw *RawMaster) GenerateMaterial(name string, amount int, spawnType SpawnType) gloader.GameComponentList {
	materialIdx, ok := rw.MaterialIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	cl := gloader.GameComponentList{}
	cl.Material = &gc.Material{Amount: amount}
	material := rw.Raws.Materials[materialIdx]
	cl.Name = &gc.Name{Name: material.Name}
	cl.Description = &gc.Description{Description: material.Description}
	if spawnType == SpawnInBackpack {
		cl.InBackpack = &gc.InBackpack{}
	}

	return cl
}

func (rw *RawMaster) GenerateRecipe(name string) gloader.GameComponentList {
	recipeIdx, ok := rw.RecipeIndex[name]
	if !ok {
		log.Fatalf("キーが存在しない: %s", name)
	}
	recipe := rw.Raws.Recipes[recipeIdx]
	cl := gloader.GameComponentList{}
	cl.Name = &gc.Name{Name: recipe.Name}
	cl.Recipe = &gc.Recipe{}
	for _, input := range recipe.Inputs {
		cl.Recipe.Inputs = append(cl.Recipe.Inputs, gc.RecipeInput{Name: input.Name, Amount: input.Amount})
	}

	// マッチしたitemの定義から持ってくる
	item := rw.GenerateItem(recipe.Name, SpawnInBackpack)
	cl.Description = &gc.Description{Description: item.Description.Description}
	if item.Weapon != nil {
		cl.Weapon = item.Weapon
	}
	if item.Consumable != nil {
		cl.Consumable = item.Consumable
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
