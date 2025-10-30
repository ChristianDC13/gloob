package lexer

var Keywords = map[string]TokenType{
	"var":      TokenTypeVar,
	"const":    TokenTypeConst,
	"function": TokenTypeFunction,
	"loop":     TokenTypeLoop,
	"if":       TokenTypeIf,
	"else":     TokenTypeElse,
	"return":   TokenTypeReturn,
	"break":    TokenTypeBreak,
	"continue": TokenTypeContinue,
	"import":   TokenTypeImport,
	"true":     TokenTypeTrue,
	"false":    TokenTypeFalse,
	"yes":      TokenTypeYes,
	"no":       TokenTypeNo,
	"on":       TokenTypeOn,
	"off":      TokenTypeOff,
	"from":     TokenTypeFrom,
	"to":       TokenTypeTo,
	"null":     TokenTypeNull,
	"fun":      TokenTypeFunction,
}
