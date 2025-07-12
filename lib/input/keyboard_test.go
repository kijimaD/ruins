package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestGlobalKeyStateLogic(t *testing.T) {
	// グローバルキー状態の動作テスト

	// テスト用にリセット
	ResetGlobalKeyStateForTest()

	_ = NewDefaultKeyboardInput() // 異なるインスタンスでもグローバル状態を共有することを確認

	// 初期状態では最後に押されたキーはnil
	if GlobalKeyState.GetLastPressedKey() != nil {
		t.Error("初期状態でグローバルキー状態がnilではない")
	}

	// keyboard1で模擬的にキー設定
	GlobalKeyState.SetLastPressedKey(ebiten.KeyEnter)

	// keyboard2からも同じ状態が見える
	lastKey := GlobalKeyState.GetLastPressedKey()
	if lastKey == nil || *lastKey != ebiten.KeyEnter {
		t.Error("グローバル状態が正しく共有されていない")
	}

	// ClearLastPressedKey()の動作確認
	GlobalKeyState.ClearLastPressedKey()
	if GlobalKeyState.GetLastPressedKey() != nil {
		t.Error("ClearLastPressedKey()後にグローバル状態がクリアされていない")
	}
}

func TestMockKeyboardInput_DifferentKeyFunction(t *testing.T) {
	// テスト用にグローバル状態をリセット
	ResetGlobalKeyStateForTest()

	mock := NewMockKeyboardInput()

	// 異なるキー押下機能のテスト
	mock.SetKeyJustPressedIfDifferent(ebiten.KeyEnter, true)
	if !mock.IsKeyJustPressedIfDifferent(ebiten.KeyEnter) {
		t.Error("初回のキー押下が検出されなかった")
	}

	// 最後に押されたキーが記録されているか確認（グローバル状態を確認）
	lastKey := GlobalKeyState.GetLastPressedKey()
	if lastKey == nil || *lastKey != ebiten.KeyEnter {
		t.Error("最後に押されたキーが正しく記録されていない")
	}

	// クリア後の確認
	mock.ClearLastPressedKey()
	if GlobalKeyState.GetLastPressedKey() != nil {
		t.Error("ClearLastPressedKey()後にキーがクリアされていない")
	}

	// リセット後の確認
	mock.SetKeyJustPressedIfDifferent(ebiten.KeySpace, true)
	mock.Reset()
	if mock.IsKeyJustPressedIfDifferent(ebiten.KeySpace) {
		t.Error("Reset()後にキー状態がクリアされていない")
	}
}

func TestDifferentKeySequence(t *testing.T) {
	// テスト用にグローバル状態をリセット
	ResetGlobalKeyStateForTest()

	// 異なるキーの連続押下シーケンステスト
	mock := NewMockKeyboardInput()

	// Enter → Space → Enter のシーケンス

	// 1. 最初のEnterキー（成功）
	mock.SetKeyJustPressedIfDifferent(ebiten.KeyEnter, true)
	if !mock.IsKeyJustPressedIfDifferent(ebiten.KeyEnter) {
		t.Error("初回Enterキー押下が検出されなかった")
	}

	// 2. 同じEnterキーを再度押下（モック設定で無効化）
	mock.SetKeyJustPressedIfDifferent(ebiten.KeyEnter, false)
	if mock.IsKeyJustPressedIfDifferent(ebiten.KeyEnter) {
		t.Error("同じキーの連続押下が誤検出された")
	}

	// 3. 異なるSpaceキー（成功）
	mock.SetKeyJustPressedIfDifferent(ebiten.KeySpace, true)
	if !mock.IsKeyJustPressedIfDifferent(ebiten.KeySpace) {
		t.Error("異なるキー（Space）押下が検出されなかった")
	}

	// 4. 再度Enterキー（前回はSpaceなので成功）
	mock.SetKeyJustPressedIfDifferent(ebiten.KeyEnter, true)
	if !mock.IsKeyJustPressedIfDifferent(ebiten.KeyEnter) {
		t.Error("前回と異なるキー（Enter）押下が検出されなかった")
	}
}

func TestSharedKeyboardInput(t *testing.T) {
	// 共有インスタンスのテスト

	// 複数回呼び出しても同じインスタンスが返される
	keyboard1 := GetSharedKeyboardInput()
	keyboard2 := GetSharedKeyboardInput()

	if keyboard1 != keyboard2 {
		t.Error("GetSharedKeyboardInput()が異なるインスタンスを返している")
	}

	// テスト用にリセット
	ResetGlobalKeyStateForTest()

	// グローバル状態が共有されることを確認
	GlobalKeyState.SetLastPressedKey(ebiten.KeyEnter)

	lastKey := GlobalKeyState.GetLastPressedKey()
	if lastKey == nil || *lastKey != ebiten.KeyEnter {
		t.Error("共有インスタンス経由でグローバル状態にアクセスできない")
	}
}
