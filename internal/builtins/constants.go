package builtins

import (
	"gloob-interpreter/internal/parser"
	"gloob-interpreter/internal/scope"
	"gloob-interpreter/internal/values"
)

// SetupConstants adds all built-in constants to the scope
func SetupConstants(s *scope.Scope) {
	s.Declare("null", &values.NullValue{Type: parser.NodeTypeNull}, true)
	s.Declare("pi", &values.NumericValue{Type: parser.NodeTypeNumeric, Value: 3.141592653589793}, true)
}
