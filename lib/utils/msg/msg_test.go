package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	q := NewQueue()
	q.events = append(q.events, &msg{
		body: []rune("こんにちは"),
	})
	q.events = append(q.events, &msg{
		body: []rune("こんばんは"),
	})
	q.Exec()
	assert.Equal(t, "こんにちは", q.buf)
	q.Next()
	q.Exec()
	assert.Equal(t, "こんにちはこんばんは", q.buf)
}

func TestMsg(t *testing.T) {
	q := NewQueue()
	q.events = append(q.events, &msg{
		body: []rune("こんにちは"),
	})
	q.exec()
	q.exec()
	assert.Equal(t, "こん", q.buf)
	q.exec()
	q.exec()
	q.exec()
	assert.Equal(t, "こんにちは", q.buf)
	q.exec()
	assert.Equal(t, "こんにちは", q.buf)
	q.exec()
}

func TestWait(t *testing.T) {
	q := NewQueue()
	q.events = append(q.events, &msg{
		body: []rune("東京"),
	})
	q.events = append(q.events, &flush{})
	q.events = append(q.events, &msg{
		body: []rune("京都"),
	})
	q.exec()
	q.exec()
	assert.Equal(t, "東京", q.buf)
	q.exec()
	assert.Equal(t, "東京", q.buf)
	q.Next() // flush
	q.exec()
	assert.Equal(t, "", q.buf)
	q.Next()
	q.exec()
	q.exec()
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

	q := NewQueue()
	q.events = e.events

	q.Exec()
	assert.Equal(t, "こんにちは...", q.buf)
	q.Next()
	q.Exec()
	// (flush実行)
	q.Next()
	q.Exec()
	assert.Equal(t, "今日はいかがですか", q.buf)
}
