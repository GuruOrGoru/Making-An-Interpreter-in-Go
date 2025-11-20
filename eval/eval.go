package eval

import (
	"github.com/guruorgoru/goru-verbal-interpreter/ast"
	"github.com/guruorgoru/goru-verbal-interpreter/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.IfStatement:
		return evalIfStatement(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return inputBoolToBoolObj(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		switch node.Operator {
		case "!":
			return evalBangOp(right)
		case "-":
			return evalNegateOp(right)
		default:
			return NULL
		}
	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		switch {
		case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:

			return evalIntegerInfixOp(node.Operator, left, right)
		case node.Operator == "==":
			return inputBoolToBoolObj(left == right)
		case node.Operator == "!=":
			return inputBoolToBoolObj(left != right)
		default:
			return NULL
		}
	case *ast.ReturnStatement:
		var val object.Object
		if node.ReturnValue != nil {
			val = Eval(node.ReturnValue)
		} else {
			val = NULL
		}
		return &object.ReturnValue{Value: val}
	}
	return NULL
}

func evalIntegerInfixOp(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return inputBoolToBoolObj(leftVal < rightVal)
	case ">":
		return inputBoolToBoolObj(leftVal > rightVal)
	case "==":
		return inputBoolToBoolObj(leftVal == rightVal)
	case "!=":
		return inputBoolToBoolObj(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalBangOp(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalNegateOp(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalStmt(stmt []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmt {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalIfStatement(ie *ast.IfStatement) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func inputBoolToBoolObj(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJECT {
			return result
		}
	}
	return result
}
