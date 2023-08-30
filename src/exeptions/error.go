package exeptions

import "fmt"

type CompilerError struct {
	File    string
	Line    int
	Column  int
	Message string
}

func (e *CompilerError) Error() string {
	return fmt.Sprintf("\033[31mError at:\033[0m %s:%d:%d: %s", e.File, e.Line, e.Column, e.Message)
}
