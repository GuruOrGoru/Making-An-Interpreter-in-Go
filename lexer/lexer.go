package lexer

import "github.com/guruorgoru/goru-verbal-interpreter/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.ch = 0
	} else {
		lex.ch = lex.input[lex.readPosition]
	}
	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	switch lex.ch {
	case '=':
		tok = token.Token{Type: token.ASSIGN, Literal: string(lex.ch)}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(lex.ch)}
	case '(':
		tok = token.Token{Type: token.LEFTPARENTHESIS, Literal: string(lex.ch)}
	case ')':
		tok = token.Token{Type: token.RIGHTPARENTHESIS, Literal: string(lex.ch)}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(lex.ch)}
	case '{':
		tok = token.Token{Type: token.LEFTBRACES, Literal: string(lex.ch)}
	case '}':
		tok = token.Token{Type: token.RIGHTBRACES, Literal: string(lex.ch)}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(lex.ch)}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	lex.readChar()
	return tok
}
