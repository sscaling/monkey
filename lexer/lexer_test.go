package lexer

import (
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

	input := `let x = 10`
	l := New(input)

	e := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.EOF, ""},
	}

	for i, tt := range e {
		//		fmt.Printf("item %d, expected '%v'", i, e)

		tok := l.NextToken()
		//	fmt.Printf(", token: '%v'\n", tt)
		if tt.expectedLiteral != tok.Literal {
			t.Fatalf("Unexpected literal '%v' found, expected '%v' at position %d\n", tok.Literal, tt.expectedLiteral, i)
		}

		if tt.expectedType != tok.Type {
			t.Fatalf("Unexpected type '%v' found, expected '%v' at position %d\n", tok.Type, tt.expectedType, i)
		}
	}
}
