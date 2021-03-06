package lexer

import (
	"fmt"
)

// ---------- Helper functions: ---------- //

// Determines if a character is valid as an identifier character.
func isValidChar(char rune) bool {
	switch char {
	case '\\', '.', '(', ')', '=', '\n':
		return false
	}

	return true
}

// ---------- TokenType type: ---------- //

type TokenType int

const (
	// Single-character
	LAMBDA TokenType = iota
	DOT
	LEFT_PAREN
	RIGHT_PAREN
	EQUAL

	// Keywords
	LET

	// Literals
	IDENTIFIER

	NEWLINE
	EOF
)

// TokenType lookup. Return a string interpretation of token type.
func (t TokenType) typeString() string {
	switch t {
	case 0:
		return "LAMBDA"
	case 1:
		return "DOT"
	case 2:
		return "LEFT_PAREN"
	case 3:
		return "RIGHT_PAREN"
	case 4:
		return "EQUAL"
	case 5:
		return "LET"
	case 6:
		return "IDENTIFIER"
	case 7:
		return "NEWLINE"
	case 8:
		return "EOF"
	}

	return "INVALID"
}

// ---------- Token type: ---------- //

// Literals stores as empty interface, use type assertions when parsing
type Token struct {
	TType   TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// Print an instance of a token.
func (tok Token) String() string {
	if tok.Literal == nil {
		return fmt.Sprintf("{%s, '%s', %d}", tok.TType.typeString(), tok.Lexeme, tok.Line)
	} else {
		return fmt.Sprintf("{%s, '%s', %v, %d}", tok.TType.typeString(), tok.Lexeme, tok.Literal, tok.Line)
	}
}

// Prints a list of tokens in a readable manner as {Token_Type, lexeme, (literal), line}
func PrintTokens(tokens []Token) {
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

// ---------- Lexer type: ---------- //

type lexer struct {
	start    int
	current  int
	line     int
	tokens   []Token
	source   string
	hadError bool
	repl     bool
}

// Lexer constructor, initializes default values
func NewLexer(code string, replFlag bool) lexer {
	l := lexer{}
	l.start = 0
	l.current = 0
	l.line = 1
	l.source = code
	l.hadError = false
	l.repl = replFlag

	return l
}

// Error handling:
func (l *lexer) throwError(message string) {
	fmt.Printf("[Line: %v] %s\n", l.line, message)
	l.hadError = true
}

// Checks if current position has reaced the end of the source
func (l *lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

// Consumes the current character and returns it
func (l *lexer) advance() rune {
	l.current++
	return rune(l.source[l.current-1])
}

// Peeks at next character without consuming it
func (l *lexer) peek() rune {
	if !l.isAtEnd() {
		return rune(l.source[l.current])
	} else {
		return '\n'
	}
}

// Adds a new Token instance to l.tokens using input type and literal, and infered lexeme and line
func (l *lexer) addToken(tType TokenType, literal interface{}) {
	if l.repl == true {
		l.tokens = append(l.tokens, Token{tType, l.source[l.start:l.current], literal, 0})
	} else {
		l.tokens = append(l.tokens, Token{tType, l.source[l.start:l.current], literal, l.line})
	}
}

func (l *lexer) getWord() {
	// Advance to end of word
	for isValidChar(l.peek()) && l.peek() != ' ' && !l.isAtEnd() {
		l.advance()
	}

	// Store word lexeme
	word := l.source[l.start:l.current]

	switch word {
	case "let":
		l.addToken(LET, word)
	default:
		l.addToken(IDENTIFIER, word)
	}
}

func (l *lexer) scanToken() {
	char := l.advance()

	switch char {

	// Single-character tokens
	case '\\':
		l.addToken(LAMBDA, nil)
	case '.':
		l.addToken(DOT, nil)
	case '(':
		l.addToken(LEFT_PAREN, nil)
	case ')':
		l.addToken(RIGHT_PAREN, nil)
	case '=':
		l.addToken(EQUAL, nil)

	// Line comment
	case '#':
		for l.peek() != '\n' {
			l.advance()
		}

	// Skip over whitespace
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		if len(l.tokens) > 0 {
			if l.tokens[len(l.tokens)-1].TType != NEWLINE {
				l.tokens = append(l.tokens, Token{NEWLINE, "NEWLINE", nil, l.line})
			}
		}

	// Multi-character tokens
	default:
		l.getWord()
	}
}

// Scan all tokens in a given input.
func (l *lexer) ScanTokens() ([]Token, bool) {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}

	l.tokens = append(l.tokens, Token{EOF, "EOF", nil, l.line})

	return l.tokens, l.hadError
}
