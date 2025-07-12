// 拠点でのコマンド選択画面
package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type HomeMenuState struct {
	ui            *ebitenui.UI
	trans         *states.Transition
	menu          *menu.Menu
	uiBuilder     *menu.MenuUIBuilder
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

func (st *HomeMenuState) OnPause(world w.World) {}

func (st *HomeMenuState) OnResume(world w.World) {
	st.updateMemberContainer(world)
	// フォーカス状態を更新
	if st.uiBuilder != nil && st.menu != nil {
		st.uiBuilder.UpdateFocus(st.menu)
	}
}

func (st *HomeMenuState) OnStart(world w.World) {
	// デバッグ用データ初期化（初回のみ）
	worldhelper.InitDebugData(world)

	if st.keyboardInput == nil {
		st.keyboardInput = input.NewDefaultKeyboardInput()
	}

	bg := (*world.Resources.SpriteSheets)["bg_cup1"]
	st.bg = bg.Texture.Image

	st.initMenu(world)
	st.ui = st.initUI(world)
}

func (st *HomeMenuState) OnStop(world w.World) {}

func (st *HomeMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	return states.Transition{Type: states.TransNone}
}

func (st *HomeMenuState) Draw(world w.World, screen *ebiten.Image) {
	if st.bg != nil {
		screen.DrawImage(st.bg, &ebiten.DrawImageOptions{})
	}

	st.ui.Draw(screen)
}

// initMenu はメニューコンポーネントを初期化する
func (st *HomeMenuState) initMenu(world w.World) {
	// メニュー項目の定義（homeMenuTransから変換）
	items := []menu.MenuItem{
		{
			ID:          "departure",
			Label:       "出発",
			Description: "遺跡に出発する",
			UserData:    states.Transition{Type: states.TransPush, NewStates: []states.State{&DungeonSelectState{}}},
		},
		{
			ID:          "craft",
			Label:       "合成",
			Description: "アイテムを合成する",
			UserData:    states.Transition{Type: states.TransPush, NewStates: []states.State{&CraftMenuState{}}},
		},
		{
			ID:          "replace",
			Label:       "入替",
			Description: "仲間を入れ替える(未実装)",
			UserData:    states.Transition{Type: states.TransNone},
		},
		{
			ID:          "inventory",
			Label:       "所持",
			Description: "所持品を確認する",
			UserData:    states.Transition{Type: states.TransSwitch, NewStates: []states.State{&InventoryMenuState{}}},
		},
		{
			ID:          "equipment",
			Label:       "装備",
			Description: "装備を変更する",
			UserData:    states.Transition{Type: states.TransSwitch, NewStates: []states.State{&EquipMenuState{}}},
		},
		{
			ID:          "exit",
			Label:       "終了",
			Description: "タイトル画面に戻る",
			UserData:    states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}},
		},
	}

	// メニューの設定
	config := menu.MenuConfig{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.MenuCallbacks{
		OnSelect: func(index int, item menu.MenuItem) {
			// 選択されたアイテムのUserDataからTransitionを取得
			if trans, ok := item.UserData.(states.Transition); ok {
				st.trans = &trans
			}
		},
		OnCancel: func() {
			// Escapeキーが押された時の処理（タイトル画面に戻る）
			st.trans = &states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
		},
		OnFocusChange: func(oldIndex, newIndex int) {
			// フォーカス変更時に説明文を更新
			st.updateActionDescription(world, newIndex)
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
		OnHover: func(index int, item menu.MenuItem) {
			// ホバー時に説明文を更新
			st.updateActionDescription(world, index)
		},
	}

	// メニューを作成
	st.menu = menu.NewMenu(config, callbacks)

	// UIビルダーを作成
	st.uiBuilder = menu.NewMenuUIBuilder(world)
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
	st.actionDescContainer.AddChild(eui.NewMenuText(items[index].Description, world))
	st.updateMemberContainer(world)
}

// ================

func (st *HomeMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()
	st.memberContainer = eui.NewRowContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)

	st.actionDescContainer = eui.NewRowContainer()
	st.actionDescContainer.AddChild(eui.NewMenuText(" ", world))

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
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))
}
