package menu

import (
	"fmt"
	"testing"

	"github.com/kijimaD/ruins/lib/inputmapper"
	"github.com/stretchr/testify/require"
)

func TestMenuDoAction(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		items         []Item
		initialIndex  int
		action        inputmapper.ActionID
		expectedIndex int
	}{
		{
			name: "ActionMenuDownで次へ移動",
			items: []Item{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:  0,
			action:        inputmapper.ActionMenuDown,
			expectedIndex: 1,
		},
		{
			name: "ActionMenuUpで前へ移動",
			items: []Item{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:  1,
			action:        inputmapper.ActionMenuUp,
			expectedIndex: 0,
		},
		{
			name: "無効な項目をスキップ",
			items: []Item{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2", Disabled: true},
				{ID: "3", Label: "Item 3"},
			},
			initialIndex:  0,
			action:        inputmapper.ActionMenuDown,
			expectedIndex: 2, // Item 2はDisabledなのでスキップしてItem 3へ
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config := Config{
				Items:          tt.items,
				InitialIndex:   tt.initialIndex,
				WrapNavigation: true,
				Orientation:    Vertical,
			}

			menu := NewMenu(config, Callbacks{})

			// DoActionでメニュー操作
			require.NoError(t, menu.DoAction(tt.action))

			// 結果を確認
			if menu.GetFocusedIndex() != tt.expectedIndex {
				t.Errorf("expected focused index %d, got %d", tt.expectedIndex, menu.GetFocusedIndex())
			}
		})
	}
}

func TestMenuDoActionSelect(t *testing.T) {
	t.Parallel()
	items := []Item{
		{ID: "1", Label: "Item 1"},
		{ID: "2", Label: "Item 2"},
	}

	var selectedIndex int
	var selectedItem Item
	config := Config{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    Vertical,
	}

	callbacks := Callbacks{
		OnSelect: func(index int, item Item) error {
			selectedIndex = index
			selectedItem = item
			return nil
		},
	}

	menu := NewMenu(config, callbacks)

	// ActionMenuSelectで選択
	require.NoError(t, menu.DoAction(inputmapper.ActionMenuSelect))

	// 結果を確認
	if selectedIndex != 0 {
		t.Errorf("expected selected index 0, got %d", selectedIndex)
	}
	if selectedItem.ID != "1" {
		t.Errorf("expected selected item ID '1', got '%s'", selectedItem.ID)
	}
}

func TestMenuDoActionCancel(t *testing.T) {
	t.Parallel()
	items := []Item{
		{ID: "1", Label: "Item 1"},
	}

	var cancelCalled bool
	config := Config{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    Vertical,
	}

	callbacks := Callbacks{
		OnCancel: func() {
			cancelCalled = true
		},
	}

	menu := NewMenu(config, callbacks)

	// ActionMenuCancelでキャンセル
	require.NoError(t, menu.DoAction(inputmapper.ActionMenuCancel))

	// 結果を確認
	if !cancelCalled {
		t.Error("expected OnCancel to be called")
	}
}

func TestMenuWrapNavigation(t *testing.T) {
	t.Parallel()

	t.Run("循環ナビゲーション有効", func(t *testing.T) {
		t.Parallel()
		config := Config{
			Items: []Item{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			InitialIndex:   2,
			WrapNavigation: true,
			Orientation:    Vertical,
		}

		menu := NewMenu(config, Callbacks{})

		// 最後の項目から下矢印で最初に戻る
		require.NoError(t, menu.DoAction(inputmapper.ActionMenuDown))

		if menu.GetFocusedIndex() != 0 {
			t.Errorf("expected wrap to first item, got index %d", menu.GetFocusedIndex())
		}
	})

	t.Run("循環ナビゲーション無効", func(t *testing.T) {
		t.Parallel()
		config := Config{
			Items: []Item{
				{ID: "1", Label: "Item 1"},
				{ID: "2", Label: "Item 2"},
				{ID: "3", Label: "Item 3"},
			},
			InitialIndex:   2,
			WrapNavigation: false,
			Orientation:    Vertical,
		}

		menu := NewMenu(config, Callbacks{})

		// 最後の項目から下矢印で停止
		require.NoError(t, menu.DoAction(inputmapper.ActionMenuDown))

		if menu.GetFocusedIndex() != 2 {
			t.Errorf("expected to stay at last item, got index %d", menu.GetFocusedIndex())
		}
	})
}

func TestMenuDisabledItems(t *testing.T) {
	t.Parallel()

	config := Config{
		Items: []Item{
			{ID: "1", Label: "Item 1"},
			{ID: "2", Label: "Item 2", Disabled: true},
			{ID: "3", Label: "Item 3"},
		},
		InitialIndex:   1, // 無効なアイテム
		WrapNavigation: true,
		Orientation:    Vertical,
	}

	menu := NewMenu(config, Callbacks{})

	// 初期フォーカスが無効なアイテムを避けて最初の有効なアイテムに移動することを確認
	if menu.GetFocusedIndex() != 0 {
		t.Errorf("expected initial focus to skip disabled item, got index %d", menu.GetFocusedIndex())
	}

	// 無効なアイテムで選択しても何も起こらないことを確認
	selectionCalled := false
	menu.callbacks.OnSelect = func(_ int, _ Item) error {
		selectionCalled = true
		return nil
	}

	menu.focusedIndex = 1 // 強制的に無効なアイテムにフォーカス
	require.NoError(t, menu.DoAction(inputmapper.ActionMenuSelect))

	if selectionCalled {
		t.Error("expected OnSelect not to be called for disabled item")
	}
}

func TestMenuPagination(t *testing.T) {
	t.Parallel()

	items := make([]Item, 10)
	for i := 0; i < 10; i++ {
		items[i] = Item{ID: fmt.Sprintf("%d", i), Label: fmt.Sprintf("Item %d", i)}
	}

	config := Config{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    Vertical,
		ItemsPerPage:   3,
	}

	menu := NewMenu(config, Callbacks{})

	// 初期ページ
	if menu.GetCurrentPage() != 1 {
		t.Errorf("expected current page 1, got %d", menu.GetCurrentPage())
	}

	if menu.GetTotalPages() != 4 {
		t.Errorf("expected total pages 4, got %d", menu.GetTotalPages())
	}

	// 表示可能な項目数の確認
	visibleItems, indices := menu.GetVisibleItems()
	if len(visibleItems) != 3 {
		t.Errorf("expected 3 visible items, got %d", len(visibleItems))
	}
	if indices[0] != 0 || indices[1] != 1 || indices[2] != 2 {
		t.Errorf("expected indices [0,1,2], got %v", indices)
	}

	// 複数回下矢印でページをまたぐ
	require.NoError(t, menu.DoAction(inputmapper.ActionMenuDown)) // 0 -> 1
	require.NoError(t, menu.DoAction(inputmapper.ActionMenuDown)) // 1 -> 2
	require.NoError(t, menu.DoAction(inputmapper.ActionMenuDown)) // 2 -> 3 (次のページ)

	if menu.GetFocusedIndex() != 3 {
		t.Errorf("expected focused index 3, got %d", menu.GetFocusedIndex())
	}

	if menu.GetCurrentPage() != 2 {
		t.Errorf("expected current page 2, got %d", menu.GetCurrentPage())
	}
}
