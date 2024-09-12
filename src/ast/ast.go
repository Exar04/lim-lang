package ast

import (
	"bytes"
	"fmt"
	"limLang/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
	GetTreeFormat() string
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

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
func (p *Program) GetTreeFormat() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.GetTreeFormat())
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
func (i *Identifier) GetTreeFormat() string { return "" }
func (i *Identifier) String() string {
	var str string
	str = fmt.Sprintf("%s %s", i.HoldsVarType.Literal, i.TokenLiteral())

	return str
}

type IntStatement struct {
	Token token.Token // token.INT
	Name  *Identifier
	Value Expression
}

func (is *IntStatement) expressionNode()       {}
func (is *IntStatement) statementNode()        {}
func (is *IntStatement) TokenLiteral() string  { return is.Token.Literal }
func (is *IntStatement) GetTreeFormat() string { return "" }
func (is *IntStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.TokenLiteral())
	out.WriteString(is.Name.String())
	out.WriteString(" = ")
	if is.Value != nil {
		out.WriteString(is.Value.String())
	}
	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()       {}
func (il *IntegerLiteral) TokenLiteral() string  { return il.Token.Literal }
func (il *IntegerLiteral) String() string        { return il.Token.Literal }
func (il *IntegerLiteral) GetTreeFormat() string { return il.Token.Literal }

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
	return out.String()
}
func (rs *ReturnStatement) GetTreeFormat() string { return "" }

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
func (es *ExpressionStatement) GetTreeFormat() string { return "" }

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
func (pe *PrefixExpression) GetTreeFormat() string { return "" }

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
func (ie *InfixExpression) GetTreeFormat() string { return "" }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()       {}
func (b *Boolean) TokenLiteral() string  { return b.Token.Literal }
func (b *Boolean) String() string        { return b.Token.Literal }
func (b *Boolean) GetTreeFormat() string { return "" }

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()        {}
func (bs *BlockStatement) TokenLiteral() string  { return bs.Token.Literal }
func (bs *BlockStatement) GetTreeFormat() string { return "" }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(fmt.Sprint("", s.String(), "\n"))
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
func (is *IfStatement) statementNode()        {}
func (is *IfStatement) TokenLiteral() string  { return is.Token.Literal }
func (is *IfStatement) GetTreeFormat() string { return "" }
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
func (fl *FunctionStatement) statementNode()        {}
func (fl *FunctionStatement) TokenLiteral() string  { return fl.Token.Literal }
func (fl *FunctionStatement) GetTreeFormat() string { return "" }
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

func (ce *CallExpression) expressionNode()       {}
func (ce *CallExpression) TokenLiteral() string  { return ce.Token.Literal }
func (ce *CallExpression) GetTreeFormat() string { return "" }
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

type BoolStatement struct {
	Token token.Token // token.BOOL
	Name  *Identifier
	Value Expression
}

func (is *BoolStatement) expressionNode()       {}
func (is *BoolStatement) statementNode()        {}
func (is *BoolStatement) GetTreeFormat() string { return "" }
func (is *BoolStatement) TokenLiteral() string {
	return is.Token.Literal
}
func (is *BoolStatement) String() string {
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

type StringVal struct {
	Token token.Token
	Value string
}

func (s *StringVal) expressionNode()       {}
func (s *StringVal) TokenLiteral() string  { return s.Token.Literal }
func (s *StringVal) String() string        { return s.Value }
func (s *StringVal) GetTreeFormat() string { return "" }

type StringStatement struct {
	Token token.Token // token.STRING
	Name  *Identifier
	Value Expression
}

func (ss *StringStatement) expressionNode()       {}
func (ss *StringStatement) statementNode()        {}
func (ss *StringStatement) GetTreeFormat() string { return "" }
func (ss *StringStatement) TokenLiteral() string {
	return ss.Token.Literal
}
func (ss *StringStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ss.TokenLiteral())
	out.WriteString(ss.Name.String())
	out.WriteString("= ")
	if ss.Value != nil {
		out.WriteString(ss.Value.String())
	}
	// out.WriteString(";")
	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token // ARRAY token
	Name     *Identifier
	Type     token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()       {}
func (ss *ArrayLiteral) statementNode()        {}
func (al *ArrayLiteral) TokenLiteral() string  { return al.Token.Literal }
func (al *ArrayLiteral) GetTreeFormat() string { return "" }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // the [ token
	Ident *Identifier
	Index Expression
}

func (ie *IndexExpression) expressionNode()       {}
func (ss *IndexExpression) statementNode()        {}
func (al *IndexExpression) GetTreeFormat() string { return "" }
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Ident.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
