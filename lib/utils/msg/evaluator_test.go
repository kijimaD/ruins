package msg

import (
	"fmt"
	"testing"
)

func TestEval(t *testing.T) {
	input := `こんにちは[r]世界[p]
←無視される改行たたたたた。
←有効な改行`

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	e := Eval(program)
	v, _ := e.(*msg)
	fmt.Printf("%#v\n", string(v.body))
}
