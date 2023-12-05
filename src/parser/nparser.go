package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Program struct {
	Statements []*Statement `parser:"@@*"`
}

func (pg *Program) Accept(v Visitor) error {
	return v.VisitProgram(pg)
}

type Statement struct {
	VarDeclarationStatement *VarDeclarationStatement `parser:"   @@"`
	AssignmentStatement     *AssignmentStatement     `parser:" | @@"`
	FnDeclarationStatement  *FnDeclarationStatement  `parser:" | @@"`
	ImportStatement         *ImportStatement         `parser:" | @@"`
	BlockStatement          *BlockStatement          `parser:" | @@"`
	IfStatement             *IfStatement             `parser:" | @@"`
	WhileStatement          *WhileStatement          `parser:" | @@"`
	ForStatement            *ForStatement            `parser:" | @@"`
	BreakStatement          *BreakStatement          `parser:" | @@"`
	ReturnStatement         *ReturnStatement         `parser:" | @@"`
	ContinueStatement       *ContinueStatement       `parser:" | @@"`
	ExitStatement           *ExitStatement           `parser:" | @@"`
	CallStatement           *CallStatement           `parser:" | @@"`
}

func (pg *Statement) Accept(v Visitor) error {
	if pg.VarDeclarationStatement != nil {
		return v.VisitVarDeclarationStatement(pg.VarDeclarationStatement)
	} else if pg.AssignmentStatement != nil {
		return v.VisitAssignmentStatement(pg.AssignmentStatement)
	} else if pg.FnDeclarationStatement != nil {
		return v.VisitFnDeclarationStatement(pg.FnDeclarationStatement)
	} else if pg.ImportStatement != nil {
		return v.VisitImportStatement(pg.ImportStatement)
	} else if pg.BlockStatement != nil {
		return v.VisitBlockStatement(pg.BlockStatement)
	} else if pg.IfStatement != nil {
		return v.VisitIfStatement(pg.IfStatement)
	} else if pg.WhileStatement != nil {
		return v.VisitWhileStatement(pg.WhileStatement)
	} else if pg.ForStatement != nil {
		return v.VisitForStatement(pg.ForStatement)
	} else if pg.BreakStatement != nil {
		return v.VisitBreakStatement(pg.BreakStatement)
	} else if pg.ReturnStatement != nil {
		return v.VisitReturnStatement(pg.ReturnStatement)
	} else if pg.ContinueStatement != nil {
		return v.VisitContinueStatement(pg.ContinueStatement)
	} else if pg.ExitStatement != nil {
		return v.VisitExitStatement(pg.ExitStatement)
	}
	return v.VisitStatement(pg)
}

type ImportStatement struct {
	Imports []*string `parser:"'import' '{' @Ident (',' @Ident)* '}'"`
	File    string    `parser:"'from' @String"`
}

func (varDec *ImportStatement) Accept(v Visitor) error {
	return v.VisitImportStatement(varDec)
}

type VarDeclarationStatement struct {
	Pos lexer.Position
	Dec *VarDeclarationExpression `parser:"@@ ';'"`
}

func (varDec *VarDeclarationStatement) Accept(v Visitor) error {
	return v.VisitVarDeclarationStatement(varDec)
}

type CallStatement struct {
	Pos  lexer.Position
	Expr *CallExpression `parser:"@@ ';'"`
}

func (varDec *CallStatement) Accept(v Visitor) error {
	return v.VisitCallStatement(varDec)
}

type AssignmentStatement struct {
	Pos lexer.Position
	Ass *AssignmentExpression `parser:"@@ ';'"`
}

func (assSt *AssignmentStatement) Accept(v Visitor) error {
	return v.VisitAssignmentStatement(assSt)
}

type FnDeclarationStatement struct {
	Pos        lexer.Position
	Identifier string          `parser:"'fn' @Ident"`
	Params     []*Param        `parser:"'(' @@? (',' @@)* ')'"`
	Type       Type            `parser:"':' @@"`
	Body       *BlockStatement `parser:"@@"`
}

func (t *FnDeclarationStatement) Accept(v Visitor) error {
	return v.VisitFnDeclarationStatement(t)
}

type IfStatement struct {
	Pos        lexer.Position
	Test       *Expression     `parser:"'if' '(' @@ ')'"`
	Consequent *BlockStatement `parser:"@@"`
	Alternate  *BlockStatement `parser:"('else' @@)?"`
}

func (t *IfStatement) Accept(v Visitor) error {
	return v.VisitIfStatement(t)
}

type WhileStatement struct {
	Pos  lexer.Position
	Test *Expression     `parser:"'while' '(' @@ ')'"`
	Body *BlockStatement `parser:"@@"`
}

func (t *WhileStatement) Accept(v Visitor) error {
	return v.VisitWhileStatement(t)
}

type ForStatement struct {
	Pos  lexer.Position
	Init *VarDeclarationExpression `parser:"'for' '(' @@? "`
	Test *Expression               `parser:"';' @@?"`
	Inc  *AssignmentExpression     `parser:"';' @@? ')'"`
	Body *BlockStatement           `parser:"@@"`
}

func (t *ForStatement) Accept(v Visitor) error {
	return v.VisitForStatement(t)
}

type BreakStatement struct {
	Value string `parser:"'break' ';'"`
}

func (t *BreakStatement) Accept(v Visitor) error {
	return v.VisitBreakStatement(t)
}

type ContinueStatement struct {
	Value string `parser:"'continue' ';'"`
}

func (t *ContinueStatement) Accept(v Visitor) error {
	return v.VisitContinueStatement(t)
}

type ReturnStatement struct {
	Value Expression `parser:"'return' @@ ';'"`
}

func (t *ReturnStatement) Accept(v Visitor) error {
	return v.VisitReturnStatement(t)
}

type ExitStatement struct {
	Value Expression `parser:"'exit' '(' @@ ')' ';'"`
}

func (t *ExitStatement) Accept(v Visitor) error {
	return v.VisitExitStatement(t)
}

type BlockStatement struct {
	Instructions []*Statement `parser:"'{' @@* '}'"`
}

func (t *BlockStatement) Accept(v Visitor) error {
	return v.VisitBlockStatement(t)
}

type AssignmentExpression struct {
	Identifier       string            `parser:"@Ident"`
	Operator         string            `parser:"@('='|'+='|'-='|'*='|'/=' |'|=')"`
	Expression       *Expression       `parser:"@@?"`
	ArrayInitializer *ArrayInitializer `parser:"@@?"`
}

func (t *AssignmentExpression) Accept(v Visitor) error {
	return v.VisitAssignmentExpression(t)
}

type VarDeclarationExpression struct {
	Specifier        string            `parser:"@('let'|'const')"`
	Identifier       string            `parser:"@Ident"`
	Type             *Type             `parser:"':' @@?"`
	Expression       *Expression       `parser:"('=' @@)?"`
	ArrayInitializer *ArrayInitializer `parser:"('=' @@)?"`
}

func (t *VarDeclarationExpression) Accept(v Visitor) error {
	return v.VisitVarDeclarationExpression(t)
}

type ArrayInitializer struct {
	Values []Value `parser:"'[' @@? (',' @@)* ']'"`
}

func (t *ArrayInitializer) Accept(v Visitor) error {
	return v.VisitArrayInitializer(t)
}

type Expression struct {
	Pos   lexer.Position
	Value LogicalOrExpression `parser:"@@"`
}

func (t *Expression) Accept(v Visitor) error {
	return v.VisitExpression(t)
}

type LogicalOrExpression struct {
	Left  LogicalAndExpression `parser:"@@"`
	Right *LogicalOrExpression `parser:"('||' @@)?"`
}

func (t *LogicalOrExpression) Accept(v Visitor) error {
	return v.VisitLogicalOrExpression(t)
}

type LogicalAndExpression struct {
	Left  EqualityExpression    `parser:"@@"`
	Right *LogicalAndExpression `parser:"('&&' @@)?"`
}

func (t *LogicalAndExpression) Accept(v Visitor) error {
	return v.VisitLogicalAndExpression(t)
}

type EqualityExpression struct {
	Left     RelationalExpression `parser:"@@"`
	Operator *string              `parser:"@( '==' | '!=' )?"`
	Right    *EqualityExpression  `parser:"@@?"`
}

func (t *EqualityExpression) Accept(v Visitor) error {
	return v.VisitEqualityExpression(t)
}

type RelationalExpression struct {
	Left     AdditiveExpression    `parser:"@@"`
	Operator *string               `parser:"@( '<' | '>' | '<=' | '>=' )?"`
	Right    *RelationalExpression `parser:"@@?"`
}

func (t *RelationalExpression) Accept(v Visitor) error {
	return v.VisitRelationalExpression(t)
}

type AdditiveExpression struct {
	Left     MultiplicativeExpression `parser:"@@"`
	Operator *string                  `parser:"@( '+' | '-' )?"`
	Right    *AdditiveExpression      `parser:"@@?"`
}

func (t *AdditiveExpression) Accept(v Visitor) error {
	return v.VisitAdditiveExpression(t)
}

type MultiplicativeExpression struct {
	Left     UnaryExpression           `parser:"@@"`
	Operator *string                   `parser:"@( '*' | '/' | '%' )?"`
	Right    *MultiplicativeExpression `parser:"@@?"`
}

func (t *MultiplicativeExpression) Accept(v Visitor) error {
	return v.VisitMultiplicativeExpression(t)
}

type UnaryExpression struct {
	Value *PrimaryExpression `parser:"@@"`
}

func (t *UnaryExpression) Accept(v Visitor) error {
	return v.VisitUnaryExpression(t)
}

type PrimaryExpression struct {
	Call        *CallExpression `parser:"@@"`
	Literal     *Value          `parser:"| @@"`
	ArrayLookup *ArrayLookup    `parser:"| @@"`
	Expression  *Expression     `parser:"| '(' @@ ')'"`
	Identifier  *string         `parser:"| @Ident"`
}

func (t *PrimaryExpression) Accept(v Visitor) error {
	return v.VisitPrimaryExpression(t)
}

type CallExpression struct {
	FunctionName string        `parser:"@Ident"`
	Arguments    []*Expression `parser:"'(' @@? (',' @@)* ')'"`
}

func (t *CallExpression) Accept(v Visitor) error {
	return v.VisitCallExpression(t)
}

type ArrayLookup struct {
	Identifier string     `parser:"@Ident"`
	Index      Expression `parser:"'[' @@ ']'"`
}

func (t *ArrayLookup) Accept(v Visitor) error {
	return v.VisitArrayLookup(t)
}

// Utils
type Value interface{ value() }

type Float struct {
	Value float64 `parser:"@Float"`
}

func (f Float) value() {}

type Int struct {
	Value int `parser:"@Int"`
}

func (f Int) value() {}

type String struct {
	Value string `parser:"@String"`
}

func (f String) value() {}

type Bool struct {
	Value string `parser:"@('true' | 'false')"`
}

func (f Bool) value() {}

type Param struct {
	Pos        lexer.Position
	Identifier string `parser:"@Ident"`
	Type       Type   `parser:"':' @@"`
}

type Type struct {
	Pos   lexer.Position
	Base  string `parser:"@( 'int8' | 'int16' | 'int32' | 'int64' | 'uint8' | 'uint16' | 'uint32' | 'uint64' | 'bool' | 'void' | 'float32' | 'float64' )"`
	Array bool   `parser:"@('[' ']')*"`
}

func NewNParser() *participle.Parser[Program] {
	graphQLLexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: `(?:#|//)[^\n]*\n?`},
		{Name: "Ident", Pattern: `[a-zA-Z]\w*`},
		{Name: "Float", Pattern: `(?:\d*\.\d+|\d+\.\d*)`}, // Pattern for Float
		{Name: "Int", Pattern: `\d+`},                     // Pattern for Int
		{Name: "String", Pattern: `'(\\'|[^'])*'`},        // Pattern for String      // Two-character tokens
		{Name: "Punct", Pattern: `\+=|-=|/=|\*=|<=|>=|==|!=|&&|\|\||[-[!@#$%^&*()+_={}\[\]\|:;"'<,>.?/]`},
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	})

	parser := participle.MustBuild[Program](
		participle.Lexer(graphQLLexer),
		participle.Elide("Comment", "Whitespace"),
		participle.UseLookahead(2),
		participle.Union[Value](Float{}, Int{}, String{}, Bool{}),
	)
	return parser
}

type Visitor interface {
	VisitProgram(t *Program) error
	VisitVarDeclarationStatement(t *VarDeclarationStatement) error
	VisitAssignmentStatement(t *AssignmentStatement) error
	VisitFnDeclarationStatement(t *FnDeclarationStatement) error
	VisitIfStatement(t *IfStatement) error
	VisitWhileStatement(t *WhileStatement) error
	VisitForStatement(t *ForStatement) error
	VisitBreakStatement(t *BreakStatement) error
	VisitContinueStatement(t *ContinueStatement) error
	VisitReturnStatement(t *ReturnStatement) error
	VisitBlockStatement(t *BlockStatement) error
	VisitAssignmentExpression(t *AssignmentExpression) error
	VisitVarDeclarationExpression(t *VarDeclarationExpression) error
	VisitExpression(t *Expression) error
	VisitLogicalOrExpression(t *LogicalOrExpression) error
	VisitLogicalAndExpression(t *LogicalAndExpression) error
	VisitEqualityExpression(t *EqualityExpression) error
	VisitRelationalExpression(t *RelationalExpression) error
	VisitAdditiveExpression(t *AdditiveExpression) error
	VisitMultiplicativeExpression(t *MultiplicativeExpression) error
	VisitUnaryExpression(t *UnaryExpression) error
	VisitPrimaryExpression(t *PrimaryExpression) error
	VisitCallExpression(t *CallExpression) error
	VisitImportStatement(t *ImportStatement) error
	VisitCallStatement(t *CallStatement) error
	VisitStatement(t *Statement) error
	VisitExitStatement(t *ExitStatement) error
	VisitArrayLookup(t *ArrayLookup) error
	VisitArrayInitializer(t *ArrayInitializer) error
}
