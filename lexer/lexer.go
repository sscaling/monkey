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

func (l *Lexer) NextToken() token.Token {

	t := token.Token{}

	switch {
	case '(' == l.ch:
		t = token.Token{Type: token.LPAREN}
	case ')' == l.ch:
		t = token.Token{Type: token.RPAREN}
	}

	l.readChar()

	return t
}
