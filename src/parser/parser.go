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
	start := p.this()
	tk := p.this()
	switch tk.Type {
	case lexer.TK_INTEGER:
		val, _ := strconv.ParseInt(tk.Literal, 10, 64)
		lit := &Literal{
			BaseNode: BaseNode{
				Start: start,
				End:   start,
			},
			Value: val,
		}
		p.advance()
		return lit, nil
	case lexer.TK_FLOAT:
		val, _ := strconv.ParseFloat(tk.Literal, 64)
		lit := &Literal{
			BaseNode: BaseNode{
				Start: start,
				End:   start,
			},
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
		if p.next.Type == lexer.TK_LEFT_PAREN {
			start := p.this()
			ident := p.this().Literal
			p.advance()
			p.advance()
			var params []Node
			for {
				if p.this().Type == lexer.TK_RIGHT_PAREN {
					break
				}
				expr, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				params = append(params, expr)
				if p.this().Type == lexer.TK_COMMA {
					p.advance()
				}
			}
			end := p.this()
			p.advance()
			return &CallExpression{
				BaseNode: BaseNode{
					Start: start,
					End:   end,
				},
				Identifier: ident,
				Args:       params,
			}, nil
		} else {
			lookUp := &VariableLookup{
				Id: p.this().Literal,
			}
			p.advance() // Goes to next token
			return lookUp, nil
		}
	}
	return nil, &exeptions.CompilerError{
		File:    p.next.Position.Filename,
		Line:    p.next.Position.Line,
		Column:  p.next.Position.Column,
		Message: fmt.Sprintf("ParsingError: Expected expression, but got '%s'.", p.next.Literal),
	}
}

func (p *Parser) parseTerm() (Node, error) {
	start := p.this()
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
		if p.this().Type == lexer.TK_STAR || p.this().Type == lexer.TK_SLASH || p.this().Type == lexer.TK_MODULO {
			op := p.this().Literal
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			end := p.this()
			left = &BinaryExpression{
				BaseNode: BaseNode{
					Start: start,
					End:   end,
				},
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
	start := p.this()
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
			end := p.this()
			left = &BinaryExpression{
				BaseNode: BaseNode{
					Start: start,
					End:   end,
				},
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
	start := p.this()
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
			end := p.this()
			left = &BinaryExpression{
				BaseNode: BaseNode{
					Start: start,
					End:   end,
				},
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
	start := p.this()
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
	end := p.this()
	exitStmt := &ExitStatment{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Argument: node,
	}
	return exitStmt, nil

}
func (p *Parser) parseReturnStatement() (Node, error) {
	start := p.this()
	p.advance() // Skips return
	if p.this().Type == lexer.TK_SEMICOLON {
		return &ReturnStatment{
			BaseNode: BaseNode{
				Start: start,
				End:   start,
			},
			Argument: nil,
		}, nil
	}
	node, parsingError := p.parseExpression()
	if parsingError != nil {
		return nil, parsingError
	}
	end := p.this()

	return &ReturnStatment{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Argument: node,
	}, nil
}

func (p *Parser) parseBreakStatement() (Node, error) {
	start := p.this()
	p.advance() // Skips break
	end := p.this()
	_, err := p.match(lexer.TK_SEMICOLON, ";")
	if err != nil {
		return nil, err
	}
	return &BreakStatment{
		BaseNode{
			Start: start,
			End:   end,
		},
	}, nil

}

func (p *Parser) parseDeleclarationStatement() (Node, error) {
	start := p.this()
	// TODO: No const implement yet
	p.advance() // skips "var/val"
	startDec := p.this()
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
	end := p.this()
	dec := &VariableDeclaration{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Kind: "var",
		Declaration: &VariableDeclarator{
			BaseNode: BaseNode{
				Start: startDec,
				End:   end,
			},
			Id:   ident,
			Init: node,
		},
	}

	return dec, nil

}

func (p *Parser) parseAssignmentStatement() (Node, error) {
	start := p.this()
	ident := p.current.Literal
	p.advance()
	var value Node
	var err error
	switch p.current.Literal {
	case "=":
		p.advance()
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	case "+=":
		p.advance()
		tmp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		value = &BinaryExpression{
			Left: &VariableLookup{
				Id: ident,
			},
			Right:    tmp,
			Operator: "+",
		}
	case "-=":
		p.advance()
		tmp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		value = &BinaryExpression{
			Left: &VariableLookup{
				Id: ident,
			},
			Right:    tmp,
			Operator: "-",
		}
	case "&=":
		p.advance()
		tmp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		value = &BinaryExpression{
			Left: &VariableLookup{
				Id: ident,
			},
			Right:    tmp,
			Operator: "&",
		}
	case "|=":
		p.advance()
		tmp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		value = &BinaryExpression{
			Left: &VariableLookup{
				Id: ident,
			},
			Right:    tmp,
			Operator: "|",
		}
	case "*=":
		p.advance()
		tmp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		value = &BinaryExpression{
			Left: &VariableLookup{
				Id: ident,
			},
			Right:    tmp,
			Operator: "*",
		}
	default:
		return nil, &exeptions.CompilerError{
			File:    p.this().Position.Filename,
			Line:    p.this().Position.Line,
			Column:  p.this().Position.Column,
			Message: fmt.Sprintf("Encountenterd unknow assignment operator '%s'.", p.this().Literal),
		}
	}
	end := p.this()
	assignment := &AssignmentStatement{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Identifier: ident,
		Value:      value,
	}
	return assignment, nil
}

func (p *Parser) parseBlockStatement() (Node, error) {
	start := p.this()
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
	end := p.this()
	p.advance()
	return &BlockStatement{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Instructions: stmts,
	}, nil
}

func (p *Parser) parseWhileStatement() (Node, error) {
	start := p.this()
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
	end := p.this()
	return &WhileStatment{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Test: test,
		Body: body,
	}, nil
}

func (p *Parser) parseForStatement() (Node, error) {
	start := p.this()
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
	end := p.this()

	return &ForStatment{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Init:   init,
		Test:   test,
		Update: update,
		Body:   body,
	}, nil
}

func (p *Parser) parseIfStatement() (Node, error) {
	start := p.this()
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
	var end = p.this()
	return &IfStatment{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Test:       test,
		Consequent: consequent,
		Alternate:  alternate,
	}, nil
}

func (p *Parser) parseParams() (Node, error) {
	var args []string
	for {
		if p.this().Type == lexer.TK_RIGHT_PAREN {
			break
		}
		ident, err := p.match(lexer.TK_IDENTIFIER, "identifier")
		if err != nil {
			return nil, err
		}
		args = append(args, ident.Literal)
		if p.this().Type == lexer.TK_COMMA {
			p.advance()
		}
	}
	return &Params{
		Args: args,
	}, nil
}

func (p *Parser) parseFunctionDeclaration() (Node, error) {
	start := p.this()
	p.advance()                                              // Skips fn
	ident, err := p.match(lexer.TK_IDENTIFIER, "identifier") // Skips function name
	if err != nil {
		return nil, err
	}
	if _, err := p.match(lexer.TK_LEFT_PAREN, "("); err != nil {
		return nil, err
	}
	var params Node
	if p.this().Type != lexer.TK_RIGHT_PAREN {
		params, err = p.parseParams()
	}
	if _, err := p.match(lexer.TK_RIGHT_PAREN, ")"); err != nil {
		return nil, err
	}
	ret, err := p.match(lexer.TK_IDENTIFIER, "return type") // Skips return type
	if err != nil {
		return nil, err
	}
	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	end := p.this()
	return &FunctionDeclaration{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Function: &Function{
			BaseNode: BaseNode{
				Start: start,
				End:   end,
			},
			Id:         ident.Literal,
			Parameters: params,
			ReturnType: ret.Literal,
			Body:       body,
		},
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
	case lexer.TK_RETURN:
		stmt, err := p.parseReturnStatement()
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
	case lexer.TK_BREAK:
		return p.parseBreakStatement()
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
	case lexer.TK_FN:
		return p.parseFunctionDeclaration()
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
	start := p.current
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
	end := p.current
	prog := &Program{
		BaseNode: BaseNode{
			Start: start,
			End:   end,
		},
		Instructions: nodes,
	}
	return prog, nil
}

func (p *Parser) Parse() (Node, error) {
	p.current = p.tk.GetToken()
	p.next = p.tk.GetToken()
	return p.parseProgram()
}
