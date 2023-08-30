package compiler_test

import (
	"bsc/src/compiler"
	"bsc/src/lexer"
	"bsc/src/parser"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestBsPrograms(t *testing.T) {
	t.Run("Should compile and exit with exitcode 0", func(t *testing.T) {
		// given
		fileName := "t0_empty"
		file, _ := os.Open("./examples/" + fileName + ".bs")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		parser := parser.NewParser(tokenizer)
		ast, _ := parser.Parse()
		compiler := compiler.NewNASMElf64Compiler(ast)

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
				if exitError.ExitCode() != 0 {
					t.Fatalf("Expected exit code 0, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 69", func(t *testing.T) {
		// given
		fileName := "t1_literals"
		file, _ := os.Open("./examples/" + fileName + ".bs")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		parser := parser.NewParser(tokenizer)
		ast, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
		}
		compiler := compiler.NewNASMElf64Compiler(ast)

		// when
		err = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 69 {
					t.Fatalf("Expected exit code 69, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 10", func(t *testing.T) {
		// given
		fileName := "t2_variable_dec"
		file, _ := os.Open("./examples/" + fileName + ".bs")
		defer file.Close()
		tokenizer := lexer.NewTokenizer(file)
		parser := parser.NewParser(tokenizer)
		ast, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
		}
		compiler := compiler.NewNASMElf64Compiler(ast)

		// when
		err = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 10 {
					t.Fatalf("Expected exit code 10, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
}
