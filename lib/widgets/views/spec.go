package views

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// UpdateSpec は性能表示コンテナを更新する
func UpdateSpec(world w.World, targetContainer *widget.Container, entity ecs.Entity) {
	targetContainer.RemoveChildren()

	{
		if entity.HasComponent(world.Components.Stackable) {
			v := world.Components.Stackable.Get(entity).(*gc.Stackable)
			amount := fmt.Sprintf("%d 個", v.Count)
			targetContainer.AddChild(styled.NewBodyText(amount, consts.TextColor, world.Resources.UIResources))
		}

		if entity.HasComponent(world.Components.Value) {
			v := world.Components.Value.Get(entity).(*gc.Value)
			value := fmt.Sprintf("◆ %d", v.Value)
			targetContainer.AddChild(styled.NewBodyText(value, consts.TextColor, world.Resources.UIResources))
		}

		if entity.HasComponent(world.Components.Attack) {
			attack := world.Components.Attack.Get(entity).(*gc.Attack)
			targetContainer.AddChild(styled.NewBodyText(attack.AttackCategory.String(), consts.TextColor, world.Resources.UIResources))

			damage := fmt.Sprintf("%s %s", consts.DamageLabel, strconv.Itoa(attack.Damage))
			targetContainer.AddChild(styled.NewBodyText(damage, consts.TextColor, world.Resources.UIResources))

			accuracy := fmt.Sprintf("%s %s", consts.AccuracyLabel, strconv.Itoa(attack.Accuracy))
			targetContainer.AddChild(styled.NewBodyText(accuracy, consts.TextColor, world.Resources.UIResources))

			attackCount := fmt.Sprintf("%s %s", consts.AttackCountLabel, strconv.Itoa(attack.AttackCount))
			targetContainer.AddChild(styled.NewBodyText(attackCount, consts.TextColor, world.Resources.UIResources))

			if attack.Element != gc.ElementTypeNone {
				targetContainer.AddChild(damageAttrText(world, attack.Element, attack.Element.String()))
			}
		}
		if entity.HasComponent(world.Components.Wearable) {
			wearable := world.Components.Wearable.Get(entity).(*gc.Wearable)
			equipmentCategory := fmt.Sprintf("%s %s", consts.EquimentCategoryLabel, wearable.EquipmentCategory)
			targetContainer.AddChild(styled.NewBodyText(equipmentCategory, consts.TextColor, world.Resources.UIResources))

			defense := fmt.Sprintf("%s %+d", consts.DefenseLabel, wearable.Defense)
			targetContainer.AddChild(styled.NewBodyText(defense, consts.TextColor, world.Resources.UIResources))
			addEquipBonus(targetContainer, wearable.EquipBonus, world)
		}
		if entity.HasComponent(world.Components.Card) {
			card := world.Components.Card.Get(entity).(*gc.Card)
			cost := fmt.Sprintf("コスト %d", card.Cost)
			targetContainer.AddChild(styled.NewBodyText(cost, consts.TextColor, world.Resources.UIResources))
		}
	}
}

// UpdateSpecFromSpec はEntitySpecから性能表示コンテナを更新する
// エンティティを生成せずに性能を表示できる
func UpdateSpecFromSpec(world w.World, targetContainer *widget.Container, spec gc.EntitySpec) {
	targetContainer.RemoveChildren()

	if spec.Value != nil {
		value := fmt.Sprintf("◆ %d", spec.Value.Value)
		targetContainer.AddChild(styled.NewBodyText(value, consts.TextColor, world.Resources.UIResources))
	}

	if spec.Attack != nil {
		targetContainer.AddChild(styled.NewBodyText(spec.Attack.AttackCategory.String(), consts.TextColor, world.Resources.UIResources))

		damage := fmt.Sprintf("%s %s", consts.DamageLabel, strconv.Itoa(spec.Attack.Damage))
		targetContainer.AddChild(styled.NewBodyText(damage, consts.TextColor, world.Resources.UIResources))

		accuracy := fmt.Sprintf("%s %s", consts.AccuracyLabel, strconv.Itoa(spec.Attack.Accuracy))
		targetContainer.AddChild(styled.NewBodyText(accuracy, consts.TextColor, world.Resources.UIResources))

		attackCount := fmt.Sprintf("%s %s", consts.AttackCountLabel, strconv.Itoa(spec.Attack.AttackCount))
		targetContainer.AddChild(styled.NewBodyText(attackCount, consts.TextColor, world.Resources.UIResources))

		if spec.Attack.Element != gc.ElementTypeNone {
			targetContainer.AddChild(damageAttrText(world, spec.Attack.Element, spec.Attack.Element.String()))
		}
	}

	if spec.Wearable != nil {
		equipmentCategory := fmt.Sprintf("%s %s", consts.EquimentCategoryLabel, spec.Wearable.EquipmentCategory)
		targetContainer.AddChild(styled.NewBodyText(equipmentCategory, consts.TextColor, world.Resources.UIResources))

		defense := fmt.Sprintf("%s %+d", consts.DefenseLabel, spec.Wearable.Defense)
		targetContainer.AddChild(styled.NewBodyText(defense, consts.TextColor, world.Resources.UIResources))
		addEquipBonus(targetContainer, spec.Wearable.EquipBonus, world)
	}

	if spec.Card != nil {
		cost := fmt.Sprintf("コスト %d", spec.Card.Cost)
		targetContainer.AddChild(styled.NewBodyText(cost, consts.TextColor, world.Resources.UIResources))
	}
}

// damageAttrText は属性によって色付けする
func damageAttrText(world w.World, dat gc.ElementType, str string) *widget.Text {
	res := world.Resources.UIResources
	var text *widget.Text
	switch dat {
	case gc.ElementTypeFire:
		text = styled.NewBodyText(str, consts.FireColor, res)
	case gc.ElementTypeThunder:
		text = styled.NewBodyText(str, consts.ThunderColor, res)
	case gc.ElementTypeChill:
		text = styled.NewBodyText(str, consts.ChillColor, res)
	case gc.ElementTypePhoton:
		text = styled.NewBodyText(str, consts.PhotonColor, res)
	default:
		text = styled.NewBodyText(str, consts.TextColor, res)
	}

	return text
}

// addEquipBonus は装備ボーナスを表示する
func addEquipBonus(targetContainer *widget.Container, equipBonus gc.EquipBonus, world w.World) {
	if equipBonus.Vitality != 0 {
		vitality := fmt.Sprintf("%s %+d", consts.VitalityLabel, equipBonus.Vitality)
		targetContainer.AddChild(styled.NewBodyText(vitality, consts.TextColor, world.Resources.UIResources))
	}

	if equipBonus.Strength != 0 {
		strength := fmt.Sprintf("%s %+d", consts.StrengthLabel, equipBonus.Strength)
		targetContainer.AddChild(styled.NewBodyText(strength, consts.TextColor, world.Resources.UIResources))
	}

	if equipBonus.Sensation != 0 {
		sensation := fmt.Sprintf("%s %+d", consts.SensationLabel, equipBonus.Sensation)
		targetContainer.AddChild(styled.NewBodyText(sensation, consts.TextColor, world.Resources.UIResources))
	}

	if equipBonus.Dexterity != 0 {
		dexterity := fmt.Sprintf("%s %+d", consts.DexterityLabel, equipBonus.Dexterity)
		targetContainer.AddChild(styled.NewBodyText(dexterity, consts.TextColor, world.Resources.UIResources))
	}

	if equipBonus.Agility != 0 {
		agility := fmt.Sprintf("%s %+d", consts.AgilityLabel, equipBonus.Agility)
		targetContainer.AddChild(styled.NewBodyText(agility, consts.TextColor, world.Resources.UIResources))
	}
}

// AddMemberStatusText はメンバーの名前とHPを簡易テキスト表示で追加する
func AddMemberStatusText(targetContainer *widget.Container, entity ecs.Entity, world w.World) {
	if !entity.HasComponent(world.Components.Name) || !entity.HasComponent(world.Components.Pools) {
		return
	}

	name := world.Components.Name.Get(entity).(*gc.Name)
	pools := world.Components.Pools.Get(entity).(*gc.Pools)

	targetContainer.AddChild(styled.NewMenuText(name.Name, world.Resources.UIResources))
	targetContainer.AddChild(styled.NewBodyText(fmt.Sprintf("%s %3d/%3d", consts.HPLabel, pools.HP.Current, pools.HP.Max), consts.TextColor, world.Resources.UIResources))
}

// AddMemberBars はメンバーの名前、HP/SPバー、レベルを詳細表示で追加する
func AddMemberBars(targetContainer *widget.Container, entity ecs.Entity, world w.World) {
	if !entity.HasComponent(world.Components.Name) || !entity.HasComponent(world.Components.Pools) {
		return
	}

	name := world.Components.Name.Get(entity).(*gc.Name)
	pools := world.Components.Pools.Get(entity).(*gc.Pools)
	res := world.Resources.UIResources

	memberContainer := styled.NewVerticalContainer()

	// 名前
	memberContainer.AddChild(styled.NewMenuText(name.Name, world.Resources.UIResources))

	// HPラベル
	hpLabel := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s %3d/%3d", consts.HPLabel, pools.HP.Current, pools.HP.Max), &res.Text.SmallFace, consts.TextColor),
	)
	memberContainer.AddChild(hpLabel)

	// HPプログレスバー
	hpProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(140, 16),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter},
			),
		),
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{0, 200, 0, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		widget.ProgressBarOpts.TrackPadding(&widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
		widget.ProgressBarOpts.Values(0, pools.HP.Max, pools.HP.Current),
	)
	memberContainer.AddChild(hpProgressbar)

	// SPラベル
	spLabel := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s %3d/%3d", consts.SPLabel, pools.SP.Current, pools.SP.Max), &res.Text.SmallFace, consts.TextColor),
	)
	memberContainer.AddChild(spLabel)

	// SPプログレスバー
	spProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(140, 16),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter},
			),
		),
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{255, 200, 0, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		widget.ProgressBarOpts.TrackPadding(&widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
		widget.ProgressBarOpts.Values(0, pools.SP.Max, pools.SP.Current),
	)
	memberContainer.AddChild(spProgressbar)

	targetContainer.AddChild(memberContainer)
}
