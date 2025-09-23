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

func TestBuilderChaining(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("ビルダーチェーンの複雑な組み合わせ", func(t *testing.T) {
		t.Parallel()

		callbackExecuted := false
		window := NewBuilder(world).
			Message("複雑なメッセージテスト").
			Speaker("テストキャラクター").
			Size(800, 400).
			Position(100, 200).
			Choice("はい", func() { callbackExecuted = true }).
			Choice("いいえ", func() {}).
			ChoiceWithDescription("詳細", "詳細な説明", func() {}).
			OnClose(func() {}).
			Build()

		assert.NotNil(t, window)
		assert.Equal(t, "複雑なメッセージテスト", window.content.Text)
		assert.Equal(t, "テストキャラクター", window.content.SpeakerName)
		assert.Equal(t, 800, window.config.Size.Width)
		assert.Equal(t, 400, window.config.Size.Height)
		assert.Equal(t, 100, window.config.Position.X)
		assert.Equal(t, 200, window.config.Position.Y)
		assert.Len(t, window.content.Choices, 3)

		// 選択肢の実行テスト
		window.selectChoice(0)
		assert.True(t, callbackExecuted)
	})

	t.Run("空のメッセージでもビルド可能", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("").
			Build()

		assert.NotNil(t, window)
		assert.Equal(t, "", window.content.Text)
	})

	t.Run("デフォルト値の上書き", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("テスト").
			Size(400, 300).
			Size(500, 350). // 上書き
			Position(50, 100).
			Position(75, 125). // 上書き
			Build()

		assert.Equal(t, 500, window.config.Size.Width)
		assert.Equal(t, 350, window.config.Size.Height)
		assert.Equal(t, 75, window.config.Position.X)
		assert.Equal(t, 125, window.config.Position.Y)
	})
}

func TestBuilderValidation(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("無効なサイズでもビルド可能", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("テスト").
			Size(-100, -50). // 負の値
			Build()

		assert.NotNil(t, window)
		// 負の値でもビルダーはエラーにしない（内部で適切に処理される）
		assert.Equal(t, -100, window.config.Size.Width)
		assert.Equal(t, -50, window.config.Size.Height)
	})

	t.Run("重複した選択肢の追加", func(t *testing.T) {
		t.Parallel()

		action1Called := false
		action2Called := false

		window := NewBuilder(world).
			Message("重複テスト").
			Choice("同じテキスト", func() { action1Called = true }).
			Choice("同じテキスト", func() { action2Called = true }).
			Build()

		assert.Len(t, window.content.Choices, 2)
		assert.Equal(t, "同じテキスト", window.content.Choices[0].Text)
		assert.Equal(t, "同じテキスト", window.content.Choices[1].Text)

		// それぞれ独立したアクションが実行される
		window.selectChoice(0)
		assert.True(t, action1Called)
		assert.False(t, action2Called)
	})
}

func TestBuilderChoiceFeatures(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("選択肢なしのメッセージ", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("選択肢なし").
			Build()

		assert.Len(t, window.content.Choices, 0)
	})

	t.Run("選択肢のみ（空のアクション）", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("選択肢テスト").
			Choice("選択肢1", nil).
			Choice("選択肢2", func() {}).
			Build()

		assert.Len(t, window.content.Choices, 2)
		assert.NotPanics(t, func() {
			window.selectChoice(0) // nilアクションでもパニックしない
		})
	})

	t.Run("多数の選択肢", func(t *testing.T) {
		t.Parallel()

		builder := NewBuilder(world).Message("多数の選択肢")

		// 10個の選択肢を追加
		for i := 0; i < 10; i++ {
			builder = builder.Choice("選択肢"+string(rune('A'+i)), func() {})
		}

		window := builder.Build()
		assert.Len(t, window.content.Choices, 10)
	})

	t.Run("選択肢の説明文", func(t *testing.T) {
		t.Parallel()

		window := NewBuilder(world).
			Message("説明付き選択肢").
			ChoiceWithDescription("攻撃", "敵に物理攻撃を仕掛ける", func() {}).
			ChoiceWithDescription("魔法", "魔法による攻撃を行う", func() {}).
			ChoiceWithDescription("防御", "次のターンまでダメージを軽減する", func() {}).
			Build()

		assert.Len(t, window.content.Choices, 3)
		assert.Equal(t, "攻撃", window.content.Choices[0].Text)
		assert.Equal(t, "敵に物理攻撃を仕掛ける", window.content.Choices[0].Description)
		assert.Equal(t, "魔法", window.content.Choices[1].Text)
		assert.Equal(t, "魔法による攻撃を行う", window.content.Choices[1].Description)
	})
}

func TestBuilderCallbacks(t *testing.T) {
	t.Parallel()

	world := createTestWorld(t)

	t.Run("OnCloseコールバックの実行タイミング", func(t *testing.T) {
		t.Parallel()

		closeCallCount := 0
		window := NewBuilder(world).
			Message("コールバックテスト").
			OnClose(func() {
				closeCallCount++
			}).
			Build()

		assert.Equal(t, 0, closeCallCount)

		window.Close()
		assert.Equal(t, 1, closeCallCount)

		// 複数回閉じても1回だけ実行される
		window.Close()
		assert.Equal(t, 1, closeCallCount)
	})

	t.Run("複数のOnClose設定（最後が有効）", func(t *testing.T) {
		t.Parallel()

		callback1Called := false
		callback2Called := false

		window := NewBuilder(world).
			Message("複数コールバック").
			OnClose(func() { callback1Called = true }).
			OnClose(func() { callback2Called = true }). // 上書き
			Build()

		window.Close()
		assert.False(t, callback1Called) // 上書きされて実行されない
		assert.True(t, callback2Called)  // 最後に設定されたものが実行される
	})
}
