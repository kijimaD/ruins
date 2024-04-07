package msg

import (
	"fmt"
	"log"
	"time"
)

type Evaluator struct {
	Events []Event
}

func NewEvaluator(node Node) *Evaluator {
	e := Evaluator{}
	e.Eval(node)

	return &e
}

func (e *Evaluator) Eval(node Node) Event {
	switch node := node.(type) {
	case *Program:
		return e.evalProgram(node)
	case *ExpressionStatement:
		return e.Eval(node.Expression)
	case *FunctionLiteral:
		var eve Event
		switch node.FuncName.Value {
		case CMD_FLUSH:
			eve = &flush{}
		case CMD_LINE_END_WAIT:
			eve = &lineEndWait{}
		case CMD_IMAGE:
			eve = &ChangeBg{Source: node.Parameters.Map["source"]}
		case CMD_WAIT:
			duration, err := time.ParseDuration(fmt.Sprintf("%sms", node.Parameters.Map["time"]))
			if err != nil {
				log.Fatal(err)
			}
			eve = &wait{durationMsec: duration}
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
