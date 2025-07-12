package input

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GlobalKeyState はグローバルなキー状態を管理する
var GlobalKeyState = &globalKeyState{
	lastPressedKey: nil,
}

// globalKeyState はアプリケーション全体で共有されるキー状態
type globalKeyState struct {
	lastPressedKey *ebiten.Key
	mutex          sync.RWMutex
}

// SetLastPressedKey は最後に押されたキーを設定する
func (g *globalKeyState) SetLastPressedKey(key ebiten.Key) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.lastPressedKey = &key
}

// GetLastPressedKey は最後に押されたキーを取得する
func (g *globalKeyState) GetLastPressedKey() *ebiten.Key {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.lastPressedKey
}

// ClearLastPressedKey は最後に押されたキーをクリアする
func (g *globalKeyState) ClearLastPressedKey() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.lastPressedKey = nil
}

// ResetForTest はテスト用にグローバル状態をリセットする
func ResetGlobalKeyStateForTest() {
	GlobalKeyState.ClearLastPressedKey()
}

// KeyboardInput はキーボード入力を抽象化するインターフェース
type KeyboardInput interface {
	IsKeyJustPressed(key ebiten.Key) bool
	IsKeyPressed(key ebiten.Key) bool
	IsKeyJustPressedIfDifferent(key ebiten.Key) bool // 前回と異なるキーの場合のみtrue
	ClearLastPressedKey()                            // 最後に押されたキーをクリア
}

// DefaultKeyboardInput はEbitenのキーボード入力をラップする実装
type DefaultKeyboardInput struct {
	// グローバル状態を使用するため、ローカルフィールドは不要
}

var (
	// sharedKeyboardInput は共有されるキーボード入力インスタンス
	sharedKeyboardInput KeyboardInput
	once                sync.Once
)

// GetSharedKeyboardInput は共有されるキーボード入力インスタンスを返す
func GetSharedKeyboardInput() KeyboardInput {
	once.Do(func() {
		sharedKeyboardInput = &DefaultKeyboardInput{}
	})
	return sharedKeyboardInput
}

func NewDefaultKeyboardInput() KeyboardInput {
	return &DefaultKeyboardInput{}
}

func (d *DefaultKeyboardInput) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

func (d *DefaultKeyboardInput) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// IsKeyJustPressedIfDifferent は前回と異なるキーが押された場合のみtrueを返す
func (d *DefaultKeyboardInput) IsKeyJustPressedIfDifferent(key ebiten.Key) bool {
	if !inpututil.IsKeyJustPressed(key) {
		return false
	}
	
	// グローバル状態から前回のキーを取得
	lastKey := GlobalKeyState.GetLastPressedKey()
	
	// 前回と同じキーの場合は無効
	if lastKey != nil && *lastKey == key {
		return false
	}
	
	// 前回と異なるキー（または初回）の場合は有効
	GlobalKeyState.SetLastPressedKey(key)
	return true
}

// ClearLastPressedKey は最後に押されたキーの記録をクリアする
func (d *DefaultKeyboardInput) ClearLastPressedKey() {
	GlobalKeyState.ClearLastPressedKey()
}
