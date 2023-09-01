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
	Ast            parser.Node
	Out            strings.Builder
	StackSize      int
	Variables      map[string]*variable
	VariableScopes []map[string]*variable
	LoopCount      int
}

func NewNASMElf64Compiler(ast parser.Node) *NASMElf64Compiler {
	var builder strings.Builder
	return &NASMElf64Compiler{
		Ast:            ast,
		Out:            builder,
		StackSize:      0,
		Variables:      make(map[string]*variable),
		VariableScopes: []map[string]*variable{make(map[string]*variable)},
		LoopCount:      0,
	}
}

func (c *NASMElf64Compiler) findVariable(id string) (*variable, bool) {
	for i := len(c.VariableScopes) - 1; i >= 0; i-- {
		if variable, exists := c.VariableScopes[i][id]; exists {
			return variable, true
		}
	}
	return nil, false
}
func (c *NASMElf64Compiler) push(reg string) {
	c.Out.WriteString(fmt.Sprintf("    push %s\n", reg))
	c.StackSize++
}

func (c *NASMElf64Compiler) pop(reg string) {
	c.Out.WriteString(fmt.Sprintf("    pop %s\n", reg))
	c.StackSize--
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
	switch be.Operator {
	case "+":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    add rax, rbx\n")
		c.push("rax")
		return nil
	case "-":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    sub rax, rbx\n")
		c.push("rax")
		return nil
	case "&":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    and rax, rbx\n")
		c.push("rax")
		return nil
	case "|":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    or rax, rbx\n")
		c.push("rax")
		return nil
	case "*":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    mul rbx\n")
		c.push("rax")
		return nil
	case "==":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    cmp rax, rbx\n")
		c.Out.WriteString("    sete al\n") // Set AL if equal
		c.Out.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case "!=":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    cmp rax, rbx\n")
		c.Out.WriteString("    setne al\n") // Set AL if not equal
		c.Out.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case "<":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    cmp rax, rbx\n")
		c.Out.WriteString("    setl al\n") // Set AL if less than
		c.Out.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case ">":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    cmp rax, rbx\n")
		c.Out.WriteString("    setg al\n")
		c.Out.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case "<=":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    cmp rax, rbx\n")
		c.Out.WriteString("    setle al\n")
		c.Out.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case ">=":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.Out.WriteString("    cmp rax, rbx\n")
		c.Out.WriteString("    setge al\n")
		c.Out.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	default:
		return &exeptions.CompilerError{
			File:    "Bla",
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("CompileError: Unkown BinaryExpression Operator: '%s'", be.Operator),
		}
	}

}

func (c *NASMElf64Compiler) VisitConditionalExpression(ce *parser.ConditionalExpression) error {
	return nil
}

func (c *NASMElf64Compiler) VisitExpressionStatement(es *parser.ExpressionStatement) error {
	return nil
}

func (c *NASMElf64Compiler) VisitVariableDeclarator(vd *parser.VariableDeclarator) error {
	currentScope := c.VariableScopes[len(c.VariableScopes)-1]
	if _, exists := currentScope[vd.Id]; exists {
		return &exeptions.CompilerError{
			File:    "Bla",
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("CompileError: Cannot declare a variable that already exists: '%s'", vd.Id),
		}
	}
	if err := vd.Init.Accept(c); err != nil {
		return err
	}
	currentScope[vd.Id] = &variable{
		StackLocation: c.StackSize,
	}
	return nil
}

func (c *NASMElf64Compiler) VisitVariableDeclaration(vd *parser.VariableDeclaration) error {
	return vd.Declaration.Accept(c)
}

func (c *NASMElf64Compiler) VisitVariableLookup(vl *parser.VariableLookup) error {
	variable, exists := c.findVariable(vl.Id)
	if !exists {
		return &exeptions.CompilerError{
			File:    "Bla",
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("CompileError: Undeclared Identifier '%s'", vl.Id),
		}
	}
	c.Out.WriteString(fmt.Sprintf("    mov rax, QWORD [rbp - %d]\n", (variable.StackLocation)*8))
	c.push("rax")
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
	c.LoopCount += 1
	tmpCount := c.LoopCount
	tmpStackCount := c.StackSize

	c.Out.WriteString(fmt.Sprintf(".L%d:\n", tmpCount))

	if err := ws.Test.Accept(c); err != nil {
		return err
	}
	c.pop("rax")
	c.Out.WriteString("    test rax, rax\n")
	c.Out.WriteString(fmt.Sprintf("    jz .E%d\n", tmpCount))

	if err := ws.Body.Accept(c); err != nil {
		return err
	}
	fmt.Println(tmpStackCount, c.StackSize)
	for ; tmpStackCount < c.StackSize; c.StackSize-- {
		c.pop("rax")
	}
	c.Out.WriteString(fmt.Sprintf("    jmp .L%d\n", tmpCount))
	c.Out.WriteString(fmt.Sprintf(".E%d:\n", tmpCount))
	return nil
}

func (c *NASMElf64Compiler) VisitBlockStatement(bs *parser.BlockStatement) error {
	c.VariableScopes = append(c.VariableScopes, make(map[string]*variable))
	for _, node := range bs.Instructions {
		if err := node.Accept(c); err != nil {
			return err
		}
	}
	c.VariableScopes = c.VariableScopes[:len(c.VariableScopes)-1]
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
	variable, exists := c.findVariable(as.Identifier)
	if !exists {
		return &exeptions.CompilerError{
			File:    "Bla",
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("CompileError: Cannot assign to an unassigned variable '%s'", as.Identifier),
		}
	}
	as.Value.Accept(c)
	c.pop("rax")
	c.Out.WriteString(fmt.Sprintf("    mov QWORD [rbp - %d], rax\n", (variable.StackLocation)*8))
	return nil
}

func (c *NASMElf64Compiler) VisitProgram(p *parser.Program) error {
	c.Out.WriteString("global _start\n")
	c.Out.WriteString("_start:\n")
	c.Out.WriteString("    push rbp\n")     // Save current base pointer
	c.Out.WriteString("    mov rbp, rsp\n") // Set new base pointer to current stack pointer
	for _, stmt := range p.Instructions {
		if err := stmt.Accept(c); err != nil {
			return err
		}
	}
	c.Out.WriteString("    mov rsp, rbp\n")
	c.Out.WriteString("    pop rbp\n")
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
