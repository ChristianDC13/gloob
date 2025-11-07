package errors

import (
	"fmt"
	"gloob-interpreter/internal/colors"
	"gloob-interpreter/internal/lexer"
	"os"
	"strings"
)

// Error message constants for parser (syntax errors)
const (
	ErrExpectedIdentifier      = "An identifier was expected here dude 😎"
	ErrExpectedEqual           = "An equal sign was expected here dude 😎"
	ErrExpectedColon           = "A colon was expected after the key :v"
	ErrExpectedOpenCurly       = "Expected opening curly brackets"
	ErrExpectedCloseCurly      = "Expected closing curly brackets"
	ErrExpectedOpenParen       = "Expected opening parentheses"
	ErrExpectedCloseParen      = "Expected closing parentheses"
	ErrExpectedCloseSquare     = "Expected closing square brackets"
	ErrExpectedFunctionName    = "Expected function name"
	ErrConstMustHaveValue      = "A constant declaration must have a value 🤔"
	ErrExpectedIdentifierParam = "Expected an identifier here 👀"
	ErrUnexpectedToken         = "Unexpected token '%s'. Are you sure you typed it correctly? 🤔"
)

// Error message constants for runtime (interpreter errors)
const (
	ErrVariableNotFound           = "Variable '%s' not found. Are you sure you typed it correctly? 🤔"
	ErrVariableAlreadyDeclared    = "Variable '%s' already declared"
	ErrVariableNotInitialized     = "Variable '%s' is not initialized. Are you sure you declared it? 🤔"
	ErrConstantCannotBeAssigned   = "Constant '%s' cannot be assigned to because it is, how can i say it to you? It is a constant 😒"
	ErrDivisionByZero             = "You know you cannot divide by zero, what are you trying to prove? 😒"
	ErrUnknownOperator            = "Unknown operator: '%s', i don't know what to tell you 🫣"
	ErrUnknownOperatorWithString  = "Unknown operator: '%s', with string operands"
	ErrInvalidOperandTypes        = "Invalid operand types for binary expression: %s %s %s"
	ErrInvalidLeftOperand         = "Invalid left operand type for binary expression: %s"
	ErrInvalidRightOperand        = "Invalid right operand type for binary expression: %s"
	ErrCannotAccessProperty       = "Cannot access property '%s' on non-object type: %s"
	ErrPropertyNotFound           = "Property '%s' not found on object"
	ErrCannotAssignProperty       = "Cannot assign property '%s' on non-object type: %s"
	ErrCannotIndexNonArray        = "Cannot index non-array type: %s"
	ErrIndexMustBeNumeric         = "Index must be numeric"
	ErrArrayIndexOutOfBounds      = "Array index out of bounds: %d (array length: %d)"
	ErrStringIndexOutOfBounds     = "String index out of bounds: %d (string length: %d)"
	ErrCannotIndexType            = "Cannot index type: %s"
	ErrInvalidNativeFunction      = "Invalid native function type"
	ErrFunctionArgCountMismatch   = "Function '%s' expects %d arguments, got %d"
	ErrCannotCallNonFunction      = "Cannot call non-function value: %s"
	ErrUnknownNodeType            = "Unknown node type: '%s', i don't know what to tell you 🫣"
	ErrRangeLoopNeedsNumeric      = "Range loop requires numeric values for 'from' and 'to'"
	ErrRangeLoopIncrementNumeric  = "Range loop increment must be numeric"
	ErrForEachNeedsArray          = "For-each loop requires an array, got %s"
	ErrCannotCompareTypes         = "Cannot compare %s and %s with operator %s"
	ErrUnknownComparisonOperator  = "Unknown comparison operator: %s"
	ErrUnknownLogicalOperator     = "Unknown logical operator: %s"
	ErrCannotUseOperatorWithNull  = "Cannot use operator %s with null values"
	ErrInvalidIdentifierForAssign = "Invalid identifier type for variable assignment: %s"
)

// SyntaxError prints a detailed syntax error with file context and exits.
func SyntaxError(token lexer.Token, sourceCode string, message string) {
	// Print the error header with file location
	fmt.Printf("\n%s %s\n", colors.Red("Syntax Error:"), message)

	if token.Filename != "" {
		fmt.Printf("%s  at %s:%d:%d\n", colors.Blue("-->"), token.Filename, token.Line, token.ColumnStart)
	} else {
		fmt.Printf("%s  at line %d, column %d\n", colors.Blue("-->"), token.Line, token.ColumnStart)
	}

	// Get the line from source code
	lines := strings.Split(sourceCode, "\n")
	if token.Line > 0 && token.Line <= len(lines) {
		lineContent := lines[token.Line-1]

		// Print line number and content
		fmt.Printf("%s\n", colors.Blue(fmt.Sprintf("   %d | ", token.Line)))
		fmt.Printf("   %d | %s\n", token.Line, lineContent)

		// Print the pointer to the error location
		padding := strings.Repeat(" ", token.ColumnStart-1)
		underline := strings.Repeat("^", max(1, token.ColumnEnd-token.ColumnStart+1))
		fmt.Printf("%s %s%s\n", colors.Blue("     |"), padding, colors.Red(underline))
	}

	fmt.Println()
	os.Exit(1)
}

// RuntimeError prints a detailed runtime error with file context if available and exits.
func RuntimeError(token *lexer.Token, sourceCode string, message string) {
	// Print the error header
	fmt.Printf("\n%s %s\n", colors.Red("Runtime Error:"), message)

	// If we have token information, show file location
	if token != nil {
		if token.Filename != "" {
			fmt.Printf("%s  at %s:%d:%d\n", colors.Blue("-->"), token.Filename, token.Line, token.ColumnStart)
		} else {
			fmt.Printf("%s  at line %d, column %d\n", colors.Blue("-->"), token.Line, token.ColumnStart)
		}

		// Get the line from source code if available
		if sourceCode != "" {
			lines := strings.Split(sourceCode, "\n")
			if token.Line > 0 && token.Line <= len(lines) {
				lineContent := lines[token.Line-1]

				// Print line number and content
				fmt.Printf("%s\n", colors.Blue(fmt.Sprintf("   %d | ", token.Line)))
				fmt.Printf("   %d | %s\n", token.Line, lineContent)

				// Print the pointer to the error location
				padding := strings.Repeat(" ", token.ColumnStart-1)
				underline := strings.Repeat("^", max(1, token.ColumnEnd-token.ColumnStart+1))
				fmt.Printf("%s %s%s\n", colors.Blue("     |"), padding, colors.Red(underline))
			}
		}
	}

	fmt.Println()
	os.Exit(1)
}

// Helper function to get max of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
