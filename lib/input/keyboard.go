package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// KeyboardInput はキーボード入力を抽象化するインターフェース
type KeyboardInput interface {
	IsKeyJustPressed(key ebiten.Key) bool
	IsKeyPressed(key ebiten.Key) bool
}

// DefaultKeyboardInput はEbitenのキーボード入力をラップする実装
type DefaultKeyboardInput struct{}

func NewDefaultKeyboardInput() KeyboardInput {
	return &DefaultKeyboardInput{}
}

func (d *DefaultKeyboardInput) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

func (d *DefaultKeyboardInput) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
