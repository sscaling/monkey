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
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMI_COLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMI_COLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMI_COLON, ";"},

		{token.NOT, "!"},
		{token.MINUS, "-"},
		{token.DIVIDE, "/"},
		{token.MULTIPLY, "*"},
		{token.INTEGER, "5"},
		{token.SEMI_COLON, ";"},

		{token.INTEGER, "5"},
		{token.LESS_THAN, "<"},
		{token.INTEGER, "10"},
		{token.GREATER_THAN, ">"},
		{token.INTEGER, "5"},
		{token.SEMI_COLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INTEGER, "5"},
		{token.LESS_THAN, "<"},
		{token.INTEGER, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMI_COLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMI_COLON, ";"},
		{token.RBRACE, "}"},

		{token.INTEGER, "10"},
		{token.EQUALS, "=="},
		{token.INTEGER, "10"},
		{token.SEMI_COLON, ";"},

		{token.INTEGER, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INTEGER, "9"},
		{token.SEMI_COLON, ";"},

		{token.EOF, ""},
	}

	for _, tt := range e {
		//		fmt.Printf("item %d, expected '%v'", i, e)

		tok := l.NextToken()
		//	fmt.Printf(", token: '%v'\n", tt)
		if tt.expectedLiteral != tok.Literal {
			t.Fatalf("Unexpected literal '%v' found, expected '%v' at position %d:%d\n", tok.Literal, tt.expectedLiteral, tok.Line, tok.Column)
		}

		if tt.expectedType != tok.Type {
			t.Fatalf("Unexpected type '%v' found, expected '%v' at position %d:%d\n", tok.Type, tt.expectedType, tok.Line, tok.Column)
		}
	}
}
