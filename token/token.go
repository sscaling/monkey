package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	LPAREN  = "("
	RPAREN  = ")"
	EOF     = ""
	ILLEGAL = ""
)
