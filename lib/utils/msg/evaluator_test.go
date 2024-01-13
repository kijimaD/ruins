package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	input := `こんにちは[l]世界[p]
←無視される改行たたたたた。
←有効な改行`

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	e := Evaluator{}
	e.Eval(program)
	results := []string{}
	for _, e := range e.Events {
		switch event := e.(type) {
		case *msgEmit:
			results = append(results, string(event.body))
		case *flush:
			results = append(results, "flush")
		case *lineEndWait:
			results = append(results, "lineEndWait")
		}
	}
	expect := []string{
		"こんにちは",
		"lineEndWait",
		"世界",
		"flush",
		"←無視される改行たたたたた。\n←有効な改行",
	}
	assert.Equal(t, expect, results)
}
