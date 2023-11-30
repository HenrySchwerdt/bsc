package compiler_test

import (
	"bsc/src/compiler"
	"bsc/src/lexer"
	"bsc/src/parser"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestBsCPrograms(t *testing.T) {
	files, _ := ioutil.ReadDir("./examples")
	for _, file := range files {
		exitCode, err := extractExitCode(file.Name())
		folderName := trimExtension(file.Name())
		if err != nil {
			t.Fatal(err)
		}
		t.Run(file.Name(), func(t *testing.T) {
			// given
			fileName := folderName
			file, _ := os.Open("./examples/" + fileName + ".bs")
			defer file.Close()
			tokenizer := lexer.NewTokenizer(file)
			parser := parser.NewParser(tokenizer)
			ast, _ := parser.Parse()
			compiler := compiler.NewCCompiler(ast)

			// when
			err := compiler.Compile(fileName, "out")
			if err != nil {
				t.Fatalf("Compilation error: %s", err)
			}
			// then
			cmd := exec.Command("./" + fileName + "/out")
			if err := cmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					fmt.Println(exitError.ExitCode())
					if exitError.ExitCode() != exitCode {
						t.Fatalf("Expected exit code %d, but got: %d", exitCode, exitError.ExitCode())
					}
				} else {
					t.Fatalf("Could not run the program: %s", err)
				}
			}
			os.RemoveAll(fileName)
		})
	}
}
