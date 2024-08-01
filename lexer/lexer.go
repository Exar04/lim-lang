package lexer

import (
	"fmt"
	"lim-lang/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte

	// CurrentLine       int
	// PreviousLineStart int
	// CurrentLineStart int
	// PosG              int
	// errLine              string
	AtCharNumFromCurrentLine int
	illegalTokenPosition     int
	// tokenErrors []tokenErr // this will store all the illegal tokens and their line numbers { lineNum, tokenName, line }
}

// type tokenErr struct {
// 	lineNum   int
// 	line      string
// 	tokenName string
// }

func New(intput string) *Lexer {
	l := &Lexer{input: intput}
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
	// l.PosG = l.position //////// L take
	l.readPosition += 1
	if l.ch != '\t' {
		l.AtCharNumFromCurrentLine += 1
	}

	if l.ch == '\n' {
		l.AtCharNumFromCurrentLine = 0
		l.ch = 10
		// fmt.Println(l.ch)
	}
}

func (l *Lexer) ReadErrorLine() string {
	cln := l.position - l.AtCharNumFromCurrentLine + 1
	var errStr string
	fmt.Println(cln)
	for l.input[cln] != '\n' {
		// l.errLine += string(l.input[cln])
		errStr += string(l.input[cln])

		cln += 1

		if cln == len(l.input) {
			return errStr
		}
	}
	return errStr
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '\n':
		tok = token.Token{Type: token.ENDOFLINE, Literal: string(l.ch)}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
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
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	// case '\'':
	// 	fmt.Println("for handling chars")
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		l.readChar()
		// tok = newToken(token.Value_STRING, l.readString())
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
			// l.ReadErrorLine()
			// l.illegalTokenPosition = l.position - l.CurrentLineStart
			l.illegalTokenPosition = l.AtCharNumFromCurrentLine
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) realNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
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
	for isLetter(l.ch) {
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

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
