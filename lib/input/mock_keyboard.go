package input

import "github.com/hajimehoshi/ebiten/v2"

// MockKeyboardInput はテスト用のモックキーボード入力実装
type MockKeyboardInput struct {
	pressedKeys            map[ebiten.Key]bool
	justPressedKeys        map[ebiten.Key]bool
	justPressedIfDifferent map[ebiten.Key]bool
}

func NewMockKeyboardInput() *MockKeyboardInput {
	return &MockKeyboardInput{
		pressedKeys:            make(map[ebiten.Key]bool),
		justPressedKeys:        make(map[ebiten.Key]bool),
		justPressedIfDifferent: make(map[ebiten.Key]bool),
	}
}

func (m *MockKeyboardInput) IsKeyJustPressed(key ebiten.Key) bool {
	return m.justPressedKeys[key]
}

func (m *MockKeyboardInput) IsKeyPressed(key ebiten.Key) bool {
	return m.pressedKeys[key]
}

func (m *MockKeyboardInput) IsKeyJustPressedIfDifferent(key ebiten.Key) bool {
	return m.justPressedIfDifferent[key]
}

// SetKeyJustPressed はテスト用にキーの状態を設定する
func (m *MockKeyboardInput) SetKeyJustPressed(key ebiten.Key, pressed bool) {
	m.justPressedKeys[key] = pressed
}

// SetKeyPressed はテスト用にキーの状態を設定する
func (m *MockKeyboardInput) SetKeyPressed(key ebiten.Key, pressed bool) {
	m.pressedKeys[key] = pressed
}

// SetKeyJustPressedIfDifferent はテスト用に異なるキー検出の状態を設定
func (m *MockKeyboardInput) SetKeyJustPressedIfDifferent(key ebiten.Key, pressed bool) {
	m.justPressedIfDifferent[key] = pressed
}

// Reset は全てのキー状態をリセットする
func (m *MockKeyboardInput) Reset() {
	m.pressedKeys = make(map[ebiten.Key]bool)
	m.justPressedKeys = make(map[ebiten.Key]bool)
	m.justPressedIfDifferent = make(map[ebiten.Key]bool)
}
