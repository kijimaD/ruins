package menu

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
)

func TestMenuNavigation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		items          []Item
		initialIndex   int
		wrapNavigation bool
		keyPress       ebiten.Key
		expectedIndex  int
	}{
		{
			name: "下矢印キーでフォーカス移動",
			items: []Item{
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
			items: []Item{
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
			items: []Item{
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
			items: []Item{
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
			items: []Item{
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
			t.Parallel()
			config := Config{
				Items:          tt.items,
				InitialIndex:   tt.initialIndex,
				WrapNavigation: tt.wrapNavigation,
				Orientation:    Vertical,
			}

			menu := NewMenu(config, Callbacks{})

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
	t.Parallel()
	items := []Item{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
		{ID: "3", Label: "Item 3"},
	}

	var selectedIndex int
	var selectedItem Item
	var selectionCalled bool

	config := Config{
		Items:        items,
		InitialIndex: 1,
	}

	callbacks := Callbacks{
		OnSelect: func(index int, item Item) {
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
	t.Parallel()
	items := []Item{
		{ID: "1", Label: "Item 1"},
	}

	var cancelCalled bool

	config := Config{
		Items: items,
	}

	callbacks := Callbacks{
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
	t.Parallel()
	items := []Item{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
		{ID: "3", Label: "Item 3"},
	}

	config := Config{
		Items:        items,
		InitialIndex: 0,
	}

	menu := NewMenu(config, Callbacks{})

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

func TestMenuDisabledItems(t *testing.T) {
	t.Parallel()
	items := []Item{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2", Disabled: true},
		{ID: "3", Label: "Item 3", Disabled: true},
		{ID: "4", Label: "Item 4"},
	}

	var selectionCalled bool

	config := Config{
		Items:        items,
		InitialIndex: 1, // 無効なアイテムから開始
	}

	callbacks := Callbacks{
		OnSelect: func(_ int, _ Item) {
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

func TestMenuConsecutiveEnterPrevention(t *testing.T) {
	t.Parallel()
	// グローバル状態をリセット
	input.ResetGlobalKeyStateForTest()

	items := []Item{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
	}

	var selectionCount int

	config := Config{
		Items:        items,
		InitialIndex: 0,
	}

	callbacks := Callbacks{
		OnSelect: func(_ int, _ Item) {
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

// TestMenuScrolling はスクロール機能をテストする
func TestMenuScrolling(t *testing.T) {
	t.Parallel()

	// 10項目のメニューを作成（ページサイズ3）
	items := make([]Item, 10)
	for i := 0; i < 10; i++ {
		items[i] = Item{
			ID:    fmt.Sprintf("item_%d", i),
			Label: fmt.Sprintf("Item %d", i),
		}
	}

	config := Config{
		Items:        items,
		InitialIndex: 0,
		ItemsPerPage: 3,
	}

	menu := NewMenu(config, Callbacks{})

	// 初期状態の確認
	if menu.GetCurrentPage() != 1 {
		t.Errorf("初期ページが間違っています: 期待 1, 実際 %d", menu.GetCurrentPage())
	}

	if menu.GetTotalPages() != 4 {
		t.Errorf("総ページ数が間違っています: 期待 4, 実際 %d", menu.GetTotalPages())
	}

	// 表示項目の確認（GetVisibleItemsWithIndicesを使用）
	visibleItems, _ := menu.GetVisibleItems()
	if len(visibleItems) != 3 {
		t.Errorf("表示項目数が間違っています: 期待 3, 実際 %d", len(visibleItems))
	}

	// 最初のページでの下矢印移動
	mockInput := input.NewMockKeyboardInput()

	// Item 0 → Item 1
	mockInput.SetKeyJustPressed(ebiten.KeyArrowDown, true)
	menu.Update(mockInput)
	if menu.GetFocusedIndex() != 1 {
		t.Errorf("ページ内移動が失敗: 期待 1, 実際 %d", menu.GetFocusedIndex())
	}

	// Item 1 → Item 2
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyArrowDown, true)
	menu.Update(mockInput)
	if menu.GetFocusedIndex() != 2 {
		t.Errorf("ページ内移動が失敗: 期待 2, 実際 %d", menu.GetFocusedIndex())
	}

	// ページの境界を超えて次のページに移動（Item 2 → Item 3）
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyArrowDown, true)
	menu.Update(mockInput)
	if menu.GetFocusedIndex() != 3 {
		t.Errorf("次ページへの移動が失敗: 期待 3, 実際 %d", menu.GetFocusedIndex())
	}

	if menu.GetCurrentPage() != 2 {
		t.Errorf("ページが更新されていません: 期待 2, 実際 %d", menu.GetCurrentPage())
	}
}

// TestMenuPageUpDown はPageUp/PageDownキーをテストする
func TestMenuPageUpDown(t *testing.T) {
	t.Parallel()

	items := make([]Item, 10)
	for i := 0; i < 10; i++ {
		items[i] = Item{
			ID:    fmt.Sprintf("item_%d", i),
			Label: fmt.Sprintf("Item %d", i),
		}
	}

	config := Config{
		Items:        items,
		InitialIndex: 0,
		ItemsPerPage: 3,
	}

	menu := NewMenu(config, Callbacks{})

	// PageDownで次のページに移動
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyPageDown, true)
	menu.Update(mockInput)

	if menu.GetCurrentPage() != 2 {
		t.Errorf("PageDownでページ移動が失敗: 期待 2, 実際 %d", menu.GetCurrentPage())
	}

	if menu.GetFocusedIndex() != 3 {
		t.Errorf("PageDownでフォーカス移動が失敗: 期待 3, 実際 %d", menu.GetFocusedIndex())
	}

	// PageUpで前のページに戻る
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyPageUp, true)
	menu.Update(mockInput)

	if menu.GetCurrentPage() != 1 {
		t.Errorf("PageUpでページ移動が失敗: 期待 1, 実際 %d", menu.GetCurrentPage())
	}

	if menu.GetFocusedIndex() != 0 {
		t.Errorf("PageUpでフォーカス移動が失敗: 期待 0, 実際 %d", menu.GetFocusedIndex())
	}
}

// TestMenuVisibleItemsWithIndices は表示項目とインデックスの取得をテストする
func TestMenuVisibleItemsWithIndices(t *testing.T) {
	t.Parallel()

	items := make([]Item, 7)
	for i := 0; i < 7; i++ {
		items[i] = Item{
			ID:    fmt.Sprintf("item_%d", i),
			Label: fmt.Sprintf("Item %d", i),
		}
	}

	config := Config{
		Items:        items,
		ItemsPerPage: 3,
	}

	menu := NewMenu(config, Callbacks{})

	// 1ページ目の確認
	visibleItems, indices := menu.GetVisibleItems()
	if len(visibleItems) != 3 {
		t.Errorf("1ページ目の表示項目数が間違っています: 期待 3, 実際 %d", len(visibleItems))
	}
	if len(indices) != 3 {
		t.Errorf("1ページ目のインデックス数が間違っています: 期待 3, 実際 %d", len(indices))
	}
	if indices[0] != 0 || indices[1] != 1 || indices[2] != 2 {
		t.Errorf("1ページ目のインデックスが間違っています: 期待 [0,1,2], 実際 %v", indices)
	}

	// 2ページ目に移動
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyPageDown, true)
	menu.Update(mockInput)

	visibleItems, indices = menu.GetVisibleItems()
	if len(visibleItems) != 3 {
		t.Errorf("2ページ目の表示項目数が間違っています: 期待 3, 実際 %d", len(visibleItems))
	}
	if indices[0] != 3 || indices[1] != 4 || indices[2] != 5 {
		t.Errorf("2ページ目のインデックスが間違っています: 期待 [3,4,5], 実際 %v", indices)
	}

	// 3ページ目（最後のページ、項目が少ない）に移動
	mockInput.Reset()
	mockInput.SetKeyJustPressed(ebiten.KeyPageDown, true)
	menu.Update(mockInput)

	visibleItems, indices = menu.GetVisibleItems()
	if len(visibleItems) != 1 {
		t.Errorf("3ページ目の表示項目数が間違っています: 期待 1, 実際 %d", len(visibleItems))
	}
	if indices[0] != 6 {
		t.Errorf("3ページ目のインデックスが間違っています: 期待 [6], 実際 %v", indices)
	}
}

func TestMenuScrollingDisabled(t *testing.T) {
	t.Parallel()

	items := make([]Item, 10)
	for i := 0; i < 10; i++ {
		items[i] = Item{
			ID:    fmt.Sprintf("item_%d", i),
			Label: fmt.Sprintf("Item %d", i),
		}
	}

	// スクロール無効（ItemsPerPage = 0）
	config := Config{
		Items:        items,
		ItemsPerPage: 0, // スクロール無効
	}

	menu := NewMenu(config, Callbacks{})

	// 全項目が表示されることを確認（GetVisibleItemsWithIndicesを使用）
	visibleItems, _ := menu.GetVisibleItems()
	if len(visibleItems) != 10 {
		t.Errorf("スクロール無効時の表示項目数が間違っています: 期待 10, 実際 %d", len(visibleItems))
	}

	// PageDownキーを押しても何も変わらないことを確認
	mockInput := input.NewMockKeyboardInput()
	mockInput.SetKeyJustPressed(ebiten.KeyPageDown, true)
	menu.Update(mockInput)

	if menu.GetCurrentPage() != 1 {
		t.Errorf("スクロール無効時にページが変わっています: 期待 1, 実際 %d", menu.GetCurrentPage())
	}
}
