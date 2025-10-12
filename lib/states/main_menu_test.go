package states

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/testutil"
)

func TestMainMenuNavigation(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	state.keyboardInput = mockInput

	// メニューを初期化（Worldは簡易版を使用）
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// 初期状態の確認
	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("初期フォーカスが不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}

	// 下矢印キーでフォーカス移動
	mockInput.SetKeyJustPressed(ebiten.KeyArrowDown, true)
	state.menu.Update(mockInput)

	if state.menu.GetFocusedIndex() != 1 {
		t.Errorf("下矢印キー後のフォーカス位置が不正: 期待 1, 実際 %d", state.menu.GetFocusedIndex())
	}

	// 上矢印キーでフォーカス移動
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyArrowUp, true)
	state.menu.Update(mockInput)

	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("上矢印キー後のフォーカス位置が不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}
}

func TestMainMenuCircularNavigation(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	state.keyboardInput = mockInput

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// 最後の項目にフォーカス移動
	itemCount := len(state.menu.GetItems())
	state.menu.SetFocusedIndex(itemCount - 1)

	// 下矢印キーで循環して最初に戻る
	mockInput.SetKeyJustPressed(ebiten.KeyArrowDown, true)
	state.menu.Update(mockInput)

	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("循環移動後のフォーカス位置が不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}
}

func TestMainMenuSelection(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	state.keyboardInput = mockInput

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// 「終了」項目にフォーカス移動（インデックス2）
	state.menu.SetFocusedIndex(2)

	// Enterキーで選択（セッションベース）
	mockInput.SimulateEnterPressRelease()
	state.menu.Update(mockInput)

	// トランジションが設定されることを確認
	if state.GetTransition() == nil {
		t.Error("トランジションが設定されていない")
	} else if state.GetTransition().Type != es.TransQuit {
		t.Errorf("期待されるトランジション: TransQuit, 実際: %v", state.GetTransition().Type)
	}
}

func TestMainMenuCancel(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	state.keyboardInput = mockInput

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// Escapeキーでキャンセル
	mockInput.SetKeyJustPressed(ebiten.KeyEscape, true)
	state.menu.Update(mockInput)

	// トランジションが設定されることを確認
	if state.GetTransition() == nil {
		t.Error("トランジションが設定されていない")
	} else if state.GetTransition().Type != es.TransQuit {
		t.Errorf("期待されるトランジション: TransQuit, 実際: %v", state.GetTransition().Type)
	}
}

func TestMainMenuItems(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// メニュー項目の確認
	items := state.menu.GetItems()
	expectedItems := []string{"town", "load", "exit"}

	if len(items) != len(expectedItems) {
		t.Errorf("メニュー項目数が不正: 期待 %d, 実際 %d", len(expectedItems), len(items))
	}

	for i, expectedID := range expectedItems {
		if i < len(items) && items[i].ID != expectedID {
			t.Errorf("メニュー項目ID[%d]が不正: 期待 %s, 実際 %s", i, expectedID, items[i].ID)
		}
	}

	// ラベルの確認
	expectedLabels := []string{"開始", "読込", "終了"}
	for i, expectedLabel := range expectedLabels {
		if i < len(items) && items[i].Label != expectedLabel {
			t.Errorf("メニュー項目ラベル[%d]が不正: 期待 %s, 実際 %s", i, expectedLabel, items[i].Label)
		}
	}
}

func TestMainMenuTabNavigation(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	state.keyboardInput = mockInput

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// Tabキーでフォーカス移動
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	mockInput.SetKeyPressed(ebiten.KeyShift, false)
	state.menu.Update(mockInput)

	if state.menu.GetFocusedIndex() != 1 {
		t.Errorf("Tabキー後のフォーカス位置が不正: 期待 1, 実際 %d", state.menu.GetFocusedIndex())
	}

	// Shift+Tabキーでフォーカス移動
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	mockInput.SetKeyPressed(ebiten.KeyShift, true)
	state.menu.Update(mockInput)

	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("Shift+Tabキー後のフォーカス位置が不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}
}
