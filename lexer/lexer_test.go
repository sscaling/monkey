package lexer

import (
	"fmt"
	"testing"

	"github.com/sscaling/monkey/token"
)

func TestBasicLex(t *testing.T) {

	l := New("()")

	fmt.Printf("%v\n", l)
	tt := l.NextToken()
	fmt.Printf("Token %v\n", tt)

	if tt.Type != token.LPAREN {
		t.FailNow()
	}

	tt = l.NextToken()
	fmt.Printf("Token %v\n", tt)
	if tt.Type != token.RPAREN {
		t.FailNow()
	}
}
