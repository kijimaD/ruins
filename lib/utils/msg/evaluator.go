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
	case *FunctionLiteral:
		var eve event
		switch node.FuncName.Value {
		case CMD_FLUSH:
			eve = &flush{}
		case CMD_LINE_END_WAIT:
			eve = &lineEndWait{}
		case CMD_IMAGE:
			eve = &notImplement{}
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

func (e *Evaluator) evalProgram(program *Program) event {
	var result event

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}
