package ir

import (
	"fmt"
	"strings"
)

const IDENTATION = "  "

type IR interface {
	ToString() string
}

type BQCModule struct {
	Functions []IR
	Data      []IR
}

func (bqc BQCModule) ToString() string {
	builder := strings.Builder{}
	for _, ir := range bqc.Functions {
		builder.WriteString(ir.ToString())
	}
	return builder.String()
}

type BQCFunction struct {
	Export bool
	Name   string
	Type   BQCType
	Params []Paramter
	Body   []IR
}

func (bqc BQCFunction) ToString() string {
	builder := strings.Builder{}
	if bqc.Export {
		builder.WriteString("export ")
	}
	if bqc.Name == "main" {
		builder.WriteString(fmt.Sprintf("function %s $%s(", bqc.Type.ToString(), bqc.Name))
	} else {
		builder.WriteString(fmt.Sprintf("function %s $%s(", bqc.Type.ToString(), "bs_"+bqc.Name))
	}
	for i, param := range bqc.Params {
		builder.WriteString(param.ToString())
		if i < len(bqc.Params)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(") {\n")
	builder.WriteString("@start\n")
	for _, ir := range bqc.Body {
		builder.WriteString(ir.ToString())
	}
	builder.WriteString("}\n")
	return builder.String()
}

type VariableDeclaration struct {
	Name  string
	Type  BQCType
	Value string
}

func (vd VariableDeclaration) ToString() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s%s =l alloc4 %d\n", IDENTATION, "%p_"+vd.Name, vd.Type.Bytes))
	builder.WriteString(fmt.Sprintf("%sstore%s %s, %s\n", IDENTATION, vd.Type.ToString(), "%"+vd.Value, "%p_"+vd.Name))
	return builder.String()
}

type AssigmentDeclaration struct {
	Name  string
	Type  BQCType
	Value string
}

func (ad AssigmentDeclaration) ToString() string {
	return fmt.Sprintf("%sstore%s %s, %s\n", IDENTATION, ad.Type.ToString(), "%"+ad.Value, "%p_"+ad.Name)
}

type ModifyAssigmentDeclaration struct {
	Name     string
	Type     BQCType
	Value    string
	Tmp      string
	Operator string
}

func (iad ModifyAssigmentDeclaration) ToString() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s%s =%s load%s %s\n", IDENTATION, "%"+iad.Tmp, iad.Type.ToString(), iad.Type.ToString(), "%p_"+iad.Name))
	builder.WriteString(fmt.Sprintf("%s%s =%s %s %s, %s\n", IDENTATION, "%"+iad.Tmp, iad.Type.ToString(), iad.Operator, "%"+iad.Tmp, "%"+iad.Value))
	builder.WriteString(fmt.Sprintf("%sstore%s %s, %s\n", IDENTATION, iad.Type.ToString(), "%"+iad.Tmp, "%p_"+iad.Name))
	return builder.String()
}

type BQCIf struct {
	Count     int
	Condition []IR
	TmpName   string
	When      []IR
	Then      []IR
}

func (bqc BQCIf) ToString() string {
	builder := strings.Builder{}
	wlabel := fmt.Sprintf("@w_%d", bqc.Count)
	tlabel := fmt.Sprintf("@t_%d", bqc.Count)
	elabel := fmt.Sprintf("@e_%d", bqc.Count)
	for _, ir := range bqc.Condition {
		builder.WriteString(ir.ToString())
	}
	if len(bqc.Then) > 0 {
		builder.WriteString(fmt.Sprintf("%sjnz %s, %s, %s\n", IDENTATION, "%"+bqc.TmpName, wlabel, tlabel))
	} else {
		builder.WriteString(fmt.Sprintf("%sjnz %s, %s, %s\n", IDENTATION, "%"+bqc.TmpName, wlabel, elabel))
	}
	builder.WriteString(wlabel + "\n")
	for _, ir := range bqc.When {
		builder.WriteString(ir.ToString())
	}
	if len(bqc.Then) > 0 {
		builder.WriteString(fmt.Sprintf("%sjmp %s\n", IDENTATION, elabel))
		builder.WriteString(tlabel + "\n")
		for _, ir := range bqc.Then {
			builder.WriteString(ir.ToString())
		}
	}
	builder.WriteString(elabel + "\n")
	return builder.String()
}

type BQCWhile struct {
	Count     int
	Condition []IR
	TmpName   string
	Body      []IR
}

func (bqc BQCWhile) ToString() string {
	builder := strings.Builder{}
	wlabel := fmt.Sprintf("@w_%d", bqc.Count)
	slabel := fmt.Sprintf("@s_%d", bqc.Count)
	elabel := fmt.Sprintf("@e_%d", bqc.Count)
	builder.WriteString(wlabel + "\n")
	for _, ir := range bqc.Condition {
		builder.WriteString(ir.ToString())
	}
	builder.WriteString(fmt.Sprintf("%sjnz %s, %s, %s\n", IDENTATION, "%"+bqc.TmpName, slabel, elabel))
	builder.WriteString(slabel + "\n")
	for _, ir := range bqc.Body {
		builder.WriteString(ir.ToString())
	}
	builder.WriteString(fmt.Sprintf("%sjmp %s\n", IDENTATION, wlabel))
	builder.WriteString(elabel + "\n")
	return builder.String()
}

type BQCFor struct {
	Count     int
	Condition []IR
	Inc       []IR
	TmpName   string
	Body      []IR
}

func (bqc BQCFor) ToString() string {
	builder := strings.Builder{}
	flabel := fmt.Sprintf("@f_%d", bqc.Count)
	blabel := fmt.Sprintf("@b_%d", bqc.Count)
	clabel := fmt.Sprintf("@c_%d", bqc.Count)
	elabel := fmt.Sprintf("@e_%d", bqc.Count)
	builder.WriteString(flabel + "\n")
	for _, ir := range bqc.Condition {
		builder.WriteString(ir.ToString())
	}
	builder.WriteString(fmt.Sprintf("%sjnz %s, %s, %s\n", IDENTATION, "%"+bqc.TmpName, blabel, elabel))
	builder.WriteString(blabel + "\n")
	for _, ir := range bqc.Body {
		builder.WriteString(ir.ToString())
	}
	builder.WriteString(clabel + "\n")
	for _, ir := range bqc.Inc {
		builder.WriteString(ir.ToString())
	}
	builder.WriteString(fmt.Sprintf("%sjmp %s\n", IDENTATION, flabel))
	builder.WriteString(elabel + "\n")
	return builder.String()
}

type BQCFunctionCall struct {
	Name     string
	IsExtern bool
	Args     []Paramter
	NonVoid  bool
	Tmp      string
	Type     BQCType
}

func (bqc BQCFunctionCall) ToString() string {
	builder := strings.Builder{}
	if bqc.NonVoid {
		builder.WriteString(fmt.Sprintf("%s%s =%s ", IDENTATION, "%"+bqc.Tmp, bqc.Type.ToString()))
	} else {
		builder.WriteString(IDENTATION)
	}
	if bqc.IsExtern || bqc.Name == "main" {
		builder.WriteString(fmt.Sprintf("call $%s(", bqc.Name))
	} else {
		builder.WriteString(fmt.Sprintf("call $%s(", "bs_"+bqc.Name))
	}
	for i, arg := range bqc.Args {
		builder.WriteString(arg.ToString())
		if i < len(bqc.Args)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")\n")
	return builder.String()
}

func (bqc BQCFunctionCall) GetTmp() string {
	return bqc.Tmp
}

type BQCReturn struct {
	Tmp      string
	HasValue bool
}

func (bqc BQCReturn) ToString() string {
	if !bqc.HasValue {
		return fmt.Sprintf("%sret\n", IDENTATION)
	}
	return fmt.Sprintf("%sret %s\n", IDENTATION, "%"+bqc.Tmp)
}

type BQCType struct {
	Rep   string
	Bytes int
}

func (bqc BQCType) ToString() string {
	return bqc.Rep
}

type BQCLiteral struct {
	Value string
	Type  BQCType
	Tmp   string
}

func (bqc BQCLiteral) ToString() string {
	return fmt.Sprintf("%s%s =%s copy %s\n", IDENTATION, "%"+bqc.Tmp, bqc.Type.ToString(), bqc.Value)
}

func (bqc BQCLiteral) GetTmp() string {
	return bqc.Tmp
}

type BQCVariable struct {
	Name string
	Type BQCType
	Tmp  string
}

func (bqc BQCVariable) ToString() string {
	return fmt.Sprintf("%s%s =%s load%s %s\n", IDENTATION, "%"+bqc.Tmp, bqc.Type.ToString(), bqc.Type.ToString(), "%p_"+bqc.Name)
}

func (bqc BQCVariable) GetTmp() string {
	return bqc.Tmp
}

type BQCBinary struct {
	Op   string
	TmpL string
	TmpR string
	Tmp  string
	Type BQCType
}

func (bqc BQCBinary) ToString() string {
	return fmt.Sprintf("%s%s =%s %s %s, %s\n", IDENTATION, "%"+bqc.Tmp, bqc.Type.ToString(), bqc.Op, "%"+bqc.TmpL, "%"+bqc.TmpR)
}

func (bqc BQCBinary) GetTmp() string {
	return bqc.Tmp
}

type IRTmp interface {
	IR
	GetTmp() string
}

type Paramter struct {
	Name string
	Type BQCType
}

func (p Paramter) ToString() string {
	return fmt.Sprintf("%s %s", p.Type.ToString(), "%"+p.Name)
}

type BQCJump struct {
	Label string
}

func (bqc BQCJump) ToString() string {
	return fmt.Sprintf("%sjmp %s\n", IDENTATION, bqc.Label)
}

type StringData struct {
	Name  string
	Value string
}

func (sd StringData) ToString() string {
	return fmt.Sprintf(`data $%s = {b "%s" ,b 0}\n`, "%p_"+sd.Name, sd.Value)
}
