package lexer

import (
	"testing"

	"github.com/guruorgoru/goru-verbal-interpreter/token"
)

func TestNextToken(T *testing.T) {
	input := `
	manau a = 5;
	manau b = 10;
	manau add = karya(x, y) {
		x + y;
	}

	manau result = add(a, b)
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// manau a = 5;
		{token.LET, "manau"},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// manau b = 10;
		{token.LET, "manau"},
		{token.IDENTIFIER, "b"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// manau add = karya(x, y) { x + y; }
		{token.LET, "manau"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "karya"},
		{token.LEFTPARENTHESIS, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHTPARENTHESIS, ")"},
		{token.LEFTBRACES, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHTBRACES, "}"},

		// manau result = add(a, b)
		{token.LET, "manau"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFTPARENTHESIS, "("},
		{token.IDENTIFIER, "a"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "b"},
		{token.RIGHTPARENTHESIS, ")"},
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
