package msg

func Eval(node Node) event {
	switch node := node.(type) {

	// 文
	case *Program:
		return evalProgram(node)
	case *ExpressionStatement:
		return Eval(node.Expression)
	case *TextLiteral:
		return &msg{body: []rune(node.Value)}
	}

	return nil
}

func evalProgram(program *Program) event {
	var result event

	for _, statement := range program.Statements {
		result = Eval(statement)
	}
	return result
}
