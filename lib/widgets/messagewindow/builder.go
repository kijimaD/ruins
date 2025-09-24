package messagewindow

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	w "github.com/kijimaD/ruins/lib/world"
)

// Choice は選択肢を表す
type Choice struct {
	Text     string
	Action   func() // 選択時の処理
	Disabled bool   // 無効化フラグ
}

// MessageContent はメッセージの内容
type MessageContent struct {
	Text        string
	Choices     []Choice // 選択肢システム
	SpeakerName string   // 話者名（会話時）
}

// Builder はメッセージウィンドウを構築するためのビルダー
type Builder struct {
	config   Config
	content  MessageContent
	world    w.World
	onClose  func()
	onChoice func(choice Choice) // 選択肢コールバック
}

// NewBuilder は新しいBuilderを作成する
func NewBuilder(world w.World) *Builder {
	return &Builder{
		config: DefaultConfig(),
		world:  world,
	}
}

// Message はメッセージテキストを設定する
func (b *Builder) Message(text string) *Builder {
	b.content.Text = text
	return b
}

// Speaker は話者名を設定する（会話メッセージ用）
func (b *Builder) Speaker(name string) *Builder {
	b.content.SpeakerName = name
	return b
}

// Choice は選択肢を追加する
func (b *Builder) Choice(text string, action func()) *Builder {
	choice := Choice{
		Text:   text,
		Action: action,
	}
	b.content.Choices = append(b.content.Choices, choice)
	return b
}

// Size はウィンドウサイズを設定する
func (b *Builder) Size(width, height int) *Builder {
	b.config.Size = WindowSize{Width: width, Height: height}
	return b
}

// Position はウィンドウ位置を設定する
func (b *Builder) Position(x, y int) *Builder {
	b.config.Position = Position{X: x, Y: y}
	b.config.Center = false
	return b
}

// Center は画面中央配置を設定する
func (b *Builder) Center() *Builder {
	b.config.Center = true
	return b
}

// BackgroundColor は背景色を設定する
func (b *Builder) BackgroundColor(c color.Color) *Builder {
	b.config.WindowStyle.BackgroundColor = c
	return b
}

// TextColor はテキスト色を設定する
func (b *Builder) TextColor(c color.RGBA) *Builder {
	b.config.TextStyle.Color = c
	return b
}

// SkippableKeys はスキップ可能キーを設定する
func (b *Builder) SkippableKeys(keys ...ebiten.Key) *Builder {
	b.config.SkippableKeys = keys
	return b
}

// OnClose は閉じる時のコールバックを設定する
func (b *Builder) OnClose(callback func()) *Builder {
	b.onClose = callback
	return b
}

// OnChoice は選択肢選択時のコールバックを設定する
func (b *Builder) OnChoice(callback func(choice Choice)) *Builder {
	b.onChoice = callback
	return b
}

// Config はカスタム設定を適用する
func (b *Builder) Config(config Config) *Builder {
	b.config = config
	return b
}

// Build はメッセージウィンドウを構築する
func (b *Builder) Build() *Window {
	return &Window{
		config:   b.config,
		content:  b.content,
		world:    b.world,
		onClose:  b.onClose,
		onChoice: b.onChoice,
		isOpen:   true,
	}
}
