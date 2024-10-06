package lex

import (
	"unicode/utf8"
)

// A lexer is a function that processes the input and returns the next state
type StateFn func(*Lexer) StateFn

// Maximum depth of the lexer, if reached, the lexer will stop processing
const MAX_DEPTH_DEFAULT = 1000

// Used to represent the end of the input
const EOF = -1


// A lexer is used to tokenize a string into a series of tokens
type Lexer struct {
	Items       chan *Token // channel of scanned items
	input       string      // the string being scanned
	line        int         // current line number
	start       Pos         // start position of this item
	position    Pos         // current position in the input
	initalState StateFn     // initial state of the lexer
	tokens      []Token     // list of tokens
	depth       int         // current depth of the lexer
	maxDepth    int         // maximum depth of the lexer
}

// Create a new lexer that will tokenize the given input
func NewLexer(input string) *Lexer {
	return &Lexer{
		Items:       make(chan *Token, 200),
		input:       input,
		line:        1,
		start:       0,
		position:    0,
		initalState: lexStart,
		tokens:      []Token{},
		depth:       0,
		maxDepth:    MAX_DEPTH_DEFAULT,
	}
}

// Useful for parsing partial input
func (l *Lexer) SetInitialState(state StateFn) {
	l.initalState = state
}

// Tokenize the input string and return the tokens
func (l *Lexer) Tokenize() []Token {
	for state := l.initalState; state != nil; {
		l.depth++
		if l.depth >= l.maxDepth {
			l.emitError("Max depth reached while parsing input")
			break
		}
		state = state(l)
	}
	close(l.Items)
	return l.tokens
}

// Get the next rune in the input, incrementing the position
// and line number as necessary
// Returns EOF if the end of the input has been reached
func (l *Lexer) next() rune {
	if int(l.position) >= len(l.input) {
		return EOF
	}
	var rune, size = utf8.DecodeRuneInString(l.input[l.position:])
	l.position += Pos(size)
	if IsNewLine(rune) {
		l.line++
	}
	return rune
}

// Emit a token with the given type and value
func (l *Lexer) emit(tokenType TokenType) {
	var token = l.toToken(tokenType)
	l.Items <- token
	l.start = l.position
	l.tokens = append(l.tokens, *token)
}

// Convert the current input into a token with the given type
func (l *Lexer) toToken(tokenType TokenType) *Token {
	return NewToken(
		tokenType,
		l.start,
		l.position,
		l.line,
		l.input[l.start:l.position],
	)
}

// Emit an error token with the given message
func (l *Lexer) emitError(error string) StateFn {
	l.Items <- NewToken(
		TokenError,
		l.start,
		l.position,
		l.line,
		error,
	)
	return nil
}
