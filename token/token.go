package token

import "fmt"

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	Position int
	Line     int
	Column   int
}

func (t Token) String() string {
	return string(t.Type)
}

func (t Token) Pretty() string {
	return fmt.Sprintf("%s '%s' [%d:%d]", t.Type, t.Literal, t.Line, t.Column)
}

const (
	IDENT   = "IDENT" // foo, bar etc
	INTEGER = "INTEGER"

	ASSIGN       = "="
	EQUALS       = "=="
	BANG         = "!"
	NOT_EQUAL    = "!="
	PLUS         = "+"
	MINUS        = "-"
	MULTIPLY     = "*"
	DIVIDE       = "/"
	LESS_THAN    = "<"
	GREATER_THAN = ">"

	COMMA      = ","
	SEMI_COLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdentifier(identifier string) TokenType {
	tok, ok := keywords[identifier]
	if ok {
		return tok
	}

	return IDENT
}
