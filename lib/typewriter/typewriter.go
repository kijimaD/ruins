package typewriter

import (
	"time"
	"unicode/utf8"
)

// Typewriter は文字送り表示を制御する汎用コンポーネント
type Typewriter struct {
	config      Config
	currentText string
	displayText string
	position    int // 現在の文字位置（byte単位）
	charIndex   int // 現在の文字インデックス（文字単位）
	totalChars  int // 総文字数
	lastUpdate  time.Time
	state       State
	startTime   time.Time // 開始時刻

	// イベントハンドラー
	onChar     func(char string, index int)
	onComplete func()
	onSkip     func()
}

type State int

const (
	StateIdle State = iota
	StateTyping
	StateComplete
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
	t.charIndex = 0
	t.totalChars = utf8.RuneCountInString(text)
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
		char, size := utf8.DecodeRuneInString(t.currentText[t.position:])
		charStr := t.currentText[t.position : t.position+size]

		t.displayText += charStr
		t.position += size
		t.charIndex++
		t.lastUpdate = now

		// 文字表示イベント
		if t.onChar != nil {
			t.onChar(charStr, t.charIndex)
		}

		// 特殊文字による待機時間調整
		if delay := t.getSpecialCharDelay(char); delay > 0 {
			t.lastUpdate = t.lastUpdate.Add(delay)
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
		t.charIndex = t.totalChars
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
	if t.totalChars == 0 {
		return 0.0  // テキストがない場合は0.0を返す
	}
	return float64(t.charIndex) / float64(t.totalChars)
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
	t.charIndex = 0
	t.totalChars = 0
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

// getSpecialCharDelay は特殊文字による待機時間調整
func (t *Typewriter) getSpecialCharDelay(char rune) time.Duration {
	switch char {
	case '。', '！', '？':
		return t.config.PunctuationDelay
	case '、', '，':
		return t.config.CommaDelay
	case '\n':
		return t.config.NewlineDelay
	default:
		return 0
	}
}
