package tabmenu

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
)

func TestTabMenuNavigation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{"タブ切り替え（Tab/Shift+Tab）", testTabSwitching},
		{"アイテム選択（上下矢印キー）", testItemNavigation},
		{"循環ナビゲーション", testWrapNavigation},
		{"選択機能", testSelection},
		{"キャンセル機能", testCancel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.test(t)
		})
	}
}

func testTabSwitching(t *testing.T) {
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
	tabMenu.Update()
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
	tabMenu.Update()
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("Shift+Tab後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
}

func testTabSwitchingWithTabKey(t *testing.T) {
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
	tabMenu.Update()
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
	tabMenu.Update()
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
	tabMenu.Update()
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 2 {
		t.Errorf("Shift+Tab循環後のタブインデックスが不正: 期待 2, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// 最後のタブでTab → 最初のタブに循環
	mockInput.SetKeyPressed(ebiten.KeyShift, false) // Shiftキーを離す
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	tabMenu.Update()
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("Tab循環後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
}

func testItemNavigation(t *testing.T) {
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
		OnItemChange: func(_ int, _, _ int, _ menu.Item) {
			itemChangeCount++
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
	tabMenu.Update()
	mockInput.Reset()

	if tabMenu.GetCurrentItemIndex() != 1 {
		t.Errorf("下矢印後のアイテムインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentItemIndex())
	}
	if itemChangeCount != 1 {
		t.Errorf("アイテム変更コールバック回数が不正: 期待 1, 実際 %d", itemChangeCount)
	}

	// 上矢印でアイテム1に戻る
	mockInput.SetKeyJustPressed(ebiten.KeyArrowUp, true)
	tabMenu.Update()
	mockInput.Reset()

	if tabMenu.GetCurrentItemIndex() != 0 {
		t.Errorf("上矢印後のアイテムインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentItemIndex())
	}
}

func testWrapNavigation(t *testing.T) {
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
	tabMenu.Update()
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("Shift+Tab循環後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// 最後のタブでTab → 最初のタブに循環
	mockInput.SetKeyPressed(ebiten.KeyShift, false) // Shiftキーを離す
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	tabMenu.Update()
	mockInput.Reset()

	if tabMenu.GetCurrentTabIndex() != 0 {
		t.Errorf("Tab循環後のタブインデックスが不正: 期待 0, 実際 %d", tabMenu.GetCurrentTabIndex())
	}
}

func testSelection(t *testing.T) {
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
		OnSelectItem: func(_, _ int, _ TabItem, item menu.Item) {
			selectedItem = item
		},
	}

	mockInput := input.NewMockKeyboardInput()
	tabMenu := NewTabMenu(config, callbacks, mockInput)

	// Enterキーで選択（セッションベース）
	mockInput.SimulateEnterPressRelease()
	tabMenu.Update()

	if selectedItem.ID != "item1" {
		t.Errorf("選択されたアイテムが不正: 期待 item1, 実際 %s", selectedItem.ID)
	}
}

func testCancel(t *testing.T) {
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
	tabMenu.Update()

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
	tabMenu.SetTabIndex(1)
	if tabMenu.GetCurrentTabIndex() != 1 {
		t.Errorf("設定後のタブインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentTabIndex())
	}

	// アイテムインデックスの設定
	tabMenu.SetTabIndex(0) // タブ1に戻す
	tabMenu.SetItemIndex(1)
	if tabMenu.GetCurrentItemIndex() != 1 {
		t.Errorf("設定後のアイテムインデックスが不正: 期待 1, 実際 %d", tabMenu.GetCurrentItemIndex())
	}
}
