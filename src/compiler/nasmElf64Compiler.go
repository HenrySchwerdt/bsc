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
type function struct {
	Params []string
	Ret    string
}

type NASMElf64Compiler struct {
	Ast            parser.Node
	StackManager   *StackManager
	Out            strings.Builder
	OutFn          strings.Builder
	LoopCount      int
	ConditionCount int
	InLoop         bool
	InFunction     bool
}

func NewNASMElf64Compiler(ast parser.Node) *NASMElf64Compiler {
	var builder strings.Builder
	var builderfn strings.Builder
	return &NASMElf64Compiler{
		Ast:            ast,
		StackManager:   NewStackManager(),
		Out:            builder,
		OutFn:          builderfn,
		LoopCount:      0,
		ConditionCount: 0,
		InLoop:         false,
		InFunction:     false,
	}
}

func (c *NASMElf64Compiler) WriteString(str string) {
	if c.InFunction {
		c.OutFn.WriteString(str)
	} else {
		c.Out.WriteString(str)
	}
}

func (c *NASMElf64Compiler) GetOut() *strings.Builder {
	if c.InFunction {
		return &c.OutFn
	} else {
		return &c.Out
	}
}

func (c *NASMElf64Compiler) push(reg string) {
	c.WriteString(fmt.Sprintf("    push %s\n", reg))
	c.StackManager.Push(8)
}

func (c *NASMElf64Compiler) pop(reg string) {
	c.WriteString(fmt.Sprintf("    pop %s\n", reg))
	c.StackManager.Pop(8)
}

func (c *NASMElf64Compiler) VisitLiteral(l *parser.Literal) error {
	c.WriteString(fmt.Sprintf("    mov rax, %d\n", l.Value.(int64)))
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
		c.WriteString("    add rax, rbx\n")
		c.push("rax")
		return nil
	case "-":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    sub rax, rbx\n")
		c.push("rax")
		return nil
	case "&":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    and rax, rbx\n")
		c.push("rax")
		return nil
	case "|":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    or rax, rbx\n")
		c.push("rax")
		return nil
	case "*":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    mul rbx\n")
		c.push("rax")
		return nil
	case "==":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    cmp rax, rbx\n")
		c.WriteString("    sete al\n") // Set AL if equal
		c.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case "!=":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    cmp rax, rbx\n")
		c.WriteString("    setne al\n") // Set AL if not equal
		c.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case "<":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    cmp rax, rbx\n")
		c.WriteString("    setl al\n") // Set AL if less than
		c.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case ">":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    cmp rax, rbx\n")
		c.WriteString("    setg al\n")
		c.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case "<=":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    cmp rax, rbx\n")
		c.WriteString("    setle al\n")
		c.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	case ">=":
		be.Left.Accept(c)
		be.Right.Accept(c)
		c.pop("rbx")
		c.pop("rax")
		c.WriteString("    cmp rax, rbx\n")
		c.WriteString("    setge al\n")
		c.WriteString("    movzx rax, al\n")
		c.push("rax")
		return nil
	default:
		return &exeptions.CompilerError{
			File:    be.Start.Position.Filename,
			Line:    be.Start.Position.Line,
			Column:  be.Start.Position.Column,
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
	if _, exists := c.StackManager.GetVariableOffset(vd.Id); exists {
		return &exeptions.CompilerError{
			File:    vd.Start.Position.Filename,
			Line:    vd.Start.Position.Line,
			Column:  vd.Start.Position.Column,
			Message: fmt.Sprintf("CompileError: Cannot declare a variable that already exists: '%s'", vd.Id),
		}
	}
	if err := vd.Init.Accept(c); err != nil {
		return err
	}
	c.StackManager.AddVariable(vd.Id, 8)
	return nil
}

func (c *NASMElf64Compiler) VisitVariableDeclaration(vd *parser.VariableDeclaration) error {
	return vd.Declaration.Accept(c)
}

func (c *NASMElf64Compiler) VisitVariableLookup(vl *parser.VariableLookup) error {
	variable, exists := c.StackManager.GetVariableOffset(vl.Id)
	if !exists {
		return &exeptions.CompilerError{
			File:    vl.Start.Position.Filename,
			Line:    vl.Start.Position.Line,
			Column:  vl.Start.Position.Column,
			Message: fmt.Sprintf("CompileError: Undeclared Identifier '%s'", vl.Id),
		}
	}
	if variable > 0 {
		c.WriteString(fmt.Sprintf("    mov rax, QWORD [rbp + %d]\n", variable))
	} else {
		c.WriteString(fmt.Sprintf("    mov rax, QWORD [rbp - %d]\n", variable*-1))
	}

	c.push("rax")
	return nil
}
func (c *NASMElf64Compiler) VisitCallExpression(ce *parser.CallExpression) error {
	for _, arg := range ce.Args {
		arg.Accept(c)
	}
	c.WriteString(fmt.Sprintf("    call %s\n", ce.Identifier))
	// TODO cleaner solution for that instead of poping all the args from the stack
	for i := 0; i < len(ce.Args); i++ {
		c.pop("rbx")
	}
	c.push("rax") // to put the return value on the stack like the other expressions
	return nil
}

func (c *NASMElf64Compiler) VisitFunction(fn *parser.Function) error {
	c.InFunction = true
	c.WriteString(fmt.Sprintf("%s:\n", fn.Id)) // function label
	c.StackManager.StartNewStackFrame(c.GetOut())
	fn.Parameters.Accept(c)
	fn.Body.Accept(c)
	c.StackManager.CloseCurrentStackFrame(c.GetOut())
	c.WriteString("    ret\n")

	c.InFunction = false
	return nil
}

func (c *NASMElf64Compiler) VisitFunctionDeclaration(fd *parser.FunctionDeclaration) error {
	return fd.Function.Accept(c)
}

func (c *NASMElf64Compiler) VisitIfStatment(is *parser.IfStatment) error {
	c.ConditionCount += 1

	if err := is.Test.Accept(c); err != nil {
		return err
	}
	tmpStackCount := c.StackManager.StackSize / 8
	c.pop("rax")
	tmpCondCount := c.ConditionCount

	c.WriteString("    test rax, rax\n")
	c.WriteString(fmt.Sprintf("    jz .CE%d\n", tmpCondCount))
	is.Consequent.Accept(c)
	c.WriteString(fmt.Sprintf("    jmp .CF%d\n", tmpCondCount))
	c.WriteString(fmt.Sprintf(".CE%d:\n", tmpCondCount))
	if is.Alternate != nil {
		is.Alternate.Accept(c)
	}
	c.WriteString(fmt.Sprintf(".CF%d:\n", tmpCondCount))
	for tmpStackCount < c.StackManager.StackSize/8 {
		c.pop("rax")
	}
	return nil
}

func (c *NASMElf64Compiler) VisitReturnStatment(rs *parser.ReturnStatment) error {
	rs.Argument.Accept(c)
	c.pop("rax")
	c.WriteString("    mov rsp, rbp\n")
	c.WriteString("    pop rbp\n")
	c.WriteString("    ret\n")
	return nil
}

func (c *NASMElf64Compiler) VisitBreakStatment(bs *parser.BreakStatment) error {
	if !c.InLoop {
		return &exeptions.CompilerError{
			File:    bs.Start.Position.Filename,
			Line:    bs.Start.Position.Line,
			Column:  bs.Start.Position.Column,
			Message: "Cannot use 'break' statement outside of a loop.",
		}
	}
	c.Out.WriteString(fmt.Sprintf("    jmp .E%d\n", c.LoopCount))
	return nil
}

func (c *NASMElf64Compiler) VisitForStatment(fs *parser.ForStatment) error {
	tmpStackCount := c.StackManager.StackSize / 8
	c.LoopCount += 1
	tmpCount := c.LoopCount
	if err := fs.Init.Accept(c); err != nil {
		return err
	}
	c.WriteString(fmt.Sprintf(".L%d:\n", tmpCount))

	if err := fs.Test.Accept(c); err != nil {
		return err
	}
	c.pop("rax")
	c.WriteString("    test rax, rax\n")
	c.WriteString(fmt.Sprintf("    jz .E%d\n", tmpCount))
	c.InLoop = true
	if err := fs.Body.Accept(c); err != nil {
		return err
	}
	c.InLoop = false
	if err := fs.Update.Accept(c); err != nil {
		return err
	}

	c.WriteString(fmt.Sprintf("    jmp .L%d\n", tmpCount))
	c.WriteString(fmt.Sprintf(".E%d:\n", tmpCount))
	for tmpStackCount < c.StackManager.StackSize/8 {
		c.pop("rax")
	}
	return nil
}

func (c *NASMElf64Compiler) VisitWhileStatment(ws *parser.WhileStatment) error {
	c.LoopCount += 1
	tmpCount := c.LoopCount
	tmpStackCount := c.StackManager.StackSize / 8
	c.WriteString(fmt.Sprintf(".L%d:\n", tmpCount))

	if err := ws.Test.Accept(c); err != nil {
		return err
	}
	c.pop("rax")
	c.WriteString("    test rax, rax\n")
	c.WriteString(fmt.Sprintf("    jz .E%d\n", tmpCount))
	c.InLoop = true
	if err := ws.Body.Accept(c); err != nil {
		return err
	}
	c.InLoop = false
	for tmpStackCount < c.StackManager.StackSize/8 {
		c.pop("rax")
	}
	c.WriteString(fmt.Sprintf("    jmp .L%d\n", tmpCount))
	c.WriteString(fmt.Sprintf(".E%d:\n", tmpCount))

	return nil
}

func (c *NASMElf64Compiler) VisitBlockStatement(bs *parser.BlockStatement) error {
	for _, node := range bs.Instructions {
		if err := node.Accept(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *NASMElf64Compiler) VisitExitStatment(es *parser.ExitStatment) error {
	if err := es.Argument.Accept(c); err != nil {
		return err
	}
	c.WriteString("    mov rax, 60\n")
	c.pop("rdi")
	c.WriteString("    syscall\n")
	return nil
}

func (c *NASMElf64Compiler) VisitAssignmentStatement(as *parser.AssignmentStatement) error {
	variable, exists := c.StackManager.GetVariableOffset(as.Identifier)
	if !exists {
		return &exeptions.CompilerError{
			File:    as.Start.Position.Filename,
			Line:    as.Start.Position.Line,
			Column:  as.Start.Position.Column,
			Message: fmt.Sprintf("CompileError: Cannot assign to an unassigned variable '%s'", as.Identifier),
		}
	}
	as.Value.Accept(c)
	c.pop("rax")
	if variable > 0 {
		c.WriteString(fmt.Sprintf("    mov QWORD [rbp + %d], rax\n", variable))
	} else {
		c.WriteString(fmt.Sprintf("    mov QWORD [rbp - %d], rax\n", variable*-1))
	}

	return nil
}

func (c *NASMElf64Compiler) VisitProgram(p *parser.Program) error {
	c.WriteString("global _start\n")
	c.WriteString("_start:\n")
	c.StackManager.StartNewStackFrame(c.GetOut())
	for _, stmt := range p.Instructions {
		if err := stmt.Accept(c); err != nil {
			return err
		}
	}
	c.StackManager.CloseCurrentStackFrame(c.GetOut())
	c.WriteString("    mov rax, 60\n")
	c.WriteString("    mov rdi, 0\n")
	c.WriteString("    syscall\n")
	return nil
}

func (c *NASMElf64Compiler) VisitParams(pa *parser.Params) error {

	for i, param := range pa.Args {
		c.StackManager.AddParam(param, len(pa.Args), i, 8)
	}
	return nil
}

func (c *NASMElf64Compiler) Compile(outDir, outFile string) error {
	if err := (c.Ast).Accept(c); err != nil {
		return err
	}
	os.RemoveAll(outDir)
	os.Mkdir(outDir, 0755)
	ioutil.WriteFile(fmt.Sprintf("%s/%s.asm", outDir, outFile), []byte(c.Out.String()+c.OutFn.String()), 0755)

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
