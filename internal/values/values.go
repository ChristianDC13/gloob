package values

import (
	"fmt"
	"gloob-interpreter/internal/colors"
	"gloob-interpreter/internal/parser"
)

// RuntimeValue is the interface that all runtime values must implement.
// This allows the interpreter to work with different types of values uniformly.
// All runtime values must be able to report their type for type checking and dispatch.
type RuntimeValue interface {
	NodeType() parser.NodeType
}

// NumericValue represents number values at runtime.
// Examples: 42, 3.14, -10
type NumericValue struct {
	Type  parser.NodeType `json:"type"`  // Always NodeTypeNumeric
	Value float64         `json:"value"` // The actual numeric value
}

func (n *NumericValue) NodeType() parser.NodeType {
	return parser.NodeTypeNumeric
}

func (n *NumericValue) String() string {
	return fmt.Sprintf("%g", n.Value)
}

// BooleanValue represents boolean values at runtime.
// Examples: true, false
type BooleanValue struct {
	Type  parser.NodeType `json:"type"`  // Always NodeTypeBoolean
	Value bool            `json:"value"` // The actual boolean value
}

func (b *BooleanValue) NodeType() parser.NodeType {
	return parser.NodeTypeBoolean
}

func (b *BooleanValue) String() string {
	return fmt.Sprintf("%t", b.Value)
}

// NullValue represents the null value at runtime.
// Used when variables are uninitialized or explicitly set to null.
type NullValue struct {
	Type parser.NodeType `json:"type"` // Always NodeTypeNull
}

func (n *NullValue) NodeType() parser.NodeType {
	return parser.NodeTypeNull
}

func (n *NullValue) String() string {
	return "null"
}

// StringValue represents string values at runtime.
// Examples: "hello", "world", ""
type StringValue struct {
	Type  parser.NodeType `json:"type"`  // Always NodeTypeString
	Value string          `json:"value"` // The actual string content
}

func (s *StringValue) NodeType() parser.NodeType {
	return parser.NodeTypeString
}

func (s *StringValue) String() string {
	return s.Value
}

// NodeVariableDeclaration represents a variable declaration at runtime.
// This is used to track variable declarations and their values for debugging/display.
type NodeVariableDeclaration struct {
	Type  parser.NodeType `json:"type"`  // Always NodeTypeVariableDeclaration
	Name  string          `json:"name"`  // Variable name
	Value RuntimeValue    `json:"value"` // Variable value
}

func (n *NodeVariableDeclaration) NodeType() parser.NodeType {
	return parser.NodeTypeVariableDeclaration
}

func (n *NodeVariableDeclaration) String() string {
	return fmt.Sprintf("var %s = %s", n.Name, n.Value)
}

// ObjectValue represents object values at runtime.
// Objects are collections of key-value pairs where keys are strings and values are RuntimeValues.
// Examples: { name: "John", age: 30 }, { nested: { value: 42 } }
type ObjectValue struct {
	Type       parser.NodeType         `json:"type"`       // Always NodeTypeObject
	Properties map[string]RuntimeValue `json:"properties"` // Key-value pairs
}

func (o *ObjectValue) NodeType() parser.NodeType {
	return parser.NodeTypeObject
}

func (o *ObjectValue) String() string {
	return "\n" + o.stringWithIndent(0)
}

// stringWithIndent creates a formatted string representation of the object with proper indentation.
// This is used for pretty-printing objects with nested structures.
func (o *ObjectValue) stringWithIndent(indentLevel int) string {
	indent := ""
	for i := 0; i < indentLevel; i++ {
		indent += "    "
	}

	result := colors.White("{\n")
	first := true
	for key, value := range o.Properties {
		if !first {
			result += ",\n"
		}
		result += indent + "    " + colors.White(fmt.Sprintf("%s: ", key))

		var valueColor = func(value RuntimeValue) string {
			if value.NodeType() == parser.NodeTypeNumeric {
				return colors.Yellow(fmt.Sprintf("%s", value))
			}
			if value.NodeType() == parser.NodeTypeBoolean {
				return colors.Blue(fmt.Sprintf("%s", value))
			}
			if value.NodeType() == parser.NodeTypeNull {
				return colors.Red(fmt.Sprintf("%s", value))
			}

			if value.NodeType() == parser.NodeTypeString {
				return colors.Green(fmt.Sprintf("\"%s\"", value))
			}
			return colors.White(fmt.Sprintf("%s", value))
		}

		// Handle nested objects with proper indentation
		if objValue, ok := value.(*ObjectValue); ok {
			result += objValue.stringWithIndent(indentLevel + 1)
		} else {
			result += valueColor(value)
		}
		first = false
	}
	result += "\n" + indent + "}"
	return result
}

// FunctionValue represents user-defined functions at runtime.
// These are functions defined in Gloob code that can be called with arguments.
// Examples: function greet(name) { return "Hello " + name }
type FunctionValue struct {
	Type       parser.NodeType    `json:"type"`       // Always NodeTypeFunctionDeclaration
	Identifier string             `json:"identifier"` // Function name
	Parameters []string           `json:"parameters"` // Parameter names
	Body       []parser.Statement `json:"body"`       // Function body statements
	Scope      interface{}        `json:"scope"`      // Closure scope (captured variables) - will be set to *scope.Scope
}

func (f *FunctionValue) NodeType() parser.NodeType {
	return parser.NodeTypeFunctionDeclaration
}

func (f *FunctionValue) String() string {
	return fmt.Sprintf("function %s(%s) { %s }", f.Identifier, f.Parameters, f.Body)
}

// CollectionValue represents array/collection values at runtime.
// This is currently not implemented but reserved for future array support.
type CollectionValue struct {
	Type  parser.NodeType `json:"type"`  // Always NodeTypeCollection
	Value []RuntimeValue  `json:"value"` // Array elements
}

func (c *CollectionValue) NodeType() parser.NodeType {
	return parser.NodeTypeCollection
}

func (c *CollectionValue) String() string {
	return fmt.Sprintf("%v", c.Value)
}

// BreakValue is a special runtime value that signals a break statement.
// This is used internally by the interpreter to exit loops.
type BreakValue struct {
	Type parser.NodeType `json:"type"` // Always NodeTypeBreakExpression
}

func (b *BreakValue) NodeType() parser.NodeType {
	return parser.NodeTypeBreakExpression
}

func (b *BreakValue) String() string {
	return "break"
}

// ReturnValue is a special value that signals a return from a function.
// It wraps the actual return value.
type ReturnValue struct {
	Type  parser.NodeType `json:"type"`  // Always NodeTypeReturnValue
	Value RuntimeValue    `json:"value"` // The value being returned
}

func (r *ReturnValue) NodeType() parser.NodeType {
	return parser.NodeTypeReturnValue
}

func (r *ReturnValue) String() string {
	if r.Value == nil {
		return "return"
	}
	return fmt.Sprintf("return %s", r.Value)
}

// ArrayValue represents an array at runtime.
// Arrays are 1-based indexed in Gloob.
type ArrayValue struct {
	Type     parser.NodeType `json:"type"`     // Always NodeTypeArray
	Elements []RuntimeValue  `json:"elements"` // Array elements
}

func (a *ArrayValue) NodeType() parser.NodeType {
	return parser.NodeTypeArray
}

func (a *ArrayValue) String() string {
	return fmt.Sprintf("%v", a.Elements)
}

// NativeFunctionValue represents built-in functions at runtime.
// These are functions implemented in Go that are available globally.
// Examples: print(), type(), len(), input()
type NativeFunctionValue struct {
	Type       parser.NodeType                                           `json:"type"` // Always NodeTypeNativeFunction
	Expression func(args []RuntimeValue, scope interface{}) RuntimeValue // The Go function to call
}

func (n *NativeFunctionValue) NodeType() parser.NodeType {
	return parser.NodeTypeNativeFunction
}

func (n *NativeFunctionValue) String() string {
	return "function"
}
