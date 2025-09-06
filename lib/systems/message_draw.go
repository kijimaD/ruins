package systems

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/styles"
	w "github.com/kijimaD/ruins/lib/world"
)

const (
	logAreaHeight = 120 // ログエリアの高さ（余裕を持たせて大きめに）
	maxLogLines   = 5   // 表示する最大行数
	logAreaMargin = 8   // 余白
	lineHeight    = 20  // 1行の高さ
	yPadding      = 8   // 下端の追加パディング
)

var (
	messageUI        *ebitenui.UI // ログメッセージ用のUI
	lastMessageCount int          // 前回のメッセージ数を記録
)

// DrawMessages はログメッセージを画面下部に描画する
func DrawMessages(world w.World, screen *ebiten.Image) {
	// 画面サイズを取得
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// ログエリアの位置とサイズを計算（画面下部、横幅いっぱい）
	logAreaX := 0
	logAreaWidth := screenWidth

	// シンプルに固定サイズで計算
	fixedHeight := logAreaMargin*2 + maxLogLines*lineHeight + yPadding*2
	logAreaY := screenHeight - fixedHeight

	// 背景を描画（固定サイズ）
	drawMessageBackground(screen, logAreaX, logAreaY, logAreaWidth, fixedHeight)

	// UIが初期化されていない場合は初期化
	if messageUI == nil {
		initMessageUI(world)
	}

	// ログメッセージが更新されている場合はUIを再構築
	updateMessageUI(world)

	// UIを更新
	messageUI.Update()

	// オフスクリーンサイズ（固定）
	offscreenWidth := logAreaWidth - logAreaMargin*2
	offscreenHeight := fixedHeight - logAreaMargin*2

	offscreen := ebiten.NewImage(offscreenWidth, offscreenHeight)
	messageUI.Draw(offscreen)

	// 描画位置を調整（枠内に正確に配置）
	op := &ebiten.DrawImageOptions{}
	drawY := logAreaY + logAreaMargin
	op.GeoM.Translate(float64(logAreaX+logAreaMargin), float64(drawY))
	screen.DrawImage(offscreen, op)
}

// initMessageUI はメッセージUI用の初期化を行う
func initMessageUI(world w.World) {
	// 初期状態でメッセージを取得
	messages := gamelog.FieldLog.Get()

	// ログ用コンテナを作成（シンプルな縦並び）
	logContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(3),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    2,
					Bottom: 2, // パディングを最小限に
					Left:   2,
					Right:  2,
				}),
			),
		),
	)

	// 実際のメッセージを追加
	startIndex := 0
	if len(messages) > maxLogLines {
		startIndex = len(messages) - maxLogLines
	}

	for _, message := range messages[startIndex:] {
		if message == "" {
			continue
		}
		messageWidget := eui.NewListItemText(message, styles.TextColor, false, world)
		logContainer.AddChild(messageWidget)
	}

	// メッセージがない場合
	if len(messages) == 0 {
		placeholderWidget := eui.NewListItemText("ログメッセージなし", styles.ForegroundColor, false, world)
		logContainer.AddChild(placeholderWidget)
	}

	// UIを初期化（シンプルに）
	messageUI = &ebitenui.UI{
		Container: logContainer,
	}

	// 初期メッセージ数を設定
	lastMessageCount = len(messages)
}

// updateMessageUI はログメッセージが更新された場合にUIを再構築する
func updateMessageUI(world w.World) {
	messages := gamelog.FieldLog.Get()
	currentMessageCount := len(messages)

	// メッセージ数が変わっていない場合は更新不要
	if currentMessageCount == lastMessageCount {
		return
	}

	// ログ用コンテナを作成（シンプルな縦並び）
	logContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(3),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    2,
					Bottom: 2, // パディングを最小限に
					Left:   2,
					Right:  2,
				}),
			),
		),
	)

	// 最新のメッセージのみを表示（最大maxLogLines行）
	startIndex := 0
	if len(messages) > maxLogLines {
		startIndex = len(messages) - maxLogLines
	}

	// メッセージを追加
	for _, message := range messages[startIndex:] {
		if message == "" {
			continue
		}
		messageWidget := eui.NewListItemText(message, styles.TextColor, false, world)
		logContainer.AddChild(messageWidget)
	}

	// メッセージがない場合のプレースホルダー
	if len(messages) == 0 {
		placeholderWidget := eui.NewListItemText("ログメッセージなし", styles.ForegroundColor, false, world)
		logContainer.AddChild(placeholderWidget)
	}

	// UIを更新（シンプルに）
	messageUI.Container = logContainer

	// メッセージ数を更新
	lastMessageCount = currentMessageCount
}

// drawMessageBackground はログエリアの背景を描画する
func drawMessageBackground(screen *ebiten.Image, x, y, width, height int) {
	// 枠線を描画（白色）
	vector.StrokeRect(screen,
		float32(x),
		float32(y),
		float32(width),
		float32(height),
		2,                              // 線の太さ
		color.RGBA{255, 255, 255, 255}, // 白色の枠線
		false)

	// 内側の背景を描画（黒色、半透明）
	vector.DrawFilledRect(screen,
		float32(x+2),
		float32(y+2),
		float32(width-4),
		float32(height-4),
		color.RGBA{0, 0, 0, 200}, // 半透明の黒背景
		false)
}
