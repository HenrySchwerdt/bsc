package compiler

import (
	"bsc/src/exeptions"
	"bsc/src/parser"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type variable struct {
	StackLocation int
}

type NASMElf64Compiler struct {
	Ast       parser.Node
	Out       strings.Builder
	StackSize int
	Variables map[string]*variable
}

func NewNASMElf64Compiler(ast parser.Node) *NASMElf64Compiler {
	var builder strings.Builder
	return &NASMElf64Compiler{
		Ast:       ast,
		Out:       builder,
		StackSize: 0,
		Variables: make(map[string]*variable),
	}
}

func (c *NASMElf64Compiler) push(reg string) {
	c.Out.WriteString(fmt.Sprintf("    push %s\n", reg))
	c.StackSize += 1
}

func (c *NASMElf64Compiler) pop(reg string) {
	c.Out.WriteString(fmt.Sprintf("    pop %s\n", reg))
	c.StackSize -= 1
}

func (c *NASMElf64Compiler) VisitLiteral(l *parser.Literal) error {
	c.Out.WriteString(fmt.Sprintf("    mov rax, %d\n", l.Value.(int64)))
	c.push("rax")
	return nil
}

func (c *NASMElf64Compiler) VisitIdentifier(id *parser.Identifier) error {
	return nil
}

func (c *NASMElf64Compiler) VisitUnaryExpression(ue *parser.UnaryExpression) error {
	return nil
}

func (c *NASMElf64Compiler) VisitBinaryExpression(be *parser.BinaryExpression) error {
	be.Left.Accept(c)
	be.Right.Accept(c)
	c.pop("rax")
	c.pop("rbx")
	c.Out.WriteString("    add rax, rbx\n")
	c.push("rax")
	return nil
}

func (c *NASMElf64Compiler) VisitConditionalExpression(ce *parser.ConditionalExpression) error {
	return nil
}

func (c *NASMElf64Compiler) VisitExpressionStatement(es *parser.ExpressionStatement) error {
	return nil
}

func (c *NASMElf64Compiler) VisitVariableDeclarator(vd *parser.VariableDeclarator) error {
	if _, exists := c.Variables[vd.Id]; exists {
		return &exeptions.CompilerError{
			File:    "Bla",
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("CompileError: Cannot declare a variable that already exists: '%s'", vd.Id),
		}
	}
	c.Variables[vd.Id] = &variable{
		StackLocation: c.StackSize,
	}
	return vd.Init.Accept(c)
}

func (c *NASMElf64Compiler) VisitVariableDeclaration(vd *parser.VariableDeclaration) error {
	return vd.Declaration.Accept(c)
}

func (c *NASMElf64Compiler) VisitVariableLookup(vl *parser.VariableLookup) error {
	if _, exists := c.Variables[vl.Id]; !exists {
		return &exeptions.CompilerError{
			File:    "Bla",
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("CompileError: Undeclared Identifier '%s'", vl.Id),
		}
	}
	c.push(fmt.Sprintf("QWORD [rsp + %d]", (c.StackSize-c.Variables[vl.Id].StackLocation-1)*8))
	return nil
}

func (c *NASMElf64Compiler) VisitFunction(fn *parser.Function) error {
	return nil
}

func (c *NASMElf64Compiler) VisitFunctionDeclaration(fd *parser.FunctionDeclaration) error {
	return nil
}

func (c *NASMElf64Compiler) VisitForStatment(fs *parser.ForStatment) error {
	return nil
}

func (c *NASMElf64Compiler) VisitIfStatment(is *parser.IfStatment) error {
	return nil
}

func (c *NASMElf64Compiler) VisitReturnStatment(rs *parser.ReturnStatment) error {
	return nil
}

func (c *NASMElf64Compiler) VisitWhileStatment(ws *parser.WhileStatment) error {
	return nil
}

func (c *NASMElf64Compiler) VisitBlockStatement(bs *parser.BlockStatement) error {
	return nil
}

func (c *NASMElf64Compiler) VisitExitStatment(es *parser.ExitStatment) error {
	if err := es.Argument.Accept(c); err != nil {
		return err
	}
	c.Out.WriteString("    mov rax, 60\n")
	c.pop("rdi")
	c.Out.WriteString("    syscall\n")
	return nil
}

func (c *NASMElf64Compiler) VisitAssignmentStatement(as *parser.AssignmentStatement) error {
	return nil
}

func (c *NASMElf64Compiler) VisitProgram(p *parser.Program) error {
	c.Out.WriteString("global _start\n")
	c.Out.WriteString("_start:\n")
	for _, stmt := range p.Instructions {
		if err := stmt.Accept(c); err != nil {
			return err
		}
	}
	c.Out.WriteString("    mov rax, 60\n")
	c.Out.WriteString("    mov rdi, 0\n")
	c.Out.WriteString("    syscall\n")
	return nil
}

func (c *NASMElf64Compiler) Compile(outDir, outFile string) error {
	if err := (c.Ast).Accept(c); err != nil {
		return err
	}
	os.RemoveAll(outDir)
	os.Mkdir(outDir, 0755)
	ioutil.WriteFile(fmt.Sprintf("%s/%s.asm", outDir, outFile), []byte(c.Out.String()), 0755)

	cmdNasm := exec.Command("nasm", "-felf64", fmt.Sprintf("%s/%s.asm", outDir, outFile))
	err := cmdNasm.Run()
	if err != nil {
		return fmt.Errorf("Failed to run nasm: %v", err)
	}

	cmdLd := exec.Command("ld", "-o", fmt.Sprintf("%s/%s", outDir, outFile), fmt.Sprintf("%s/%s.o", outDir, outFile))

	var stderr bytes.Buffer
	cmdLd.Stderr = &stderr
	err = cmdLd.Run()
	if err != nil {
		return fmt.Errorf("Failed to run ld: %v\nLinker output:\n%s", err, stderr.String())
	}

	return nil
}
