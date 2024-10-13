package gamelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatest(t *testing.T) {
	{
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")
		assert.Equal(t, []string{"1", "2"}, ss.Latest(5))
		assert.Equal(t, []string{"2"}, ss.Latest(1))
	}
	{
		ss := SafeSlice{}
		ss.Append("1")
		ss.Append("2")
		ss.Append("3")
		ss.Append("4")
		ss.Append("5")
		assert.Equal(t, []string{"3", "4", "5"}, ss.Latest(3))
		assert.Equal(t, []string{"5"}, ss.Latest(1))
	}
}

func TestFlush(t *testing.T) {
	ss := SafeSlice{}
	ss.Append("1")
	ss.Append("2")
	assert.Equal(t, []string{"1", "2"}, ss.Latest(5))
	ss.Flush()
	assert.Equal(t, []string{}, ss.Latest(5))
}
