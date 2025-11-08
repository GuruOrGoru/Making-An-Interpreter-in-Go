package parser

import (
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

	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nilj")
	}

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
			t.Fatalf("statement.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		}
		letStmt, ok := statement.(*ast.LetStatement)
		if !ok {
			t.Fatalf("statement not *ast.LetStatement. got=%T", statement)
		}

		if letStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("letStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, letStmt.Name.Value)
		}

	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	firta 5;
	firta 10;
	firta 993322;
	`
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "firta" {
			t.Errorf("returnStmt.TokenLiteral not 'firta', got %q", statement.TokenLiteral())
		}
	}
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
