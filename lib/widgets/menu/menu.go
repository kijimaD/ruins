package menu

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
)

// Item はメニュー項目を表す
type Item struct {
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

// Config はメニューの設定
type Config struct {
	Items          []Item
	InitialIndex   int
	WrapNavigation bool        // 端で循環するか
	Orientation    Orientation // Vertical or Horizontal
	// ペジネーション設定
	ItemsPerPage      int  // 1ページに表示する項目数（0=制限なし）
	ShowPageIndicator bool // ページインジケーターを表示するか
}

// Callbacks はメニューのコールバック
type Callbacks struct {
	OnSelect      func(index int, item Item)
	OnCancel      func()
	OnFocusChange func(oldIndex, newIndex int)
	OnHover       func(index int, item Item)
}

// Menu は共通メニューコンポーネント
type Menu struct {
	config    Config
	callbacks Callbacks

	// 基本状態
	focusedIndex int
	hoveredIndex int

	// ペジネーション状態
	currentPage    int  // 現在のページ（0ベース）
	needsUIRebuild bool // UI再構築が必要かどうか

	// UI要素
	container   *widget.Container
	itemWidgets []widget.PreferredSizeLocateableWidget
	uiBuilder   *UIBuilder // UIビルダーの参照を保持

	// 入力
	keyboardInput input.KeyboardInput
}

// NewMenu はメニューを作成する
func NewMenu(config Config, callbacks Callbacks) *Menu {
	m := &Menu{
		config:       config,
		callbacks:    callbacks,
		focusedIndex: config.InitialIndex,
		hoveredIndex: -1,
	}

	// ページ設定の初期化
	m.initializePagination()

	// 初期インデックスの検証
	if m.focusedIndex < 0 || m.focusedIndex >= len(config.Items) ||
		(len(config.Items) > 0 && config.Items[m.focusedIndex].Disabled) {
		m.focusedIndex = m.findFirstEnabled()
	}

	// フォーカスされた項目に基づいて現在のページを設定
	m.updatePageFromFocus()

	return m
}

// Update はメニューの状態を更新する
func (m *Menu) Update(keyboardInput input.KeyboardInput) {
	m.keyboardInput = keyboardInput
	m.handleKeyboard()
	m.updateFocus()
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
		m.updatePageFromFocus() // ページを更新
		if m.callbacks.OnFocusChange != nil {
			m.callbacks.OnFocusChange(oldIndex, m.focusedIndex)
		}
	}
}

// GetItems はメニュー項目を返す
func (m *Menu) GetItems() []Item {
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

// SetUIBuilder はUIビルダーを設定する
func (m *Menu) SetUIBuilder(builder *UIBuilder) {
	m.uiBuilder = builder
}

// handleKeyboard はキーボード入力を処理する
func (m *Menu) handleKeyboard() {
	if len(m.config.Items) == 0 {
		return
	}

	oldIndex := m.focusedIndex
	handled := false

	// 基本的なナビゲーション（スクロール対応）
	if m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown) {
		m.navigateNext()
		handled = true
	} else if m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp) {
		m.navigatePrevious()
		handled = true
	}

	// ページ移動（PageUp/PageDownまたはCtrl+Up/Down）
	if m.keyboardInput.IsKeyJustPressed(ebiten.KeyPageUp) ||
		(m.keyboardInput.IsKeyPressed(ebiten.KeyControl) && m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp)) {
		m.navigatePageUp()
		handled = true
	} else if m.keyboardInput.IsKeyJustPressed(ebiten.KeyPageDown) ||
		(m.keyboardInput.IsKeyPressed(ebiten.KeyControl) && m.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown)) {
		m.navigatePageDown()
		handled = true
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

	// 選択（Enterは押下-押上ワンセット）
	enterPressed := m.keyboardInput.IsEnterJustPressedOnce()

	if enterPressed {
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

	// 次の有効な項目を探す
	currentIndex := m.focusedIndex
	for i := 1; i < itemCount; i++ {
		nextIndex := currentIndex + i

		// 範囲外の場合
		if nextIndex >= itemCount {
			if m.config.WrapNavigation {
				nextIndex = nextIndex % itemCount
			} else {
				break // 循環しない場合は停止
			}
		}

		if !m.config.Items[nextIndex].Disabled {
			oldIndex := m.focusedIndex
			m.focusedIndex = nextIndex
			m.updatePageFromFocus() // ページを更新

			// フォーカス変更を通知
			if m.callbacks.OnFocusChange != nil {
				m.callbacks.OnFocusChange(oldIndex, m.focusedIndex)
			}
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

	// 前の有効な項目を探す
	currentIndex := m.focusedIndex
	for i := 1; i < itemCount; i++ {
		prevIndex := currentIndex - i

		// 範囲外の場合
		if prevIndex < 0 {
			if m.config.WrapNavigation {
				prevIndex = (prevIndex + itemCount) % itemCount
			} else {
				break // 循環しない場合は停止
			}
		}

		if !m.config.Items[prevIndex].Disabled {
			oldIndex := m.focusedIndex
			m.focusedIndex = prevIndex
			m.updatePageFromFocus() // ページを更新

			// フォーカス変更を通知
			if m.callbacks.OnFocusChange != nil {
				m.callbacks.OnFocusChange(oldIndex, m.focusedIndex)
			}
			break
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
	// UI再構築が必要かチェック
	if m.needsUIRebuild && m.uiBuilder != nil {
		m.rebuildUI()
		m.needsUIRebuild = false
	}

	// フォーカス更新
	if m.uiBuilder != nil {
		m.uiBuilder.UpdateFocus(m)
	}
}

// rebuildUI はUI全体を再構築する
func (m *Menu) rebuildUI() {
	if m.container == nil || m.uiBuilder == nil {
		return
	}

	// コンテナをクリアして再構築
	m.container.RemoveChildren()

	// ページインジケーターを追加
	if m.config.ShowPageIndicator && m.config.ItemsPerPage > 0 && m.GetTotalPages() > 1 {
		pageIndicator := m.uiBuilder.CreatePageIndicator(m)
		m.container.AddChild(pageIndicator)
	}

	// 現在のページの項目を追加
	m.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)
	visibleItems, indices := m.GetVisibleItems()

	for i, item := range visibleItems {
		originalIndex := indices[i]
		btn := m.uiBuilder.CreateMenuButton(m, originalIndex, item)
		m.container.AddChild(btn)
		m.itemWidgets = append(m.itemWidgets, btn)
	}
}

// ================ スクロール関連メソッド ================

// initializePagination はページネーション設定を初期化する
func (m *Menu) initializePagination() {
	m.currentPage = 0
	m.needsUIRebuild = false
}

// updatePageFromFocus はフォーカスに基づいてページを更新する
func (m *Menu) updatePageFromFocus() {
	if m.config.ItemsPerPage <= 0 {
		return // スクロール無効
	}

	newPage := m.focusedIndex / m.config.ItemsPerPage
	if newPage != m.currentPage {
		m.currentPage = newPage
		m.needsUIRebuild = true // UI再構築をマーク
	}
}

// navigatePageDown は次のページに移動
func (m *Menu) navigatePageDown() {
	if m.config.ItemsPerPage <= 0 {
		return
	}

	itemCount := len(m.config.Items)
	nextPageStart := (m.currentPage + 1) * m.config.ItemsPerPage

	if nextPageStart >= itemCount {
		return // 最後のページ
	}

	// 次のページの最初の有効な項目にフォーカス
	for i := nextPageStart; i < itemCount && i < nextPageStart+m.config.ItemsPerPage; i++ {
		if !m.config.Items[i].Disabled {
			oldIndex := m.focusedIndex
			m.focusedIndex = i
			m.currentPage = i / m.config.ItemsPerPage
			m.needsUIRebuild = true

			if m.callbacks.OnFocusChange != nil {
				m.callbacks.OnFocusChange(oldIndex, m.focusedIndex)
			}
			return
		}
	}
}

// navigatePageUp は前のページに移動
func (m *Menu) navigatePageUp() {
	if m.config.ItemsPerPage <= 0 {
		return
	}

	if m.currentPage <= 0 {
		return // 既に最初のページ
	}

	prevPageStart := (m.currentPage - 1) * m.config.ItemsPerPage
	prevPageEnd := prevPageStart + m.config.ItemsPerPage

	// 前のページの最初の有効な項目にフォーカス
	for i := prevPageStart; i < prevPageEnd && i < len(m.config.Items); i++ {
		if !m.config.Items[i].Disabled {
			oldIndex := m.focusedIndex
			m.focusedIndex = i
			m.currentPage = i / m.config.ItemsPerPage
			m.needsUIRebuild = true

			if m.callbacks.OnFocusChange != nil {
				m.callbacks.OnFocusChange(oldIndex, m.focusedIndex)
			}
			return
		}
	}
}

// GetCurrentPage は現在のページ番号を返す（1ベース）
func (m *Menu) GetCurrentPage() int {
	return m.currentPage + 1
}

// GetTotalPages は総ページ数を返す
func (m *Menu) GetTotalPages() int {
	if m.config.ItemsPerPage <= 0 {
		return 1
	}
	return (len(m.config.Items) + m.config.ItemsPerPage - 1) / m.config.ItemsPerPage
}

// GetVisibleItems は現在のページで表示される項目とその元のインデックスを返す
func (m *Menu) GetVisibleItems() ([]Item, []int) {
	if m.config.ItemsPerPage <= 0 {
		indices := make([]int, len(m.config.Items))
		for i := range indices {
			indices[i] = i
		}
		return m.config.Items, indices
	}

	start := m.currentPage * m.config.ItemsPerPage
	end := start + m.config.ItemsPerPage
	if end > len(m.config.Items) {
		end = len(m.config.Items)
	}

	visibleItems := m.config.Items[start:end]
	indices := make([]int, len(visibleItems))
	for i := range indices {
		indices[i] = start + i
	}

	return visibleItems, indices
}

// GetPageIndicatorText はページインジケーターのテキストを返す
func (m *Menu) GetPageIndicatorText() string {
	if !m.config.ShowPageIndicator || m.config.ItemsPerPage <= 0 {
		return ""
	}

	if m.GetTotalPages() <= 1 {
		return ""
	}

	arrows := ""

	// 前のページがある場合は上矢印を追加
	if m.HasPreviousPage() {
		arrows += " ↑"
	} else {
		arrows += " 　"
	}

	// 次のページがある場合は下矢印を追加
	if m.HasNextPage() {
		arrows += " ↓"
	} else {
		arrows += " 　"
	}

	return fmt.Sprintf("%d/%d%s", m.GetCurrentPage(), m.GetTotalPages(), arrows)
}

// HasPreviousPage は前のページがあるかを返す
func (m *Menu) HasPreviousPage() bool {
	return m.currentPage > 0
}

// HasNextPage は次のページがあるかを返す
func (m *Menu) HasNextPage() bool {
	if m.config.ItemsPerPage <= 0 {
		return false
	}
	nextPageStart := (m.currentPage + 1) * m.config.ItemsPerPage
	return nextPageStart < len(m.config.Items)
}
