package gamelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatest(t *testing.T) {
	t.Run("数を指定して新しい順に取得できる", func(t *testing.T) {
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")

		assert.Equal(t, []string{"2"}, ss.Latest(1))
	})
	t.Run("長さを超えて指定するとある分だけ返す", func(t *testing.T) {
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")

		assert.Equal(t, []string{"1", "2"}, ss.Latest(5))
	})
}

func TestPop(t *testing.T) {
	t.Run("取得できる", func(t *testing.T) {
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")

		assert.Equal(t, []string{"2"}, ss.Pop(1))
	})
	t.Run("取得できる", func(t *testing.T) {
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")
		ss.Append("3")
		ss.Append("4")
		ss.Append("5")

		assert.Equal(t, []string{"4", "5"}, ss.Pop(2))
	})
	t.Run("長さを超えて指定するとある分だけ返す", func(t *testing.T) {
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")

		assert.Equal(t, []string{"1", "2"}, ss.Pop(5))
	})
	t.Run("取得した分は消える", func(t *testing.T) {
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")
		ss.Append("3")
		ss.Append("4")
		ss.Append("5")

		ss.Pop(2)
		assert.Equal(t, 3, len(ss.Latest(10)))
		ss.Pop(1)
		assert.Equal(t, 2, len(ss.Latest(10)))
		ss.Pop(0)
		assert.Equal(t, 2, len(ss.Latest(10)))
		ss.Pop(2)
		assert.Equal(t, 0, len(ss.Latest(10)))
		ss.Pop(2)
		assert.Equal(t, 0, len(ss.Latest(10)))
	})
}

func TestFlush(t *testing.T) {
	t.Run("リセットできる", func(t *testing.T) {
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")

		assert.Equal(t, []string{"1", "2"}, ss.Latest(5))
		ss.Flush()
		assert.Equal(t, []string{}, ss.Latest(5))
	})
}
