package parser

import (
	"fmt"
	"limLang/ast"
	"limLang/lexer"
	"limLang/token"
	"testing"
)

func TestIntStatements(t *testing.T) {
	input := `
	int x = 5;
	int foobar = 838383;
	int diss = -8 * -(4 + 4);
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
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

func TestBoolStatements(t *testing.T) {
	input := `
	bool y = true
	bool x = false 
	bool a = !false 
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	fmt.Println(program.Statements)
	if program == nil {
		t.Fatal("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	// tests := []struct {
	// 	expectedIdentifier    string
	// 	expectedIdentiferType string
	// }{
	// 	{"y", token.BOOL},
	// }

	// for i, tt := range tests {
	// 	stmt := program.Statements[i]
	// }
}

func TestStringStatements(t *testing.T) {
	input := `
	string a = "data"
	string b = "data 2";
	string c = "l2"
	string diss = "Hello"+"World";
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	fmt.Println(program.Statements)
	if program == nil {
		t.Fatal("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	// tests := []struct {
	// 	expectedIdentifier    string
	// 	expectedIdentiferType string
	// }{
	// 	{"y", token.BOOL},
	// }

	// for i, tt := range tests {
	// 	stmt := program.Statements[i]
	// }
}

func TestDefineStatements(t *testing.T) {
	input := `
	x := 32 + 3 * 2
	y := true
	str := "yash"
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	fmt.Println(program.Statements)
	if program == nil {
		t.Fatal("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	// tests := []struct {
	// 	expectedIdentifier    string
	// 	expectedIdentiferType string
	// }{
	// 	{"x", token.INT},
	// 	{"y", token.BOOL},
	// 	{"str", token.STRING},
	// }

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
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar", "foobar"},
		{"return 4*3+2", "(4*3)+2"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
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
func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// [...]
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"(3 < 5 == true) == false",
			"(((3 < 5) == true) == false)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		if len(program.Statements) != 1 {

			t.Fatalf("program.Statements does not contain 1 statement. got=%d %s",
				len(program.Statements), program.Statements)
		}

		stmt := program.Statements[0]
		expstm, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if !(expstm.String() == tt.expected) {
			t.Fatal("")
		}
	}
}

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
		int z = 333*3+2;
	}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	fmt.Println(program)
}
func TestTernaryOperatorStatement(t *testing.T) {

}
func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn Add(int x, string y) int { 
		int z = x + y
		return z
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
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

func TestCallExpression(t *testing.T) {
	input := `
	add(3*3+2,Ani)
	sub(ali(b), bli(a))
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != 2 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			2, len(program.Statements))
	}

	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fmt.Println(stmt1)
	fmt.Println(stmt2)
}

func TestParsingArray(t *testing.T) {
	// input := "int arr = [1, 2 * 2, 3 + 3]"
	// input := "int []arr = {1, 2 * 2, 3 + 3}"
	// l := lexer.New(input)
	// p := New(l)
	// program := p.ParseProgram()
	// stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	// array, ok := stmt.Expression.(*ast.ArrayLiteral)
	// if !ok {
	// 	t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	// }
	// if len(array.Elements) != 3 {
	// 	t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	// }
}

// structs are not yet supported
func TestStruct(t *testing.T) {
	// input := `
	// struct student {
	// 	name string
	// 	age int
	// 	gender bool
	// }
	// `

	// l := lexer.New(input)
	// p := New(l)
	// program := p.ParseProgram()
	// if len(program.Statements) != 1 {
	// 	t.Fatalf("program.Body does not contain %d statements. got=%d\n",
	// 		1, len(program.Statements))
	// }

	// stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)

	// if !ok {
	// 	t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	// }
	// fmt.Println(stmt1)
}
