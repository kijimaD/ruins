package hud

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/widgets/messagelog"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
)

const (
	defaultLogAreaHeight = 120 // ログエリアの高さ（余裕を持たせて大きめに）
	defaultMaxLogLines   = 5   // 表示する最大行数
	defaultLogAreaMargin = 8   // 余白
	defaultLineHeight    = 20  // 1行の高さ
	defaultYPadding      = 8   // 下端の追加パディング
)

// MessageArea はHUDメッセージエリア
type MessageArea struct {
	widget  *messagelog.Widget
	enabled bool
}

// NewMessageArea は新しいHUDMessageAreaを作成する
func NewMessageArea(world w.World) *MessageArea {
	// MessageLogWidget設定
	widgetConfig := messagelog.WidgetConfig{
		MaxLines:   defaultMaxLogLines,
		LineHeight: defaultLineHeight,
		Spacing:    3,
		Padding: messagelog.Insets{
			Top:    2,
			Bottom: 2,
			Left:   2,
			Right:  2,
		},
	}

	// MessageLogWidgetを作成
	widget := messagelog.NewWidget(widgetConfig, world)

	// デフォルトでFieldLogを使用
	widget.SetStore(gamelog.FieldLog)

	return &MessageArea{
		widget:  widget,
		enabled: true,
	}
}

// SetStore はログストアを設定する
func (area *MessageArea) SetStore(store *gamelog.SafeSlice) {
	if area.widget == nil {
		return
	}
	area.widget.SetStore(store)
}

// SetEnabled は有効/無効を設定する
func (area *MessageArea) SetEnabled(enabled bool) {
	area.enabled = enabled
}

// Update はメッセージエリアを更新する
func (area *MessageArea) Update(_ w.World) {
	if !area.enabled || area.widget == nil {
		return
	}

	area.widget.Update()
}

// Draw はメッセージエリアを描画する
func (area *MessageArea) Draw(world w.World, screen *ebiten.Image) {
	if !area.enabled || area.widget == nil {
		return
	}

	// 画面サイズを取得
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// ログエリアの位置とサイズを計算（画面下部、横幅いっぱい）
	logAreaX := 0
	logAreaWidth := screenWidth

	// シンプルに固定サイズで計算
	fixedHeight := defaultLogAreaMargin*2 + defaultMaxLogLines*defaultLineHeight + defaultYPadding*2
	logAreaY := screenHeight - fixedHeight

	// 背景を描画（固定サイズ）
	styled.DrawFramedBackground(screen, logAreaX, logAreaY, logAreaWidth, fixedHeight, styled.DefaultMessageBackgroundStyle())

	// オフスクリーンサイズ（固定）
	offscreenWidth := logAreaWidth - defaultLogAreaMargin*2
	offscreenHeight := fixedHeight - defaultLogAreaMargin*2

	// メッセージウィジェットを描画
	drawX := logAreaX + defaultLogAreaMargin
	drawY := logAreaY + defaultLogAreaMargin

	area.widget.Draw(screen, drawX, drawY, offscreenWidth, offscreenHeight)
}
