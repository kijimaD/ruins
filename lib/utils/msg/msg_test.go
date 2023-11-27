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

func TestSingle(t *testing.T) {
	input := `こんにちは...`
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
}
