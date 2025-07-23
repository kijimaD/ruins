package typewriter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypewriterBasicFlow(t *testing.T) {
	t.Parallel()

	config := Config{
		CharDelay:   50 * time.Millisecond,
		SkipEnabled: true,
	}
	tw := New(config)

	completed := false
	tw.OnComplete(func() {
		completed = true
	})

	tw.Start("Hello")

	// 初期状態確認
	assert.True(t, tw.IsTyping())
	assert.False(t, tw.IsComplete())
	assert.Equal(t, "", tw.GetDisplayText())
	assert.Equal(t, 0.0, tw.GetProgress())

	// 文字送りシミュレーション
	for i := 0; i < 20 && !completed; i++ {
		tw.Update()
		time.Sleep(60 * time.Millisecond)
	}

	assert.True(t, completed)
	assert.False(t, tw.IsTyping())
	assert.True(t, tw.IsComplete())
	assert.Equal(t, "Hello", tw.GetDisplayText())
	assert.Equal(t, 1.0, tw.GetProgress())
}

func TestTypewriterSkip(t *testing.T) {
	t.Parallel()

	tw := New(BattleConfig())

	skipped := false
	completed := false

	tw.OnSkip(func() {
		skipped = true
	})

	tw.OnComplete(func() {
		completed = true
	})

	tw.Start("Long message")

	// 1文字表示後にスキップ
	tw.Update()
	assert.True(t, tw.IsTyping())

	tw.Skip()

	assert.True(t, skipped)
	assert.True(t, completed)
	assert.True(t, tw.IsComplete())
	assert.False(t, tw.IsTyping())
	assert.Equal(t, "Long message", tw.GetDisplayText())
	assert.Equal(t, 1.0, tw.GetProgress())
}

func TestTypewriterProgress(t *testing.T) {
	t.Parallel()

	tw := New(Config{CharDelay: 10 * time.Millisecond})
	tw.Start("ABC")

	assert.Equal(t, 0.0, tw.GetProgress())

	// 1文字表示
	time.Sleep(15 * time.Millisecond)
	tw.Update()
	assert.InDelta(t, 0.33, tw.GetProgress(), 0.1)

	// 2文字表示
	time.Sleep(15 * time.Millisecond)
	tw.Update()
	assert.InDelta(t, 0.66, tw.GetProgress(), 0.1)

	// 3文字表示（完了）
	time.Sleep(15 * time.Millisecond)
	updated := tw.Update()
	// 最後の文字が表示された後、次のUpdateで完了状態になる
	if updated {
		time.Sleep(15 * time.Millisecond)
		tw.Update()
	}
	assert.Equal(t, 1.0, tw.GetProgress())
	assert.True(t, tw.IsComplete())
}

func TestTypewriterUTF8Support(t *testing.T) {
	t.Parallel()

	tw := New(Config{CharDelay: 10 * time.Millisecond})
	tw.Start("こんにちは")

	// 日本語5文字の進行確認
	for i := 0; i < 10; i++ {
		time.Sleep(15 * time.Millisecond)
		if !tw.Update() {
			break
		}
	}

	assert.True(t, tw.IsComplete())
	assert.Equal(t, "こんにちは", tw.GetDisplayText())
	assert.Equal(t, 1.0, tw.GetProgress())
}

func TestTypewriterPauseResume(t *testing.T) {
	t.Parallel()

	config := Config{
		CharDelay:    50 * time.Millisecond,
		PauseEnabled: true,
	}
	tw := New(config)

	tw.Start("Test")
	tw.Update() // 1文字表示

	// 一時停止
	tw.Pause()
	assert.True(t, tw.IsPaused())
	assert.False(t, tw.IsTyping())

	// 一時停止中は進行しない
	oldDisplay := tw.GetDisplayText()
	time.Sleep(60 * time.Millisecond)
	tw.Update()
	assert.Equal(t, oldDisplay, tw.GetDisplayText())

	// 再開
	tw.Resume()
	assert.False(t, tw.IsPaused())
	assert.True(t, tw.IsTyping())
}

func TestTypewriterReset(t *testing.T) {
	t.Parallel()

	tw := New(BattleConfig())
	tw.Start("Test message")

	// 少し待ってからUpdate（文字が表示されるまで）
	time.Sleep(60 * time.Millisecond)
	tw.Update()

	// リセット前の状態確認
	assert.True(t, len(tw.GetDisplayText()) > 0)
	assert.True(t, tw.IsTyping())

	// リセット
	tw.Reset()

	// リセット後の状態確認
	assert.Equal(t, "", tw.GetDisplayText())
	assert.False(t, tw.IsTyping())
	assert.False(t, tw.IsComplete())
	assert.Equal(t, 0.0, tw.GetProgress())
	// StateIdleかつ文字数が0なので、progressは正常に0になる
}

func TestTypewriterSpecialCharDelay(t *testing.T) {
	t.Parallel()

	tw := New(Config{
		CharDelay:        10 * time.Millisecond,
		PunctuationDelay: 100 * time.Millisecond,
		CommaDelay:       50 * time.Millisecond,
	})

	// 句読点での待機時間テスト
	tw.Start("こんにちは。")

	// 通常文字の表示
	for i := 0; i < 5; i++ {
		time.Sleep(15 * time.Millisecond)
		tw.Update()
	}

	// 句読点前の状態
	assert.Equal(t, "こんにちは", tw.GetDisplayText())

	// 句読点表示（より長い待機時間が必要）
	time.Sleep(15 * time.Millisecond)
	updated := tw.Update()

	// 句読点表示後は完了状態になる
	if updated {
		// 句読点の追加遅延後に再度更新を試す
		time.Sleep(110 * time.Millisecond)
		tw.Update()
	}

	assert.Equal(t, "こんにちは。", tw.GetDisplayText())
	assert.True(t, tw.IsComplete())
}

func TestConfigPresets(t *testing.T) {
	t.Parallel()

	// 各設定プリセットが正常に作成できることを確認
	configs := map[string]Config{
		"Default": DefaultConfig(),
		"Fast":    FastConfig(),
		"Slow":    SlowConfig(),
		"Battle":  BattleConfig(),
		"Dialog":  DialogConfig(),
	}

	for name, config := range configs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Greater(t, config.CharDelay, time.Duration(0))
			assert.GreaterOrEqual(t, config.PunctuationDelay, time.Duration(0))
			assert.GreaterOrEqual(t, config.CommaDelay, time.Duration(0))
			assert.GreaterOrEqual(t, config.NewlineDelay, time.Duration(0))

			// プリセットでタイプライターが作成できることを確認
			tw := New(config)
			assert.NotNil(t, tw)
		})
	}
}

func TestTypewriterCallbacks(t *testing.T) {
	t.Parallel()

	tw := New(Config{CharDelay: 10 * time.Millisecond})

	charCount := 0
	var lastChar string
	var lastIndex int

	tw.OnChar(func(char string, index int) {
		charCount++
		lastChar = char
		lastIndex = index
	})

	completed := false
	tw.OnComplete(func() {
		completed = true
	})

	tw.Start("AB")

	// 1文字目
	time.Sleep(15 * time.Millisecond)
	tw.Update()
	assert.Equal(t, 1, charCount)
	assert.Equal(t, "A", lastChar)
	assert.Equal(t, 1, lastIndex)

	// 2文字目
	time.Sleep(15 * time.Millisecond)
	updated := tw.Update()
	assert.Equal(t, 2, charCount)
	assert.Equal(t, "B", lastChar)
	assert.Equal(t, 2, lastIndex)

	// 最後の文字が表示された後、次のUpdateで完了状態になる
	if updated {
		time.Sleep(15 * time.Millisecond)
		tw.Update()
	}
	assert.True(t, completed)
}

// MockKeyboardInput はテスト用のキーボード入力モック
type MockKeyboardInput struct {
	enterPressed bool
}

func (m *MockKeyboardInput) IsEnterJustPressedOnce() bool {
	pressed := m.enterPressed
	m.enterPressed = false // 一度だけ押されたことにする
	return pressed
}

func (m *MockKeyboardInput) SetEnterPressed(pressed bool) {
	m.enterPressed = pressed
}

func TestMessageHandlerBasicFlow(t *testing.T) {
	t.Parallel()

	mockInput := &MockKeyboardInput{}
	handler := NewMessageHandler(BattleConfig(), mockInput)

	uiUpdated := false
	completed := false

	handler.SetOnUpdateUI(func(text string) {
		uiUpdated = true
	})

	handler.SetOnComplete(func() bool {
		completed = true
		return true
	})

	handler.Start("Hello")

	// 初期状態確認
	assert.True(t, handler.IsTyping())
	assert.False(t, handler.IsComplete())
	assert.Equal(t, "", handler.GetDisplayText())

	// 1文字表示
	time.Sleep(60 * time.Millisecond)
	shouldComplete := handler.Update()
	assert.False(t, shouldComplete)
	assert.True(t, uiUpdated)
	assert.True(t, len(handler.GetDisplayText()) > 0)

	// Enterキーでスキップ
	mockInput.SetEnterPressed(true)
	shouldComplete = handler.Update()
	assert.False(t, shouldComplete) // スキップなので完了は次回
	assert.True(t, handler.IsComplete())
	assert.Equal(t, "Hello", handler.GetDisplayText())

	// 完了状態でEnterキーを押す
	mockInput.SetEnterPressed(true)
	shouldComplete = handler.Update()
	assert.True(t, shouldComplete)
	assert.True(t, completed)
}

func TestMessageHandlerUpdateFlow(t *testing.T) {
	t.Parallel()

	mockInput := &MockKeyboardInput{}
	handler := NewMessageHandler(Config{CharDelay: 30 * time.Millisecond}, mockInput)

	charCallbackCount := 0
	handler.SetOnChar(func(char string, index int) {
		charCallbackCount++
	})

	handler.Start("Test")

	// 通常の更新サイクル
	for i := 0; i < 10 && !handler.IsComplete(); i++ {
		time.Sleep(40 * time.Millisecond)
		handler.Update()
	}

	assert.True(t, handler.IsComplete())
	assert.Equal(t, "Test", handler.GetDisplayText())
	assert.Equal(t, 4, charCallbackCount) // "Test" = 4文字
}

func TestMessageHandlerSkipFlow(t *testing.T) {
	t.Parallel()

	mockInput := &MockKeyboardInput{}
	handler := NewMessageHandler(BattleConfig(), mockInput)

	skipCalled := false
	handler.SetOnSkip(func() {
		skipCalled = true
	})

	handler.Start("Long message for skip test")

	// 文字送り開始
	time.Sleep(60 * time.Millisecond)
	handler.Update()
	assert.True(t, handler.IsTyping())

	// Enterキーでスキップ
	mockInput.SetEnterPressed(true)
	handler.Update()

	assert.True(t, skipCalled)
	assert.True(t, handler.IsComplete())
	assert.Equal(t, "Long message for skip test", handler.GetDisplayText())
}

func TestMessageHandlerEnterWait(t *testing.T) {
	t.Parallel()

	mockInput := &MockKeyboardInput{}
	handler := NewMessageHandler(Config{CharDelay: 10 * time.Millisecond}, mockInput)

	completed := false
	handler.SetOnComplete(func() bool {
		completed = true
		return true // 完了を示す
	})

	handler.Start("Hello")

	// 文字送り完了まで待つ
	for i := 0; i < 20 && !handler.IsComplete(); i++ {
		time.Sleep(15 * time.Millisecond)
		shouldComplete := handler.Update()
		assert.False(t, shouldComplete) // まだ完了しない
	}

	// 完了状態になっているはず
	assert.True(t, handler.IsComplete())
	assert.Equal(t, "Hello", handler.GetDisplayText())
	assert.False(t, completed) // まだコールバックは呼ばれていない

	// Enterキーを押すと完了コールバックが呼ばれ、shouldCompleteがtrueになる
	mockInput.SetEnterPressed(true)
	shouldComplete := handler.Update()

	assert.True(t, shouldComplete) // 完了を通知
	assert.True(t, completed)      // コールバックが呼ばれた
}

func TestMessageHandlerTypingSkip(t *testing.T) {
	t.Parallel()

	mockInput := &MockKeyboardInput{}
	handler := NewMessageHandler(Config{
		CharDelay:   50 * time.Millisecond,
		SkipEnabled: true,
	}, mockInput) // 適度に遅い設定

	handler.Start("This is a long message for skip testing")

	// 文字が表示されるまで少し待つ
	time.Sleep(60 * time.Millisecond)
	handler.Update()
	assert.True(t, handler.IsTyping())

	// 少なくとも1文字は表示されているはず
	displayText := handler.GetDisplayText()
	assert.True(t, len(displayText) > 0)
	assert.True(t, len(displayText) < len("This is a long message for skip testing"))

	// タイピング中にEnterキーを押すとスキップ
	mockInput.SetEnterPressed(true)
	shouldComplete := handler.Update()

	assert.False(t, shouldComplete)      // スキップ時は即座には完了しない
	assert.True(t, handler.IsComplete()) // ただし状態は完了になる
	assert.Equal(t, "This is a long message for skip testing", handler.GetDisplayText())

	// 完了状態でもう一度Enterキーを押すと完了
	mockInput.SetEnterPressed(true)
	shouldComplete = handler.Update()
	assert.True(t, shouldComplete) // 今度は完了を通知
}
