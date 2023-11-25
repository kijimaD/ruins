package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	str := `こんにちは。\n
今日は晴れです。[n]
ところで。[l]
どうなりました?`

	page1 := page{lines: []line{
		line{str: "こんにちは。\n今日は晴れです"},
	}}
	page2 := page{lines: []line{
		line{str: "ところで"},
		line{str: "どうなりました?"},
	}}
	expect := MsgBuilder{
		pages: []page{
			page1,
			page2,
		},
	}

	assert.Equal(t, expect, New(str))
}
