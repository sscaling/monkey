package lexer

import (
	"fmt"

	"github.com/sscaling/monkey/token"
)

type Lexer struct {
	input        string
	position     int // need to track both position and currPosition for slicing the input data
	line         int
	column       int
	readPosition int
	ch           byte
}

func New(program string) *Lexer {
	l := &Lexer{input: program, line: 1, column: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// if previous character was a new line, reset position counters
	if l.ch == '\n' {
		l.line += 1
		l.column = 0
	}

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]

		// for any character, other than carriage return, assume we increment one column
		// in the input. NOTE: this may not be suitable for tabs?
		if l.ch != '\r' {
			l.column += 1
		}
	}
	// Always set the position and increase the readPosition as this is used for slicing the input data
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peakChar() byte {
	if (l.readPosition) >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) eatWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {

	l.eatWhitespace()

	var t token.Token
	t.Line = l.line
	t.Column = l.column

	switch {
	case '=' == l.ch:
		if l.peakChar() == '=' {
			l.readChar()
			t.Type = token.EQUALS
			t.Literal = "=="
		} else {
			t = newToken(token.ASSIGN, l)
		}
	case '!' == l.ch:
		if l.peakChar() == '=' {
			l.readChar()
			t.Type = token.NOT_EQUAL
			t.Literal = "!="
		} else {
			t = newToken(token.BANG, l)
		}
	case '+' == l.ch:
		t = newToken(token.PLUS, l)
	case '(' == l.ch:
		t = newToken(token.LPAREN, l)
	case ')' == l.ch:
		t = newToken(token.RPAREN, l)
	case '{' == l.ch:
		t = newToken(token.LBRACE, l)
	case '}' == l.ch:
		t = newToken(token.RBRACE, l)
	case '+' == l.ch:
		t = newToken(token.PLUS, l)
	case '-' == l.ch:
		t = newToken(token.MINUS, l)
	case '*' == l.ch:
		t = newToken(token.MULTIPLY, l)
	case '/' == l.ch:
		t = newToken(token.DIVIDE, l)
	case '<' == l.ch:
		t = newToken(token.LESS_THAN, l)
	case '>' == l.ch:
		t = newToken(token.GREATER_THAN, l)
	case ',' == l.ch:
		t = newToken(token.COMMA, l)
	case ';' == l.ch:
		t = newToken(token.SEMI_COLON, l)
	case 0 == l.ch:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isInteger(l.ch) {
			t.Literal = l.readInteger()
			t.Type = token.INTEGER
			return t
		} else if isLetter(l.ch) {
			// return immediately as readIdentifier has already moved onto the next position
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdentifier(t.Literal)
			return t
		} else {
			t = newToken(token.ILLEGAL, l)
		}
	}

	l.readChar()

	return t
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (l *Lexer) readInteger() string {
	start := l.position
	for isInteger(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func isInteger(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func newToken(t token.TokenType, l *Lexer) token.Token {
	return token.Token{Type: t, Literal: string(l.ch), Line: l.line, Column: l.column}
}

func (l *Lexer) Debug() {
	charOrWhitespace := func(b rune) byte {
		if isWhitespace(byte(b)) {
			return byte(' ')
		} else {
			return byte(b)
		}
	}

	line := 1
	column := 0
	fmt.Printf("\n%05d: ", line)
	for _, x := range l.input {
		if x == '\n' {
			line += 1
			column = 0
			fmt.Printf("\n%05d: ", line)
		} else if column > 0 && column%10 == 0 {
			fmt.Printf("|%04d|", column)
		}

		fmt.Printf("%c", charOrWhitespace(x))

		column += 1
	}
}
