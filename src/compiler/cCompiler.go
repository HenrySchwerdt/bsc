package compiler

import (
	"bsc/src/parser"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type CCompiler struct {
	Ast            parser.Program
	OutMain        strings.Builder
	OutOuter       strings.Builder
	OutIncludes    strings.Builder
	LoopCount      int
	ConditionCount int
	InLoop         bool
	InFunction     bool
}

func NewCCompiler(ast *parser.Program) *CCompiler {
	var builder strings.Builder
	var builderfn strings.Builder
	var includes strings.Builder
	return &CCompiler{
		Ast:            *ast,
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

func convertType(typea string) (string, error) {
	switch typea {
	case "int8":
		return "char", nil
	case "int16":
		return "short", nil
	case "int32":
		return "int", nil
	case "int64":
		return "long long", nil
	case "uint8":
		return "unsigned char", nil
	case "uint16":
		return "unsigned short", nil
	case "uint32":
		return "unsigned int", nil
	case "uint64":
		return "unsigned long long", nil
	case "float32":
		return "float", nil
	case "float64":
		return "double", nil
	case "void":
		return "void", nil
	}
	return "", errors.New(fmt.Sprintf("No such type %s", typea))

}

func (c *CCompiler) VisitProgram(t *parser.Program) error {
	c.getBuffer().WriteString("int main(){")
	for _, stmt := range t.Statements {
		err := stmt.Accept(c)
		if err != nil {
			return err
		}
	}
	c.getBuffer().WriteString("return 0;}")
	return nil
}
func (c *CCompiler) VisitVarDeclarationStatement(t *parser.VarDeclarationStatement) error {
	err := t.Dec.Accept(c)
	if err != nil {
		return err
	}
	c.getBuffer().WriteString(";")
	return nil
}
func (c *CCompiler) VisitAssignmentStatement(t *parser.AssignmentStatement) error {
	err := t.Ass.Accept(c)
	if err != nil {
		return err
	}
	c.getBuffer().WriteString(";")
	return nil
}
func (c *CCompiler) VisitFnDeclarationStatement(t *parser.FnDeclarationStatement) error {
	c.InFunction = true
	cType, err := convertType(t.Type.Base)
	if err != nil {
		return err
	}
	c.getBuffer().WriteString(fmt.Sprintf("%s %s(", cType, t.Identifier))

	for i := 0; i < len(t.Params)-1; i++ {
		param := t.Params[i]
		ciType, err := convertType(param.Type.Base)
		if err != nil {
			return err
		}

		if param.Type.Array {
			c.getBuffer().WriteString(fmt.Sprintf("%s* %s,", ciType, param.Identifier))
		} else {
			c.getBuffer().WriteString(fmt.Sprintf("%s %s,", ciType, param.Identifier))
		}
	}
	param := t.Params[len(t.Params)-1]
	ciType, err := convertType(param.Type.Base)
	if err != nil {
		return err
	}
	if param.Type.Array {
		c.getBuffer().WriteString(fmt.Sprintf("%s* %s)", ciType, param.Identifier))
	} else {
		c.getBuffer().WriteString(fmt.Sprintf("%s %s)", ciType, param.Identifier))
	}
	err = t.Body.Accept(c)
	if err != nil {
		return err
	}
	c.InFunction = false
	return nil
}
func (c *CCompiler) VisitIfStatement(t *parser.IfStatement) error {
	c.getBuffer().WriteString("if(")
	err := t.Test.Accept(c)
	c.getBuffer().WriteString(")")
	if t.Consequent != nil {
		err = t.Consequent.Accept(c)
	}

	if t.Alternate != nil {
		c.getBuffer().WriteString("else")
		err = t.Alternate.Accept(c)
	}
	if err != nil {
		return err
	}
	return nil
}
func (c *CCompiler) VisitWhileStatement(t *parser.WhileStatement) error {
	c.InLoop = true
	c.getBuffer().WriteString("while(")
	err := t.Test.Accept(c)
	c.getBuffer().WriteString(")")
	err = t.Body.Accept(c)
	if err != nil {
		return err
	}
	c.InLoop = false
	return nil
}
func (c *CCompiler) VisitForStatement(t *parser.ForStatement) error {
	c.InLoop = true
	var err error
	c.getBuffer().WriteString("for(")
	if t.Init != nil {
		err = t.Init.Accept(c)
	}
	c.getBuffer().WriteString(";")
	if t.Test != nil {
		err = t.Test.Accept(c)
	}
	c.getBuffer().WriteString(";")
	if t.Inc != nil {
		err = t.Inc.Accept(c)
	}
	c.getBuffer().WriteString(")")
	err = t.Body.Accept(c)
	if err != nil {
		return err
	}
	c.InLoop = false
	return nil
}

func (c *CCompiler) VisitArrayLookup(t *parser.ArrayLookup) error {
	c.getBuffer().WriteString(fmt.Sprintf("%s[", t.Identifier))
	t.Index.Accept(c)
	c.getBuffer().WriteString("]")
	return nil
}

func (c *CCompiler) VisitArrayInitializer(t *parser.ArrayInitializer) error {
	return nil
}

func (c *CCompiler) VisitBreakStatement(t *parser.BreakStatement) error {
	if c.InLoop {
		c.getBuffer().WriteString("break;")
	} else {
		return errors.New("Cannot use break outside a loop")
	}
	return nil
}
func (c *CCompiler) VisitContinueStatement(t *parser.ContinueStatement) error {
	if c.InLoop {
		c.getBuffer().WriteString("continue;")
	} else {
		return errors.New("Cannot use continue outside a loop")
	}
	return nil
}
func (c *CCompiler) VisitReturnStatement(t *parser.ReturnStatement) error {
	if c.InFunction {
		c.getBuffer().WriteString("return ")
		err := t.Value.Accept(c)
		if err != nil {
			return err
		}
		c.getBuffer().WriteString(";")
	} else {
		return errors.New("Cannot use return outside a function")
	}
	return nil
}

func (c *CCompiler) VisitExitStatement(t *parser.ExitStatement) error {
	c.getBuffer().WriteString("exit(")
	t.Value.Accept(c)
	c.getBuffer().WriteString(");")
	return nil
}

func (c *CCompiler) VisitBlockStatement(t *parser.BlockStatement) error {
	c.getBuffer().WriteString("{")
	if t.Instructions == nil {
		c.getBuffer().WriteString("}")
		return nil
	}
	for _, stmt := range t.Instructions {
		err := stmt.Accept(c)
		if err != nil {
			return err
		}
	}
	c.getBuffer().WriteString("}")
	return nil
}
func (c *CCompiler) VisitAssignmentExpression(t *parser.AssignmentExpression) error {
	c.getBuffer().WriteString(fmt.Sprintf("%s%s", t.Identifier, t.Operator))
	err := t.Expression.Accept(c)
	if err != nil {
		return err
	}
	return nil
}
func (c *CCompiler) VisitVarDeclarationExpression(t *parser.VarDeclarationExpression) error {
	cType, _ := convertType(t.Type.Base)
	if t.ArrayInitializer != nil {
		c.getBuffer().WriteString(fmt.Sprintf("%s* %s=malloc(%d * sizeof(%s));", cType, t.Identifier, len(t.ArrayInitializer.Values), cType))
		for i, value := range t.ArrayInitializer.Values {
			c.getBuffer().WriteString(fmt.Sprintf("%s[%d]=", t.Identifier, i))
			switch v := (value).(type) {
			case parser.Bool:
				c.getBuffer().WriteString(fmt.Sprintf("%s", v.Value))
			case parser.Float:
				c.getBuffer().WriteString(fmt.Sprintf("%f", v.Value))
			case parser.Int:
				c.getBuffer().WriteString(fmt.Sprintf("%d", v.Value))
			case parser.String:
				c.getBuffer().WriteString(fmt.Sprintf("%s", v.Value))
			default:
				return errors.New("Error")
			}
			if i != len(t.ArrayInitializer.Values)-1 {
				c.getBuffer().WriteString(";")
			}

		}
	} else {
		cType, err := convertType(t.Type.Base)
		if err != nil {
			return err
		}
		c.getBuffer().WriteString(fmt.Sprintf("%s %s=", cType, t.Identifier))
		err = t.Expression.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *CCompiler) VisitExpression(t *parser.Expression) error {
	err := t.Value.Accept(c)
	if err != nil {
		return err
	}
	return nil
}
func (c *CCompiler) VisitLogicalOrExpression(t *parser.LogicalOrExpression) error {
	err := t.Left.Accept(c)
	if t.Right != nil {
		c.getBuffer().WriteString("||")
		err = t.Right.Accept(c)
	}
	return err
}
func (c *CCompiler) VisitLogicalAndExpression(t *parser.LogicalAndExpression) error {
	err := t.Left.Accept(c)
	if t.Right != nil {
		c.getBuffer().WriteString("&&")
		err = t.Right.Accept(c)
	}
	return err
}
func (c *CCompiler) VisitEqualityExpression(t *parser.EqualityExpression) error {
	err := t.Left.Accept(c)
	if t.Operator != nil {
		c.getBuffer().WriteString(fmt.Sprintf("%s", *t.Operator))
		err = t.Right.Accept(c)
	}
	return err
}
func (c *CCompiler) VisitRelationalExpression(t *parser.RelationalExpression) error {
	err := t.Left.Accept(c)
	if t.Operator != nil {
		c.getBuffer().WriteString(fmt.Sprintf("%s", *t.Operator))
		err = t.Right.Accept(c)
	}
	return err
}
func (c *CCompiler) VisitAdditiveExpression(t *parser.AdditiveExpression) error {
	err := t.Left.Accept(c)
	if t.Operator != nil {
		c.getBuffer().WriteString(fmt.Sprintf("%s", *t.Operator))
		err = t.Right.Accept(c)
	}

	return err
}
func (c *CCompiler) VisitMultiplicativeExpression(t *parser.MultiplicativeExpression) error {
	err := t.Left.Accept(c)
	if t.Operator != nil {
		c.getBuffer().WriteString(fmt.Sprintf("%s", *t.Operator))
		err = t.Right.Accept(c)
	}
	return err
}
func (c *CCompiler) VisitUnaryExpression(t *parser.UnaryExpression) error {
	return t.Value.Accept(c)
}
func (c *CCompiler) VisitPrimaryExpression(t *parser.PrimaryExpression) error {
	if t.Call != nil {
		t.Call.Accept(c)
	} else if t.Expression != nil {
		t.Accept(c)
	} else if t.Identifier != nil {
		c.getBuffer().WriteString(fmt.Sprintf("%s", *t.Identifier))
	} else if t.Literal != nil {
		switch v := (*t.Literal).(type) {
		case parser.Bool:
			c.getBuffer().WriteString(fmt.Sprintf("%s", v.Value))
		case parser.Float:
			c.getBuffer().WriteString(fmt.Sprintf("%f", v.Value))
		case parser.Int:
			c.getBuffer().WriteString(fmt.Sprintf("%d", v.Value))
		case parser.String:
			c.getBuffer().WriteString(fmt.Sprintf("%s", v.Value))
		default:
			return errors.New("Error")
		}
	} else if t.ArrayLookup != nil {
		t.ArrayLookup.Accept(c)
	}
	return nil
}
func (c *CCompiler) VisitCallExpression(t *parser.CallExpression) error {
	c.getBuffer().WriteString(fmt.Sprintf("%s(", t.FunctionName))
	for i := 0; i < len(t.Arguments)-1; i++ {
		t.Arguments[i].Accept(c)
		c.getBuffer().WriteString(",")
	}
	t.Arguments[len(t.Arguments)-1].Accept(c)
	c.getBuffer().WriteString(")")
	return nil
}
func (c *CCompiler) VisitImportStatement(t *parser.ImportStatement) error {
	return nil
}
func (c *CCompiler) VisitCallStatement(t *parser.CallStatement) error {
	return t.Expr.Accept(c)
}
func (c *CCompiler) VisitStatement(t *parser.Statement) error {
	return t.Accept(c)
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
