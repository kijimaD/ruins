package states

import (
	"testing"

	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/inputmapper"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/require"
)

func TestMainMenuNavigation(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// メニューを初期化（Worldは簡易版を使用）
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// 初期状態の確認
	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("初期フォーカスが不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}

	// ActionMenuDownでフォーカス移動
	require.NoError(t, state.menu.DoAction(inputmapper.ActionMenuDown))

	if state.menu.GetFocusedIndex() != 1 {
		t.Errorf("ActionMenuDown後のフォーカス位置が不正: 期待 1, 実際 %d", state.menu.GetFocusedIndex())
	}

	// ActionMenuUpでフォーカス移動
	require.NoError(t, state.menu.DoAction(inputmapper.ActionMenuUp))

	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("ActionMenuUp後のフォーカス位置が不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}
}

func TestMainMenuCircularNavigation(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// 最後の項目にフォーカス移動
	itemCount := len(state.menu.GetItems())
	state.menu.SetFocusedIndex(itemCount - 1)

	// ActionMenuDownで循環して最初に戻る
	require.NoError(t, state.menu.DoAction(inputmapper.ActionMenuDown))

	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("循環移動後のフォーカス位置が不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}
}

func TestMainMenuSelection(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// 「終了」項目にフォーカス移動（インデックス2）
	state.menu.SetFocusedIndex(2)

	// ActionMenuSelectで選択
	require.NoError(t, state.menu.DoAction(inputmapper.ActionMenuSelect))

	// トランジションが設定されることを確認
	transition := state.GetTransition()
	if transition == nil || transition.Type != es.TransQuit {
		if transition == nil {
			t.Error("トランジションが設定されていない")
		} else {
			t.Errorf("期待されるトランジション: TransQuit, 実際: %v", transition.Type)
		}
	}
}

func TestMainMenuCancel(t *testing.T) {
	t.Parallel()
	// MainMenuStateを作成
	state := &MainMenuState{}

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// ActionMenuCancelでキャンセル
	require.NoError(t, state.menu.DoAction(inputmapper.ActionMenuCancel))

	// トランジションが設定されることを確認
	transition := state.GetTransition()
	if transition == nil || transition.Type != es.TransQuit {
		if transition == nil {
			t.Error("トランジションが設定されていない")
		} else {
			t.Errorf("期待されるトランジション: TransQuit, 実際: %v", transition.Type)
		}
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

	// メニューを初期化
	world := testutil.InitTestWorld(t)
	state.initMenu(world)

	// ActionMenuDownでフォーカス移動（Tabキーと同じ動作）
	require.NoError(t, state.menu.DoAction(inputmapper.ActionMenuDown))

	if state.menu.GetFocusedIndex() != 1 {
		t.Errorf("ActionMenuDown後のフォーカス位置が不正: 期待 1, 実際 %d", state.menu.GetFocusedIndex())
	}

	// ActionMenuUpでフォーカス移動（Shift+Tabキーと同じ動作）
	require.NoError(t, state.menu.DoAction(inputmapper.ActionMenuUp))

	if state.menu.GetFocusedIndex() != 0 {
		t.Errorf("ActionMenuUp後のフォーカス位置が不正: 期待 0, 実際 %d", state.menu.GetFocusedIndex())
	}
}
