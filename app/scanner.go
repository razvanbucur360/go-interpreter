package main

import (
	"fmt"
	"strconv"
)

// Scanner performs lexical analysis to convert source code into tokens
type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
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

// ScanTokens scans all tokens in the source
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	// Add EOF token
	s.tokens = append(s.tokens, *NewToken(s.line, EOF, nil, ""))
	return s.tokens
}

// isAtEnd checks if we've reached the end of the source
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// scanToken scans a single token
func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case '.':
		s.addToken(DOT)
	case ',':
		s.addToken(COMMA)
	case ';':
		s.addToken(SEMICOLON)
	case '+':
		s.addToken(PLUS)
	case '-':
		s.addToken(MINUS)
	case '*':
		s.addToken(STAR)
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
		break
	case '"':
		s.string()
		break
	default:
		// Use the error function instead of just printing
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			error(s.line, fmt.Sprintf("Unexpected character: %c", c))
		}
	}
}

// advance consumes the next character in the source
func (s *Scanner) advance() byte {
	if s.isAtEnd() {
		return 0
	}
	result := s.source[s.current]
	s.current++
	return result
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	} else {
		return s.source[s.current]
	}
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		error(s.line, "Unterminated string.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]

	s.addTokenLiteral(STRING, value)

}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	numStr := s.source[s.start:s.current]

	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		error(s.line, "Invalid number.")
		return
	}

	s.addTokenLiteral(NUMBER, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumneric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	TokenType, exists := keywords[text]
	if !exists {
		TokenType = IDENTIFIER
	}

	s.addToken(TokenType)
}

func (s *Scanner) isDigit(expected byte) bool {
	return expected >= '0' && expected <= '9'
}

func (s *Scanner) isAlpha(expected byte) bool {
	return (expected >= 'a' && expected <= 'z') || (expected >= 'A' && expected <= 'Z' || expected == '_')
}

func (s *Scanner) isAlphaNumneric(expected byte) bool {
	return s.isDigit(expected) || s.isAlpha(expected)
}

// addToken adds a token with no Literal value
func (s *Scanner) addToken(TokenType TokenType) {
	s.addTokenLiteral(TokenType, nil)
}

// addTokenLiteral adds a token with a Literal value
func (s *Scanner) addTokenLiteral(TokenType TokenType, Literal Object) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *NewToken(s.line, TokenType, Literal, text))
}
