package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/sokotwo/lib/components"
	"github.com/kijimaD/sokotwo/lib/effects"
	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	selection     int
	inventoryMenu []ecs.Entity
	menuLen       int
}

// State interface ================

func (st *InventoryMenuState) OnPause(world w.World) {}

func (st *InventoryMenuState) OnResume(world w.World) {}

func (st *InventoryMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.inventoryMenu = append(st.inventoryMenu, loader.AddEntities(world, prefabs.Menu.InventoryMenu)...)
}

func (st *InventoryMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.inventoryMenu...)
}

func (st *InventoryMenuState) Update(world w.World) states.Transition {
	effects.RunEffectQueue(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&CampMenuState{}}}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	st.menuLen = 0
	itemList := ""
	descriptions := []string{}
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Name,
		gameComponents.Description,
		gameComponents.InBackpack,
		gameComponents.Consumable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		st.menuLen += 1

		name := gameComponents.Name.Get(entity).(*gc.Name)
		description := gameComponents.Description.Get(entity).(*gc.Description)
		descriptions = append(descriptions, description.Description)
		itemList += fmt.Sprintf("%s \n", name.Name)
	}))

	world.Manager.Join(
		world.Components.Engine.Text,
		world.Components.Engine.UITransform,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		switch text.ID {
		case "description":
			if len(descriptions) != 0 {
				text.Text = descriptions[st.selection]
			}
		case "item_list":
			text.Text = itemList
		case "cursor":
			ui := world.Components.Engine.UITransform.Get(entity).(*ec.UITransform)
			ui.Translation.Y = 500 - st.selection*32 - 10
		}
	}))

	return updateMenu(st, world)
}

// Menu Interface ================

func (st *InventoryMenuState) getSelection() int {
	return st.selection
}

func (st *InventoryMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *InventoryMenuState) confirmSelection(world w.World) states.Transition {
	gameComponents := world.Components.Game.(*gc.Components)
	var members []ecs.Entity
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		members = append(members, entity)
	}))

	switch st.selection {
	// アイテムを選択できるようにする
	case 0:
		// TODO: 仮で先頭の仲間固定にしている。ターゲットを選べるようにする
		effects.AddEffect(nil, effects.Damage{Amount: 10}, effects.Single{Target: members[0]})

		return states.Transition{Type: states.TransNone}
	}
	return states.Transition{Type: states.TransNone}
}

func (st *InventoryMenuState) getMenuIDs() []string {
	return []string{""}
}

func (st *InventoryMenuState) getCursorMenuIDs() []string {
	l := 0
	if st.menuLen == 0 {
		l = 1
	} else {
		l = st.menuLen
	}
	return make([]string, l)
}
