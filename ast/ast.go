package ast

import (
	"bytes"
	"fmt"
	"lim-lang/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}
type Expression interface {
	Node
	expressionNode()
}
type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token        token.Token
	HoldsVarType token.Token
	Value        string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	str := fmt.Sprintf("%s %s ", i.HoldsVarType.Literal, i.TokenLiteral())

	return str
}

type IntStatement struct {
	Token token.Token // token.INT
	Name  *Identifier
	Value Expression
}

func (is *IntStatement) expressionNode() {}
func (is *IntStatement) statementNode()  {}
func (is *IntStatement) TokenLiteral() string {
	return is.Token.Literal
}
func (is *IntStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.TokenLiteral())
	out.WriteString(is.Name.String())
	out.WriteString("= ")
	if is.Value != nil {
		out.WriteString(is.Value.String())
	}
	// out.WriteString(";")
	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	// out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(fmt.Sprint("    ", s.String(), "\n"))
	}

	return out.String()
}

type IfStatement struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	NextCase    *IfStatement
}

// func (ie *IfExpression) expressionNode()      {}
func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	isRootNode := true
	cNode := is
	for cNode.NextCase != nil {
		out.WriteString(is.Token.Literal)
		if is.Condition != nil {
			if !isRootNode {
				out.WriteString(" else")
			}
			out.WriteString(" ")
			out.WriteString(is.Condition.String())
		}
		out.WriteString(" {\n")
		out.WriteString(is.Consequence.String())
		out.WriteString("}")
		isRootNode = false
		cNode = cNode.NextCase
	}
	out.WriteString(cNode.Token.Literal)
	if cNode.Condition != nil {
		if !isRootNode {
			out.WriteString(" else")
		}
		out.WriteString(" ")
		out.WriteString(cNode.Condition.String())
	}
	out.WriteString(" {\n")
	out.WriteString(cNode.Consequence.String())
	out.WriteString("}")

	return out.String()
}

type FunctionStatement struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	ReturnType token.Token
	FnName     string
	Body       *BlockStatement
}

// func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionStatement) statementNode()       {}
func (fl *FunctionStatement) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fl.FnName)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.ReturnType.Literal)
	out.WriteString(" {")
	out.WriteString("\n")
	out.WriteString(fl.Body.String())
	out.WriteString("}")

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
