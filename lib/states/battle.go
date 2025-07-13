package states

import (
	"fmt"
	"image"
	"log"
	"math/rand/v2"

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
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/euiext"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/systems"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type BattleState struct {
	states.BaseState
	ui *ebitenui.UI

	// 現在のサブステート
	phase battlePhase
	// 1個前のサブステート。サブステート変更を検知するのに使う
	prevPhase battlePhase
	// テキスト送り待ち状態
	isWaitClick bool
	// 味方パーティ
	party worldhelper.Party

	// 背景
	bg *ebiten.Image
	// 敵表示コンテナ
	enemyListContainer *widget.Container
	// 各フェーズでの選択表示に使うコンテナ
	selectContainer *widget.Container
	// 味方表示コンテナ
	memberContainer *widget.Container
	// 結果ウィンドウ
	resultWindow *widget.Window

	// 選択中のアイテム
	selectedItem ecs.Entity
	// カードの説明表示コンテナ
	cardSpecContainer *widget.Container

	// キーボード選択用フィールド
	currentSelection int   // 現在の選択インデックス
	selectionItems   []any // 選択可能な項目リスト
}

func (st BattleState) String() string {
	return "Battle"
}

// State interface ================

var _ es.State = &BattleState{}

func (st *BattleState) OnPause(world w.World) {}

func (st *BattleState) OnResume(world w.World) {}

func (st *BattleState) OnStart(world w.World) {
	_ = worldhelper.SpawnEnemy(world, "軽戦車")
	_ = worldhelper.SpawnEnemy(world, "火の玉")

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
	world.Manager.Join(
		gameComponents.BattleCommand,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))
	gamelog.BattleLog.Flush()

	// FIXME: state transition: popで削除されてくれない。stateインスタンスが使い回されているように見える
	st.phase = nil
}

func (st *BattleState) Update(world w.World) states.Transition {
	st.ui.Update()

	// キーボード選択処理（ログ表示中は除く）
	if !st.isWaitClick {
		st.handleKeyboardSelection(world)
	}

	// ステートが変わった最初の1回だけ実行される
	// TODO: 複雑化しているので、ここもstates packageを使ってやったほうがいい
	if st.phase != st.prevPhase {
		switch v := st.phase.(type) {
		case *phaseChoosePolicy:
			var err error
			st.party, err = worldhelper.NewParty(world, components.FactionAlly)
			if err != nil {
				log.Fatal(err)
			}
			st.reloadPolicy(world)
		case *phaseChooseAction:
			st.reloadAction(world, v)
		case *phaseChooseTarget:
			st.reloadTarget(world, v)
		case *phaseEnemyActionSelect:
			gameComponents := world.Components.Game.(*gc.Components)

			// マスタとして事前に生成されたカードエンティティをメモしておく
			masterCardEntityMap := map[string]ecs.Entity{}
			world.Manager.Join(
				gameComponents.Name,
				gameComponents.Item,
				gameComponents.Card,
				gameComponents.ItemLocationNone,
			).Visit(ecs.Visit(func(entity ecs.Entity) {
				name := gameComponents.Name.Get(entity).(*gc.Name)
				masterCardEntityMap[name.Name] = entity
			}))

			// 敵のコマンドを投入する
			// UIはない。1回だけ実行して敵コマンドを投入し、次のステートにいく
			world.Manager.Join(
				gameComponents.FactionEnemy,
				gameComponents.Attributes,
				gameComponents.CommandTable,
			).Visit(ecs.Visit(func(entity ecs.Entity) {
				// テーブル取得
				ctComponent := gameComponents.CommandTable.Get(entity).(*gc.CommandTable)
				rawMaster := world.Resources.RawMaster.(*raw.RawMaster)
				ct := rawMaster.GetCommandTable(ctComponent.Name)
				name := ct.SelectByWeight()

				// プレイヤーキャラから選択
				allys := []ecs.Entity{}
				world.Manager.Join(
					gameComponents.Name,
					gameComponents.FactionAlly,
					gameComponents.Pools,
				).Visit(ecs.Visit(func(entity ecs.Entity) {
					allys = append(allys, entity)
				}))
				targetEntity := allys[rand.IntN(len(allys))]

				// 攻撃カードによって対象を選択
				cl := loader.EntityComponentList{}
				cl.Game = append(cl.Game, components.GameComponentList{
					BattleCommand: &gc.BattleCommand{
						Owner:  entity,
						Target: targetEntity,
						Way:    masterCardEntityMap[name],
					},
				})
				loader.AddEntities(world, cl)
			}))

			st.phase = &phaseExecute{}
		case *phaseExecute:
		case *phaseResult:
		case *phaseGameOver:
		}
		st.prevPhase = st.phase
	}

	// 毎回実行される
	switch v := st.phase.(type) {
	case nil:
		// 戦闘ステート開始直後に実行される
		st.phase = &phaseChoosePolicy{}
	case *phaseChoosePolicy:
	case *phaseChooseAction:
	case *phaseChooseTarget:
	case *phaseEnemyActionSelect:
	case *phaseExecute:
		effects.RunEffectQueue(world)
		st.updateEnemyListContainer(world)
		st.reloadExecute(world)
		st.reloadMsg(world)
		st.updateMemberContainer(world)

		switch systems.BattleExtinctionSystem(world) {
		case systems.BattleExtinctionNone:
		case systems.BattleExtinctionAlly:
			gamelog.BattleLog.Append("全滅した。")
			st.phase = &phaseGameOver{}
			return states.Transition{Type: states.TransNone}
		case systems.BattleExtinctionMonster:
			gamelog.BattleLog.Append("敵を全滅させた。")
			st.phase = &phaseResult{}
			return states.Transition{Type: states.TransNone}
		}

		gameComponents := world.Components.Game.(*gc.Components)

		commandCount := 0
		world.Manager.Join(
			gameComponents.BattleCommand,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			commandCount += 1
		}))
		if commandCount > 0 {
			// 未処理のコマンドがまだ残っている
			st.isWaitClick = true
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				gs.BattleCommandSystem(world)
				st.isWaitClick = false
			}
			return states.Transition{Type: states.TransNone}
		}

		// 処理完了
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			st.phase = &phaseChoosePolicy{}
			st.isWaitClick = false
			gamelog.BattleLog.Flush()
		}
	case *phaseResult:
		st.reloadMsg(world)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch v.actionCount {
			case 0:
				dropResult := systems.BattleDropSystem(world)
				st.resultWindow = st.initResultWindow(world, dropResult)
				st.ui.AddWindow(st.resultWindow)
			default:
				return states.Transition{Type: states.TransPop}
			}
			v.actionCount += 1
		}
	case *phaseGameOver:
		st.reloadMsg(world)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameOverState{}}}
		}
	}

	return st.ConsumeTransition()
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
	st.enemyListContainer = eui.NewRowContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxWidth:  400,
				MaxHeight: 200,
			}),
			widget.WidgetOpts.MinSize(0, 600),
		),
	)
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

// 敵一覧を更新する
func (st *BattleState) updateEnemyListContainer(world w.World) {
	st.enemyListContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Attributes,
		gameComponents.Pools,
		gameComponents.Render,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		{
			pools := gameComponents.Pools.Get(entity).(*gc.Pools)
			if pools.HP.Current == 0 {
				return
			}
		}
		container := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		)
		{
			render := gameComponents.Render.Get(entity).(*gc.Render)
			sheets := (*world.Resources.SpriteSheets)[render.BattleBody.SheetName]
			graphic := widget.NewGraphic(
				widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Stretch: true,
				})),
				widget.GraphicOpts.Image(sheets.Texture.Image),
			)
			container.AddChild(graphic)
		}
		{
			name := gameComponents.Name.Get(entity).(*gc.Name)
			pools := gameComponents.Pools.Get(entity).(*gc.Pools)
			text := fmt.Sprintf("%s\n%3d/%3d", name.Name, pools.HP.Current, pools.HP.Max)
			container.AddChild(eui.NewMenuText(text, world))
		}

		st.enemyListContainer.AddChild(container)
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

	// キーボード選択の初期化
	st.selectionItems = entries
	st.currentSelection = 0
	res := world.Resources.UIResources
	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			v, ok := e.(policyEntry)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}

			// 選択中の項目にマーカーを表示
			label := string(v)
			for i, item := range st.selectionItems {
				if item == e && i == st.currentSelection {
					label = "▶ " + label
					break
				}
			}
			return label
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {}),
		euiext.ListOpts.EntryButtonOpts(),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			entry := args.Entry.(policyEntry)
			switch entry {
			case policyEntryAttack:
				members := []ecs.Entity{}
				worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
					members = append(members, entity)
				})
				st.phase = &phaseChooseAction{owner: *st.party.Value()}
			case policyEntryItem:
				// TODO: 未実装
			case policyEntryEscape:
				st.SetTransition(states.Transition{Type: states.TransPop})
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
	usableCards := []any{}   // 実際にはecs.Entityが入る。Listで受け取るのが[]anyだからそうしている
	unusableCards := []any{} // 実際にはecs.Entityが入る。Listで受け取るのが[]anyだからそうしている
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.ItemLocationEquipped,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		card := gameComponents.Card.Get(entity).(*gc.Card)
		ownerPools := gameComponents.Pools.Get(currentPhase.owner).(*gc.Pools)
		equipped := gameComponents.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
		if currentPhase.owner == equipped.Owner {
			if ownerPools.SP.Current >= card.Cost {
				// 使用可能
				usableCards = append(usableCards, entity)
			} else {
				// 使用不可
				unusableCards = append(unusableCards, entity)
			}
		}
	}))

	// 装備がなくても詰まないようにデフォルトの攻撃手段を追加する
	// TODO: わかりにくいのでコンポーネント化したほうがいいかも
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Item,
		gameComponents.Card,
		gameComponents.ItemLocationNone,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		if name.Name == "体当たり" {
			usableCards = append(usableCards, entity)
		}
	}))

	// キーボード選択の初期化（使用可能カードのみ）
	st.selectionItems = usableCards
	st.currentSelection = 0
	if len(usableCards) > 0 {
		// 最初のカードの詳細を表示
		entity := usableCards[0].(ecs.Entity)
		st.selectedItem = entity
		st.cardSpecContainer.RemoveChildren()

		res := world.Resources.UIResources
		transContainer := eui.NewVerticalContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(700, 120),
			),
			widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
		)
		views.UpdateSpec(world, transContainer, entity)
		st.cardSpecContainer.AddChild(transContainer)
	}

	res := world.Resources.UIResources
	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			v, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := gameComponents.Name.Get(v).(*gc.Name)
			card := gameComponents.Card.Get(v).(*gc.Card)

			// 選択中の項目にマーカーを表示
			label := fmt.Sprintf("%s (%d)", name.Name, card.Cost)
			for i, item := range st.selectionItems {
				if item == e && i == st.currentSelection {
					label = "▶ " + label
					break
				}
			}
			return label
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
		usableCards,
		opts,
		world,
	)
	st.selectContainer.AddChild(list)

	if len(unusableCards) > 0 {
		notAvaliableList := eui.NewList(
			unusableCards,
			opts,
			world,
		)
		notAvaliableList.GetWidget().Disabled = true
		st.selectContainer.AddChild(notAvaliableList)
	}

	{
		name := gameComponents.Name.Get(*st.party.Value()).(*gc.Name)
		st.selectContainer.AddChild(eui.NewMenuText(name.Name, world))
	}
}

// ================

func (st *BattleState) reloadTarget(world w.World, currentPhase *phaseChooseTarget) {
	gameComponents := world.Components.Game.(*gc.Components)

	// 生きている敵をリストアップ
	enemies := []any{}
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 生きている敵のみ対象とする
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		if pools.HP.Current == 0 {
			return
		}

		// 敵をリストに追加
		enemies = append(enemies, entity)
	}))

	// キーボード選択の初期化
	st.selectionItems = enemies
	st.currentSelection = 0

	// UIを表示（カーソルも含む）
	st.reloadTargetUI(world, currentPhase)
}

// ================

func (st *BattleState) reloadMsg(world w.World) {
	st.selectContainer.RemoveChildren()
	st.cardSpecContainer.RemoveChildren()

	entries := []any{}
	for _, e := range gamelog.BattleLog.Get() {
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
		gameComponents.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))
}

func (st *BattleState) initResultWindow(world w.World, dropResult systems.DropResult) *widget.Window {
	res := world.Resources.UIResources
	const width = 800
	const height = 400
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	content := eui.NewWindowContainer(world)
	// TODO: 経験値をプラスする
	// EXPが0~100まであり、100に到達するとレベルを1上げ、EXPを0に戻す
	// 獲得経験値は、相手の種別ランクとレベル差によって決まる
	content.AddChild(widget.NewText(widget.TextOpts.Text("経験", res.Text.TitleFace, styles.TextColor)))
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
		gameComponents.Attributes,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entryContainer := eui.NewRowContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(200, 0),
			),
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    0,
					Bottom: 0,
					Left:   20,
				}),
			)),
		)
		content.AddChild(entryContainer)

		name := gameComponents.Name.Get(entity).(*gc.Name)
		entryContainer.AddChild(
			widget.NewText(
				widget.TextOpts.Text(name.Name, res.Text.Face, styles.TextColor),
				widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)),
			),
		)
		xpBefore := dropResult.XPBefore[entity]
		xpAfter := dropResult.XPAfter[entity]
		entryContainer.AddChild(
			widget.NewText(
				widget.TextOpts.Text(fmt.Sprintf("%d → %d", xpBefore, xpAfter), res.Text.Face, styles.TextColor),
				widget.TextOpts.Position(widget.TextPositionEnd, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)),
			),
		)
		if dropResult.IsLevelUp[entity] {
			entryContainer.AddChild(
				widget.NewText(
					widget.TextOpts.Text("Lv ↑", res.Text.Face, styles.TextColor),
					widget.TextOpts.Position(widget.TextPositionEnd, widget.TextPositionCenter),
					widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)),
				),
			)
		}
	}))

	content.AddChild(widget.NewText(widget.TextOpts.Text("物品", res.Text.TitleFace, styles.TextColor)))
	for _, mn := range dropResult.MaterialNames {
		text := fmt.Sprintf("  %s", mn)
		content.AddChild(widget.NewText(widget.TextOpts.Text(text, res.Text.Face, styles.TextColor)))
	}
	resultWindow := widget.NewWindow(
		widget.WindowOpts.Contents(content),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.MinSize(width, height),
		widget.WindowOpts.MaxSize(width, height),
	)
	rect := image.Rect(0, 0, screenWidth/2+width/2, screenHeight/2+height/2)
	rect = rect.Add(image.Point{screenWidth/2 - width/2, screenHeight/2 - height/2})
	resultWindow.SetLocation(rect)

	return resultWindow
}

// handleKeyboardSelection はキーボードでの選択処理を行う
func (st *BattleState) handleKeyboardSelection(world w.World) {
	// 選択可能な項目がない場合は何もしない
	if len(st.selectionItems) == 0 {
		return
	}

	// 上下キーで選択を移動
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		st.currentSelection--
		if st.currentSelection < 0 {
			st.currentSelection = len(st.selectionItems) - 1
		}
		st.updateSelectionDisplay(world)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		st.currentSelection++
		if st.currentSelection >= len(st.selectionItems) {
			st.currentSelection = 0
		}
		st.updateSelectionDisplay(world)
	}

	// エンターキーで決定
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		st.executeSelection(world)
	}
}

// updateSelectionDisplay は選択状態の表示を更新する
func (st *BattleState) updateSelectionDisplay(world w.World) {
	switch p := st.phase.(type) {
	case *phaseChoosePolicy:
		// ポリシー選択のUIを再描画（キーボード選択初期化をスキップ）
		st.reloadPolicyUI(world)
	case *phaseChooseAction:
		// アクション選択のUIを再描画（キーボード選択初期化をスキップ）
		st.reloadActionUI(world, p)
	case *phaseChooseTarget:
		// ターゲット選択のUIを再描画（キーボード選択初期化をスキップ）
		st.reloadTargetUI(world, p)
	}
}

// executeSelection は現在の選択を実行する
func (st *BattleState) executeSelection(world w.World) {
	if st.currentSelection < 0 || st.currentSelection >= len(st.selectionItems) {
		return
	}

	switch p := st.phase.(type) {
	case *phaseChoosePolicy:
		entry := st.selectionItems[st.currentSelection].(policyEntry)
		switch entry {
		case policyEntryAttack:
			members := []ecs.Entity{}
			worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
				members = append(members, entity)
			})
			st.phase = &phaseChooseAction{owner: *st.party.Value()}
		case policyEntryItem:
			// TODO: 未実装
		case policyEntryEscape:
			st.SetTransition(states.Transition{Type: states.TransPop})
		}
	case *phaseChooseAction:
		cardEntity := st.selectionItems[st.currentSelection].(ecs.Entity)
		st.phase = &phaseChooseTarget{
			owner: p.owner,
			way:   cardEntity,
		}
		st.cardSpecContainer.RemoveChildren()
	case *phaseChooseTarget:
		entity := st.selectionItems[st.currentSelection].(ecs.Entity)

		cl := loader.EntityComponentList{}
		cl.Game = append(cl.Game, components.GameComponentList{
			BattleCommand: &gc.BattleCommand{
				Owner:  p.owner,
				Target: entity,
				Way:    p.way,
			},
		})
		loader.AddEntities(world, cl)

		err := st.party.Next()
		if err == nil {
			// 次のメンバー
			st.phase = &phaseChooseAction{owner: *st.party.Value()}
		} else {
			// 全員分完了
			st.phase = &phaseEnemyActionSelect{}
			gs.BattleCommandSystem(world) // 初回実行
		}
	}
}

// reloadPolicyUI はポリシー選択UIのみを再描画する（キーボード選択初期化なし）
func (st *BattleState) reloadPolicyUI(world w.World) {
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

			// 選択中の項目にマーカーを表示
			label := string(v)
			for i, item := range st.selectionItems {
				if item == e && i == st.currentSelection {
					label = "▶ " + label
					break
				}
			}
			return label
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {}),
		euiext.ListOpts.EntryButtonOpts(),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			entry := args.Entry.(policyEntry)
			switch entry {
			case policyEntryAttack:
				members := []ecs.Entity{}
				worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
					members = append(members, entity)
				})
				st.phase = &phaseChooseAction{owner: *st.party.Value()}
			case policyEntryItem:
				// TODO: 未実装
			case policyEntryEscape:
				st.SetTransition(states.Transition{Type: states.TransPop})
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

// reloadActionUI はアクション選択UIのみを再描画する（キーボード選択初期化なし）
func (st *BattleState) reloadActionUI(world w.World, currentPhase *phaseChooseAction) {
	st.selectContainer.RemoveChildren()

	gameComponents := world.Components.Game.(*gc.Components)
	usableCards := []any{}

	// 使用可能カードを再収集（st.selectionItemsから取得）
	for _, item := range st.selectionItems {
		usableCards = append(usableCards, item)
	}

	res := world.Resources.UIResources
	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			v, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := gameComponents.Name.Get(v).(*gc.Name)
			card := gameComponents.Card.Get(v).(*gc.Card)

			// 選択中の項目にマーカーを表示
			label := fmt.Sprintf("%s (%d)", name.Name, card.Cost)
			for i, item := range st.selectionItems {
				if item == e && i == st.currentSelection {
					label = "▶ " + label
					break
				}
			}
			return label
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
			transContainer := eui.NewVerticalContainer(
				widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(700, 120)),
				widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
			)
			views.UpdateSpec(world, transContainer, entity)
			st.cardSpecContainer.AddChild(transContainer)
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
			st.cardSpecContainer.RemoveChildren()
		}),
		euiext.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.ImageTrans)),
	}
	list := eui.NewList(
		usableCards,
		opts,
		world,
	)
	st.selectContainer.AddChild(list)

	// 現在選択されているカードの詳細を表示
	if st.currentSelection >= 0 && st.currentSelection < len(st.selectionItems) {
		entity := st.selectionItems[st.currentSelection].(ecs.Entity)
		st.selectedItem = entity
		st.cardSpecContainer.RemoveChildren()

		transContainer := eui.NewVerticalContainer(
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(700, 120)),
			widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
		)
		views.UpdateSpec(world, transContainer, entity)
		st.cardSpecContainer.AddChild(transContainer)
	}

	{
		name := gameComponents.Name.Get(*st.party.Value()).(*gc.Name)
		st.selectContainer.AddChild(eui.NewMenuText(name.Name, world))
	}
}

// reloadTargetUI はターゲット選択UIのみを再描画する（キーボード選択初期化なし）
func (st *BattleState) reloadTargetUI(world w.World, currentPhase *phaseChooseTarget) {
	st.selectContainer.RemoveChildren()
	st.cardSpecContainer.RemoveChildren()

	gameComponents := world.Components.Game.(*gc.Components)

	// 敵エンティティを再収集（st.selectionItemsから取得）
	for i, item := range st.selectionItems {
		entity := item.(ecs.Entity)

		vc := eui.NewVerticalContainer()
		st.cardSpecContainer.AddChild(vc)

		// 選択中の項目を視覚的に区別
		buttonText := "選択"
		if i == st.currentSelection {
			buttonText = "▶ " + buttonText
		}

		btn := eui.NewButton(
			buttonText,
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

				err := st.party.Next()
				if err == nil {
					st.phase = &phaseChooseAction{owner: *st.party.Value()}
				} else {
					st.phase = &phaseEnemyActionSelect{}
					gs.BattleCommandSystem(world)
				}
			}),
		)
		vc.AddChild(btn)

		name := gameComponents.Name.Get(entity).(*gc.Name)
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		text := fmt.Sprintf("%s\n%3d/%3d", name.Name, pools.HP.Current, pools.HP.Max)
		vc.AddChild(eui.NewMenuText(text, world))
	}
}
