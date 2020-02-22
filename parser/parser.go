package parser

import (
	"fmt"
	"lambda/ast"
	"lambda/lexer"
)

// ---------- Parser type: ---------- //

type parser struct {
	tokens      []lexer.Token
	context     []string
	current     int
	errFlag     bool
}

// Parser constructor, initializes default vaules
func NewParser(tokens []lexer.Token) parser {
	p := parser{}
	p.tokens = tokens
	p.current = 0

	return p
}

// ---------- Node creator methods: ---------- //

func (p *parser) Parse() (ast.Term, bool) {
	ast, err := p.term(), p.errFlag

	return ast, err
}

func (p *parser) term() ast.Term {
	if p.peek().TType == lexer.LAMBDA {
		p.advance()

		// Creating a new identifier places it at the end of the context stack.
		param := p.advance().Lexeme

		// Push the identifier onto the context stack.
		p.context = append(p.context, param)

		p.consume(lexer.DOT, "Expect '.' after function parameter.")
		body := p.term()

		// Return the abstraction.
		return ast.Abstraction{param, body}
	}

	return p.application()
}

func (p *parser) application() ast.Term {
	left := p.application()

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

		// Lookup the identifier in the context stack. Attach it's De Bruijn index to the object.
		index, ok := contextIndex(p.peek(), p.context)

		if ok {
			// Attach the identifier's distance from it's declaration in the context stack
			term := ast.Identifier{p.advance(), len(p.context) - index, false}
			return term, true
		} else {
			// If it's not in the context stack, it's a free variables.
			// Free variables wrapped in n lambdas are given index n (using p.lamdaDepth)
			term := ast.Identifier{p.advance(), -1, true}

			// Do not push free variables to the context stack? Sure.

			return term, true
		}

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

// Return the index of an identifer in a slice. Used for De Bruijn Index calculation.
func contextIndex(obj lexer.Token, slice []string) (int, bool) {
	for k, v := range slice {
		if obj.Lexeme == v {
			return k + 1, true
		}
	}

	return -1, false
}
