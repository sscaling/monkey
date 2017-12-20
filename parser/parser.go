package parser

import (
	"github.com/sscaling/monkey/ast"
	"github.com/sscaling/monkey/lexer"
	"github.com/sscaling/monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
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
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	return nil
}

func (p *Parser) parseStatement() ast.Statement {

	// TODO: change to switch
	if p.curToken.Type == token.LET {
		return p.parseLetStatement()
	}

	return nil
}

// let <ident> = <expr>;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	s := &ast.LetStatement{Token: p.curToken}

	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}

	s.Name = identifier

	if p.peekToken.Type == token.ASSIGN {
		p.nextToken()
	} else {
		// error
		return nil
	}

	// FIXME: Value implementation
	for p.curToken.Type != token.SEMI_COLON {
		p.nextToken()
	}

	s.Value = p.parseExpression()

	return s
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	// error?

	return nil
}

func (p *Parser) parseExpression() ast.Expression {
	// what expressions?

	return nil
}
