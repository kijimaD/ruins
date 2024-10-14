package views

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/utils/consts"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 一人分のHPバーを表示する
func AddMemberBar(world w.World, targetContainer *widget.Container, entity ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	memberContainer := eui.NewVerticalContainer()

	name := gameComponents.Name.Get(entity).(*gc.Name)
	memberContainer.AddChild(eui.NewMenuText(name.Name, world))

	pools := gameComponents.Pools.Get(entity).(*gc.Pools)
	memberContainer.AddChild(eui.NewMenuText(fmt.Sprintf("%s %3d/%3d", consts.HPLabel, pools.HP.Current, pools.HP.Max), world))

	res := world.Resources.UIResources
	hProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(140, 20),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter},
			),
		),
		widget.ProgressBarOpts.Images(
			res.ProgressBar.TrackImage,
			res.ProgressBar.FillImage,
		),
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    3,
			Bottom: 3,
			Left:   2,
			Right:  2,
		}),
		widget.ProgressBarOpts.Values(0, pools.HP.Max, pools.HP.Current),
	)
	memberContainer.AddChild(hProgressbar)
	memberContainer.AddChild(eui.NewMenuText(fmt.Sprintf("LV %d", pools.Level), world))

	targetContainer.AddChild(memberContainer)
}
