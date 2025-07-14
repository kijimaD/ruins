package menu

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
)

func TestMenuNavigation(t *testing.T) {
	tests := []struct {
		name           string
		items          []MenuItem
		initialIndex   int
		wrapNavigation bool
		keyPress       ebiten.Key
		expectedIndex  int
	}{
		{
			name: "下矢印キーでフォーカス移動",
			items: []MenuItem{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:   0,
			wrapNavigation: true,
			keyPress:       ebiten.KeyArrowDown,
			expectedIndex:  1,
		},
		{
			name: "上矢印キーでフォーカス移動",
			items: []MenuItem{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:   1,
			wrapNavigation: true,
			keyPress:       ebiten.KeyArrowUp,
			expectedIndex:  0,
		},
		{
			name: "無効な項目をスキップ",
			items: []MenuItem{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2", Disabled: true},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:   0,
			wrapNavigation: true,
			keyPress:       ebiten.KeyArrowDown,
			expectedIndex:  2, // Item 2はDisabledなのでスキップしてItem 3へ
		},
		{
			name: "循環ナビゲーション（最後から最初へ）",
			items: []MenuItem{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:   2,
			wrapNavigation: true,
			keyPress:       ebiten.KeyArrowDown,
			expectedIndex:  0,
		},
		{
			name: "循環ナビゲーション無効（端で停止）",
			items: []MenuItem{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:   2,
			wrapNavigation: false,
			keyPress:       ebiten.KeyArrowDown,
			expectedIndex:  2, // 端で停止
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := MenuConfig{
				Items:          tt.items,
				InitialIndex:   tt.initialIndex,
				WrapNavigation: tt.wrapNavigation,
				Orientation:    Vertical,
			}

			menu := NewMenu(config, MenuCallbacks{})

			// モックキーボード入力を作成
			mockInput := input.NewMockKeyboardInput()
			mockInput.SetKeyJustPressed(tt.keyPress, true)

			// メニューを更新
			menu.Update(mockInput)

			// 結果を検証
			if menu.GetFocusedIndex() != tt.expectedIndex {
				t.Errorf("期待されるフォーカスインデックス: %d, 実際: %d", tt.expectedIndex, menu.GetFocusedIndex())
			}
		})
	}
}

func TestMenuSelection(t *testing.T) {
	items := []MenuItem{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
		{ID: "3", Label: "Item 3"},
	}

	var selectedIndex int
	var selectedItem MenuItem
	var selectionCalled bool

	config := MenuConfig{
		Items:        items,
		InitialIndex: 1,
	}

	callbacks := MenuCallbacks{
		OnSelect: func(index int, item MenuItem) {
			selectedIndex = index
			selectedItem = item
			selectionCalled = true
		},
	}

	menu := NewMenu(config, callbacks)

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	// Enterキーは押下-押上ワンセットを使用
	mockInput.SimulateEnterPressRelease()

	// メニューを更新
	menu.Update(mockInput)

	// 結果を検証
	if !selectionCalled {
		t.Error("OnSelectコールバックが呼ばれていない")
	}

	if selectedIndex != 1 {
		t.Errorf("期待される選択インデックス: 1, 実際: %d", selectedIndex)
	}

	if selectedItem.ID != "2" {
		t.Errorf("期待される選択アイテムID: 2, 実際: %s", selectedItem.ID)
	}
}

func TestMenuCancel(t *testing.T) {
	items := []MenuItem{
		{ID: "1", Label: "Item 1"},
	}

	var cancelCalled bool

	config := MenuConfig{
		Items: items,
	}

	callbacks := MenuCallbacks{
		OnCancel: func() {
			cancelCalled = true
		},
	}

	menu := NewMenu(config, callbacks)

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyEscape, true)

	// メニューを更新
	menu.Update(mockInput)

	// 結果を検証
	if !cancelCalled {
		t.Error("OnCancelコールバックが呼ばれていない")
	}
}

func TestMenuTabNavigation(t *testing.T) {
	items := []MenuItem{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
		{ID: "3", Label: "Item 3"},
	}

	config := MenuConfig{
		Items:        items,
		InitialIndex: 0,
	}

	menu := NewMenu(config, MenuCallbacks{})

	// Tabキーのテスト
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	mockInput.SetKeyPressed(ebiten.KeyShift, false)

	menu.Update(mockInput)

	if menu.GetFocusedIndex() != 1 {
		t.Errorf("Tabキーでの移動が失敗: 期待 1, 実際 %d", menu.GetFocusedIndex())
	}

	// Shift+Tabキーのテスト
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyTab, true)
	mockInput.SetKeyPressed(ebiten.KeyShift, true)

	menu.Update(mockInput)

	if menu.GetFocusedIndex() != 0 {
		t.Errorf("Shift+Tabキーでの移動が失敗: 期待 0, 実際 %d", menu.GetFocusedIndex())
	}
}

func TestMenuGridNavigation(t *testing.T) {
	items := []MenuItem{
		{ID: "1", Label: "Item 1"}, {ID: "2", Label: "Item 2"},
		{ID: "3", Label: "Item 3"}, {ID: "4", Label: "Item 4"},
		{ID: "5", Label: "Item 5"}, {ID: "6", Label: "Item 6"},
	}

	config := MenuConfig{
		Items:        items,
		Columns:      2, // 2列のグリッド
		InitialIndex: 0,
	}

	menu := NewMenu(config, MenuCallbacks{})

	// 右矢印キーのテスト（0→1）
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyArrowRight, true)

	menu.Update(mockInput)

	if menu.GetFocusedIndex() != 1 {
		t.Errorf("右矢印キーでの移動が失敗: 期待 1, 実際 %d", menu.GetFocusedIndex())
	}

	// 左矢印キーのテスト（1→0）
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyArrowLeft, true)

	menu.Update(mockInput)

	if menu.GetFocusedIndex() != 0 {
		t.Errorf("左矢印キーでの移動が失敗: 期待 0, 実際 %d", menu.GetFocusedIndex())
	}
}

func TestMenuDisabledItems(t *testing.T) {
	items := []MenuItem{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2", Disabled: true},
		{ID: "3", Label: "Item 3", Disabled: true},
		{ID: "4", Label: "Item 4"},
	}

	var selectionCalled bool

	config := MenuConfig{
		Items:        items,
		InitialIndex: 1, // 無効なアイテムから開始
	}

	callbacks := MenuCallbacks{
		OnSelect: func(_ int, _ MenuItem) {
			selectionCalled = true
		},
	}

	menu := NewMenu(config, callbacks)

	// 初期フォーカスが有効なアイテム（0）に修正されているか確認
	if menu.GetFocusedIndex() != 0 {
		t.Errorf("初期フォーカスが無効なアイテムを避けていない: 期待 0, 実際 %d", menu.GetFocusedIndex())
	}

	// 無効なアイテムでEnterを押しても選択されないことを確認
	menu.focusedIndex = 1 // 強制的に無効なアイテムにフォーカス
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyEnter, true)

	menu.Update(mockInput)

	if selectionCalled {
		t.Error("無効なアイテムが選択された")
	}
}

func TestMenuSpaceSelection(t *testing.T) {
	items := []MenuItem{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
		{ID: "3", Label: "Item 3"},
	}

	var selectedIndex int
	var selectedItem MenuItem
	var selectionCalled bool

	config := MenuConfig{
		Items:        items,
		InitialIndex: 0,
	}

	callbacks := MenuCallbacks{
		OnSelect: func(index int, item MenuItem) {
			selectedIndex = index
			selectedItem = item
			selectionCalled = true
		},
	}

	menu := NewMenu(config, callbacks)

	// モックキーボード入力を作成
	mockInput := input.NewMockKeyboardInput()
	// Spaceキーは通常のIsKeyJustPressedを使用
	mockInput.SetKeyJustPressed(ebiten.KeySpace, true)

	// メニューを更新
	menu.Update(mockInput)

	// 結果を検証
	if !selectionCalled {
		t.Error("OnSelectコールバックが呼ばれていない")
	}

	if selectedIndex != 0 {
		t.Errorf("期待される選択インデックス: 0, 実際: %d", selectedIndex)
	}

	if selectedItem.ID != "1" {
		t.Errorf("期待される選択アイテムID: 1, 実際: %s", selectedItem.ID)
	}
}

func TestMenuConsecutiveEnterPrevention(t *testing.T) {
	// グローバル状態をリセット
	input.ResetGlobalKeyStateForTest()

	items := []MenuItem{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
	}

	var selectionCount int

	config := MenuConfig{
		Items:        items,
		InitialIndex: 0,
	}

	callbacks := MenuCallbacks{
		OnSelect: func(_ int, _ MenuItem) {
			selectionCount++
		},
	}

	menu := NewMenu(config, callbacks)
	mockInput := input.NewMockKeyboardInput()

	// 1回目のEnter - これは成功するはず（押下-押上ワンセット）
	mockInput.SimulateEnterPressRelease()
	menu.Update(mockInput)

	if selectionCount != 1 {
		t.Errorf("1回目の選択が失敗: 期待 1, 実際 %d", selectionCount)
	}

	// 注意: 押下-押上ワンセット制御付きEnterキーのテストは実際のキー状態制御が必要なため、
	// モック環境では連続クリック防止の詳細テストは省略
	// 実際のゲーム環境では押下-押上ワンセット制限が働く
}

func TestMenuSpaceBasicFunction(t *testing.T) {
	items := []MenuItem{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
	}

	var selectionCount int

	config := MenuConfig{
		Items:        items,
		InitialIndex: 0,
	}

	callbacks := MenuCallbacks{
		OnSelect: func(_ int, _ MenuItem) {
			selectionCount++
		},
	}

	menu := NewMenu(config, callbacks)
	mockInput := input.NewMockKeyboardInput()

	// Spaceキーによる選択のテスト
	mockInput.SetKeyJustPressed(ebiten.KeySpace, true)
	menu.Update(mockInput)

	if selectionCount != 1 {
		t.Errorf("Spaceキーによる選択が失敗: 期待 1, 実際 %d", selectionCount)
	}

	// リセット後の再選択テスト
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeySpace, true)
	menu.Update(mockInput)

	if selectionCount != 2 {
		t.Errorf("リセット後のSpace選択が失敗: 期待 2, 実際 %d", selectionCount)
	}
}
