package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
)

// DungeonSelectState はダンジョン選択のゲームステート
type DungeonSelectState struct {
	es.BaseState
	ui                   *ebitenui.UI
	menu                 *menu.Menu
	uiBuilder            *menu.UIBuilder
	keyboardInput        input.KeyboardInput
	dungeonDescContainer *widget.Container
}

func (st DungeonSelectState) String() string {
	return "DungeonSelect"
}

// State interface ================

var _ es.State = &DungeonSelectState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *DungeonSelectState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *DungeonSelectState) OnResume(_ w.World) {
	// フォーカス状態を更新
	if st.uiBuilder != nil && st.menu != nil {
		st.uiBuilder.UpdateFocus(st.menu)
	}
}

// OnStart はステートが開始される際に呼ばれる
func (st *DungeonSelectState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	st.initMenu(world)
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *DungeonSelectState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *DungeonSelectState) Update(_ w.World) es.Transition {
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}}
	}

	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *DungeonSelectState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

// initMenu はメニューコンポーネントを初期化する
func (st *DungeonSelectState) initMenu(world w.World) {
	// メニュー項目の定義（dungeonSelectTransから変換）
	items := []menu.Item{}
	for _, data := range dungeonSelectTrans {
		items = append(items, menu.Item{
			ID:          data.label, // ラベルをIDとして使用
			Label:       data.label,
			Description: data.desc,
			UserData:    data.trans,
		})
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
			// Escapeキーの処理はUpdate()で直接行うため、ここでは何もしない
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
func (st *DungeonSelectState) updateActionDescription(world w.World, index int) {
	if st.dungeonDescContainer == nil || st.menu == nil {
		return
	}

	items := st.menu.GetItems()
	if index < 0 || index >= len(items) {
		return
	}

	st.dungeonDescContainer.RemoveChildren()
	st.dungeonDescContainer.AddChild(eui.NewMenuText(items[index].Description, world))
}

func (st *DungeonSelectState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.BlackColor)),
	)

	// メニューと説明文を横並びにするためのコンテナ
	horizontalContainer := eui.NewRowContainer()

	// メニューのUIを構築
	menuContainer := st.uiBuilder.BuildUI(st.menu)

	// 説明文用のコンテナ
	st.dungeonDescContainer = eui.NewVerticalContainer()
	st.dungeonDescContainer.AddChild(eui.NewMenuText(" ", world))

	// 左側にメニュー、右側に説明文を配置
	horizontalContainer.AddChild(
		menuContainer,
		st.dungeonDescContainer,
	)

	rootContainer.AddChild(horizontalContainer)

	// 初期状態の説明文を設定
	st.updateActionDescription(world, st.menu.GetFocusedIndex())

	return &ebitenui.UI{Container: rootContainer}
}

// MEMO: まだtransは全部同じ
var dungeonSelectTrans = []struct {
	label string
	desc  string
	trans es.Transition
}{
	{
		label: "森の遺跡",
		desc:  "鬱蒼とした森の奥地にある遺跡",
		trans: es.Transition{Type: es.TransReplace, NewStateFuncs: []es.StateFactory{NewDungeonStateWithDepth(1)}},
	},
	{
		label: "山の遺跡",
		desc:  "切り立った山の洞窟にある遺跡",
		trans: es.Transition{Type: es.TransReplace, NewStateFuncs: []es.StateFactory{NewDungeonStateWithDepth(1)}},
	},
	{
		label: "塔の遺跡",
		desc:  "雲にまで届く塔を持つ遺跡",
		trans: es.Transition{Type: es.TransReplace, NewStateFuncs: []es.StateFactory{NewDungeonStateWithDepth(1)}},
	},
	{
		label: "拠点メニューに戻る",
		desc:  "",
		trans: es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}},
	},
}
