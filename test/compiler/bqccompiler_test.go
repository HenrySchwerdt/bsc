package compiler_test

import (
	"bsc/src/compiler"
	"bsc/src/parser"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestQBEPrograms(t *testing.T) {
	t.Run("t0_empty_0", func(t *testing.T) {
		// given
		fileName := "t0_empty_0"
		file, _ := os.Open("./examples/t0_empty_0.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)

		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 0 {
					t.Fatalf("Expected exit code %d, but got: %d", 0, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t1_literals_69", func(t *testing.T) {
		// given
		fileName := "t1_literals_69"
		file, _ := os.Open("./examples/t1_literals_69.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 69 {
					t.Fatalf("Expected exit code %d, but got: %d", 69, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t2_variable_dec_10", func(t *testing.T) {
		// given
		fileName := "t2_variable_dec_10"
		file, _ := os.Open("./examples/t2_variable_dec_10.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 10 {
					t.Fatalf("Expected exit code %d, but got: %d", 10, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t3_addition_24", func(t *testing.T) {
		// given
		fileName := "t3_addition_24"
		file, _ := os.Open("./examples/t3_addition_24.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 24 {
					t.Fatalf("Expected exit code %d, but got: %d", 24, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t4_sub_2", func(t *testing.T) {
		// given
		fileName := "t4_sub_2"
		file, _ := os.Open("./examples/t4_sub_2.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 2 {
					t.Fatalf("Expected exit code %d, but got: %d", 2, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t5_mul_45", func(t *testing.T) {
		// given
		fileName := "t5_mul_45"
		file, _ := os.Open("./examples/t5_mul_45.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 45 {
					t.Fatalf("Expected exit code %d, but got: %d", 45, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t6_equality_1", func(t *testing.T) {
		// given
		fileName := "t6_equality_1"
		file, _ := os.Open("./examples/t6_equality_1.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 1 {
					t.Fatalf("Expected exit code %d, but got: %d", 1, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t7_break_5", func(t *testing.T) {
		// given
		fileName := "t7_break_5"
		file, _ := os.Open("./examples/t7_break_5.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 5 {
					t.Fatalf("Expected exit code %d, but got: %d", 5, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t8_while_147", func(t *testing.T) {
		// given
		fileName := "t8_while_147"
		file, _ := os.Open("./examples/t8_while_147.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 147 {
					t.Fatalf("Expected exit code %d, but got: %d", 147, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t8_while_147", func(t *testing.T) {
		// given
		fileName := "t8_while_147"
		file, _ := os.Open("./examples/t8_while_147.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 147 {
					t.Fatalf("Expected exit code %d, but got: %d", 147, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})
	t.Run("t9_fib_34", func(t *testing.T) {
		// given
		fileName := "t9_fib_34"
		file, _ := os.Open("./examples/t9_fib_34.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 34 {
					t.Fatalf("Expected exit code %d, but got: %d", 34, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t10_if_100", func(t *testing.T) {
		// given
		fileName := "t10_if_100"
		file, _ := os.Open("./examples/t10_if_100.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 100 {
					t.Fatalf("Expected exit code %d, but got: %d", 100, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t11_for_96", func(t *testing.T) {
		// given
		fileName := "t11_for_96"
		file, _ := os.Open("./examples/t11_for_96.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 96 {
					t.Fatalf("Expected exit code %d, but got: %d", 96, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t12_function_8", func(t *testing.T) {
		// given
		fileName := "t12_function_8"
		file, _ := os.Open("./examples/t12_function_8.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 8 {
					t.Fatalf("Expected exit code %d, but got: %d", 8, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t13_factorial_24", func(t *testing.T) {
		// given
		fileName := "t13_factorial_24"
		file, _ := os.Open("./examples/t13_factorial_24.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
		// when
		_ = compiler.Compile(fileName, "out")
		if err != nil {
			t.Fatalf("Compilation error: %s", err)
		}
		// then
		cmd := exec.Command("./" + fileName + "/out")
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError.ExitCode())
				if exitError.ExitCode() != 24 {
					t.Fatalf("Expected exit code %d, but got: %d", 24, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t14_gcd_14", func(t *testing.T) {
		// given
		fileName := "t14_gcd_14"
		file, _ := os.Open("./examples/t14_gcd_14.bs")
		defer file.Close()
		parser := parser.NewNParser()
		ast, err := parser.Parse(file.Name(), file)
		if err != nil {
			t.Fatalf("Parsing error: %s", err)
		}
		compiler := compiler.NewBQCCompiler(ast)
		compiler.StdLibPath = "../../"
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
				if exitError.ExitCode() != 14 {
					t.Fatalf("Expected exit code %d, but got: %d", 14, exitError.ExitCode())
				}
			} else {
				t.Fatalf("Could not run the program: %s", err)
			}
		}
		os.RemoveAll(fileName)
	})

	t.Run("t16_hello_world", func(t *testing.T) {
		// given
		// fileName := "t16_hello_world"
		file, err := os.Open("./examples/t16_hello_world.bs")
		fmt.Println(file, err)
		defer file.Close()
		parser := parser.NewNParser()
		ast, _ := parser.Parse(file.Name(), file)
		s, err := json.Marshal(ast)
		// compiler := compiler.NewBQCCompiler(ast)
		fmt.Println(s, err)
		// // when
		// _ = compiler.Compile(fileName, "out")
		// if err != nil {
		// 	t.Fatalf("Compilation error: %s", err)
		// }
		// // then
		// cmd := exec.Command("./" + fileName + "/out")
		// if err := cmd.Run(); err != nil {
		// 	if exitError, ok := err.(*exec.ExitError); ok {
		// 		fmt.Println(exitError.ExitCode())
		// 		if exitError.ExitCode() != 14 {
		// 			t.Fatalf("Expected exit code %d, but got: %d", 14, exitError.ExitCode())
		// 		}
		// 	} else {
		// 		t.Fatalf("Could not run the program: %s", err)
		// 	}
		// }
		// os.RemoveAll(fileName)
	})

}
