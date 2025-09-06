package typewriter

import (
	"strings"

	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
)

// UIConfig はUI作成時の設定
type UIConfig struct {
	MaxVisibleLines int                  // 最大表示行数（デフォルト: 4）
	FixedWidth      int                  // 固定幅（デフォルト: 720）
	LineHeight      int                  // 行の高さ（デフォルト: 25）
	TextFace        *text.Face           // テキストのフォントフェイス
	TextColor       color.Color          // テキスト色
	PanelImage      *image.NineSlice     // 背景パネル画像
	ArrowImage      *widget.GraphicImage // 矢印画像
}

// DefaultUIConfig はデフォルトのUI設定を返す
func DefaultUIConfig() UIConfig {
	return UIConfig{
		MaxVisibleLines: 4,
		FixedWidth:      720,
		LineHeight:      25,
		TextColor:       colornames.White,
	}
}

// MessageUIBuilder はメッセージUI作成を統合管理
type MessageUIBuilder struct {
	messageHandler *MessageHandler
	config         UIConfig
	container      *widget.Container
}

// NewMessageUIBuilder は新しいUIビルダーを作成
func NewMessageUIBuilder(messageHandler *MessageHandler, config UIConfig) *MessageUIBuilder {
	builder := &MessageUIBuilder{
		messageHandler: messageHandler,
		config:         config,
	}

	// UI更新フックを設定
	messageHandler.SetOnUpdateUI(func(_ string) {
		builder.container = builder.createContainer()
	})

	// 初回コンテナ作成
	builder.container = builder.createContainer()

	return builder
}

// GetContainer はメッセージ表示用のコンテナを取得
func (b *MessageUIBuilder) GetContainer() *widget.Container {
	return b.container
}

// Update はUIの更新処理（アニメーション更新など）
func (b *MessageUIBuilder) Update() {
	// プロンプトアニメーション更新（コンテナ再構築が必要）
	if b.messageHandler != nil && b.messageHandler.IsWaitingForInput() {
		b.container = b.createContainer()
	}
}

// createContainer はメッセージ表示用のコンテナを作成
func (b *MessageUIBuilder) createContainer() *widget.Container {
	// 固定サイズの計算
	fixedHeight := b.config.LineHeight * b.config.MaxVisibleLines
	fixedWidth := b.config.FixedWidth

	// メッセージ全体のコンテナ（水平レイアウト）
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(5), // 水平方向のスペース
		)),
	)

	// メッセージ表示用のコンテナ（固定サイズ）
	messageContainerOpts := []widget.ContainerOpt{
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(fixedWidth, fixedHeight), // 固定サイズを直接指定
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	}

	// 背景画像があれば追加
	if b.config.PanelImage != nil {
		messageContainerOpts = append(messageContainerOpts,
			widget.ContainerOpts.BackgroundImage(b.config.PanelImage))
	}

	messageContainer := widget.NewContainer(messageContainerOpts...)

	// テキスト表示用の垂直コンテナ（メッセージコンテナ内に配置）
	textContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionStart,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	// タイプライター表示中のテキストを取得
	currentDisplayText := ""
	if b.messageHandler != nil {
		currentDisplayText = b.messageHandler.GetDisplayText()
	}

	// 現在表示中のテキストを行に分割
	currentLines := strings.Split(currentDisplayText, "\n")

	// 表示する行を制限（最新のmaxVisibleLines行のみ）
	displayLines := currentLines
	if len(currentLines) > b.config.MaxVisibleLines {
		// 古い行を削除して最新の行のみ表示
		displayLines = currentLines[len(currentLines)-b.config.MaxVisibleLines:]
	}

	// 各行をテキストウィジェットとして追加
	for _, lineText := range displayLines {
		// テキストフェイスがnilの場合はスキップ
		if b.config.TextFace == nil {
			continue
		}

		textWidget := widget.NewText(
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				MaxWidth: fixedWidth - 20, // テキストの最大幅を制限（パディング分を考慮）
			})),
			widget.TextOpts.Text(lineText, b.config.TextFace, b.config.TextColor),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		)

		textContainer.AddChild(textWidget)
	}

	// テキストコンテナをメッセージコンテナに追加
	messageContainer.AddChild(textContainer)

	// メッセージコンテナをコンテナに追加
	container.AddChild(messageContainer)

	// プロンプト表示が必要な場合、右側に配置
	if b.messageHandler != nil && b.messageHandler.IsWaitingForInput() && b.config.ArrowImage != nil {
		// プロンプトコンテナを作成
		promptContainer := b.messageHandler.CreatePromptContainer(b.config.ArrowImage)
		if promptContainer != nil {
			// プロンプト用の右側コンテナ（メッセージコンテナと同じ高さ）
			promptSideContainer := widget.NewContainer(
				widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
				})),
				widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			)

			// プロンプトをコンテナの下端に配置（固定高さのボトムに合わせる）
			promptWrapper := widget.NewContainer(
				widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionStart,
					Padding: &widget.Insets{
						Top: fixedHeight - 25, // メッセージコンテナの高さ - プロンプト高さ
					},
				})),
				widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			)

			promptWrapper.AddChild(promptContainer)
			promptSideContainer.AddChild(promptWrapper)
			container.AddChild(promptSideContainer)
		}
	}

	return container
}
