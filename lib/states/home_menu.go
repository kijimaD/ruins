// Package states は拠点でのコマンド選択画面
package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/config"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/widgets/common"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// HomeMenuState は拠点メニューのゲームステート
type HomeMenuState struct {
	es.BaseState
	ui            *ebitenui.UI
	menu          *menu.Menu
	uiBuilder     *menu.UIBuilder
	keyboardInput input.KeyboardInput

	// 背景
	bg *ebiten.Image

	// メンバー一覧コンテナ
	memberContainer *widget.Container
	// 選択肢アクションの説明コンテナ
	actionDescContainer *widget.Container
}

func (st HomeMenuState) String() string {
	return "HomeMenu"
}

// State interface ================

var _ es.State = &HomeMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *HomeMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *HomeMenuState) OnResume(world w.World) {
	st.updateMemberContainer(world)
	// フォーカス状態を更新
	if st.uiBuilder != nil && st.menu != nil {
		st.uiBuilder.UpdateFocus(st.menu)
	}
}

// OnStart はステートが開始される際に呼ばれる
func (st *HomeMenuState) OnStart(world w.World) {
	// デバッグ用データ初期化（初回のみ）
	worldhelper.InitDebugData(world)

	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	bg := (*world.Resources.SpriteSheets)["bg_cup1"]
	st.bg = bg.Texture.Image

	st.initMenu(world)
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *HomeMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *HomeMenuState) Update(_ w.World) es.Transition {
	cfg := config.MustGet()
	if cfg.Debug && inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewDebugMenuState}}
	}

	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *HomeMenuState) Draw(_ w.World, screen *ebiten.Image) {
	if st.bg != nil {
		screen.DrawImage(st.bg, &ebiten.DrawImageOptions{})
	}

	st.ui.Draw(screen)
}

// initMenu はメニューコンポーネントを初期化する
func (st *HomeMenuState) initMenu(world w.World) {
	// メニュー項目の定義（homeMenuTransから変換）
	items := []menu.Item{
		{
			ID:          "departure",
			Label:       "出発",
			Description: "遺跡に出発する",
			UserData:    es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewDungeonSelectState}},
		},
		{
			ID:          "craft",
			Label:       "合成",
			Description: "アイテムを合成する",
			UserData:    es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewCraftMenuState}},
		},
		{
			ID:          "replace",
			Label:       "入替",
			Description: "仲間を入れ替える",
			UserData:    es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewPartySetupState}},
		},
		{
			ID:          "inventory",
			Label:       "所持",
			Description: "所持品を確認する",
			UserData:    es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewInventoryMenuState}},
		},
		{
			ID:          "equipment",
			Label:       "装備",
			Description: "装備を変更する",
			UserData:    es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewEquipMenuState}},
		},
		{
			ID:          "save",
			Label:       "書込",
			Description: "ゲームを保存する",
			UserData:    es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewSaveMenuState}},
		},
		{
			ID:          "exit",
			Label:       "終了",
			Description: "タイトル画面に戻る",
			UserData:    es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}},
		},
	}

	// メニューの設定
	config := menu.Config{
		Items:             items,
		InitialIndex:      0,
		WrapNavigation:    true,
		Orientation:       menu.Vertical,
		ItemsPerPage:      10,
		ShowPageIndicator: true,
	}

	// コールバックの設定
	callbacks := menu.Callbacks{
		OnSelect: func(_ int, item menu.Item) {
			// 選択されたアイテムのUserDataからTransitionを取得
			if trans, ok := item.UserData.(es.Transition); ok {
				st.SetTransition(trans)
			}
		},
		OnCancel: func() {
			// Escapeキーが押された時の処理（タイトル画面に戻る）
			st.SetTransition(es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}})
		},
		OnFocusChange: func(_, newIndex int) {
			// フォーカス変更時に説明文を更新
			st.updateActionDescription(world, newIndex)
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
		OnHover: func(index int, _ menu.Item) {
			// ホバー時に説明文を更新
			st.updateActionDescription(world, index)
		},
	}

	// メニューを作成
	st.menu = menu.NewMenu(config, callbacks)

	// UIビルダーを作成
	st.uiBuilder = menu.NewUIBuilder(world)
}

// updateActionDescription は選択された項目の説明文を更新する
func (st *HomeMenuState) updateActionDescription(world w.World, index int) {
	if st.actionDescContainer == nil || st.menu == nil {
		return
	}

	items := st.menu.GetItems()
	if index < 0 || index >= len(items) {
		return
	}

	st.actionDescContainer.RemoveChildren()
	st.actionDescContainer.AddChild(common.NewMenuText(items[index].Description, world))
	st.updateMemberContainer(world)
}

// ================

func (st *HomeMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := common.NewVerticalContainer()
	st.memberContainer = common.NewRowContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)

	st.actionDescContainer = common.NewRowContainer()
	st.actionDescContainer.AddChild(common.NewMenuText(" ", world))

	// メニューのUIを構築してコンテナに追加
	menuContainer := st.uiBuilder.BuildUI(st.menu)

	rootContainer.AddChild(
		st.memberContainer,
		st.actionDescContainer,
		menuContainer,
	)

	st.updateMemberContainer(world)
	// 初期状態の説明文を設定
	st.updateActionDescription(world, st.menu.GetFocusedIndex())

	return &ebitenui.UI{Container: rootContainer}
}

// メンバー一覧を更新する
func (st *HomeMenuState) updateMemberContainer(world w.World) {
	st.memberContainer.RemoveChildren()
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))
}
