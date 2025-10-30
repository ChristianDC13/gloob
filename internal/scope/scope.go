package scope

import (
	"fmt"
	"gloob-interpreter/internal/colors"
	"gloob-interpreter/internal/values"
	"os"
)

type Scope struct {
	parent    *Scope
	variables map[string]values.RuntimeValue
	constants map[string]struct{}
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]values.RuntimeValue),
		constants: make(map[string]struct{}),
	}
}

func (s *Scope) Declare(name string, value values.RuntimeValue, isConstant bool) values.RuntimeValue {
	if _, ok := s.variables[name]; ok {
		fmt.Printf("Variable %s already declared\n", name)
		os.Exit(1)
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
	if _, ok := scope.constants[name]; ok {
		fmt.Printf("Constant %s cannot be assigned to because it is, how can i say it to you? It is a constant 😒\n", colors.Red(name))
		os.Exit(1)
		return nil
	}
	if scope == nil {
		fmt.Printf("Variable %s not found. Are you sure you typed it correctly? 🤔\n", colors.Red(name))
		os.Exit(1)
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

	fmt.Printf("Variable %s not found. Are you sure you typed it correctly? 🤔\n", colors.Red(name))
	os.Exit(1)
	return nil
}

func (s *Scope) Get(name string) values.RuntimeValue {
	scope := s.Resolve(name)
	if scope == nil {
		fmt.Printf("Variable %s not found. Are you sure you typed it correctly? 🤔\n", colors.Red(name))
		os.Exit(1)
		return nil
	}
	value := scope.variables[name]
	if value == nil {
		fmt.Printf("Variable %s is not initialized. Are you sure you declared it? 🤔\n", colors.Red(name))
		os.Exit(1)
		return nil
	}
	return value
}

func (s *Scope) GetVariables() map[string]values.RuntimeValue {
	return s.variables
}
