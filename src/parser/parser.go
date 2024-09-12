package parser

import (
	"fmt"
	"limLang/ast"
	"limLang/lexer"
	"limLang/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
	// postfixParseFn func() ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	curLineNum int

	err struct{}

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		if p.curToken.Type == token.ILLEGAL {
			// handle the illegal token
			return nil
		}
		if p.curToken.Type == token.ENDOFLINE {
			p.curLineNum += 1
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.TRUE] = p.parseBoolean
	p.prefixParseFns[token.FALSE] = p.parseBoolean
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression
	p.prefixParseFns[token.STRING] = p.parseStringLiteral
	// p.prefixParseFns[token.LBRACK] = p.parseArrayLiteral

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.SLASH] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[token.EQ] = p.parseInfixExpression
	p.infixParseFns[token.NOT_EQ] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression
	p.infixParseFns[token.LPAREN] = p.parseCallExpression

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	if p.curToken.Type == token.ENDOFLINE {
		p.curLineNum += 1
		p.nextToken()
	}
}

func (p *Parser) parseStatement() ast.Statement {
	// right now we are skipping all the \n tokens
	for p.curToken.Type == token.ENDOFLINE {
		p.nextToken()
	}
	switch p.curToken.Type {
	case token.Keyword_INT:
		if p.peekToken.Type == token.LBRACK {
			return p.parseArrayLiteral()
		}
		return p.parseIntStatement()
	case token.Keyword_BOOL:
		return p.parseBoolStatement()
	case token.Keyword_STRING:
		return p.parseStringStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENT:
		if p.peekToken.Type == token.DEFINE {
			return p.parseDefineStatement()
		} else if p.peekToken.Type == token.LBRACK {
			return p.parseIndexExpression()
			// } else if p.peekToken.Type == token.ASSIGN {
			// 	return p.parseReassignStatement()
		}
		fallthrough
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDefineStatement() ast.Statement {
	ident := p.curToken
	p.nextToken()
	p.nextToken()
	switch p.curToken.Type {
	case token.INT, token.PLUS, token.MINUS:
		stmt := &ast.IntStatement{Token: token.Token{Type: token.INT, Literal: "int"}}
		stmt.Name = &ast.Identifier{Token: ident, Value: ident.Literal}

		stmt.Value = p.parseExpression(LOWEST)
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt

	case token.STRING:
		stmt := &ast.StringStatement{Token: token.Token{Type: token.STRING, Literal: "string"}}
		stmt.Name = &ast.Identifier{Token: ident, Value: ident.Literal}

		stmt.Value = &ast.StringVal{Token: token.Token{Type: token.STRING, Literal: "string"}, Value: p.curToken.Literal}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	case token.TRUE, token.FALSE, token.BANG:
		stmt := &ast.BoolStatement{Token: token.Token{Type: token.BOOL, Literal: "bool"}}
		stmt.Name = &ast.Identifier{Token: ident, Value: ident.Literal}

		stmt.Value = p.parseExpression(LOWEST)
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}
	return nil
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringVal{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		// msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		// p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseIndexExpression() *ast.IndexExpression {
	exp := &ast.IndexExpression{Token: p.curToken, Ident: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}}
	p.nextToken()
	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACK) {
		return nil
	}
	return exp
}

func (p *Parser) parseArrayLiteral() *ast.ArrayLiteral {
	// int []arr = [1, 2 * 2, 3 + 3] // this is the structure of our array in lim lang
	array := &ast.ArrayLiteral{Token: token.Token{Type: token.ARRAY, Literal: token.ARRAY}}
	// now current token will be the statement type either int, bool, string or custome sturct type
	array.Type = p.curToken

	// now we will get [ paran next
	if !p.expectPeek(token.LBRACK) {
		return nil
	}
	// now we will get ] paran next
	if !p.expectPeek(token.RBRACK) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	array.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()

	// now our current token is [ paran
	array.Elements = p.parseExpressionList(token.RBRACK)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()

	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	// fmt.Println(p.curToken, p.peekToken)
	return list
}

func (p *Parser) parseIntStatement() *ast.IntStatement {
	stmt := &ast.IntStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		if p.peekTokenIs(token.SEMICOLON) {
			stmt.Value = &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "int"}, Value: 0}
			return stmt
		}
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBoolStatement() *ast.BoolStatement {
	stmt := &ast.BoolStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		if p.peekTokenIs(token.SEMICOLON) {
			stmt.Value = &ast.Boolean{Token: token.Token{Type: token.BOOL, Literal: "bool"}, Value: false}
			return stmt
		}
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseStringStatement() *ast.StringStatement {
	stmt := &ast.StringStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		if p.peekTokenIs(token.SEMICOLON) {
			stmt.Value = &ast.StringVal{Token: token.Token{Type: token.STRING, Literal: "string"}, Value: ""}
			return stmt
		}
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	rst := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	rst.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return rst
}
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	rootIfStmt := &ast.IfStatement{Token: p.curToken, NextCase: nil}

	p.nextToken()
	rootIfStmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	rootIfStmt.Consequence = p.parseBlockStatement()
	leafIfStmt := rootIfStmt

	for p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.peekTokenIs(token.IF) && p.peekTokenIs(token.LBRACE) {
			elseStmt := &ast.IfStatement{Token: p.curToken, Condition: nil, NextCase: nil}

			if !p.peekTokenIs(token.LBRACE) {
				return nil
			}
			p.nextToken()
			elseStmt.Consequence = p.parseBlockStatement()
			leafIfStmt.NextCase = elseStmt
			return rootIfStmt

		} else if !p.peekTokenIs(token.IF) && !p.peekTokenIs(token.LBRACE) {
			// this should give us a syntax error
			return nil
		}
		p.nextToken()
		nextIfStmt := &ast.IfStatement{Token: p.curToken, NextCase: nil}
		p.nextToken()

		nextIfStmt.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(token.LBRACE) {
			// p.nextToken()
			// fmt.Println(p.curToken)
			// throw an error
			return nil
		}
		nextIfStmt.Consequence = p.parseBlockStatement()
		leafIfStmt.NextCase = nextIfStmt
		leafIfStmt = nextIfStmt
	}
	return rootIfStmt
}
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	lit := &ast.FunctionStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	lit.FnName = p.curToken.Literal
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if p.peekTokenIs(token.Keyword_INT) {
		p.nextToken()
		lit.ReturnType = p.curToken
	} else if p.peekTokenIs(token.Keyword_BOOL) {
		p.nextToken()
		lit.ReturnType = p.curToken
	} else if p.peekTokenIs(token.Keyword_STRING) {
		p.nextToken()
		lit.ReturnType = p.curToken
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident1 := &ast.Identifier{}
	if p.curToken.Type == token.Keyword_INT || p.curToken.Type == token.Keyword_BOOL || p.curToken.Type == token.Keyword_STRING {
		ident1.HoldsVarType = p.curToken
		p.nextToken()
	} else {
		fmt.Println("there was error while parsing function paramerters!")
		p.nextToken()
	}
	ident1.Token = p.curToken
	ident1.Value = p.curToken.Literal
	identifiers = append(identifiers, ident1)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident2 := &ast.Identifier{}
		if p.curToken.Type == token.Keyword_INT || p.curToken.Type == token.Keyword_BOOL || p.curToken.Type == token.Keyword_STRING {
			ident2.HoldsVarType = p.curToken
			p.nextToken()
		} else {
			fmt.Println("there was error while parsing function paramerters!")
			p.nextToken()
		}
		ident2.Token = p.curToken
		ident2.Value = p.curToken.Literal
		identifiers = append(identifiers, ident2)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	// exp.Arguments = p.parseExpressionList(token.LPAREN)
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		// p.peekError(t)
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	// msg := fmt.Sprintf("no prefix parse function for %s found", t)
	// This error function is still incomplete and doesn't return a proper error
	// fmt.Printf("no prefix parse function for %s found\n", t)
	p.gotError(t)
}

func (p *Parser) gotError(t token.TokenType) {
	fmt.Println("error while parsing", t)
}
