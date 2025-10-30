package lexer

import "unicode"

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func CaptureToken(literal string, tokenType TokenType, line int, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}

type Lexer struct {
	input string
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) Tokenize() []Token {
	tokens := []Token{}
	chars := []rune(l.input)

	line := 1
	column := 1

	for len(chars) > 0 {
		ch := chars[0]

		// handle whitespace
		if ch == ' ' || ch == '\t' || ch == '\r' {
			chars = chars[1:]
			continue
		}

		if ch == '\n' {
			tokens = append(tokens, CaptureToken("\n", TokenTypeNewline, line, column))
			line++
			column = 1
			chars = chars[1:]
			continue
		}

		column++

		tokenType := TokenTypeUnknown
		literal := string(ch)

		if unicode.IsLetter(ch) {
			literal = ""
			for len(chars) > 0 && (unicode.IsLetter(chars[0]) || unicode.IsDigit(chars[0])) {
				literal += string(chars[0])
				chars = chars[1:]
			}
			if isKeyW, tokenType := isKeyword(literal); isKeyW {
				tokens = append(tokens, CaptureToken(literal, tokenType, line, column))
				continue
			}
			tokenType = TokenTypeIdentifier
			tokens = append(tokens, CaptureToken(literal, tokenType, line, column))
			continue
		}

		// Handle negative numbers: check if '-' is followed by a digit
		if ch == '-' && len(chars) > 1 && unicode.IsDigit(chars[1]) {
			literal = string(ch)
			chars = chars[1:] // consume the '-'
			// Continue to parse as number
			for len(chars) > 0 && (unicode.IsDigit(chars[0]) || chars[0] == '.') {
				literal += string(chars[0])
				chars = chars[1:]
			}
			tokenType = TokenTypeNumber
			tokens = append(tokens, CaptureToken(literal, tokenType, line, column))
			continue
		}

		if unicode.IsDigit(ch) {
			literal = ""
			for len(chars) > 0 && (unicode.IsDigit(chars[0]) || chars[0] == '.') {
				literal += string(chars[0])
				chars = chars[1:]
			}
			tokenType = TokenTypeNumber
			tokens = append(tokens, CaptureToken(literal, tokenType, line, column))
			continue
		}

		// Handle string literals (both single and double quotes)
		if ch == '"' || ch == '\'' {
			quoteChar := ch
			literal = ""
			chars = chars[1:] // consume opening quote

			for len(chars) > 0 && chars[0] != quoteChar {
				literal += string(chars[0])
				chars = chars[1:]
			}

			if len(chars) == 0 {
				// Unterminated string
				tokens = append(tokens, CaptureToken(literal, TokenTypeUnknown, line, column))
				continue
			}

			chars = chars[1:] // consume closing quote
			tokenType = TokenTypeString
			tokens = append(tokens, CaptureToken(literal, tokenType, line, column))
			continue
		}

		switch ch {
		case '=':
			// Check for == operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = "=="
				tokenType = TokenTypeEqualEqual
				chars = chars[1:] // consume second =
			} else {
				tokenType = TokenTypeEqual
			}
		case '!':
			// Check for != operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = "!="
				tokenType = TokenTypeNotEqual
				chars = chars[1:] // consume =
			} else {
				tokenType = TokenTypeExclamation
			}
		case '>':
			// Check for >= operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = ">="
				tokenType = TokenTypeGreaterThanEqual
				chars = chars[1:] // consume =
			} else {
				tokenType = TokenTypeGreaterThan
			}
		case '<':
			// Check for <= operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = "<="
				tokenType = TokenTypeLessThanEqual
				chars = chars[1:] // consume =
			} else {
				tokenType = TokenTypeLessThan
			}
		case '+', '-', '*', '%':
			tokenType = TokenTypeOperator
		case '/':
			if len(chars) > 1 && chars[1] == '/' {
				literal = "//"
				tokenType = TokenTypeComment
				chars = chars[1:] // consume /
			} else {
				tokenType = TokenTypeOperator
			}
		case '(':
			tokenType = TokenTypeOpenParentheses
		case ')':
			tokenType = TokenTypeCloseParentheses
		case '{':
			tokenType = TokenTypeOpenCurlyBrackets
		case '}':
			tokenType = TokenTypeCloseCurlyBrackets
		case '[':
			tokenType = TokenTypeOpenSquareBrackets
		case ']':
			tokenType = TokenTypeCloseSquareBrackets
		case ':':
			tokenType = TokenTypeColon
		case ';':
			tokenType = TokenTypeSemicolon
		case ',':
			tokenType = TokenTypeComma
		case '.':
			tokenType = TokenTypeDot
		case '&':
			if len(chars) > 1 && chars[1] == '&' {
				literal = "&&"
				tokenType = TokenTypeAnd
				chars = chars[1:] // consume second &
			} else {
				tokenType = TokenTypeAmpersand
			}
		case '|':
			if len(chars) > 1 && chars[1] == '|' {
				literal = "||"
				tokenType = TokenTypeOr
				chars = chars[1:] // consume second |
			} else {
				tokenType = TokenTypePipe
			}
		default:
			tokenType = TokenTypeUnknown
		}

		tokens = append(tokens, CaptureToken(literal, tokenType, line, column))
		chars = chars[1:]
	}

	tokens = append(tokens, CaptureToken("EOF", TokenTypeEOF, line, column))

	return tokens
}

func isKeyword(literal string) (bool, TokenType) {
	if keyword, ok := Keywords[literal]; ok {
		return true, keyword
	}
	return false, TokenTypeUnknown
}
