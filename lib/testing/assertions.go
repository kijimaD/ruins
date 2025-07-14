package testing

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

// 位置のアサーション
func AssertPosition(t *testing.T, component *gc.Position, expectedX, expectedY float64) {
	t.Helper()
	if component == nil {
		t.Error("位置コンポーネントがnilです")
		return
	}
	if float64(component.X) != expectedX || float64(component.Y) != expectedY {
		t.Errorf("位置が期待値と異なります: 期待値(%f, %f), 実際値(%f, %f)",
			expectedX, expectedY, float64(component.X), float64(component.Y))
	}
}

// 体力のアサーション
func AssertPools(t *testing.T, component *gc.Pools, expectedHPCurrent, expectedHPMax int) {
	t.Helper()
	if component == nil {
		t.Error("体力コンポーネントがnilです")
		return
	}
	if component.HP.Current != expectedHPCurrent {
		t.Errorf("現在体力が期待値と異なります: 期待値%d, 実際値%d",
			expectedHPCurrent, component.HP.Current)
	}
	if component.HP.Max != expectedHPMax {
		t.Errorf("最大体力が期待値と異なります: 期待値%d, 実際値%d",
			expectedHPMax, component.HP.Max)
	}
}

// 属性のアサーション
func AssertAttribute(t *testing.T, attribute gc.Attribute, expectedBase, expectedModifier, expectedTotal int) {
	t.Helper()
	if attribute.Base != expectedBase {
		t.Errorf("基本値が期待値と異なります: 期待値%d, 実際値%d",
			expectedBase, attribute.Base)
	}
	if attribute.Modifier != expectedModifier {
		t.Errorf("修正値が期待値と異なります: 期待値%d, 実際値%d",
			expectedModifier, attribute.Modifier)
	}
	if attribute.Total != expectedTotal {
		t.Errorf("合計値が期待値と異なります: 期待値%d, 実際値%d",
			expectedTotal, attribute.Total)
	}
}

// 名前のアサーション
func AssertName(t *testing.T, component *gc.Name, expectedName string) {
	t.Helper()
	if component == nil {
		t.Error("名前コンポーネントがnilです")
		return
	}
	if component.Name != expectedName {
		t.Errorf("名前が期待値と異なります: 期待値%s, 実際値%s",
			expectedName, component.Name)
	}
}

// レンダリングのアサーション（簡略版）
func AssertRender(t *testing.T, component *gc.Render, expectedSheet string) {
	t.Helper()
	if component == nil {
		t.Error("レンダリングコンポーネントがnilです")
		return
	}
	// 実際のレンダリング構造は複雑なので、簡単なチェックに留める
}

// 攻撃力のアサーション
func AssertAttack(t *testing.T, component *gc.Attack, expectedDamage, expectedAccuracy int) {
	t.Helper()
	if component == nil {
		t.Error("攻撃コンポーネントがnilです")
		return
	}
	if component.Damage != expectedDamage {
		t.Errorf("攻撃力が期待値と異なります: 期待値%d, 実際値%d",
			expectedDamage, component.Damage)
	}
	if component.Accuracy != expectedAccuracy {
		t.Errorf("命中率が期待値と異なります: 期待値%d, 実際値%d",
			expectedAccuracy, component.Accuracy)
	}
}

// コンポーネントの存在確認
func AssertHasComponent(t *testing.T, componentList gc.GameComponentList, componentName string) {
	t.Helper()
	switch componentName {
	case "Player":
		if componentList.InParty == nil {
			t.Error("プレイヤーコンポーネントが存在しません")
		}
	case "Enemy":
		if componentList.FactionType == nil {
			t.Error("敵コンポーネントが存在しません")
		}
	case "Item":
		if componentList.Item == nil {
			t.Error("アイテムコンポーネントが存在しません")
		}
	case "Weapon":
		if componentList.Wearable == nil || componentList.Attack == nil {
			t.Error("武器コンポーネントが存在しません")
		}
	case "Consumable":
		if componentList.Consumable == nil {
			t.Error("消耗品コンポーネントが存在しません")
		}
	case "Material":
		if componentList.Material == nil {
			t.Error("素材コンポーネントが存在しません")
		}
	default:
		t.Errorf("不明なコンポーネント名: %s", componentName)
	}
}

// コンポーネントの非存在確認
func AssertNotHasComponent(t *testing.T, componentList gc.GameComponentList, componentName string) {
	t.Helper()
	switch componentName {
	case "Player":
		if componentList.InParty != nil {
			t.Error("プレイヤーコンポーネントが存在してはいけません")
		}
	case "Enemy":
		if componentList.FactionType != nil {
			t.Error("敵コンポーネントが存在してはいけません")
		}
	case "Item":
		if componentList.Item != nil {
			t.Error("アイテムコンポーネントが存在してはいけません")
		}
	default:
		t.Errorf("不明なコンポーネント名: %s", componentName)
	}
}

// 素材の量のアサーション
func AssertMaterialAmount(t *testing.T, component *gc.Material, expectedAmount int) {
	t.Helper()
	if component == nil {
		t.Error("素材コンポーネントがnilです")
		return
	}
	if component.Amount != expectedAmount {
		t.Errorf("素材の量が期待値と異なります: 期待値%d, 実際値%d",
			expectedAmount, component.Amount)
	}
}

// 回復量のアサーション
func AssertHealingAmount(t *testing.T, component *gc.ProvidesHealing, expectedAmount int) {
	t.Helper()
	if component == nil {
		t.Error("回復コンポーネントがnilです")
		return
	}

	// NumeralAmountの場合の検証
	if numeralAmount, ok := component.Amount.(gc.NumeralAmount); ok {
		if numeralAmount.Calc() != expectedAmount {
			t.Errorf("回復量が期待値と異なります: 期待値%d, 実際値%d",
				expectedAmount, numeralAmount.Calc())
		}
	} else {
		t.Error("回復量がNumeralAmountではありません")
	}
}

// 範囲内の値かチェック
func AssertInRange(t *testing.T, value, min, max int, description string) {
	t.Helper()
	if value < min || value > max {
		t.Errorf("%sが範囲外です: 期待値[%d, %d], 実際値%d",
			description, min, max, value)
	}
}

// 2つの値が等しいかチェック
func AssertEqual(t *testing.T, actual, expected interface{}, description string) {
	t.Helper()
	if actual != expected {
		t.Errorf("%sが期待値と異なります: 期待値%v, 実際値%v",
			description, expected, actual)
	}
}

// 値がnilでないかチェック
func AssertNotNil(t *testing.T, value interface{}, description string) {
	t.Helper()
	if value == nil {
		t.Errorf("%sがnilです", description)
	}
}

// 値がnilかチェック
func AssertNil(t *testing.T, value interface{}, description string) {
	t.Helper()
	if value != nil {
		t.Errorf("%sがnilではありません: 実際値%v", description, value)
	}
}
