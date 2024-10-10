package states

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/euiext"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/utils/mathutil"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type BattleState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	// 現在のサブステート
	phase battlePhase
	// 1個前のサブステート。サブステート変更を検知するのに使う
	prevPhase battlePhase
	// 現在処理中のメンバーのインデックス。メンバーごとにコマンドを発行するため
	curMemberIndex int

	// 敵表示コンテナ
	enemyListContainer *widget.Container
	// 各フェーズでの選択表示に使うコンテナ
	selectContainer *widget.Container
}

func (st BattleState) String() string {
	return "Battle"
}

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

	// ステートが変わった最初の1回だけ実行される
	if st.phase != st.prevPhase {
		switch v := st.phase.(type) {
		case *phaseChoosePolicy:
			st.reloadPolicy(world)
		case *phaseChooseAction:
			st.reloadAction(world, v)
		case *phaseChooseTarget:
			st.reloadTarget(world, v)
		case *phaseExecute:
		case *phaseResult:
		}
		st.prevPhase = st.phase
	}

	// 毎回実行される
	switch st.phase.(type) {
	case *phaseChoosePolicy:
	case *phaseChooseAction:
	case *phaseChooseTarget:
	case *phaseExecute:
		effects.RunEffectQueue(world)
		gs.BattleCommandSystem(world)
		st.updateEnemyListContainer(world)
		st.reloadExecute(world)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			st.phase = &phaseChoosePolicy{}
		}
	case *phaseResult:
	}

	return states.Transition{Type: states.TransNone}
}

func (st *BattleState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *BattleState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalTransContainer()
	st.enemyListContainer = st.initEnemyContainer()
	st.updateEnemyListContainer(world)
	st.selectContainer = eui.NewVerticalContainer()
	st.reloadPolicy(world)
	rootContainer.AddChild(st.enemyListContainer)
	rootContainer.AddChild(st.selectContainer)

	return &ebitenui.UI{Container: rootContainer}
}

// 中央寄せのコンテナ
func (st *BattleState) initEnemyContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.DebugColor)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    10,
					Bottom: 10,
					Left:   10,
					Right:  10,
				}),
			)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxWidth:  200,
				MaxHeight: 200,
			}),
			widget.WidgetOpts.MinSize(0, 0),
		),
	)
}

// 敵一覧を更新する
func (st *BattleState) updateEnemyListContainer(world w.World) {
	st.enemyListContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Enemy,
		gameComponents.Pools,
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
	policyEntryItem   policyEntry = "道具"
	policyEntryEscape policyEntry = "逃走"
)

func (st *BattleState) reloadPolicy(world w.World) {
	st.selectContainer.RemoveChildren()

	entries := []any{
		policyEntryAttack,
		policyEntryItem,
		policyEntryEscape,
	}
	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e interface{}) string {
			v, ok := e.(policyEntry)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			return string(v)
		}),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			entry := args.Entry.(policyEntry)
			switch entry {
			case policyEntryAttack:
				members := []ecs.Entity{}
				simple.InPartyMember(world, func(entity ecs.Entity) {
					members = append(members, entity)
				})
				owner := members[st.curMemberIndex]
				st.phase = &phaseChooseAction{owner: owner}
			case policyEntryItem:
				// TODO: 未実装
			case policyEntryEscape:
				st.trans = &states.Transition{Type: states.TransPop}
			default:
				log.Fatal("unexpected entry detect!")
			}
		}),
	}
	list := eui.NewList(
		entries,
		[]widget.ButtonOpt{
			widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
				fmt.Println("Cursor Entered: " + args.Button.Text().Label)
			}),
		},
		opts,
		world,
	)
	st.selectContainer.AddChild(list)
}

// ================

func (st *BattleState) reloadAction(world w.World, currentPhase *phaseChooseAction) {
	st.selectContainer.RemoveChildren()
	st.updateEnemyListContainer(world)

	gameComponents := world.Components.Game.(*gc.Components)
	equipCards := []any{} // 実際にはecs.Entityが入る。Listで受け取るのが[]anyだからそうしている
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Equipped,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := gameComponents.Equipped.Get(entity).(*gc.Equipped)
		if currentPhase.owner == equipped.Owner {
			equipCards = append(equipCards, entity)
		}
	}))

	// 装備がなくても詰まないようにデフォルトの攻撃手段を追加する
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Card,
		gameComponents.Equipped.Not(),
		gameComponents.InBackpack.Not(),
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipCards = append(equipCards, entity)
	}))

	{
		members := []ecs.Entity{}
		simple.InPartyMember(world, func(entity ecs.Entity) {
			members = append(members, entity)
		})
		member := members[st.curMemberIndex]
		name := simple.GetName(world, member)
		st.selectContainer.AddChild(eui.NewMenuText(name.Name, world))
	}
	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e interface{}) string {
			v, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := simple.GetName(world, v)
			return name.Name
		}),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			cardEntity, ok := args.Entry.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			card := simple.GetCard(world, cardEntity)
			if card == nil {
				log.Fatal("unexpected error: entityがcardを保持していない")
			}
			st.phase = &phaseChooseTarget{
				owner: currentPhase.owner,
				way:   cardEntity,
			}
		}),
	}
	list := eui.NewList(
		equipCards,
		[]widget.ButtonOpt{
			widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
				fmt.Println("Cursor Entered: " + args.Button.Text().Label)
			}),
		},
		opts,
		world,
	)
	st.selectContainer.AddChild(list)
}

// ================

func (st *BattleState) reloadTarget(world w.World, currentPhase *phaseChooseTarget) {
	st.enemyListContainer.RemoveChildren()
	st.selectContainer.RemoveChildren()

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Enemy,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 敵キャラごとにターゲット選択ボタンを作成する
		vc := eui.NewVerticalContainer()
		st.enemyListContainer.AddChild(vc)

		name := gameComponents.Name.Get(entity).(*gc.Name)
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		text := fmt.Sprintf("%s\n%3d/%3d", name.Name, pools.HP.Current, pools.HP.Max)
		vc.AddChild(eui.NewMenuText(text, world))

		btn := eui.NewItemButton(
			"選択",
			func(args *widget.ButtonClickedEventArgs) {
				cl := loader.EntityComponentList{}
				cl.Game = append(cl.Game, components.GameComponentList{
					BattleCommand: &gc.BattleCommand{
						Owner:  currentPhase.owner,
						Target: entity,
						Way:    currentPhase.way,
					},
				})
				loader.AddEntities(world, cl)

				// みんな出揃ったらExecuteに、揃ってなかったらメンバーをインクリメントしてコマンド選択ステートへ
				members := []ecs.Entity{}
				simple.InPartyMember(world, func(entity ecs.Entity) {
					members = append(members, entity)
				})
				if st.curMemberIndex >= len(members)-1 {
					st.curMemberIndex = 0
					st.phase = &phaseExecute{}
				} else {
					st.curMemberIndex = mathutil.Min(st.curMemberIndex+1, len(members)-1)
					st.phase = &phaseChooseAction{owner: members[st.curMemberIndex]}
				}

			},
			world,
		)
		vc.AddChild(btn)
	}))
}

// ================

func (st *BattleState) reloadExecute(world w.World) {
	st.updateEnemyListContainer(world)
}
