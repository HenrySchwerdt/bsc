package parser

import (
	"bsc/src/lexer"
)

type Node interface {
	Accept(Visitor) error
}

type BaseNode struct {
	Start lexer.Token
	End   lexer.Token
}

type Literal struct {
	BaseNode
	Value interface{}
}

func (l *Literal) Accept(v Visitor) error {
	return v.VisitLiteral(l)
}

type Identifier struct {
	BaseNode
	Name string `json:"name"`
}

func (id *Identifier) Accept(v Visitor) error {
	return v.VisitIdentifier(id)
}

type UnaryExpression struct {
	BaseNode
	Argument Node   `json:"arg"`
	Operator string `json:"op"`
}

func (ue *UnaryExpression) Accept(v Visitor) error {
	return v.VisitUnaryExpression(ue)
}

type Params struct {
	BaseNode
	Args []string
}

func (pa *Params) Accept(v Visitor) error {
	return v.VisitParams(pa)
}

type BinaryExpression struct {
	BaseNode
	Left     Node   `json:"left"`
	Right    Node   `json:"right"`
	Operator string `json:"op"`
}

func (be *BinaryExpression) Accept(v Visitor) error {
	return v.VisitBinaryExpression(be)
}

type ConditionalExpression struct {
	BaseNode
	Test       Node `json:"test"`
	Consequent Node `json:"cons"`
	Alternate  Node `json:"alt"`
}

func (ce *ConditionalExpression) Accept(v Visitor) error {
	return v.VisitConditionalExpression(ce)
}

type CallExpression struct {
	BaseNode
	Identifier string
	Args       []Node
}

func (ce *CallExpression) Accept(v Visitor) error {
	return v.VisitCallExpression(ce)
}

type ExpressionStatement struct {
	BaseNode
	Expression Node `json:"expression"`
}

func (es *ExpressionStatement) Accept(v Visitor) error {
	return v.VisitExpressionStatement(es)
}

type VariableDeclarator struct {
	BaseNode
	Id   string `json:"id"`
	Init Node   `json:"init"`
}

func (vd *VariableDeclarator) Accept(v Visitor) error {
	return v.VisitVariableDeclarator(vd)
}

type VariableDeclaration struct {
	BaseNode
	Kind        string `json:"kind"`
	Declaration Node   `json:"declaration"`
}

func (vd *VariableDeclaration) Accept(v Visitor) error {
	return v.VisitVariableDeclaration(vd)
}

type VariableLookup struct {
	BaseNode
	Id string `json:"id"`
}

func (vl *VariableLookup) Accept(v Visitor) error {
	return v.VisitVariableLookup(vl)
}

type Function struct {
	BaseNode
	Id         string `json:"id"`
	Parameters Node   `json:"params"`
	ReturnType string `json:"ret"`
	Body       Node   `json:"body"`
}

func (fn *Function) Accept(v Visitor) error {
	return v.VisitFunction(fn)
}

type FunctionDeclaration struct {
	BaseNode
	Function Node `json:"function"`
}

func (fd *FunctionDeclaration) Accept(v Visitor) error {
	return v.VisitFunctionDeclaration(fd)
}

type ForStatment struct {
	BaseNode
	Init   Node `json:"init"`
	Test   Node `json:"test"`
	Update Node `json:"update"`
	Body   Node `json:"body"`
}

func (fs *ForStatment) Accept(v Visitor) error {
	return v.VisitForStatment(fs)
}

type WhileStatment struct {
	BaseNode
	Test Node `json:"test"`
	Body Node `json:"body"`
}

func (ws *WhileStatment) Accept(v Visitor) error {
	return v.VisitWhileStatment(ws)
}

type IfStatment struct {
	BaseNode
	Test       Node `json:"test"`
	Consequent Node `json:"cons"`
	Alternate  Node `json:"alt"`
}

func (is *IfStatment) Accept(v Visitor) error {
	return v.VisitIfStatment(is)
}

type AssignmentStatement struct {
	BaseNode
	Identifier string `json:"ident"`
	Value      Node   `json:"value"`
}

func (as *AssignmentStatement) Accept(v Visitor) error {
	return v.VisitAssignmentStatement(as)
}

type ReturnStatment struct {
	BaseNode
	Argument Node `json:"arg"`
}

func (rs *ReturnStatment) Accept(v Visitor) error {
	return v.VisitReturnStatment(rs)
}

type ExitStatment struct {
	BaseNode
	Argument Node `json:"arg"`
}

func (es *ExitStatment) Accept(v Visitor) error {
	return v.VisitExitStatment(es)
}

type BreakStatment struct {
	BaseNode
}

func (bs *BreakStatment) Accept(v Visitor) error {
	return v.VisitBreakStatment(bs)
}

type BlockStatement struct {
	BaseNode
	Instructions []Node `json:"instructions"`
}

func (bs *BlockStatement) Accept(v Visitor) error {
	return v.VisitBlockStatement(bs)
}

type Program struct {
	BaseNode
	Instructions []Node `json:"instructions"`
}

func (p *Program) Accept(v Visitor) error {
	return v.VisitProgram(p)
}

type Visitor interface {
	VisitLiteral(l *Literal) error
	VisitIdentifier(id *Identifier) error
	VisitUnaryExpression(ue *UnaryExpression) error
	VisitBinaryExpression(be *BinaryExpression) error
	VisitConditionalExpression(ce *ConditionalExpression) error
	VisitExpressionStatement(es *ExpressionStatement) error
	VisitVariableDeclarator(vd *VariableDeclarator) error
	VisitVariableDeclaration(vd *VariableDeclaration) error
	VisitVariableLookup(vl *VariableLookup) error
	VisitFunction(fn *Function) error
	VisitFunctionDeclaration(fd *FunctionDeclaration) error
	VisitForStatment(fs *ForStatment) error
	VisitIfStatment(is *IfStatment) error
	VisitReturnStatment(rs *ReturnStatment) error
	VisitWhileStatment(ws *WhileStatment) error
	VisitBlockStatement(bs *BlockStatement) error
	VisitExitStatment(es *ExitStatment) error
	VisitBreakStatment(bs *BreakStatment) error
	VisitAssignmentStatement(as *AssignmentStatement) error
	VisitCallExpression(ce *CallExpression) error
	VisitParams(pa *Params) error
	VisitProgram(p *Program) error
}
