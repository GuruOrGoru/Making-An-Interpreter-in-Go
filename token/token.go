package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENTIFIER =  "IDENTIFIER"
	INT = "INT"

	PLUS = "+"
	ASSIGN = "="

	COMMA = ","
	SEMICOLON = ";"

	LEFTPARENTHESIS = "("
	RIGHTPARENTHESIS = ")"
	LEFTBRACES = "{"
	RIGHTBRACES = "}"

	FUNCTION = "FUNCTION"
	LET = "LET"
)
