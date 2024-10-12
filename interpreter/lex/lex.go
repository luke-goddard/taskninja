package lex

import (
	"fmt"
	"unicode/utf8"

	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

// A lexer is a function that processes the input and returns the next state
type StateFn func(*Lexer) StateFn

// Maximum depth of the lexer, if reached, the lexer will stop processing
const MAX_DEPTH_DEFAULT = 1000

// Used to represent the end of the input
const EOF = -1

// A lexer is used to tokenize a string into a series of tokens
type Lexer struct {
	errors      *manager.ErrorManager
	input       string        // the string being scanned
	line        int           // current line number
	start       token.Pos     // start position of this item
	position    token.Pos     // current position in the input
	initalState StateFn       // initial state of the lexer
	tokens      []token.Token // list of tokens
	depth       int           // current depth of the lexer
	maxDepth    int           // maximum depth of the lexer
	seenCommand bool          // whether a command has been seen
}

// Create a new lexer that will tokenize the given input
func NewLexer(manager *manager.ErrorManager) *Lexer {
	return &Lexer{
		errors:      manager,
		line:        1,
		start:       0,
		position:    0,
		initalState: lexStart,
		tokens:      []token.Token{},
		depth:       0,
		maxDepth:    MAX_DEPTH_DEFAULT,
		seenCommand: false,
	}
}

// Useful for parsing partial input
func (l *Lexer) SetInitialState(state StateFn) {
	l.initalState = state
}

func (l *Lexer) SetInput(input string) *Lexer {
	l.input = input
	return l
}

func (l *Lexer) Reset() *Lexer {
	l.line = 1
	l.start = 0
	l.position = 0
	l.tokens = []token.Token{}
	l.depth = 0
	l.seenCommand = false
	return l
}

// Tokenize the input string and return the tokens
func (l *Lexer) Tokenize() ([]token.Token, []manager.ErrorTranspiler) {
	for state := l.initalState; state != nil; {
		l.depth++
		if l.depth >= l.maxDepth {
			l.emitError("Max depth reached while parsing input")
			break
		}
		state = state(l)
	}
	return l.tokens, l.errors.LexErrors()
}

// Get the next rune in the input, incrementing the position
// and line number as necessary
// Returns EOF if the end of the input has been reached
func (l *Lexer) next() rune {
	if int(l.position) >= len(l.input) {
		return EOF
	}
	var rune, size = utf8.DecodeRuneInString(l.input[l.position:])
	l.position += token.Pos(size)
	if IsNewLine(rune) {
		l.line++
	}
	return rune
}

// Emit a token with the given type and value
func (l *Lexer) emit(tokenType token.TokenType) {
	var token = l.toToken(tokenType)
	l.start = l.position
	l.tokens = append(l.tokens, *token)
}

// Convert the current input into a token with the given type
func (l *Lexer) toToken(tokenType token.TokenType) *token.Token {
	return token.NewToken(
		tokenType,
		l.start,
		l.position,
		l.line,
		l.input[l.start:l.position],
	)
}

// Emit an error token with the given message
func (l *Lexer) emitError(message string) StateFn {
	var token = token.NewToken(
		token.Error,
		l.start,
		l.position,
		l.line,
		message,
	)
	l.errors.EmitLex(message, token)
	return nil
}

// Returns nill, indicating that the lexer should stop processing
func (l *Lexer) errorf(format string, args ...interface{}) StateFn {
	var message = fmt.Sprintf(format, args...)
	return l.emitError(message)
}
