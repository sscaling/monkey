package lexer

import (
	"fmt"
	"testing"

	"github.com/sscaling/monkey/token"
)

func TestIllegalChar(t *testing.T) {
	tt := New("^").NextToken().Type
	if tt != token.ILLEGAL {
		t.Fatalf("Expected illegal token for '^'")
	}
}

func TestBasicLex(t *testing.T) {

	input := "()"
	l := New(input)

	expected := []token.TokenType{
		token.LPAREN,
		token.RPAREN,
		token.EOF,
	}

	for i, e := range expected {
		fmt.Printf("item %d, expected '%v'", i, e)

		tt := l.NextToken().Type
		fmt.Printf(", token: '%v'\n", tt)
		if tt != e {
			t.Fatalf("Unexpected token '%v' found at position %d\n", tt, i)
		}
	}
}
