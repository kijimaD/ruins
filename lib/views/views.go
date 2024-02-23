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
		case components.Material:
			var amount string
			if v.Amount != 0 {
				amount = fmt.Sprintf("%d 個", v.Amount)
				targetContainer.AddChild(eui.NewBodyText(amount, styles.TextColor, world))
			}
		case components.Weapon:
			var accuracy string
			if v.Accuracy != 0 {
				accuracy = fmt.Sprintf("命中 %s", strconv.Itoa(v.Accuracy))
				targetContainer.AddChild(eui.NewBodyText(accuracy, styles.TextColor, world))
			}
			var baseDamage string
			if v.BaseDamage != 0 {
				baseDamage = fmt.Sprintf("攻撃力 %s", strconv.Itoa(v.BaseDamage))
				targetContainer.AddChild(eui.NewBodyText(baseDamage, styles.TextColor, world))
			}
			var consumption string
			if v.EnergyConsumption != 0 {
				consumption = fmt.Sprintf("消費SP %s", strconv.Itoa(v.EnergyConsumption))
				targetContainer.AddChild(eui.NewBodyText(consumption, styles.TextColor, world))
			}
			if attr := v.DamageAttr.String(); attr != "" && attr != components.DamageAttrNone.String() {
				targetContainer.AddChild(eui.NewBodyText(fmt.Sprintf("<%s>", attr), styles.TextColor, world))
			}
		}
	}

	return targetContainer
}
