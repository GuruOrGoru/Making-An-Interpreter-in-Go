package lexer

import (
	"github.com/guruorgoru/goru-verbal-interpreter/token"
)

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

	for lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\n' || lex.ch == '\r' {
		lex.readChar()
	}

	switch lex.ch {
	case '=':
		if lex.peekAtNextChar() == '=' {
			currentChar := lex.ch
			lex.readChar()
			literal := string(currentChar) + string(lex.ch)
			tok = token.Token{Type: token.EQUALS, Literal: literal}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(lex.ch)}
		}
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
	case '!':
		if lex.peekAtNextChar() == '=' {
			currentChar := lex.ch
			lex.readChar()
			literal := string(currentChar) + string(lex.ch)
			tok = token.Token{Type: token.NOTEQUALS, Literal: literal}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(lex.ch)}
		}
	case '-':
		tok = token.Token{Type: token.MINUS, Literal: string(lex.ch)}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(lex.ch)}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(lex.ch)}
	case '<':
		tok = token.Token{Type: token.LESSERTHAN, Literal: string(lex.ch)}
	case '>':
		tok = token.Token{Type: token.GREATERTHAN, Literal: string(lex.ch)}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lex.ch) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookForIdentifier(tok.Literal)
			return tok
		} else if isNumber(lex.ch) {
			tok.Literal = lex.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(lex.ch)}
		}
	}

	lex.readChar()
	return tok
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (lex *Lexer) readIdentifier() string {
	position := lex.position

	for isLetter(lex.ch) {
		lex.readChar()
	}

	return lex.input[position:lex.position]
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lex *Lexer) readNumber() string {
	position := lex.position

	for isNumber(lex.ch) {
		lex.readChar()
	}

	return lex.input[position:lex.position]
}

func (lex *Lexer) peekAtNextChar() byte {
	if lex.readPosition >= len(lex.input) {
		return 0
	}
	return lex.input[lex.readPosition]
}
