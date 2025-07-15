package testing

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
)

// CreateTestPlayer はテスト用のプレイヤーコンポーネントを作成する
func CreateTestPlayer(t *testing.T) gc.GameComponentList {
	t.Helper()

	return gc.GameComponentList{
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

// CreateTestEnemy はテスト用の敵コンポーネントを作成する
func CreateTestEnemy(t *testing.T, name string) gc.GameComponentList {
	t.Helper()

	return gc.GameComponentList{
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

// CreateTestItem はテスト用のアイテムコンポーネントを作成する
func CreateTestItem(t *testing.T, name string, itemType TestItemType) gc.GameComponentList {
	t.Helper()

	base := gc.GameComponentList{
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

// TestItemType はテスト用アイテムタイプの定義
type TestItemType int

const (
	// TestItemTypeWeapon は武器タイプを表す
	TestItemTypeWeapon TestItemType = iota
	// TestItemTypeConsumable は消耗品タイプを表す
	TestItemTypeConsumable
	// TestItemTypeMaterial は素材タイプを表す
	TestItemTypeMaterial
)

// CreateTestRawMaster はテスト用のMasterを作成する
func CreateTestRawMaster(t *testing.T) raw.Master {
	t.Helper()

	return raw.Master{
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

// CreateTestBattleScenario はテスト用のバトルシナリオを作成する
func CreateTestBattleScenario(t *testing.T) TestBattleScenario {
	t.Helper()

	return TestBattleScenario{
		Player: CreateTestPlayer(t),
		Enemies: []gc.GameComponentList{
			CreateTestEnemy(t, "スライム"),
			CreateTestEnemy(t, "ゴブリン"),
		},
		Items: []gc.GameComponentList{
			CreateTestItem(t, "テスト剣", TestItemTypeWeapon),
			CreateTestItem(t, "テスト薬草", TestItemTypeConsumable),
		},
	}
}

// TestBattleScenario はテスト用バトルシナリオの構造体
type TestBattleScenario struct {
	Player  gc.GameComponentList
	Enemies []gc.GameComponentList
	Items   []gc.GameComponentList
}
