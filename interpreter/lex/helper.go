package lex

import (
	"unicode"
	"unicode/utf8"

	"github.com/luke-goddard/taskninja/interpreter/token"
)

// Returns true if the rune is a new line character
func IsNewLine(r rune) bool {
	return r == '\n' || r == '\r'
}

// Returns true if the rune is a whitespace character or a new line character
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || IsNewLine(r)
}

// Returns true if the rune is a letter
func IsAlphabet(r rune) bool {
	return unicode.IsLetter(r)
}

// Returns true if the rune is a letter or a number
func IsAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

// Returns true if the rune is a number
func IsNumber(r rune) bool {
	return unicode.IsNumber(r)
}

// Get the current input
func (l *Lexer) current() string {
	return l.input[l.start:l.position]
}

// Ignore the current input
func (l *Lexer) ignore() {
	l.start = l.position
}

// Ignore the next rune in the input
func (lex *Lexer) skip() {
	lex.next()
	lex.ignore()
}

// Backup one rune
// This can be called multiple times and will not go past the beginning
// of the input. This also accounts for unicode characters with different size
func (l *Lexer) backup() {
	if l.position < 0 {
		return
	}
	var r, size = utf8.DecodeLastRuneInString(l.input[:l.position])
	l.position -= token.Pos(size)
	if IsNewLine(r) {
		l.line--
	}
}

// Peek at the next rune in the input without consuming it
func (l *Lexer) peek() rune {
	var start = l.start
	var pos = l.position
	var r = l.next()
	if l.position != pos {
		l.backup()
	}
	if pos != l.position {
		panic("peek and backup are not symmetric in position")
	}
	if start != l.start {
		panic("peek and backup are not symmetric in start")
	}
	return r
}

// Carry on scanning the input until the next none whitespace character
func (l *Lexer) skipWhitespace() {
	for {
		var r = l.next()
		if r == EOF {
			break
		}
		if !IsWhitespace(r) {
			l.backup()
			break
		}
	}
	l.ignore()
}

// Carry on scanning the input until the the predicate matches
// returns the rune that matched the predicate or EOF
func (l *Lexer) readUntil(predicate func(rune) bool) rune {
	for {
		var r = l.next()
		if r == EOF {
			return r
		}
		if predicate(r) {
			l.backup()
			return r
		}
	}
}

// Carry on scanning until we reach a whitespace character including new lines
func (l *Lexer) readUntilWhitespace() {
	l.readUntil(IsWhitespace)
}

// Get the previous rune in the input
func (l *Lexer) prev() rune {
	if l.position <= 0 {
		return EOF
	}
	var r, _ = utf8.DecodeLastRuneInString(l.input[:l.position-1])
	return r
}
