package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	IDENT   = "IDENT"
	INTEGER = "INTEGER"
	EQUALS  = "="
	LPAREN  = "("
	RPAREN  = ")"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)
