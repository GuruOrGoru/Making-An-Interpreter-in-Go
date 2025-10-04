package lexer

import (
	"testing"

	"github.com/guruorgoru/goru-verbal-interpreter/token"
)

func TestNextToken(T *testing.T) {
	input := `+=(){};,`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.ASSIGN, "="},
		{token.LEFTPARENTHESIS, "("},
		{token.RIGHTPARENTHESIS, ")"},
		{token.LEFTBRACES, "{"},
		{token.RIGHTBRACES, "}"},
		{token.SEMICOLON, ";"},
		{token.COMMA, ","},
		{token.EOF, ""},
	}

	lex := New(input)

	for i, test := range tests {
		tok := lex.NextToken()

		if tok.Type != test.expectedType {
			T.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, test.expectedType, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			T.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, test.expectedLiteral, tok.Literal)
		}
	}
}
