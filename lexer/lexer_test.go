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

	input := `let five = 5;
	let ten = 10;

	let add = fn(x, y) {
		x + y;
	};
	`
	/*
		let result = add(five, ten);
		!-/*5;
		5 < 10 > 5;

		if (5 < 10) {
			return true;
		} else {
			return false;
		}

		10 == 10;
		10 != 9;
		`
	*/
	l := New(input)

	l.Debug()

	e := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMI_COLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMI_COLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.RBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMI_COLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMI_COLON, ";"},
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
