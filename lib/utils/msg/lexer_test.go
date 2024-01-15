package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `こんにちは[l]あああ
←改行した。[p]
[image source="test.png" page="fore"]
[wait time="100"]`
	l := NewLexer(input)

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{
			expectedType:    TEXT,
			expectedLiteral: "こんにちは",
		},
		{
			expectedType:    LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    IDENT,
			expectedLiteral: "l",
		},
		{
			expectedType:    RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    TEXT,
			expectedLiteral: "あああ\n←改行した。",
		},
		{
			expectedType:    LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    IDENT,
			expectedLiteral: "p",
		},
		{
			expectedType:    RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    IDENT,
			expectedLiteral: "image",
		},
		{
			expectedType:    IDENT,
			expectedLiteral: "source",
		},
		{
			expectedType:    EQUAL,
			expectedLiteral: "=",
		},
		{
			expectedType:    STRING,
			expectedLiteral: "test.png",
		},
		{
			expectedType:    IDENT,
			expectedLiteral: "page",
		},
		{
			expectedType:    EQUAL,
			expectedLiteral: "=",
		},
		{
			expectedType:    STRING,
			expectedLiteral: "fore",
		},
		{
			expectedType:    RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    LBRACKET,
			expectedLiteral: "[",
		},
		{
			expectedType:    IDENT,
			expectedLiteral: "wait",
		},
		{
			expectedType:    IDENT,
			expectedLiteral: "time",
		},
		{
			expectedType:    EQUAL,
			expectedLiteral: "=",
		},
		{
			expectedType:    STRING,
			expectedLiteral: "100",
		},
		{
			expectedType:    RBRACKET,
			expectedLiteral: "]",
		},
		{
			expectedType:    EOF,
			expectedLiteral: "",
		},
	}

	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}
