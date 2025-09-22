package states

import (
	"log"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/colors"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
)

// DungeonMenuState はダンジョン内メニューのゲームステート
type DungeonMenuState struct {
	es.BaseState[w.World]
	ui            *ebitenui.UI
	menu          *menu.Menu
	uiBuilder     *menu.UIBuilder
	keyboardInput input.KeyboardInput
}

func (st DungeonMenuState) String() string {
	return "DungeonMenu"
}

// State interface ================

var _ es.State[w.World] = &DungeonMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *DungeonMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *DungeonMenuState) OnResume(_ w.World) {
	// フォーカス状態を更新
	if st.uiBuilder != nil && st.menu != nil {
		st.uiBuilder.UpdateFocus(st.menu)
	}
}

// OnStart はステートが開始される際に呼ばれる
func (st *DungeonMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	st.initMenu(world)
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *DungeonMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *DungeonMenuState) Update(_ w.World) es.Transition[w.World] {
	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *DungeonMenuState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

// initMenu はメニューコンポーネントを初期化する
func (st *DungeonMenuState) initMenu(world w.World) {
	// メニュー項目の定義
	items := []menu.Item{
		{
			ID:          "close",
			Label:       TextClose,
			Description: "メニューを閉じる",
			UserData: func() {
				st.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			},
		},
		{
			// デバッグ用。本来の脱出パッドと同じ処理で拠点に戻る
			ID:          "exit",
			Label:       "脱出",
			Description: "脱出する",
			UserData: func() {
				// DungeonリソースにStateEventを設定
				gameResources := world.Resources.Dungeon
				gameResources.SetStateEvent(resources.StateEventWarpEscape)

				// DungeonStateに戻る
				// メニューから戻ってDungeonStateにいかないと、DungeonStateのOnStopが呼ばれず、エンティティの解放漏れが起こる
				st.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			},
		},
	}

	// メニュー設定
	config := menu.Config{
		Items:             items,
		InitialIndex:      0,
		WrapNavigation:    true,
		Orientation:       menu.Vertical,
		ItemsPerPage:      20,
		ShowPageIndicator: true,
	}

	// コールバック設定
	callbacks := menu.Callbacks{
		OnSelect: func(_ int, item menu.Item) {
			switch userData := item.UserData.(type) {
			case func():
				userData()
			default:
				log.Fatal("想定していないデータ形式が指定された")
			}
		},
		OnCancel: func() {
			st.SetTransition(es.Transition[w.World]{Type: es.TransPop})
		},
		OnFocusChange: func(_, _ int) {
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
	}

	st.menu = menu.NewMenu(config, callbacks)
}

func (st *DungeonMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := styled.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(colors.BlackColor)),
	)

	// UIビルダーを使用してメニューUIを構築
	st.uiBuilder = menu.NewUIBuilder(world)
	menuContainer := st.uiBuilder.BuildUI(st.menu)
	rootContainer.AddChild(menuContainer)

	return &ebitenui.UI{Container: rootContainer}
}
