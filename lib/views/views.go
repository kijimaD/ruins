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
		if material, ok := component.(components.Material); ok {
			var amount string
			if material.Amount != 0 {
				amount = fmt.Sprintf("%d 個", material.Amount)
				targetContainer.AddChild(eui.NewBodyText(amount, styles.TextColor, world))
			}
		}
		if weapon, ok := component.(components.Weapon); ok {
			var accuracy string
			if weapon.Accuracy != 0 {
				accuracy = fmt.Sprintf("命中 %s", strconv.Itoa(weapon.Accuracy))
				targetContainer.AddChild(eui.NewBodyText(accuracy, styles.TextColor, world))
			}
			var baseDamage string
			if weapon.BaseDamage != 0 {
				baseDamage = fmt.Sprintf("攻撃力 %s", strconv.Itoa(weapon.BaseDamage))
				targetContainer.AddChild(eui.NewBodyText(baseDamage, styles.TextColor, world))
			}
			var consumption string
			if weapon.EnergyConsumption != 0 {
				consumption = fmt.Sprintf("消費SP %s", strconv.Itoa(weapon.EnergyConsumption))
				targetContainer.AddChild(eui.NewBodyText(consumption, styles.TextColor, world))
			}
		}
	}

	return targetContainer
}
