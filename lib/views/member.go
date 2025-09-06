package views

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/widgets/common"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AddMemberBar は一人分のHPバーを表示する
func AddMemberBar(world w.World, targetContainer *widget.Container, entity ecs.Entity) {
	res := world.Resources.UIResources
	memberContainer := common.NewVerticalContainer()

	name := world.Components.Name.Get(entity).(*gc.Name)
	pools := world.Components.Pools.Get(entity).(*gc.Pools)
	memberContainer.AddChild(common.NewMenuText(name.Name, world))
	hpLabel := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s %3d/%3d", consts.HPLabel, pools.HP.Current, pools.HP.Max), res.Text.SmallFace, styles.TextColor),
	)
	memberContainer.AddChild(hpLabel)

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

	spLabel := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s %3d/%3d", consts.SPLabel, pools.SP.Current, pools.SP.Max), res.Text.SmallFace, styles.TextColor),
	)
	memberContainer.AddChild(spLabel)

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
	levelLabel := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Lv %d", pools.Level), res.Text.SmallFace, styles.TextColor),
	)
	memberContainer.AddChild(levelLabel)

	targetContainer.AddChild(memberContainer)
}
