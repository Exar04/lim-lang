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
	IDENT          = "IDENT"         // add, foobar, x, y, ...
	Keywork_INT    = "Keywork_INT"   // 1343456
	Keywork_FLOAT  = "Keywork_FLOAT" // 12.34
	Keywork_BOOL   = "Keywork_BOOL"  // true, false
	Keywork_STRING = "Keywork_STRING"

	INT    = "INT"
	FLOAT  = "FLOAT"
	BOOL   = "BOOL"
	STRING = "STRING"

	// Value_INT    = "Value_INT"
	// Value_FLOAT  = "Value_FLOAT"
	// Value_STRING = "Value_STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	LTEQ   = "<="
	GTEQ   = ">="
	EQ     = "=="
	NOT_EQ = "!="

	ADD_ASSIGN = "+="
	SUB_ASSIGN = "-="
	MUL_ASSIGN = "*="
	QUO_ASSIGN = "/="
	REM_ASSIGN = "%="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	// IFELSE   = "IFELSE"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	NULL   = "NULL"

	PRINT = "PRINT"
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

	"int":    Keywork_INT,
	"bool":   Keywork_BOOL,
	"float":  Keywork_FLOAT,
	"string": Keywork_STRING,

	// "int":    INT,
	// "bool":   BOOL,
	// "float":  FLOAT,
	// "string": STRING,
	"print": PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
