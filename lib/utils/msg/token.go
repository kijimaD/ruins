package msg

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	TEXT = "TEXT"

	LBRACKET = "["
	RBRACKET = "]"
)
