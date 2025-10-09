package raw

import (
	"fmt"
	"image/color"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/ruins/assets"
	gc "github.com/kijimaD/ruins/lib/components"
)

// Master はローデータを管理し、効率的な検索のためのインデックスを提供する
type Master struct {
	Raws              Raws
	ItemIndex         map[string]int
	RecipeIndex       map[string]int
	MemberIndex       map[string]int
	CommandTableIndex map[string]int
	DropTableIndex    map[string]int
	SpriteSheetIndex  map[string]int
	TileIndex         map[string]int
	PropIndex         map[string]int
}

// Raws は全てのローデータを格納する構造体
type Raws struct {
	Items         []Item
	Recipes       []Recipe
	Members       []Member
	CommandTables []CommandTable
	DropTables    []DropTable
	SpriteSheets  []SpriteSheet
	Tiles         []TileRaw
	Props         []PropRaw
}

// Item はアイテムのローデータ
type Item struct {
	Name            string
	Description     string
	SpriteSheetName string
	SpriteKey       string
	Value           *int
	InflictsDamage  *int
	Stackable       *bool // スタック可能かどうか
	Consumable      *Consumable
	ProvidesHealing *ProvidesHealing
	Wearable        *Wearable
	EquipBonus      *EquipBonus
	Card            *Card
	Attack          *Attack
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

// Recipe はレシピの情報
type Recipe struct {
	Name   string
	Inputs []RecipeInput
}

// RecipeInput は合成の元になる素材
type RecipeInput struct {
	Name   string
	Amount int
}

// Member はメンバーの情報
type Member struct {
	Name            string
	Player          *bool
	Attributes      Attributes
	SpriteSheetName string
	SpriteKey       string
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
	rw.RecipeIndex = map[string]int{}
	rw.MemberIndex = map[string]int{}
	rw.CommandTableIndex = map[string]int{}
	rw.DropTableIndex = map[string]int{}
	rw.SpriteSheetIndex = map[string]int{}
	rw.TileIndex = map[string]int{}
	rw.PropIndex = map[string]int{}

	metaData, err := toml.Decode(entityMetadataContent, &rw.Raws)
	if err != nil {
		return Master{}, fmt.Errorf("TOML decode error: %w", err)
	}
	// 未知のキーがあった場合はエラーにする
	undecoded := metaData.Undecoded()
	if len(undecoded) > 0 {
		return Master{}, fmt.Errorf("unknown keys found in TOML: %v", undecoded)
	}

	for i, item := range rw.Raws.Items {
		rw.ItemIndex[item.Name] = i
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
	for i, tile := range rw.Raws.Tiles {
		rw.TileIndex[tile.Name] = i
	}
	for i, prop := range rw.Raws.Props {
		rw.PropIndex[prop.Name] = i
	}

	return rw, nil
}

// NewItemSpec は指定された名前のアイテムのEntitySpecを生成する
func (rw *Master) NewItemSpec(name string, locationType *gc.ItemLocationType) (gc.EntitySpec, error) {
	itemIdx, ok := rw.ItemIndex[name]
	if !ok {
		return gc.EntitySpec{}, NewKeyNotFoundError(name, "ItemIndex")
	}
	item := rw.Raws.Items[itemIdx]

	entitySpec := gc.EntitySpec{}
	if locationType != nil {
		entitySpec.ItemLocationType = locationType
	}
	entitySpec.Item = &gc.Item{}
	entitySpec.Name = &gc.Name{Name: item.Name}
	entitySpec.Description = &gc.Description{Description: item.Description}
	entitySpec.SpriteRender = &gc.SpriteRender{
		SpriteSheetName: item.SpriteSheetName,
		SpriteKey:       item.SpriteKey,
		Depth:           gc.DepthNumRug,
	}

	if item.Consumable != nil {
		if err := gc.TargetGroupType(item.Consumable.TargetGroup).Valid(); err != nil {
			return gc.EntitySpec{}, fmt.Errorf("%s: %w", "invalid target group type", err)
		}
		if err := gc.TargetNumType(item.Consumable.TargetNum).Valid(); err != nil {
			return gc.EntitySpec{}, fmt.Errorf("%s: %w", "invalid target num type", err)
		}
		targetType := gc.TargetType{
			TargetGroup: gc.TargetGroupType(item.Consumable.TargetGroup),
			TargetNum:   gc.TargetNumType(item.Consumable.TargetNum),
		}

		if err := gc.UsableSceneType(item.Consumable.UsableScene).Valid(); err != nil {
			return gc.EntitySpec{}, fmt.Errorf("%s: %w", "invalid usable scene type", err)
		}
		entitySpec.Consumable = &gc.Consumable{
			UsableScene: gc.UsableSceneType(item.Consumable.UsableScene),
			TargetType:  targetType,
		}
	}

	if item.ProvidesHealing != nil {
		if err := item.ProvidesHealing.ValueType.Valid(); err != nil {
			return gc.EntitySpec{}, fmt.Errorf("%s: %w", "invalid value type", err)
		}
		switch item.ProvidesHealing.ValueType {
		case PercentageType:
			entitySpec.ProvidesHealing = &gc.ProvidesHealing{Amount: gc.RatioAmount{Ratio: item.ProvidesHealing.Ratio}}
		case NumeralType:
			entitySpec.ProvidesHealing = &gc.ProvidesHealing{Amount: gc.NumeralAmount{Numeral: item.ProvidesHealing.Amount}}
		}
	}
	if item.InflictsDamage != nil {
		entitySpec.InflictsDamage = &gc.InflictsDamage{Amount: *item.InflictsDamage}
	}

	if item.Card != nil {
		if err := gc.TargetGroupType(item.Card.TargetGroup).Valid(); err != nil {
			return gc.EntitySpec{}, fmt.Errorf("%s: %w", "invalid card target group type", err)
		}
		if err := gc.TargetNumType(item.Card.TargetNum).Valid(); err != nil {
			return gc.EntitySpec{}, fmt.Errorf("%s: %w", "invalid card target num type", err)
		}

		entitySpec.Card = &gc.Card{
			TargetType: gc.TargetType{
				TargetGroup: gc.TargetGroupType(item.Card.TargetGroup),
				TargetNum:   gc.TargetNumType(item.Card.TargetNum),
			},
			Cost: item.Card.Cost,
		}
	}

	if item.Attack != nil {
		if err := gc.ElementType(item.Attack.Element).Valid(); err != nil {
			return gc.EntitySpec{}, err
		}
		if err := gc.AttackType(item.Attack.AttackCategory).Valid(); err != nil {
			return gc.EntitySpec{}, err
		}

		entitySpec.Attack = &gc.Attack{
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
			return gc.EntitySpec{}, err
		}
		entitySpec.Wearable = &gc.Wearable{
			Defense:           item.Wearable.Defense,
			EquipmentCategory: gc.EquipmentType(item.Wearable.EquipmentCategory),
			EquipBonus:        bonus,
		}
	}

	if item.Value != nil {
		entitySpec.Value = &gc.Value{Value: *item.Value}
	}

	return entitySpec, nil
}

// NewRecipeSpec は指定された名前のレシピのEntitySpecを生成する
func (rw *Master) NewRecipeSpec(name string) (gc.EntitySpec, error) {
	recipeIdx, ok := rw.RecipeIndex[name]
	if !ok {
		return gc.EntitySpec{}, NewKeyNotFoundError(name, "RecipeIndex")
	}
	recipe := rw.Raws.Recipes[recipeIdx]
	entitySpec := gc.EntitySpec{}
	entitySpec.Name = &gc.Name{Name: recipe.Name}
	entitySpec.Recipe = &gc.Recipe{}
	for _, input := range recipe.Inputs {
		entitySpec.Recipe.Inputs = append(entitySpec.Recipe.Inputs, gc.RecipeInput{Name: input.Name, Amount: input.Amount})
	}

	// 説明文や分類のため、マッチしたitemの定義から持ってくる
	// マスターデータのため位置を指定しない
	itemSpec, err := rw.NewItemSpec(recipe.Name, nil)
	if err != nil {
		return gc.EntitySpec{}, fmt.Errorf("%s: %w", "failed to generate item for recipe", err)
	}
	entitySpec.Description = &gc.Description{Description: itemSpec.Description.Description}
	if itemSpec.Card != nil {
		entitySpec.Card = itemSpec.Card
	}
	if itemSpec.Attack != nil {
		entitySpec.Attack = itemSpec.Attack
	}
	if itemSpec.Wearable != nil {
		entitySpec.Wearable = itemSpec.Wearable
	}
	if itemSpec.Consumable != nil {
		entitySpec.Consumable = itemSpec.Consumable
	}
	if itemSpec.Value != nil {
		entitySpec.Value = itemSpec.Value
	}

	return entitySpec, nil
}

// generateFighter は指定された名前の戦闘員のゲームコンポーネントを生成する(敵・味方共通)
func (rw *Master) generateFighter(name string) (gc.EntitySpec, error) {
	memberIdx, ok := rw.MemberIndex[name]
	if !ok {
		return gc.EntitySpec{}, fmt.Errorf("キーが存在しない: %s", name)
	}
	member := rw.Raws.Members[memberIdx]

	entitySpec := gc.EntitySpec{}
	entitySpec.Name = &gc.Name{Name: member.Name}
	entitySpec.TurnBased = &gc.TurnBased{AP: gc.Pool{Current: 100, Max: 100}} // TODO: Attributesから計算する
	entitySpec.SpriteRender = &gc.SpriteRender{
		SpriteSheetName: member.SpriteSheetName,
		SpriteKey:       member.SpriteKey,
		Depth:           gc.DepthNumPlayer,
	}
	entitySpec.Attributes = &gc.Attributes{
		Vitality:  gc.Attribute{Base: member.Attributes.Vitality},
		Strength:  gc.Attribute{Base: member.Attributes.Strength},
		Sensation: gc.Attribute{Base: member.Attributes.Sensation},
		Dexterity: gc.Attribute{Base: member.Attributes.Dexterity},
		Agility:   gc.Attribute{Base: member.Attributes.Agility},
		Defense:   gc.Attribute{Base: member.Attributes.Defense},
	}
	entitySpec.Pools = &gc.Pools{}
	entitySpec.EquipmentChanged = &gc.EquipmentChanged{}
	if member.Player != nil && *member.Player {
		entitySpec.Player = &gc.Player{}
	}

	commandTableIdx, ok := rw.CommandTableIndex[name]
	if ok {
		commandTable := rw.Raws.CommandTables[commandTableIdx]
		entitySpec.CommandTable = &gc.CommandTable{Name: commandTable.Name}
	}

	dropTableIdx, ok := rw.DropTableIndex[name]
	if ok {
		dropTable := rw.Raws.DropTables[dropTableIdx]
		entitySpec.DropTable = &gc.DropTable{Name: dropTable.Name}
	}

	return entitySpec, nil
}

// NewPlayerSpec は指定された名前のプレイヤーのEntitySpecを生成する
func (rw *Master) NewPlayerSpec(name string) (gc.EntitySpec, error) {
	entitySpec, err := rw.generateFighter(name)
	if err != nil {
		return gc.EntitySpec{}, err
	}
	entitySpec.FactionType = &gc.FactionAlly
	entitySpec.Player = &gc.Player{}
	entitySpec.Hunger = gc.NewHunger()
	entitySpec.LightSource = &gc.LightSource{
		Radius:  6,
		Color:   color.RGBA{R: 255, G: 200, B: 150, A: 255}, // ランタンの暖色光
		Enabled: true,
	}
	return entitySpec, nil
}

// NewEnemySpec は指定された名前の敵のEntitySpecを生成する
func (rw *Master) NewEnemySpec(name string) (gc.EntitySpec, error) {
	entitySpec, err := rw.generateFighter(name)
	if err != nil {
		return gc.EntitySpec{}, err
	}
	entitySpec.FactionType = &gc.FactionEnemy

	return entitySpec, nil
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

// TileRaw はタイルのローデータ定義
type TileRaw struct {
	Name        string
	Description string
	Walkable    bool
}

// PropRaw は置物のローデータ定義
type PropRaw struct {
	Name            string
	Description     string
	SpriteSheetName string
	SpriteKey       string
	BlockPass       bool
	BlockView       bool
	LightSource     *LightSourceRaw
}

// LightSourceRaw は光源のローデータ定義
type LightSourceRaw struct {
	Radius  gc.Tile
	Color   color.RGBA
	Enabled bool
}

// GetTile は指定された名前のタイルを取得する
func (rw *Master) GetTile(name string) (TileRaw, error) {
	tileIdx, ok := rw.TileIndex[name]
	if !ok {
		return TileRaw{}, NewKeyNotFoundError(name, "TileIndex")
	}

	return rw.Raws.Tiles[tileIdx], nil
}

// GetProp は指定された名前の置物の設定を取得する
func (rw *Master) GetProp(name string) (PropRaw, error) {
	propIdx, ok := rw.PropIndex[name]
	if !ok {
		return PropRaw{}, NewKeyNotFoundError(name, "PropIndex")
	}

	return rw.Raws.Props[propIdx], nil
}

// NewPropSpec は指定された名前の置物のEntitySpecを生成する
func (rw *Master) NewPropSpec(name string) (gc.EntitySpec, error) {
	propRaw, err := rw.GetProp(name)
	if err != nil {
		return gc.EntitySpec{}, err
	}

	entitySpec := gc.EntitySpec{}
	entitySpec.Prop = &gc.Prop{}
	entitySpec.Name = &gc.Name{Name: propRaw.Name}
	entitySpec.Description = &gc.Description{Description: propRaw.Description}
	entitySpec.SpriteRender = &gc.SpriteRender{
		SpriteSheetName: propRaw.SpriteSheetName,
		SpriteKey:       propRaw.SpriteKey,
		Depth:           gc.DepthNumRug,
	}

	if propRaw.BlockPass {
		entitySpec.BlockPass = &gc.BlockPass{}
	}
	if propRaw.BlockView {
		entitySpec.BlockView = &gc.BlockView{}
	}

	// 光源の設定
	if propRaw.LightSource != nil {
		entitySpec.LightSource = &gc.LightSource{
			Radius:  propRaw.LightSource.Radius,
			Color:   propRaw.LightSource.Color,
			Enabled: propRaw.LightSource.Enabled,
		}
	}

	return entitySpec, nil
}
