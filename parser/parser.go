package parser

import (
	"fmt"
	"strconv"

	"github.com/guruorgoru/goru-verbal-interpreter/ast"
	"github.com/guruorgoru/goru-verbal-interpreter/lexer"
	"github.com/guruorgoru/goru-verbal-interpreter/token"
)

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQUALS:          EQUALS,
	token.NOTEQUALS:       EQUALS,
	token.LESSERTHAN:      LESSGREATER,
	token.GREATERTHAN:     LESSGREATER,
	token.PLUS:            SUM,
	token.MINUS:           SUM,
	token.SLASH:           PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.LEFTPARENTHESIS: CALL,
}

type Parser struct {
	lexer            *lexer.Lexer
	currentToken     token.Token
	nextToken        token.Token
	errors           []string
	prefixParseFuncs map[token.TokenType]prefixParseFunc
	infixParseFuncs  map[token.TokenType]infixParseFunc
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(token token.TokenType, prefixFunc prefixParseFunc) {
	p.prefixParseFuncs[token] = prefixFunc
}

func (p *Parser) registerInfix(token token.TokenType, infixFunc infixParseFunc) {
	p.infixParseFuncs[token] = infixFunc
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}
	p.readNextToken()
	p.readNextToken()

	p.prefixParseFuncs = make(map[token.TokenType]prefixParseFunc)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.infixParseFuncs = make(map[token.TokenType]infixParseFunc)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQUALS, p.parseInfixExpression)
	p.registerInfix(token.NOTEQUALS, p.parseInfixExpression)
	p.registerInfix(token.LESSERTHAN, p.parseInfixExpression)

	p.registerInfix(token.GREATERTHAN, p.parseInfixExpression)

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.readNextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.readNextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}


func (p *Parser) readNextToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) nextPrecedence() int {
	if p, ok := precedences[p.nextToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
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
		return p.parseExpressionStatement()
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

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := "could not parse %q as integer"
		p.errors = append(p.errors, fmt.Sprintf(msg, p.currentToken.Literal))
		return nil
	}
	literal.Value = value
	return literal
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.readNextToken()

	for !p.currentTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{}

	statement.Expression = p.parseExpression(LOWEST)

	if p.nextToken.Type == token.SEMICOLON {
		p.readNextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFuncs[p.currentToken.Type]

	if prefix == nil {
		msg := fmt.Sprintf("no prefix parse function for %v", p.currentToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExpression := prefix()

	for (p.nextToken.Type != token.SEMICOLON) && precedence < p.nextPrecedence() {
		infix := p.infixParseFuncs[p.nextToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.readNextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
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
