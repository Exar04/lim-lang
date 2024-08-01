package parser

import (
	"fmt"
	"lim-lang/ast"
	"lim-lang/lexer"
	"lim-lang/token"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.GetErrors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		fmt.Println(msg.errMessage)
		fmt.Println(msg.errStr)
		fmt.Println(msg.errPositionMarker)
	}
	t.FailNow()
}

func TestIntStatements(t *testing.T) {
	input := `
	int x = 5;
    int foobar = 838383;
	int diss = -8 * -(4 + 4);
   `

	// 	failingInput := `
	//    int x 5;
	//    int = 10;
	//    int 838383;
	//    `
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	fmt.Println(program.Statements)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier    string
		expectedIdentiferType string
	}{
		{"x", token.INT},
		{"foobar", token.INT},
		{"diss", token.INT},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testIntStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testIntStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "int" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	intStmt, ok := s.(*ast.IntStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if intStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, intStmt.Name.Value)
		return false
	}

	if intStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, intStmt.Name)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	// 	input := `
	// 	return 4;
	// 	return dataVar;
	// 	return 4*3+2;
	//    `

	// 	failingInput := `
	// return int;
	// return string;
	// return -=/;
	//    `
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
		{"return 4*3+2", "(4*3)+2"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {

			t.Fatalf("program.Statements does not contain 1 statement. got=%d %s",
				len(program.Statements), program.Statements)
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		// checking if return value is correct is still remaining!
		// we should check if it is a valid return, it should be int,string,bool value or literal, and it shouldn't be a keyword
		// if !testReturnValue(t, returnStmt.ReturnValue, tt.expectedValue) {
		// 	return
		// }
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn Add(int x, string y) int { 
		int z = x + y
		return z
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	// fmt.Println(program)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fmt.Println(stmt)
	if stmt.FnName != "Add" {
		t.Fatalf("function name should have been Add got `%s`", stmt.FnName)
	}
	if stmt.ReturnType.Literal != "int" {
		t.Fatalf("function should return int type got `%s`", stmt.ReturnType.Literal)
	}
}

// this is still incomplete
func TestIfStatement(t *testing.T) {
	input := `
	if true	{
		int a = 5;
	} else if 5 < 9 {
		int b = 4;
	} else if 9 > 3 {
		int z = 50;
	} else {
		int u = 2;
		int z = 333;
	}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	fmt.Println(program)
}
