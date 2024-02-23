package views

import (
	"fmt"
	"strconv"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
)

func UpdateSpec(world w.World, targetContainer *widget.Container, cs []any) *widget.Container {
	targetContainer.RemoveChildren()

	// FIXME: 必要なすべてのstate全体で必要なコンポーネントが渡されてくるので、たとえば消耗品のタブ時でも武器のcase内は実行される
	for _, component := range cs {
		switch v := component.(type) {
		case components.Material:
			if v.Amount != 0 {
				amount := fmt.Sprintf("%d 個", v.Amount)
				targetContainer.AddChild(eui.NewBodyText(amount, styles.TextColor, world))
			}
		case components.Weapon:
			if v.Accuracy != 0 {
				accuracy := fmt.Sprintf("命中 %s", strconv.Itoa(v.Accuracy))
				targetContainer.AddChild(eui.NewBodyText(accuracy, styles.TextColor, world))
			}
			if v.BaseDamage != 0 {
				baseDamage := fmt.Sprintf("攻撃力 %s", strconv.Itoa(v.BaseDamage))
				targetContainer.AddChild(eui.NewBodyText(baseDamage, styles.TextColor, world))
			}
			if v.EnergyConsumption != 0 {
				consumption := fmt.Sprintf("消費SP %s", strconv.Itoa(v.EnergyConsumption))
				targetContainer.AddChild(eui.NewBodyText(consumption, styles.TextColor, world))
			}
			if attr := v.DamageAttr.String(); attr != "" && attr != components.DamageAttrNone.String() {
				text := damageAttrText(world, v.DamageAttr, attr)
				targetContainer.AddChild(text)
			}
		case components.Wearable:
			if v.BaseDefense != 0 {
				baseDefense := fmt.Sprintf("防御力 %s", strconv.Itoa(v.BaseDefense))
				targetContainer.AddChild(eui.NewBodyText(baseDefense, styles.TextColor, world))
			}
			if attr := v.EquipmentSlot.String(); attr != "" {
				equipmentSlot := fmt.Sprintf("部位 %s", v.EquipmentSlot)
				targetContainer.AddChild(eui.NewBodyText(equipmentSlot, styles.TextColor, world))
			}
		}
	}

	return targetContainer
}

// 属性によって色付けする
func damageAttrText(world w.World, dat components.DamageAttrType, str string) *widget.Text {
	var text *widget.Text
	switch dat {
	case components.DamageAttrFire:
		text = eui.NewBodyText(str, styles.FireColor, world)
	case components.DamageAttrThunder:
		text = eui.NewBodyText(str, styles.ThunderColor, world)
	case components.DamageAttrChill:
		text = eui.NewBodyText(str, styles.ChillColor, world)
	case components.DamageAttrPhoton:
		text = eui.NewBodyText(str, styles.PhotonColor, world)
	default:
		text = eui.NewBodyText(str, styles.TextColor, world)
	}

	return text
}
