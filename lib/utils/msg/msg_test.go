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

	noneOpt := submitOpt{}
	str, result := p.submit(noneOpt)
	assert.Equal(t, "a", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(noneOpt)
	assert.Equal(t, "b", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(noneOpt)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitFail, result)

	nextOpt := submitOpt{lineNext: true}
	str, result = p.submit(nextOpt)
	assert.Equal(t, "c", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(nextOpt)
	assert.Equal(t, "d", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = p.submit(nextOpt)
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

	noneOpt := submitOpt{}
	lineOpt := submitOpt{lineNext: true}
	pageOpt := submitOpt{pageNext: true}

	str, result := m.submit(noneOpt)
	assert.Equal(t, "a", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(noneOpt)
	assert.Equal(t, "b", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(noneOpt)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitFail, result)

	str, result = m.submit(lineOpt)
	assert.Equal(t, "c", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(noneOpt)
	assert.Equal(t, "d", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(noneOpt)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitFail, result)

	str, result = m.submit(lineOpt)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitFail, result)

	str, result = m.submit(pageOpt)
	assert.Equal(t, "e", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(noneOpt)
	assert.Equal(t, "f", str)
	assert.Equal(t, SubmitLineOK, result)

	str, result = m.submit(noneOpt)
	assert.Equal(t, "", str)
	assert.Equal(t, SubmitFail, result)
}

// func TestMsgbuf(t *testing.T) {
// 	m := Msg{
// 		pages: []page{
// 			{
// 				lines: []line{
// 					{str: "ab"},
// 					{str: "cd"},
// 				},
// 			},
// 			{
// 				lines: []line{
// 					{str: "ef"},
// 					{str: "gh"},
// 				},
// 			},
// 		},
// 	}

// 	noneOpt := submitOpt{}
// 	// pageOpt := submitOpt{pageNext: true}

// 	str := m.Buf(noneOpt)
// 	assert.Equal(t, "a", str)

// 	str = m.Buf(noneOpt)
// 	assert.Equal(t, "ab", str)

// 	// フラグなしだと先に進まない
// 	str = m.Buf(noneOpt)
// 	assert.Equal(t, "ab", str)

// 	// str = m.Buf()
// 	// assert.Equal(t, "abc", str)

// 	// str = m.Buf()
// 	// assert.Equal(t, "abcd", str)

// 	// // フラグなしだと先に進まない
// 	// str = m.Buf()
// 	// assert.Equal(t, "abcd", str)

// 	// // ページをまたぐとフラッシュされる
// 	// str = m.Buf()
// 	// assert.Equal(t, "e", str)
// 	// str = m.Buf()
// 	// assert.Equal(t, "ef", str)
// }
