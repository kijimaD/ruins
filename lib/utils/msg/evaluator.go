package msg

type Evaluator struct {
	Events []event
}

func (e *Evaluator) Eval(node Node) event {
	switch node := node.(type) {
	case *Program:
		return e.evalProgram(node)
	case *ExpressionStatement:
		return e.Eval(node.Expression)
	case *CmdExpression:
		var event event
		switch node.Cmd.String() {
		case "p":
			event = &flush{}
			e.Events = append(e.Events, event)
		case "l":
			event = &lineEndWait{}
			e.Events = append(e.Events, event)
		}
		return event
	case *TextLiteral:
		m := &msgEmit{body: []rune(node.Value)}
		e.Events = append(e.Events, m)
		return m
	}

	return nil
}

func (e *Evaluator) evalProgram(program *Program) event {
	var result event

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}
