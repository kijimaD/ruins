package menu

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
)

// MenuItem はメニュー項目を表す
type MenuItem struct {
	ID          string
	Label       string
	Disabled    bool
	Icon        *ebiten.Image
	Description string      // ツールチップや説明文
	UserData    interface{} // 任意のデータを保持
}

// Orientation はメニューの向き
type Orientation int

const (
	// Vertical は縦方向配置
	Vertical Orientation = iota
	// Horizontal は横方向配置
	Horizontal
)

// MenuConfig はメニューの設定
type MenuConfig struct {
	Items          []MenuItem
	InitialIndex   int
	WrapNavigation bool        // 端で循環するか
	Orientation    Orientation // Vertical or Horizontal
	Columns        int         // グリッド表示時の列数（0=リスト表示）
}

// MenuCallbacks はメニューのコールバック
type MenuCallbacks struct {
	OnSelect      func(index int, item MenuItem)
	OnCancel      func()
	OnFocusChange func(oldIndex, newIndex int)
	OnHover       func(index int, item MenuItem)
}

// Menu は共通メニューコンポーネント
type Menu struct {
	config    MenuConfig
	callbacks MenuCallbacks

	// 状態
	focusedIndex int
	hoveredIndex int

	// UI要素
	container   *widget.Container
	itemWidgets []widget.PreferredSizeLocateableWidget

	// 入力
	keyboardInput input.KeyboardInput
}

// NewMenu はメニューを作成する
func NewMenu(config MenuConfig, callbacks MenuCallbacks) *Menu {
	m := &Menu{
		config:       config,
		callbacks:    callbacks,
		focusedIndex: config.InitialIndex,
		hoveredIndex: -1,
	}

	// 初期インデックスの検証
	if m.focusedIndex < 0 || m.focusedIndex >= len(config.Items) ||
		(len(config.Items) > 0 && config.Items[m.focusedIndex].Disabled) {
		m.focusedIndex = m.findFirstEnabled()
	}

	return m
}

// Update はメニューの状態を更新する
func (m *Menu) Update(keyboardInput input.KeyboardInput) {
	m.keyboardInput = keyboardInput
	m.handleKeyboard()
}

// GetFocusedIndex は現在フォーカスされている項目のインデックスを返す
func (m *Menu) GetFocusedIndex() int {
	return m.focusedIndex
}

// SetFocusedIndex はフォーカスする項目を設定する
func (m *Menu) SetFocusedIndex(index int) {
	if index >= 0 && index < len(m.config.Items) && !m.config.Items[index].Disabled {
		oldIndex := m.focusedIndex
		m.focusedIndex = index
		m.updateFocus()
		if m.callbacks.OnFocusChange != nil {
			m.callbacks.OnFocusChange(oldIndex, m.focusedIndex)
		}
	}
}

// GetItems はメニュー項目を返す
func (m *Menu) GetItems() []MenuItem {
	return m.config.Items
}

// GetContainer はUIコンテナを返す
func (m *Menu) GetContainer() *widget.Container {
	return m.container
}

// SetContainer はUIコンテナを設定する
func (m *Menu) SetContainer(container *widget.Container) {
	m.container = container
}

// handleKeyboard はキーボード入力を処理する
func (m *Menu) handleKeyboard() {
	if len(m.config.Items) == 0 {
		return
	}

	oldIndex := m.focusedIndex
	handled := false

	// 基本的なナビゲーション
	if m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown) {
		m.navigateNext()
		handled = true
	} else if m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp) {
		m.navigatePrevious()
		handled = true
	}

	// グリッド表示時の左右移動
	if m.config.Columns > 0 {
		if m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowRight) {
			m.navigateRight()
			handled = true
		} else if m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			m.navigateLeft()
			handled = true
		}
	}

	// Tab/Shift+Tab
	if m.keyboardInput.IsKeyJustPressed(ebiten.KeyTab) {
		if m.keyboardInput.IsKeyPressed(ebiten.KeyShift) {
			m.navigatePrevious()
		} else {
			m.navigateNext()
		}
		handled = true
	}

	// 選択（Enterは押下-押上ワンセット、Spaceは通常）
	enterPressed := m.keyboardInput.IsEnterJustPressedOnce()
	spacePressed := m.keyboardInput.IsKeyJustPressed(ebiten.KeySpace)

	if enterPressed || spacePressed {
		m.selectCurrent()
		handled = true
	}

	// キャンセル
	if m.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		if m.callbacks.OnCancel != nil {
			m.callbacks.OnCancel()
		}
		handled = true
	}

	// フォーカス変更の通知
	if handled && oldIndex != m.focusedIndex && m.callbacks.OnFocusChange != nil {
		m.callbacks.OnFocusChange(oldIndex, m.focusedIndex)
	}
}

// navigateNext は次の有効な項目に移動する
func (m *Menu) navigateNext() {
	itemCount := len(m.config.Items)
	if itemCount == 0 {
		return
	}

	startIndex := m.focusedIndex
	for i := 0; i < itemCount; i++ {
		nextIndex := (startIndex + i + 1) % itemCount

		// 循環しない場合は端で止まる
		if !m.config.WrapNavigation && nextIndex <= startIndex {
			break
		}

		if !m.config.Items[nextIndex].Disabled {
			m.focusedIndex = nextIndex
			m.updateFocus()
			break
		}
	}
}

// navigatePrevious は前の有効な項目に移動する
func (m *Menu) navigatePrevious() {
	itemCount := len(m.config.Items)
	if itemCount == 0 {
		return
	}

	startIndex := m.focusedIndex
	for i := 0; i < itemCount; i++ {
		prevIndex := (startIndex - i - 1 + itemCount) % itemCount

		// 循環しない場合は端で止まる
		if !m.config.WrapNavigation && prevIndex >= startIndex {
			break
		}

		if !m.config.Items[prevIndex].Disabled {
			m.focusedIndex = prevIndex
			m.updateFocus()
			break
		}
	}
}

// navigateRight はグリッド表示時の右移動
func (m *Menu) navigateRight() {
	if m.config.Columns <= 0 {
		return
	}

	itemCount := len(m.config.Items)
	currentRow := m.focusedIndex / m.config.Columns
	currentCol := m.focusedIndex % m.config.Columns

	// 右の列に移動
	if currentCol < m.config.Columns-1 {
		nextIndex := currentRow*m.config.Columns + currentCol + 1
		if nextIndex < itemCount && !m.config.Items[nextIndex].Disabled {
			m.focusedIndex = nextIndex
			m.updateFocus()
		}
	}
}

// navigateLeft はグリッド表示時の左移動
func (m *Menu) navigateLeft() {
	if m.config.Columns <= 0 {
		return
	}

	currentRow := m.focusedIndex / m.config.Columns
	currentCol := m.focusedIndex % m.config.Columns

	// 左の列に移動
	if currentCol > 0 {
		nextIndex := currentRow*m.config.Columns + currentCol - 1
		if !m.config.Items[nextIndex].Disabled {
			m.focusedIndex = nextIndex
			m.updateFocus()
		}
	}
}

// selectCurrent は現在の項目を選択する
func (m *Menu) selectCurrent() {
	if m.focusedIndex < 0 || m.focusedIndex >= len(m.config.Items) {
		return
	}

	item := m.config.Items[m.focusedIndex]
	if !item.Disabled && m.callbacks.OnSelect != nil {
		m.callbacks.OnSelect(m.focusedIndex, item)
	}
}

// findFirstEnabled は最初の有効な項目のインデックスを返す
func (m *Menu) findFirstEnabled() int {
	for i, item := range m.config.Items {
		if !item.Disabled {
			return i
		}
	}
	return 0
}

// updateFocus はフォーカス状態を更新する（UIがある場合に使用）
func (m *Menu) updateFocus() {
	// UIコンテナが設定されている場合のフォーカス更新処理
	// この部分は後でUIビルダーと連携して実装
}
