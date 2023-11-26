package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// 	str := `こんにちは。\n
	// 今日は晴れです。[n]
	// ところで。[l]
	// どうなりました?`

	//	page1 := page{lines: []line{
	//		line{str: "こんにちは。\n今日は晴れです"},
	//	}}
	//	page2 := page{lines: []line{
	//		line{str: "ところで"},
	//		line{str: "どうなりました?"},
	//	}}
	//	expect := MsgBuilder{
	//		pages: []page{
	//			page1,
	//			page2,
	//		},
	//	}
	//
	// assert.Equal(t, expect, NewLexer(str))
}

func TestLineSubmit(t *testing.T) {
	l := line{
		str: "abc",
	}

	tests := []struct {
		expectValue string
		expectCont  bool
	}{
		{
			expectValue: "a",
			expectCont:  true,
		},
		{
			expectValue: "b",
			expectCont:  true,
		},
		{
			expectValue: "c",
			expectCont:  true,
		},
		{
			expectValue: "",
			expectCont:  false,
		},
	}

	for _, tt := range tests {
		str, cont := l.submit()
		assert.Equal(t, tt.expectValue, str)
		assert.Equal(t, tt.expectCont, cont)
	}
}

func TestPageSubmit(t *testing.T) {
	p := page{
		lines: []line{
			{str: "ab"},
			{str: "cd"},
		},
	}

	tests := []struct {
		expectValue string
		expectCont  bool
	}{
		{
			expectValue: "a",
			expectCont:  true,
		},
		{
			expectValue: "b",
			expectCont:  true,
		},
		{
			expectValue: "c",
			expectCont:  true,
		},
		{
			expectValue: "d",
			expectCont:  true,
		},
		{
			expectValue: "",
			expectCont:  false,
		},
	}

	for _, tt := range tests {
		str, cont := p.submit()
		assert.Equal(t, tt.expectValue, str)
		assert.Equal(t, tt.expectCont, cont)
	}
}

func TestMsgSubmit(t *testing.T) {
	m := Msg{
		pages: []page{
			{
				lines: []line{
					{str: "ab"},
					{str: "cd"},
				},
			},
			{
				lines: []line{
					{str: "ef"},
					{str: "gh"},
				},
			},
		},
	}

	tests := []struct {
		expectValue string
		expectCont  bool
	}{
		{
			expectValue: "a",
			expectCont:  true,
		},
		{
			expectValue: "b",
			expectCont:  true,
		},
		{
			expectValue: "c",
			expectCont:  true,
		},
		{
			expectValue: "d",
			expectCont:  true,
		},
		{
			expectValue: "e",
			expectCont:  true,
		},
		{
			expectValue: "f",
			expectCont:  true,
		},
		{
			expectValue: "g",
			expectCont:  true,
		},
		{
			expectValue: "h",
			expectCont:  true,
		},
		{
			expectValue: "",
			expectCont:  false,
		},
	}

	for _, tt := range tests {
		str, cont := m.submit()
		assert.Equal(t, tt.expectValue, str)
		assert.Equal(t, tt.expectCont, cont)
	}
}

func TestMsgbuf(t *testing.T) {
	m := Msg{
		pages: []page{
			{
				lines: []line{
					{str: "ab"},
					{str: "cd"},
				},
			},
			{
				lines: []line{
					{str: "ef"},
					{str: "gh"},
				},
			},
		},
	}

	str := m.Buf()
	assert.Equal(t, "a", str)

	str = m.Buf()
	assert.Equal(t, "ab", str)

	// フラグなしだと先に進まない
	str = m.Buf()
	assert.Equal(t, "ab", str)

	str = m.Buf()
	assert.Equal(t, "abc", str)

	str = m.Buf()
	assert.Equal(t, "abcd", str)

	// フラグなしだと先に進まない
	str = m.Buf()
	assert.Equal(t, "abcd", str)
}
