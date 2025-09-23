package states

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/messagewindow"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
)

// MessageWindowState はメッセージウィンドウを表示するステート
type MessageWindowState struct {
	es.BaseState[w.World]
	ui             *ebitenui.UI
	menu           *menu.Menu
	menuBuilder    *menu.UIBuilder
	keyboardInput  input.KeyboardInput
	messageWindow  *messagewindow.Window
	showingMessage bool
}

func (st MessageWindowState) String() string {
	return "MessageWindow"
}

// State interface ================

var _ es.State[w.World] = &MessageWindowState{}

// NewMessageWindowState は新しいMessageWindowStateを作成する
func NewMessageWindowState() es.State[w.World] {
	return &MessageWindowState{}
}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MessageWindowState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *MessageWindowState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *MessageWindowState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *MessageWindowState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *MessageWindowState) Update(_ w.World) es.Transition[w.World] {
	// メッセージウィンドウが表示中の場合
	if st.showingMessage && st.messageWindow != nil {
		st.messageWindow.Update()

		// メッセージウィンドウが閉じられた場合
		if st.messageWindow.IsClosed() {
			st.showingMessage = false
			st.messageWindow = nil
		}
	} else {
		// メニューの更新
		st.menu.Update(st.keyboardInput)
		st.ui.Update()
	}

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *MessageWindowState) Draw(_ w.World, screen *ebiten.Image) {
	// 背景メニューを描画
	st.ui.Draw(screen)

	// メッセージウィンドウが表示中の場合はその上に描画
	if st.showingMessage && st.messageWindow != nil {
		st.messageWindow.Draw(screen)
	}
}

// ================

func (st *MessageWindowState) initUI(world w.World) *ebitenui.UI {
	rootContainer := styled.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(consts.BlackColor)),
	)

	// Menuコンポーネントを作成
	st.createMenu(world)

	// MenuのUIを構築
	st.menuBuilder = menu.NewUIBuilder(world)
	menuContainer := st.menuBuilder.BuildUI(st.menu)
	rootContainer.AddChild(menuContainer)

	return &ebitenui.UI{Container: rootContainer}
}

// createMenu はメニューを作成する
func (st *MessageWindowState) createMenu(world w.World) {
	testItems := []struct {
		label string
		f     func(world w.World)
	}{
		{
			label: "基本メッセージ",
			f: func(world w.World) {
				st.showBasicMessage(world)
			},
		},
		{
			label: "ストーリーメッセージ",
			f: func(world w.World) {
				st.showStoryMessage(world)
			},
		},
		{
			label: "会話メッセージ",
			f: func(world w.World) {
				st.showDialogMessage(world)
			},
		},
		{
			label: "システムメッセージ",
			f: func(world w.World) {
				st.showSystemMessage(world)
			},
		},
		{
			label: "長いメッセージ",
			f: func(world w.World) {
				st.showLongMessage(world)
			},
		},
		{
			label: "選択肢テスト",
			f: func(world w.World) {
				st.showChoiceMessage(world)
			},
		},
		{
			label: "戻る",
			f: func(_ w.World) {
				st.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			},
		},
	}

	items := make([]menu.Item, len(testItems))
	for i, data := range testItems {
		items[i] = menu.Item{
			ID:       data.label,
			Label:    data.label,
			UserData: i,
		}
	}

	config := menu.Config{
		Items:             items,
		InitialIndex:      0,
		WrapNavigation:    true,
		Orientation:       menu.Vertical,
		ItemsPerPage:      10,
		ShowPageIndicator: true,
	}

	callbacks := menu.Callbacks{
		OnSelect: func(index int, _ menu.Item) {
			if index >= 0 && index < len(testItems) {
				testItems[index].f(world)
			}
		},
		OnCancel: func() {
			st.SetTransition(es.Transition[w.World]{Type: es.TransPop})
		},
		OnFocusChange: func(_, _ int) {
			if st.menuBuilder != nil {
				st.menuBuilder.UpdateFocus(st.menu)
			}
		},
	}

	st.menu = menu.NewMenu(config, callbacks)
}

// showBasicMessage は基本的なメッセージを表示する
func (st *MessageWindowState) showBasicMessage(world w.World) {
	st.messageWindow = messagewindow.NewBuilder(world).
		Message("これは基本的なメッセージです。\nEnter、Escape、Spaceキーで閉じることができます。").
		Build()
	st.showingMessage = true
}

// showStoryMessage はストーリーメッセージを表示する
func (st *MessageWindowState) showStoryMessage(world w.World) {
	st.messageWindow = messagewindow.NewBuilder(world).
		Message("遠い昔、魔法と剣の世界で...\n\n勇敢な冒険者が伝説の遺跡を発見した。").
		Type(messagewindow.TypeStory).
		Build()
	st.showingMessage = true
}

// showDialogMessage は会話メッセージを表示する
func (st *MessageWindowState) showDialogMessage(world w.World) {
	st.messageWindow = messagewindow.NewBuilder(world).
		Message("ようこそ、勇敢な冒険者よ！\nこの町で何かお手伝いできることはありますか？").
		Speaker("村長").
		Type(messagewindow.TypeDialog).
		Build()
	st.showingMessage = true
}

// showSystemMessage はシステムメッセージを表示する
func (st *MessageWindowState) showSystemMessage(world w.World) {
	st.messageWindow = messagewindow.NewBuilder(world).
		Message("ゲームが自動保存されました。\n\nシステム: セーブデータが正常に作成されました。").
		Type(messagewindow.TypeSystem).
		Build()
	st.showingMessage = true
}

// showLongMessage は長いメッセージを表示する
func (st *MessageWindowState) showLongMessage(world w.World) {
	longText := `これは非常に長いメッセージのテストです。

メッセージウィンドウは自動的にサイズを調整し、
長いテキストでも適切に表示されることを確認しています。

複数行のテキストと改行が正しく処理されること、
そしてウィンドウの背景やボーダーが適切に描画されることを
このテストで検証できます。

日本語のテキストも問題なく表示されるはずです。
句読点、記号、数字123なども含めて確認してみましょう。`

	st.messageWindow = messagewindow.NewBuilder(world).
		Message(longText).
		Size(700, 400).
		Build()
	st.showingMessage = true
}

// showChoiceMessage は選択肢機能のテストメッセージを表示する
func (st *MessageWindowState) showChoiceMessage(world w.World) {
	st.messageWindow = messagewindow.NewBuilder(world).
		Message("選択肢システムのテストです。\n\n矢印キーで移動、Enterで選択してください。").
		Choice("勇敢に戦う", func() {
			st.showChoiceResult(world, "勇敢に戦う", "あなたは勇気を振り絞って敵に立ち向かった！\n\n戦闘の結果、見事に勝利を収めました。\n経験値とアイテムを獲得しました！")
		}).
		Choice("慎重に交渉する", func() {
			st.showChoiceResult(world, "慎重に交渉する", "あなたは冷静に話し合いを試みた。\n\n相手も理解を示し、平和的に解決できました。\n友好関係が向上しました！")
		}).
		Choice("素早く逃走する", func() {
			st.showChoiceResult(world, "素早く逃走する", "あなたは迅速にその場を離れた。\n\n無事に危険を回避できましたが、\n何も得ることはできませんでした。")
		}).
		Build()
	st.showingMessage = true
}

// showChoiceResult は選択肢の結果メッセージを表示する
func (st *MessageWindowState) showChoiceResult(world w.World, choiceTitle, resultText string) {
	st.messageWindow = messagewindow.NewBuilder(world).
		Message(fmt.Sprintf("【%s】\n\n%s", choiceTitle, resultText)).
		Type(messagewindow.TypeEvent).
		OnClose(func() {
			// 結果メッセージの後にフォローアップメッセージを表示
			st.showFollowUpMessage(world)
		}).
		Build()
	st.showingMessage = true
}

// showFollowUpMessage はフォローアップメッセージを表示する
func (st *MessageWindowState) showFollowUpMessage(world w.World) {
	st.messageWindow = messagewindow.NewBuilder(world).
		Message("選択肢システムのテストが完了しました。\n\nEnterキーでメニューに戻ります。").
		Type(messagewindow.TypeSystem).
		Build()
	st.showingMessage = true
}
