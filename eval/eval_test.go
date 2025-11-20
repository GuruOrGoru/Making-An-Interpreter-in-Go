package eval

import (
	"fmt"
	"testing"

	"github.com/guruorgoru/goru-verbal-interpreter/lexer"
	"github.com/guruorgoru/goru-verbal-interpreter/object"
	"github.com/guruorgoru/goru-verbal-interpreter/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluatedProgram := testEval(tt.input)
		testDeezInts(t, evaluatedProgram, tt.expected)
	}
}

func TestEvalBoleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"satya", true},
		{"jhuth", false},
		{"1 > 5", false},
		{"420 > 69", true},
		{"5 != 5", false},
		{"1 == 1", true},
		{"satya == satya", true},
		{"jhuth == jhuth", true},
		{"satya == jhuth", false},
		{"satya != jhuth", true},
		{"jhuth != satya", true},
		{"(1 < 2) == satya", true},
		{"(1 < 2) == jhuth", false},
		{"(1 > 2) == satya", false},
		{"(1 > 2) == jhuth", true},
	}

	for _, tt := range tests {
		evaluatedProgram := testEval(tt.input)
		testDeezBools(t, evaluatedProgram, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"yadi (satya) { 10 }", 10},
		{"yadi (jhuth) { 10 }", nil},
		{"yadi (1) { 10 }", 10},
		{"yadi (1 < 2) { 10 }", 10},
		{"yadi (1 > 2) { 10 }", nil},
		{"yadi (1 > 2) { 10 } else { 20 }", 20},
		{"yadi (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testDeezInts(t, evaluated, int64(integer))
		} else {
			testDeezNulls(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"firta 10;", 10},
		{"firta 10; 9;", 10},
		{"9; firta 2 + 5; 9;", 7},
		{"yadi (10 > 1) { yadi (10 > 1) { firta 10; } firta 1; }", 10},
	}
	for _, tt := range tests {
		evalt := testEval(tt.input)
		testDeezInts(t, evalt, tt.expected)
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"manau a = 5; a;", 5},
		{"manau a = 5 * 5; a;", 25},
		{"manau a = 5; manau b = a; b;", 5},
		{"manau a = 5; manau b = a; manau c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testDeezInts(t, testEval(tt.input), tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		panic(fmt.Sprintf("parser errors: %v", p.Errors()))
	}

	return Eval(program, object.NewEnvironment())
}

func testDeezInts(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not an intege, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has stupidly wrong value, got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testDeezBools(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not an boolean, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has stupidly wrong value, got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func testDeezNulls(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + satya;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + satya; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-satya",
			"unknown operator: -BOOLEAN",
		},
		{
			"satya + jhuth;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; satya + jhuth; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"yadi (10 > 1) { satya + jhuth; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
yadi (10 > 1) {
yadi (10 > 1) {
firta satya + jhuth;
}
firta 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}
