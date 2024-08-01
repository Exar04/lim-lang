package lexer

import (
	"fmt"
	"lim-lang/token"
	"strings"
	"testing"
)

func TestNextToken(t *testing.T) {

	nextTokenVariableInitilization(t)
	nextTokenIfStatement(t)
	nextTokenFn(t)
	nextTokenOperators(t)

	CheckIllegalTok(t)
}

// This method is not working properly maybe instead of checking for illegal token errors in lexer i should just check them while parsing
func CheckIllegalTok(t *testing.T) {
	input := `int data;
	int data = 35;
	hehehheheh\
	noss .
	miso .
	`
	l := New(input)
	for l.position < len(l.input) {
		tok := l.NextToken()
		if tok.Type == token.ILLEGAL {

			// fmt.Println()
			fmt.Println(strings.ReplaceAll(l.ReadErrorLine(), "\t", " "))
			// fmt.Println(strings.ReplaceAll(l.errLine, "\t", " "))
			npl := l.illegalTokenPosition
			// npl := l.CaluclateIllegleTokenPosition()
			for npl > 0 {
				fmt.Print("-")
				npl -= 1
			}
			fmt.Print("^\n")

			fmt.Println("Error : ", tok.Literal, tok.Type)
		}
	}
}

func nextTokenVariableInitilization(t *testing.T) {
	input := `int data;
	int data = 52;
	bool data = false;	
	string data = "thisIsStr";
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// int data;
		{token.Keywork_INT, "int"},
		{token.IDENT, "data"},
		{token.SEMICOLON, ";"},

		// int data = 52;
		{token.Keywork_INT, "int"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.INT, "52"},
		{token.SEMICOLON, ";"},

		// bool data = false;
		{token.Keywork_BOOL, "bool"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},

		// string data = "thisIsStr";
		{token.Keywork_STRING, "string"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.STRING, "thisIsStr"},
		{token.SEMICOLON, ";"},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func nextTokenIfStatement(t *testing.T) {
	input := `
	if 5 < 10 {
		return true
	} else if{
		return true 
 	}else{
		return false
	}
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.RBRACE, "}"},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func nextTokenFn(t *testing.T) {
	input := `
	fn printNum(){
		int data = 54;
		print(data);
	}
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.FUNCTION, "fn"},
		{token.IDENT, "printNum"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.Keywork_INT, "int"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.INT, "54"},
		{token.SEMICOLON, ";"},
		{token.PRINT, "print"},
		{token.LPAREN, "("},
		{token.IDENT, "data"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func nextTokenOperators(t *testing.T) {
	input := `
	!-/*5
   	5 < 10 > 5

	10 <= 10
	10 >= 10
	10 == 10
	10 != 9
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},

		{token.INT, "10"},
		{token.LTEQ, "<="},
		{token.INT, "10"},

		{token.INT, "10"},
		{token.GTEQ, ">="},
		{token.INT, "10"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func nextTokenConst(t *testing.T) {
	input := `
	const int data = 4;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CONST, "const"},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func nextTokenPointers(t *testing.T) {
	input := `
	&data;
	*data;
	string na = "yash"
	string *sa  = &na
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
