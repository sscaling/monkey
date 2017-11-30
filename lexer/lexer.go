package lexer

import "github.com/sscaling/monkey/token"

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

func (l *Lexer) NextToken() *token.Token {

	var t *token.Token

	switch {
	case 0 == l.ch:
		// Return immediately, do not continue processing
		return newToken(token.EOF, l.ch)
	case '(' == l.ch:
		t = newToken(token.LPAREN, l.ch)
	case ')' == l.ch:
		t = newToken(token.RPAREN, l.ch)
	default:
		t = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()

	return t
}

func newToken(t token.TokenType, literal byte) *token.Token {
	return &token.Token{Type: t, Literal: string(literal)}
}
