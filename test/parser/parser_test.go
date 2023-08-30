package parser_test

import (
	"bsc/src/lexer"
	"bsc/src/parser"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestSingleCharTokens(t *testing.T) {
	t.Run("Should correctly identify all single character tokens.", func(t *testing.T) {
		// given
		file, _ := os.Open("./examples/simple_test.bs")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		parser := parser.NewParser(tokenizer)
		// when
		ast, parseErr := parser.Parse()
		fmt.Println(parseErr)
		jsonData, err := json.MarshalIndent(ast, "", "	")
		fmt.Println(err, string(jsonData))
	})
}
