package ast

import (
	"testing"

	"limLang/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&IntStatement{
				Token: token.Token{Type: token.INT, Literal: "int"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "int myVar =  anotherVar" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
