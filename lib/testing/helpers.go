package testing

import (
	"testing"

	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
)

// テスト用のプレイヤーコンポーネントを作成する
func CreateTestPlayer(t *testing.T) components.GameComponentList {
	t.Helper()

	return components.GameComponentList{
		Name:     &gc.Name{Name: "テストプレイヤー"},
		Position: &gc.Position{X: 100, Y: 100},
		Pools: &gc.Pools{
			HP:    gc.Pool{Current: 100, Max: 100},
			SP:    gc.Pool{Current: 50, Max: 50},
			XP:    0,
			Level: 1,
		},
		Attributes: &gc.Attributes{
			Vitality:  gc.Attribute{Base: 10, Modifier: 0, Total: 10},
			Strength:  gc.Attribute{Base: 10, Modifier: 0, Total: 10},
			Sensation: gc.Attribute{Base: 10, Modifier: 0, Total: 10},
			Dexterity: gc.Attribute{Base: 10, Modifier: 0, Total: 10},
			Agility:   gc.Attribute{Base: 10, Modifier: 0, Total: 10},
			Defense:   gc.Attribute{Base: 5, Modifier: 0, Total: 5},
		},
		InParty: &gc.InParty{},
	}
}

// テスト用の敵コンポーネントを作成する
func CreateTestEnemy(t *testing.T, name string) components.GameComponentList {
	t.Helper()

	return components.GameComponentList{
		Name:     &gc.Name{Name: name},
		Position: &gc.Position{X: 200, Y: 200},
		Pools: &gc.Pools{
			HP:    gc.Pool{Current: 50, Max: 50},
			SP:    gc.Pool{Current: 20, Max: 20},
			XP:    0,
			Level: 1,
		},
		Attributes: &gc.Attributes{
			Vitality:  gc.Attribute{Base: 8, Modifier: 0, Total: 8},
			Strength:  gc.Attribute{Base: 8, Modifier: 0, Total: 8},
			Sensation: gc.Attribute{Base: 8, Modifier: 0, Total: 8},
			Dexterity: gc.Attribute{Base: 8, Modifier: 0, Total: 8},
			Agility:   gc.Attribute{Base: 8, Modifier: 0, Total: 8},
			Defense:   gc.Attribute{Base: 3, Modifier: 0, Total: 3},
		},
		FactionType: &gc.FactionEnemy,
	}
}

// テスト用のアイテムコンポーネントを作成する
func CreateTestItem(t *testing.T, name string, itemType TestItemType) components.GameComponentList {
	t.Helper()

	base := components.GameComponentList{
		Item:             &gc.Item{},
		Name:             &gc.Name{Name: name},
		Description:      &gc.Description{Description: "テスト用アイテム"},
		ItemLocationType: &gc.ItemLocationInBackpack,
	}

	switch itemType {
	case TestItemTypeWeapon:
		base.Wearable = &gc.Wearable{
			Defense:           0,
			EquipmentCategory: gc.EquipmentTorso,
		}
		base.Attack = &gc.Attack{
			Accuracy:       80,
			Damage:         10,
			AttackCount:    1,
			Element:        gc.ElementTypeNone,
			AttackCategory: gc.AttackSword,
		}
	case TestItemTypeConsumable:
		base.Consumable = &gc.Consumable{
			UsableScene: gc.UsableSceneBattle,
			TargetType: gc.TargetType{
				TargetGroup: gc.TargetGroupAlly,
				TargetNum:   gc.TargetSingle,
			},
		}
		base.ProvidesHealing = &gc.ProvidesHealing{
			Amount: gc.NumeralAmount{Numeral: 30},
		}
	case TestItemTypeMaterial:
		base.Material = &gc.Material{Amount: 10}
	}

	return base
}

// テスト用アイテムタイプの定義
type TestItemType int

const (
	TestItemTypeWeapon TestItemType = iota
	TestItemTypeConsumable
	TestItemTypeMaterial
)

// テスト用のRawMasterを作成する
func CreateTestRawMaster(t *testing.T) raw.RawMaster {
	t.Helper()

	return raw.RawMaster{
		Raws: raw.Raws{
			Items: []raw.Item{
				{
					Name:        "テスト剣",
					Description: "テスト用の剣",
					Wearable: &raw.Wearable{
						Defense:           0,
						EquipmentCategory: "weapon",
					},
					Attack: &raw.Attack{
						Damage:         10,
						Accuracy:       80,
						AttackCount:    1,
						Element:        "none",
						AttackCategory: "sword",
					},
				},
				{
					Name:        "テスト薬草",
					Description: "テスト用の薬草",
					Consumable: &raw.Consumable{
						UsableScene: "battle",
						TargetGroup: "ally",
						TargetNum:   "single",
					},
					ProvidesHealing: &raw.ProvidesHealing{
						ValueType: raw.NumeralType,
						Amount:    30,
					},
				},
			},
			Materials: []raw.Material{
				{
					Name:        "テスト鉱石",
					Description: "テスト用の鉱石",
				},
			},
		},
		ItemIndex: map[string]int{
			"テスト剣":  0,
			"テスト薬草": 1,
		},
		MaterialIndex: map[string]int{
			"テスト鉱石": 0,
		},
	}
}

// テスト用のバトルシナリオを作成する
func CreateTestBattleScenario(t *testing.T) TestBattleScenario {
	t.Helper()

	return TestBattleScenario{
		Player: CreateTestPlayer(t),
		Enemies: []components.GameComponentList{
			CreateTestEnemy(t, "スライム"),
			CreateTestEnemy(t, "ゴブリン"),
		},
		Items: []components.GameComponentList{
			CreateTestItem(t, "テスト剣", TestItemTypeWeapon),
			CreateTestItem(t, "テスト薬草", TestItemTypeConsumable),
		},
	}
}

// テスト用バトルシナリオの構造体
type TestBattleScenario struct {
	Player  components.GameComponentList
	Enemies []components.GameComponentList
	Items   []components.GameComponentList
}
