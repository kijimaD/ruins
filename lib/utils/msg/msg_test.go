package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	q := NewQueue()
	msg := &msg{
		body: []rune("こんにちは"),
	}
	q.events = append(q.events, msg)
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
