package tabmenu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
)

// TabItem はタブの項目を定義する
type TabItem struct {
	ID    string
	Label string
	Items []menu.MenuItem
}

// TabMenuConfig はタブメニューの設定
type TabMenuConfig struct {
	Tabs              []TabItem
	InitialTabIndex   int
	InitialItemIndex  int
	WrapNavigation    bool // タブ/アイテム両方で端循環するか
	OnlyDifferentKeys bool // 前回と異なるキーのみ受け付けるか
}

// TabMenuCallbacks はタブメニューのコールバック
type TabMenuCallbacks struct {
	OnSelectItem func(tabIndex int, itemIndex int, tab TabItem, item menu.MenuItem)
	OnCancel     func()
	OnTabChange  func(oldTabIndex, newTabIndex int, tab TabItem)
	OnItemChange func(tabIndex int, oldItemIndex, newItemIndex int, item menu.MenuItem)
}

// TabMenu はタブ付きメニューコンポーネント
type TabMenu struct {
	config    TabMenuConfig
	callbacks TabMenuCallbacks

	// 状態
	currentTabIndex  int
	currentItemIndex int
	keyboardInput    input.KeyboardInput
}

// NewTabMenu は新しいTabMenuを作成する
func NewTabMenu(config TabMenuConfig, callbacks TabMenuCallbacks, keyboardInput input.KeyboardInput) *TabMenu {
	return &TabMenu{
		config:           config,
		callbacks:        callbacks,
		currentTabIndex:  config.InitialTabIndex,
		currentItemIndex: config.InitialItemIndex,
		keyboardInput:    keyboardInput,
	}
}

// Update はタブメニューを更新する
func (tm *TabMenu) Update() bool {
	handled := false

	// タブ切り替え（左右矢印キー）
	if tm.handleTabNavigation() {
		handled = true
	}

	// アイテム選択（上下矢印キー）
	if tm.handleItemNavigation() {
		handled = true
	}

	// 選択（Enter/Space）
	if tm.handleSelection() {
		handled = true
	}

	// キャンセル（Escape）
	if tm.handleCancel() {
		handled = true
	}

	return handled
}

// handleTabNavigation はタブ切り替えを処理する
func (tm *TabMenu) handleTabNavigation() bool {
	var leftPressed, rightPressed, tabPressed, shiftTabPressed bool

	// ナビゲーションキー（矢印キー、Tab/Shift+Tab）はOnlyDifferentKeysの対象外とする
	// これにより一貫したナビゲーション体験を提供する
	if tm.config.OnlyDifferentKeys {
		leftPressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeyArrowLeft)
		rightPressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeyArrowRight)
	} else {
		leftPressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowLeft)
		rightPressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowRight)
	}

	// Tab/Shift+Tabキー
	if tm.config.OnlyDifferentKeys {
		tabPressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeyTab)
	} else {
		tabPressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeyTab)
	}
	
	// Shift+Tabの判定
	if tabPressed {
		shiftPressed := tm.keyboardInput.IsKeyPressed(ebiten.KeyShift)
		if shiftPressed {
			shiftTabPressed = true
			tabPressed = false // Shift+Tabの場合は通常のTabとして扱わない
		}
	}

	if leftPressed || shiftTabPressed {
		tm.navigateToPreviousTab()
		return true
	} else if rightPressed || tabPressed {
		tm.navigateToNextTab()
		return true
	}

	return false
}

// handleItemNavigation はアイテム選択を処理する
func (tm *TabMenu) handleItemNavigation() bool {
	var upPressed, downPressed bool

	if tm.config.OnlyDifferentKeys {
		upPressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeyArrowUp)
		downPressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeyArrowDown)
	} else {
		upPressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp)
		downPressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown)
	}

	if upPressed {
		tm.navigateToPreviousItem()
		return true
	} else if downPressed {
		tm.navigateToNextItem()
		return true
	}

	return false
}

// handleSelection は選択を処理する
func (tm *TabMenu) handleSelection() bool {
	var enterPressed, spacePressed bool

	if tm.config.OnlyDifferentKeys {
		enterPressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeyEnter)
		spacePressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeySpace)
	} else {
		enterPressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeyEnter)
		spacePressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeySpace)
	}

	if enterPressed || spacePressed {
		tm.selectCurrentItem()
		return true
	}

	return false
}

// handleCancel はキャンセルを処理する
func (tm *TabMenu) handleCancel() bool {
	var escapePressed bool

	if tm.config.OnlyDifferentKeys {
		escapePressed = tm.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeyEscape)
	} else {
		escapePressed = tm.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape)
	}

	if escapePressed {
		if tm.callbacks.OnCancel != nil {
			tm.callbacks.OnCancel()
		}
		return true
	}

	return false
}

// navigateToPreviousTab は前のタブに移動する
func (tm *TabMenu) navigateToPreviousTab() {
	oldIndex := tm.currentTabIndex

	if tm.currentTabIndex > 0 {
		tm.currentTabIndex--
	} else if tm.config.WrapNavigation {
		tm.currentTabIndex = len(tm.config.Tabs) - 1
	}

	if oldIndex != tm.currentTabIndex {
		tm.currentItemIndex = 0 // タブ変更時はアイテムインデックスをリセット
		if tm.callbacks.OnTabChange != nil {
			tm.callbacks.OnTabChange(oldIndex, tm.currentTabIndex, tm.config.Tabs[tm.currentTabIndex])
		}
		tm.notifyItemChange(0, 0) // 新しいタブの最初のアイテムを通知
	}
}

// navigateToNextTab は次のタブに移動する
func (tm *TabMenu) navigateToNextTab() {
	oldIndex := tm.currentTabIndex

	if tm.currentTabIndex < len(tm.config.Tabs)-1 {
		tm.currentTabIndex++
	} else if tm.config.WrapNavigation {
		tm.currentTabIndex = 0
	}

	if oldIndex != tm.currentTabIndex {
		tm.currentItemIndex = 0 // タブ変更時はアイテムインデックスをリセット
		if tm.callbacks.OnTabChange != nil {
			tm.callbacks.OnTabChange(oldIndex, tm.currentTabIndex, tm.config.Tabs[tm.currentTabIndex])
		}
		tm.notifyItemChange(0, 0) // 新しいタブの最初のアイテムを通知
	}
}

// navigateToPreviousItem は前のアイテムに移動する
func (tm *TabMenu) navigateToPreviousItem() {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 {
		return
	}

	oldIndex := tm.currentItemIndex

	if tm.currentItemIndex > 0 {
		tm.currentItemIndex--
	} else if tm.config.WrapNavigation {
		tm.currentItemIndex = len(currentTab.Items) - 1
	}

	if oldIndex != tm.currentItemIndex {
		tm.notifyItemChange(oldIndex, tm.currentItemIndex)
	}
}

// navigateToNextItem は次のアイテムに移動する
func (tm *TabMenu) navigateToNextItem() {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 {
		return
	}

	oldIndex := tm.currentItemIndex

	if tm.currentItemIndex < len(currentTab.Items)-1 {
		tm.currentItemIndex++
	} else if tm.config.WrapNavigation {
		tm.currentItemIndex = 0
	}

	if oldIndex != tm.currentItemIndex {
		tm.notifyItemChange(oldIndex, tm.currentItemIndex)
	}
}

// selectCurrentItem は現在のアイテムを選択する
func (tm *TabMenu) selectCurrentItem() {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 || tm.currentItemIndex >= len(currentTab.Items) {
		return
	}

	currentItem := currentTab.Items[tm.currentItemIndex]

	if tm.callbacks.OnSelectItem != nil {
		tm.callbacks.OnSelectItem(tm.currentTabIndex, tm.currentItemIndex, currentTab, currentItem)
	}
}

// notifyItemChange はアイテム変更を通知する
func (tm *TabMenu) notifyItemChange(oldIndex, newIndex int) {
	if len(tm.config.Tabs) == 0 || tm.currentTabIndex >= len(tm.config.Tabs) {
		return
	}

	currentTab := tm.config.Tabs[tm.currentTabIndex]
	if len(currentTab.Items) == 0 || newIndex >= len(currentTab.Items) {
		return
	}

	if tm.callbacks.OnItemChange != nil {
		tm.callbacks.OnItemChange(tm.currentTabIndex, oldIndex, newIndex, currentTab.Items[newIndex])
	}
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
func (tm *TabMenu) GetCurrentItem() menu.MenuItem {
	currentTab := tm.GetCurrentTab()
	if len(currentTab.Items) == 0 || tm.currentItemIndex >= len(currentTab.Items) {
		return menu.MenuItem{}
	}
	return currentTab.Items[tm.currentItemIndex]
}

// SetTabIndex はタブインデックスを設定する
func (tm *TabMenu) SetTabIndex(index int) {
	if index >= 0 && index < len(tm.config.Tabs) {
		oldIndex := tm.currentTabIndex
		tm.currentTabIndex = index
		tm.currentItemIndex = 0 // タブ変更時はアイテムインデックスをリセット

		if oldIndex != tm.currentTabIndex {
			if tm.callbacks.OnTabChange != nil {
				tm.callbacks.OnTabChange(oldIndex, tm.currentTabIndex, tm.config.Tabs[tm.currentTabIndex])
			}
			tm.notifyItemChange(0, 0)
		}
	}
}

// SetItemIndex はアイテムインデックスを設定する
func (tm *TabMenu) SetItemIndex(index int) {
	currentTab := tm.GetCurrentTab()
	if index >= 0 && index < len(currentTab.Items) {
		oldIndex := tm.currentItemIndex
		tm.currentItemIndex = index

		if oldIndex != tm.currentItemIndex {
			tm.notifyItemChange(oldIndex, tm.currentItemIndex)
		}
	}
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
		if tm.currentItemIndex >= len(currentTab.Items) {
			tm.currentItemIndex = len(currentTab.Items) - 1
			if tm.currentItemIndex < 0 {
				tm.currentItemIndex = 0
			}
		}
	}
}