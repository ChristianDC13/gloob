package interpreter

import (
	"fmt"
	"gloob-interpreter/internal/colors"
	"gloob-interpreter/internal/parser"
	"gloob-interpreter/internal/scope"
	"gloob-interpreter/internal/values"
	"os"
)

// Evaluate is the main dispatch function for the runtime interpreter.
// It takes any AST node (Statement or Expression) and routes it to the appropriate
// evaluator function based on the node's type. This is the heart of the interpreter.
//
// The function uses a switch statement to determine which evaluator to call:
// - Literal values (numbers, strings, booleans, null) are converted to runtime values
// - Expressions are evaluated recursively
// - Statements are executed and may modify the global scope
//
// Examples of what gets evaluated:
// - Program: Executes all statements in sequence
// - VariableDeclaration: Declares variables in the scope
// - BinaryExpression: Performs arithmetic and comparison operations
// - Object: Creates object values with properties
// - MemberAccess: Accesses object properties
// - CallExpression: Executes function calls
// - IfStatement: Executes conditional logic
func Evaluate(node parser.Statement, s *scope.Scope) values.RuntimeValue {
	switch node.NodeType() {
	// Literal values - convert directly to runtime values
	case parser.NodeTypeNumeric:
		return &values.NumericValue{Type: node.NodeType(), Value: node.(*parser.Numeric).Value}
	case parser.NodeTypeBoolean:
		return &values.BooleanValue{Type: node.NodeType(), Value: node.(*parser.Boolean).Value}
	case parser.NodeTypeNull:
		return &values.NullValue{Type: parser.NodeTypeNull}
	case parser.NodeTypeString:
		return &values.StringValue{Type: parser.NodeTypeString, Value: node.(*parser.String).Value}

	// Expressions - evaluate recursively
	case parser.NodeTypeBinaryExpression:
		return evaluateBinaryExpression(node.(*parser.BinaryExpression), s)
	case parser.NodeTypeIdentifier:
		return evaluateIdentifier(node.(*parser.Identifier), s)
	case parser.NodeTypeObject:
		return evaluateObject(node.(*parser.Object), s)
	case parser.NodeTypeArray:
		return evaluateArray(node.(*parser.Array), s)
	case parser.NodeTypeArrayIndex:
		return evaluateArrayIndex(node.(*parser.ArrayIndex), s)
	case parser.NodeTypeMemberAccess:
		return evaluateMemberAccess(node.(*parser.MemberAccess), s)
	case parser.NodeTypeCallExpression:
		return evaluateCallExpression(node.(*parser.CallExpression), s)

	// Statements - execute and potentially modify scope
	case parser.NodeTypeProgram:
		return evaluateProgram(node.(*parser.Program), s)
	case parser.NodeTypeVariableDeclaration:
		return evaluateVariableDeclaration(node.(*parser.VariableDeclaration), node.(*parser.VariableDeclaration).Constant, s)
	case parser.NodeTypeVariableAssignment:
		return evaluateVariableAssignment(node.(*parser.VariableAssignmentExpression), s)
	case parser.NodeTypeFunctionDeclaration:
		return evaluateFunctionDeclaration(node.(*parser.FunctionDeclaration), s)
	case parser.NodeTypeIfStatement:
		return evaluateIfStatement(node.(*parser.IfStatement), s)
	case parser.NodeTypeLoopStatement:
		return evaluateLoopStatement(node.(*parser.LoopStatement), s)
	case parser.NodeTypeBreakExpression:
		return evaluateBreakExpression(node.(*parser.BreakExpression), s)
	case parser.NodeTypeReturnStatement:
		return evaluateReturnStatement(node.(*parser.ReturnStatement), s)
	// Native functions - return as-is
	case parser.NodeTypeNativeFunction:
		return node.(*values.NativeFunctionValue)

	default:
		fmt.Printf("Unknown node type: %s, i don't know what to tell you 🫣\n", colors.Red(node.NodeType()))
		os.Exit(1)
		return nil
	}
}
