package lexer

import (
	"fmt"

	"github.com/sscaling/monkey/token"
)

type Lexer struct {
	input        string
	currPosition int
	ch           byte
}

func New(program string) *Lexer {
	l := &Lexer{input: program}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.currPosition >= len(l.input) {
		l.ch = 0
		return
	}

	l.ch = l.input[l.currPosition]
	l.currPosition += 1
}

func (l *Lexer) peakChar() byte {
	if (l.currPosition + 1) >= len(l.input) {
		return 0
	}

	return l.input[l.currPosition+1]
}

func (l *Lexer) eatWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {

	var t token.Token

	l.eatWhitespace()

	switch {
	case '=' == l.ch:
		if l.peakChar() == '=' {
			l.readChar()
			t.Type = token.EQUALS
			t.Literal = "=="
		} else {
			t = newToken(token.ASSIGN, l.ch)
		}
	case '!' == l.ch:
		if l.peakChar() == '=' {
			l.readChar()
			t.Type = token.NOT_EQUAL
			t.Literal = "!="
		} else {
			t = newToken(token.NOT, l.ch)
		}
	case '+' == l.ch:
		t = newToken(token.PLUS, l.ch)
	case '(' == l.ch:
		t = newToken(token.LPAREN, l.ch)
	case ')' == l.ch:
		t = newToken(token.RPAREN, l.ch)
	case '{' == l.ch:
		t = newToken(token.LBRACE, l.ch)
	case '}' == l.ch:
		t = newToken(token.RBRACE, l.ch)
	case '+' == l.ch:
		t = newToken(token.PLUS, l.ch)
	case '-' == l.ch:
		t = newToken(token.MINUS, l.ch)
	case '*' == l.ch:
		t = newToken(token.MULTIPLY, l.ch)
	case '/' == l.ch:
		t = newToken(token.DIVIDE, l.ch)
	case '<' == l.ch:
		t = newToken(token.LESS_THAN, l.ch)
	case '>' == l.ch:
		t = newToken(token.GREATER_THAN, l.ch)
	case ',' == l.ch:
		t = newToken(token.COMMA, l.ch)
	case ';' == l.ch:
		t = newToken(token.SEMI_COLON, l.ch)
	case 0 == l.ch:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isInteger(l.ch) {
			return l.readInteger()
		} else if isLetter(l.ch) {
			// return immediately as readIdentifier has already moved onto the next position
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdentifier(t.Literal)
			return t
		} else {
			t = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return t
}

func (l *Lexer) readIdentifier() string {
	start := l.currPosition - 1
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[start:(l.currPosition - 1)]
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (l *Lexer) readInteger() token.Token {
	start := l.currPosition - 1
	for isInteger(l.ch) {
		l.readChar()
	}

	for i, x := range l.input {
		fmt.Printf("%03d: %c\n", i, x)
	}
	return token.Token{Type: token.INTEGER, Literal: l.input[start:(l.currPosition - 1)]}
}

func isInteger(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func newToken(t token.TokenType, literal byte) token.Token {
	return token.Token{Type: t, Literal: string(literal)}
}
