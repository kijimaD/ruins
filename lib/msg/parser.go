package msg

import (
	"fmt"
	"log"
)

// Parser は構文解析器を表す構造体
type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token // 現在のトークン
	peekToken Token // 次のトークン

	// 構文解析関数
	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

type (
	// 前置構文解析関数。前置演算子には「左側」が存在しない
	prefixParseFn func() Expression
	// 中置構文解析関数 n + 1
	// 引数は中置演算子の「左側」
	infixParseFn func(Expression) Expression
)

const (
	_ int = iota
	// LOWEST は最低優先度
	LOWEST
	// CMD はコマンド優先度
	CMD // [...]
)

// 優先順位テーブル。トークンタイプと優先順位を関連付ける
var precedences = map[TokenType]int{
	LBRACKET: CMD,
}

// NewParser は字句解析器を受け取って初期化する
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 前置トークン
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(TEXT, p.parseTextLiteral)
	p.registerPrefix(LBRACKET, p.parseFunctionLiteral)

	// 2つトークンを読み込む。curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

// Errors はエラーのアクセサ
func (p *Parser) Errors() []string {
	return p.errors
}

// 次のトークンに進む
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram はパースを開始する。トークンを1つずつ辿る
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// 文をパースする。トークンの型によって適用関数を変える
func (p *Parser) parseStatement() Statement {
	// 式文の構文解析を試みる
	return p.parseExpressionStatement()
}

// 式文を構文解析する
func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// 次のトークンと引数の型を比較する
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

// 構文解析関数を登録する
func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// 優先順位の処理
	// より低い優先順位のトークンに遭遇する間繰り返す
	// 優先順位が同じもしくは高いトークンに遭遇すると実行しない
	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// 次のトークンタイプに対応している優先順位を返す
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 文字列トークンをパース
func (p *Parser) parseTextLiteral() Expression {
	return &TextLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// コマンドリテラルをパース
// [image storage="test.png"]
// [p]
func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}
	p.nextToken()
	ident := Identifier{Token: p.curToken, Value: p.curToken.Literal}
	lit.FuncName = ident

	if !p.peekTokenIs(RBRACKET) {
		p.nextToken()
	}
	lit.Parameters = p.parseFunctionParameters()

	p.nextToken()

	return lit
}

// 引数をパース
func (p *Parser) parseFunctionParameters() NamedParams {
	namedParams := NamedParams{}
	namedParams.Map = map[string]string{}

	for !p.peekTokenIs(RBRACKET) {
		name := Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if !p.peekTokenIs(EQUAL) {
			log.Fatal("シンタックスエラー: EQUALがない: ", p.curToken.Literal)
		}
		p.nextToken()
		if !p.peekTokenIs(STRING) {
			log.Fatal("シンタックスエラー: STRINGがない: ", p.curToken.Literal)
		}
		p.nextToken()
		namedParams.Map[name.Value] = p.curToken.Literal

		if p.peekTokenIs(RBRACKET) {
			break
		}
		p.nextToken()
	}

	return namedParams
}
