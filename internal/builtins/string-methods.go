package builtins

import (
	"fmt"
	"gloob-interpreter/internal/parser"
	"gloob-interpreter/internal/values"
	"os"
	"strings"
)

// StringLenMethod returns the length of a string
func StringLenMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			return &values.NumericValue{
				Type:  parser.NodeTypeNumeric,
				Value: float64(len(str.Value)),
			}
		},
	}
}

// StringUpperMethod converts a string to uppercase
func StringUpperMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			return &values.StringValue{
				Type:  parser.NodeTypeString,
				Value: strings.ToUpper(str.Value),
			}
		},
	}
}

// StringLowerMethod converts a string to lowercase
func StringLowerMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			return &values.StringValue{
				Type:  parser.NodeTypeString,
				Value: strings.ToLower(str.Value),
			}
		},
	}
}

// StringTrimMethod removes leading and trailing whitespace
func StringTrimMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			return &values.StringValue{
				Type:  parser.NodeTypeString,
				Value: strings.TrimSpace(str.Value),
			}
		},
	}
}

// StringContainsMethod checks if a string contains a substring
func StringContainsMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("contains() expects 1 argument, got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			if args[0].NodeType() != parser.NodeTypeString {
				fmt.Printf("contains() expects a string argument\n")
				os.Exit(1)
				return nil
			}
			substring := args[0].(*values.StringValue).Value
			return &values.BooleanValue{
				Type:  parser.NodeTypeBoolean,
				Value: strings.Contains(str.Value, substring),
			}
		},
	}
}

// StringSplitMethod splits a string into an array by a separator
func StringSplitMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("split() expects 1 argument (separator), got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			if args[0].NodeType() != parser.NodeTypeString {
				fmt.Printf("split() expects a string separator\n")
				os.Exit(1)
				return nil
			}
			separator := args[0].(*values.StringValue).Value
			parts := strings.Split(str.Value, separator)

			// Convert string parts to RuntimeValue array
			elements := make([]values.RuntimeValue, len(parts))
			for i, part := range parts {
				elements[i] = &values.StringValue{
					Type:  parser.NodeTypeString,
					Value: part,
				}
			}

			return &values.ArrayValue{
				Type:     parser.NodeTypeArray,
				Elements: elements,
			}
		},
	}
}

// StringReplaceMethod replaces all occurrences of a substring with another
func StringReplaceMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 2 {
				fmt.Printf("replace() expects 2 arguments (old, new), got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			if args[0].NodeType() != parser.NodeTypeString || args[1].NodeType() != parser.NodeTypeString {
				fmt.Printf("replace() expects string arguments\n")
				os.Exit(1)
				return nil
			}
			oldStr := args[0].(*values.StringValue).Value
			newStr := args[1].(*values.StringValue).Value
			return &values.StringValue{
				Type:  parser.NodeTypeString,
				Value: strings.ReplaceAll(str.Value, oldStr, newStr),
			}
		},
	}
}

// StringIndexOfMethod returns the 1-based index of the first occurrence of a substring
// Returns 0 if not found (to stay consistent with 1-based indexing)
func StringIndexOfMethod(str *values.StringValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("indexOf() expects 1 argument (substring), got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			if args[0].NodeType() != parser.NodeTypeString {
				fmt.Printf("indexOf() expects a string argument\n")
				os.Exit(1)
				return nil
			}
			substring := args[0].(*values.StringValue).Value
			index := strings.Index(str.Value, substring)

			// Convert 0-based to 1-based (or 0 if not found)
			if index == -1 {
				index = 0
			} else {
				index = index + 1
			}

			return &values.NumericValue{
				Type:  parser.NodeTypeNumeric,
				Value: float64(index),
			}
		},
	}
}

// GetStringMethod returns the appropriate string method as a native function
func GetStringMethod(str *values.StringValue, methodName string) values.RuntimeValue {
	switch methodName {
	case "len":
		return StringLenMethod(str)
	case "upper":
		return StringUpperMethod(str)
	case "lower":
		return StringLowerMethod(str)
	case "trim":
		return StringTrimMethod(str)
	case "contains":
		return StringContainsMethod(str)
	case "split":
		return StringSplitMethod(str)
	case "replace":
		return StringReplaceMethod(str)
	case "indexOf":
		return StringIndexOfMethod(str)
	default:
		fmt.Printf("Unknown string method: %s\n", methodName)
		os.Exit(1)
		return nil
	}
}
