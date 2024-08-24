package lexer

import (
	"bytes"
	"fmt"
	"limLang/token"
	"os"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte

	CurrentLineNumber  int
	StartOfCurrentLine int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, CurrentLineNumber: 0, StartOfCurrentLine: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1

	if l.ch == '\n' {
		l.CurrentLineNumber += 1
		l.StartOfCurrentLine = l.readPosition + 1
	}

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	// skips single line comment
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' {
			l.readChar()
		}
	}
	// skips multi line comment
	if l.ch == '/' && l.peekChar() == '*' {
		for {
			if l.ch == '*' && l.peekChar() == '/' {
				break
			}
			l.readChar()
		}
		l.readChar()
		l.readChar()
		fmt.Println(l.ch)
	}
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.ADD_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.SUB_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.MUL_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.QUO_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.REM_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.MODULUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEFINE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.COLON, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BITWISE_OR, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BITWISE_AND, l.ch)
		}
	case '.':
		tok = newToken(token.PERIOD, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '\n':
		tok = newToken(token.ENDOFLINE, l.ch)
	case '"':
		l.readChar()
		tok.Literal = l.readString()
		tok.Type = token.STRING
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.realNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
			l.errBreaker("Got illigal token")
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
func (l *Lexer) realNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if isLetter(l.ch) {
		l.errBreaker("A number shouldn't be as a first character in a identifier name")
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readString() string {
	var str string
	for l.ch != '"' {
		str += string(l.ch)
		l.readChar()
	}
	return str
}
func (l *Lexer) skipInlineNMultiLineComment() {
	// var str string
	for l.ch != '*' && l.peekChar() != '/' {
		// str += string(l.ch)
		l.readChar()
	}
	l.readChar()
	// return str
}

func (l *Lexer) readSingleLineComment() string {
	var str string
	for l.ch != '\n' {
		str += string(l.ch)
		l.readChar()
	}
	return str
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) errBreaker(errMsg string) {
	fmt.Println("Error occured at line number :", l.CurrentLineNumber)
	fmt.Println(errMsg)
	var errLine bytes.Buffer
	start := l.StartOfCurrentLine - 1
	currChar := l.input[start]
	for currChar == '\t' || currChar == ' ' {
		start += 1
		currChar = l.input[start]
	}
	nStart := start
	for currChar != '\n' {
		errLine.WriteString(string(currChar))
		start += 1
		currChar = l.input[start]
	}
	fmt.Println(errLine.String())

	for nStart < l.readPosition-1 {
		fmt.Print("-")
		nStart += 1
	}
	fmt.Print("^")

	// right now we are just breaking the program here but later
	// we should return a illigal token to the parser
	// and parser will end program gracefully instead of doing it abruptly like we are doing it rn
	os.Exit(1)
}
