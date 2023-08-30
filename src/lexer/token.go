package lexer

import (
	"fmt"
	"text/scanner"
)

type TokenType int

const (
	TK_LEFT_PAREN TokenType = iota
	TK_RIGHT_PAREN
	TK_LEFT_BRACE
	TK_RIGHT_BRACE
	TK_LEFT_BRACKET
	TK_RIGHT_BRACKET
	TK_COMMA
	TK_DOT
	TK_COLON
	TK_SEMICOLON
	TK_QUESTION_MARK
	TK_BIT_NOT
	TK_BANG
	TK_BANG_EQUAL
	TK_MINUS
	TK_MINUS_MINUS
	TK_PLUS
	TK_PLUS_PLUS
	TK_SLASH
	TK_STAR
	TK_BIT_AND
	TK_BIT_OR
	TK_EQUAL_EQUAL
	TK_GREATER
	TK_BIT_SHIFT_RIGHT
	TK_GREATER_EQUAL
	TK_LESS
	TK_BIT_SHIFT_LEFT
	TK_LESS_EQUAL
	TK_MINUS_EQUAL
	TK_PLUS_EQUAL
	TK_SLASH_EQUAL
	TK_STAR_EQUAL
	TK_AND
	TK_BIT_AND_EQUAL
	TK_OR
	TK_BIT_OR_EQUAL
	TK_EQUAL
	TK_ARROW
	TK_IDENTIFIER
	TK_INTEGER
	TK_FLOAT
	TK_CHAR
	TK_STRING
	TK_RAWSTRING
	TK_COMMENT
	TK_CLASS
	TK_STRUCT
	TK_IF
	TK_ELSE
	TK_TRUE
	TK_FALSE
	TK_FOR
	TK_WHILE
	TK_FN
	TK_NULL
	TK_RETURN
	TK_EXIT
	TK_SUPER
	TK_THIS
	TK_VAR
	TK_VAL
	TK_ERROR
	TK_EOF
)

type Token struct {
	Type     TokenType
	Position scanner.Position
	Literal  string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %d, Position: %s, Literal: %q}", t.Type, t.Position, t.Literal)
}
