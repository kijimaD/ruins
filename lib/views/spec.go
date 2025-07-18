package views

import (
	"fmt"
	"strconv"

	"github.com/ebitenui/ebitenui/widget"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// UpdateSpec は性能表示コンテナを更新する
func UpdateSpec(world w.World, targetContainer *widget.Container, entity ecs.Entity) {
	targetContainer.RemoveChildren()

	{
		if entity.HasComponent(world.Components.Material) {
			v := world.Components.Material.Get(entity).(*gc.Material)
			amount := fmt.Sprintf("%d 個", v.Amount)
			targetContainer.AddChild(eui.NewBodyText(amount, styles.TextColor, world))
		}

		if entity.HasComponent(world.Components.Attack) {
			attack := world.Components.Attack.Get(entity).(*gc.Attack)
			targetContainer.AddChild(eui.NewBodyText(attack.AttackCategory.String(), styles.TextColor, world))

			damage := fmt.Sprintf("%s %s", consts.DamageLabel, strconv.Itoa(attack.Damage))
			targetContainer.AddChild(eui.NewBodyText(damage, styles.TextColor, world))

			accuracy := fmt.Sprintf("%s %s", consts.AccuracyLabel, strconv.Itoa(attack.Accuracy))
			targetContainer.AddChild(eui.NewBodyText(accuracy, styles.TextColor, world))

			attackCount := fmt.Sprintf("%s %s", consts.AttackCountLabel, strconv.Itoa(attack.AttackCount))
			targetContainer.AddChild(eui.NewBodyText(attackCount, styles.TextColor, world))

			if attack.Element != gc.ElementTypeNone {
				targetContainer.AddChild(damageAttrText(world, attack.Element, attack.Element.String()))
			}
		}
		if entity.HasComponent(world.Components.Wearable) {
			wearable := world.Components.Wearable.Get(entity).(*gc.Wearable)
			equipmentCategory := fmt.Sprintf("%s %s", consts.EquimentCategoryLabel, wearable.EquipmentCategory)
			targetContainer.AddChild(eui.NewBodyText(equipmentCategory, styles.TextColor, world))

			defense := fmt.Sprintf("%s %+d", consts.DefenseLabel, wearable.Defense)
			targetContainer.AddChild(eui.NewBodyText(defense, styles.TextColor, world))
			addEquipBonus(targetContainer, wearable.EquipBonus, world)
		}
		if entity.HasComponent(world.Components.Card) {
			card := world.Components.Card.Get(entity).(*gc.Card)
			cost := fmt.Sprintf("コスト %d", card.Cost)
			targetContainer.AddChild(eui.NewBodyText(cost, styles.TextColor, world))
		}
	}
}

// damageAttrText は属性によって色付けする
func damageAttrText(world w.World, dat gc.ElementType, str string) *widget.Text {
	var text *widget.Text
	switch dat {
	case gc.ElementTypeFire:
		text = eui.NewBodyText(str, styles.FireColor, world)
	case gc.ElementTypeThunder:
		text = eui.NewBodyText(str, styles.ThunderColor, world)
	case gc.ElementTypeChill:
		text = eui.NewBodyText(str, styles.ChillColor, world)
	case gc.ElementTypePhoton:
		text = eui.NewBodyText(str, styles.PhotonColor, world)
	default:
		text = eui.NewBodyText(str, styles.TextColor, world)
	}

	return text
}

// addEquipBonus は装備ボーナスを表示する
func addEquipBonus(targetContainer *widget.Container, equipBonus gc.EquipBonus, world w.World) {
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
