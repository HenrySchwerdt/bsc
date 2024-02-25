package compiler

import (
	"bsc/src/ir"
	"bsc/src/parser"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type BQCCompiler struct {
	Ast            parser.Program
	OutMain        strings.Builder
	OutOuter       strings.Builder
	Module         ir.BQCModule
	CurrentFrame   [][]ir.IR
	LoopCount      int
	ConditionCount int
	Vars           string
	ContinueLabel  string
	BreakLabel     string
	lastUsedVar    rune
	InLoop         bool
	Counter        int
	InFunction     bool
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
	var builder strings.Builder
	var builderfn strings.Builder
	return &BQCCompiler{
		Ast:            *ast,
		OutMain:        builder,
		OutOuter:       builderfn,
		Module:         ir.BQCModule{},
		CurrentFrame:   make([][]ir.IR, 0),
		LoopCount:      0,
		ConditionCount: 0,
		lastUsedVar:    'a' - 1,
		InLoop:         false,
		Counter:        0,
		InFunction:     false,
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
	c.pushFrame()
	for _, stmt := range t.Statements {
		err := stmt.Accept(c)
		if err != nil {
			return err
		}
	}
	tmpName := c.getNextVariableName()
	c.addIR(ir.BQCLiteral{
		Value: "0",
		Type:  ir.BQCType{Rep: "w", Bytes: 4},
		Tmp:   tmpName,
	})
	c.addIR(ir.BQCReturn{
		Tmp: tmpName,
	})
	frame := c.popFrame()
	c.Module.Functions = append(c.Module.Functions, ir.BQCFunction{
		Export: true,
		Name:   "main",
		Type:   ir.BQCType{Rep: "w", Bytes: 4},
		Params: make([]ir.Paramter, 0),
		Body:   frame,
	})

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
	c.Module.Functions = append(c.Module.Functions, ir.BQCFunction{
		Export: false,
		Name:   t.Identifier,
		Params: params,
		Type:   ir.BQCType{Rep: "w", Bytes: 4}, // TODO fix this
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
		Tmp: lastTmp,
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

func (c *BQCCompiler) Compile(outDir, outFile string) error {
	if err := (c.Ast).Accept(c); err != nil {
		return err
	}
	os.RemoveAll(outDir)
	os.Mkdir(outDir, 0755)
	_, err := exec.LookPath("qbe")
	if err != nil {
		fmt.Println("QBE not found in the system.")
		os.Exit(1)
	}
	os.WriteFile(fmt.Sprintf("%s/%s.ssa", outDir, outFile), []byte(c.Module.ToString()), 0755)
	cmd := exec.Command("qbe", "-o", fmt.Sprintf("%s/%s.s", outDir, outFile), fmt.Sprintf("%s/%s.ssa", outDir, outFile))
	cmd.Run()
	cmd = exec.Command("cc", "-o", fmt.Sprintf("%s/%s", outDir, outFile), fmt.Sprintf("%s/%s.s", outDir, outFile))
	cmd.Run()

	return nil
}
