package lexer

import (
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
	case '(' == l.ch:
		t = newToken(token.LPAREN, l.ch)
	case ')' == l.ch:
		t = newToken(token.RPAREN, l.ch)
	case '=' == l.ch:
		t = newToken(token.EQUALS, l.ch)
	case isInteger(l.ch):
		return l.readInteger()
	case isLetter(l.ch):
		// return immediately as readIdentifier has already moved onto the next position
		return l.readIdentifier()
	case 0 == l.ch:
		t.Literal = ""
		t.Type = token.EOF
	default:
		t = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()

	return t
}

func (l *Lexer) readIdentifier() token.Token {
	ident := make([]byte, 0)
	for isLetter(l.ch) {
		ident = append(ident, l.ch)
		l.readChar()
	}

	return token.Token{Type: token.IDENT, Literal: string(ident)}
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (l *Lexer) readInteger() token.Token {
	integer := make([]byte, 0)
	for isInteger(l.ch) {
		integer = append(integer, l.ch)
		l.readChar()
	}

	return token.Token{Type: token.INTEGER, Literal: string(integer)}
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
