package msg

import (
	"bytes"
	"strings"
)

// Node はASTノードの基底インターフェース
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement はステートメントノードのインターフェース
type Statement interface {
	Node
	statementNode()
}

// Expression は式ノードのインターフェース
type Expression interface {
	Node
	expressionNode()
}

// Program は構文解析器が生成する全てのASTのルートノードになる
type Program struct {
	Statements []Statement
}

// TokenLiteral はインターフェースで定義されている関数の1つ
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// インターフェースで定義されている関数の1つ
// 文字列表示してデバッグしやすいようにする
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ExpressionStatement は式文を表すノード
type ExpressionStatement struct {
	Token      Token      // 式の最初のトークン
	Expression Expression // 式を保持
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral はトークンのリテラル値を返す
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// CmdExpression はコマンド式を表すノード
type CmdExpression struct {
	Token      Token // '['トークン
	Expression Expression
	Cmd        Event
}

func (ie *CmdExpression) expressionNode() {}

// TokenLiteral はトークンのリテラル値を返す
func (ie *CmdExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *CmdExpression) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	out.WriteString(ie.Expression.String())
	out.WriteString("]")

	return out.String()
}

// TextLiteral はテキストリテラルを表すノード
type TextLiteral struct {
	Token Token
	Value string
}

func (sl *TextLiteral) expressionNode() {}

// TokenLiteral はトークンのリテラル値を返す
func (sl *TextLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *TextLiteral) String() string       { return sl.Token.Literal }

// FunctionLiteral は関数リテラルを表すノード
type FunctionLiteral struct {
	Token      Token
	FuncName   Identifier
	Parameters NamedParams
}

func (fl *FunctionLiteral) expressionNode() {} // fnの結果をほかの変数に代入できたりするため。代入式の一部として扱うためには、式でないといけない
// TokenLiteral はトークンのリテラル値を返す
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for k, v := range fl.Parameters.Map {
		params = append(params, k+"="+v)
	}

	out.WriteString("[")
	out.WriteString(fl.FuncName.Value)
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("]")

	return out.String()
}

// Identifier は識別子を表す
type Identifier struct {
	Token Token // token.IDENT トークン
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral はトークンの文字列表現を返す
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// NamedParams は名前付きパラメータを表す
type NamedParams struct {
	Map map[string]string
}

func (n *NamedParams) expressionNode() {}

// TokenLiteral はトークンの文字列表現を返す
func (n *NamedParams) TokenLiteral() string { return "map" }
func (n *NamedParams) String() string {
	var out bytes.Buffer

	for k, v := range n.Map {
		out.WriteString(k)
		out.WriteString(" = ")
		out.WriteString(v)
	}
	return out.String()
}
