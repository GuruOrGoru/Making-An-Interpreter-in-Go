package parser

import (
	"fmt"

	"github.com/guruorgoru/goru-verbal-interpreter/ast"
	"github.com/guruorgoru/goru-verbal-interpreter/lexer"
	"github.com/guruorgoru/goru-verbal-interpreter/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	nextToken    token.Token
	errors       []string
}

func (p *Parser) Errors() []string {
	return p.errors
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}
	p.readNextToken()
	p.readNextToken()
	return p
}

func (p *Parser) readNextToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		statement := p.ParseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.readNextToken()
	}

	return program
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectNextToken(token.ASSIGN) {
		return nil
	}

	for !p.currentTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.readNextToken()

	for !p.currentTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}

	return statement
}

func (p *Parser) expectNextToken(t token.TokenType) bool {
	if p.nextToken.Type == t {
		p.readNextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := "expected next token to be %s, got %s instead"
	p.errors = append(p.errors, fmt.Sprintf(msg, t, p.nextToken.Type))
}
