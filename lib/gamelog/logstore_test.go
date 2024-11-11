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
