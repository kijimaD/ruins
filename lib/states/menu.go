package states

import (
	"image/color"

	"github.com/kijimaD/sokotwo/lib/engine/math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type menu interface {
	getSelection() int
	setSelection(selection int)
	confirmSelection(world w.World) states.Transition
	getMenuIDs() []string
	getCursorMenuIDs() []string
}

var menuLastCursorPosition = math.VectorInt2{}

func updateMenu(menu menu, world w.World) states.Transition {
	var transition states.Transition
	selection := menu.getSelection()
	numItems := len(menu.getCursorMenuIDs())

	// Handle keyboard events
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyDown), inpututil.IsKeyJustPressed(ebiten.KeyRight):
		menu.setSelection(math.Mod(selection+1, numItems))
	case inpututil.IsKeyJustPressed(ebiten.KeyUp), inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		menu.setSelection(math.Mod(selection-1, numItems))
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace):
		return menu.confirmSelection(world)
	}

	// Set cursor color
	newSelection := menu.getSelection()
	for iCursor, id := range menu.getCursorMenuIDs() {
		world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
			text := world.Components.Engine.Text.Get(entity).(*ec.Text)
			if text.ID == id {
				text.Color = color.RGBA{0, 0, 0, 0}
				if iCursor == newSelection {
					text.Color = color.RGBA{255, 255, 255, 255}
				}
			}
		}))
	}
	return transition
}
