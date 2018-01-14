package ast

import (
	"testing"

	"github.com/sscaling/monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "MyVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	s := program.String()
	if s != "let myVar = anotherVar;" {
		t.Errorf("program.String() incorrect. got = %q", s)
	}
}
