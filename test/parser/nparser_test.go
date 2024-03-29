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
		file, err := os.Open("./examples/complex_test.bs")
		if err != nil {
			t.Error(err)
		}
		parser := parser.NewNParser()
		fmt.Println(parser.String())
		ast, err := parser.Parse("complex_test.bs", file)
		repr.Println(ast)
		fmt.Println(err)
		fmt.Println("Hallo welt")
		// when
	})
}
