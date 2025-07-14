package msg

import (
	"fmt"
	"log"
	"time"
)

// Evaluator はASTノードを評価する構造体
type Evaluator struct {
	Events []Event
}

// NewEvaluator は新しいEvaluatorを作成する
func NewEvaluator(node Node) *Evaluator {
	e := Evaluator{}
	e.Eval(node)

	return &e
}

// Eval はノードを評価してイベントを返す
func (e *Evaluator) Eval(node Node) Event {
	switch node := node.(type) {
	case *Program:
		return e.evalProgram(node)
	case *ExpressionStatement:
		return e.Eval(node.Expression)
	case *FunctionLiteral:
		var eve Event
		switch node.FuncName.Value {
		case CmdFlush:
			eve = &flush{}
		case CmdLineEndWait:
			eve = &lineEndWait{}
		case CmdImage:
			eve = &ChangeBg{Source: node.Parameters.Map["source"]}
		case CmdWait:
			duration, err := time.ParseDuration(fmt.Sprintf("%sms", node.Parameters.Map["time"]))
			if err != nil {
				log.Fatal(err)
			}
			eve = &wait{duration: duration}
		}
		e.Events = append(e.Events, eve)
		return eve
	case *TextLiteral:
		m := &msgEmit{body: []rune(node.Value)}
		e.Events = append(e.Events, m)
		return m
	}

	return nil
}

func (e *Evaluator) evalProgram(program *Program) Event {
	var result Event

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}
