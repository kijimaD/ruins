package input

import "github.com/hajimehoshi/ebiten/v2"

// MockKeyboardInput はテスト用のモックキーボード入力実装
type MockKeyboardInput struct {
	pressedKeys              map[ebiten.Key]bool
	justPressedKeys          map[ebiten.Key]bool
	justPressedIfDifferentKeys map[ebiten.Key]bool
	lastPressedKey           *ebiten.Key
}

func NewMockKeyboardInput() *MockKeyboardInput {
	return &MockKeyboardInput{
		pressedKeys:              make(map[ebiten.Key]bool),
		justPressedKeys:          make(map[ebiten.Key]bool),
		justPressedIfDifferentKeys: make(map[ebiten.Key]bool),
		lastPressedKey:           nil,
	}
}

func (m *MockKeyboardInput) IsKeyJustPressed(key ebiten.Key) bool {
	return m.justPressedKeys[key]
}

func (m *MockKeyboardInput) IsKeyPressed(key ebiten.Key) bool {
	return m.pressedKeys[key]
}

func (m *MockKeyboardInput) IsKeyJustPressedIfDifferent(key ebiten.Key) bool {
	// テスト用はローカル状態を使用
	return m.justPressedIfDifferentKeys[key]
}

func (m *MockKeyboardInput) ClearLastPressedKey() {
	// テスト用はローカル状態を使用
	m.lastPressedKey = nil
}

// SetKeyJustPressed はテスト用にキーの状態を設定する
func (m *MockKeyboardInput) SetKeyJustPressed(key ebiten.Key, pressed bool) {
	m.justPressedKeys[key] = pressed
}

// SetKeyPressed はテスト用にキーの状態を設定する
func (m *MockKeyboardInput) SetKeyPressed(key ebiten.Key, pressed bool) {
	m.pressedKeys[key] = pressed
}

// SetKeyJustPressedIfDifferent はテスト用に異なるキー押下の状態を設定する
func (m *MockKeyboardInput) SetKeyJustPressedIfDifferent(key ebiten.Key, pressed bool) {
	m.justPressedIfDifferentKeys[key] = pressed
	if pressed {
		m.lastPressedKey = &key
	}
}

// Reset は全てのキー状態をリセットする
func (m *MockKeyboardInput) Reset() {
	m.pressedKeys = make(map[ebiten.Key]bool)
	m.justPressedKeys = make(map[ebiten.Key]bool)
	m.justPressedIfDifferentKeys = make(map[ebiten.Key]bool)
	m.lastPressedKey = nil
}
