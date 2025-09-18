package raw

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/ruins/assets"
	gc "github.com/kijimaD/ruins/lib/components"
)

// Master はローデータを管理し、効率的な検索のためのインデックスを提供する
type Master struct {
	Raws              Raws
	ItemIndex         map[string]int
	MaterialIndex     map[string]int
	RecipeIndex       map[string]int
	MemberIndex       map[string]int
	CommandTableIndex map[string]int
	DropTableIndex    map[string]int
	SpriteSheetIndex  map[string]int
}

// Raws は全てのローデータを格納する構造体
type Raws struct {
	Items         []Item         `toml:"item"`
	Materials     []Material     `toml:"material"`
	Recipes       []Recipe       `toml:"recipe"`
	Members       []Member       `toml:"member"`
	CommandTables []CommandTable `toml:"command_table"`
	DropTables    []DropTable    `toml:"drop_table"`
	SpriteSheets  []SpriteSheet  `toml:"sprite_sheet"`
}

// Item はアイテムのローデータ
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

// ProvidesHealing は回復効果を提供する構造体
type ProvidesHealing struct {
	ValueType ValueType
	Amount    int
	Ratio     float64
}

// Consumable は消費可能なアイテムの設定
type Consumable struct {
	UsableScene string
	TargetGroup string
	TargetNum   string
}

// Card はカードアイテムの設定
type Card struct {
	Cost        int
	TargetGroup string
	TargetNum   string
}

// Attack は攻撃性能の設定
type Attack struct {
	Accuracy       int    // 命中率
	Damage         int    // 攻撃力
	AttackCount    int    // 攻撃回数
	Element        string // 攻撃属性
	AttackCategory string // 攻撃種別
}

// Wearable は装備可能アイテムの設定
type Wearable struct {
	Defense           int
	EquipmentCategory string
}

// EquipBonus は装備ボーナスの設定
type EquipBonus struct {
	Vitality  int
	Strength  int
	Sensation int
	Dexterity int
	Agility   int
}

// Material は素材アイテムの情報
type Material struct {
	Name        string
	Description string
}

// Recipe はレシピの情報
type Recipe struct {
	Name   string
	Inputs []RecipeInput `toml:"inputs"`
}

// RecipeInput は合成の元になる素材
type RecipeInput struct {
	Name   string
	Amount int
}

// Member はメンバーの情報
type Member struct {
	Name       string
	Job        string
	Player     *bool
	Attributes Attributes `toml:"attributes"`
}

// Attributes はキャラクターの能力値
type Attributes struct {
	Vitality  int
	Strength  int
	Sensation int
	Dexterity int
	Agility   int
	Defense   int
}

// LoadFromFile はファイルからローデータを読み込む
func LoadFromFile(path string) (Master, error) {
	bs, err := assets.FS.ReadFile(path)
	if err != nil {
		return Master{}, err
	}
	rw, err := Load(string(bs))
	if err != nil {
		return Master{}, err
	}
	return rw, nil
}

// Load は文字列からローデータを読み込む
func Load(entityMetadataContent string) (Master, error) {
	rw := Master{}
	rw.ItemIndex = map[string]int{}
	rw.MaterialIndex = map[string]int{}
	rw.RecipeIndex = map[string]int{}
	rw.MemberIndex = map[string]int{}
	rw.CommandTableIndex = map[string]int{}
	rw.DropTableIndex = map[string]int{}
	rw.SpriteSheetIndex = map[string]int{}
	_, err := toml.Decode(entityMetadataContent, &rw.Raws)
	if err != nil {
		return Master{}, err
	}

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
	for i, dropTable := range rw.Raws.DropTables {
		rw.DropTableIndex[dropTable.Name] = i
	}
	for i, spriteSheet := range rw.Raws.SpriteSheets {
		rw.SpriteSheetIndex[spriteSheet.Name] = i
	}

	return rw, nil
}

// GenerateItem は指定された名前のアイテムのゲームコンポーネントを生成する
func (rw *Master) GenerateItem(name string, locationType gc.ItemLocationType) (gc.GameComponentList, error) {
	itemIdx, ok := rw.ItemIndex[name]
	if !ok {
		return gc.GameComponentList{}, NewKeyNotFoundError(name, "ItemIndex")
	}
	item := rw.Raws.Items[itemIdx]
	cl := gc.GameComponentList{}
	cl.ItemLocationType = &locationType
	cl.Item = &gc.Item{}
	cl.Name = &gc.Name{Name: item.Name}
	cl.Description = &gc.Description{Description: item.Description}

	if item.Consumable != nil {
		if err := gc.TargetGroupType(item.Consumable.TargetGroup).Valid(); err != nil {
			return gc.GameComponentList{}, fmt.Errorf("%s: %w", "invalid target group type", err)
		}
		if err := gc.TargetNumType(item.Consumable.TargetNum).Valid(); err != nil {
			return gc.GameComponentList{}, fmt.Errorf("%s: %w", "invalid target num type", err)
		}
		targetType := gc.TargetType{
			TargetGroup: gc.TargetGroupType(item.Consumable.TargetGroup),
			TargetNum:   gc.TargetNumType(item.Consumable.TargetNum),
		}

		if err := gc.UsableSceneType(item.Consumable.UsableScene).Valid(); err != nil {
			return gc.GameComponentList{}, fmt.Errorf("%s: %w", "invalid usable scene type", err)
		}
		cl.Consumable = &gc.Consumable{
			UsableScene: gc.UsableSceneType(item.Consumable.UsableScene),
			TargetType:  targetType,
		}
	}

	if item.ProvidesHealing != nil {
		if err := item.ProvidesHealing.ValueType.Valid(); err != nil {
			return gc.GameComponentList{}, fmt.Errorf("%s: %w", "invalid value type", err)
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
			return gc.GameComponentList{}, fmt.Errorf("%s: %w", "invalid card target group type", err)
		}
		if err := gc.TargetNumType(item.Card.TargetNum).Valid(); err != nil {
			return gc.GameComponentList{}, fmt.Errorf("%s: %w", "invalid card target num type", err)
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
			return gc.GameComponentList{}, err
		}
		if err := gc.AttackType(item.Attack.AttackCategory).Valid(); err != nil {
			return gc.GameComponentList{}, err
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
		if err := gc.EquipmentType(item.Wearable.EquipmentCategory).Valid(); err != nil {
			return gc.GameComponentList{}, err
		}
		cl.Wearable = &gc.Wearable{
			Defense:           item.Wearable.Defense,
			EquipmentCategory: gc.EquipmentType(item.Wearable.EquipmentCategory),
			EquipBonus:        bonus,
		}
	}

	return cl, nil
}

// GenerateMaterial は指定された名前の素材のゲームコンポーネントを生成する
func (rw *Master) GenerateMaterial(name string, amount int, locationType gc.ItemLocationType) (gc.GameComponentList, error) {
	materialIdx, ok := rw.MaterialIndex[name]
	if !ok {
		return gc.GameComponentList{}, NewKeyNotFoundError(name, "MaterialIndex")
	}
	cl := gc.GameComponentList{}
	cl.Material = &gc.Material{Amount: amount}
	material := rw.Raws.Materials[materialIdx]
	cl.Name = &gc.Name{Name: material.Name}
	cl.Description = &gc.Description{Description: material.Description}
	cl.ItemLocationType = &locationType

	return cl, nil
}

// GenerateRecipe は指定された名前のレシピのゲームコンポーネントを生成する
func (rw *Master) GenerateRecipe(name string) (gc.GameComponentList, error) {
	recipeIdx, ok := rw.RecipeIndex[name]
	if !ok {
		return gc.GameComponentList{}, NewKeyNotFoundError(name, "RecipeIndex")
	}
	recipe := rw.Raws.Recipes[recipeIdx]
	cl := gc.GameComponentList{}
	cl.Name = &gc.Name{Name: recipe.Name}
	cl.Recipe = &gc.Recipe{}
	for _, input := range recipe.Inputs {
		cl.Recipe.Inputs = append(cl.Recipe.Inputs, gc.RecipeInput{Name: input.Name, Amount: input.Amount})
	}

	// 説明文や分類のため、マッチしたitemの定義から持ってくる
	item, err := rw.GenerateItem(recipe.Name, gc.ItemLocationInBackpack)
	if err != nil {
		return gc.GameComponentList{}, fmt.Errorf("%s: %w", "failed to generate item for recipe", err)
	}
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

	return cl, nil
}

// generateFighter は指定された名前の戦闘員のゲームコンポーネントを生成する(敵・味方共通)
func (rw *Master) generateFighter(name string) (gc.GameComponentList, error) {
	memberIdx, ok := rw.MemberIndex[name]
	if !ok {
		return gc.GameComponentList{}, fmt.Errorf("キーが存在しない: %s", name)
	}
	member := rw.Raws.Members[memberIdx]

	cl := gc.GameComponentList{}
	cl.Name = &gc.Name{Name: member.Name}
	if member.Job != "" {
		cl.Job = &gc.Job{Job: member.Job}
	}
	cl.TurnBased = &gc.TurnBased{AP: gc.Pool{Current: 100, Max: 100}} // TODO: Attributesから計算する
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
	if member.Player != nil && *member.Player {
		cl.Player = &gc.Player{}
	}

	commandTableIdx, ok := rw.CommandTableIndex[name]
	if ok {
		commandTable := rw.Raws.CommandTables[commandTableIdx]
		cl.CommandTable = &gc.CommandTable{Name: commandTable.Name}
	}

	dropTableIdx, ok := rw.DropTableIndex[name]
	if ok {
		dropTable := rw.Raws.DropTables[dropTableIdx]
		cl.DropTable = &gc.DropTable{Name: dropTable.Name}
	}

	return cl, nil
}

// GeneratePlayer は指定された名前のプレイヤーのゲームコンポーネントを生成する
func (rw *Master) GeneratePlayer(name string) (gc.GameComponentList, error) {
	cl, err := rw.generateFighter(name)
	if err != nil {
		return gc.GameComponentList{}, err
	}
	cl.FactionType = &gc.FactionAlly
	cl.Player = &gc.Player{}
	cl.Operator = &gc.Operator{}
	cl.Hunger = gc.NewHunger()
	return cl, nil
}

// GenerateEnemy は指定された名前の敵のゲームコンポーネントを生成する
func (rw *Master) GenerateEnemy(name string) (gc.GameComponentList, error) {
	cl, err := rw.generateFighter(name)
	if err != nil {
		return gc.GameComponentList{}, err
	}
	cl.FactionType = &gc.FactionEnemy

	return cl, nil
}

// GetCommandTable は指定された名前のコマンドテーブルを取得する
func (rw *Master) GetCommandTable(name string) (CommandTable, error) {
	ctIdx, ok := rw.CommandTableIndex[name]
	if !ok {
		return CommandTable{}, fmt.Errorf("キーが存在しない: %s", name)
	}
	commandTable := rw.Raws.CommandTables[ctIdx]

	return commandTable, nil
}

// GetDropTable は指定された名前のドロップテーブルを取得する
func (rw *Master) GetDropTable(name string) (DropTable, error) {
	dtIdx, ok := rw.DropTableIndex[name]
	if !ok {
		return DropTable{}, fmt.Errorf("キーが存在しない: %s", name)
	}
	dropTable := rw.Raws.DropTables[dtIdx]

	return dropTable, nil
}
