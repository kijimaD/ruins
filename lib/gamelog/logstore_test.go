package gamelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatest(t *testing.T) {
	t.Parallel()
	t.Run("古い順に取得できる", func(t *testing.T) {
		t.Parallel()
		ss := SafeSlice{}
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
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")
		ss.Append("3")
		ss.Append("4")
		ss.Append("5")

		assert.Equal(t, []string{"1", "2", "3", "4", "5"}, ss.Pop())
	})
	t.Run("取得した分は消える", func(t *testing.T) {
		t.Parallel()
		ss := SafeSlice{}
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
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")

		assert.Equal(t, []string{"1", "2"}, ss.Get())
		ss.Flush()
		assert.Equal(t, []string{}, ss.Get())
	})
}
