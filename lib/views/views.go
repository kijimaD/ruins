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

	for _, component := range cs {
		switch v := component.(type) {
		case *components.Material:
			if v == nil {
				continue
			}
			amount := fmt.Sprintf("%d 個", v.Amount)
			targetContainer.AddChild(eui.NewBodyText(amount, styles.TextColor, world))
		case *components.Weapon:
			if v == nil {
				continue
			}
			accuracy := fmt.Sprintf("命中 %s", strconv.Itoa(v.Accuracy))
			targetContainer.AddChild(eui.NewBodyText(accuracy, styles.TextColor, world))

			baseDamage := fmt.Sprintf("攻撃力 %s", strconv.Itoa(v.BaseDamage))
			targetContainer.AddChild(eui.NewBodyText(baseDamage, styles.TextColor, world))

			consumption := fmt.Sprintf("消費SP %s", strconv.Itoa(v.EnergyConsumption))
			targetContainer.AddChild(eui.NewBodyText(consumption, styles.TextColor, world))

			targetContainer.AddChild(damageAttrText(world, v.DamageAttr, v.DamageAttr.String()))
		case *components.Wearable:
			if v == nil {
				continue
			}
			baseDefense := fmt.Sprintf("防御力 %s", strconv.Itoa(v.BaseDefense))
			targetContainer.AddChild(eui.NewBodyText(baseDefense, styles.TextColor, world))

			equipmentSlot := fmt.Sprintf("部位 %s", v.EquipmentSlot)
			targetContainer.AddChild(eui.NewBodyText(equipmentSlot, styles.TextColor, world))
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
