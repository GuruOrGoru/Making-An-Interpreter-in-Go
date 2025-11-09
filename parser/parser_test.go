package parser

import (
	"fmt"
	"testing"

	"github.com/guruorgoru/goru-verbal-interpreter/ast"
	"github.com/guruorgoru/goru-verbal-interpreter/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
	manau x = 69;
	manau y = 420;
	manau foobar  = 42069;
	`

	program := parseProgram(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if statement.TokenLiteral() != "manau" {
			t.Fatalf("statement.TokenLiteral not 'manau'. got=%q", statement.TokenLiteral())
		}
		letStmt, ok := statement.(*ast.LetStatement)
		if !ok {
			t.Fatalf("statement not *ast.LetStatement. got=%T", statement)
		}

		testIdentifier(t, letStmt.Name, tt.expectedIdentifier)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	firta 5;
	firta 10;
	firta 993322;
	`
	program := parseProgram(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		testReturnStatement(t, statement, "firta")
	}
}

func TestIdentifierEpression(t *testing.T) {
	input := "foobar;"

	program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	testIdentifier(t, ident, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	testIntegerLiteral(t, literal, 5)
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!69;", "!", 69},
		{"-420;", "-", 420},
	}

	for _, tt := range prefixTests {
		program := parseProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		testPrefixExpression(t, exp, tt.operator, tt.intValue)
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"69 + 420;", 69, "+", 420},
		{"993322 - 123456;", 993322, "-", 123456},
		{"50 * 2;", 50, "*", 2},
		{"100 / 4;", 100, "/", 4},
	}

	for _, tt := range infixTests {
		program := parseProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func (p *Parser) TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		program := parseProgram(t, tt.input)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

// testing booleans

func TestBooleanExpression(t *testing.T) {
	input := []struct {
		input         string
		expectedValue bool
	}{
		{"satya;", true},
		{"jhuth;", false},
	}
	for _, tt := range input {
		program := parseProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}

		if boolean.Value != tt.expectedValue {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedValue, boolean.Value)
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func parseProgram(t *testing.T, input string) *ast.Program {
	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	return program
}

func testReturnStatement(t *testing.T, stmt ast.Statement, expectedToken string) bool {
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
		return false
	}
	if returnStmt.TokenLiteral() != expectedToken {
		t.Errorf("returnStmt.TokenLiteral not '%s', got %q", expectedToken, returnStmt.TokenLiteral())
		return false
	}
	return true
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right any) bool {
	prefixExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp is not ast.PrefixExpression. got=%T", exp)
		return false
	}
	if prefixExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%s", operator, prefixExp.Operator)
		return false
	}
	if !testLiteralExpression(t, prefixExp.Right, right) {
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any,
	operator string, right any,
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}
