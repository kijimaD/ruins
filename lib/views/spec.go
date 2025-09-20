package views

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/ruins/lib/colors"
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
		if entity.HasComponent(world.Components.Material) {
			v := world.Components.Material.Get(entity).(*gc.Material)
			amount := fmt.Sprintf("%d 個", v.Amount)
			targetContainer.AddChild(styled.NewBodyText(amount, colors.TextColor, world))
		}

		if entity.HasComponent(world.Components.Attack) {
			attack := world.Components.Attack.Get(entity).(*gc.Attack)
			targetContainer.AddChild(styled.NewBodyText(attack.AttackCategory.String(), colors.TextColor, world))

			damage := fmt.Sprintf("%s %s", consts.DamageLabel, strconv.Itoa(attack.Damage))
			targetContainer.AddChild(styled.NewBodyText(damage, colors.TextColor, world))

			accuracy := fmt.Sprintf("%s %s", consts.AccuracyLabel, strconv.Itoa(attack.Accuracy))
			targetContainer.AddChild(styled.NewBodyText(accuracy, colors.TextColor, world))

			attackCount := fmt.Sprintf("%s %s", consts.AttackCountLabel, strconv.Itoa(attack.AttackCount))
			targetContainer.AddChild(styled.NewBodyText(attackCount, colors.TextColor, world))

			if attack.Element != gc.ElementTypeNone {
				targetContainer.AddChild(damageAttrText(world, attack.Element, attack.Element.String()))
			}
		}
		if entity.HasComponent(world.Components.Wearable) {
			wearable := world.Components.Wearable.Get(entity).(*gc.Wearable)
			equipmentCategory := fmt.Sprintf("%s %s", consts.EquimentCategoryLabel, wearable.EquipmentCategory)
			targetContainer.AddChild(styled.NewBodyText(equipmentCategory, colors.TextColor, world))

			defense := fmt.Sprintf("%s %+d", consts.DefenseLabel, wearable.Defense)
			targetContainer.AddChild(styled.NewBodyText(defense, colors.TextColor, world))
			addEquipBonus(targetContainer, wearable.EquipBonus, world)
		}
		if entity.HasComponent(world.Components.Card) {
			card := world.Components.Card.Get(entity).(*gc.Card)
			cost := fmt.Sprintf("コスト %d", card.Cost)
			targetContainer.AddChild(styled.NewBodyText(cost, colors.TextColor, world))
		}
	}
}

// damageAttrText は属性によって色付けする
func damageAttrText(world w.World, dat gc.ElementType, str string) *widget.Text {
	var text *widget.Text
	switch dat {
	case gc.ElementTypeFire:
		text = styled.NewBodyText(str, colors.FireColor, world)
	case gc.ElementTypeThunder:
		text = styled.NewBodyText(str, colors.ThunderColor, world)
	case gc.ElementTypeChill:
		text = styled.NewBodyText(str, colors.ChillColor, world)
	case gc.ElementTypePhoton:
		text = styled.NewBodyText(str, colors.PhotonColor, world)
	default:
		text = styled.NewBodyText(str, colors.TextColor, world)
	}

	return text
}

// addEquipBonus は装備ボーナスを表示する
func addEquipBonus(targetContainer *widget.Container, equipBonus gc.EquipBonus, world w.World) {
	if equipBonus.Vitality != 0 {
		vitality := fmt.Sprintf("%s %+d", consts.VitalityLabel, equipBonus.Vitality)
		targetContainer.AddChild(styled.NewBodyText(vitality, colors.TextColor, world))
	}

	if equipBonus.Strength != 0 {
		strength := fmt.Sprintf("%s %+d", consts.StrengthLabel, equipBonus.Strength)
		targetContainer.AddChild(styled.NewBodyText(strength, colors.TextColor, world))
	}

	if equipBonus.Sensation != 0 {
		sensation := fmt.Sprintf("%s %+d", consts.SensationLabel, equipBonus.Sensation)
		targetContainer.AddChild(styled.NewBodyText(sensation, colors.TextColor, world))
	}

	if equipBonus.Dexterity != 0 {
		dexterity := fmt.Sprintf("%s %+d", consts.DexterityLabel, equipBonus.Dexterity)
		targetContainer.AddChild(styled.NewBodyText(dexterity, colors.TextColor, world))
	}

	if equipBonus.Agility != 0 {
		agility := fmt.Sprintf("%s %+d", consts.AgilityLabel, equipBonus.Agility)
		targetContainer.AddChild(styled.NewBodyText(agility, colors.TextColor, world))
	}
}

// AddMemberStatusText はメンバーの名前とHPを簡易テキスト表示で追加する
func AddMemberStatusText(targetContainer *widget.Container, entity ecs.Entity, world w.World) {
	if !entity.HasComponent(world.Components.Name) || !entity.HasComponent(world.Components.Pools) {
		return
	}

	name := world.Components.Name.Get(entity).(*gc.Name)
	pools := world.Components.Pools.Get(entity).(*gc.Pools)

	targetContainer.AddChild(styled.NewMenuText(name.Name, world))
	targetContainer.AddChild(styled.NewBodyText(fmt.Sprintf("%s %3d/%3d", consts.HPLabel, pools.HP.Current, pools.HP.Max), colors.TextColor, world))
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
	memberContainer.AddChild(styled.NewMenuText(name.Name, world))

	// HPラベル
	hpLabel := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s %3d/%3d", consts.HPLabel, pools.HP.Current, pools.HP.Max), res.Text.SmallFace, colors.TextColor),
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
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
		widget.ProgressBarOpts.Values(0, pools.HP.Max, pools.HP.Current),
	)
	memberContainer.AddChild(hpProgressbar)

	// SPラベル
	spLabel := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s %3d/%3d", consts.SPLabel, pools.SP.Current, pools.SP.Max), res.Text.SmallFace, colors.TextColor),
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
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
		widget.ProgressBarOpts.Values(0, pools.SP.Max, pools.SP.Current),
	)
	memberContainer.AddChild(spProgressbar)

	targetContainer.AddChild(memberContainer)
}
