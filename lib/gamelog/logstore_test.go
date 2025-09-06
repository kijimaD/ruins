package gamelog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatest(t *testing.T) {
	t.Parallel()
	t.Run("古い順に取得できる", func(t *testing.T) {
		t.Parallel()
		ss := NewSafeSlice(10) // 新しいコンストラクタを使用
		ss.Append("1")
		ss.Append("2")
		ss.Append("3")

		assert.Equal(t, []string{"1", "2", "3"}, ss.Get())
	})
}

func TestPop(t *testing.T) {
	t.Parallel()
	t.Run("取得できる", func(t *testing.T) {
		t.Parallel()
		ss := NewSafeSlice(10)
		ss.Append("1")
		ss.Append("2")
		ss.Append("3")
		ss.Append("4")
		ss.Append("5")

		assert.Equal(t, []string{"1", "2", "3", "4", "5"}, ss.Pop())
	})
	t.Run("取得した分は消える", func(t *testing.T) {
		t.Parallel()
		ss := NewSafeSlice(10)
		ss.Append("1")
		ss.Append("2")
		ss.Append("3")
		ss.Append("4")
		ss.Append("5")

		ss.Pop()
		assert.Equal(t, 0, len(ss.Get()))
	})
}

func TestFlush(t *testing.T) {
	t.Parallel()
	t.Run("リセットできる", func(t *testing.T) {
		t.Parallel()
		ss := NewSafeSlice(10)
		ss.Append("1")
		ss.Append("2")

		assert.Equal(t, []string{"1", "2"}, ss.Get())
		ss.Flush()
		assert.Equal(t, []string{}, ss.Get())
	})
}

// 新しい機能のテスト
func TestSafeSliceMaxSize(t *testing.T) {
	t.Parallel()
	t.Run("最大サイズを超えた場合、古い要素が削除される", func(t *testing.T) {
		t.Parallel()
		// 最大サイズ3のSafeSliceを作成
		sl := NewSafeSlice(3)

		// 3つの要素を追加
		sl.Append("message1")
		sl.Append("message2")
		sl.Append("message3")

		content := sl.Get()
		assert.Equal(t, 3, len(content))
		assert.Equal(t, []string{"message1", "message2", "message3"}, content)

		// 4つ目を追加（最古の要素が削除されるはず）
		sl.Append("message4")

		content = sl.Get()
		assert.Equal(t, 3, len(content))
		// 最古の要素（message1）が削除され、message2, message3, message4が残るはず
		assert.Equal(t, []string{"message2", "message3", "message4"}, content)
	})

	t.Run("0サイズ指定時はデフォルトサイズが使用される", func(t *testing.T) {
		t.Parallel()
		sl := NewSafeSlice(0)

		// デフォルトサイズまで追加
		for i := 0; i < DefaultMaxLogSize+10; i++ {
			sl.Append(fmt.Sprintf("message%d", i))
		}

		content := sl.Get()
		assert.Equal(t, DefaultMaxLogSize, len(content))
	})

	t.Run("負の値指定時はデフォルトサイズが使用される", func(t *testing.T) {
		t.Parallel()
		sl := NewSafeSlice(-5)

		// 少し多めに追加してテスト
		for i := 0; i < 20; i++ {
			sl.Append(fmt.Sprintf("message%d", i))
		}

		content := sl.Get()
		assert.Equal(t, 20, len(content)) // デフォルトサイズ以下なので全部残る
	})
}

func TestSafeSliceMemoryLeak(t *testing.T) {
	t.Parallel()
	t.Run("大量の要素追加でもメモリリークしない", func(t *testing.T) {
		t.Parallel()
		sl := NewSafeSlice(10)

		// 大量の要素を追加
		for i := 0; i < 1000; i++ {
			sl.Append(fmt.Sprintf("message%d", i))
		}

		content := sl.Get()
		// 最大サイズを超えないことを確認
		assert.Equal(t, 10, len(content))

		// 最新の10個が保持されていることを確認
		for i := 0; i < 10; i++ {
			expected := fmt.Sprintf("message%d", 990+i)
			assert.Equal(t, expected, content[i])
		}
	})
}
