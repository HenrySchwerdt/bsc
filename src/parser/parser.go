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

func (p *Parser) this() lexer.Token {
	return p.current
}

func (p *Parser) match(tokenType lexer.TokenType, keyword string) (lexer.Token, error) {
	if p.this().Type == tokenType {
		matched := p.this()
		p.advance()
		return matched, nil
	}
	return lexer.Token{}, &exeptions.CompilerError{
		File:    p.next.Position.Filename,
		Line:    p.next.Position.Line,
		Column:  p.next.Position.Column,
		Message: fmt.Sprintf("ParsingError: Expected '%s', but got '%s'.", keyword, p.this().Literal),
	}
}

func (p *Parser) parseLiteral() (Node, error) {
	tk := p.this()
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

	if p.this().Type == lexer.TK_LEFT_PAREN {
		p.advance() // Skips '('
		expr, exprErr := p.parseExpression()
		if exprErr != nil {
			return nil, exprErr
		}
		if _, err := p.match(lexer.TK_RIGHT_PAREN, ")"); err != nil {
			return nil, err
		}
		return expr, nil
	}
	if p.this().Type == lexer.TK_INTEGER || p.this().Type == lexer.TK_FLOAT {

		return p.parseLiteral()
	}
	if p.this().Type == lexer.TK_IDENTIFIER {
		lookUp := &VariableLookup{
			Id: p.this().Literal,
		}
		p.advance() // Goes to next token
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
		if p.this().Type == lexer.TK_STAR || p.this().Type == lexer.TK_SLASH {
			op := p.this().Literal
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = &BinaryExpression{
				Left:     left,
				Right:    right,
				Operator: op,
			}
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
		if p.this().Type == lexer.TK_EQUAL_EQUAL || p.this().Type == lexer.TK_GREATER || p.this().Type == lexer.TK_LESS || p.this().Type == lexer.TK_GREATER_EQUAL || p.this().Type == lexer.TK_LESS_EQUAL {

			op := p.this().Literal
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = &BinaryExpression{
				Left:     left,
				Right:    right,
				Operator: op,
			}
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
		if p.this().Type == lexer.TK_PLUS || p.this().Type == lexer.TK_MINUS {
			op := p.this().Literal
			p.advance()
			right, err := p.parseComp()
			if err != nil {
				return nil, err
			}
			left = &BinaryExpression{
				Left:     left,
				Right:    right,
				Operator: op,
			}
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) parseExitStatement() (Node, error) {
	p.advance() // Skips exit
	_, err := p.match(lexer.TK_LEFT_PAREN, "(")
	if err != nil {
		return nil, err
	}
	node, parsingError := p.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}

	_, err = p.match(lexer.TK_RIGHT_PAREN, ")")
	if err != nil {
		return nil, err
	}

	exitStmt := &ExitStatment{
		Argument: node,
	}
	return exitStmt, nil

}

func (p *Parser) parseDeleclarationStatement() (Node, error) {
	// TODO: No const implement yet
	p.advance() // skips "var/val"
	identTk, err := p.match(lexer.TK_IDENTIFIER, "identifier")
	if err != nil {
		return nil, err
	}
	ident := identTk.Literal
	if _, err := p.match(lexer.TK_EQUAL, "="); err != nil {
		return nil, err
	}
	node, parsingError := p.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	dec := &VariableDeclaration{
		Kind: "var",
		Declaration: &VariableDeclarator{
			Id:   ident,
			Init: node,
		},
	}
	return dec, nil

}

func (p *Parser) parseAssignmentStatement() (Node, error) {
	ident := p.current.Literal
	p.advance()
	if _, err := p.match(lexer.TK_EQUAL, "="); err != nil {
		return nil, err
	}
	node, parsingError := p.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}

	assignment := &AssignmentStatement{
		Identifier: ident,
		Value:      node,
	}
	return assignment, nil
}

func (p *Parser) parseBlockStatement() (Node, error) {
	stmts := make([]Node, 0)
	p.advance()
	for {
		if p.current.Type == lexer.TK_RIGHT_BRACE || p.current.Type == lexer.TK_EOF {
			break
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	p.advance()
	return &BlockStatement{
		Instructions: stmts,
	}, nil
}

func (p *Parser) parseWhileStatement() (Node, error) {
	p.advance() // Skips 'while'
	if _, err := p.match(lexer.TK_LEFT_PAREN, "("); err != nil {
		return nil, err
	}
	test, err := p.parseComp()
	if err != nil {
		return nil, err
	}
	if _, err := p.match(lexer.TK_RIGHT_PAREN, ")"); err != nil {
		return nil, err
	}
	body, err := p.parseStatement()
	return &WhileStatment{
		Test: test,
		Body: body,
	}, nil
}

func (p *Parser) parseForStatement() (Node, error) {
	p.advance()
	_, err := p.match(lexer.TK_LEFT_PAREN, "(")
	if err != nil {
		return nil, err
	}

	init, err := p.parseStatement() // Parse the initializer
	if err != nil {
		return nil, err
	}

	test, err := p.parseComp() // Parse the condition
	if err != nil {
		return nil, err
	}

	_, err = p.match(lexer.TK_SEMICOLON, ";")
	if err != nil {
		return nil, err
	}

	update, err := p.parseAssignmentStatement() // Use parseExpression instead of parseAssignmentStatement
	if err != nil {
		return nil, err
	}

	_, err = p.match(lexer.TK_RIGHT_PAREN, ")")
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &ForStatment{
		Init:   init,
		Test:   test,
		Update: update,
		Body:   body,
	}, nil
}

func (p *Parser) parseIfStatement() (Node, error) {
	p.advance()                                                  // Skips if
	if _, err := p.match(lexer.TK_LEFT_PAREN, "("); err != nil { // Matches '('
		return nil, err
	}
	test, err := p.parseComp()
	if err != nil {
		return nil, err
	}
	if _, err := p.match(lexer.TK_RIGHT_PAREN, ")"); err != nil { // Matches ')'
		return nil, err
	}
	consequent, err := p.parseStatement()
	var alternate Node
	if p.current.Type == lexer.TK_ELSE {
		p.advance()
		alternate, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return &IfStatment{
		Test:       test,
		Consequent: consequent,
		Alternate:  alternate,
	}, nil
}

func (p *Parser) parseStatement() (Node, error) {
	switch p.this().Type {
	case lexer.TK_EXIT:
		stmt, err := p.parseExitStatement()
		if _, err := p.match(lexer.TK_SEMICOLON, ";"); err != nil {
			return nil, err
		}

		return stmt, err
	case lexer.TK_VAL:
		stmt, err := p.parseDeleclarationStatement()
		if _, err := p.match(lexer.TK_SEMICOLON, ";"); err != nil {
			return nil, err
		}

		return stmt, err
	case lexer.TK_VAR:
		stmt, err := p.parseDeleclarationStatement()
		if _, err := p.match(lexer.TK_SEMICOLON, ";"); err != nil {
			return nil, err
		}

		return stmt, err
	case lexer.TK_LEFT_BRACE:
		return p.parseBlockStatement()
	case lexer.TK_WHILE:
		return p.parseWhileStatement()
	case lexer.TK_FOR:
		return p.parseForStatement()
	case lexer.TK_IF:
		return p.parseIfStatement()
	case lexer.TK_IDENTIFIER:
		stmt, err := p.parseAssignmentStatement()
		if _, err := p.match(lexer.TK_SEMICOLON, ";"); err != nil {
			return nil, err
		}
		return stmt, err
	// case lexer.TK_FN:
	// 	return p.parseFunction()
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
