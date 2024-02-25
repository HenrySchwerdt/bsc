package compiler

import (
	"bsc/src/ir"
	"bsc/src/parser"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BQCCompiler struct {
	Ast           parser.Program
	StdLibPath    string
	Imports       []ir.BQCModule
	Module        ir.BQCModule
	CurrentFrame  [][]ir.IR
	Vars          string
	ContinueLabel string
	BreakLabel    string
	lastUsedVar   rune
	InLoop        bool
	Counter       int
	InFunction    bool
}

func (c *BQCCompiler) getCurrentFrame() []ir.IR {
	return c.CurrentFrame[len(c.CurrentFrame)-1]
}

func (c *BQCCompiler) addIR(ir ir.IR) {

	c.CurrentFrame[len(c.CurrentFrame)-1] = append(c.CurrentFrame[len(c.CurrentFrame)-1], ir)
}

func (c *BQCCompiler) pushFrame() {
	c.CurrentFrame = append(c.CurrentFrame, make([]ir.IR, 0))
}

func (c *BQCCompiler) popFrame() []ir.IR {
	frame := c.CurrentFrame[len(c.CurrentFrame)-1]
	c.CurrentFrame = c.CurrentFrame[:len(c.CurrentFrame)-1]
	return frame
}

func NewBQCCompiler(ast *parser.Program) *BQCCompiler {
	return &BQCCompiler{
		Ast:          *ast,
		Imports:      make([]ir.BQCModule, 0),
		Module:       ir.BQCModule{},
		CurrentFrame: make([][]ir.IR, 0),
		lastUsedVar:  'a' - 1,
		InLoop:       false,
		Counter:      0,
		InFunction:   false,
	}
}

func (c *BQCCompiler) getNextVariableName() string {
	if c.lastUsedVar < 'z' {
		c.lastUsedVar++
		return string(c.lastUsedVar)
	}
	c.lastUsedVar = 'a'
	return string(c.lastUsedVar)
}

func (c *BQCCompiler) getCounter() int {
	c.Counter++
	return c.Counter
}

func (c *BQCCompiler) VisitProgram(t *parser.Program) error {
	module := ir.BQCModule{
		Functions: make([]ir.IR, 0),
	}
	c.Module = module
	for _, stmt := range t.Statements {
		err := stmt.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *BQCCompiler) VisitVarDeclarationStatement(t *parser.VarDeclarationStatement) error {
	t.Dec.Accept(c)
	return nil
}
func (c *BQCCompiler) VisitAssignmentStatement(t *parser.AssignmentStatement) error {
	t.Ass.Accept(c)
	return nil
}
func (c *BQCCompiler) VisitFnDeclarationStatement(t *parser.FnDeclarationStatement) error {
	fmt.Println("In function ", t.Identifier)
	c.pushFrame()
	c.InFunction = true
	params := make([]ir.Paramter, 0)
	for _, p := range t.Params {
		c.addIR(ir.VariableDeclaration{
			Name:  p.Identifier,
			Type:  ir.BQCType{Rep: "w", Bytes: 4}, // TODO fix this
			Value: p.Identifier,
		})
		params = append(params, ir.Paramter{
			Name: p.Identifier,
			Type: ir.BQCType{Rep: "w", Bytes: 4}, // TODO fix this
		})
	}
	for _, i := range t.Body.Instructions {
		i.Accept(c)
	}
	c.InFunction = false
	frame := c.popFrame()
	var ftype ir.BQCType
	if t.Type.Base == "void" {
		ftype = ir.BQCType{Rep: "", Bytes: 0}
	} else {
		ftype = ir.BQCType{Rep: "w", Bytes: 4}
	}
	frame = append(frame, ir.BQCReturn{
		HasValue: false,
	})
	c.Module.Functions = append(c.Module.Functions, ir.BQCFunction{
		Export: t.Export,
		Name:   t.Identifier,
		Params: params,
		Type:   ftype,
		Body:   frame,
	})
	return nil
}
func (c *BQCCompiler) VisitIfStatement(t *parser.IfStatement) error {
	countstmt := c.getCounter()
	c.pushFrame()
	t.Test.Accept(c)
	testFrame := c.popFrame()
	c.pushFrame()
	t.Consequent.Accept(c)
	whenFrame := c.popFrame()
	c.pushFrame()
	if t.Alternate != nil {
		t.Alternate.Accept(c)
	}
	thenFrame := c.popFrame()
	c.addIR(ir.BQCIf{
		Count:     countstmt,
		Condition: testFrame,
		TmpName:   testFrame[len(testFrame)-1].(ir.IRTmp).GetTmp(),
		When:      whenFrame,
		Then:      thenFrame,
	})
	return nil
}
func (c *BQCCompiler) VisitWhileStatement(t *parser.WhileStatement) error {
	c.InLoop = true
	countstmt := c.getCounter()
	tmpBreak := c.BreakLabel
	tmpContinue := c.ContinueLabel
	c.pushFrame()
	t.Test.Accept(c)
	testFrame := c.popFrame()
	c.BreakLabel = fmt.Sprintf("@e_%d", countstmt)
	c.ContinueLabel = fmt.Sprintf("@w_%d", countstmt)
	c.pushFrame()
	t.Body.Accept(c)
	bodyFrame := c.popFrame()
	c.BreakLabel = tmpBreak
	c.ContinueLabel = tmpContinue
	c.addIR(ir.BQCWhile{
		Count:     countstmt,
		Condition: testFrame,
		TmpName:   testFrame[len(testFrame)-1].(ir.IRTmp).GetTmp(),
		Body:      bodyFrame,
	})
	c.InLoop = false
	return nil
}
func (c *BQCCompiler) VisitForStatement(t *parser.ForStatement) error {
	c.InLoop = true
	countstmt := c.getCounter()
	t.Init.Accept(c)
	tmpBreak := c.BreakLabel
	tmpContinue := c.ContinueLabel
	c.pushFrame()
	t.Test.Accept(c)
	testFrame := c.popFrame()
	c.BreakLabel = fmt.Sprintf("@e_%d", countstmt)
	c.ContinueLabel = fmt.Sprintf("@c_%d", countstmt)
	c.pushFrame()
	t.Body.Accept(c)
	bodyFrame := c.popFrame()
	c.pushFrame()
	t.Inc.Accept(c)
	incFrame := c.popFrame()
	c.BreakLabel = tmpBreak
	c.ContinueLabel = tmpContinue
	c.addIR(ir.BQCFor{
		Count:     countstmt,
		Condition: testFrame,
		Inc:       incFrame,
		Body:      bodyFrame,
		TmpName:   testFrame[len(testFrame)-1].(ir.IRTmp).GetTmp(),
	})
	c.InLoop = false
	return nil
}

func (c *BQCCompiler) VisitArrayLookup(t *parser.ArrayLookup) error {
	// c.getBuffer().WriteString(fmt.Sprintf("%s[", t.Identifier))
	// t.Index.Accept(c)
	// c.getBuffer().WriteString("]")
	return nil
}

func (c *BQCCompiler) VisitArrayInitializer(t *parser.ArrayInitializer) error {
	return nil
}

func (c *BQCCompiler) VisitBreakStatement(t *parser.BreakStatement) error {
	if c.InLoop {
		c.addIR(ir.BQCJump{
			Label: c.BreakLabel,
		})
	} else {
		return errors.New("cannot use break outside a loop")
	}
	return nil
}
func (c *BQCCompiler) VisitContinueStatement(t *parser.ContinueStatement) error {
	if c.InLoop {
		c.addIR(ir.BQCJump{
			Label: c.ContinueLabel,
		})
	} else {
		return errors.New("cannot use continue outside a loop")
	}
	return nil
}
func (c *BQCCompiler) VisitReturnStatement(t *parser.ReturnStatement) error {
	t.Value.Accept(c)
	lastTmp := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp).GetTmp()
	c.addIR(ir.BQCReturn{
		Tmp:      lastTmp,
		HasValue: true, // TODO: fix this
	})
	return nil
}

func (c *BQCCompiler) VisitExitStatement(t *parser.ExitStatement) error {
	t.Value.Accept(c)
	lastTmp := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp).GetTmp()
	c.addIR(ir.BQCFunctionCall{
		Name: "exit",
		Args: []ir.Paramter{
			{
				Type: ir.BQCType{Rep: "w", Bytes: 4},
				Name: lastTmp,
			},
		},
	})
	return nil
}

func (c *BQCCompiler) VisitPrintStatement(t *parser.PrintStatement) error {
	t.Value.Accept(c)
	lastTmp := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp).GetTmp()
	c.addIR(ir.BQCFunctionCall{
		Name: "print",
		Args: []ir.Paramter{
			{
				Type: ir.BQCType{Rep: "w", Bytes: 4},
				Name: lastTmp,
			},
		},
	})
	return nil
}

func (c *BQCCompiler) VisitBlockStatement(t *parser.BlockStatement) error {
	for _, stmt := range t.Instructions {
		err := stmt.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *BQCCompiler) VisitAssignmentExpression(t *parser.AssignmentExpression) error {
	t.Expression.Accept(c)
	lastTmp := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp).GetTmp()
	c.addIR(ir.AssigmentDeclaration{
		Name:  t.Identifier,
		Type:  ir.BQCType{Rep: "w", Bytes: 4}, // TODO: fix this
		Value: lastTmp,
	})
	return nil
}
func (c *BQCCompiler) VisitVarDeclarationExpression(t *parser.VarDeclarationExpression) error {
	t.Expression.Accept(c)
	lastTmp := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp).GetTmp()
	c.addIR(ir.VariableDeclaration{
		Name:  t.Identifier,
		Type:  ir.BQCType{Rep: "w", Bytes: 4}, // TODO: fix this
		Value: lastTmp,
	})
	return nil
}
func (c *BQCCompiler) VisitExpression(t *parser.Expression) error {
	err := t.Value.Accept(c)
	if err != nil {
		return err
	}
	return nil
}
func (c *BQCCompiler) VisitLogicalOrExpression(t *parser.LogicalOrExpression) error {
	t.Left.Accept(c)
	if t.Right != nil {
		irL := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		t.Right.Accept(c)
		irR := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		c.addIR(ir.BQCBinary{
			Op:   "or",
			TmpL: irL.GetTmp(),
			TmpR: irR.GetTmp(),
			Tmp:  c.getNextVariableName(),
			Type: ir.BQCType{Rep: "w", Bytes: 4},
		})
	}
	return nil
}
func (c *BQCCompiler) VisitLogicalAndExpression(t *parser.LogicalAndExpression) error {
	t.Left.Accept(c)
	if t.Right != nil {
		irL := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		t.Right.Accept(c)
		irR := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		c.addIR(ir.BQCBinary{
			Op:   "and",
			TmpL: irL.GetTmp(),
			TmpR: irR.GetTmp(),
			Tmp:  c.getNextVariableName(),
			Type: ir.BQCType{Rep: "w", Bytes: 4},
		})
	}
	return nil
}
func (c *BQCCompiler) VisitEqualityExpression(t *parser.EqualityExpression) error {
	t.Left.Accept(c)
	if t.Right != nil {
		irL := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		t.Right.Accept(c)
		irR := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		tmp := c.getNextVariableName()
		var op string
		if *t.Operator == "==" {
			op = "ceqw"
		} else {
			op = "cnew"
		}
		c.addIR(ir.BQCBinary{
			Op:   op,
			TmpL: irL.GetTmp(),
			TmpR: irR.GetTmp(),
			Tmp:  tmp,
			Type: ir.BQCType{Rep: "w", Bytes: 4},
		})
	}
	return nil
}
func (c *BQCCompiler) VisitRelationalExpression(t *parser.RelationalExpression) error {
	t.Left.Accept(c)
	if t.Right != nil {
		irL := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		t.Right.Accept(c)
		irR := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		tmp := c.getNextVariableName()
		var op string
		if *t.Operator == "<" {
			op = "csltw"
		} else if *t.Operator == ">" {
			op = "csgtw"
		} else if *t.Operator == "<=" {
			op = "cslew"
		} else {
			op = "csgew"
		}
		c.addIR(ir.BQCBinary{
			Op:   op,
			TmpL: irL.GetTmp(),
			TmpR: irR.GetTmp(),
			Tmp:  tmp,
			Type: ir.BQCType{Rep: "w", Bytes: 4},
		})
	}
	return nil
}
func (c *BQCCompiler) VisitAdditiveExpression(t *parser.AdditiveExpression) error {
	t.Left.Accept(c)
	if t.Right != nil {
		irL := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		t.Right.Accept(c)
		irR := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		tmp := c.getNextVariableName()
		var op string
		if *t.Operator == "+" {
			op = "add"
		} else {
			op = "sub"
		}
		c.addIR(ir.BQCBinary{
			Op:   op,
			TmpL: irL.GetTmp(),
			TmpR: irR.GetTmp(),
			Tmp:  tmp,
			Type: ir.BQCType{Rep: "w", Bytes: 4},
		})
	}
	return nil
}
func (c *BQCCompiler) VisitMultiplicativeExpression(t *parser.MultiplicativeExpression) error {
	t.Left.Accept(c)
	if t.Right != nil {
		irL := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		t.Right.Accept(c)
		irR := c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp)
		tmp := c.getNextVariableName()
		var op string
		if *t.Operator == "*" {
			op = "mul"
		} else if *t.Operator == "/" {
			op = "div"
		} else {
			op = "rem"
		}
		c.addIR(ir.BQCBinary{
			Op:   op,
			TmpL: irL.GetTmp(),
			TmpR: irR.GetTmp(),
			Tmp:  tmp,
			Type: ir.BQCType{Rep: "w", Bytes: 4},
		})
	}
	return nil
}
func (c *BQCCompiler) VisitUnaryExpression(t *parser.UnaryExpression) error {
	return t.Value.Accept(c)
}
func (c *BQCCompiler) VisitPrimaryExpression(t *parser.PrimaryExpression) error {
	if t.Call != nil {
		return t.Call.Accept(c)
	} else if t.Expression != nil {
		return t.Accept(c)
	} else if t.Identifier != nil {
		c.addIR(ir.BQCVariable{
			Name: *t.Identifier,
			Type: ir.BQCType{Rep: "w", Bytes: 4},
			Tmp:  c.getNextVariableName(),
		})
		return nil
	} else if t.Literal != nil {
		switch v := (*t.Literal).(type) {
		case parser.Bool:
			if v.Value == "true" {
				c.addIR(ir.BQCLiteral{
					Value: "1",
					Type:  ir.BQCType{Rep: "w", Bytes: 4},
					Tmp:   c.getNextVariableName(),
				})
			} else {
				c.addIR(ir.BQCLiteral{
					Value: "0",
					Type:  ir.BQCType{Rep: "w", Bytes: 4},
					Tmp:   c.getNextVariableName(),
				})
			}
		case parser.Float:
			c.addIR(ir.BQCLiteral{
				Value: fmt.Sprintf("%d", int(v.Value)), // TODO: fix this
				Type:  ir.BQCType{Rep: "w", Bytes: 4},
				Tmp:   c.getNextVariableName(),
			})
		case parser.Int:
			c.addIR(ir.BQCLiteral{
				Value: fmt.Sprintf("%d", v.Value),
				Type:  ir.BQCType{Rep: "w", Bytes: 4},
				Tmp:   c.getNextVariableName(),
			})
		case parser.String:
			c.addIR(ir.BQCLiteral{
				Value: v.Value, // TODO: fix this
				Type:  ir.BQCType{Rep: "w", Bytes: 4},
				Tmp:   c.getNextVariableName(),
			})
		default:
			return errors.New("Error")
		}
	} else if t.ArrayLookup != nil {
		t.ArrayLookup.Accept(c)
	}
	return nil
}
func (c *BQCCompiler) VisitCallExpression(t *parser.CallExpression) error {
	tmps := make([]ir.Paramter, 0)
	for _, arg := range t.Arguments {
		arg.Accept(c)
		tmps = append(tmps, ir.Paramter{
			Name: c.getCurrentFrame()[len(c.getCurrentFrame())-1].(ir.IRTmp).GetTmp(),
			Type: ir.BQCType{Rep: "w", Bytes: 4},
		})
	}
	// TODO need to make checkup
	c.addIR(ir.BQCFunctionCall{
		Name:    t.FunctionName,
		Args:    tmps,
		NonVoid: true,
		Tmp:     c.getNextVariableName(),
		Type:    ir.BQCType{Rep: "w", Bytes: 4},
	})
	return nil
}
func (c *BQCCompiler) VisitImportStatement(t *parser.ImportStatement) error {
	file, err := os.Open(filepath.Join(c.StdLibPath, strings.Trim(t.File, "\"")+".bs"))
	if err != nil {
		return err
	}
	defer file.Close()
	parser := parser.NewNParser()
	ast, err := parser.Parse(file.Name(), file)
	if err != nil {
		return err
	}
	compiler := NewBQCCompiler(ast)
	err = compiler.compile()
	if err != nil {
		return err
	}
	c.Imports = append(c.Imports, compiler.Imports...)
	c.Imports = append(c.Imports, compiler.Module)
	return nil
}
func (c *BQCCompiler) VisitCallStatement(t *parser.CallStatement) error {
	return t.Expr.Accept(c)
}
func (c *BQCCompiler) VisitStatement(t *parser.Statement) error {
	if t.AssignmentStatement != nil {
		return t.AssignmentStatement.Accept(c)
	} else if t.VarDeclarationStatement != nil {
		return t.VarDeclarationStatement.Accept(c)
	} else if t.FnDeclarationStatement != nil {
		return t.FnDeclarationStatement.Accept(c)
	} else if t.IfStatement != nil {
		return t.IfStatement.Accept(c)
	} else if t.WhileStatement != nil {
		return t.WhileStatement.Accept(c)
	} else if t.ForStatement != nil {
		return t.ForStatement.Accept(c)
	} else if t.ReturnStatement != nil {
		return t.ReturnStatement.Accept(c)
	} else if t.BreakStatement != nil {
		return t.BreakStatement.Accept(c)
	} else if t.ContinueStatement != nil {
		return t.ContinueStatement.Accept(c)
	} else if t.BlockStatement != nil {
		return t.BlockStatement.Accept(c)
	} else if t.CallStatement != nil {
		return t.CallStatement.Accept(c)
	} else if t.ExitStatement != nil {
		return t.ExitStatement.Accept(c)
	}
	return nil
}

func (c *BQCCompiler) compile() error {
	if err := (c.Ast).Accept(c); err != nil {
		return err
	}
	return nil
}

func (c *BQCCompiler) Compile(outDir, outFile string) error {
	// Compile the program and check for errors
	err := c.compile()
	if err != nil {
		return err
	}

	// Remove existing output directory and create a new one
	os.RemoveAll(outDir)
	os.Mkdir(outDir, 0755)

	// Check if QBE is available
	_, err = exec.LookPath("qbe")
	if err != nil {
		fmt.Println("QBE not found in the system.")
		return err
	}

	var concatenatedModules bytes.Buffer
	for _, importedModule := range c.Imports {
		concatenatedModules.WriteString(importedModule.ToString())
	}
	concatenatedModules.WriteString(c.Module.ToString())

	tmpFile := filepath.Join(outDir, outFile+".ssa")
	err = os.WriteFile(tmpFile, concatenatedModules.Bytes(), 0644)
	if err != nil {
		return err
	}

	cmdQBE := exec.Command("qbe", "-o", filepath.Join(outDir, outFile+".s"), tmpFile)
	qbeOutput, err := cmdQBE.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running qbe: %v\n%s", err, string(qbeOutput))
	}

	cmdCC := exec.Command("cc", "-o", filepath.Join(outDir, outFile), filepath.Join(outDir, outFile+".s"))
	ccOutput, err := cmdCC.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running cc: %v\n%s", err, string(ccOutput))
	}

	return nil
}
