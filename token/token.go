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
	BANG = "!"
	MINUS = "-"
	ASTERISK = "*"
	SLASH = "/"

	COMMA = ","
	SEMICOLON = ";"
	LESSERTHAN = "<"
	GREATERTHAN = ">"
	EQUALS = "=="
	NOTEQUALS = "!="

	LEFTPARENTHESIS = "("
	RIGHTPARENTHESIS = ")"
	LEFTBRACES = "{"
	RIGHTBRACES = "}"

	FUNCTION = "FUNCTION"
	LET = "LET"
	IF = "IF"
	ELSE = "ELSE"
	TRUE = "TRUE"
	FALSE = "FALSE"
	RETURN = "RETURN"
) 

// Keywords contains the keywords usable in that langauge

var keywords = map[string]TokenType{
	"karya": FUNCTION,
	"manau": LET,
	"yadi": IF,
	"natra": ELSE,
	"satya": TRUE,
	"jhuth": FALSE,
	"firta": RETURN,
}

func LookForIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
