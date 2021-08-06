package main

import (
	"strconv"
	"unicode"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	content string
	Tokens  []Token

	start   int
	current int
	line    int
}

func NewScanner(content string) *Scanner {
	return &Scanner{
		content: content,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (sc *Scanner) ScanTokens() []Token {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.scanToken()
	}

	sc.addToken(EOF)
	return sc.Tokens
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.content)
}

func (sc *Scanner) scanToken() {
	switch c := sc.advance(); c {
	// Single-character tokens
	case '(':
		sc.addToken(LEFT_PAREN)
	case ')':
		sc.addToken(RIGHT_PAREN)
	case '{':
		sc.addToken(LEFT_BRACE)
	case '}':
		sc.addToken(RIGHT_BRACE)
	case ',':
		sc.addToken(COMMA)
	case '.':
		sc.addToken(DOT)
	case '-':
		sc.addToken(MINUS)
	case '+':
		sc.addToken(PLUS)
	case ';':
		sc.addToken(SEMICOLON)
	case '/':
		if sc.match('/') {
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
		} else if sc.match('*') {
			sc.multilineComment()
		} else {
			sc.addToken(SLASH)
		}
	case '*':
		sc.addToken(STAR)

	// One or two character tokens
	case '!':
		if sc.match('=') {
			sc.addToken(BANG_EQUAL)
		} else {
			sc.addToken(EQUAL)
		}
	case '=':
		if sc.match('=') {
			sc.addToken(EQUAL_EQUAL)
		} else {
			sc.addToken(EQUAL)
		}
	case '>':
		if sc.match('=') {
			sc.addToken(GREATER_EQUAL)
		} else {
			sc.addToken(GREATER)
		}
	case '<':
		if sc.match('=') {
			sc.addToken(LESS_EQUAL)
		} else {
			sc.addToken(LESS)
		}

	// Ignore whitespace
	case ' ', '\r', '\t':
		return

	// New lines
	case '\n':
		sc.line += 1
		return

	// Literals
	case '"':
		sc.string()

	default:
		// Numbers
		if unicode.IsDigit(rune(c)) {
			sc.number()
			return
		}

		// Identifiers
		if unicode.IsLetter(rune(c)) || c == '_' {
			sc.identifier()
			return
		}

		LoxError(sc.line, "Unexpected character")
	}
}

func (sc *Scanner) advance() byte {
	c := sc.content[sc.current]
	sc.current += 1
	return c
}

func (sc *Scanner) addToken(tokenType TokenType) {
	sc.addTokenWithLiteral(tokenType, nil)
}

func (sc *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	var lexeme string
	if tokenType != EOF {
		lexeme = sc.content[sc.start:sc.current]
	}

	sc.Tokens = append(sc.Tokens, NewToken(tokenType, lexeme, literal, sc.line))
}

func (sc *Scanner) match(c byte) bool {
	if sc.peek() != c {
		return false
	}

	sc.advance()
	return true
}

func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return 0
	} else {
		return sc.content[sc.current]
	}
}

func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.content) {
		return 0
	} else {
		return sc.content[sc.current+1]
	}
}

func (sc *Scanner) string() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line += 1
		}

		// we consume the character after the \ so we support things like \"
		if sc.peek() == '\\' {
			sc.advance()
		}

		// as we can advance above we want to check we can still consume characters
		if !sc.isAtEnd() {
			sc.advance()
		}
	}

	if sc.isAtEnd() {
		LoxError(sc.line, "Unterminated string")
	}

	sc.advance() // closing quote (")

	value := sc.content[sc.start+1 : sc.current-1] // trim surrounding quotes
	sc.addTokenWithLiteral(STRING, value)
}

func (sc *Scanner) number() {
	for unicode.IsDigit(rune(sc.peek())) {
		sc.advance()
	}

	// fractional part (optional)
	if sc.peek() == '.' && unicode.IsDigit(rune(sc.peekNext())) {
		sc.advance() // consume the '.'

		for unicode.IsDigit(rune(sc.peek())) {
			sc.advance()
		}
	}

	value, err := strconv.ParseFloat(sc.content[sc.start:sc.current], 64)
	if err != nil {
		LoxError(sc.line, err.Error())
	}

	sc.addTokenWithLiteral(NUMBER, value)
}

func (sc *Scanner) identifier() {
	for unicode.IsLetter(rune(sc.peek())) || unicode.IsDigit(rune(sc.peek())) || sc.peek() == '_' {
		sc.advance()
	}

	value := sc.content[sc.start:sc.current]
	if tokenType, ok := keywords[value]; ok {
		sc.addToken(tokenType)
	} else {
		sc.addToken(IDENTIFIER)
	}
}

func (sc *Scanner) multilineComment() {
	depth := 1
	for !sc.isAtEnd() && depth > 0 {
		if sc.peek() == '*' && sc.peekNext() == '/' {
			depth -= 1
			sc.advance()
		} else if sc.peek() == '/' && sc.peekNext() == '*' {
			depth += 1
			sc.advance()
		} else if sc.peek() == '\n' {
			sc.line += 1
		}

		// as we can advance above we want to prevent overflow
		if !sc.isAtEnd() {
			sc.advance()
		}
	}

	if depth > 0 {
		LoxError(sc.line, "Multiline comment was not closed")
	}
}
