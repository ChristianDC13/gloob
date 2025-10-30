package builtins

import (
	"fmt"
	"gloob-interpreter/internal/parser"
	"gloob-interpreter/internal/values"
	"os"
)

// ArrayPushMethod adds an element to the end of an array
func ArrayPushMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("push() expects 1 argument, got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			array.Elements = append(array.Elements, args[0])
			return array
		},
	}
}

// ArrayPopMethod removes and returns the last element of an array
func ArrayPopMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(array.Elements) == 0 {
				fmt.Printf("Cannot pop from empty array\n")
				os.Exit(1)
				return nil
			}
			lastIndex := len(array.Elements) - 1
			lastElement := array.Elements[lastIndex]
			array.Elements = array.Elements[:lastIndex]
			return lastElement
		},
	}
}

// ArrayLenMethod returns the length of an array
func ArrayLenMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			return &values.NumericValue{
				Type:  parser.NodeTypeNumeric,
				Value: float64(len(array.Elements)),
			}
		},
	}
}

// ArrayRemoveMethod removes an element at the specified index (1-based)
func ArrayRemoveMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("remove() expects 1 argument (index), got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			if args[0].NodeType() != parser.NodeTypeNumeric {
				fmt.Printf("remove() expects numeric index\n")
				os.Exit(1)
				return nil
			}
			index := int(args[0].(*values.NumericValue).Value)
			// Convert 1-based to 0-based
			index = index - 1
			if index < 0 || index >= len(array.Elements) {
				fmt.Printf("Array index out of bounds: %d\n", index+1)
				os.Exit(1)
				return nil
			}
			array.Elements = append(array.Elements[:index], array.Elements[index+1:]...)
			return array
		},
	}
}

// ArrayInsertMethod inserts an element at the specified index (1-based)
func ArrayInsertMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 2 {
				fmt.Printf("insert() expects 2 arguments (index, value), got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			if args[0].NodeType() != parser.NodeTypeNumeric {
				fmt.Printf("insert() expects numeric index\n")
				os.Exit(1)
				return nil
			}
			index := int(args[0].(*values.NumericValue).Value)
			// Convert 1-based to 0-based
			index = index - 1
			if index < 0 || index > len(array.Elements) {
				fmt.Printf("Array index out of bounds: %d\n", index+1)
				os.Exit(1)
				return nil
			}
			// Insert element at index
			array.Elements = append(array.Elements[:index], append([]values.RuntimeValue{args[1]}, array.Elements[index:]...)...)
			return array
		},
	}
}

// ArrayIndexOfMethod returns the 1-based index of the first occurrence of an element
// Returns 0 if not found (to stay consistent with 1-based indexing)
func ArrayIndexOfMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("indexOf() expects 1 argument (element), got %d\n", len(args))
				os.Exit(1)
				return nil
			}

			searchValue := args[0]
			for i, element := range array.Elements {
				// Simple equality check based on type and value
				if elementsEqual(element, searchValue) {
					// Return 1-based index
					return &values.NumericValue{
						Type:  parser.NodeTypeNumeric,
						Value: float64(i + 1),
					}
				}
			}

			// Not found, return 0
			return &values.NumericValue{
				Type:  parser.NodeTypeNumeric,
				Value: 0,
			}
		},
	}
}

// ArrayContainsMethod checks if an array contains an element
func ArrayContainsMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("contains() expects 1 argument (element), got %d\n", len(args))
				os.Exit(1)
				return nil
			}

			searchValue := args[0]
			for _, element := range array.Elements {
				if elementsEqual(element, searchValue) {
					return &values.BooleanValue{
						Type:  parser.NodeTypeBoolean,
						Value: true,
					}
				}
			}

			return &values.BooleanValue{
				Type:  parser.NodeTypeBoolean,
				Value: false,
			}
		},
	}
}

// ArrayJoinMethod joins array elements into a string with a separator
func ArrayJoinMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			if len(args) != 1 {
				fmt.Printf("join() expects 1 argument (separator), got %d\n", len(args))
				os.Exit(1)
				return nil
			}
			if args[0].NodeType() != parser.NodeTypeString {
				fmt.Printf("join() expects a string separator\n")
				os.Exit(1)
				return nil
			}

			separator := args[0].(*values.StringValue).Value

			if len(array.Elements) == 0 {
				return &values.StringValue{
					Type:  parser.NodeTypeString,
					Value: "",
				}
			}

			// Build the joined string
			result := fmt.Sprint(array.Elements[0])
			for i := 1; i < len(array.Elements); i++ {
				result += separator + fmt.Sprint(array.Elements[i])
			}

			return &values.StringValue{
				Type:  parser.NodeTypeString,
				Value: result,
			}
		},
	}
}

// ArrayReverseMethod reverses the array in-place
func ArrayReverseMethod(array *values.ArrayValue) *values.NativeFunctionValue {
	return &values.NativeFunctionValue{
		Type: parser.NodeTypeNativeFunction,
		Expression: func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
			// Reverse the array in-place
			for i, j := 0, len(array.Elements)-1; i < j; i, j = i+1, j-1 {
				array.Elements[i], array.Elements[j] = array.Elements[j], array.Elements[i]
			}
			return array
		},
	}
}

// elementsEqual checks if two RuntimeValues are equal
func elementsEqual(a, b values.RuntimeValue) bool {
	if a.NodeType() != b.NodeType() {
		return false
	}

	switch a.NodeType() {
	case parser.NodeTypeNumeric:
		return a.(*values.NumericValue).Value == b.(*values.NumericValue).Value
	case parser.NodeTypeString:
		return a.(*values.StringValue).Value == b.(*values.StringValue).Value
	case parser.NodeTypeBoolean:
		return a.(*values.BooleanValue).Value == b.(*values.BooleanValue).Value
	case parser.NodeTypeNull:
		return true
	default:
		// For complex types, use pointer comparison
		return a == b
	}
}

// GetArrayMethod returns the appropriate array method as a native function
func GetArrayMethod(array *values.ArrayValue, methodName string) values.RuntimeValue {
	switch methodName {
	case "push":
		return ArrayPushMethod(array)
	case "pop":
		return ArrayPopMethod(array)
	case "len":
		return ArrayLenMethod(array)
	case "remove":
		return ArrayRemoveMethod(array)
	case "insert":
		return ArrayInsertMethod(array)
	case "indexOf":
		return ArrayIndexOfMethod(array)
	case "contains":
		return ArrayContainsMethod(array)
	case "join":
		return ArrayJoinMethod(array)
	case "reverse":
		return ArrayReverseMethod(array)
	default:
		fmt.Printf("Unknown array method: %s\n", methodName)
		os.Exit(1)
		return nil
	}
}
