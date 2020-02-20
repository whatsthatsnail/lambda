package parser

import (
	"fmt"
	"lambda/ast"
	"lambda/lexer"
)

// ---------- Parser type: ---------- //

type parser struct {
	tokens  []lexer.Token
	current int
	errFlag bool
}

// Parser constructor, initializes default vaules
func NewParser(tokens []lexer.Token) parser {
	p := parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

//// Helper methods:
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

func (p *parser) consume(tType lexer.TokenType, message string) {
	if p.check(tType) {
		p.advance()
	} else {
		p.parseError(p.peek(), message)
	}
}

func (p *parser) parseError(token lexer.Token, message string) {
	p.errFlag = true
	fmt.Printf("[Line: %v] %s\n", token.Line, message)
}

//// Node creator methods:

func (p *parser) Parse() (ast.Term, bool) {
	return p.term(), p.errFlag
}

func (p *parser) term() ast.Term {
	if p.peek().TType == lexer.LAMBDA {
		p.advance()
		param := ast.Identifier{p.advance()}
		p.consume(lexer.DOT, "Expect '.' after function parameter.")
		body := p.term()

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
		term := ast.Identifier{p.advance()}
		return term, true
	}

	return ast.Abstraction{}, false
}