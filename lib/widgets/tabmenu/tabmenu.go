package tabmenu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
)

// TabItem はタブの項目を定義する
type TabItem struct {
	ID    string
	Label string
	Items []menu.Item
}

// Config はタブメニューの設定
type Config struct {
	Tabs             []TabItem
	InitialTabIndex  int
	InitialItemIndex int
	WrapNavigation   bool // タブ/アイテム両方で端循環するか
	ItemsPerPage     int  // 1ページに表示する項目数（0=制限なし）
}

// Callbacks はタブメニューのコールバック
type Callbacks struct {
	OnSelectItem func(tabIndex int, itemIndex int, tab TabItem, item menu.Item) error
	OnCancel     func()
	OnTabChange  func(oldTabIndex, newTabIndex int, tab TabItem)
	OnItemChange func(tabIndex int, oldItemIndex, newItemIndex int, item menu.Item) error
}

// TabMenu はタブ付きメニューコンポーネント
// TODO: Menuとの違いを明確にする。TabMenuはUIBuilderを持たず状態管理だけを行っている。なのでインジケーターUIを各自で実装する必要がある。いっぽうでMenuはUIを持ち、インジケーターが含まれている
type TabMenu struct {
	config    Config
	callbacks Callbacks

	// 状態
	currentTabIndex  int
	currentItemIndex int
	keyboardInput    input.KeyboardInput

	// ページネーション状態
	currentPage int // 現在のページ（0ベース）
}

// NewTabMenu は新しいTabMenuを作成する
func NewTabMenu(config Config, callbacks Callbacks, keyboardInput input.KeyboardInput) *TabMenu {
	tm := &TabMenu{
		config:           config,
		callbacks:        callbacks,
		currentTabIndex:  config.InitialTabIndex,
		currentItemIndex: config.InitialItemIndex,
		keyboardInput:    keyboardInput,
	}

	// 初期タブのアイテム数を確認してインデックスを調整
	if len(config.Tabs) > 0 && config.InitialTabIndex < len(config.Tabs) {
		initialTab := config.Tabs[config.InitialTabIndex]
		if len(initialTab.Items) == 0 {
			tm.currentItemIndex = -1
		} else if config.InitialItemIndex >= len(initialTab.Items) {
			tm.currentItemIndex = len(initialTab.Items) - 1
		} else if config.InitialItemIndex < 0 {
			tm.currentItemIndex = 0
		}
	}

	return tm
}

// Update はタブメニューを更新する
func (tm *TabMenu) Update() (bool, error) {
	// タブ切り替え（左右矢印キー）
	handled, err := tm.handleTabNavigation()
	if err != nil {
		return false, err
	}

	// アイテム選択（上下矢印キー）
	itemHandled, err := tm.handleItemNavigation()
	if err != nil {
		return false, err
	}
	if itemHandled {
		handled = true
	}

	// 選択（Enter）
	selectionHandled, err := tm.handleSelection()
	if err != nil {
		return false, err
	}
	if selectionHandled {
		handled = true
	}

	// キャンセル（Escape）
	if tm.handleCancel() {
		handled = true
	}

	return handled, nil
}

// handleTabNavigation はタブ切り替えを処理する
func (tm *TabMenu) handleTabNavigation() (bool, error) {
	tabPressed := tm.keyboardInput.IsKeyJustPressed(ebiten.KeyTab)

	// Shift+Tabの判定
	var shiftTabPressed bool
	if tabPressed {
		shiftPressed := tm.keyboardInput.IsKeyPressed(ebiten.KeyShift)
		if shiftPressed {
			shiftTabPressed = true
			tabPressed = false // Shift+Tabの場合は通常のTabとして扱わない
		}
	}

	if shiftTabPressed {
		if err := tm.navigateToPreviousTab(); err != nil {
			return false, err
		}
		return true, nil
	} else if tabPressed {
		if err := tm.navigateToNextTab(); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// handleItemNavigation はアイテム選択を処理する
func (tm *TabMenu) handleItemNavigation() (bool, error) {
	upPressed := tm.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp)
	downPressed := tm.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown)
	pageUpPressed := tm.keyboardInput.IsKeyJustPressed(ebiten.KeyPageUp)
	pageDownPressed := tm.keyboardInput.IsKeyJustPressed(ebiten.KeyPageDown)

	if upPressed {
		if err := tm.navigateToPreviousItem(); err != nil {
			return false, err
		}
		return true, nil
	} else if downPressed {
		if err := tm.navigateToNextItem(); err != nil {
			return false, err
		}
		return true, nil
	} else if pageUpPressed {
		if err := tm.navigatePageUp(); err != nil {
			return false, err
		}
		return true, nil
	} else if pageDownPressed {
		if err := tm.navigatePageDown(); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// handleSelection は選択を処理する
func (tm *TabMenu) handleSelection() (bool, error) {
	// Enterキーは押下-押上ワンセット制御を使用
	enterPressed := tm.keyboardInput.IsEnterJustPressedOnce()

	if enterPressed {
		if err := tm.selectCurrentItem(); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// handleCancel はキャンセルを処理する
func (tm *TabMenu) handleCancel() bool {
	escapePressed := tm.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape)

	if escapePressed {
		if tm.callbacks.OnCancel != nil {
			tm.callbacks.OnCancel()
		}
		return true
	}

	return false
}

// navigateToPreviousTab は前のタブに移動する
func (tm *TabMenu) navigateToPreviousTab() error {
	oldIndex := tm.currentTabIndex

	if tm.currentTabIndex > 0 {
		tm.currentTabIndex--
	} else if tm.config.WrapNavigation {
		tm.currentTabIndex = len(tm.config.Tabs) - 1
	}

	if oldIndex != tm.currentTabIndex {
		// タブ変更時はアイテムインデックスとページをリセット
		newTab := tm.config.Tabs[tm.currentTabIndex]
		if len(newTab.Items) > 0 {
			tm.currentItemIndex = 0
			tm.currentPage = 0
		} else {
			tm.currentItemIndex = -1 // 空のタブでは無効なインデックス
		}

		if tm.callbacks.OnTabChange != nil {
			tm.callbacks.OnTabChange(oldIndex, tm.currentTabIndex, newTab)
		}

		// アイテムが存在する場合のみnotifyItemChangeを呼ぶ
		if len(newTab.Items) > 0 {
			if err := tm.notifyItemChange(0, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

// navigateToNextTab は次のタブに移動する
func (tm *TabMenu) navigateToNextTab() error {
	oldIndex := tm.currentTabIndex

	if tm.currentTabIndex < len(tm.config.Tabs)-1 {
		tm.currentTabIndex++
	} else if tm.config.WrapNavigation {
		tm.currentTabIndex = 0
	}

	if oldIndex != tm.currentTabIndex {
		// タブ変更時はアイテムインデックスとページをリセット
		newTab := tm.config.Tabs[tm.currentTabIndex]
		if len(newTab.Items) > 0 {
			tm.currentItemIndex = 0
			tm.currentPage = 0
		} else {
			tm.currentItemIndex = -1 // 空のタブでは無効なインデックス
		}

		if tm.callbacks.OnTabChange != nil {
			tm.callbacks.OnTabChange(oldIndex, tm.currentTabIndex, newTab)
		}

		// アイテムが存在する場合のみnotifyItemChangeを呼ぶ
		if len(newTab.Items) > 0 {
			if err := tm.notifyItemChange(0, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

// navigateToPreviousItem は前のアイテムに移動する
func (tm *TabMenu) navigateToPreviousItem() error {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return nil
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 {
		return nil
	}

	oldIndex := tm.currentItemIndex

	// ページネーション対応
	if tm.config.ItemsPerPage > 0 {
		pageStart := tm.currentPage * tm.config.ItemsPerPage

		// ページ内での移動
		if tm.currentItemIndex < 0 {
			tm.currentItemIndex = pageStart
		} else if tm.currentItemIndex > pageStart {
			tm.currentItemIndex--
		} else if tm.currentPage > 0 {
			// 前のページへ
			tm.currentPage--
			tm.currentItemIndex = (tm.currentPage+1)*tm.config.ItemsPerPage - 1
		} else if tm.config.WrapNavigation {
			// 最後のページへ
			tm.currentPage = (len(currentTab.Items) - 1) / tm.config.ItemsPerPage
			tm.currentItemIndex = len(currentTab.Items) - 1
		}
	} else {
		// ページネーションなし
		if tm.currentItemIndex < 0 {
			tm.currentItemIndex = 0
		} else if tm.currentItemIndex > 0 {
			tm.currentItemIndex--
		} else if tm.config.WrapNavigation {
			tm.currentItemIndex = len(currentTab.Items) - 1
		}
	}

	if oldIndex != tm.currentItemIndex {
		if err := tm.notifyItemChange(oldIndex, tm.currentItemIndex); err != nil {
			return err
		}
	}
	return nil
}

// navigateToNextItem は次のアイテムに移動する
func (tm *TabMenu) navigateToNextItem() error {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return nil
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 {
		return nil
	}

	oldIndex := tm.currentItemIndex

	// ページネーション対応
	if tm.config.ItemsPerPage > 0 {
		pageStart := tm.currentPage * tm.config.ItemsPerPage
		pageEnd := pageStart + tm.config.ItemsPerPage
		if pageEnd > len(currentTab.Items) {
			pageEnd = len(currentTab.Items)
		}

		// ページ内での移動
		if tm.currentItemIndex < 0 {
			tm.currentItemIndex = pageStart
		} else if tm.currentItemIndex < pageEnd-1 {
			tm.currentItemIndex++
		} else if pageEnd < len(currentTab.Items) {
			// 次のページへ
			tm.currentPage++
			tm.currentItemIndex = tm.currentPage * tm.config.ItemsPerPage
		} else if tm.config.WrapNavigation {
			// 最初のページへ
			tm.currentPage = 0
			tm.currentItemIndex = 0
		}
	} else {
		// ページネーションなし
		if tm.currentItemIndex < 0 {
			tm.currentItemIndex = 0
		} else if tm.currentItemIndex < len(currentTab.Items)-1 {
			tm.currentItemIndex++
		} else if tm.config.WrapNavigation {
			tm.currentItemIndex = 0
		}
	}

	if oldIndex != tm.currentItemIndex {
		if err := tm.notifyItemChange(oldIndex, tm.currentItemIndex); err != nil {
			return err
		}
	}
	return nil
}

// navigatePageUp は前のページに移動する
func (tm *TabMenu) navigatePageUp() error {
	if tm.config.ItemsPerPage <= 0 {
		return nil
	}

	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return nil
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 {
		return nil
	}

	if tm.currentPage > 0 {
		tm.currentPage--
		tm.currentItemIndex = tm.currentPage * tm.config.ItemsPerPage
		if err := tm.notifyItemChange(tm.currentItemIndex, tm.currentItemIndex); err != nil {
			return err
		}
	}
	return nil
}

// navigatePageDown は次のページに移動する
func (tm *TabMenu) navigatePageDown() error {
	if tm.config.ItemsPerPage <= 0 {
		return nil
	}

	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return nil
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 {
		return nil
	}

	totalPages := (len(currentTab.Items) + tm.config.ItemsPerPage - 1) / tm.config.ItemsPerPage
	if tm.currentPage < totalPages-1 {
		tm.currentPage++
		tm.currentItemIndex = tm.currentPage * tm.config.ItemsPerPage
		if err := tm.notifyItemChange(tm.currentItemIndex, tm.currentItemIndex); err != nil {
			return err
		}
	}
	return nil
}

// selectCurrentItem は現在のアイテムを選択する
func (tm *TabMenu) selectCurrentItem() error {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return nil
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 || tm.currentItemIndex >= len(currentTab.Items) || tm.currentItemIndex < 0 {
		return nil
	}

	currentItem := currentTab.Items[tm.currentItemIndex]

	if tm.callbacks.OnSelectItem != nil {
		if err := tm.callbacks.OnSelectItem(tm.currentTabIndex, tm.currentItemIndex, currentTab, currentItem); err != nil {
			return err
		}
	}
	return nil
}

// notifyItemChange はアイテム変更を通知する
func (tm *TabMenu) notifyItemChange(oldIndex, newIndex int) error {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return nil
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 || newIndex >= len(currentTab.Items) || newIndex < 0 {
		return nil
	}

	if tm.callbacks.OnItemChange != nil {
		if err := tm.callbacks.OnItemChange(tm.currentTabIndex, oldIndex, newIndex, currentTab.Items[newIndex]); err != nil {
			return err
		}
	}
	return nil
}

// GetCurrentTabIndex は現在のタブインデックスを返す
func (tm *TabMenu) GetCurrentTabIndex() int {
	return tm.currentTabIndex
}

// GetCurrentItemIndex は現在のアイテムインデックスを返す
func (tm *TabMenu) GetCurrentItemIndex() int {
	return tm.currentItemIndex
}

// GetCurrentTab は現在のタブを返す
func (tm *TabMenu) GetCurrentTab() TabItem {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return TabItem{}
	}
	return tm.config.Tabs[tm.currentTabIndex]
}

// GetCurrentItem は現在のアイテムを返す
func (tm *TabMenu) GetCurrentItem() menu.Item {
	currentTab := tm.GetCurrentTab()
	if len(currentTab.Items) == 0 || tm.currentItemIndex >= len(currentTab.Items) || tm.currentItemIndex < 0 {
		return menu.Item{}
	}
	return currentTab.Items[tm.currentItemIndex]
}

// SetTabIndex はタブインデックスを設定する
func (tm *TabMenu) SetTabIndex(index int) error {
	if index >= 0 && index < len(tm.config.Tabs) {
		oldIndex := tm.currentTabIndex
		tm.currentTabIndex = index

		// タブ変更時はアイテムインデックスとページをリセット
		newTab := tm.config.Tabs[tm.currentTabIndex]
		if len(newTab.Items) > 0 {
			tm.currentItemIndex = 0
			tm.currentPage = 0
		} else {
			tm.currentItemIndex = -1 // 空のタブでは無効なインデックス
		}

		if oldIndex != tm.currentTabIndex {
			if tm.callbacks.OnTabChange != nil {
				tm.callbacks.OnTabChange(oldIndex, tm.currentTabIndex, newTab)
			}

			// アイテムが存在する場合のみnotifyItemChangeを呼ぶ
			if len(newTab.Items) > 0 {
				if err := tm.notifyItemChange(0, 0); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// SetItemIndex はアイテムインデックスを設定する
func (tm *TabMenu) SetItemIndex(index int) error {
	currentTab := tm.GetCurrentTab()
	if index >= 0 && index < len(currentTab.Items) {
		oldIndex := tm.currentItemIndex
		tm.currentItemIndex = index

		if oldIndex != tm.currentItemIndex {
			if err := tm.notifyItemChange(oldIndex, tm.currentItemIndex); err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateTabs はタブを更新する（動的にアイテムが変更された場合）
func (tm *TabMenu) UpdateTabs(tabs []TabItem) {
	tm.config.Tabs = tabs

	// 現在のインデックスが範囲外になった場合は調整
	if tm.currentTabIndex >= len(tabs) {
		tm.currentTabIndex = len(tabs) - 1
		if tm.currentTabIndex < 0 {
			tm.currentTabIndex = 0
		}
	}

	// 現在のタブのアイテムインデックスを調整
	if len(tabs) > 0 && tm.currentTabIndex < len(tabs) {
		currentTab := tabs[tm.currentTabIndex]
		if len(currentTab.Items) == 0 {
			tm.currentItemIndex = -1 // 空のタブでは無効なインデックス
		} else if tm.currentItemIndex >= len(currentTab.Items) {
			tm.currentItemIndex = len(currentTab.Items) - 1
		} else if tm.currentItemIndex < 0 {
			tm.currentItemIndex = 0
		}

		// UpdateTabsは内部状態の更新のみを行い、コールバックは呼ばない
		// コールバックは呼び出し元で必要に応じて実行される
	}
}

// GetVisibleItems は現在のページで表示される項目とその元のインデックスを返す
func (tm *TabMenu) GetVisibleItems() ([]menu.Item, []int) {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return []menu.Item{}, []int{}
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]

	if tm.config.ItemsPerPage <= 0 {
		indices := make([]int, len(currentTab.Items))
		for i := range indices {
			indices[i] = i
		}
		return currentTab.Items, indices
	}

	start := tm.currentPage * tm.config.ItemsPerPage
	end := start + tm.config.ItemsPerPage
	if end > len(currentTab.Items) {
		end = len(currentTab.Items)
	}

	if start >= len(currentTab.Items) {
		return []menu.Item{}, []int{}
	}

	visibleItems := currentTab.Items[start:end]
	indices := make([]int, len(visibleItems))
	for i := range indices {
		indices[i] = start + i
	}

	return visibleItems, indices
}

// GetCurrentPage は現在のページ番号を返す(表示用なので1ベース)
func (tm *TabMenu) GetCurrentPage() int {
	return tm.currentPage + 1
}

// GetTotalPages は総ページ数を返す
func (tm *TabMenu) GetTotalPages() int {
	if tm.config.ItemsPerPage <= 0 {
		return 1
	}

	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return 1
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	return (len(currentTab.Items) + tm.config.ItemsPerPage - 1) / tm.config.ItemsPerPage
}

// GetPageIndicatorText はページインジケーターのテキストを返す
func (tm *TabMenu) GetPageIndicatorText() string {
	if tm.config.ItemsPerPage <= 0 || tm.GetTotalPages() <= 1 {
		return ""
	}

	arrows := ""

	// 前のページがある場合は上矢印を追加
	if tm.HasPreviousPage() {
		arrows += " ↑"
	} else {
		arrows += " 　"
	}

	// 次のページがある場合は下矢印を追加
	if tm.HasNextPage() {
		arrows += " ↓"
	} else {
		arrows += " 　"
	}

	return fmt.Sprintf("%d/%d%s", tm.GetCurrentPage(), tm.GetTotalPages(), arrows)
}

// HasPreviousPage は前のページがあるかを返す
func (tm *TabMenu) HasPreviousPage() bool {
	return tm.currentPage > 0
}

// HasNextPage は次のページがあるかを返す
func (tm *TabMenu) HasNextPage() bool {
	if tm.config.ItemsPerPage <= 0 {
		return false
	}

	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return false
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	nextPageStart := (tm.currentPage + 1) * tm.config.ItemsPerPage
	return nextPageStart < len(currentTab.Items)
}
