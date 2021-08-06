package main

type TokenType string

const (
	// Single-character tokens
	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	LEFT_BRACE  = "LEFT_BRACE"
	RIGHT_BRACE = "RIGHT_BRACE"
	COMMA       = "COMMA"
	DOT         = "DOT"
	MINUS       = "MINUS"
	PLUS        = "PLUS"
	SEMICOLON   = "SEMICOLON"
	SLASH       = "SLASH"
	STAR        = "STAR"

	// One or two character tokens
	BANG_EQUAL    = "BANG_EQUAL"
	BANG          = "BANG"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	EQUAL         = "EQUAL"
	GREATER_EQUAL = "GREATER_EQUAL"
	GREATER       = "GREATER"
	LESS_EQUAL    = "LESS_EQUAL"
	LESS          = "LESS"

	// Literals
	STRING     = "STRING"
	NUMBER     = "NUMBER"
	IDENTIFIER = "IDENTIFIER"

	// Keywords
	AND    = "AND"
	CLASS  = "CLASS"
	ELSE   = "ELSE"
	FALSE  = "FALSE"
	FUN    = "FUN"
	FOR    = "FOR"
	IF     = "IF"
	NIL    = "NIL"
	OR     = "OR"
	PRINT  = "PRINT"
	RETURN = "RETURN"
	SUPER  = "SUPER"
	THIS   = "THIS"
	TRUE   = "TRUE"
	VAR    = "VAR"
	WHILE  = "WHILE"

	EOF = "EOF"
)
