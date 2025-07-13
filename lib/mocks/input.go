package mocks

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// 入力ハンドラーのインターフェース
type InputHandler interface {
	IsKeyPressed(key ebiten.Key) bool
	IsKeyJustPressed(key ebiten.Key) bool
	GetWheel() (float64, float64)
	GetCursorPosition() (int, int)
}

// テスト用のモック入力ハンドラー
type MockInputHandler struct {
	pressedKeys      map[ebiten.Key]bool
	justPressedKeys  map[ebiten.Key]bool
	wheelX, wheelY   float64
	cursorX, cursorY int
}

// 新しいモック入力ハンドラーを作成する
func NewMockInputHandler() *MockInputHandler {
	return &MockInputHandler{
		pressedKeys:     make(map[ebiten.Key]bool),
		justPressedKeys: make(map[ebiten.Key]bool),
	}
}

// キーが押されているかを返す
func (m *MockInputHandler) IsKeyPressed(key ebiten.Key) bool {
	return m.pressedKeys[key]
}

// キーが今フレームで押されたかを返す
func (m *MockInputHandler) IsKeyJustPressed(key ebiten.Key) bool {
	return m.justPressedKeys[key]
}

// マウスホイールの値を返す
func (m *MockInputHandler) GetWheel() (float64, float64) {
	return m.wheelX, m.wheelY
}

// カーソル位置を返す
func (m *MockInputHandler) GetCursorPosition() (int, int) {
	return m.cursorX, m.cursorY
}

// テスト用のヘルパーメソッド

// キーを押す
func (m *MockInputHandler) PressKey(key ebiten.Key) {
	m.pressedKeys[key] = true
	m.justPressedKeys[key] = true
}

// キーを離す
func (m *MockInputHandler) ReleaseKey(key ebiten.Key) {
	m.pressedKeys[key] = false
	m.justPressedKeys[key] = false
}

// 全てのキーを離す
func (m *MockInputHandler) ReleaseAllKeys() {
	m.pressedKeys = make(map[ebiten.Key]bool)
	m.justPressedKeys = make(map[ebiten.Key]bool)
}

// マウスホイールの値を設定する
func (m *MockInputHandler) SetWheel(x, y float64) {
	m.wheelX, m.wheelY = x, y
}

// カーソル位置を設定する
func (m *MockInputHandler) SetCursorPosition(x, y int) {
	m.cursorX, m.cursorY = x, y
}

// フレーム終了時の処理（JustPressedをリセット）
func (m *MockInputHandler) EndFrame() {
	m.justPressedKeys = make(map[ebiten.Key]bool)
}
