package parser

import (
	"fmt"

	"github.com/sscaling/monkey/ast"
	"github.com/sscaling/monkey/lexer"
	"github.com/sscaling/monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	pos int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Pretty())
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	p.pos++
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	currentPos := 0
	for !p.curTokenIs(token.EOF) {

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		//		fmt.Printf("p.peekToken.Type ? '%#v'. current pos %d vs p.pos %d\n", p.peekToken, currentPos, p.pos)
		if currentPos == p.pos {
			msg := fmt.Sprintf("Failed to advance tokens. Parsing must have failed! Current Token %#v, Peek Token %#v\n", p.curToken, p.peekToken)
			p.errors = append(p.errors, msg)
			return nil
		}

		currentPos = p.pos
	}

	return program
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) parseStatement() ast.Statement {
	//	fmt.Printf("parseStatement:: %+v\n", p)

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	}

	return nil
}

// let <ident> = <expr>;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	s := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	s.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.eatUntilSemiColon()

	p.nextToken()

	return s
}

// return <expr>;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	s := &ast.ReturnStatement{Token: p.curToken}

	p.eatUntilSemiColon()

	p.nextToken()

	return s
}

// FIXME: implement values
func (p *Parser) eatUntilSemiColon() {
	for !p.curTokenIs(token.SEMI_COLON) {
		p.nextToken()
	}
}
