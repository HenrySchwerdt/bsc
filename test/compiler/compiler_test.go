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
	t.Run("Should compile and exit with exitcode 24", func(t *testing.T) {
		// given
		fileName := "t3_addition"
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
				if exitError.ExitCode() != 24 {
					t.Fatalf("Expected exit code 24, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 2", func(t *testing.T) {
		// given
		fileName := "t4_sub"
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
				if exitError.ExitCode() != 2 {
					t.Fatalf("Expected exit code 2, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 45", func(t *testing.T) {
		// given
		fileName := "t5_mul"
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
				if exitError.ExitCode() != 45 {
					t.Fatalf("Expected exit code 45, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 1", func(t *testing.T) {
		// given
		fileName := "t6_equality"
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
				if exitError.ExitCode() != 1 {
					t.Fatalf("Expected exit code 1, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 147", func(t *testing.T) {
		// given
		fileName := "t8_sum_to_50"
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
				if exitError.ExitCode() != 147 {
					t.Fatalf("Expected exit code 147, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 34", func(t *testing.T) {
		// given
		fileName := "t9_fib"
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
				if exitError.ExitCode() != 34 {
					t.Fatalf("Expected exit code 34, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("Should compile and exit with exitcode 100", func(t *testing.T) {
		// given
		fileName := "t10_if"
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
				if exitError.ExitCode() != 100 {
					t.Fatalf("Expected exit code 100, but got: %d", exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
}
