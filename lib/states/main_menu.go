package states

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/inputmapper"
	"github.com/kijimaD/ruins/lib/mapplanner"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
)

// MainMenuState は新しいメニューコンポーネントを使用するメインメニュー
type MainMenuState struct {
	es.BaseState[w.World]
	ui        *ebitenui.UI
	menu      *menu.Menu
	uiBuilder *menu.UIBuilder
}

func (st MainMenuState) String() string {
	return "MainMenu"
}

// State interface ================

var _ es.State[w.World] = &MainMenuState{}
var _ es.ActionHandler[w.World] = &MainMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MainMenuState) OnPause(_ w.World) error { return nil }

// OnResume はステートが再開される際に呼ばれる
func (st *MainMenuState) OnResume(_ w.World) error { return nil }

// OnStart はステート開始時の処理を行う
func (st *MainMenuState) OnStart(world w.World) error {
	st.initMenu(world)
	st.ui = st.initUI(world)
	return nil
}

// OnStop はステートが停止される際に呼ばれる
func (st *MainMenuState) OnStop(_ w.World) error { return nil }

// Update はゲームステートの更新処理を行う
func (st *MainMenuState) Update(_ w.World) (es.Transition[w.World], error) {
	// メニューの更新（キーボード入力→Action変換は Menu 内部で実施）
	if err := st.menu.Update(); err != nil {
		return es.Transition[w.World]{Type: es.TransNone}, err
	}

	if st.ui != nil {
		st.ui.Update()
	}

	// BaseStateの共通処理を使用
	return st.ConsumeTransition(), nil
}

// Draw はスクリーンに描画する
func (st *MainMenuState) Draw(world w.World, screen *ebiten.Image) error {
	bg := (*world.Resources.SpriteSheets)["bg_title1"]
	screen.DrawImage(bg.Texture.Image, nil)

	st.ui.Draw(screen)
	return nil
}

// ================

// HandleInput はキー入力をActionに変換する
func (st *MainMenuState) HandleInput() (inputmapper.ActionID, bool) {
	// 未使用
	return "", false
}

// DoAction はActionを実行する
func (st *MainMenuState) DoAction(_ w.World, action inputmapper.ActionID) (es.Transition[w.World], error) {
	switch action {
	case inputmapper.ActionMenuCancel:
		// メインメニューでのキャンセルは終了
		return es.Transition[w.World]{Type: es.TransQuit}, nil
	default:
		return es.Transition[w.World]{}, fmt.Errorf("未知のアクション: %s", action)
	}
}

// ================

// initMenu はメニューコンポーネントを初期化する
func (st *MainMenuState) initMenu(world w.World) {
	// メニュー項目の定義
	items := []menu.Item{
		{
			ID:       "town",
			Label:    "開始",
			UserData: es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeTown))}},
		},
		{
			ID:       "load",
			Label:    "読込",
			UserData: es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewLoadMenuState}},
		},
		{
			ID:       "exit",
			Label:    "終了",
			UserData: es.Transition[w.World]{Type: es.TransQuit},
		},
	}

	// メニューの設定
	config := menu.Config{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
		ItemsPerPage:   10,
	}

	// コールバックの設定
	callbacks := menu.Callbacks{
		OnSelect: func(_ int, item menu.Item) error {
			// 選択されたアイテムのUserDataからTransitionを取得
			if trans, ok := item.UserData.(es.Transition[w.World]); ok {
				st.SetTransition(trans)
			}
			return nil
		},
		OnCancel: func() {
			// Escapeキーが押された時の処理
			st.SetTransition(es.Transition[w.World]{Type: es.TransQuit})
		},
		OnFocusChange: func(_, _ int) {
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
	}

	// メニューを作成
	st.menu = menu.NewMenu(config, callbacks)

	// UIビルダーを作成
	st.uiBuilder = menu.NewUIBuilder(world)
}

// initUI はUIを初期化する
func (st *MainMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// メニューのUIを構築してコンテナに追加
	menuContainer := st.uiBuilder.BuildUI(st.menu)

	// 深い金色/琥珀色
	amberColor := color.NRGBA{R: 255, G: 191, B: 0, A: 255}

	// ゲームタイトル「Ruins」のテキストを作成
	titleText := widget.NewText(
		widget.TextOpts.Text("Ruins", &world.Resources.UIResources.Text.HugeTitleFace, amberColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding: &widget.Insets{
					Top: 100, // 画面上部から100ピクセル下に配置
				},
			}),
		),
	)

	// バージョン表示テキストを作成
	versionInfo := []string{}
	if consts.AppVersion != "v0.0.0" {
		versionInfo = append(versionInfo, consts.AppVersion)
	}
	if consts.AppCommit != "0000000" {
		versionInfo = append(versionInfo, consts.AppCommit)
	}
	if consts.AppDate != "0000-00-00" {
		versionInfo = append(versionInfo, consts.AppDate)
	}
	versionText := widget.NewText(
		widget.TextOpts.Text(strings.Join(versionInfo, "\n"), &world.Resources.UIResources.Text.SmallFace, consts.SecondaryColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				Padding: &widget.Insets{
					Right:  20, // 画面右端から20ピクセル左に配置
					Bottom: 20, // 画面下端から20ピクセル上に配置
				},
			}),
		),
	)

	// ラッパーコンテナを作成(メニューの位置指定のため)
	wrapperContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding: &widget.Insets{
					Top: 400, // メニューを下寄りにする
				},
			}),
		),
	)

	// メニューコンテナをラッパーに追加
	wrapperContainer.AddChild(menuContainer)

	// タイトルテキスト、メニュー、バージョンテキストをrootContainerに追加
	rootContainer.AddChild(titleText)
	rootContainer.AddChild(wrapperContainer)
	rootContainer.AddChild(versionText)

	return &ebitenui.UI{Container: rootContainer}
}
