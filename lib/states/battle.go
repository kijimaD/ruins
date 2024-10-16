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
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/utils/mathutil"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	// 戦闘ログメッセージの高さ。文字分で指定する
	MessageCharBaseHeight = 10
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
	// テキスト送り待ち状態
	isWaitClick bool

	// 背景
	bg *ebiten.Image
	// 敵表示コンテナ
	enemyListContainer *widget.Container
	// 各フェーズでの選択表示に使うコンテナ
	selectContainer *widget.Container
	// 味方表示コンテナ
	memberContainer *widget.Container

	// 選択中のアイテム
	selectedItem ecs.Entity
	// カードの説明表示コンテナ
	cardSpecContainer *widget.Container
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

	bg := (*world.Resources.SpriteSheets)["bg_jungle1"]
	st.bg = bg.Texture.Image

	st.ui = st.initUI(world)
}

func (st *BattleState) OnStop(world w.World) {
	// 後片付け
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))

	gamelog.BattleLog.Flush()
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
		st.updateEnemyListContainer(world)
		st.reloadExecute(world)
		st.reloadMsg(world)

		// commandが残っていればクリック待ちにする
		commandCount := 0
		gameComponents := world.Components.Game.(*gc.Components)
		world.Manager.Join(
			gameComponents.BattleCommand,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			commandCount += 1
		}))

		if commandCount != 0 {
			// 未処理のコマンドがまだ残っている
			st.isWaitClick = true
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				gs.BattleCommandSystem(world)
				st.isWaitClick = false
			}
			return states.Transition{Type: states.TransNone}
		}

		// 選択完了
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			st.phase = &phaseChoosePolicy{}
			st.isWaitClick = false
			gamelog.BattleLog.Flush()
		}
	case *phaseResult:
	}

	return states.Transition{Type: states.TransNone}
}

func (st *BattleState) Draw(world w.World, screen *ebiten.Image) {
	if st.bg != nil {
		screen.DrawImage(st.bg, &ebiten.DrawImageOptions{})
	}

	st.ui.Draw(screen)
}

// ================

func (st *BattleState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()
	st.enemyListContainer = st.initEnemyContainer()
	st.updateEnemyListContainer(world)

	st.selectContainer = eui.NewVerticalContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 120)),
	)
	st.reloadPolicy(world)

	// 非表示にできるように背景が設定されていない
	st.cardSpecContainer = eui.NewVerticalContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(600, 120),
		),
	)

	st.memberContainer = eui.NewRowContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)
	st.updateMemberContainer(world)

	actionContainer := eui.NewRowContainer()
	actionContainer.AddChild(st.selectContainer, st.cardSpecContainer)
	rootContainer.AddChild(
		st.memberContainer,
		st.enemyListContainer,
		actionContainer,
	)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *BattleState) initEnemyContainer() *widget.Container {
	return eui.NewRowContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxWidth:  200,
				MaxHeight: 200,
			}),
			widget.WidgetOpts.MinSize(0, 600),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)
}

// 敵一覧を更新する
func (st *BattleState) updateEnemyListContainer(world w.World) {
	st.enemyListContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		{
			// とりあえず仮の画像
			tankSS := (*world.Resources.SpriteSheets)["front_tank1"]
			graphic := widget.NewGraphic(
				widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Stretch: true,
				})),
				widget.GraphicOpts.Image(tankSS.Texture.Image),
			)
			st.enemyListContainer.AddChild(graphic)
		}

		{
			name := gameComponents.Name.Get(entity).(*gc.Name)
			pools := gameComponents.Pools.Get(entity).(*gc.Pools)
			text := fmt.Sprintf("%s\n%3d/%3d", name.Name, pools.HP.Current, pools.HP.Max)
			st.enemyListContainer.AddChild(eui.NewMenuText(text, world))
		}
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
	res := world.Resources.UIResources
	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			v, ok := e.(policyEntry)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			return string(v)
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {}),
		euiext.ListOpts.EntryButtonOpts(),
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
		euiext.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.ImageTrans)),
	}
	list := eui.NewList(
		entries,
		opts,
		world,
	)
	st.selectContainer.AddChild(list)
}

// ================

func (st *BattleState) reloadAction(world w.World, currentPhase *phaseChooseAction) {
	st.selectContainer.RemoveChildren()
	st.cardSpecContainer.RemoveChildren()

	gameComponents := world.Components.Game.(*gc.Components)
	equipCards := []any{} // 実際にはecs.Entityが入る。Listで受け取るのが[]anyだからそうしている
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.ItemLocationEquipped,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := gameComponents.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
		if currentPhase.owner == equipped.Owner {
			equipCards = append(equipCards, entity)
		}
	}))

	// 装備がなくても詰まないようにデフォルトの攻撃手段を追加する
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Card,
		gameComponents.ItemLocationNone,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipCards = append(equipCards, entity)
	}))

	res := world.Resources.UIResources
	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			v, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := gameComponents.Name.Get(v).(*gc.Name)
			card := gameComponents.Card.Get(v).(*gc.Card)
			return fmt.Sprintf("%s(%d)", name.Name, card.Cost)
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {
			entity, ok := e.(ecs.Entity)
			if !ok {
				return
			}
			if st.selectedItem != entity {
				st.selectedItem = entity
			}
			st.cardSpecContainer.RemoveChildren()
			// 透明度つきの背景を設定したコンテナ。cardSpecContainerは背景が設定されておらず非表示にできるようになっている
			transContainer := eui.NewVerticalContainer(
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.MinSize(700, 120),
				),
				widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
			)
			views.UpdateSpec(world, transContainer, entity)
			st.cardSpecContainer.AddChild(transContainer)

			return
		}),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			cardEntity, ok := args.Entry.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			card := gameComponents.Card.Get(cardEntity).(*gc.Card)
			if card == nil {
				log.Fatal("unexpected error: entityがcardを保持していない")
			}
			st.phase = &phaseChooseTarget{
				owner: currentPhase.owner,
				way:   cardEntity,
			}
			st.cardSpecContainer.RemoveChildren() // 選択し終わったら消す。こうするより非表示にしたほうがいいかもしれない
		}),
		euiext.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.ImageTrans)),
	}
	list := eui.NewList(
		equipCards,
		opts,
		world,
	)
	st.selectContainer.AddChild(list)
	{
		members := []ecs.Entity{}
		simple.InPartyMember(world, func(entity ecs.Entity) {
			members = append(members, entity)
		})
		member := members[st.curMemberIndex]
		name := gameComponents.Name.Get(member).(*gc.Name)
		st.selectContainer.AddChild(eui.NewMenuText(name.Name, world))
	}
}

// ================

func (st *BattleState) reloadTarget(world w.World, currentPhase *phaseChooseTarget) {
	st.selectContainer.RemoveChildren()
	st.cardSpecContainer.RemoveChildren()

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 敵キャラごとにターゲット選択ボタンを作成する
		vc := eui.NewVerticalContainer()
		st.cardSpecContainer.AddChild(vc)

		btn := eui.NewItemButton(
			"選択",
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
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
					gs.BattleCommandSystem(world) // 初回実行。以降は全部消化するまでクリックで実行する
				} else {
					st.curMemberIndex = mathutil.Min(st.curMemberIndex+1, len(members)-1)
					st.phase = &phaseChooseAction{owner: members[st.curMemberIndex]}
				}
			}),
		)
		vc.AddChild(btn)

		name := gameComponents.Name.Get(entity).(*gc.Name)
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		text := fmt.Sprintf("%s\n%3d/%3d", name.Name, pools.HP.Current, pools.HP.Max)
		vc.AddChild(eui.NewMenuText(text, world))
	}))
}

// ================

func (st *BattleState) reloadMsg(world w.World) {
	st.selectContainer.RemoveChildren()
	st.cardSpecContainer.RemoveChildren()

	entries := []any{}
	for _, e := range gamelog.BattleLog.Latest(MessageCharBaseHeight) {
		entries = append(entries, e)
	}
	if st.isWaitClick {
		entries = append(entries, any("▼"))
	}

	res := world.Resources.UIResources
	opts := []euiext.ListOpt{
		euiext.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(world.Resources.ScreenDimensions.Width-20, 280),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
				Padding:            widget.NewInsetsSimple(50),
			}),
		)),
		euiext.ListOpts.SliderOpts(
			widget.SliderOpts.MinHandleSize(5),
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(4))),
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			v, ok := e.(string)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			return v
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {}),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {}),
		euiext.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.ImageTrans)),
	}

	list := eui.NewList(
		entries,
		opts,
		world,
	)
	st.selectContainer.AddChild(list)
}

// ================

func (st *BattleState) reloadExecute(world w.World) {
	st.updateEnemyListContainer(world)

	// 処理を書く...
}

// メンバー一覧を更新する
func (st *BattleState) updateMemberContainer(world w.World) {
	st.memberContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))
}
