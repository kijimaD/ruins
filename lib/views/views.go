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

func UpdateSpec(world w.World, targetContainer *widget.Container, weapon components.Weapon) *widget.Container {
	targetContainer.RemoveChildren()

	var accuracy string
	if weapon.Accuracy != 0 {
		accuracy = fmt.Sprintf("命中 %s", strconv.Itoa(weapon.Accuracy))
	}
	var baseDamage string
	if weapon.BaseDamage != 0 {
		baseDamage = fmt.Sprintf("攻撃力 %s", strconv.Itoa(weapon.BaseDamage))
	}
	var consumption string
	if weapon.EnergyConsumption != 0 {
		consumption = fmt.Sprintf("消費SP %s", strconv.Itoa(weapon.EnergyConsumption))
	}
	targetContainer.AddChild(eui.NewBodyText(accuracy, styles.TextColor, world))
	targetContainer.AddChild(eui.NewBodyText(baseDamage, styles.TextColor, world))
	targetContainer.AddChild(eui.NewBodyText(consumption, styles.TextColor, world))

	return targetContainer
}
