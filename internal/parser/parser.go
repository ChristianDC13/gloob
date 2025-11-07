package parser

import (
	"fmt"
	"gloob-interpreter/internal/errors"
	"gloob-interpreter/internal/lexer"
	"strconv"
)

// Parser implements a recursive descent parser for the Gloob language.
// It converts a stream of tokens into an Abstract Syntax Tree (AST).
// The parser uses proper operator precedence and handles all language constructs.
type Parser struct {
	tokens     []lexer.Token // Current stream of tokens to parse
	sourceCode string        // Original source code for error reporting
	filename   string        // Filename for error reporting
}

// NewParser creates a new parser instance with the given tokens.
func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// at returns the current token without consuming it.
// This is used for lookahead during parsing.
func (p *Parser) at() lexer.Token {
	return p.tokens[0]
}

// next consumes and returns the current token, advancing to the next one.
func (p *Parser) next() lexer.Token {
	token := p.at()
	p.tokens = p.tokens[1:]
	return token
}

// nextWithExpect consumes the current token and expects it to be of a specific type.
// If the token doesn't match the expected type, it prints an error and exits.
func (p *Parser) nextWithExpect(expected lexer.TokenType, message string) lexer.Token {
	token := p.next()
	if token.Type != expected {
		p.syntaxError(token, message)
		return lexer.Token{}
	}
	return token
}

// syntaxError prints a detailed syntax error with file context and exits.
func (p *Parser) syntaxError(token lexer.Token, message string) {
	errors.SyntaxError(token, p.sourceCode, message)
}

// notEOF checks if there are more tokens to parse.
func (p *Parser) notEOF() bool {
	return p.at().Type != lexer.TokenTypeEOF
}

// ProduceAST is the main entry point for parsing.
// It takes source code, tokenizes it, and produces a complete AST.
func (p *Parser) ProduceAST(sourceCode string) *Program {
	return p.ProduceASTWithFilename(sourceCode, "<stdin>")
}

// ProduceASTWithFilename is like ProduceAST but allows specifying a filename for error reporting.
func (p *Parser) ProduceASTWithFilename(sourceCode string, filename string) *Program {
	// Store source code and filename for error reporting
	p.sourceCode = sourceCode
	p.filename = filename

	// First, tokenize the source code
	p.tokens = lexer.NewLexer(sourceCode, filename).Tokenize()
	program := &Program{
		Statements: []Statement{},
	}

	// Parse all statements until EOF
	for p.notEOF() {
		// Skip newlines (they're not meaningful statements)
		if p.at().Type == lexer.TokenTypeNewline {
			p.next()
			continue
		}
		statement := p.parseStatement()
		program.Statements = append(program.Statements, statement)
	}

	return program
}

// parseStatement is the entry point for parsing statements.
// It determines what type of statement to parse based on the current token.
func (p *Parser) parseStatement() Statement {
	switch p.at().Type {
	case lexer.TokenTypeImport:
		return p.parseImportStatement()
	case lexer.TokenTypeVar, lexer.TokenTypeConst:
		return p.parseVariableDeclaration()
	case lexer.TokenTypeFunction:
		return p.parseFunctionDeclaration()
	case lexer.TokenTypeIf:
		return p.parseIfStatement()
	case lexer.TokenTypeLoop:
		return p.parseLoopStatement()
	case lexer.TokenTypeReturn:
		return p.parseReturnStatement()
	case lexer.TokenTypeComment:
		return p.parseCommentStatement()
	default:
		// If it's not a statement keyword, treat it as an expression
		return p.parseExpression()
	}
}

// parseImportStatement parses import statements.
// Examples: import "utils/helpers", import "math.gloob"
func (p *Parser) parseImportStatement() *ImportStatement {
	p.next() // consume 'import'

	// Expect a string literal with the file path
	pathToken := p.nextWithExpect(lexer.TokenTypeString, "Expected string path after import")

	return &ImportStatement{
		Path: pathToken.Literal,
	}
}
func (p *Parser) parseCommentStatement() *Null {
	for p.notEOF() && p.at().Type != lexer.TokenTypeNewline {
		p.next()
	}
	return &Null{}
}

// parseVariableDeclaration parses variable and constant declarations.
// Examples: var name = "value", const PI = 3.14, var x;
func (p *Parser) parseVariableDeclaration() *VariableDeclaration {
	// Determine if this is a const or var declaration
	isConstant := p.next().Type == lexer.TokenTypeConst
	identifier := p.nextWithExpect(lexer.TokenTypeIdentifier, errors.ErrExpectedIdentifier).Literal

	// Check if this is a declaration without assignment (var x; or var x\n)
	if p.at().Type == lexer.TokenTypeSemicolon || p.at().Type == lexer.TokenTypeNewline {
		if p.at().Type == lexer.TokenTypeSemicolon {
			p.next()
		}
		if isConstant {
			p.syntaxError(p.at(), errors.ErrConstMustHaveValue)
			return nil
		}

		return &VariableDeclaration{
			Constant:   isConstant,
			Identifier: identifier,
			Value:      nil,
		}
	}

	// Parse the assignment part
	p.nextWithExpect(lexer.TokenTypeEqual, errors.ErrExpectedEqual)
	value := p.parseExpression()

	// Skip optional semicolon and newlines
	if p.at().Type == lexer.TokenTypeSemicolon {
		p.next()
	}

	return &VariableDeclaration{
		Constant:   isConstant,
		Identifier: identifier,
		Value:      value,
	}
}

// parseExpression is the entry point for parsing expressions.
// It follows proper operator precedence by delegating to specific precedence levels.
func (p *Parser) parseExpression() Expression {
	return p.parseAssignmentExpression()
}

// parseAssignmentExpression handles assignment operations with lowest precedence.
// Examples: name = "value", obj.property = 42
func (p *Parser) parseAssignmentExpression() Expression {
	left := p.parseLogicalExpression()

	// Check if this is an assignment (right-associative)
	if p.at().Type == lexer.TokenTypeEqual {
		p.next()
		value := p.parseExpression() // Recursive call for right-associativity
		return &VariableAssignmentExpression{
			Identifier: left,
			Value:      value,
		}
	}

	return left
}

// parseLogicalExpression handles logical operators (&& and ||).
// Examples: a && b, x || y
func (p *Parser) parseLogicalExpression() Expression {
	left := p.parseComparisonOnlyExpression()

	// Handle logical operators (left-associative)
	for p.at().Type == lexer.TokenTypeAnd || p.at().Type == lexer.TokenTypeOr {
		operator := p.next().Literal
		right := p.parseComparisonOnlyExpression()

		left = &BinaryExpression{
			Type:     NodeTypeBinaryExpression,
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left
}

// parseComparisonOnlyExpression handles comparison operators without logical operators.
// This is used to prevent infinite recursion in parseComparisonExpression.
func (p *Parser) parseComparisonOnlyExpression() Expression {
	left := p.parseAdditiveExpression()

	// Handle multiple comparison operators (left-associative)
	for p.at().Type == lexer.TokenTypeEqualEqual || p.at().Type == lexer.TokenTypeNotEqual ||
		p.at().Type == lexer.TokenTypeGreaterThan || p.at().Type == lexer.TokenTypeGreaterThanEqual ||
		p.at().Type == lexer.TokenTypeLessThan || p.at().Type == lexer.TokenTypeLessThanEqual {
		operator := p.next().Literal
		right := p.parseAdditiveExpression()

		left = &BinaryExpression{
			Type:     NodeTypeBinaryExpression,
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left
}

// parseAdditiveExpression handles addition and subtraction operators.
// Examples: a + b, x - y, "hello" + "world"
func (p *Parser) parseAdditiveExpression() Expression {
	left := p.parseMultiplicativeExpression()

	// Handle multiple additive operators (left-associative)
	for p.at().Literal == "+" || p.at().Literal == "-" {
		operator := p.next().Literal
		right := p.parseMultiplicativeExpression()

		left = &BinaryExpression{
			Type:     NodeTypeBinaryExpression,
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left
}

// parseMultiplicativeExpression handles multiplication, division, and modulo operators.
// Examples: a * b, x / y, n % 2
func (p *Parser) parseMultiplicativeExpression() Expression {
	left := p.parsePrimaryExpression()

	// Handle multiple multiplicative operators (left-associative)
	for p.at().Literal == "/" || p.at().Literal == "*" || p.at().Literal == "%" {
		operator := p.next().Literal
		right := p.parsePrimaryExpression()

		left = &BinaryExpression{
			Type:     NodeTypeBinaryExpression,
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left
}

// parsePrimaryExpression handles the highest precedence expressions.
// These include literals, identifiers, parentheses, and object literals.
func (p *Parser) parsePrimaryExpression() Expression {
	var expr Expression

	tokenType := p.at().Type
	switch tokenType {
	case lexer.TokenTypeIdentifier:
		token := p.next()
		expr = &Identifier{
			Type:  NodeTypeIdentifier,
			Name:  token.Literal,
			Token: &token,
		}
	case lexer.TokenTypeNumber:
		value, err := strconv.ParseFloat(p.next().Literal, 64)
		if err != nil {
			panic(err)
		}
		expr = &Numeric{
			Type:  NodeTypeNumeric,
			Value: value,
		}
	case lexer.TokenTypeOpenParentheses:
		p.next()
		expr = p.parseExpression()
		p.nextWithExpect(lexer.TokenTypeCloseParentheses, errors.ErrExpectedCloseParen)
	case lexer.TokenTypeNull:
		p.next()
		expr = &Null{}
	case lexer.TokenTypeTrue, lexer.TokenTypeYes, lexer.TokenTypeOn:
		p.next()
		expr = &Boolean{Value: true}
	case lexer.TokenTypeFalse, lexer.TokenTypeNo, lexer.TokenTypeOff:
		p.next()
		expr = &Boolean{Value: false}
	case lexer.TokenTypeString:
		token := p.next()
		expr = &String{
			Type:  NodeTypeString,
			Value: token.Literal,
		}
	case lexer.TokenTypeBreak:
		p.next()
		expr = &BreakExpression{}
	case lexer.TokenTypeOpenCurlyBrackets:
		expr = p.parseObjectExpression()
	case lexer.TokenTypeOpenSquareBrackets:
		expr = p.parseArrayExpression()
	default:
		p.syntaxError(p.at(), fmt.Sprintf(errors.ErrUnexpectedToken, p.at().Literal))
		return nil
	}

	// Handle postfix operations (member access, array indexing, function calls)
	return p.parsePostfixExpression(expr)
}

// parsePostfixExpression handles member access, array indexing, and function calls
// that can be chained after any expression (e.g., "hello".len(), [1,2,3].pop(), etc.)
func (p *Parser) parsePostfixExpression(expr Expression) Expression {
	for {
		switch p.at().Type {
		case lexer.TokenTypeOpenSquareBrackets:
			expr = p.parseArrayIndex(expr)
		case lexer.TokenTypeOpenParentheses:
			expr = p.parseCallExpression(expr)
		case lexer.TokenTypeDot:
			expr = p.parseMemberAccess(expr)
		default:
			return expr
		}
	}
}

func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
	p.next()
	identifier := p.nextWithExpect(lexer.TokenTypeIdentifier, errors.ErrExpectedFunctionName)
	args := p.parseArguments()
	var params []string
	for _, arg := range args {
		if _, ok := arg.(*Identifier); !ok {
			p.syntaxError(p.at(), errors.ErrExpectedIdentifierParam)
			return nil
		}
		params = append(params, arg.(*Identifier).Name)
	}

	p.nextWithExpect(lexer.TokenTypeOpenCurlyBrackets, errors.ErrExpectedOpenCurly)
	body := p.parseBlock()
	return &FunctionDeclaration{
		Identifier: identifier.Literal,
		Parameters: params,
		Body:       body,
	}

}

func (p *Parser) parseBlock() []Statement {
	statements := []Statement{}
	for p.notEOF() && p.at().Type != lexer.TokenTypeCloseCurlyBrackets {
		// Skip newlines
		if p.at().Type == lexer.TokenTypeNewline {
			p.next()
			continue
		}
		statement := p.parseStatement()
		statements = append(statements, statement)
	}
	p.nextWithExpect(lexer.TokenTypeCloseCurlyBrackets, errors.ErrExpectedCloseCurly)
	return statements
}

func (p *Parser) parseArguments() []Expression {
	p.nextWithExpect(lexer.TokenTypeOpenParentheses, errors.ErrExpectedOpenParen)
	arguments := []Expression{}
	for p.notEOF() && p.at().Type != lexer.TokenTypeCloseParentheses {
		// Skip newlines
		if p.at().Type == lexer.TokenTypeNewline {
			p.next()
			continue
		}
		if p.at().Type == lexer.TokenTypeComma {
			p.next()
			continue
		}
		argument := p.parseExpression()
		arguments = append(arguments, argument)
	}
	p.nextWithExpect(lexer.TokenTypeCloseParentheses, errors.ErrExpectedCloseParen)
	return arguments
}

// parseObjectExpression parses object literals.
// Examples: { name: "John", age: 30 }, { }, { nested: { value: 42 } }
func (p *Parser) parseObjectExpression() Expression {
	// If it's not an object literal, delegate to additive expressions
	if p.at().Type != lexer.TokenTypeOpenCurlyBrackets {
		return p.parseAdditiveExpression()
	}

	p.next() // consume the opening brace
	properties := []Property{}

	// Parse properties until closing brace
	for p.notEOF() && p.at().Type != lexer.TokenTypeCloseCurlyBrackets {
		// Skip newlines
		if p.at().Type == lexer.TokenTypeNewline {
			p.next()
			continue
		}
		key := p.nextWithExpect(lexer.TokenTypeIdentifier, errors.ErrExpectedIdentifier).Literal
		p.nextWithExpect(lexer.TokenTypeColon, errors.ErrExpectedColon)
		value := p.parseExpression()
		properties = append(properties, Property{Key: key, Value: value})

		// Skip comma if present
		if p.at().Type == lexer.TokenTypeComma {
			p.next()
		}
	}

	p.nextWithExpect(lexer.TokenTypeCloseCurlyBrackets, errors.ErrExpectedCloseCurly)
	return &Object{Properties: properties}
}

// parseArrayExpression parses array literals.
// Examples: [1, 2, 3], ["hello", "world"], []
func (p *Parser) parseArrayExpression() Expression {
	// If it's not an array literal, delegate to additive expressions
	if p.at().Type != lexer.TokenTypeOpenSquareBrackets {
		return p.parseAdditiveExpression()
	}

	p.next() // consume the opening bracket
	elements := []Expression{}

	// Parse elements until closing bracket
	for p.notEOF() && p.at().Type != lexer.TokenTypeCloseSquareBrackets {
		// Skip newlines
		if p.at().Type == lexer.TokenTypeNewline {
			p.next()
			continue
		}
		if p.at().Type == lexer.TokenTypeComma {
			p.next()
			continue
		}
		element := p.parseExpression()
		elements = append(elements, element)

		// Skip comma if present
		if p.at().Type == lexer.TokenTypeComma {
			p.next()
		}
	}

	p.nextWithExpect(lexer.TokenTypeCloseSquareBrackets, errors.ErrExpectedCloseSquare)
	return &Array{Elements: elements}
}

// parseArrayIndex handles array element access.
// Examples: arr[1], arr[i + 1]
func (p *Parser) parseArrayIndex(array Expression) Expression {
	p.next() // consume the opening bracket
	index := p.parseExpression()
	p.nextWithExpect(lexer.TokenTypeCloseSquareBrackets, errors.ErrExpectedCloseSquare)

	return &ArrayIndex{
		ArrayExpression: array,
		Index:           index,
	}
}

// parseMemberAccess handles property access.
// Examples: obj.name, person.address, str.len
func (p *Parser) parseMemberAccess(object Expression) Expression {
	p.next() // consume the dot
	property := p.nextWithExpect(lexer.TokenTypeIdentifier, errors.ErrExpectedIdentifier).Literal

	return &MemberAccess{
		Object:   object,
		Property: property,
	}
}

func (p *Parser) parseCallExpression(callee Expression) *CallExpression {
	p.nextWithExpect(lexer.TokenTypeOpenParentheses, errors.ErrExpectedOpenParen)

	args := []Expression{}

	// Parse arguments
	for p.notEOF() && p.at().Type != lexer.TokenTypeCloseParentheses {
		// Skip newlines
		if p.at().Type == lexer.TokenTypeNewline {
			p.next()
			continue
		}

		arg := p.parseExpression()
		args = append(args, arg)

		// Check for comma separator
		if p.at().Type == lexer.TokenTypeComma {
			p.next()
		}
	}

	p.nextWithExpect(lexer.TokenTypeCloseParentheses, errors.ErrExpectedCloseParen)

	return &CallExpression{
		Type:   NodeTypeCallExpression,
		Callee: callee,
		Args:   args,
	}
}

func (p *Parser) parseIfStatement() *IfStatement {
	p.next() // consume 'if'

	condition := p.parseExpression()

	p.nextWithExpect(lexer.TokenTypeOpenCurlyBrackets, errors.ErrExpectedOpenCurly)
	body := p.parseBlock()

	ifStatement := &IfStatement{
		Condition: condition,
		Body:      body,
		ElseIfs:   []ElseIfClause{},
		ElseBody:  []Statement{},
	}

	// Parse elseif clauses
	for p.notEOF() && p.at().Type == lexer.TokenTypeElse {
		p.next() // consume 'else'

		// Check if it's an elseif (has a condition)
		if p.at().Type == lexer.TokenTypeIf {
			p.next() // consume 'if'
			elseifCondition := p.parseExpression()
			p.nextWithExpect(lexer.TokenTypeOpenCurlyBrackets, errors.ErrExpectedOpenCurly)
			elseifBody := p.parseBlock()

			elseifClause := ElseIfClause{
				Condition: elseifCondition,
				Body:      elseifBody,
			}
			ifStatement.ElseIfs = append(ifStatement.ElseIfs, elseifClause)
		} else {
			// It's an else clause
			p.nextWithExpect(lexer.TokenTypeOpenCurlyBrackets, errors.ErrExpectedOpenCurly)
			ifStatement.ElseBody = p.parseBlock()
			break
		}
	}

	return ifStatement
}

func (p *Parser) parseLoopStatement() *LoopStatement {
	p.next() // consume 'loop'

	// Check if this is an infinite loop (no condition, directly follows with {)
	if p.at().Type == lexer.TokenTypeOpenCurlyBrackets {
		// Infinite loop - no condition
		p.next() // consume the opening brace
		body := p.parseBlock()
		return &LoopStatement{
			Condition: nil,
			Body:      body,
		}
	}

	// Check if this is a range loop or for-each loop (loop <var> from ...)
	if p.at().Type == lexer.TokenTypeIdentifier && len(p.tokens) > 4 && p.tokens[1].Type == lexer.TokenTypeFrom {
		loopVar := p.next().Literal // consume identifier (e.g., "i" or "element")
		p.nextWithExpect(lexer.TokenTypeFrom, "Expected 'from' after loop variable")
		from := p.parseExpression()

		// Check if this is a range loop (has 'to') or for-each loop (goes directly to {)
		if p.at().Type == lexer.TokenTypeTo {
			// Range loop: loop i from X to Y [; increment]
			p.next() // consume 'to'
			to := p.parseExpression()

			// Check if there's an optional increment
			var increment Expression
			if p.at().Type == lexer.TokenTypeColon {
				p.next() // consume colon
				increment = p.parseExpression()
			}

			p.nextWithExpect(lexer.TokenTypeOpenCurlyBrackets, errors.ErrExpectedOpenCurly)
			body := p.parseBlock()

			return &LoopStatement{
				LoopVar:   loopVar,
				From:      from,
				To:        to,
				Increment: increment,
				Body:      body,
			}
		} else {
			// For-each loop: loop element from arr { }
			p.nextWithExpect(lexer.TokenTypeOpenCurlyBrackets, errors.ErrExpectedOpenCurly)
			body := p.parseBlock()

			return &LoopStatement{
				LoopVar:   loopVar,
				From:      from, // This is the iterable (array)
				IsForEach: true,
				Body:      body,
			}
		}
	}

	// Traditional condition-based loop
	condition := p.parseExpression()
	p.nextWithExpect(lexer.TokenTypeOpenCurlyBrackets, errors.ErrExpectedOpenCurly)
	body := p.parseBlock()

	return &LoopStatement{
		Condition: condition,
		Body:      body,
	}
}

// parseReturnStatement parses return statements.
// Examples: return, return 42, return x + y
func (p *Parser) parseReturnStatement() *ReturnStatement {
	p.next() // consume 'return'

	// Check if return has a value or is bare
	// If the next token is a closing curly brace or newline, it's a bare return
	if p.at().Type == lexer.TokenTypeCloseCurlyBrackets || p.at().Type == lexer.TokenTypeNewline || p.at().Type == lexer.TokenTypeEOF {
		return &ReturnStatement{
			Value: nil,
		}
	}

	// Parse the return value expression
	value := p.parseExpression()

	return &ReturnStatement{
		Value: value,
	}
}
