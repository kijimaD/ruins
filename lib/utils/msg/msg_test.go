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
		expectCont  submitResult
	}{
		{
			expectValue: "a",
			expectCont:  SubmitLineOK,
		},
		{
			expectValue: "b",
			expectCont:  SubmitLineOK,
		},
		{
			expectValue: "c",
			expectCont:  SubmitLineOK,
		},
		{
			expectValue: "",
			expectCont:  SubmitLineFinish,
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

	str, result := p.submit(submitStop)
	assert.Equal(t, "a", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(submitStop)
	assert.Equal(t, "b", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(submitStop)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitLineFinish, result)

	str, result = p.submit(submitLineNext)
	assert.Equal(t, "c", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(submitLineNext)
	assert.Equal(t, "d", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(submitLineNext)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitPageFinish, result)
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
				},
			},
		},
	}

	str, result := m.submit(submitStop)
	assert.Equal(t, "a", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(submitStop)
	assert.Equal(t, "b", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(submitStop)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitLineFinish, result)

	str, result = m.submit(submitLineNext)
	assert.Equal(t, "c", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(submitStop)
	assert.Equal(t, "d", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(submitStop)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitLineFinish, result)

	str, result = m.submit(submitLineNext)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitPageFinish, result)

	str, result = m.submit(submitPageNext)
	assert.Equal(t, "e", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(submitStop)
	assert.Equal(t, "f", str)
	assert.Equal(t, SubmitLineOK, result)

	// これがあると状態が変わって後続の返すステータスが変わる。ビミョーな挙動なので直す
	// str, result = m.submit(submitPageNext)
	// assert.Equal(t, "", str)
	// assert.Equal(t, SubmitMsgFinish, result)

	str, result = m.submit(submitStop)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitLineFinish, result)
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

	str := m.Buf(submitStop)
	assert.Equal(t, "a", str)

	str = m.Buf(submitStop)
	assert.Equal(t, "ab", str)

	// フラグなしだと先に進まない
	str = m.Buf(submitStop)
	assert.Equal(t, "ab\n", str)

	str = m.Buf(submitLineNext)
	assert.Equal(t, "ab\nc", str)

	str = m.Buf(submitStop)
	assert.Equal(t, "ab\ncd", str)

	// ここでLINE_FINISHが出る。何もなくても、一度進まないといけない
	// 次に何もないときには、一気にPAGE_FINISHにしたいな
	str = m.Buf(submitStop)
	assert.Equal(t, "ab\ncd\n", str)

	// ここでPAGE_FINISHが出る。
	str = m.Buf(submitLineNext)
	assert.Equal(t, "", str) // nextしても改行だけのときは、無効にしたいな...

	str = m.Buf(submitPageNext)
	assert.Equal(t, "e", str)

	str = m.Buf(submitStop)
	assert.Equal(t, "ef", str)

	str = m.Buf(submitLineNext)
	assert.Equal(t, "efg", str)

	str = m.Buf(submitStop)
	assert.Equal(t, "efgh", str)

	str = m.Buf(submitStop)
	assert.Equal(t, "efgh\n", str)
}
