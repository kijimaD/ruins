package states

import (
	"testing"

	"github.com/ebitenui/ebitenui/widget"
)

func TestMainMenuFocusMovement(t *testing.T) {
	state := &MainMenuState{}

	// menuButtonsを手動で初期化
	state.menuButtons = make([]*widget.Button, 4) // 4つのメニュー項目
	state.focusIndex = 0

	initialFocus := state.focusIndex

	// 下に移動
	state.moveFocusDown()
	expectedIndex := (initialFocus + 1) % len(state.menuButtons)
	if state.focusIndex != expectedIndex {
		t.Errorf("下移動後のフォーカス位置が不正: 期待 %d, 実際 %d", expectedIndex, state.focusIndex)
	}

	// 上に移動
	state.moveFocusUp()
	if state.focusIndex != initialFocus {
		t.Errorf("上移動後のフォーカス位置が不正: 期待 %d, 実際 %d", initialFocus, state.focusIndex)
	}
}

func TestMainMenuCircularNavigation(t *testing.T) {
	state := &MainMenuState{}

	// menuButtonsを手動で初期化
	state.menuButtons = make([]*widget.Button, 4) // 4つのメニュー項目
	state.focusIndex = len(state.menuButtons) - 1 // 最後の要素

	// 下に移動（循環して最初に戻る）
	state.moveFocusDown()
	if state.focusIndex != 0 {
		t.Errorf("循環移動後のフォーカス位置が不正: 期待 0, 実際 %d", state.focusIndex)
	}

	// 上に移動（循環して最後に戻る）
	state.moveFocusUp()
	if state.focusIndex != len(state.menuButtons)-1 {
		t.Errorf("逆循環移動後のフォーカス位置が不正: 期待 %d, 実際 %d", 
			len(state.menuButtons)-1, state.focusIndex)
	}
}

func TestMainMenuSelectionExecution(t *testing.T) {
	state := &MainMenuState{}

	// menuButtonsを手動で初期化
	state.menuButtons = make([]*widget.Button, 4) // 4つのメニュー項目
	state.focusIndex = 1 // "拠点"
	
	// 選択実行
	state.executeCurrentSelection()

	// トランジションが設定されているかチェック
	if state.trans == nil {
		t.Error("選択実行後にトランジションが設定されていない")
	}

	if state.trans != nil && state.trans.Type != mainMenuTrans[1].trans.Type {
		t.Errorf("期待されるトランジションタイプ: %v, 実際: %v", 
			mainMenuTrans[1].trans.Type, state.trans.Type)
	}
}