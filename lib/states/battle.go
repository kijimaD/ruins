package states

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type BattleState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	phase *battlePhase

	// 敵表示コンテナ
	enemyListContainer *widget.Container
	// 各フェーズでの選択表示に使うコンテナ
	selectContainer *widget.Container
}

func (st BattleState) String() string {
	return "Battle"
}

type battlePhase int

var (
	// 開戦 / 逃走
	phaseChoosePolicy battlePhase = 0
	// 各キャラの行動選択
	phaseChooseAction battlePhase = 1
	// 行動の対象選択
	phaseChooseTarget battlePhase = 2
	// 戦闘実行
	phaseExecute battlePhase = 3
	// リザルト画面
	phaseResult battlePhase = 4
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

	return states.Transition{Type: states.TransNone}
}

func (st *BattleState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)

	switch st.phase {
	case &phaseChoosePolicy:
		st.reloadPolicy(world)
	case &phaseChooseAction:
		st.reloadAction(world)
	case &phaseChooseTarget:
	case &phaseExecute:
	case &phaseResult:
	}
	if st.phase != nil {
		st.phase = nil
	}
}

// ================

func (st *BattleState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalTransContainer()
	st.enemyListContainer = eui.NewRowContainer()
	st.updateEnemyListContainer(world)
	st.selectContainer = eui.NewRowContainer()
	st.reloadPolicy(world)
	rootContainer.AddChild(st.enemyListContainer)
	rootContainer.AddChild(st.selectContainer)
	return &ebitenui.UI{Container: rootContainer}
}

// 敵一覧を更新する
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

// ================

type policyEntry string

const (
	policyEntryAttack policyEntry = "攻撃"
	policyEntryEscape policyEntry = "逃走"
)

func (st *BattleState) reloadPolicy(world w.World) {
	st.selectContainer.RemoveChildren()

	entries := []any{
		policyEntryAttack,
		policyEntryEscape,
	}
	list := eui.NewList(
		entries,
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			v, ok := e.(policyEntry)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			return string(v)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(policyEntry)
			switch entry {
			case policyEntryAttack:
				st.phase = &phaseChooseAction
			case policyEntryEscape:
				st.trans = &states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
			default:
				log.Fatal("unexpected entry detect!")
			}
		}),
		world,
	)
	st.selectContainer.AddChild(list)
}

// ================

func (st *BattleState) reloadAction(world w.World) {
	st.selectContainer.RemoveChildren()

	members := []ecs.Entity{}
	simple.InPartyMember(world, func(entity ecs.Entity) {
		members = append(members, entity)
	})
	// とりあえず先頭のメンバーだけ。本来は命令する対象による
	owner := members[0]
	gameComponents := world.Components.Game.(*gc.Components)
	eqs := []any{} // 実際にはecs.Entityが入る
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Equipped,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := gameComponents.Equipped.Get(entity).(*gc.Equipped)
		if owner == equipped.Owner {
			eqs = append(eqs, entity)
		}
	}))
	list := eui.NewList(
		eqs,
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			v, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := simple.GetName(world, v)
			return name.Name
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			v, ok := args.Entry.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := simple.GetName(world, v)
			// TODO: ここでpushする
			fmt.Println(name.Name)
		}),
		world,
	)
	st.selectContainer.AddChild(list)
}
