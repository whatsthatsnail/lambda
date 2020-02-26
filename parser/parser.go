package parser

import (
	"fmt"
	"lambda/ast"
	"lambda/lexer"
)

// ---------- Parser type: ---------- //

type parser struct {
	tokens     []lexer.Token
	statements []ast.Statement
	current    int
	errFlag    bool
}

// Parser constructor, initializes default vaules
func NewParser(tokens []lexer.Token) parser {
	p := parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

// ---------- Node creator methods: ---------- //

func (p *parser) Parse() ([]ast.Statement, bool) {
	for !p.isAtEnd() {
		p.statements = append(p.statements, p.statement())
	}

	if p.tokens[0].TType == lexer.EOF {
		p.errFlag = true
	}

	return p.statements, p.errFlag
}

func (p *parser) statement() ast.Statement {
	var stmt interface{}

	switch p.peek().TType {
	case lexer.LET:
		p.advance()
		stmt = p.definition()
	default:
		stmt = p.term()
	}

	p.consume(lexer.NEWLINE, "Expect newline after statment.")

	return stmt.(ast.Statement)
}

func (p *parser) definition() ast.Statement {
	id := ast.Identifier{p.advance().Lexeme}

	p.consume(lexer.EQUAL, "Expect '=' after definition identifier.")

	term := p.term()

	return ast.Definition{id, term}
}

func (p *parser) term() ast.Term {
	if p.peek().TType == lexer.LAMBDA {
		p.advance()

		param, _ := p.atom()

		p.consume(lexer.DOT, "Expect '.' after function parameter.")
		body := p.term()

		// Return the abstraction.
		return ast.Abstraction{param, body}
	}

	return p.application()
}

func (p *parser) application() ast.Term {
	left, _ := p.atom()

	right, ok := p.atom()
	for ok {
		left = ast.Application{left, right}
		right, ok = p.atom()
	}

	return left
}

func (p *parser) atom() (ast.Term, bool) {
	if p.peek().TType == lexer.LEFT_PAREN {
		p.advance()
		term := p.term()
		p.consume(lexer.RIGHT_PAREN, "Expect closing ')' after term.")
		return term, true
	} else if p.peek().TType == lexer.IDENTIFIER {
		id := ast.Identifier{p.advance().Lexeme}
		if p.peek().TType == lexer.EQUAL {
			p.consume(lexer.EQUAL, "Expect '=' after implicit definition.")
			return ast.Definition{id, p.term()}, false
		}
		return id, true
	}

	return ast.Abstraction{}, false
}

// ---------- Helper methods: ---------- //

// Return current token without advancing.
func (p *parser) peek() lexer.Token {
	return p.tokens[p.current]
}

// Check if the current position is the last token (an EOF token).
func (p *parser) isAtEnd() bool {
	return p.peek().TType == lexer.EOF
}

// Return the token directly before the current position.
func (p *parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

// Advance the current position and return the current token.
func (p *parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
		return p.previous()
	} else {
		return p.tokens[p.current]
	}
}

// Compare the type of the current token to a given TokenType.
func (p *parser) check(tType lexer.TokenType) bool {
	if !p.isAtEnd() {
		return p.peek().TType == tType
	} else {
		return false
	}
}

// Advance over a specified token type. Throw an error if the actual token doesn't match.
func (p *parser) consume(tType lexer.TokenType, message string) {
	if p.check(tType) {
		p.advance()
	} else {
		p.parseError(p.peek(), message)
	}
}

// Print an error message at the current line.
func (p *parser) parseError(token lexer.Token, message string) {
	p.errFlag = true
	fmt.Printf("[Line: %v] %s\n", token.Line, fmt.Sprintf("%v -- %s", p.peek().Lexeme, message))
}
