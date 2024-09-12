package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT          = "IDENT"          // add, foobar, x, y, ...
	Keyword_INT    = "Keyword_INT"    // 1343456
	Keyword_FLOAT  = "Keyword_FLOAT"  // 12.34
	Keyword_BOOL   = "Keyword_BOOL"   // true, false
	Keyword_STRING = "Keyword_STRING" // "lala la la la"

	INT    = "INT"
	FLOAT  = "FLOAT"
	BOOL   = "BOOL"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULUS  = "%"

	// Bitwise Operators
	BITWISE_AND = "&"
	BITWISE_OR  = "|"

	LT     = "<"
	GT     = ">"
	LTEQ   = "<="
	GTEQ   = ">="
	EQ     = "=="
	NOT_EQ = "!="
	OR     = "||"
	AND    = "&&"

	ADD_ASSIGN = "+="
	SUB_ASSIGN = "-="
	MUL_ASSIGN = "*="
	QUO_ASSIGN = "/="
	REM_ASSIGN = "%="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	PERIOD = "."
	ARROW  = "->"
	DEFINE = ":="

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACK = "["
	RBRACK = "]"

	ARRAY = "ARRAY"

	// Keywords
	FUNCTION  = "FUNCTION"
	CONST     = "CONST"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	NULL      = "NULL"
	ENDOFLINE = "ENDOFLINE"
	STRUCT    = "STRUCT"
	// BREAK     = "BREAK"
	// CONTINUE  = "CONTINUE"
	// FOR = "FOR"

	// PRINT = "PRINT"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"const":  CONST,
	"null":   NULL,

	"int":    Keyword_INT,
	"bool":   Keyword_BOOL,
	"float":  Keyword_FLOAT,
	"string": Keyword_STRING,

	"struct": STRUCT,

	// "print": PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
