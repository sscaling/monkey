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

	pos int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
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
			fmt.Errorf("Failed to advance tokens. Parsing must have failed! Current Token %#v, Peek Token %#v\n", p.curToken, p.peekToken)
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

	fmt.Errorf("Expected '%v', got '%v'\n", t, p.peekToken.Type)
	return false
}

func (p *Parser) parseStatement() ast.Statement {
	//	fmt.Printf("parseStatement:: %+v\n", p)

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	}

	return nil
}

// let <ident> = <expr>;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	s := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		fmt.Errorf("Expected token.IDENT, got '%v'\n", p.peekToken.Type)
		return nil
	}

	s.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// FIXME: Value implementation
	for !p.curTokenIs(token.SEMI_COLON) {
		p.nextToken()
	}

	p.nextToken()

	return s
}

func (p *Parser) parseExpression() ast.Expression {
	// what expressions?

	return nil
}
