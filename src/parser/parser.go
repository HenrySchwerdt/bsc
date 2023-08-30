package parser

import (
	"bsc/src/exeptions"
	"bsc/src/lexer"
	"fmt"
	"strconv"
)

type Parser struct {
	tk      *lexer.Tokenizer
	current lexer.Token
	next    lexer.Token
}

func NewParser(tk *lexer.Tokenizer) *Parser {
	return &Parser{
		tk: tk,
	}
}

func (p *Parser) advance() lexer.Token {
	tmp := p.current
	p.current = p.next
	p.next = p.tk.GetToken()
	return tmp
}

func (p *Parser) peek() lexer.Token {
	return p.next
}

func (p *Parser) match(tokenType lexer.TokenType, keyword string) error {
	if p.peek().Type == tokenType {
		p.advance()
		return nil
	}
	return &exeptions.CompilerError{
		File:    p.next.Position.Filename,
		Line:    p.next.Position.Line,
		Column:  p.next.Position.Column,
		Message: fmt.Sprintf("ParsingError: Expected '%s', but got '%s'.", keyword, p.next.Literal),
	}
}

func (p *Parser) parseLiteral() (Node, error) {
	tk := p.peek()
	switch tk.Type {
	case lexer.TK_INTEGER:
		val, _ := strconv.ParseInt(tk.Literal, 10, 64)
		lit := &Literal{
			Value: val,
		}
		p.advance()
		return lit, nil
	case lexer.TK_FLOAT:
		val, _ := strconv.ParseFloat(tk.Literal, 64)
		lit := &Literal{
			Value: val,
		}
		p.advance()
		return lit, nil
	default:
		return nil, &exeptions.CompilerError{
			File:    p.next.Position.Filename,
			Line:    p.next.Position.Line,
			Column:  p.next.Position.Column,
			Message: fmt.Sprintf("ParsingError: Expected literal, but got '%s'.", p.next.Literal),
		}
	}
}

func (p *Parser) parseFactor() (Node, error) {
	if p.peek().Type == lexer.TK_LEFT_PAREN {
		p.advance()
		expr, exprErr := p.parseExpression()
		if exprErr != nil {
			return nil, exprErr
		}
		if err := p.match(lexer.TK_RIGHT_PAREN, ")"); err != nil {
			return nil, err
		}
		return expr, nil
	}
	if p.peek().Type == lexer.TK_INTEGER || p.peek().Type == lexer.TK_FLOAT {
		return p.parseLiteral()
	}
	if p.peek().Type == lexer.TK_IDENTIFIER {
		lookUp := &VariableLookup{
			Id: p.peek().Literal,
		}
		p.advance()
		return lookUp, nil
	}
	return nil, &exeptions.CompilerError{
		File:    p.next.Position.Filename,
		Line:    p.next.Position.Line,
		Column:  p.next.Position.Column,
		Message: fmt.Sprintf("ParsingError: Expected expression, but got '%s'.", p.next.Literal),
	}
}

func (p *Parser) parseTerm() (Node, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
		if p.peek().Type == lexer.TK_STAR || p.peek().Type == lexer.TK_SLASH {
			op := p.advance().Literal
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			binExpr := &BinaryExpression{
				Left:     left,
				Right:    right,
				Operator: op,
			}
			return binExpr, nil
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) parseComp() (Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		if p.peek().Type == lexer.TK_EQUAL_EQUAL || p.peek().Type == lexer.TK_GREATER || p.peek().Type == lexer.TK_LESS || p.peek().Type == lexer.TK_GREATER_EQUAL || p.peek().Type == lexer.TK_LESS_EQUAL {
			op := p.advance().Literal
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			binExpr := &BinaryExpression{
				Left:     left,
				Right:    right,
				Operator: op,
			}
			return binExpr, nil
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) parseExpression() (Node, error) {
	left, err := p.parseComp()
	if err != nil {
		return nil, err
	}
	for {
		if p.peek().Type == lexer.TK_PLUS || p.peek().Type == lexer.TK_MINUS {
			op := p.advance().Literal
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			binExpr := &BinaryExpression{
				Left:     left,
				Right:    right,
				Operator: op,
			}
			return binExpr, nil
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) parseExitStatement() (Node, error) {
	err := p.match(lexer.TK_LEFT_PAREN, "(")
	if err != nil {
		return nil, err
	}
	node, parsingError := p.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}

	err = p.match(lexer.TK_RIGHT_PAREN, ")")
	if err != nil {
		return nil, err
	}

	err = p.match(lexer.TK_SEMICOLON, ";")
	if err != nil {
		return nil, err
	}

	exitStmt := &ExitStatment{
		Argument: node,
	}
	p.advance()
	return exitStmt, nil

}

func (p *Parser) parseDeleclarationStatement() (Node, error) {
	// TODO: No const implement yet
	if err := p.match(lexer.TK_IDENTIFIER, "identifier"); err != nil {
		return nil, err
	}
	ident := p.current.Literal
	if err := p.match(lexer.TK_EQUAL, "="); err != nil {
		return nil, err
	}
	node, parsingError := p.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	if err := p.match(lexer.TK_SEMICOLON, ";"); err != nil {
		return nil, err
	}
	dec := &VariableDeclaration{
		Kind: "var",
		Declaration: &VariableDeclarator{
			Id:   ident,
			Init: node,
		},
	}
	p.advance()
	return dec, nil

}

func (p *Parser) parseAssignmentStatement() (Node, error) {
	ident := p.current.Literal

	if err := p.match(lexer.TK_EQUAL, "="); err != nil {
		return nil, err
	}
	node, parsingError := p.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}

	if err := p.match(lexer.TK_SEMICOLON, ";"); err != nil {
		return nil, err
	}
	p.advance()
	p.advance()

	assignment := &AssignmentStatement{
		Identifier: ident,
		Value:      node,
	}
	return assignment, nil
}

func (p *Parser) parseStatement() (Node, error) {
	switch p.current.Type {
	case lexer.TK_EXIT:
		return p.parseExitStatement()
	case lexer.TK_VAL:
		return p.parseDeleclarationStatement()
	case lexer.TK_VAR:
		return p.parseDeleclarationStatement()
	case lexer.TK_IDENTIFIER:
		return p.parseAssignmentStatement()
	default:
		return nil, &exeptions.CompilerError{
			File:    p.next.Position.Filename,
			Line:    p.next.Position.Line,
			Column:  p.next.Position.Column,
			Message: fmt.Sprintf("ParsingError: Expected statement, but got '%s'.", p.current.Literal),
		}
	}
}

func (p *Parser) parseProgram() (Node, error) {
	nodes := make([]Node, 0)

	for {
		if p.next.Type == lexer.TK_EOF {
			break
		}
		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	prog := &Program{
		Instructions: nodes,
	}
	return prog, nil
}

func (p *Parser) Parse() (Node, error) {
	p.current = p.tk.GetToken()
	p.next = p.tk.GetToken()
	return p.parseProgram()
}
