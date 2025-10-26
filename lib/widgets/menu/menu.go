package menu

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/inputmapper"
)

// Item はメニュー項目を表す
type Item struct {
	ID               string
	Label            string
	AdditionalLabels []string // 追加表示項目（個数、価格など）右側に表示される
	Disabled         bool
	UserData         interface{} // 任意のデータを保持
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
	ItemsPerPage int // 1ページに表示する項目数（0=制限なし）
}

// Callbacks はメニューのコールバック
type Callbacks struct {
	OnSelect      func(index int, item Item) error
	OnCancel      func()
	OnFocusChange func(oldIndex, newIndex int)
}

// Menu は共通メニューコンポーネント
type Menu struct {
	config    Config
	callbacks Callbacks

	// 基本状態
	focusedIndex int

	// ペジネーション状態
	currentPage int // 現在のページ（0ベース）

	// UI要素
	container   *widget.Container
	itemWidgets []widget.PreferredSizeLocateableWidget
	uiBuilder   *UIBuilder // UIビルダーの参照を保持
}

// NewMenu はメニューを作成する
func NewMenu(config Config, callbacks Callbacks) *Menu {
	m := &Menu{
		config:       config,
		callbacks:    callbacks,
		focusedIndex: config.InitialIndex,
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

// Update はキーボード入力を待ち受けて、Actionに変換してメニュー操作を実行する
func (m *Menu) Update() error {
	keyboardInput := input.GetSharedKeyboardInput()
	if action, ok := m.translateKeyToAction(keyboardInput); ok {
		return m.DoAction(action)
	}
	return nil
}

// translateKeyToAction はキーボード入力をActionに変換する
func (m *Menu) translateKeyToAction(keyboardInput input.KeyboardInput) (inputmapper.ActionID, bool) {
	// 下矢印キー
	if keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return inputmapper.ActionMenuDown, true
	}

	// 上矢印キー
	if keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return inputmapper.ActionMenuUp, true
	}

	// Tabキー
	if keyboardInput.IsKeyJustPressed(ebiten.KeyTab) {
		if keyboardInput.IsKeyPressed(ebiten.KeyShift) {
			return inputmapper.ActionMenuUp, true
		}
		return inputmapper.ActionMenuDown, true
	}

	// Enterキー（セッションベース検出で複数回実行を防止）
	if keyboardInput.IsEnterJustPressedOnce() {
		return inputmapper.ActionMenuSelect, true
	}

	// Escapeキー
	if keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		return inputmapper.ActionMenuCancel, true
	}

	return "", false
}

// DoAction はActionを受け取ってメニュー操作を実行する
func (m *Menu) DoAction(action inputmapper.ActionID) error {
	switch action {
	case inputmapper.ActionMenuDown:
		m.navigateNext()
	case inputmapper.ActionMenuUp:
		m.navigatePrevious()
	case inputmapper.ActionMenuSelect:
		return m.selectCurrent()
	case inputmapper.ActionMenuCancel:
		if m.callbacks.OnCancel != nil {
			m.callbacks.OnCancel()
		}
	}
	return nil
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
func (m *Menu) selectCurrent() error {
	if m.focusedIndex < 0 || m.focusedIndex >= len(m.config.Items) {
		return nil
	}

	item := m.config.Items[m.focusedIndex]
	if !item.Disabled && m.callbacks.OnSelect != nil {
		return m.callbacks.OnSelect(m.focusedIndex, item)
	}
	return nil
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

// ================ スクロール関連メソッド ================

// initializePagination はページネーション設定を初期化する
func (m *Menu) initializePagination() {
	m.currentPage = 0
}

// updatePageFromFocus はフォーカスに基づいてページを更新する
func (m *Menu) updatePageFromFocus() {
	if m.config.ItemsPerPage <= 0 {
		return // スクロール無効
	}

	newPage := m.focusedIndex / m.config.ItemsPerPage
	if newPage != m.currentPage {
		m.currentPage = newPage
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
	if m.config.ItemsPerPage <= 0 || m.GetTotalPages() <= 1 {
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
