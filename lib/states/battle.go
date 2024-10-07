package states

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type BattleState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	phase battlePhase

	// 敵表示コンテナ
	enemyListContainer *widget.Container
	// 各フェーズでの選択表示に使うリスト
	actionList *widget.List
}

func (st BattleState) String() string {
	return "Battle"
}

type battlePhase int

const (
	// 開戦 / 逃走
	phaseChoosePolicy battlePhase = iota
	// 各キャラの行動選択
	phaseChooseAction
	// 行動の対象選択
	phaseChooseTarget
	// 戦闘実行
	phaseExecute
	// リザルト画面
	phaseResult
)

// State interface ================

func (st *BattleState) OnPause(world w.World) {}

func (st *BattleState) OnResume(world w.World) {}

func (st *BattleState) OnStart(world w.World) {
	enemy := spawner.SpawnEnemy(world, "軽戦車")
	_ = gs.EquipmentChangedSystem(world) // これをしないとHP/SPが設定されない
	effects.AddEffect(nil, effects.Healing{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Single{Target: enemy})
	effects.RunEffectQueue(world)

	st.ui = st.initUI(world)
}

func (st *BattleState) OnStop(world w.World) {
	// 後片付け
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Enemy,
		gameComponents.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))
}

func (st *BattleState) Update(world w.World) states.Transition {
	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	st.ui.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
	}

	switch st.phase {
	case phaseChoosePolicy:
	case phaseChooseAction:
	case phaseChooseTarget:
	case phaseExecute:
	case phaseResult:
	}

	return states.Transition{Type: states.TransNone}
}

func (st *BattleState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)

	switch st.phase {
	case phaseChoosePolicy:
		st.reloadPolicy(world)
	case phaseChooseAction:
	case phaseChooseTarget:
	case phaseExecute:
	case phaseResult:
	}
}

// ================

func (st *BattleState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalTransContainer()
	st.enemyListContainer = eui.NewRowContainer()
	st.updateEnemyListContainer(world)
	st.reloadPolicy(world)
	rootContainer.AddChild(st.enemyListContainer)
	rootContainer.AddChild(st.actionList)
	return &ebitenui.UI{Container: rootContainer}
}

// メンバー一覧を更新する
func (st *BattleState) updateEnemyListContainer(world w.World) {
	st.enemyListContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Enemy,
		gameComponents.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		text := fmt.Sprintf("%s\n%3d/%3d", name.Name, pools.HP.Current, pools.HP.Max)
		st.enemyListContainer.AddChild(eui.NewMenuText(text, world))
	}))
}

type policyEntry string

const (
	policyEntryAttack policyEntry = "攻撃"
	policyEntryEscape policyEntry = "逃走"
)

func (st *BattleState) reloadPolicy(world w.World) {
	entries := []any{
		policyEntryAttack,
		policyEntryEscape,
	}
	st.actionList = widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
				Padding:            widget.NewInsetsSimple(50),
			}),
		)),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}, eui.LoadButtonImage()),
			widget.SliderOpts.MinHandleSize(5),
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntryFontFace(*eui.LoadFont(world)),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{R: 0, G: 255, B: 0, A: 255},
			Unselected:                 color.NRGBA{R: 254, G: 255, B: 255, A: 255},
			SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255},
			SelectingBackground:        color.NRGBA{R: 130, G: 130, B: 130, A: 255},
			SelectingFocusedBackground: color.NRGBA{R: 130, G: 140, B: 170, A: 255},
			SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255},
			FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255},
			DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			DisabledSelectedBackground: color.NRGBA{R: 100, G: 100, B: 100, A: 255},
		}),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			v, ok := e.(policyEntry)
			if !ok {
				log.Fatal("unexpected entry label selected!")
			}
			return string(v)
		}),
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(policyEntry)
			switch entry {
			case policyEntryAttack:
				st.phase = phaseChooseAction
			case policyEntryEscape:
				st.trans = &states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
			}
		}),
	)
}
