package interpreter

import (
	"fmt"
	"gloob-interpreter/internal/builtins"
	"gloob-interpreter/internal/colors"
	"gloob-interpreter/internal/parser"
	"gloob-interpreter/internal/scope"
	"gloob-interpreter/internal/values"
	"os"
	"strings"
)

func evaluateBinaryExpression(node *parser.BinaryExpression, s *scope.Scope) values.RuntimeValue {
	left := Evaluate(node.Left, s)
	right := Evaluate(node.Right, s)

	// Handle comparison operators
	if isComparisonOperator(node.Operator) {
		return evaluateComparisonExpression(node.Operator, left, right, s)
	}

	if left.NodeType() == parser.NodeTypeString && node.Operator == "*" && right.NodeType() == parser.NodeTypeNumeric {
		return evaluateStringMultiplication(left.(*values.StringValue), right.(*values.NumericValue), s)
	}

	if left.NodeType() == parser.NodeTypeString || right.NodeType() == parser.NodeTypeString {
		return evaluateStringBinaryExpression(node.Operator, left, right, s)
	}

	if left.NodeType() != parser.NodeTypeNumeric || right.NodeType() != parser.NodeTypeNumeric {
		fmt.Printf("Invalid operand types for binary expression: %s %s %s\n", left.NodeType(), node.Operator, right.NodeType())
		os.Exit(1)
		return nil
	}

	leftNumeric, ok := left.(*values.NumericValue)
	if !ok {
		fmt.Printf("Invalid left operand type for binary expression: %s\n", left.NodeType())
		os.Exit(1)
		return nil
	}
	rightNumeric, ok := right.(*values.NumericValue)
	if !ok {
		fmt.Printf("Invalid right operand type for binary expression: %s\n", right.NodeType())
		os.Exit(1)
		return nil
	}
	return evaluateNumericBinaryExpression(node.Operator, leftNumeric, rightNumeric, s)
}

func evaluateStringMultiplication(left *values.StringValue, right *values.NumericValue, s *scope.Scope) values.RuntimeValue {
	return &values.StringValue{Type: parser.NodeTypeString, Value: strings.Repeat(left.Value, int(right.Value))}
}

func evaluateStringBinaryExpression(operator string, left values.RuntimeValue, right values.RuntimeValue, s *scope.Scope) values.RuntimeValue {
	switch operator {
	case "+":
		return &values.StringValue{Type: parser.NodeTypeString, Value: fmt.Sprintf("%v%v", left, right)}
	}
	fmt.Printf("Unknown operator: %s, with string operands\n", colors.Red(operator))
	os.Exit(1)
	return nil
}

func evaluateNumericBinaryExpression(operator string, left *values.NumericValue, right *values.NumericValue, s *scope.Scope) values.RuntimeValue {
	switch operator {
	case "+":
		return &values.NumericValue{Type: parser.NodeTypeNumeric, Value: left.Value + right.Value}
	case "-":
		return &values.NumericValue{Type: parser.NodeTypeNumeric, Value: left.Value - right.Value}
	case "*":
		return &values.NumericValue{Type: parser.NodeTypeNumeric, Value: left.Value * right.Value}
	case "/":
		if right.Value == 0 {
			fmt.Printf("You know you cannot divide by zero, what are you trying to prove? 😒 \n")
			os.Exit(1)
			return nil
		}
		return &values.NumericValue{Type: parser.NodeTypeNumeric, Value: left.Value / right.Value}
	case "%":
		return &values.NumericValue{Type: parser.NodeTypeNumeric, Value: float64(int(left.Value) % int(right.Value))}

	}
	fmt.Printf("Unknown operator: %s, i don't know what to tell you 🫣\n", colors.Red(operator))
	os.Exit(1)
	return nil
}

func evaluateProgram(program *parser.Program, s *scope.Scope) values.RuntimeValue {

	var lastEvaluated values.RuntimeValue = nil

	for _, statement := range program.Statements {
		lastEvaluated = Evaluate(statement, s)
	}

	return lastEvaluated
}

func evaluateIdentifier(node *parser.Identifier, s *scope.Scope) values.RuntimeValue {
	return s.Get(node.Name)
}

func evaluateVariableDeclaration(node *parser.VariableDeclaration, isConstant bool, s *scope.Scope) values.RuntimeValue {

	var value values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}
	if node.Value != nil {
		value = Evaluate(node.Value, s)
	}
	s.Declare(node.Identifier, value, isConstant)
	return &values.NodeVariableDeclaration{
		Type:  node.NodeType(),
		Name:  node.Identifier,
		Value: value,
	}
}

func evaluateVariableAssignment(node *parser.VariableAssignmentExpression, s *scope.Scope) values.RuntimeValue {
	if node.Identifier.NodeType() == parser.NodeTypeIdentifier {
		// Regular variable assignment
		identifier := node.Identifier.(*parser.Identifier)
		value := Evaluate(node.Value, s)
		s.Assign(identifier.Name, value)
		return value
	} else if node.Identifier.NodeType() == parser.NodeTypeMemberAccess {
		// Member access assignment (e.g., obj.property = value)
		return evaluateMemberAccessAssignment(node.Identifier.(*parser.MemberAccess), node.Value, s)
	} else if node.Identifier.NodeType() == parser.NodeTypeArrayIndex {
		// Array index assignment (e.g., arr[1] = value)
		return evaluateArrayIndexAssignment(node.Identifier.(*parser.ArrayIndex), node.Value, s)
	} else {
		fmt.Printf("Invalid identifier type for variable assignment: %s\n", node.Identifier.NodeType())
		os.Exit(1)
		return nil
	}
}

func evaluateObject(node *parser.Object, s *scope.Scope) values.RuntimeValue {
	properties := make(map[string]values.RuntimeValue)

	for _, property := range node.Properties {
		value := Evaluate(property.Value, s)
		properties[property.Key] = value
	}

	return &values.ObjectValue{
		Type:       parser.NodeTypeObject,
		Properties: properties,
	}
}

func evaluateMemberAccess(node *parser.MemberAccess, s *scope.Scope) values.RuntimeValue {
	object := Evaluate(node.Object, s)

	// Handle array methods
	if object.NodeType() == parser.NodeTypeArray {
		return builtins.GetArrayMethod(object.(*values.ArrayValue), node.Property)
	}

	// Handle string methods
	if object.NodeType() == parser.NodeTypeString {
		return builtins.GetStringMethod(object.(*values.StringValue), node.Property)
	}

	// Handle object properties
	if object.NodeType() != parser.NodeTypeObject {
		fmt.Printf("Cannot access property '%s' on non-object type: %s\n", node.Property, object.NodeType())
		os.Exit(1)
		return nil
	}

	objValue := object.(*values.ObjectValue)
	if value, exists := objValue.Properties[node.Property]; exists {
		return value
	}

	fmt.Printf("Property '%s' not found on object\n", node.Property)
	os.Exit(1)
	return nil
}

func evaluateMemberAccessAssignment(node *parser.MemberAccess, value parser.Expression, s *scope.Scope) values.RuntimeValue {
	object := Evaluate(node.Object, s)

	if object.NodeType() != parser.NodeTypeObject {
		fmt.Printf("Cannot assign property '%s' on non-object type: %s\n", node.Property, object.NodeType())
		os.Exit(1)
		return nil
	}

	objValue := object.(*values.ObjectValue)
	assignedValue := Evaluate(value, s)
	objValue.Properties[node.Property] = assignedValue

	return assignedValue
}

func evaluateArrayIndexAssignment(node *parser.ArrayIndex, value parser.Expression, s *scope.Scope) values.RuntimeValue {
	// Evaluate the array expression
	arrayValue := Evaluate(node.ArrayExpression, s)

	// Check if it's actually an array
	if arrayValue.NodeType() != parser.NodeTypeArray {
		fmt.Printf("Cannot index non-array type: %s\n", arrayValue.NodeType())
		os.Exit(1)
		return nil
	}

	// Evaluate the index
	indexValue := Evaluate(node.Index, s)
	if indexValue.NodeType() != parser.NodeTypeNumeric {
		fmt.Printf("Array index must be numeric\n")
		os.Exit(1)
		return nil
	}

	array := arrayValue.(*values.ArrayValue)
	index := int(indexValue.(*values.NumericValue).Value)

	// Arrays are 1-based in Gloob, convert to 0-based
	index = index - 1

	// Check bounds
	if index < 0 || index >= len(array.Elements) {
		fmt.Printf("Array index out of bounds: %d (array length: %d)\n", index+1, len(array.Elements))
		os.Exit(1)
		return nil
	}

	// Assign the value
	assignedValue := Evaluate(value, s)
	array.Elements[index] = assignedValue

	return assignedValue
}

func evaluateArray(node *parser.Array, s *scope.Scope) values.RuntimeValue {
	elements := make([]values.RuntimeValue, len(node.Elements))

	for i, element := range node.Elements {
		elements[i] = Evaluate(element, s)
	}

	return &values.ArrayValue{
		Type:     parser.NodeTypeArray,
		Elements: elements,
	}
}

func evaluateArrayIndex(node *parser.ArrayIndex, s *scope.Scope) values.RuntimeValue {
	// Evaluate the expression (could be array or string)
	value := Evaluate(node.ArrayExpression, s)

	// Evaluate the index
	indexValue := Evaluate(node.Index, s)
	if indexValue.NodeType() != parser.NodeTypeNumeric {
		fmt.Printf("Index must be numeric\n")
		os.Exit(1)
		return nil
	}

	index := int(indexValue.(*values.NumericValue).Value)
	// Convert 1-based to 0-based
	index = index - 1

	// Handle string indexing
	if value.NodeType() == parser.NodeTypeString {
		str := value.(*values.StringValue)

		// Check bounds
		if index < 0 || index >= len(str.Value) {
			fmt.Printf("String index out of bounds: %d (string length: %d)\n", index+1, len(str.Value))
			os.Exit(1)
			return nil
		}

		// Return single character as a string
		return &values.StringValue{
			Type:  parser.NodeTypeString,
			Value: string(str.Value[index]),
		}
	}

	// Handle array indexing
	if value.NodeType() == parser.NodeTypeArray {
		array := value.(*values.ArrayValue)

		// Check bounds
		if index < 0 || index >= len(array.Elements) {
			fmt.Printf("Array index out of bounds: %d (array length: %d)\n", index+1, len(array.Elements))
			os.Exit(1)
			return nil
		}

		return array.Elements[index]
	}

	// Not an array or string
	fmt.Printf("Cannot index type: %s\n", value.NodeType())
	os.Exit(1)
	return nil
}

func evaluateCallExpression(node *parser.CallExpression, s *scope.Scope) values.RuntimeValue {
	// Evaluate the callee (function identifier)
	calleeValue := Evaluate(node.Callee, s)

	// Check if it's a native function
	if calleeValue.NodeType() == parser.NodeTypeNativeFunction {

		// Cast to NativeFunctionValue
		nativeFunc, ok := calleeValue.(*values.NativeFunctionValue)
		if !ok {
			fmt.Printf("Invalid native function type\n")
			os.Exit(1)
			return nil
		}

		// Evaluate all arguments
		args := make([]values.RuntimeValue, len(node.Args))
		for i, arg := range node.Args {
			args[i] = Evaluate(arg, s)
		}

		// Call the native function
		return nativeFunc.Expression(args, s)
	}

	if calleeValue.NodeType() == parser.NodeTypeFunctionDeclaration {
		fun := calleeValue.(*values.FunctionValue)

		// Check parameter count
		if len(node.Args) != len(fun.Parameters) {
			fmt.Printf("Function %s expects %d arguments, got %d\n", fun.Identifier, len(fun.Parameters), len(node.Args))
			os.Exit(1)
			return nil
		}

		// Evaluate all arguments
		args := make([]values.RuntimeValue, len(node.Args))
		for i, arg := range node.Args {
			args[i] = Evaluate(arg, s)
		}

		// Create function scope
		funScope := scope.NewScope(fun.Scope.(*scope.Scope))

		// Declare parameters in function scope
		for i, paramName := range fun.Parameters {
			funScope.Declare(paramName, args[i], false)
		}

		// Execute function body
		var result values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}
		for _, statement := range fun.Body {
			result = Evaluate(statement, funScope)

			// Check if a return statement was executed
			if result.NodeType() == parser.NodeTypeReturnValue {
				// Unwrap and return the actual value
				return result.(*values.ReturnValue).Value
			}
		}

		// Implicit return: return the last expression's value
		return result
	}

	fmt.Printf("Cannot call non-function value: %s\n", calleeValue.NodeType())
	os.Exit(1)
	return nil
}

func evaluateFunctionDeclaration(node *parser.FunctionDeclaration, s *scope.Scope) values.RuntimeValue {
	fun := &values.FunctionValue{
		Type:       parser.NodeTypeFunctionDeclaration,
		Identifier: node.Identifier,
		Parameters: node.Parameters,
		Body:       node.Body,
		Scope:      s,
	}
	s.Declare(node.Identifier, fun, false)
	return fun
}

// Helper function to check if an operator is a comparison operator
func isComparisonOperator(operator string) bool {
	switch operator {
	case "==", "!=", ">", ">=", "<", "<=", "&&", "||":
		return true
	default:
		return false
	}
}

// Evaluate comparison expressions
func evaluateComparisonExpression(operator string, left values.RuntimeValue, right values.RuntimeValue, s *scope.Scope) values.RuntimeValue {
	// Handle logical operators first (they have special behavior)
	if operator == "&&" || operator == "||" {
		return evaluateLogicalExpression(operator, left, right, s)
	}

	// Handle string comparisons
	if left.NodeType() == parser.NodeTypeString && right.NodeType() == parser.NodeTypeString {
		return evaluateStringComparison(operator, left.(*values.StringValue), right.(*values.StringValue))
	}

	// Handle numeric comparisons
	if left.NodeType() == parser.NodeTypeNumeric && right.NodeType() == parser.NodeTypeNumeric {
		return evaluateNumericComparison(operator, left.(*values.NumericValue), right.(*values.NumericValue))
	}

	// Handle boolean comparisons
	if left.NodeType() == parser.NodeTypeBoolean && right.NodeType() == parser.NodeTypeBoolean {
		return evaluateBooleanComparison(operator, left.(*values.BooleanValue), right.(*values.BooleanValue))
	}

	// Handle null comparisons
	if left.NodeType() == parser.NodeTypeNull && right.NodeType() == parser.NodeTypeNull {
		return evaluateNullComparison(operator)
	}

	// Mixed type comparisons (only == and != are allowed)
	if operator == "==" {
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: false}
	}
	if operator == "!=" {
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: true}
	}

	fmt.Printf("Cannot compare %s and %s with operator %s\n", left.NodeType(), right.NodeType(), operator)
	os.Exit(1)
	return nil
}

// evaluateLogicalExpression handles logical operators && and ||
func evaluateLogicalExpression(operator string, left values.RuntimeValue, right values.RuntimeValue, s *scope.Scope) values.RuntimeValue {
	// Coerce both operands to boolean values
	leftBool := isTruthy(left)
	rightBool := isTruthy(right)

	switch operator {
	case "&&":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: leftBool && rightBool}
	case "||":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: leftBool || rightBool}
	default:
		fmt.Printf("Unknown logical operator: %s\n", operator)
		os.Exit(1)
		return nil
	}
}

func evaluateStringComparison(operator string, left *values.StringValue, right *values.StringValue) values.RuntimeValue {
	switch operator {
	case "==":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value == right.Value}
	case "!=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value != right.Value}
	case ">":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value > right.Value}
	case ">=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value >= right.Value}
	case "<":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value < right.Value}
	case "<=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value <= right.Value}
	default:
		fmt.Printf("Unknown string comparison operator: %s\n", operator)
		os.Exit(1)
		return nil
	}
}

func evaluateNumericComparison(operator string, left *values.NumericValue, right *values.NumericValue) values.RuntimeValue {
	switch operator {
	case "==":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value == right.Value}
	case "!=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value != right.Value}
	case ">":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value > right.Value}
	case ">=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value >= right.Value}
	case "<":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value < right.Value}
	case "<=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value <= right.Value}
	default:
		fmt.Printf("Unknown numeric comparison operator: %s\n", operator)
		os.Exit(1)
		return nil
	}
}

func evaluateBooleanComparison(operator string, left *values.BooleanValue, right *values.BooleanValue) values.RuntimeValue {
	switch operator {
	case "==":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value == right.Value}
	case "!=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value != right.Value}
	case ">":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value && !right.Value}
	case ">=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value || !right.Value}
	case "<":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: !left.Value && right.Value}
	case "<=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: !left.Value || right.Value}
	case "&&":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value && right.Value}
	case "||":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: left.Value || right.Value}

	default:
		fmt.Printf("Unknown boolean comparison operator: %s\n", operator)
		os.Exit(1)
		return nil
	}
}

func evaluateNullComparison(operator string) values.RuntimeValue {
	switch operator {
	case "==":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: true}
	case "!=":
		return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: false}
	default:
		fmt.Printf("Cannot use operator %s with null values\n", operator)
		os.Exit(1)
		return nil
	}
}

func evaluateIfStatement(node *parser.IfStatement, s *scope.Scope) values.RuntimeValue {
	// Evaluate the condition
	conditionValue := Evaluate(node.Condition, s)

	// Check if condition is truthy
	if isTruthy(conditionValue) {
		// Execute if body
		var result values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}
		for _, statement := range node.Body {
			result = Evaluate(statement, s)
		}
		return result
	}

	// Check elseif clauses
	for _, elseifClause := range node.ElseIfs {
		elseifValue := Evaluate(elseifClause.Condition, s)
		if isTruthy(elseifValue) {
			var result values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}
			for _, statement := range elseifClause.Body {
				result = Evaluate(statement, s)
			}
			return result
		}
	}

	// Execute else body if it exists
	if len(node.ElseBody) > 0 {
		var result values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}
		for _, statement := range node.ElseBody {
			result = Evaluate(statement, s)
		}
		return result
	}

	// Return null if no condition was met and no else clause
	return &values.NullValue{Type: parser.NodeTypeNull}
}

// Helper function to determine if a value is truthy
func isTruthy(value values.RuntimeValue) bool {
	switch v := value.(type) {
	case *values.BooleanValue:
		return v.Value
	case *values.NumericValue:
		return v.Value != 0
	case *values.StringValue:
		return v.Value != ""
	case *values.NullValue:
		return false
	default:
		return true // Objects, functions, etc. are truthy
	}
}

func evaluateLoopStatement(node *parser.LoopStatement, s *scope.Scope) values.RuntimeValue {
	var result values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}

	// Check if this is a for-each loop
	if node.IsForEach {
		return evaluateForEachLoop(node, s)
	}

	// Check if this is a range loop
	if node.LoopVar != "" {
		return evaluateRangeLoop(node, s)
	}

	// Check if this is an infinite loop (no condition)
	if node.Condition == nil {
		// Infinite loop - treat condition as always true
		for {
			// Execute loop body
			for _, statement := range node.Body {
				result = Evaluate(statement, s)
				// Check if break was executed
				if result.NodeType() == parser.NodeTypeBreakExpression {
					return &values.NullValue{Type: parser.NodeTypeNull}
				}
			}
		}
	}

	// Loop with condition
	// Evaluate the condition
	conditionValue := Evaluate(node.Condition, s)

	// Continue looping while the condition is truthy
	for isTruthy(conditionValue) {
		// Execute loop body
		for _, statement := range node.Body {
			result = Evaluate(statement, s)
			// Check if break was executed
			if result.NodeType() == parser.NodeTypeBreakExpression {
				return &values.NullValue{Type: parser.NodeTypeNull}
			}
		}

		// Re-evaluate the condition to check if we should continue
		conditionValue = Evaluate(node.Condition, s)
	}

	return result
}

// evaluateRangeLoop executes a range-based loop (loop i from X to Y)
func evaluateRangeLoop(node *parser.LoopStatement, s *scope.Scope) values.RuntimeValue {
	fromValue := Evaluate(node.From, s)
	toValue := Evaluate(node.To, s)

	// Validate types
	if fromValue.NodeType() != parser.NodeTypeNumeric || toValue.NodeType() != parser.NodeTypeNumeric {
		fmt.Printf("Range loop requires numeric values for 'from' and 'to'\n")
		os.Exit(1)
		return nil
	}

	fromNumeric := fromValue.(*values.NumericValue)
	toNumeric := toValue.(*values.NumericValue)

	// Determine increment (default is 1)
	increment := 1.0
	if node.Increment != nil {
		incValue := Evaluate(node.Increment, s)
		if incValue.NodeType() != parser.NodeTypeNumeric {
			fmt.Printf("Range loop increment must be numeric\n")
			os.Exit(1)
			return nil
		}
		increment = incValue.(*values.NumericValue).Value
	}

	// Check if loop variable already exists, if not declare it
	// We need to manually check since Declare will exit if variable exists
	scopeVars := s.GetVariables()
	_, exists := scopeVars[node.LoopVar]
	if !exists {
		scopeVars[node.LoopVar] = &values.NumericValue{Type: parser.NodeTypeNumeric, Value: fromNumeric.Value}
	} else {
		// Variable exists, just update its value
		scopeVars[node.LoopVar] = &values.NumericValue{Type: parser.NodeTypeNumeric, Value: fromNumeric.Value}
	}

	var result values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}

	// Determine loop direction based on increment sign if provided, otherwise from/to values
	// If increment is provided, it determines direction. Otherwise default to forward if from <= to
	goingForward := true
	if node.Increment != nil {
		goingForward = increment > 0
	} else {
		goingForward = fromNumeric.Value <= toNumeric.Value
	}

	current := fromNumeric.Value

	// Execute loop
	for {
		// Check termination based on direction
		if goingForward && current > toNumeric.Value {
			break
		}
		if !goingForward && current < toNumeric.Value {
			break
		}

		// Update loop variable value
		s.Assign(node.LoopVar, &values.NumericValue{Type: parser.NodeTypeNumeric, Value: current})

		// Execute loop body
		for _, statement := range node.Body {
			result = Evaluate(statement, s)
			// Check if break was executed
			if result.NodeType() == parser.NodeTypeBreakExpression {
				return &values.NullValue{Type: parser.NodeTypeNull}
			}
		}

		current += increment
	}

	return result
}

// evaluateForEachLoop executes a for-each loop (loop element from arr { })
func evaluateForEachLoop(node *parser.LoopStatement, s *scope.Scope) values.RuntimeValue {
	// Evaluate the iterable (should be an array)
	iterableValue := Evaluate(node.From, s)

	// Validate that it's an array
	if iterableValue.NodeType() != parser.NodeTypeArray {
		fmt.Printf("For-each loop requires an array, got %s\n", iterableValue.NodeType())
		os.Exit(1)
		return nil
	}

	arrayValue := iterableValue.(*values.ArrayValue)
	var result values.RuntimeValue = &values.NullValue{Type: parser.NodeTypeNull}

	// Iterate over each element in the array
	for _, element := range arrayValue.Elements {
		// Check if loop variable already exists, if not declare it
		scopeVars := s.GetVariables()
		_, exists := scopeVars[node.LoopVar]
		if !exists {
			scopeVars[node.LoopVar] = element
		} else {
			// Variable exists, just update its value
			scopeVars[node.LoopVar] = element
		}

		// Execute loop body
		for _, statement := range node.Body {
			result = Evaluate(statement, s)
			// Check if break was executed
			if result.NodeType() == parser.NodeTypeBreakExpression {
				return &values.NullValue{Type: parser.NodeTypeNull}
			}
		}
	}

	return result
}

func evaluateBreakExpression(_ *parser.BreakExpression, _ *scope.Scope) values.RuntimeValue {
	return &values.BreakValue{Type: parser.NodeTypeBreakExpression}
}

func evaluateReturnStatement(node *parser.ReturnStatement, s *scope.Scope) values.RuntimeValue {
	// If return has no value, return null
	if node.Value == nil {
		return &values.ReturnValue{
			Type:  parser.NodeTypeReturnValue,
			Value: &values.NullValue{Type: parser.NodeTypeNull},
		}
	}

	// Evaluate the return value
	value := Evaluate(node.Value, s)

	return &values.ReturnValue{
		Type:  parser.NodeTypeReturnValue,
		Value: value,
	}
}
