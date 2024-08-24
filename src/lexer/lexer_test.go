package lexer

import (
	"limLang/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	// nextTokenVariableInitilization(t)
	nextTokenIfStatement(t)
	nextTokenFn(t)
	nextTokenOperators(t)

	nextTokenConst(t)
	nextTokenStruct(t)

	// testIlligalName()
	// testIlligalSymbol()
}

func testIlligalSymbol() {
	input := `int data;
	int data = 52;


	\bool da\ta = false\

	string data = "thisIsStr"\
	`
	l := New(input)
	tok := l.NextToken()
	for tok.Type != token.EOF {
		tok = l.NextToken()
	}
}

func testIlligalName() {
	input := `
	int data;

	int 889dal = 3;
	`
	l := New(input)
	tok := l.NextToken()
	for tok.Type != token.EOF {
		tok = l.NextToken()
	}
}

func TestVariableInitilizationTokens(t *testing.T) {
	input := `int data;
	int data = 52;
	int da42 = 2;
	bool data = false;
	string data = "thisIsStr";
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// int data;
		{token.Keyword_INT, "int"},
		{token.IDENT, "data"},
		{token.SEMICOLON, ";"},
		{token.ENDOFLINE, "\n"},

		// int data = 52;
		{token.Keyword_INT, "int"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.INT, "52"},
		{token.SEMICOLON, ";"},
		{token.ENDOFLINE, "\n"},

		// bool data = false;
		{token.Keyword_BOOL, "bool"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.ENDOFLINE, "\n"},

		// string data = "thisIsStr";
		{token.Keyword_STRING, "string"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.STRING, "thisIsStr"},
		{token.SEMICOLON, ";"},
		{token.ENDOFLINE, "\n"},
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
	input := `if 5 < 10 {
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
		{token.ENDOFLINE, "\n"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.ENDOFLINE, "\n"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.LBRACE, "{"},
		{token.ENDOFLINE, "\n"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.ENDOFLINE, "\n"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.ENDOFLINE, "\n"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.ENDOFLINE, "\n"},
		{token.RBRACE, "}"},
		{token.ENDOFLINE, "\n"},
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
	input := `fn printNum(){ int data = 54; }`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.FUNCTION, "fn"},
		{token.IDENT, "printNum"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.Keyword_INT, "int"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.INT, "54"},
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
	input := `!-/*7
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
		{token.INT, "7"},
		{token.ENDOFLINE, "\n"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.ENDOFLINE, "\n"},

		{token.INT, "10"},
		{token.LTEQ, "<="},
		{token.INT, "10"},
		{token.ENDOFLINE, "\n"},

		{token.INT, "10"},
		{token.GTEQ, ">="},
		{token.INT, "10"},
		{token.ENDOFLINE, "\n"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.ENDOFLINE, "\n"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.ENDOFLINE, "\n"},
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
	input := `const int data = 4;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CONST, "const"},
		{token.Keyword_INT, "int"},
		{token.IDENT, "data"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
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

func nextTokenStruct(t *testing.T) {
	input := `
    struct student {
        name string
        age int
        isAlive bool
    }`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ENDOFLINE, "\n"},
		{token.STRUCT, "struct"},
		{token.IDENT, "student"},
		{token.LBRACE, "{"},
		{token.ENDOFLINE, "\n"},
		{token.IDENT, "name"},
		{token.Keyword_STRING, "string"},
		{token.ENDOFLINE, "\n"},
		{token.IDENT, "age"},
		{token.Keyword_INT, "int"},
		{token.ENDOFLINE, "\n"},
		{token.IDENT, "isAlive"},
		{token.Keyword_BOOL, "bool"},
		{token.ENDOFLINE, "\n"},
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

func TestComments(t *testing.T) {
	input := `
	// this is the single line comment
	int /* Inline comment */data
	/*
	This is
	multiline comment
	*/
	int nData`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ENDOFLINE, "\n"},
		{token.ENDOFLINE, "\n"},
		{token.Keyword_INT, "int"},
		{token.IDENT, "data"},
		{token.ENDOFLINE, "\n"},
		{token.ENDOFLINE, "\n"},
		{token.Keyword_INT, "int"},
		{token.IDENT, "nData"},
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
