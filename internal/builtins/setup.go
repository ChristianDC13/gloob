package builtins

import "gloob-interpreter/internal/scope"

// SetupBuiltins sets up all built-in constants and native functions
func SetupBuiltins(s *scope.Scope) {
	SetupConstants(s)
	SetupNativeFunctions(s)
}
