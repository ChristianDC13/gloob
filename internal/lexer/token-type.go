package lexer

type TokenType string

const (
	// Operators and delimiters
	TokenTypeEqual               TokenType = "EQUAL"
	TokenTypeEqualEqual          TokenType = "EQUAL_EQUAL"
	TokenTypeNotEqual            TokenType = "NOT_EQUAL"
	TokenTypeGreaterThan         TokenType = "GREATER_THAN"
	TokenTypeGreaterThanEqual    TokenType = "GREATER_THAN_EQUAL"
	TokenTypeLessThan            TokenType = "LESS_THAN"
	TokenTypeLessThanEqual       TokenType = "LESS_THAN_EQUAL"
	TokenTypeOperator            TokenType = "OPERATOR"
	TokenTypeOpenParentheses     TokenType = "OPEN_PARENTHESES"
	TokenTypeCloseParentheses    TokenType = "CLOSE_PARENTHESES"
	TokenTypeOpenCurlyBrackets   TokenType = "OPEN_CURLY_BRACKETS"
	TokenTypeCloseCurlyBrackets  TokenType = "CLOSE_CURLY_BRACKETS"
	TokenTypeOpenSquareBrackets  TokenType = "OPEN_SQUARE_BRACKETS"
	TokenTypeCloseSquareBrackets TokenType = "CLOSE_SQUARE_BRACKETS"
	TokenTypeColon               TokenType = "COLON"
	TokenTypeSemicolon           TokenType = "SEMICOLON"
	TokenTypeAmpersand           TokenType = "AMPERSAND"
	TokenTypeAnd                 TokenType = "AND"
	TokenTypeOr                  TokenType = "OR"
	TokenTypeDot                 TokenType = "DOT"
	TokenTypeComma               TokenType = "COMMA"
	TokenTypePipe                TokenType = "PIPE"
	TokenTypeExclamation         TokenType = "EXCLAMATION"
	TokenTypeNewline             TokenType = "NEWLINE"
	TokenTypeComment             TokenType = "COMMENT"

	// Tokens
	TokenTypeNumber     TokenType = "NUMBER"
	TokenTypeIdentifier TokenType = "IDENTIFIER"
	TokenTypeUnknown    TokenType = "UNKNOWN"
	TokenTypeString     TokenType = "STRING"
	TokenTypeBoolean    TokenType = "BOOLEAN"
	TokenTypeNull       TokenType = "NULL"

	// Keywords
	TokenTypeFunction TokenType = "FUNCTION"
	TokenTypeLoop     TokenType = "LOOP"
	TokenTypeIf       TokenType = "IF"
	TokenTypeElse     TokenType = "ELSE"
	TokenTypeReturn   TokenType = "RETURN"
	TokenTypeBreak    TokenType = "BREAK"
	TokenTypeContinue TokenType = "CONTINUE"
	TokenTypeImport   TokenType = "IMPORT"
	TokenTypeVar      TokenType = "VAR"
	TokenTypeConst    TokenType = "CONST"
	TokenTypeFrom     TokenType = "FROM"
	TokenTypeTo       TokenType = "TO"
	TokenTypeTrue     TokenType = "TRUE"
	TokenTypeFalse    TokenType = "FALSE"
	TokenTypeYes      TokenType = "YES"
	TokenTypeNo       TokenType = "NO"
	TokenTypeOn       TokenType = "ON"
	TokenTypeOff      TokenType = "OFF"

	// Special tokens
	TokenTypeEOF TokenType = "EOF"
)
