package main

import "fmt"

// TokenType defines the type of token
type TokenType string

const (
	// Single-character tokens
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	DOT         TokenType = "DOT"
	SEMICOLON   TokenType = "SEMICOLON"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	COMMA       TokenType = "COMMA"
	STAR        TokenType = "STAR"

	// One or two character tokens
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	SLASH         TokenType = "SLASH"

	// Literals
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"
	IDENTIFIER TokenType = "IDENTIFIER"

	// Keywords
	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FOR    TokenType = "FOR"
	FUN    TokenType = "FUN"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"

	// End of file
	EOF TokenType = "EOF"
)

// Object represents any Literal value that a token can have
type Object interface{}

// Token represents a token in the source code
type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   Object
	Line      int
}

// NewToken creates a new token
func NewToken(Line int, TokenType TokenType, Literal Object, Lexeme string) *Token {
	t := new(Token)
	t.Line = Line
	t.TokenType = TokenType
	t.Literal = Literal
	t.Lexeme = Lexeme
	return t
}

// String returns a string representation of the token
func (t *Token) String() string {
	return fmt.Sprintf("%v %s %v", t.TokenType, t.Lexeme, t.Literal)
}
