package evaluator

import (
	"lim-lang/lexer"
	"lim-lang/object"
	"lim-lang/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2", 16},
		{"5 + 2 * 10", 25},
		{"3 * (3+3) + 10", 28},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got =%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got =%d, expected=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1<2", true},
		{"1>2", false},
		{"1<1", false},
		{"1>1", false},
		{"1==2", false},
		{"1==1", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangBoolOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}
func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true {  10 }", 10},
		{"if false { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 } else if true { 20 }", 20},
		{"if 1 > 2 { 10 } else if false { 20 } else if true { 30 }", 30},
		{"if 1 < 2 { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if 10 > 1 {
			if 10 > 1 {
				return 10 
			}
			return 1
		}`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

// func TestErrorHandling(t *testing.T) {
// 	tests := []struct {
// 		input           string
// 		expectedMessage string
// 	}{
// 		{
// 			"5 + true;",
// 			"type mismatch: INTEGER + BOOLEAN",
// 		},
// 		{
// 			"5 + true; 5;",
// 			"type mismatch: INTEGER + BOOLEAN",
// 		},
// 		{
// 			"-true",
// 			"unknown operator: -BOOLEAN",
// 		},
// 		{
// 			"true + false;",
// 			"unknown operator: BOOLEAN + BOOLEAN",
// 		},
// 		{
// 			"true + false + true + false;",
// 			"unknown operator: BOOLEAN + BOOLEAN",
// 		},
// 		{
// 			"5; true + false; 5",
// 			"unknown operator: BOOLEAN + BOOLEAN",
// 		},
// 		{
// 			"if (10 > 1) { true + false; }",
// 			"unknown operator: BOOLEAN + BOOLEAN",
// 		},
// 		{
// 			`
// 	if (10 > 1) {
// 	  if (10 > 1) {
// 		return true + false;
// 	  }

// 	  return 1;
// 	}
// 	`,
// 			"unknown operator: BOOLEAN + BOOLEAN",
// 		},
// 		{
// 			"foobar",
// 			"identifier not found: foobar",
// 		},
// 	}

// 	for _, tt := range tests {
// 		evaluated := testEval(tt.input)

// 		errObj, ok := evaluated.(*object.Error)
// 		if !ok {
// 			t.Errorf("no error object returned. got=%T(%+v)",
// 				evaluated, evaluated)
// 			continue
// 		}

// 		if errObj.Message != tt.expectedMessage {
// 			t.Errorf("wrong error message. expected=%q, got=%q",
// 				tt.expectedMessage, errObj.Message)
// 		}
// 	}
// }

func TestIntStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"int a = 5; a;", 5},
		{"int a = 5+1; a", 6},
		{"int a = 5; int b = a; b;", 5},
		{"int a = 5; int b = a; int c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn Add(int x) int { 
		x+2
	}`
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "int x" {
		t.Fatalf("parameter is not 'int x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "( x + 2)\n"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"fn identity(int x) { x; }; identity(5);", 5},
		{"fn identity(int x) { return x; }; identity(5);", 5},
		{"fn doTwoX(x) { x * 2; }; doTowX(5);", 10},
		{"fn add(int x, int y) { x + y; }; add(5, 5);", 10},
		{"fn add(int x, int y) { x + y; }; add(5 + 5, add(5, 5));", 20},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
