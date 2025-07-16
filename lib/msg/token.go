package msg

// TokenType はトークンの種類を表す型
type TokenType string

// Token はトークンを表す構造体
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// ILLEGAL は不正なトークン
	ILLEGAL = "ILLEGAL"
	// EOF はファイル終端
	EOF = "EOF"

	// STRING は識別子 + リテラル。数値や変数名など、予約語ではないもの。
	STRING = "STRING"
	// IDENT は識別子
	IDENT = "IDENT"
	// TEXT はテキスト
	TEXT = "TEXT"

	// LBRACKET は左角括弧
	LBRACKET = "["
	// RBRACKET は右角括弧
	RBRACKET = "]"
	// COMMA はカンマ
	COMMA = ","
	// EQUAL は等号
	EQUAL = "="

	// CmdFlush はフラッシュコマンド
	CmdFlush = "p"
	// CmdLineEndWait は行末待機コマンド
	CmdLineEndWait = "l"
	// CmdImage はイメージコマンド
	CmdImage = "image"
	// CmdWait は待機コマンド
	CmdWait = "wait"
)

// 予約語
var keywords = map[string]TokenType{}

// LookupIdent は予約語の場合はその種類を、それ以外の場合はIDENTを返す
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
