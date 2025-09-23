package messagewindow

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/maingame"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestWorld はテスト用のワールドを作成する
func createTestWorld(t *testing.T) w.World {
	t.Helper()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)
	return world
}

func TestMessageWindowBuilder(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("基本的なメッセージウィンドウの作成", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("テストメッセージです").
			Build()

		assert.NotNil(t, window)
		assert.True(t, window.IsOpen())
		assert.False(t, window.IsClosed())
		assert.Equal(t, "テストメッセージです", window.content.Text)
	})

	t.Run("メッセージタイプの設定", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("ストーリーメッセージ").
			Type(TypeStory).
			Build()

		assert.Equal(t, TypeStory, window.content.Type)
	})

	t.Run("話者名の設定", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("こんにちは！").
			Speaker("村人A").
			Build()

		assert.Equal(t, "村人A", window.content.SpeakerName)
	})

	t.Run("ウィンドウサイズとポジションの設定", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("カスタムサイズのウィンドウ").
			Size(800, 300).
			Position(100, 50).
			Build()

		assert.Equal(t, 800, window.config.Size.Width)
		assert.Equal(t, 300, window.config.Size.Height)
		assert.Equal(t, 100, window.config.Position.X)
		assert.Equal(t, 50, window.config.Position.Y)
		assert.False(t, window.config.Center)
	})

	t.Run("中央配置の設定", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("中央配置のウィンドウ").
			Center().
			Build()

		assert.True(t, window.config.Center)
	})

	t.Run("スキップ可能キーの設定", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("カスタムキー設定").
			SkippableKeys(ebiten.KeyZ, ebiten.KeyX).
			Build()

		assert.Equal(t, []ebiten.Key{ebiten.KeyZ, ebiten.KeyX}, window.config.SkippableKeys)
	})
}

func TestMessageWindowState(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("ウィンドウの開閉状態", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("状態テスト").
			Build()

		// 初期状態は開いている
		assert.True(t, window.IsOpen())
		assert.False(t, window.IsClosed())

		// 閉じる
		window.Close()
		assert.False(t, window.IsOpen())
		assert.True(t, window.IsClosed())
	})

	t.Run("閉じるコールバックの実行", func(t *testing.T) {
		t.Parallel()

		callbackCalled := false
		window := NewBuilder(world).
			Message("コールバックテスト").
			OnClose(func() {
				callbackCalled = true
			}).
			Build()

		window.Close()
		assert.True(t, callbackCalled)
	})
}

func TestMessageWindowUpdate(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("更新処理の基本動作", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("更新テスト").
			Build()

		// 更新処理を呼び出す（エラーが発生しないことを確認）
		require.NotPanics(t, func() {
			window.Update()
		})

		assert.True(t, window.initialized)
	})

	t.Run("閉じた状態での更新処理", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("閉じた状態のテスト").
			Build()

		window.Close()

		// 閉じた状態でも更新処理でエラーが発生しないことを確認
		require.NotPanics(t, func() {
			window.Update()
		})
	})
}

func TestMessageWindowFutureFeatures(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("選択肢システム", func(t *testing.T) {
		t.Parallel()

		actionCalled := false
		window := NewBuilder(world).
			Message("選択してください").
			Choice("選択肢1", func() {
				actionCalled = true
			}).
			Choice("選択肢2", func() {}).
			Build()

		assert.Len(t, window.content.Choices, 2)
		assert.Equal(t, "選択肢1", window.content.Choices[0].Text)
		assert.Equal(t, "選択肢2", window.content.Choices[1].Text)

		// 選択肢を実行
		window.selectChoice(0)
		assert.True(t, actionCalled)
		assert.False(t, window.IsOpen()) // 選択後はウィンドウが閉じる
	})

	t.Run("説明付き選択肢", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("詳細な選択肢").
			ChoiceWithDescription("攻撃", "敵に攻撃を仕掛ける", func() {}).
			Build()

		assert.Equal(t, "攻撃", window.content.Choices[0].Text)
		assert.Equal(t, "敵に攻撃を仕掛ける", window.content.Choices[0].Description)
	})
}

func TestMessageType(t *testing.T) {
	t.Parallel()

	t.Run("メッセージタイプの文字列表現", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "Story", TypeStory.String())
		assert.Equal(t, "Event", TypeEvent.String())
		assert.Equal(t, "Dialog", TypeDialog.String())
		assert.Equal(t, "System", TypeSystem.String())
	})
}

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("デフォルト設定の確認", func(t *testing.T) {
		t.Parallel()

		config := DefaultConfig()

		assert.Equal(t, 600, config.Size.Width)
		assert.Equal(t, 200, config.Size.Height)
		assert.True(t, config.Center)
		assert.True(t, config.ActionStyle.ShowCloseButton)
		assert.Contains(t, config.SkippableKeys, ebiten.KeyEnter)
		assert.Contains(t, config.SkippableKeys, ebiten.KeyEscape)
	})
}
