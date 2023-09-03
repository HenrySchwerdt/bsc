package lexer

import (
	"os"
	"text/scanner"
)

type Tokenizer struct {
	s *scanner.Scanner
}

func NewTokenizer(file *os.File) *Tokenizer {
	var s scanner.Scanner
	s.Position.Filename = file.Name()
	s.Init(file)
	return &Tokenizer{
		s: &s,
	}
}

func (tk *Tokenizer) Peek() rune {
	return tk.s.Peek()
}

func (tk *Tokenizer) GetToken() Token {
	for {
		tok := tk.s.Scan()
		switch tok {
		case scanner.Ident:
			literal := tk.s.TokenText()
			switch literal {
			case "class":
				return Token{Type: TK_CLASS, Position: tk.s.Position, Literal: literal}
			case "struct":
				return Token{Type: TK_STRUCT, Position: tk.s.Position, Literal: literal}
			case "if":
				return Token{Type: TK_IF, Position: tk.s.Position, Literal: literal}
			case "else":
				return Token{Type: TK_ELSE, Position: tk.s.Position, Literal: literal}
			case "true":
				return Token{Type: TK_TRUE, Position: tk.s.Position, Literal: literal}
			case "false":
				return Token{Type: TK_FALSE, Position: tk.s.Position, Literal: literal}
			case "for":
				return Token{Type: TK_FOR, Position: tk.s.Position, Literal: literal}
			case "while":
				return Token{Type: TK_WHILE, Position: tk.s.Position, Literal: literal}
			case "fn":
				return Token{Type: TK_FN, Position: tk.s.Position, Literal: literal}
			case "null":
				return Token{Type: TK_NULL, Position: tk.s.Position, Literal: literal}
			case "return":
				return Token{Type: TK_RETURN, Position: tk.s.Position, Literal: literal}
			case "exit":
				return Token{Type: TK_EXIT, Position: tk.s.Position, Literal: literal}
			case "break":
				return Token{Type: TK_BREAK, Position: tk.s.Position, Literal: literal}
			case "super":
				return Token{Type: TK_SUPER, Position: tk.s.Position, Literal: literal}
			case "this":
				return Token{Type: TK_THIS, Position: tk.s.Position, Literal: literal}
			case "var":
				return Token{Type: TK_VAR, Position: tk.s.Position, Literal: literal}
			case "val":
				return Token{Type: TK_VAL, Position: tk.s.Position, Literal: literal}
			default:
				return Token{Type: TK_IDENTIFIER, Position: tk.s.Position, Literal: literal}
			}
		case scanner.Int:
			return Token{Type: TK_INTEGER, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case scanner.Float:
			return Token{Type: TK_FLOAT, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case scanner.Char:
			return Token{Type: TK_CHAR, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case scanner.String:
			return Token{Type: TK_STRING, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case scanner.RawString:
			return Token{Type: TK_RAWSTRING, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case scanner.EOF:
			return Token{Type: TK_EOF, Position: tk.s.Position, Literal: ""}
		case '(':
			return Token{Type: TK_LEFT_PAREN, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case ')':
			return Token{Type: TK_RIGHT_PAREN, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '{':
			return Token{Type: TK_LEFT_BRACE, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '}':
			return Token{Type: TK_RIGHT_BRACE, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '[':
			return Token{Type: TK_LEFT_BRACKET, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case ']':
			return Token{Type: TK_RIGHT_BRACKET, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case ',':
			return Token{Type: TK_COMMA, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '.':
			return Token{Type: TK_DOT, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case ':':
			return Token{Type: TK_COLON, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case ';':
			return Token{Type: TK_SEMICOLON, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '?':
			return Token{Type: TK_QUESTION_MARK, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '~':
			return Token{Type: TK_BIT_NOT, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '!':
			if tk.Peek() == '=' {
				tk.s.Next()
				return Token{Type: TK_BANG_EQUAL, Position: tk.s.Position, Literal: "!="}
			}
			return Token{Type: TK_BANG, Position: tk.s.Position, Literal: "="}
		case '&':
			next := tk.Peek()
			switch next {
			case '=':
				tk.s.Next()
				return Token{Type: TK_BIT_AND_EQUAL, Position: tk.s.Position, Literal: "&="}
			case '&':
				tk.s.Next()
				return Token{Type: TK_AND, Position: tk.s.Position, Literal: "&&"}
			}
			return Token{Type: TK_BIT_AND, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '|':
			next := tk.Peek()
			switch next {
			case '=':
				tk.s.Next()
				return Token{Type: TK_BIT_OR_EQUAL, Position: tk.s.Position, Literal: "|="}
			case '|':
				tk.s.Next()
				return Token{Type: TK_OR, Position: tk.s.Position, Literal: "||"}
			}
			return Token{Type: TK_BIT_OR, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '>':
			next := tk.Peek()
			switch next {
			case '=':
				tk.s.Next()
				return Token{Type: TK_GREATER_EQUAL, Position: tk.s.Position, Literal: ">="}
			case '>':
				tk.s.Next()
				return Token{Type: TK_BIT_SHIFT_RIGHT, Position: tk.s.Position, Literal: ">>"}
			}
			return Token{Type: TK_GREATER, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '<':
			next := tk.Peek()
			switch next {
			case '=':
				tk.s.Next()
				return Token{Type: TK_LESS_EQUAL, Position: tk.s.Position, Literal: "<="}
			case '<':
				tk.s.Next()
				return Token{Type: TK_BIT_SHIFT_LEFT, Position: tk.s.Position, Literal: "<<"}
			}
			return Token{Type: TK_LESS, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '=':
			next := tk.Peek()
			switch next {
			case '=':
				tk.s.Next()
				return Token{Type: TK_EQUAL_EQUAL, Position: tk.s.Position, Literal: "=="}
			case '>':
				tk.s.Next()
				return Token{Type: TK_ARROW, Position: tk.s.Position, Literal: "=>"}
			}
			return Token{Type: TK_EQUAL, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '+':
			next := tk.Peek()
			switch next {
			case '+':
				tk.s.Next()
				return Token{Type: TK_PLUS_PLUS, Position: tk.s.Position, Literal: "++"}
			case '=':
				tk.s.Next()
				return Token{Type: TK_PLUS_EQUAL, Position: tk.s.Position, Literal: "+="}
			}
			return Token{Type: TK_PLUS, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '-':
			next := tk.Peek()
			switch next {
			case '-':
				tk.s.Next()
				return Token{Type: TK_MINUS_MINUS, Position: tk.s.Position, Literal: "++"}
			case '=':
				tk.s.Next()
				return Token{Type: TK_MINUS_EQUAL, Position: tk.s.Position, Literal: "-="}
			}
			return Token{Type: TK_MINUS, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '*':
			if tk.Peek() == '=' {
				tk.s.Next()
				return Token{Type: TK_STAR_EQUAL, Position: tk.s.Position, Literal: "*="}
			}
			return Token{Type: TK_STAR, Position: tk.s.Position, Literal: tk.s.TokenText()}
		case '/':
			if tk.Peek() == '=' {
				tk.s.Next()
				return Token{Type: TK_SLASH_EQUAL, Position: tk.s.Position, Literal: "/="}
			}
			return Token{Type: TK_SLASH, Position: tk.s.Position, Literal: tk.s.TokenText()}
		}

	}
}
