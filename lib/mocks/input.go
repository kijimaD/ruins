package mocks

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandler は入力ハンドラーのインターフェース
type InputHandler interface {
	IsKeyPressed(key ebiten.Key) bool
	IsKeyJustPressed(key ebiten.Key) bool
	GetWheel() (float64, float64)
	GetCursorPosition() (int, int)
}

// MockInputHandler はテスト用のモック入力ハンドラー
type MockInputHandler struct {
	pressedKeys      map[ebiten.Key]bool
	justPressedKeys  map[ebiten.Key]bool
	wheelX, wheelY   float64
	cursorX, cursorY int
}

// NewMockInputHandler は新しいモック入力ハンドラーを作成する
func NewMockInputHandler() *MockInputHandler {
	return &MockInputHandler{
		pressedKeys:     make(map[ebiten.Key]bool),
		justPressedKeys: make(map[ebiten.Key]bool),
	}
}

// IsKeyPressed はキーが押されているかを返す
func (m *MockInputHandler) IsKeyPressed(key ebiten.Key) bool {
	return m.pressedKeys[key]
}

// IsKeyJustPressed はキーが今フレームで押されたかを返す
func (m *MockInputHandler) IsKeyJustPressed(key ebiten.Key) bool {
	return m.justPressedKeys[key]
}

// GetWheel はマウスホイールの値を返す
func (m *MockInputHandler) GetWheel() (float64, float64) {
	return m.wheelX, m.wheelY
}

// GetCursorPosition はカーソル位置を返す
func (m *MockInputHandler) GetCursorPosition() (int, int) {
	return m.cursorX, m.cursorY
}

// テスト用のヘルパーメソッド

// PressKey はキーを押す
func (m *MockInputHandler) PressKey(key ebiten.Key) {
	m.pressedKeys[key] = true
	m.justPressedKeys[key] = true
}

// ReleaseKey はキーを離す
func (m *MockInputHandler) ReleaseKey(key ebiten.Key) {
	m.pressedKeys[key] = false
	m.justPressedKeys[key] = false
}

// ReleaseAllKeys は全てのキーを離す
func (m *MockInputHandler) ReleaseAllKeys() {
	m.pressedKeys = make(map[ebiten.Key]bool)
	m.justPressedKeys = make(map[ebiten.Key]bool)
}

// SetWheel はマウスホイールの値を設定する
func (m *MockInputHandler) SetWheel(x, y float64) {
	m.wheelX, m.wheelY = x, y
}

// SetCursorPosition はカーソル位置を設定する
func (m *MockInputHandler) SetCursorPosition(x, y int) {
	m.cursorX, m.cursorY = x, y
}

// EndFrame はフレーム終了時の処理（JustPressedをリセット）
func (m *MockInputHandler) EndFrame() {
	m.justPressedKeys = make(map[ebiten.Key]bool)
}
