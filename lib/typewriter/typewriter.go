package typewriter

import (
	"math"
	"time"
	"unicode/utf8"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Typewriter は文字送り表示を制御する汎用の構造
type Typewriter struct {
	config      Config
	currentText string // 表示したい文字列
	displayText string // 表示中の文字列
	position    int    // 現在のバイト位置（UTF-8対応）
	lastUpdate  time.Time
	state       State
	startTime   time.Time // 開始時刻

	// イベントハンドラー
	onChar     func(char string, index int)
	onComplete func()
	onSkip     func()
}

// State はタイプライターの状態を表す
type State int

const (
	// StateIdle はアイドル状態
	StateIdle State = iota
	// StateTyping はタイピング中状態
	StateTyping
	// StateComplete は完了状態
	StateComplete
	// StatePaused は一時停止状態
	StatePaused
)

// New は新しいTypewriterインスタンスを作成
func New(config Config) *Typewriter {
	return &Typewriter{
		config: config,
		state:  StateIdle,
	}
}

// Start はタイプライター表示を開始
func (t *Typewriter) Start(text string) {
	t.currentText = text
	t.displayText = ""
	t.position = 0
	t.lastUpdate = time.Now()
	t.startTime = time.Now()
	t.state = StateTyping
}

// Update は時間経過による状態更新を行う
func (t *Typewriter) Update() bool {
	if t.state != StateTyping {
		return false
	}

	now := time.Now()
	if now.Sub(t.lastUpdate) < t.config.CharDelay {
		return false
	}

	// 次の文字を表示
	if t.position < len(t.currentText) {
		// UTF-8対応の文字取得
		_, size := utf8.DecodeRuneInString(t.currentText[t.position:])
		charStr := t.currentText[t.position : t.position+size]

		t.displayText += charStr
		t.position += size
		t.lastUpdate = now

		// 文字表示イベント
		if t.onChar != nil {
			charIndex := utf8.RuneCountInString(t.displayText)
			t.onChar(charStr, charIndex)
		}

		return true
	}

	// 完了
	t.state = StateComplete
	if t.onComplete != nil {
		t.onComplete()
	}
	return false
}

// Skip は文字送りをスキップして完了状態にする
func (t *Typewriter) Skip() {
	if t.state == StateTyping && t.config.SkipEnabled {
		t.displayText = t.currentText
		t.position = len(t.currentText)
		t.state = StateComplete

		if t.onSkip != nil {
			t.onSkip()
		}
		if t.onComplete != nil {
			t.onComplete()
		}
	}
}

// Pause は文字送りを一時停止
func (t *Typewriter) Pause() {
	if t.state == StateTyping && t.config.PauseEnabled {
		t.state = StatePaused
	}
}

// Resume は文字送りを再開
func (t *Typewriter) Resume() {
	if t.state == StatePaused {
		t.state = StateTyping
		t.lastUpdate = time.Now()
	}
}

// GetDisplayText は現在表示中のテキストを取得
func (t *Typewriter) GetDisplayText() string {
	return t.displayText
}

// IsTyping は文字送り中かどうかを返す
func (t *Typewriter) IsTyping() bool {
	return t.state == StateTyping
}

// IsComplete は完了状態かどうかを返す
func (t *Typewriter) IsComplete() bool {
	return t.state == StateComplete
}

// IsPaused は一時停止中かどうかを返す
func (t *Typewriter) IsPaused() bool {
	return t.state == StatePaused
}

// GetProgress は進行状況（0.0-1.0）を返す
func (t *Typewriter) GetProgress() float64 {
	totalChars := utf8.RuneCountInString(t.currentText)
	if totalChars == 0 {
		return 0.0 // テキストがない場合は0.0を返す
	}
	charIndex := utf8.RuneCountInString(t.displayText)
	return float64(charIndex) / float64(totalChars)
}

// GetElapsedTime は開始からの経過時間を返す
func (t *Typewriter) GetElapsedTime() time.Duration {
	return time.Since(t.startTime)
}

// Reset は状態をリセット
func (t *Typewriter) Reset() {
	t.currentText = ""
	t.displayText = ""
	t.position = 0
	t.state = StateIdle
}

// OnChar は1文字表示時のコールバックを設定
func (t *Typewriter) OnChar(callback func(char string, index int)) {
	t.onChar = callback
}

// OnComplete は完了時のコールバックを設定
func (t *Typewriter) OnComplete(callback func()) {
	t.onComplete = callback
}

// OnSkip はスキップ時のコールバックを設定
func (t *Typewriter) OnSkip(callback func()) {
	t.onSkip = callback
}

// MessageHandler はメッセージ表示とタイプライターを統合管理する
// typewriterをラップしてコントロールできるようにする
type MessageHandler struct {
	typewriter *Typewriter

	// フック関数
	onUpdateUI func(text string)            // UI更新時のフック
	onComplete func() bool                  // 完了時のフック（戻り値で状態遷移を指示）
	onSkip     func()                       // スキップ時のフック
	onChar     func(char string, index int) // 文字表示時のフック

	// 入力処理用インターフェース
	keyboardInput KeyboardInput

	// プロンプトアニメーション用
	promptAnimationTime time.Time
	promptAmplitude     int
	promptWidget        *widget.Graphic
}

// KeyboardInput はキーボード入力を抽象化するインターフェース
type KeyboardInput interface {
	IsEnterJustPressedOnce() bool
}

// NewMessageHandler は新しいMessageHandlerを作成
func NewMessageHandler(config Config, keyboardInput KeyboardInput) *MessageHandler {
	handler := &MessageHandler{
		typewriter:          New(config),
		keyboardInput:       keyboardInput,
		promptAnimationTime: time.Now(),
		promptAmplitude:     5, // 上下移動の幅（ピクセル）
	}

	// タイプライター側のイベントハンドラーを設定
	handler.typewriter.OnChar(func(char string, index int) {
		if handler.onChar != nil {
			handler.onChar(char, index)
		}
		// UI更新も文字表示時に実行
		if handler.onUpdateUI != nil {
			handler.onUpdateUI(handler.typewriter.GetDisplayText())
		}
	})

	handler.typewriter.OnComplete(func() {
		// UI更新
		if handler.onUpdateUI != nil {
			handler.onUpdateUI(handler.typewriter.GetDisplayText())
		}
		// 完了時にプロンプトアニメーション開始
		handler.promptAnimationTime = time.Now()
	})

	handler.typewriter.OnSkip(func() {
		if handler.onSkip != nil {
			handler.onSkip()
		}
	})

	return handler
}

// Start はタイプライター表示を開始
func (h *MessageHandler) Start(text string) {
	h.typewriter.Start(text)
	// 新しいテキスト開始時にプロンプトウィジェットをリセット
	h.promptWidget = nil
	// 開始時にもUI更新
	if h.onUpdateUI != nil {
		h.onUpdateUI(h.typewriter.GetDisplayText())
	}
}

// Update はメッセージハンドラーの更新処理（入力処理も含む）
func (h *MessageHandler) Update() (shouldComplete bool) {
	// タイプライター更新
	h.typewriter.Update()

	// 入力処理（Enterキーまたはマウスクリック）
	enterPressed := h.keyboardInput != nil && h.keyboardInput.IsEnterJustPressedOnce()
	mouseClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	if enterPressed || mouseClicked {
		if h.typewriter.IsTyping() {
			// タイピング中なら文字送りスキップ
			h.typewriter.Skip()
			return false
		} else if h.typewriter.IsComplete() {
			// 完了していたら完了を通知
			if h.onComplete != nil {
				return h.onComplete()
			}
			return true
		}
	}

	return false
}

// IsTyping は文字送り中かどうかを返す
func (h *MessageHandler) IsTyping() bool {
	return h.typewriter.IsTyping()
}

// IsComplete は完了状態かどうかを返す
func (h *MessageHandler) IsComplete() bool {
	return h.typewriter.IsComplete()
}

// SetOnUpdateUI はUI更新時のフックを設定
func (h *MessageHandler) SetOnUpdateUI(callback func(text string)) {
	h.onUpdateUI = callback
}

// SetOnComplete は完了時のフックを設定
func (h *MessageHandler) SetOnComplete(callback func() bool) {
	h.onComplete = callback
}

// SetOnSkip はスキップ時のフックを設定
func (h *MessageHandler) SetOnSkip(callback func()) {
	h.onSkip = callback
}

// SetOnChar は文字表示時のフックを設定
func (h *MessageHandler) SetOnChar(callback func(char string, index int)) {
	h.onChar = callback
}

// IsWaitingForInput は入力待ち状態かどうかを返す
func (h *MessageHandler) IsWaitingForInput() bool {
	return h.typewriter.IsComplete()
}

// GetDisplayText は現在表示中のテキストを取得
func (h *MessageHandler) GetDisplayText() string {
	return h.typewriter.GetDisplayText()
}

// GetPromptWidget はプロンプトウィジェットを取得
func (h *MessageHandler) GetPromptWidget() *widget.Graphic {
	return h.promptWidget
}

// GetPromptYOffset はプロンプトアニメーションのY座標オフセットを取得
func (h *MessageHandler) GetPromptYOffset() int {
	if !h.IsWaitingForInput() {
		return 0
	}

	elapsedTime := time.Since(h.promptAnimationTime).Seconds()
	animationCycle := 2.0 // 2秒で1周期
	yOffset := int(math.Sin(elapsedTime*2*math.Pi/animationCycle) * float64(h.promptAmplitude))
	return yOffset
}

// CreatePromptContainer はプロンプト用のコンテナを作成（UIリソースが必要）
func (h *MessageHandler) CreatePromptContainer(arrowImage *widget.GraphicImage) *widget.Container {
	if !h.IsWaitingForInput() {
		return nil
	}

	yOffset := h.GetPromptYOffset()

	// promptWidgetが未作成の場合のみ作成
	if h.promptWidget == nil {
		h.promptWidget = widget.NewGraphic(
			widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
			widget.GraphicOpts.Image(arrowImage.Idle),
		)
	}

	// プロンプト用のコンテナを作成（アニメーション位置調整のため）
	promptContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionStart,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{
			Top: yOffset, // アニメーション分の上下移動
		}))),
	)

	promptContainer.AddChild(h.promptWidget)

	return promptContainer
}
