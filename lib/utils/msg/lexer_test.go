package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `こんにちは[r]あああ
←改行した。[p]`
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
			expectedType:    TEXT,
			expectedLiteral: "r",
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
			expectedType:    TEXT,
			expectedLiteral: "p",
		},
		{
			expectedType:    RBRACKET,
			expectedLiteral: "]",
		},
	}

	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}
