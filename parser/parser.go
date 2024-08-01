package parser

// currently we can parse Int, return, expression, if statements
// const statements still to be implemented
// implementing addition types bool, string, float is still remaining

//  Tests for all the functions is still remaining

import (
	"bytes"
	"fmt"
	"lim-lang/ast"
	"lim-lang/lexer"
	"lim-lang/token"
	"strconv"
	"strings"
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
)

type Parser struct {
	l *lexer.Lexer

	errors []parseError

	curToken  token.Token
	peekToken token.Token

	curLineNum int

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}
type parseError struct {
	errStr            string
	errMessage        string
	errPositionMarker string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []parseError{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.TRUE] = p.parseBoolean
	p.prefixParseFns[token.FALSE] = p.parseBoolean
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression
	// p.prefixParseFns[token.IF] = p.parseIfExpression
	// p.prefixParseFns[token.FUNCTION] = p.parseFunctionLiteral

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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	if p.curToken.Type == token.ENDOFLINE {
		p.curLineNum += 1
		p.nextToken()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		if p.curToken.Type == token.ILLEGAL {
			// p.peekError(p.curToken.Type)
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
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Keywork_INT:
		return p.parseIntStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.Keywork_BOOL:
		return p.parseBoolStatement()
	// case token.Keywork_STRING:
	// 	return p.parseStringStatement()
	// case token.Keywork_FLOAT:
	// 	return p.parseFloatStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
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
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	rst := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	rst.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return rst
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
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	// msg := fmt.Sprintf("no prefix parse function for %s found", t)
	// This error function is still incomplete and doesn't return a proper error
	// fmt.Printf("no prefix parse function for %s found\n", t)
	p.peekError(t)
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
		// fmt.Println("tf end: ", p.curToken, p.peekToken)
	}
	return rootIfStmt
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

	if p.peekTokenIs(token.Keywork_INT) {
		p.nextToken()
		lit.ReturnType = p.curToken
	} else if p.peekTokenIs(token.Keywork_BOOL) {
		p.nextToken()
		lit.ReturnType = p.curToken
	} else if p.peekTokenIs(token.Keywork_STRING) {
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
	if p.curToken.Type == token.Keywork_INT || p.curToken.Type == token.Keywork_BOOL || p.curToken.Type == token.Keywork_STRING {
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
		if p.curToken.Type == token.Keywork_INT || p.curToken.Type == token.Keywork_BOOL || p.curToken.Type == token.Keywork_STRING {
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
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(ExpectedToken token.TokenType) {
	var epm string
	i := p.l.AtCharNumFromCurrentLine - len(p.curToken.Literal)
	for i > 0 {
		epm += "-"
		i -= 1
	}
	for l := 0; l < len(p.curToken.Literal); l++ {
		epm += "^"
	}
	var erMess string
	if ExpectedToken == token.ILLEGAL {
		erMess = fmt.Sprint("Encountered illegal token ", p.curToken.Literal, " at line ", p.curLineNum)
	} else {
		erMess = fmt.Sprint("\n", "Expected ' ", ExpectedToken, " ' at line ", p.curLineNum)
	}
	pErr := parseError{
		errStr:            strings.ReplaceAll(p.l.ReadErrorLine(), "\t", " "),
		errMessage:        erMess,
		errPositionMarker: epm,
	}
	p.errors = append(p.errors, pErr)
}

func (p *Parser) GetErrors() []parseError {
	return p.errors
}

func (p *Parser) GetErrorsStr() []string {
	var eMsgs []string
	for _, msg := range p.errors {
		fmt.Println(msg.errMessage)
		fmt.Println(msg.errStr)
		fmt.Println(msg.errPositionMarker)
		var eMsg bytes.Buffer
		eMsg.WriteString(msg.errMessage)
		eMsg.WriteString("\n")
		eMsg.WriteString(msg.errStr)
		eMsg.WriteString("\n")
		eMsg.WriteString(msg.errPositionMarker)
		eMsg.WriteString("\n")
		eMsgs = append(eMsgs, eMsg.String())
	}
	return eMsgs
}
