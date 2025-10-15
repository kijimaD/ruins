package tabmenu

import (
	"testing"

	"github.com/kijimaD/ruins/lib/inputmapper"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/stretchr/testify/require"
)

func TestTabSwitching(t *testing.T) {
	t.Parallel()
	tabs := []TabItem{
		{ID: "tab1", Label: "タブ1", Items: []menu.Item{{ID: "item1", Label: "アイテム1"}}},
		{ID: "tab2", Label: "タブ2", Items: []menu.Item{{ID: "item2", Label: "アイテム2"}}},
		{ID: "tab3", Label: "タブ3", Items: []menu.Item{{ID: "item3", Label: "アイテム3"}}},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
		WrapNavigation:   true,
	}

	tabChangeCount := 0
	callbacks := Callbacks{
		OnTabChange: func(_, _ int, _ TabItem) {
			tabChangeCount++
		},
	}

	tabMenu := NewTabMenu(config, callbacks)

	// 初期状態の確認
	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("初期タブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// ActionMenuRightでタブ2に移動
	err := tabMenu.DoAction(inputmapper.ActionMenuRight)
	require.NoError(t, err)

	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("ActionMenuRight後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
	if tabChangeCount != 1 {
		t.Errorf("タブ変更コールバック回数が不正: 期待 1, 実際 %d", tabChangeCount)
	}

	// ActionMenuLeftでタブ1に戻る
	err = tabMenu.DoAction(inputmapper.ActionMenuLeft)
	require.NoError(t, err)

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("ActionMenuLeft後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
	if tabChangeCount != 2 {
		t.Errorf("タブ変更コールバック回数が不正: 期待 2, 実際 %d", tabChangeCount)
	}
}

func TestItemNavigation(t *testing.T) {
	t.Parallel()
	tabs := []TabItem{
		{
			ID:    "tab1",
			Label: "タブ1",
			Items: []menu.Item{
				{ID: "item1", Label: "アイテム1"},
				{ID: "item2", Label: "アイテム2"},
				{ID: "item3", Label: "アイテム3"},
			},
		},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
		WrapNavigation:   true,
	}

	itemChangeCount := 0
	callbacks := Callbacks{
		OnItemChange: func(_ int, _, _ int, _ menu.Item) error {
			itemChangeCount++
			return nil
		},
	}

	tabMenu := NewTabMenu(config, callbacks)

	// 初期状態の確認
	if tabMenu.GetCurrentItemIndex() != 0 {
		t.Errorf("初期アイテムインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentItemIndex())
	}

	// ActionMenuDownでアイテム2に移動
	err := tabMenu.DoAction(inputmapper.ActionMenuDown)
	require.NoError(t, err)

	if tabMenu.GetCurrentItemIndex() != 1 {
		t.Errorf("ActionMenuDown後のアイテムインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentItemIndex())
	}
	if itemChangeCount != 1 {
		t.Errorf("アイテム変更コールバック回数が不正: 期待 1, 実際 %d", itemChangeCount)
	}

	// ActionMenuUpでアイテム1に戻る
	err = tabMenu.DoAction(inputmapper.ActionMenuUp)
	require.NoError(t, err)

	if tabMenu.GetCurrentItemIndex() != 0 {
		t.Errorf("ActionMenuUp後のアイテムインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentItemIndex())
	}
}

func TestWrapNavigation(t *testing.T) {
	t.Parallel()
	tabs := []TabItem{
		{ID: "tab1", Label: "タブ1", Items: []menu.Item{{ID: "item1", Label: "アイテム1"}}},
		{ID: "tab2", Label: "タブ2", Items: []menu.Item{{ID: "item2", Label: "アイテム2"}}},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
		WrapNavigation:   true,
	}

	tabMenu := NewTabMenu(config, Callbacks{})

	// 最初のタブでActionMenuLeft → 最後のタブに循環
	err := tabMenu.DoAction(inputmapper.ActionMenuLeft)
	require.NoError(t, err)

	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("ActionMenuLeft循環後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// 最後のタブでActionMenuRight → 最初のタブに循環
	err = tabMenu.DoAction(inputmapper.ActionMenuRight)
	require.NoError(t, err)

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("ActionMenuRight循環後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
}

func TestSelection(t *testing.T) {
	t.Parallel()
	tabs := []TabItem{
		{
			ID:    "tab1",
			Label: "タブ1",
			Items: []menu.Item{
				{ID: "item1", Label: "アイテム1", UserData: "data1"},
			},
		},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
	}

	var selectedItem menu.Item
	callbacks := Callbacks{
		OnSelectItem: func(_, _ int, _ TabItem, item menu.Item) error {
			selectedItem = item
			return nil
		},
	}

	tabMenu := NewTabMenu(config, callbacks)

	// ActionMenuSelectで選択
	err := tabMenu.DoAction(inputmapper.ActionMenuSelect)
	require.NoError(t, err)

	if selectedItem.ID != "item1" {
		t.Errorf("選択されたアイテムが不正: 期待 item1, 実際 %s", selectedItem.ID)
	}
}

func TestCancel(t *testing.T) {
	t.Parallel()
	tabs := []TabItem{
		{ID: "tab1", Label: "タブ1", Items: []menu.Item{{ID: "item1", Label: "アイテム1"}}},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
	}

	cancelCalled := false
	callbacks := Callbacks{
		OnCancel: func() {
			cancelCalled = true
		},
	}

	tabMenu := NewTabMenu(config, callbacks)

	// ActionMenuCancelでキャンセル
	err := tabMenu.DoAction(inputmapper.ActionMenuCancel)
	require.NoError(t, err)

	if !cancelCalled {
		t.Error("OnCancelが呼ばれていない")
	}
}

func TestTabMenuGetters(t *testing.T) {
	t.Parallel()
	tabs := []TabItem{
		{
			ID:    "tab1",
			Label: "タブ1",
			Items: []menu.Item{
				{ID: "item1", Label: "アイテム1"},
				{ID: "item2", Label: "アイテム2"},
			},
		},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 1,
	}

	tabMenu := NewTabMenu(config, Callbacks{})

	// 現在のタブとアイテムの確認
	currentTab := tabMenu.GetCurrentTab()
	if currentTab.ID != "tab1" {
		t.Errorf("現在のタブが不正: 期待 tab1, 実際 %s", currentTab.ID)
	}

	currentItem := tabMenu.GetCurrentItem()
	if currentItem.ID != "item2" {
		t.Errorf("現在のアイテムが不正: 期待 item2, 実際 %s", currentItem.ID)
	}
}

func TestTabMenuSetters(t *testing.T) {
	t.Parallel()
	tabs := []TabItem{
		{ID: "tab1", Label: "タブ1", Items: []menu.Item{{ID: "item1", Label: "アイテム1"}, {ID: "item2", Label: "アイテム2"}}},
		{ID: "tab2", Label: "タブ2", Items: []menu.Item{{ID: "item3", Label: "アイテム3"}}},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
	}

	tabMenu := NewTabMenu(config, Callbacks{})

	// タブインデックスの設定
	require.NoError(t, tabMenu.SetTabIndex(1))
	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("設定後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// アイテムインデックスの設定
	require.NoError(t, tabMenu.SetTabIndex(0)) // タブ1に戻す
	require.NoError(t, tabMenu.SetItemIndex(1))
	if tabMenu.GetCurrentItemIndex() != 1 {
		t.Errorf("設定後のアイテムインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentItemIndex())
	}
}
