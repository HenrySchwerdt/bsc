package compiler

import (
	"bsc/src/exeptions"
	"bsc/src/parser"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type CCompiler struct {
	Ast            parser.Node
	OutMain        strings.Builder
	OutOuter       strings.Builder
	OutIncludes    strings.Builder
	LoopCount      int
	ConditionCount int
	InLoop         bool
	InFunction     bool
}

func NewCCompiler(ast parser.Node) *CCompiler {
	var builder strings.Builder
	var builderfn strings.Builder
	var includes strings.Builder
	return &CCompiler{
		Ast:            ast,
		OutMain:        builder,
		OutOuter:       builderfn,
		OutIncludes:    includes,
		LoopCount:      0,
		ConditionCount: 0,
		InLoop:         false,
		InFunction:     false,
	}
}

func (c *CCompiler) getBuffer() *strings.Builder {
	if c.InFunction {
		return &c.OutOuter
	}
	return &c.OutMain
}

func (c *CCompiler) VisitLiteral(l *parser.Literal) error {
	c.getBuffer().WriteString(fmt.Sprintf("%d", l.Value.(int64)))
	return nil
}

func (c *CCompiler) VisitIdentifier(id *parser.Identifier) error {

	return nil
}

func (c *CCompiler) VisitUnaryExpression(ue *parser.UnaryExpression) error {
	return nil
}

func (c *CCompiler) VisitBinaryExpression(be *parser.BinaryExpression) error {
	be.Left.Accept(c)
	c.getBuffer().WriteString(be.Operator)
	be.Right.Accept(c)
	return nil
}

func (c *CCompiler) VisitConditionalExpression(ce *parser.ConditionalExpression) error {
	return nil
}

func (c *CCompiler) VisitExpressionStatement(es *parser.ExpressionStatement) error {
	return nil
}

func (c *CCompiler) VisitVariableDeclarator(vd *parser.VariableDeclarator) error {
	c.getBuffer().WriteString(fmt.Sprintf("int %s=", vd.Id))
	if err := vd.Init.Accept(c); err != nil {
		return err
	}
	c.getBuffer().WriteString(";")
	return nil
}

func (c *CCompiler) VisitVariableDeclaration(vd *parser.VariableDeclaration) error {
	return vd.Declaration.Accept(c)
}

func (c *CCompiler) VisitVariableLookup(vl *parser.VariableLookup) error {
	c.getBuffer().WriteString(vl.Id)
	return nil
}
func (c *CCompiler) VisitCallExpression(ce *parser.CallExpression) error {
	c.getBuffer().WriteString(fmt.Sprintf("%s(", ce.Identifier))
	if len(ce.Args) == 0 {
		return nil
	}
	for _, arg := range ce.Args[:len(ce.Args)-1] {
		arg.Accept(c)
		c.getBuffer().WriteString(",")
	}
	ce.Args[len(ce.Args)-1].Accept(c)
	c.getBuffer().WriteString(")")
	return nil
}

func (c *CCompiler) VisitFunction(fn *parser.Function) error {
	// Todo forward declare in some occasions
	c.InFunction = true
	c.getBuffer().WriteString(fmt.Sprintf("int %s(", fn.Id))
	fn.Parameters.Accept(c)
	c.getBuffer().WriteString("){")
	fn.Body.Accept(c)
	c.getBuffer().WriteString("}")
	c.InFunction = false
	return nil
}

func (c *CCompiler) VisitFunctionDeclaration(fd *parser.FunctionDeclaration) error {
	return fd.Function.Accept(c)
}

func (c *CCompiler) VisitIfStatment(is *parser.IfStatment) error {
	c.getBuffer().WriteString("if(")
	is.Test.Accept(c)
	c.getBuffer().WriteString(")")
	is.Consequent.Accept(c)
	if is.Alternate != nil {
		c.getBuffer().WriteString("else")
		is.Alternate.Accept(c)
	}
	return nil
}

func (c *CCompiler) VisitReturnStatment(rs *parser.ReturnStatment) error {

	c.getBuffer().WriteString("return ")
	rs.Argument.Accept(c)
	c.getBuffer().WriteString(";")
	return nil
}

func (c *CCompiler) VisitBreakStatment(bs *parser.BreakStatment) error {
	if !c.InLoop {
		return &exeptions.CompilerError{
			File:    bs.Start.Position.Filename,
			Line:    bs.Start.Position.Line,
			Column:  bs.Start.Position.Column,
			Message: "Cannot use 'break' statement outside of a loop.",
		}
	}
	c.getBuffer().WriteString("break;")
	return nil
}

func (c *CCompiler) VisitForStatment(fs *parser.ForStatment) error {
	c.getBuffer().WriteString("for(")
	if fs.Init != nil {
		fs.Init.Accept(c)
	}
	if fs.Test != nil {
		fs.Test.Accept(c)
	}
	c.getBuffer().WriteString(";")
	if fs.Update != nil {
		fs.Update.Accept(c)
		// TODO fix this
		str := c.getBuffer().String()
		if len(str) > 0 {
			// Remove the last character
			str = str[:len(str)-1]
			c.getBuffer().Reset()
			c.getBuffer().WriteString(str)
		}
	}

	c.getBuffer().WriteString(")")
	c.InLoop = true
	fs.Body.Accept(c)
	c.InLoop = false
	return nil
}

func (c *CCompiler) VisitWhileStatment(ws *parser.WhileStatment) error {

	c.getBuffer().WriteString("while(")
	ws.Test.Accept(c)
	c.getBuffer().WriteString(")")
	c.InLoop = true
	ws.Body.Accept(c)
	c.InLoop = false
	return nil
}

func (c *CCompiler) VisitBlockStatement(bs *parser.BlockStatement) error {
	c.getBuffer().WriteString("{")
	for _, ins := range bs.Instructions {
		ins.Accept(c)
	}
	c.getBuffer().WriteString("}")
	return nil
}

func (c *CCompiler) VisitExitStatment(es *parser.ExitStatment) error {
	c.getBuffer().WriteString("exit(")
	es.Argument.Accept(c)
	c.getBuffer().WriteString(");")
	return nil
}

func (c *CCompiler) VisitAssignmentStatement(as *parser.AssignmentStatement) error {
	c.getBuffer().WriteString(fmt.Sprintf("%s=", as.Identifier))
	as.Value.Accept(c)
	c.getBuffer().WriteString(";")
	return nil
}

func (c *CCompiler) VisitProgram(p *parser.Program) error {

	c.getBuffer().WriteString("int main(){")

	for _, ins := range p.Instructions {
		ins.Accept(c)
	}
	c.getBuffer().WriteString("return 0;")
	c.getBuffer().WriteString("}")
	return nil
}

func (c *CCompiler) VisitParams(pa *parser.Params) error {
	if len(pa.Args) == 0 {
		return nil
	}
	for _, arg := range pa.Args[:len(pa.Args)-1] {
		c.getBuffer().WriteString(fmt.Sprintf("int %s, ", arg))
	}
	c.getBuffer().WriteString(fmt.Sprintf("int %s", pa.Args[len(pa.Args)-1]))

	return nil
}

func (c *CCompiler) Compile(outDir, outFile string) error {
	c.OutIncludes.WriteString("#include <stdlib.h>\n")
	if err := (c.Ast).Accept(c); err != nil {
		return err
	}
	os.RemoveAll(outDir)
	os.Mkdir(outDir, 0755)
	gccPath, err := exec.LookPath("gcc")
	if err != nil {
		fmt.Println("GCC not found in the system.")
		os.Exit(1)
	}
	ioutil.WriteFile(fmt.Sprintf("%s/%s.c", outDir, outFile), []byte(c.OutIncludes.String()+c.OutOuter.String()+c.OutMain.String()), 0755)
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		_, err := io.WriteString(writer, c.OutIncludes.String()+c.OutOuter.String()+c.OutMain.String())
		if err != nil {
			fmt.Printf("Error writing to pipe: %s\n", err)
			os.Exit(1)
		}
	}()

	cmd := exec.Command(gccPath, "-x", "c", "-", "-o", fmt.Sprintf("%s/%s", outDir, outFile))
	cmd.Stdin = reader

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error while compiling C code: %s\n", err)
		os.Exit(1)
	}

	return nil
}
