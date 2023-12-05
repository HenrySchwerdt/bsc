package parser_test

import (
	"bsc/src/parser"
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/repr"
)

func TestShouldParseEverything(t *testing.T) {
	t.Run("Should bla.", func(t *testing.T) {
		// given
		file, _ := os.Open("./examples/complex_test.bs")
		parser := parser.NewNParser()
		ast, err := parser.Parse("complex_test.bs", file)
		fmt.Println(len(ast.Statements))
		repr.Println(ast)
		fmt.Println(err)
		// when
	})
}
