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
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/entities"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/euiext"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// BattleState は戦闘シーンのゲームステート
type BattleState struct {
	es.BaseState
	ui            *ebitenui.UI
	keyboardInput input.KeyboardInput

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

	// メニューコンポーネント
	currentMenu   *menu.Menu
	menuUIBuilder *menu.MenuUIBuilder
}

func (st BattleState) String() string {
	return "Battle"
}

// State interface ================

var _ es.State = &BattleState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *BattleState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *BattleState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *BattleState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	// MenuUIBuilderを初期化
	st.menuUIBuilder = menu.NewMenuUIBuilder(world)

	_ = worldhelper.SpawnEnemy(world, "軽戦車")
	_ = worldhelper.SpawnEnemy(world, "火の玉")

	bg := (*world.Resources.SpriteSheets)["bg_jungle1"]
	st.bg = bg.Texture.Image

	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *BattleState) OnStop(world w.World) {
	// 後片付け
	world.Manager.Join(
		world.Components.Name,
		world.Components.FactionEnemy,
		world.Components.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))
	world.Manager.Join(
		world.Components.BattleCommand,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))

	// バトルログをクリア
	gamelog.BattleLog.Flush()
}

// Update はゲームステートの更新処理を行う
func (st *BattleState) Update(world w.World) es.Transition {
	st.ui.Update()

	// キーボード入力処理
	st.handleKeyboardInput()

	// フェーズ初期化処理
	st.handlePhaseInitialization(world)

	// フェーズ更新処理
	return st.handlePhaseUpdate(world)
}

// handleKeyboardInput はキーボード選択処理を行う（選択フェーズのみ、ログ表示中は除く）
func (st *BattleState) handleKeyboardInput() {
	if !st.isWaitClick && st.currentMenu != nil {
		switch st.phase.(type) {
		case *phaseChoosePolicy, *phaseChooseAction, *phaseChooseTarget:
			st.currentMenu.Update(st.keyboardInput)
		}
	}
}

// handlePhaseInitialization はフェーズ変更時の初期化処理を行う
func (st *BattleState) handlePhaseInitialization(world w.World) {
	if st.phase == st.prevPhase {
		return
	}

	switch v := st.phase.(type) {
	case *phaseEnemyEncounter:
		st.initEnemyEncounterPhase(world)
	case *phaseChoosePolicy:
		st.initChoosePolicyPhase(world)
	case *phaseChooseAction:
		st.reloadAction(world, v)
	case *phaseChooseTarget:
		st.reloadTarget(world, v)
	case *phaseEnemyActionSelect:
		st.handleEnemyActionSelect(world)
	case *phaseExecute:
	case *phaseResult:
	case *phaseGameOver:
	}
	st.prevPhase = st.phase
}

// initChoosePolicyPhase は政策選択フェーズの初期化を行う
func (st *BattleState) initChoosePolicyPhase(world w.World) {
	var err error
	st.party, err = worldhelper.NewParty(world, gc.FactionAlly)
	if err != nil {
		log.Fatal(err)
	}
	st.reloadPolicy(world)
}

// handleEnemyActionSelect は敵アクション選択フェーズの処理を行う
func (st *BattleState) handleEnemyActionSelect(world w.World) {
	// マスタとして事前に生成されたカードエンティティをメモしておく
	masterCardEntityMap := st.buildMasterCardEntityMap(world)

	// 敵のコマンドを投入する
	st.processEnemyCommands(world, masterCardEntityMap)

	st.phase = &phaseExecute{}
}

// buildMasterCardEntityMap はマスターカードエンティティのマップを構築する
func (st *BattleState) buildMasterCardEntityMap(world w.World) map[string]ecs.Entity {
	masterCardEntityMap := map[string]ecs.Entity{}
	world.Manager.Join(
		world.Components.Name,
		world.Components.Item,
		world.Components.Card,
		world.Components.ItemLocationNone,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := world.Components.Name.Get(entity).(*gc.Name)
		masterCardEntityMap[name.Name] = entity
	}))
	return masterCardEntityMap
}

// processEnemyCommands は敵のコマンド処理を行う
func (st *BattleState) processEnemyCommands(world w.World, masterCardEntityMap map[string]ecs.Entity) {
	world.Manager.Join(
		world.Components.FactionEnemy,
		world.Components.Attributes,
		world.Components.CommandTable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		st.processEnemyCommand(world, entity, masterCardEntityMap)
	}))
}

// processEnemyCommand は単体の敵エンティティのコマンド処理を行う
func (st *BattleState) processEnemyCommand(world w.World, entity ecs.Entity, masterCardEntityMap map[string]ecs.Entity) {
	// テーブル取得
	ctComponent := world.Components.CommandTable.Get(entity).(*gc.CommandTable)
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	ct, err := rawMaster.GetCommandTable(ctComponent.Name)
	if err != nil {
		panic(fmt.Sprintf("GetCommandTable failed: %v", err))
	}
	name := ct.SelectByWeight()

	// プレイヤーキャラから選択
	allys := st.getAllyEntities(world)
	targetEntity := allys[rand.IntN(len(allys))]

	// 攻撃カードによって対象を選択
	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		BattleCommand: &gc.BattleCommand{
			Owner:  entity,
			Target: targetEntity,
			Way:    masterCardEntityMap[name],
		},
	})
	entities.AddEntities(world, cl)
}

// getAllyEntities は味方エンティティのリストを取得する
func (st *BattleState) getAllyEntities(world w.World) []ecs.Entity {
	allys := []ecs.Entity{}
	world.Manager.Join(
		world.Components.Name,
		world.Components.FactionAlly,
		world.Components.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		allys = append(allys, entity)
	}))
	return allys
}

// handlePhaseUpdate は毎回実行されるフェーズ更新処理を行う
func (st *BattleState) handlePhaseUpdate(world w.World) es.Transition {
	switch v := st.phase.(type) {
	case nil:
		st.phase = &phaseEnemyEncounter{}
	case *phaseEnemyEncounter:
		return st.handleEnemyEncounterPhase(world)
	case *phaseChoosePolicy:
	case *phaseChooseAction:
	case *phaseChooseTarget:
	case *phaseEnemyActionSelect:
	case *phaseExecute:
		return st.handleExecutePhase(world)
	case *phaseResult:
		return st.handleResultPhase(world, v)
	case *phaseGameOver:
		return st.handleGameOverPhase(world)
	}

	return st.ConsumeTransition()
}

// handleExecutePhase は実行フェーズの処理を行う
func (st *BattleState) handleExecutePhase(world w.World) es.Transition {
	st.updateEnemyListContainer(world)
	st.reloadExecute(world)
	st.reloadMsg(world)
	st.updateMemberContainer(world)

	// 戦闘終了判定
	if transition := st.checkBattleExtinction(world); transition.Type != es.TransNone {
		return transition
	}

	// コマンド実行処理
	return st.handleCommandExecution(world)
}

// checkBattleExtinction は戦闘終了条件をチェックする
func (st *BattleState) checkBattleExtinction(world w.World) es.Transition {
	switch gs.BattleExtinctionSystem(world) {
	case gs.BattleExtinctionNone:
		return es.Transition{Type: es.TransNone}
	case gs.BattleExtinctionAlly:
		gamelog.BattleLog.Append("全滅した。")
		st.phase = &phaseGameOver{}
		return es.Transition{Type: es.TransNone}
	case gs.BattleExtinctionMonster:
		gamelog.BattleLog.Append("敵を全滅させた。")
		st.phase = &phaseResult{}
		return es.Transition{Type: es.TransNone}
	default:
		return es.Transition{Type: es.TransNone}
	}
}

// handleCommandExecution はコマンド実行処理を行う
func (st *BattleState) handleCommandExecution(world w.World) es.Transition {
	commandCount := st.countBattleCommands(world)
	if commandCount > 0 {
		// 未処理のコマンドがまだ残っている
		st.isWaitClick = true
		if st.keyboardInput.IsEnterJustPressedOnce() {
			gs.BattleCommandSystem(world)
			st.isWaitClick = false
		}
		return es.Transition{Type: es.TransNone}
	}

	// 処理完了
	if st.keyboardInput.IsEnterJustPressedOnce() {
		st.phase = &phaseChoosePolicy{}
		st.isWaitClick = false
		gamelog.BattleLog.Flush()
	}
	return es.Transition{Type: es.TransNone}
}

// countBattleCommands は戦闘コマンド数をカウントする
func (st *BattleState) countBattleCommands(world w.World) int {
	commandCount := 0
	world.Manager.Join(
		world.Components.BattleCommand,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		commandCount++
	}))
	return commandCount
}

// handleResultPhase は結果フェーズの処理を行う
func (st *BattleState) handleResultPhase(world w.World, v *phaseResult) es.Transition {
	st.reloadMsg(world)

	if st.keyboardInput.IsEnterJustPressedOnce() {
		switch v.actionCount {
		case 0:
			dropResult := gs.BattleDropSystem(world)
			st.resultWindow = st.initResultWindow(world, dropResult)
			st.ui.AddWindow(st.resultWindow)
		default:
			return es.Transition{Type: es.TransPop}
		}
		v.actionCount++
	}
	return es.Transition{Type: es.TransNone}
}

// handleGameOverPhase はゲームオーバーフェーズの処理を行う
func (st *BattleState) handleGameOverPhase(world w.World) es.Transition {
	st.reloadMsg(world)

	if st.keyboardInput.IsEnterJustPressedOnce() {
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewGameOverState}}
	}
	return es.Transition{Type: es.TransNone}
}

// Draw はゲームステートの描画処理を行う
func (st *BattleState) Draw(_ w.World, screen *ebiten.Image) {
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
	world.Manager.Join(
		world.Components.Name,
		world.Components.FactionEnemy,
		world.Components.Attributes,
		world.Components.Pools,
		world.Components.Render,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		{
			pools := world.Components.Pools.Get(entity).(*gc.Pools)
			if pools.HP.Current == 0 {
				return
			}
		}
		container := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		)
		{
			render := world.Components.Render.Get(entity).(*gc.Render)
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
			name := world.Components.Name.Get(entity).(*gc.Name)
			pools := world.Components.Pools.Get(entity).(*gc.Pools)
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

	// MenuItemを作成
	items := []menu.MenuItem{
		{
			ID:       "attack",
			Label:    string(policyEntryAttack),
			UserData: policyEntryAttack,
		},
		{
			ID:       "item",
			Label:    string(policyEntryItem),
			UserData: policyEntryItem,
			Disabled: true, // TODO: 未実装のため無効化
		},
		{
			ID:       "escape",
			Label:    string(policyEntryEscape),
			UserData: policyEntryEscape,
		},
	}

	// Menuの設定
	config := menu.MenuConfig{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.MenuCallbacks{
		OnSelect: func(_ int, item menu.MenuItem) {
			entry := item.UserData.(policyEntry)
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
				st.SetTransition(es.Transition{Type: es.TransPop})
			}
		},
		OnFocusChange: func(_, _ int) {
			st.menuUIBuilder.UpdateFocus(st.currentMenu)
		},
	}

	// Menuを作成してUIを構築
	st.currentMenu = menu.NewMenu(config, callbacks)
	menuContainer := st.menuUIBuilder.BuildUI(st.currentMenu)
	st.selectContainer.AddChild(menuContainer)
}

// ================

func (st *BattleState) reloadAction(world w.World, currentPhase *phaseChooseAction) {
	st.selectContainer.RemoveChildren()
	st.cardSpecContainer.RemoveChildren()

	usableCards := []ecs.Entity{}
	unusableCards := []ecs.Entity{}
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationEquipped,
		world.Components.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		card := world.Components.Card.Get(entity).(*gc.Card)
		ownerPools := world.Components.Pools.Get(currentPhase.owner).(*gc.Pools)
		equipped := world.Components.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
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
		world.Components.Name,
		world.Components.Item,
		world.Components.Card,
		world.Components.ItemLocationNone,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := world.Components.Name.Get(entity).(*gc.Name)
		if name.Name == "体当たり" {
			usableCards = append(usableCards, entity)
		}
	}))

	// MenuItemを作成
	items := make([]menu.MenuItem, len(usableCards))
	for i, entity := range usableCards {
		name := world.Components.Name.Get(entity).(*gc.Name)
		card := world.Components.Card.Get(entity).(*gc.Card)
		items[i] = menu.MenuItem{
			ID:       fmt.Sprintf("card_%d", entity),
			Label:    fmt.Sprintf("%s (%d)", name.Name, card.Cost),
			UserData: entity,
		}
	}

	// 使用不可カードも追加（無効化状態で）
	for _, entity := range unusableCards {
		name := world.Components.Name.Get(entity).(*gc.Name)
		card := world.Components.Card.Get(entity).(*gc.Card)
		items = append(items, menu.MenuItem{
			ID:       fmt.Sprintf("card_disabled_%d", entity),
			Label:    fmt.Sprintf("%s (%d)", name.Name, card.Cost),
			UserData: entity,
			Disabled: true,
		})
	}

	if len(items) == 0 {
		return
	}

	// Menuの設定
	config := menu.MenuConfig{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.MenuCallbacks{
		OnSelect: func(_ int, item menu.MenuItem) {
			cardEntity := item.UserData.(ecs.Entity)
			card := world.Components.Card.Get(cardEntity).(*gc.Card)
			if card == nil {
				log.Fatal("unexpected error: entityがcardを保持していない")
			}
			st.phase = &phaseChooseTarget{
				owner: currentPhase.owner,
				way:   cardEntity,
			}
			st.cardSpecContainer.RemoveChildren()
		},
		OnFocusChange: func(_, newIndex int) {
			if newIndex >= 0 && newIndex < len(items) {
				entity := items[newIndex].UserData.(ecs.Entity)
				st.selectedItem = entity
				st.updateCardSpec(world, entity)
			}
			st.menuUIBuilder.UpdateFocus(st.currentMenu)
		},
	}

	// Menuを作成してUIを構築
	st.currentMenu = menu.NewMenu(config, callbacks)
	menuContainer := st.menuUIBuilder.BuildUI(st.currentMenu)
	st.selectContainer.AddChild(menuContainer)

	// 初期状態でカードの詳細を表示
	if len(usableCards) > 0 {
		st.selectedItem = usableCards[0]
		st.updateCardSpec(world, usableCards[0])
	}

	// プレイヤー名を表示
	{
		name := world.Components.Name.Get(*st.party.Value()).(*gc.Name)
		st.selectContainer.AddChild(eui.NewMenuText(name.Name, world))
	}
}

// updateCardSpec はカードの詳細情報を更新する
func (st *BattleState) updateCardSpec(world w.World, entity ecs.Entity) {
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

// ================

func (st *BattleState) reloadTarget(world w.World, currentPhase *phaseChooseTarget) {
	st.selectContainer.RemoveChildren()
	st.cardSpecContainer.RemoveChildren()

	// 生きている敵をリストアップ
	enemies := []ecs.Entity{}
	world.Manager.Join(
		world.Components.Name,
		world.Components.FactionEnemy,
		world.Components.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 生きている敵のみ対象とする
		pools := world.Components.Pools.Get(entity).(*gc.Pools)
		if pools.HP.Current == 0 {
			return
		}

		// 敵をリストに追加
		enemies = append(enemies, entity)
	}))

	if len(enemies) == 0 {
		return
	}

	// MenuItemを作成
	items := make([]menu.MenuItem, len(enemies))
	for i, entity := range enemies {
		name := world.Components.Name.Get(entity).(*gc.Name)
		items[i] = menu.MenuItem{
			ID:       fmt.Sprintf("enemy_%d", entity),
			Label:    name.Name,
			UserData: entity,
		}
	}

	// Menuの設定
	config := menu.MenuConfig{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.MenuCallbacks{
		OnSelect: func(_ int, item menu.MenuItem) {
			targetEntity := item.UserData.(ecs.Entity)
			cl := entities.ComponentList{}
			cl.Game = append(cl.Game, gc.GameComponentList{
				BattleCommand: &gc.BattleCommand{
					Owner:  currentPhase.owner,
					Target: targetEntity,
					Way:    currentPhase.way,
				},
			})
			entities.AddEntities(world, cl)

			err := st.party.Next()
			if err == nil {
				st.phase = &phaseChooseAction{owner: *st.party.Value()}
			} else {
				st.phase = &phaseEnemyActionSelect{}
				gs.BattleCommandSystem(world)
			}
		},
		OnFocusChange: func(_, _ int) {
			st.menuUIBuilder.UpdateFocus(st.currentMenu)
		},
	}

	// Menuを作成してUIを構築
	st.currentMenu = menu.NewMenu(config, callbacks)
	menuContainer := st.menuUIBuilder.BuildUI(st.currentMenu)
	st.selectContainer.AddChild(menuContainer)
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
		euiext.ListOpts.EntryEnterFunc(func(_ any) {}),
		euiext.ListOpts.EntrySelectedHandler(func(_ *euiext.ListEntrySelectedEventArgs) {}),
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
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
		world.Components.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))
}

func (st *BattleState) initResultWindow(world w.World, dropResult gs.DropResult) *widget.Window {
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
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
		world.Components.Attributes,
		world.Components.Pools,
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

		name := world.Components.Name.Get(entity).(*gc.Name)
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

// ================
// 敵遭遇フェーズ

// initEnemyEncounterPhase は敵遭遇フェーズの初期化を行う
func (st *BattleState) initEnemyEncounterPhase(_ w.World) {
	// 「敵が現れた」メッセージをログに追加
	gamelog.BattleLog.Append("敵が現れた。")

	// クリック待ち状態にする
	st.isWaitClick = true
}

// handleEnemyEncounterPhase は敵遭遇フェーズの更新処理を行う
func (st *BattleState) handleEnemyEncounterPhase(world w.World) es.Transition {
	// メッセージを表示
	st.reloadMsg(world)

	// エンターキーが押されたら次のフェーズに進む
	if st.keyboardInput.IsEnterJustPressedOnce() {
		st.isWaitClick = false
		gamelog.BattleLog.Flush() // メッセージをクリア
		st.phase = &phaseChoosePolicy{}
	}

	return es.Transition{Type: es.TransNone}
}
