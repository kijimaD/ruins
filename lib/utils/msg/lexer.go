package msg

type Lexer struct {
	input        string
	position     int // 現在検査中のバイトchの位置
	readPosition int // 入力における次の位置
	ch           byte
	OnIdent      bool
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 次の1文字を読んでinput文字列の現在位置を進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCIIコードの"NUL"文字に対応している
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 現在の1文字を読みこんでトークンを返す
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '[':
		tok = newToken(LBRACKET, l.ch)
		l.OnIdent = true
	case ']':
		tok = newToken(RBRACKET, l.ch)
		l.OnIdent = false
	case '=':
		tok = newToken(EQUAL, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if l.OnIdent {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal) // 予約語
			return tok
		} else {
			tok.Literal = l.readText()
			tok.Type = TEXT
			return tok
		}
	}

	l.readChar()
	return tok
}

// トークンを初期化する
func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// 予約語を読み込み
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 半角スペースを読み飛ばす
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// のぞき見(peek)。readChar()の、文字解析器を進めないバージョン。先読みだけを行う
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition] // 次の位置を返す
	}
}

func (l *Lexer) readText() string {
	position := l.position
	for {
		l.readChar()
		if l.ch == '[' || l.ch == ']' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// 英字か判定する
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
