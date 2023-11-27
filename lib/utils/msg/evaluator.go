package msg

type Evaluator struct {
	events []event
}

func (e *Evaluator) Eval(node Node) event {
	switch node := node.(type) {
	case *Program:
		return e.evalProgram(node)
	case *ExpressionStatement:
		return e.Eval(node.Expression)
	case *CmdExpression:
		m := &flush{}
		e.events = append(e.events, m)
		return m
	case *TextLiteral:
		m := &msg{body: []rune(node.Value)}
		e.events = append(e.events, m)
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
