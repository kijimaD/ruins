package states

import (
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/typewriter"
	w "github.com/kijimaD/ruins/lib/world"
)

// MessageState はメッセージ表示用のステート
type MessageState struct {
	es.BaseState
	ui            *ebitenui.UI
	keyboardInput input.KeyboardInput

	text     string
	textFunc *func() string

	// タイプライター機能
	messageHandler *typewriter.MessageHandler

	// UIウィジェット参照（テキスト更新用）
	textWidget *widget.Text

	// 複数行表示管理。最大表示する行数
	maxVisibleLines int
}

func (st MessageState) String() string {
	return "Message"
}

// State interface ================

var _ es.State = &MessageState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MessageState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *MessageState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *MessageState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	// 複数行表示初期化
	st.maxVisibleLines = 4 // 最大4行表示

	// タイプライター初期化
	if st.messageHandler == nil {
		// MessageHandlerを初期化
		st.messageHandler = typewriter.NewMessageHandler(typewriter.BattleConfig(), st.keyboardInput)

		// フックを設定
		st.setupMessageHandlerHooks(world)

		if st.text != "" {
			st.messageHandler.Start(st.text)
		}
	}

	// 初回のUI作成
	if st.ui == nil {
		st.ui = st.createUI(world)
	}
}

// OnStop はステートが停止される際に呼ばれる
func (st *MessageState) OnStop(_ w.World) {}

// setupMessageHandlerHooks はMessageHandlerのフックを設定
func (st *MessageState) setupMessageHandlerHooks(world w.World) {
	// UI更新フック
	st.messageHandler.SetOnUpdateUI(func(_ string) {
		// タイプライター使用時はUIを再作成して表示を更新
		st.ui = st.createUIWithOffset(world)
	})

	// 完了フック - MessageHandlerからの戻り値でUpdate側で制御するため、ここでは追加処理のみ
	st.messageHandler.SetOnComplete(func() bool {
		// 完了時の追加処理があればここに記述
		return true
	})
}

// Update はメッセージステートの更新処理を行う
func (st *MessageState) Update(world w.World) es.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition{Type: es.TransQuit}
	}

	// タイプライター処理
	if st.messageHandler != nil {
		// MessageHandlerに処理を委譲し、完了状態をチェック
		shouldComplete := st.messageHandler.Update()
		if shouldComplete {
			return es.Transition{Type: es.TransPop}
		}
	}

	// textFunc による動的テキスト更新
	if st.textFunc != nil {
		f := *st.textFunc
		newText := f()
		st.textFunc = nil

		// 新しいテキストの場合
		if newText != st.text {
			st.text = newText

			if st.messageHandler != nil {
				st.messageHandler.Start(st.text)
			}
		}
	}

	// プロンプトアニメーション更新（UI再構築が必要）
	if st.messageHandler != nil && st.messageHandler.IsWaitingForInput() {
		// UI再構築（アニメーション位置更新のため）
		st.ui = st.createUIWithOffset(world)
	}

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *MessageState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *MessageState) createUI(world w.World) *ebitenui.UI {
	// 常にタイプライター用の複数行対応UIを使用
	return st.createUIWithOffset(world)
}

// createUIWithOffset は複数行表示対応のUIを作成
func (st *MessageState) createUIWithOffset(world w.World) *ebitenui.UI {
	// 固定サイズの計算（関数の最初で計算）
	lineHeight := 25 // 1行あたりの高さ（概算）
	fixedHeight := lineHeight * st.maxVisibleLines
	fixedWidth := 720 // メッセージコンテナの固定幅

	// 全体のコンテナ（水平レイアウト）
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    50,
				Left:   20,
				Right:  20,
				Bottom: 5,
			}),
			widget.RowLayoutOpts.Spacing(5), // 水平方向のスペース
		)),
	)

	res := world.Resources.UIResources

	// メッセージ表示用のコンテナ（固定サイズ）
	messageContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(fixedWidth, fixedHeight), // 固定サイズを直接指定
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(res.Panel.Image), // 背景を追加
	)

	// テキスト表示用の垂直コンテナ（メッセージコンテナ内に配置）
	textContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionStart,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	// タイプライター表示中のテキストを取得
	currentDisplayText := st.text
	if st.messageHandler != nil {
		currentDisplayText = st.messageHandler.GetDisplayText()
	}

	// 現在表示中のテキストを行に分割
	currentLines := strings.Split(currentDisplayText, "\n")

	// 表示する行を制限（最新のmaxVisibleLines行のみ）
	displayLines := currentLines
	if len(currentLines) > st.maxVisibleLines {
		// 古い行を削除して最新の行のみ表示
		displayLines = currentLines[len(currentLines)-st.maxVisibleLines:]
	}

	// 各行をテキストウィジェットとして追加（プロンプトは含めない）
	for i, lineText := range displayLines {
		textWidget := widget.NewText(
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				MaxWidth: fixedWidth - 20, // テキストの最大幅を制限（パディング分を考慮）
			})),
			widget.TextOpts.Text(lineText, res.Text.Face, styles.TextColor),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		)

		// 最初のテキストウィジェットを参照として保持
		if i == 0 {
			st.textWidget = textWidget
		}

		textContainer.AddChild(textWidget)
	}

	// テキストコンテナをメッセージコンテナに追加
	messageContainer.AddChild(textContainer)

	// メッセージコンテナをルートに追加
	rootContainer.AddChild(messageContainer)

	// プロンプト表示が必要な場合、右側に配置
	if st.messageHandler != nil && st.messageHandler.IsWaitingForInput() {
		// プロンプトコンテナを作成
		promptContainer := st.messageHandler.CreatePromptContainer(res.ComboButton.Graphic)
		if promptContainer != nil {
			// プロンプト用の右側コンテナ（メッセージコンテナと同じ高さ）
			promptSideContainer := widget.NewContainer(
				widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
				})),
				widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			)

			// プロンプトをコンテナの下端に配置（固定高さのボトムに合わせる）
			promptWrapper := widget.NewContainer(
				widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionStart,
					Padding: widget.Insets{
						Top: fixedHeight - 25, // メッセージコンテナの高さ - プロンプト高さ
					},
				})),
				widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			)

			promptWrapper.AddChild(promptContainer)
			promptSideContainer.AddChild(promptWrapper)
			rootContainer.AddChild(promptSideContainer)
		}
	}

	return &ebitenui.UI{Container: rootContainer}
}
