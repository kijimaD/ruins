package tabmenu

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/stretchr/testify/require"
)

func TestTabSwitching(t *testing.T) {
	t.Parallel()
	// テスト用のタブとアイテムを作成
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

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, callbacks, mockInput)

	// 初期状態の確認
	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("初期タブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// Tabキーでタブ2に移動
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	_, err := tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("Tab後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
	if tabChangeCount != 1 {
		t.Errorf("タブ変更コールバック回数が不正: 期待 1, 実際 %d", tabChangeCount)
	}

	// Shift+Tabでタブ1に戻る
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	mockInput.SetKeyPressed(ebiten.KeyShift, true)
	_, err = tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("Shift+Tab後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
}

func TestTabSwitchingWithTabKey(t *testing.T) {
	t.Parallel()
	// テスト用のタブとアイテムを作成
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

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, callbacks, mockInput)

	// 初期状態の確認
	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("初期タブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// Tabキーでタブ2に移動
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	_, err := tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("Tab後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
	if tabChangeCount != 1 {
		t.Errorf("タブ変更コールバック回数が不正: 期待 1, 実際 %d", tabChangeCount)
	}

	// Shift+Tabでタブ1に戻る（前回のTabキーをリセット）
	mockInput.SetKeyPressed(ebiten.KeyShift, true)
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	_, err = tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("Shift+Tab後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
	if tabChangeCount != 2 {
		t.Errorf("タブ変更コールバック回数が不正: 期待 2, 実際 %d", tabChangeCount)
	}

	// 最初のタブでShift+Tab → 最後のタブに循環
	mockInput.SetKeyPressed(ebiten.KeyShift, true)
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	_, err = tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 2 {
		t.Errorf("Shift+Tab循環後のタブインデックスが不正: 期待 2, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// 最後のタブでTab → 最初のタブに循環
	mockInput.SetKeyPressed(ebiten.KeyShift, false) // Shiftキーを離す
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	_, err = tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("Tab循環後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
}

func TestItemNavigation(t *testing.T) {
	t.Parallel()
	// 複数アイテムを持つタブを作成
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

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, callbacks, mockInput)

	// 初期状態の確認
	if tabMenu.GetCurrentItemIndex() != 0 {
		t.Errorf("初期アイテムインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentItemIndex())
	}

	// 下矢印でアイテム2に移動
	mockInput.SetKeyJustPressed(ebiten.KeyArrowDown, true)
	_, err := tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentItemIndex() != 1 {
		t.Errorf("下矢印後のアイテムインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentItemIndex())
	}
	if itemChangeCount != 1 {
		t.Errorf("アイテム変更コールバック回数が不正: 期待 1, 実際 %d", itemChangeCount)
	}

	// 上矢印でアイテム1に戻る
	mockInput.SetKeyJustPressed(ebiten.KeyArrowUp, true)
	_, err = tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentItemIndex() != 0 {
		t.Errorf("上矢印後のアイテムインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentItemIndex())
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

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, Callbacks{}, mockInput)

	// 最初のタブでShift+Tab → 最後のタブに循環
	mockInput.SetKeyPressed(ebiten.KeyShift, true)
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	_, err := tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("Shift+Tab循環後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// 最後のタブでTab → 最初のタブに循環
	mockInput.SetKeyPressed(ebiten.KeyShift, false) // Shiftキーを離す
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	_, err = tabMenu.Update()
	require.NoError(t, err)
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("Tab循環後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
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

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, callbacks, mockInput)

	// Enterキーで選択（セッションベース）
	mockInput.SimulateEnterPressRelease()
	_, err := tabMenu.Update()
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

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, callbacks, mockInput)

	// Escapeキーでキャンセル
	mockInput.SetKeyJustPressed(ebiten.KeyEscape, true)
	_, err := tabMenu.Update()
	require.NoError(t, err)

	if !cancelCalled {
		t.Error("キャンセルコールバックが呼ばれていない")
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
		{ID: "tab2", Label: "タブ2", Items: []menu.Item{{ID: "item3", Label: "アイテム3"}}},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 1,
	}

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, Callbacks{}, mockInput)

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
		{
			ID:    "tab1",
			Label: "タブ1",
			Items: []menu.Item{
				{ID: "item1", Label: "アイテム1"},
				{ID: "item2", Label: "アイテム2"},
			},
		},
		{ID: "tab2", Label: "タブ2", Items: []menu.Item{{ID: "item3", Label: "アイテム3"}}},
	}

	config := Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
	}

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, Callbacks{}, mockInput)

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
