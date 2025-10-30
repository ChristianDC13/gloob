package parser

import "fmt"

// NodeType represents the type of an AST node.
// This is used to identify what kind of language construct a node represents
// and to route it to the appropriate evaluator in the runtime.
type NodeType string

const (
	// Program node represents the root of the AST containing all statements
	NodeTypeProgram NodeType = "PROGRAM"

	// Literal value nodes
	NodeTypeNumeric NodeType = "NUMERIC" // Number literals (e.g., 42, 3.14)
	NodeTypeBoolean NodeType = "BOOLEAN" // Boolean literals (true, false)
	NodeTypeString  NodeType = "STRING"  // String literals (e.g., "hello")
	NodeTypeNull    NodeType = "NULL"    // Null value

	// Identifier and expression nodes
	NodeTypeIdentifier       NodeType = "IDENTIFIER"        // Variable/function names
	NodeTypeBinaryExpression NodeType = "BINARY_EXPRESSION" // Binary operations (+, -, *, /, ==, etc.)
	NodeTypeUnaryExpression  NodeType = "UNARY_EXPRESSION"  // Unary operations (not implemented yet)

	// Object-related nodes
	NodeTypeObject       NodeType = "OBJECT"        // Object literals { key: value }
	NodeTypeProperty     NodeType = "PROPERTY"      // Object properties
	NodeTypeMemberAccess NodeType = "MEMBER_ACCESS" // Property access (obj.property)

	// Function-related nodes
	NodeTypeCallExpression      NodeType = "CALL_EXPRESSION"      // Function calls func(args)
	NodeTypeFunctionDeclaration NodeType = "FUNCTION_DECLARATION" // Function definitions
	NodeTypeNativeFunction      NodeType = "NATIVE_FUNCTION"      // Built-in functions

	// Variable-related nodes
	NodeTypeVariableDeclaration NodeType = "VARIABLE_DECLARATION" // var/const declarations
	NodeTypeVariableAssignment  NodeType = "VARIABLE_ASSIGNMENT"  // Variable assignments (var = value)

	// Control flow nodes
	NodeTypeIfStatement     NodeType = "IF_STATEMENT"     // if statements
	NodeTypeElseIfClause    NodeType = "ELSE_IF_CLAUSE"   // elseif clauses
	NodeTypeLoopStatement   NodeType = "LOOP_STATEMENT"   // loop statements
	NodeTypeBreakExpression NodeType = "BREAK_EXPRESSION" // break statements
	NodeTypeReturnStatement NodeType = "RETURN_STATEMENT" // return statements
	NodeTypeReturnValue     NodeType = "RETURN_VALUE"     // return value (runtime marker)

	// Import nodes
	NodeTypeImportStatement NodeType = "IMPORT_STATEMENT" // import statements

	// Collection nodes
	NodeTypeArray      NodeType = "ARRAY"       // Array literals [1, 2, 3]
	NodeTypeArrayIndex NodeType = "ARRAY_INDEX" // Array indexing arr[1]
	NodeTypeCollection NodeType = "COLLECTION"  // Generic collections (not implemented)
)

// Statement represents any executable statement in the language.
// All statements must implement the NodeType() method for runtime dispatch.
type Statement interface {
	NodeType() NodeType
}

// Expression represents any expression that evaluates to a value.
// Expressions can be used as values in assignments, function calls, etc.
type Expression interface {
	NodeType() NodeType
}

// Program is the root node of the AST.
// It contains all the statements that make up a Gloob program.
type Program struct {
	Statements []Statement // All statements in the program
}

func (p *Program) NodeType() NodeType {
	return NodeTypeProgram
}

// VariableDeclaration represents variable and constant declarations.
// Examples: var name = "value", const PI = 3.14
type VariableDeclaration struct {
	Constant   bool       // true for const, false for var
	Identifier string     // Variable name
	Value      Expression // Initial value (can be nil for var without assignment)
}

func (v *VariableDeclaration) NodeType() NodeType {
	return NodeTypeVariableDeclaration
}

// BinaryExpression represents binary operations like arithmetic and comparison.
// Examples: a + b, x > y, name == "test"
type BinaryExpression struct {
	Type     NodeType   `json:"type"`     // Node type (always BINARY_EXPRESSION)
	Left     Expression `json:"left"`     // Left operand
	Operator string     `json:"operator"` // Operator (+, -, *, /, ==, !=, >, <, etc.)
	Right    Expression `json:"right"`    // Right operand
}

func (b *BinaryExpression) NodeType() NodeType {
	return NodeTypeBinaryExpression
}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator, b.Right)
}

// Identifier represents variable and function names.
// Examples: name, age, calculateSum
type Identifier struct {
	Type NodeType `json:"type"` // Node type (always IDENTIFIER)
	Name string   `json:"name"` // The identifier name
}

func (i *Identifier) NodeType() NodeType {
	return NodeTypeIdentifier
}

func (i *Identifier) String() string {
	return i.Name
}

// Numeric represents number literals.
// Examples: 42, 3.14, -10
type Numeric struct {
	Type  NodeType `json:"type"`  // Node type (always NUMERIC)
	Value float64  `json:"value"` // The numeric value
}

func (n *Numeric) NodeType() NodeType {
	return NodeTypeNumeric
}

func (n *Numeric) String() string {
	return fmt.Sprintf("%g", n.Value)
}

// Null represents the null value.
// Used when a variable is declared without initialization or explicitly set to null.
type Null struct {
	// No fields needed - null is just a marker
}

func (n *Null) NodeType() NodeType {
	return NodeTypeNull
}

func (n *Null) String() string {
	return "null"
}

// Boolean represents boolean literals.
// Examples: true, false, yes, no, on, off
type Boolean struct {
	Value bool // The boolean value
}

func (b *Boolean) NodeType() NodeType {
	return NodeTypeBoolean
}

func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// String represents string literals.
// Examples: "hello", 'world', "multi-line string"
type String struct {
	Type  NodeType `json:"type"`  // Node type (always STRING)
	Value string   `json:"value"` // The string content
}

func (s *String) NodeType() NodeType {
	return NodeTypeString
}

func (s *String) String() string {
	return s.Value
}

// VariableAssignmentExpression represents assignment operations.
// Examples: name = "value", obj.property = 42
type VariableAssignmentExpression struct {
	Identifier Expression // Can be Identifier or MemberAccess
	Value      Expression // The value being assigned
}

func (a *VariableAssignmentExpression) NodeType() NodeType {
	return NodeTypeVariableAssignment
}

func (a *VariableAssignmentExpression) String() string {
	return fmt.Sprintf("%s = %s", a.Identifier, a.Value)
}

// Object represents object literals.
// Examples: { name: "John", age: 30 }, { }
type Object struct {
	Properties []Property `json:"properties"` // List of key-value pairs
}

func (o *Object) NodeType() NodeType {
	return NodeTypeObject
}

func (o *Object) String() string {
	return fmt.Sprintf("{%s}", o.Properties)
}

// Property represents a key-value pair in an object.
// Examples: name: "John", age: 30
type Property struct {
	Key   string     `json:"key"`   // Property name
	Value Expression `json:"value"` // Property value
}

func (p *Property) NodeType() NodeType {
	return NodeTypeProperty
}

func (p *Property) String() string {
	return fmt.Sprintf("%s: %s", p.Key, p.Value)
}

// MemberAccess represents property access on objects.
// Examples: obj.name, person.address.city
type MemberAccess struct {
	Object   Expression // The object being accessed
	Property string     // The property name
}

func (m *MemberAccess) NodeType() NodeType {
	return NodeTypeMemberAccess
}

func (m *MemberAccess) String() string {
	return fmt.Sprintf("%s.%s", m.Object, m.Property)
}

// CallExpression represents function calls.
// Examples: print("hello"), add(5, 3), obj.method()
type CallExpression struct {
	Type   NodeType     `json:"type"`   // Node type (always CALL_EXPRESSION)
	Callee Expression   `json:"callee"` // Function being called (Identifier or MemberAccess)
	Args   []Expression `json:"args"`   // Function arguments
}

func (c *CallExpression) NodeType() NodeType {
	return NodeTypeCallExpression
}

func (c *CallExpression) String() string {
	return fmt.Sprintf("%s(%s)", c.Callee, c.Args)
}

// FunctionDeclaration represents function definitions.
// Examples: function greet(name) { return "Hello " + name }
type FunctionDeclaration struct {
	Identifier string      // Function name
	Parameters []string    // Parameter names
	Body       []Statement // Function body statements
}

func (f *FunctionDeclaration) NodeType() NodeType {
	return NodeTypeFunctionDeclaration
}

func (f *FunctionDeclaration) String() string {
	return fmt.Sprintf("function %s(%s) { %s }", f.Identifier, f.Parameters, f.Body)
}

// ElseIfClause represents elseif conditions in if statements.
// Examples: elseif (age >= 13) { print("Teenager") }
type ElseIfClause struct {
	Condition Expression  // The condition to evaluate
	Body      []Statement // Statements to execute if condition is true
}

func (e *ElseIfClause) NodeType() NodeType {
	return NodeTypeElseIfClause
}

func (e *ElseIfClause) String() string {
	return fmt.Sprintf("elseif %s { %s }", e.Condition, e.Body)
}

// IfStatement represents conditional execution.
// Examples: if (age >= 18) { print("Adult") } else { print("Minor") }
type IfStatement struct {
	Condition Expression     // The condition to evaluate
	Body      []Statement    // Statements to execute if condition is true
	ElseIfs   []ElseIfClause // Additional elseif conditions
	ElseBody  []Statement    // Statements to execute if all conditions are false
}

func (i *IfStatement) NodeType() NodeType {
	return NodeTypeIfStatement
}

func (i *IfStatement) String() string {
	return fmt.Sprintf("if %s { %s }", i.Condition, i.Body)
}

// LoopStatement represents different types of loops.
// Examples: loop condition { do something }, loop { infinite loop },
//
//	loop i from 1 to 100 { }, loop i from 0 to 10; 2 { }
type LoopStatement struct {
	Condition Expression  // The condition to evaluate (nil for infinite/range/for-each loops)
	Body      []Statement // Statements to execute

	// Range loop fields (nil for condition-based/for-each loops)
	LoopVar   string     // Loop variable name (e.g., "i" for range, "element" for for-each)
	From      Expression // Start value for range loop OR iterable for for-each loop
	To        Expression // End value for range loop (nil for for-each)
	Increment Expression // Optional increment (nil means increment by 1, only for range loops)

	// For-each loop indicator
	IsForEach bool // True if this is a for-each loop (loop element from arr)
}

func (l *LoopStatement) NodeType() NodeType {
	return NodeTypeLoopStatement
}

func (l *LoopStatement) String() string {
	if l.LoopVar != "" {
		// Range loop
		return fmt.Sprintf("loop %s from %s to %s { %s }", l.LoopVar, l.From, l.To, l.Body)
	}
	return fmt.Sprintf("loop %s { %s }", l.Condition, l.Body)
}

// BreakExpression represents break statements.
// Examples: break
type BreakExpression struct {
}

func (b *BreakExpression) NodeType() NodeType {
	return NodeTypeBreakExpression
}

func (b *BreakExpression) String() string {
	return "break"
}

// ReturnStatement represents return statements.
// Examples: return, return value, return x + y
type ReturnStatement struct {
	Value Expression // The value to return (nil for bare "return")
}

func (r *ReturnStatement) NodeType() NodeType {
	return NodeTypeReturnStatement
}

func (r *ReturnStatement) String() string {
	if r.Value == nil {
		return "return"
	}
	return fmt.Sprintf("return %s", r.Value)
}

// ImportStatement represents an import declaration.
// Example: import "utils/helpers"
type ImportStatement struct {
	Path string // The path to the file to import
}

func (i *ImportStatement) NodeType() NodeType {
	return NodeTypeImportStatement
}

func (i *ImportStatement) String() string {
	return fmt.Sprintf("import \"%s\"", i.Path)
}

// Array represents array literals.
// Examples: [1, 2, 3], ["hello", "world"]
type Array struct {
	Elements []Expression // Elements in the array
}

func (a *Array) NodeType() NodeType {
	return NodeTypeArray
}

func (a *Array) String() string {
	return fmt.Sprintf("[%v]", a.Elements)
}

// ArrayIndex represents array element access.
// Examples: arr[1], arr[i + 1]
type ArrayIndex struct {
	ArrayExpression Expression // The array expression
	Index           Expression // The index expression
}

func (a *ArrayIndex) NodeType() NodeType {
	return NodeTypeArrayIndex
}

func (a *ArrayIndex) String() string {
	return fmt.Sprintf("%s[%s]", a.ArrayExpression, a.Index)
}
