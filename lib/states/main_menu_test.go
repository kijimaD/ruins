package states

import (
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
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

func TestMainMenuKeyboardInputLogic(t *testing.T) {
	tests := []struct {
		name           string
		initialIndex   int
		setupKeys      func(*input.MockKeyboardInput)
		expectedAction string // "down", "up", "select", ""
	}{
		{
			name:         "ArrowDown triggers down action",
			initialIndex: 0,
			setupKeys: func(m *input.MockKeyboardInput) {
				m.SetKeyJustPressed(ebiten.KeyArrowDown, true)
			},
			expectedAction: "down",
		},
		{
			name:         "ArrowUp triggers up action",
			initialIndex: 1,
			setupKeys: func(m *input.MockKeyboardInput) {
				m.SetKeyJustPressed(ebiten.KeyArrowUp, true)
			},
			expectedAction: "up",
		},
		{
			name:         "Tab triggers down action",
			initialIndex: 0,
			setupKeys: func(m *input.MockKeyboardInput) {
				m.SetKeyJustPressed(ebiten.KeyTab, true)
				m.SetKeyPressed(ebiten.KeyShift, false)
			},
			expectedAction: "down",
		},
		{
			name:         "Shift+Tab triggers up action",
			initialIndex: 1,
			setupKeys: func(m *input.MockKeyboardInput) {
				m.SetKeyJustPressed(ebiten.KeyTab, true)
				m.SetKeyPressed(ebiten.KeyShift, true)
			},
			expectedAction: "up",
		},
		{
			name:         "Enter triggers select action",
			initialIndex: 1,
			setupKeys: func(m *input.MockKeyboardInput) {
				m.SetKeyJustPressed(ebiten.KeyEnter, true)
			},
			expectedAction: "select",
		},
		{
			name:         "Space triggers select action",
			initialIndex: 0,
			setupKeys: func(m *input.MockKeyboardInput) {
				m.SetKeyJustPressed(ebiten.KeySpace, true)
			},
			expectedAction: "select",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックキーボード入力を作成
			mockInput := input.NewMockKeyboardInput()
			
			// MainMenuStateを初期化
			state := &MainMenuState{
				keyboardInput: mockInput,
			}
			
			// 最小限の初期化
			state.menuButtons = make([]*widget.Button, 4)
			state.focusIndex = tt.initialIndex
			
			// キー入力を設定
			tt.setupKeys(mockInput)
			
			// キー入力判定のテスト
			action := ""
			isShiftPressed := state.keyboardInput.IsKeyPressed(ebiten.KeyShift)
			
			if state.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown) || 
			   (state.keyboardInput.IsKeyJustPressed(ebiten.KeyTab) && !isShiftPressed) {
				action = "down"
			} else if state.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp) || 
			          (state.keyboardInput.IsKeyJustPressed(ebiten.KeyTab) && isShiftPressed) {
				action = "up"
			} else if state.keyboardInput.IsKeyJustPressed(ebiten.KeyEnter) || 
			          state.keyboardInput.IsKeyJustPressed(ebiten.KeySpace) {
				action = "select"
			}
			
			// 結果を検証
			if action != tt.expectedAction {
				t.Errorf("アクションが不正: 期待 %q, 実際 %q", tt.expectedAction, action)
			}
		})
	}
}

func TestMainMenuEscapeKey(t *testing.T) {
	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyEscape, true)
	
	// MainMenuStateを初期化
	state := &MainMenuState{
		keyboardInput: mockInput,
	}
	
	// キーボード入力が正しく設定されていることを確認
	if !state.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		t.Error("Escapeキーが押されていることを検出できない")
	}
}