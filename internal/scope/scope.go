package scope

import (
	"fmt"
	"gloob-interpreter/internal/errors"
	"gloob-interpreter/internal/lexer"
	"gloob-interpreter/internal/values"
)

type Scope struct {
	parent     *Scope
	variables  map[string]values.RuntimeValue
	constants  map[string]struct{}
	sourceCode string // Source code for error reporting
}

func NewScope(parent *Scope) *Scope {
	scope := &Scope{
		parent:    parent,
		variables: make(map[string]values.RuntimeValue),
		constants: make(map[string]struct{}),
	}
	// Inherit source code from parent if available
	if parent != nil {
		scope.sourceCode = parent.sourceCode
	}
	return scope
}

// SetSourceCode sets the source code for error reporting
func (s *Scope) SetSourceCode(sourceCode string) {
	s.sourceCode = sourceCode
}

func (s *Scope) Declare(name string, value values.RuntimeValue, isConstant bool) values.RuntimeValue {
	if _, ok := s.variables[name]; ok {
		errors.RuntimeError(nil, "", fmt.Sprintf(errors.ErrVariableAlreadyDeclared, name))
		return nil
	}
	if isConstant {
		s.constants[name] = struct{}{}
	}
	s.variables[name] = value
	return value
}

func (s *Scope) Assign(name string, value values.RuntimeValue) values.RuntimeValue {
	scope := s.Resolve(name)
	if scope == nil {
		errors.RuntimeError(nil, "", fmt.Sprintf(errors.ErrVariableNotFound, name))
		return nil
	}
	if _, ok := scope.constants[name]; ok {
		errors.RuntimeError(nil, "", fmt.Sprintf(errors.ErrConstantCannotBeAssigned, name))
		return nil
	}
	scope.variables[name] = value
	return value
}

func (s *Scope) Resolve(name string) *Scope {
	_, ok := s.variables[name]
	if ok {
		return s
	}
	if s.parent != nil {
		return s.parent.Resolve(name)
	}

	return nil // Don't error here, let the caller handle it
}

func (s *Scope) Get(name string) values.RuntimeValue {
	scope := s.Resolve(name)
	if scope == nil {
		errors.RuntimeError(nil, "", fmt.Sprintf(errors.ErrVariableNotFound, name))
		return nil
	}
	value := scope.variables[name]
	if value == nil {
		errors.RuntimeError(nil, "", fmt.Sprintf(errors.ErrVariableNotInitialized, name))
		return nil
	}
	return value
}

// GetWithToken gets a variable value and reports errors with token information
func (s *Scope) GetWithToken(name string, token *lexer.Token) values.RuntimeValue {
	scope := s.Resolve(name)
	if scope == nil {
		errors.RuntimeError(token, s.sourceCode, fmt.Sprintf(errors.ErrVariableNotFound, name))
		return nil
	}
	value := scope.variables[name]
	if value == nil {
		errors.RuntimeError(token, s.sourceCode, fmt.Sprintf(errors.ErrVariableNotInitialized, name))
		return nil
	}
	return value
}

func (s *Scope) GetVariables() map[string]values.RuntimeValue {
	return s.variables
}
