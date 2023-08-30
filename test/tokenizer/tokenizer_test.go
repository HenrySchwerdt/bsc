package tokenizer_test

import (
	"bsc/src/lexer"
	"os"
	"reflect"
	"testing"
)

func TestSingleCharTokens(t *testing.T) {
	t.Run("Should correctly identify all single character tokens.", func(t *testing.T) {
		// given
		file, _ := os.Open("./examples/singleCharTokens.txt")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		// when
		tokenTypes := make([]lexer.TokenType, 0)
		for {
			tk := tokenizer.GetToken()
			if tk.Type == lexer.TK_EOF {
				tokenTypes = append(tokenTypes, tk.Type)
				break
			} else {
				tokenTypes = append(tokenTypes, tk.Type)
			}
		}
		// then
		expectedTokens := []lexer.TokenType{lexer.TK_LEFT_PAREN, lexer.TK_RIGHT_PAREN,
			lexer.TK_LEFT_BRACE, lexer.TK_RIGHT_BRACE,
			lexer.TK_LEFT_BRACKET, lexer.TK_RIGHT_BRACKET,
			lexer.TK_COMMA,
			lexer.TK_DOT,
			lexer.TK_BIT_NOT,
			lexer.TK_BANG,
			lexer.TK_COLON,
			lexer.TK_SEMICOLON,
			lexer.TK_QUESTION_MARK,
			lexer.TK_MINUS,
			lexer.TK_PLUS,
			lexer.TK_STAR,
			lexer.TK_SLASH,
			lexer.TK_EQUAL,
			lexer.TK_BIT_OR,
			lexer.TK_BIT_AND,
			lexer.TK_LESS,
			lexer.TK_GREATER,
			lexer.TK_EOF,
		}
		if len(expectedTokens) != len(tokenTypes) {
			t.Errorf("got len %d wanted len %d", len(tokenTypes), len(expectedTokens))
		}
		if !reflect.DeepEqual(tokenTypes, expectedTokens) {
			t.Errorf("got %v wanted %v", tokenTypes, expectedTokens)
		}
	})
}

func TestDoubleCharTokens(t *testing.T) {
	t.Run("Should correctly identify all double character tokens.", func(t *testing.T) {
		// given
		file, _ := os.Open("./examples/doubleCharTokens.txt")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		// when
		tokenTypes := make([]lexer.TokenType, 0)
		for {
			tk := tokenizer.GetToken()
			if tk.Type == lexer.TK_EOF {
				tokenTypes = append(tokenTypes, tk.Type)
				break
			} else {
				tokenTypes = append(tokenTypes, tk.Type)
			}
		}
		// then
		expectedTokens := []lexer.TokenType{lexer.TK_EQUAL_EQUAL, lexer.TK_PLUS_PLUS,
			lexer.TK_MINUS_MINUS, lexer.TK_PLUS_EQUAL,
			lexer.TK_MINUS_EQUAL, lexer.TK_SLASH_EQUAL,
			lexer.TK_STAR_EQUAL,
			lexer.TK_BIT_OR_EQUAL,
			lexer.TK_BIT_AND_EQUAL,
			lexer.TK_AND,
			lexer.TK_OR,
			lexer.TK_GREATER_EQUAL,
			lexer.TK_LESS_EQUAL,
			lexer.TK_ARROW,
			lexer.TK_BIT_SHIFT_LEFT,
			lexer.TK_BIT_SHIFT_RIGHT,
			lexer.TK_EOF,
		}
		if len(expectedTokens) != len(tokenTypes) {
			t.Errorf("got len %d wanted len %d", len(tokenTypes), len(expectedTokens))
		}
		if !reflect.DeepEqual(tokenTypes, expectedTokens) {
			t.Errorf("got %v wanted %v", tokenTypes, expectedTokens)
		}
	})
}

func TestKeywordTokens(t *testing.T) {
	t.Run("Should correctly identify all keyword tokens.", func(t *testing.T) {
		// given
		file, _ := os.Open("./examples/keywordTokens.txt")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		// when
		tokenTypes := make([]lexer.TokenType, 0)
		for {
			tk := tokenizer.GetToken()
			if tk.Type == lexer.TK_EOF {
				tokenTypes = append(tokenTypes, tk.Type)
				break
			} else {
				tokenTypes = append(tokenTypes, tk.Type)
			}
		}
		// then
		expectedTokens := []lexer.TokenType{lexer.TK_CLASS, lexer.TK_STRUCT,
			lexer.TK_IF, lexer.TK_ELSE,
			lexer.TK_TRUE, lexer.TK_FALSE,
			lexer.TK_FOR,
			lexer.TK_WHILE,
			lexer.TK_FN,
			lexer.TK_NULL,
			lexer.TK_RETURN,
			lexer.TK_SUPER,
			lexer.TK_THIS,
			lexer.TK_VAL,
			lexer.TK_VAR,
			lexer.TK_EXIT,
			lexer.TK_EOF,
		}
		if len(expectedTokens) != len(tokenTypes) {
			t.Errorf("got len %d wanted len %d", len(tokenTypes), len(expectedTokens))
		}
		if !reflect.DeepEqual(tokenTypes, expectedTokens) {
			t.Errorf("got %v wanted %v", tokenTypes, expectedTokens)
		}
	})
}

func TestValueTokens(t *testing.T) {
	t.Run("Should correctly identify all value tokens.", func(t *testing.T) {
		// given
		file, _ := os.Open("./examples/valueTokens.txt")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		// when
		tokenTypes := make([]lexer.TokenType, 0)
		for {
			tk := tokenizer.GetToken()
			if tk.Type == lexer.TK_EOF {
				tokenTypes = append(tokenTypes, tk.Type)
				break
			} else {
				tokenTypes = append(tokenTypes, tk.Type)
			}
		}
		// then
		expectedTokens := []lexer.TokenType{
			lexer.TK_STRING,
			lexer.TK_INTEGER,
			lexer.TK_FLOAT,
			lexer.TK_FLOAT,
			lexer.TK_FLOAT,
			lexer.TK_CHAR,
			lexer.TK_EOF,
		}
		if len(expectedTokens) != len(tokenTypes) {
			t.Errorf("got len %d wanted len %d", len(tokenTypes), len(expectedTokens))
		}
		if !reflect.DeepEqual(tokenTypes, expectedTokens) {
			t.Errorf("got %v wanted %v", tokenTypes, expectedTokens)
		}
	})
}

func TestIdentifierTokens(t *testing.T) {
	t.Run("Should correctly identify all identifier tokens.", func(t *testing.T) {
		// given
		file, _ := os.Open("./examples/indentifierTokens.txt")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		// when
		tokenTypes := make([]lexer.TokenType, 0)
		for {
			tk := tokenizer.GetToken()
			if tk.Type == lexer.TK_EOF {
				tokenTypes = append(tokenTypes, tk.Type)
				break
			} else {
				tokenTypes = append(tokenTypes, tk.Type)
			}
		}
		// then
		expectedTokens := []lexer.TokenType{
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_IDENTIFIER,
			lexer.TK_EOF,
		}
		if len(expectedTokens) != len(tokenTypes) {
			t.Errorf("got len %d wanted len %d", len(tokenTypes), len(expectedTokens))
		}
		if !reflect.DeepEqual(tokenTypes, expectedTokens) {
			t.Errorf("got %v wanted %v", tokenTypes, expectedTokens)
		}
	})
}
