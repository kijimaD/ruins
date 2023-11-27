package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	q := Queue{active: true}
	q.events = append(q.events, &msgEmit{
		body: []rune("こんにちは"),
	})
	q.Exec()
	q.Exec()
	assert.Equal(t, "こん", q.buf)
	q.Exec()
	q.Exec()
	q.Exec()
	assert.Equal(t, "こんにちは", q.buf)
	q.Exec()
	assert.Equal(t, "こんにちは", q.buf)
	q.Exec()
}

func TestWait(t *testing.T) {
	q := Queue{active: true}
	q.events = append(q.events, &msgEmit{
		body: []rune("東京"),
	})
	q.events = append(q.events, &flush{})
	q.events = append(q.events, &msgEmit{
		body: []rune("京都"),
	})
	q.Exec()
	q.Exec()
	assert.Equal(t, "東京", q.buf)
	q.Exec()
	assert.Equal(t, "東京", q.buf)
	q.Pop() // flush
	q.Exec()
	assert.Equal(t, "", q.buf)
	q.Exec()
	q.Exec()
	assert.Equal(t, "京都", q.buf)
}

func TestBuilder(t *testing.T) {
	input := `こんにちは...[r]
今日はいかがですか`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	e := Evaluator{}
	e.Eval(program)

	q := NewQueue(e.Events)
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	assert.Equal(t, "こんにちは...", q.buf)
	q.Pop()
	// (flush実行)
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	assert.Equal(t, "今日はいかがですか", q.buf)
}

// 改行を自動挿入できる
func TestNewLine(t *testing.T) {
	input := `こんにちは[r]
ああああああああああああああああああああ`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	e := Evaluator{}
	e.Eval(program)

	q := NewQueue(e.Events)
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	assert.Equal(t, "こんにちは", q.buf)
	q.Pop()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	assert.Equal(t, "ああああああああああああああ\nあああああ", q.buf)
}

// 意図的な改行で自動改行カウントをリセットする
func TestNewLineResetCount(t *testing.T) {
	input := `こんにちは[r]
ああああああああああ
ああああああああああ`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	e := Evaluator{}
	e.Eval(program)

	q := NewQueue(e.Events)
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	assert.Equal(t, "こんにちは", q.buf)
	q.Pop()
	// 意図的に挿入した改行2つ分Execが増える
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	q.Exec()
	assert.Equal(t, "ああああああああああ\nああああああああああ", q.buf)
}
