package msg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	input := `こんにちは[l]世界[p]
←無視される改行たたたたた。
←有効な改行
[image source="test.png"]
[wait time="100"]`

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
		case *ChangeBg:
			results = append(results, fmt.Sprintf("changeBg source=%s", event.Source))
		case *wait:
			results = append(results, fmt.Sprintf("wait time=%s", event.durationMsec))
		}
	}
	expect := []string{
		"こんにちは",
		"lineEndWait",
		"世界",
		"flush",
		"←無視される改行たたたたた。\n←有効な改行\n",
		"changeBg source=test.png",
		"wait time=100ms",
	}
	assert.Equal(t, expect, results)
}
