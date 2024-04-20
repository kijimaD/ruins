package views

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/utils/consts"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func AddMemberBar(world w.World, targetContainer *widget.Container, entity ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)

	name := gameComponents.Name.Get(entity).(*gc.Name)
	targetContainer.AddChild(eui.NewMenuText(name.Name, world))

	pools := gameComponents.Pools.Get(entity).(*gc.Pools)
	targetContainer.AddChild(eui.NewMenuText(fmt.Sprintf("%s %3d/%3d", consts.HPLabel, pools.HP.Current, pools.HP.Max), world))

	hProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(140, 20),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the progress images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		// min, max, current
		widget.ProgressBarOpts.Values(0, pools.HP.Max, pools.HP.Current),
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
	)
	targetContainer.AddChild(hProgressbar)
}
