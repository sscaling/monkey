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
	for p.peekToken.Type != token.EOF {

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

func (p *Parser) parseStatement() ast.Statement {
	//	fmt.Printf("parseStatement:: %+v\n", p)

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

	//fmt.Println("parsed identifier")
	//fmt.Printf("parseLetStatement:: %+v\n", p)
	s.Name = identifier

	if p.peekToken.Type == token.ASSIGN {
		p.nextToken()
		//fmt.Println("parsed ASSIGN")
	} else {
		// error
		return nil
	}

	// FIXME: Value implementation
	for p.curToken.Type != token.SEMI_COLON {
		p.nextToken()
	}

	p.nextToken()

	//fmt.Printf("Parsed to semi-colon %+v\n", p)
	//fmt.Printf("Built statement %#v\n", s)
	//s.Value = p.parseExpression()

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
