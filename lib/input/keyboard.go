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

// InitializeWithDummyKey はアプリケーション開始時にキー状態を初期化する
func InitializeWithDummyKey() {
	GlobalKeyState.mutex.Lock()
	defer GlobalKeyState.mutex.Unlock()

	// 初回のEnterキー重複実行を防ぐため、Enterを既に押された状態にする
	enterKey := ebiten.KeyEnter
	GlobalKeyState.lastPressedKey = &enterKey
}

// KeyboardInput はキーボード入力を抽象化するインターフェース
type KeyboardInput interface {
	IsKeyJustPressed(key ebiten.Key) bool
	IsKeyPressed(key ebiten.Key) bool
	IsKeyJustPressedIfDifferent(key ebiten.Key) bool // 前回と異なるキーの場合のみtrue
}

// sharedKeyboardInput はシングルトンのキーボード入力実装
type sharedKeyboardInput struct{}

var (
	// keyboardInstance は共有されるキーボード入力インスタンス
	keyboardInstance KeyboardInput
	once             sync.Once
)

// GetSharedKeyboardInput は共有されるキーボード入力インスタンスを返す
func GetSharedKeyboardInput() KeyboardInput {
	once.Do(func() {
		keyboardInstance = &sharedKeyboardInput{}
	})
	return keyboardInstance
}

func (s *sharedKeyboardInput) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

func (s *sharedKeyboardInput) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// IsKeyJustPressedIfDifferent は前回と異なるキーが押された場合のみtrueを返す
func (s *sharedKeyboardInput) IsKeyJustPressedIfDifferent(key ebiten.Key) bool {
	if !inpututil.IsKeyJustPressed(key) {
		return false
	}

	// グローバル状態から前回のキーを取得
	lastKey := GlobalKeyState.GetLastPressedKey()

	// 前回と同じキーの場合は無効
	if lastKey != nil && *lastKey == key {
		return false
	}

	// 前回と異なるキーの場合は有効
	GlobalKeyState.SetLastPressedKey(key)
	return true
}
