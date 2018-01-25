package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/sscaling/monkey/ast"
	"github.com/sscaling/monkey/lexer"
	"github.com/sscaling/monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	pos int

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) Summary() string {
	return fmt.Sprintf("Current: %s {%s}, Peek: %s {%s} [%d:%d]\n", p.curToken, p.curToken.Literal, p.peekToken, p.peekToken.Literal, p.curToken.Line, p.curToken.Column)
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.EQUALS, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

var precedences = map[token.TokenType]int{
	token.EQUALS:       EQUALS,
	token.NOT_EQUAL:    EQUALS,
	token.LESS_THAN:    LESSGREATER,
	token.GREATER_THAN: LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.DIVIDE:       PRODUCT,
	token.MULTIPLY:     PRODUCT,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Pretty())
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
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

		//		fmt.Printf("\nParsed statement\n%#v\n\n", stmt)

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

		// should this be here or after currentPos assignment
		p.nextToken()
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

var debug bool = false
var indent int = 0

func tabs() string {
	return strings.Join(make([]string, indent*4), " ")
}

func (p *Parser) parseStatement() ast.Statement {
	if debug {
		fmt.Printf("%sparseStatement:: %s\n", tabs(), p.Summary())
	}

	indent++
	defer func() { indent-- }()

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}

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

	return s
}

// return <expr>;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	s := &ast.ReturnStatement{Token: p.curToken}

	p.eatUntilSemiColon()

	return s
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if debug {
		fmt.Printf("%sparseExpressionStatement %s\n", tabs(), p.Summary())
	}

	s := &ast.ExpressionStatement{Token: p.curToken}

	indent++
	s.Expression = p.parseExpression(LOWEST)
	indent--

	//	fmt.Printf("Parsed expression '%#v'.\nParser state: '%#v'\n", s.Expression, p)

	if p.peekTokenIs(token.SEMI_COLON) {
		// make semi-colon next token
		p.nextToken()

		// FIXME:
		p.eatUntilSemiColon()
	}

	return s
}

func fnName(f interface{}) string {
	//	 github.com/sscaling/monkey/parser.(*Parser).(github.com/sscaling/monkey/parser.parsePrefixExpression)-fm
	n := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	p := regexp.MustCompile(`^.*\.(.*)\).*$`)
	return p.ReplaceAllString(n, "$1")
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if debug {
		fmt.Printf("%sparseExpression(%d) : %s\n", tabs(), precedence, p.Summary())
	}

	prefix := p.prefixParseFns[p.curToken.Type]

	if debug {
		fmt.Printf("%sFound prefix fn = %s.(%T)\n", tabs(), fnName(prefix), prefix)
	}

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	indent++
	leftExp := prefix()
	indent--

	for !p.peekTokenIs(token.SEMI_COLON) && precedence < p.peekPrecedence() {
		//		fmt.Printf("  traversing tokens looking for infix operations, curToken %v\n", p.curToken)
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		indent++
		leftExp = infix(leftExp)
		indent--
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	if debug {
		fmt.Printf("%sFound identifier %q\n", tabs(), p.curToken.Literal)
	}
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	if debug {
		fmt.Printf("%sFound integer %q\n", tabs(), p.curToken.Literal)
	}
	e := &ast.IntegerLiteral{Token: p.curToken}

	// if base == 0, the prefix of the string determines the base (i.e. 0x etc)
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	e.Value = value

	return e
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	if debug {
		fmt.Printf("%sparsePrefixExpression %s\n", tabs(), p.Summary())
	}

	p.nextToken()

	//	fmt.Printf("parsed prefix expression. Parser: %v\n", p)

	indent++
	expr.Right = p.parseExpression(PREFIX)
	indent--

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()

	if debug {
		fmt.Printf("%sparseInfixExpression %s\n", tabs(), p.Summary())
	}

	p.nextToken()

	indent++
	expr.Right = p.parseExpression(precedence)
	indent--

	return expr
}

// FIXME: implement values
func (p *Parser) eatUntilSemiColon() {
	for !p.curTokenIs(token.SEMI_COLON) {
		p.nextToken()
	}
}
