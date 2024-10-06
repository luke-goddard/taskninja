package lex

import (
	"unicode/utf8"
)

type StateFn func(*Lexer) StateFn

const MAX_DEPTH_DEFAULT = 1000

type Lexer struct {
	Items       chan *Token // channel of scanned items
	input       string      // the string being scanned
	line        int         // current line number
	start       Pos         // start position of this item
	position    Pos         // current position in the input
	initalState StateFn
	tokens      []Token
	depth       int
	maxDepth    int
}

const EOF = -1

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

func (l *Lexer) Tokenize() []Token {
	for state := l.initalState; state != nil; {
    l.depth ++
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
