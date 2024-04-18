package views

import (
	"fmt"
	"strconv"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/utils/consts"
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
		case *components.Attack:
			if v == nil {
				continue
			}
			targetContainer.AddChild(eui.NewBodyText(v.AttackCategory.String(), styles.TextColor, world))

			accuracy := fmt.Sprintf("%s %s", consts.AccuracyLabel, strconv.Itoa(v.Accuracy))
			targetContainer.AddChild(eui.NewBodyText(accuracy, styles.TextColor, world))

			damage := fmt.Sprintf("%s %s", consts.DamageLabel, strconv.Itoa(v.Damage))
			targetContainer.AddChild(eui.NewBodyText(damage, styles.TextColor, world))

			attackCount := fmt.Sprintf("%s %s", consts.AttackCountLabel, strconv.Itoa(v.AttackCount))
			targetContainer.AddChild(eui.NewBodyText(attackCount, styles.TextColor, world))

			if v.Element != components.ElementTypeNone {
				targetContainer.AddChild(damageAttrText(world, v.Element, v.Element.String()))
			}
		case *components.Wearable:
			if v == nil {
				continue
			}
			equipmentCategory := fmt.Sprintf("%s %s", consts.EquimentCategoryLabel, v.EquipmentCategory)
			targetContainer.AddChild(eui.NewBodyText(equipmentCategory, styles.TextColor, world))

			defense := fmt.Sprintf("%s %s", consts.DefenseLabel, strconv.Itoa(v.Defense))
			targetContainer.AddChild(eui.NewBodyText(defense, styles.TextColor, world))
			addEquipBonus(targetContainer, v.EquipBonus, world)
		case *components.Card:
			if v == nil {
				continue
			}
			cost := fmt.Sprintf("コスト %d", v.Cost)
			targetContainer.AddChild(eui.NewBodyText(cost, styles.TextColor, world))
		}
	}

	return targetContainer
}

// 属性によって色付けする
func damageAttrText(world w.World, dat components.ElementType, str string) *widget.Text {
	var text *widget.Text
	switch dat {
	case components.ElementTypeFire:
		text = eui.NewBodyText(str, styles.FireColor, world)
	case components.ElementTypeThunder:
		text = eui.NewBodyText(str, styles.ThunderColor, world)
	case components.ElementTypeChill:
		text = eui.NewBodyText(str, styles.ChillColor, world)
	case components.ElementTypePhoton:
		text = eui.NewBodyText(str, styles.PhotonColor, world)
	default:
		text = eui.NewBodyText(str, styles.TextColor, world)
	}

	return text
}

func addEquipBonus(targetContainer *widget.Container, equipBonus components.EquipBonus, world w.World) {
	if equipBonus.Vitality != 0 {
		vitality := fmt.Sprintf("%s %+d", consts.VitalityLabel, equipBonus.Vitality)
		targetContainer.AddChild(eui.NewBodyText(vitality, styles.TextColor, world))
	}

	if equipBonus.Strength != 0 {
		strength := fmt.Sprintf("%s %+d", consts.StrengthLabel, equipBonus.Strength)
		targetContainer.AddChild(eui.NewBodyText(strength, styles.TextColor, world))
	}

	if equipBonus.Sensation != 0 {
		sensation := fmt.Sprintf("%s %+d", consts.SensationLabel, equipBonus.Sensation)
		targetContainer.AddChild(eui.NewBodyText(sensation, styles.TextColor, world))
	}

	if equipBonus.Dexterity != 0 {
		dexterity := fmt.Sprintf("%s %+d", consts.DexterityLabel, equipBonus.Dexterity)
		targetContainer.AddChild(eui.NewBodyText(dexterity, styles.TextColor, world))
	}

	if equipBonus.Agility != 0 {
		agility := fmt.Sprintf("%s %+d", consts.AgilityLabel, equipBonus.Agility)
		targetContainer.AddChild(eui.NewBodyText(agility, styles.TextColor, world))
	}
}
