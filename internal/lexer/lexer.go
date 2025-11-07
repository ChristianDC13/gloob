package lexer

import "unicode"

type Token struct {
	Type        TokenType
	Literal     string
	Line        int
	ColumnStart int
	ColumnEnd   int
	Filename    string
}

func CaptureToken(literal string, tokenType TokenType, line int, columnStart int, columnEnd int, filename string) Token {
	return Token{
		Type:        tokenType,
		Literal:     literal,
		Line:        line,
		ColumnStart: columnStart,
		ColumnEnd:   columnEnd,
		Filename:    filename,
	}
}

type Lexer struct {
	input    string
	filename string
}

func NewLexer(input string, filename string) *Lexer {
	return &Lexer{
		input:    input,
		filename: filename,
	}
}

func (l *Lexer) Tokenize() []Token {
	tokens := []Token{}
	chars := []rune(l.input)

	line := 1
	column := 1

	for len(chars) > 0 {
		ch := chars[0]
		columnStart := column

		// handle whitespace
		if ch == ' ' || ch == '\t' || ch == '\r' {
			chars = chars[1:]
			column++
			continue
		}

		if ch == '\n' {
			tokens = append(tokens, CaptureToken("\n", TokenTypeNewline, line, columnStart, column, l.filename))
			line++
			column = 1
			chars = chars[1:]
			continue
		}

		tokenType := TokenTypeUnknown
		literal := string(ch)

		if unicode.IsLetter(ch) {
			literal = ""
			for len(chars) > 0 && (unicode.IsLetter(chars[0]) || unicode.IsDigit(chars[0])) {
				literal += string(chars[0])
				chars = chars[1:]
				column++
			}
			if isKeyW, tokenType := isKeyword(literal); isKeyW {
				tokens = append(tokens, CaptureToken(literal, tokenType, line, columnStart, column-1, l.filename))
				continue
			}
			tokenType = TokenTypeIdentifier
			tokens = append(tokens, CaptureToken(literal, tokenType, line, columnStart, column-1, l.filename))
			continue
		}

		// Handle negative numbers: check if '-' is followed by a digit
		if ch == '-' && len(chars) > 1 && unicode.IsDigit(chars[1]) {
			literal = string(ch)
			chars = chars[1:] // consume the '-'
			column++
			// Continue to parse as number
			for len(chars) > 0 && (unicode.IsDigit(chars[0]) || chars[0] == '.') {
				literal += string(chars[0])
				chars = chars[1:]
				column++
			}
			tokenType = TokenTypeNumber
			tokens = append(tokens, CaptureToken(literal, tokenType, line, columnStart, column-1, l.filename))
			continue
		}

		if unicode.IsDigit(ch) {
			literal = ""
			for len(chars) > 0 && (unicode.IsDigit(chars[0]) || chars[0] == '.') {
				literal += string(chars[0])
				chars = chars[1:]
				column++
			}
			tokenType = TokenTypeNumber
			tokens = append(tokens, CaptureToken(literal, tokenType, line, columnStart, column-1, l.filename))
			continue
		}

		// Handle string literals (both single and double quotes)
		if ch == '"' || ch == '\'' {
			quoteChar := ch
			literal = ""
			chars = chars[1:] // consume opening quote
			column++

			for len(chars) > 0 && chars[0] != quoteChar {
				literal += string(chars[0])
				chars = chars[1:]
				column++
			}

			if len(chars) == 0 {
				// Unterminated string
				tokens = append(tokens, CaptureToken(literal, TokenTypeUnknown, line, columnStart, column-1, l.filename))
				continue
			}

			chars = chars[1:] // consume closing quote
			column++
			tokenType = TokenTypeString
			tokens = append(tokens, CaptureToken(literal, tokenType, line, columnStart, column-1, l.filename))
			continue
		}

		switch ch {
		case '=':
			// Check for == operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = "=="
				tokenType = TokenTypeEqualEqual
				chars = chars[1:] // consume second =
				column++
			} else {
				tokenType = TokenTypeEqual
			}
		case '!':
			// Check for != operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = "!="
				tokenType = TokenTypeNotEqual
				chars = chars[1:] // consume =
				column++
			} else {
				tokenType = TokenTypeExclamation
			}
		case '>':
			// Check for >= operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = ">="
				tokenType = TokenTypeGreaterThanEqual
				chars = chars[1:] // consume =
				column++
			} else {
				tokenType = TokenTypeGreaterThan
			}
		case '<':
			// Check for <= operator
			if len(chars) > 1 && chars[1] == '=' {
				literal = "<="
				tokenType = TokenTypeLessThanEqual
				chars = chars[1:] // consume =
				column++
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
				column++
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
				column++
			} else {
				tokenType = TokenTypeAmpersand
			}
		case '|':
			if len(chars) > 1 && chars[1] == '|' {
				literal = "||"
				tokenType = TokenTypeOr
				chars = chars[1:] // consume second |
				column++
			} else {
				tokenType = TokenTypePipe
			}
		default:
			tokenType = TokenTypeUnknown
		}

		column++
		tokens = append(tokens, CaptureToken(literal, tokenType, line, columnStart, column-1, l.filename))
		chars = chars[1:]
	}

	tokens = append(tokens, CaptureToken("EOF", TokenTypeEOF, line, column, column, l.filename))

	return tokens
}

func isKeyword(literal string) (bool, TokenType) {
	if keyword, ok := Keywords[literal]; ok {
		return true, keyword
	}
	return false, TokenTypeUnknown
}
