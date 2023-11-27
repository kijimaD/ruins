package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingIndexExpressions(t *testing.T) {
	input := `こんにちは[r]世界[p]
←無視される改行たたたたた。
←有効な改行`

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	{
		stmt, ok := program.Statements[0].(*ExpressionStatement)
		assert.True(t, ok)
		textLit, ok := stmt.Expression.(*TextLiteral)
		assert.True(t, ok)
		assert.Equal(t, "こんにちは", textLit.Value)
	}
	{
		stmt, ok := program.Statements[1].(*ExpressionStatement)
		assert.True(t, ok)
		cmdExp, ok := stmt.Expression.(*CmdExpression)
		assert.True(t, ok)
		assert.Equal(t, "[r]", cmdExp.String())
	}
	{
		stmt, ok := program.Statements[2].(*ExpressionStatement)
		assert.True(t, ok)
		textLit, ok := stmt.Expression.(*TextLiteral)
		assert.True(t, ok)
		assert.Equal(t, "世界", textLit.Value)
	}
	{
		stmt, ok := program.Statements[3].(*ExpressionStatement)
		assert.True(t, ok)
		cmdExp, ok := stmt.Expression.(*CmdExpression)
		assert.True(t, ok)
		assert.Equal(t, "[p]", cmdExp.String())
	}
	{
		stmt, ok := program.Statements[4].(*ExpressionStatement)
		assert.True(t, ok)
		textLit, ok := stmt.Expression.(*TextLiteral)
		assert.True(t, ok)
		assert.Equal(t, "←無視される改行たたたたた。\n←有効な改行", textLit.Value)
	}
}
