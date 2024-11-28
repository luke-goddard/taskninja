package lex

import "github.com/luke-goddard/taskninja/interpreter/token"

// +HOME or -HOME
func lexTag(l *Lexer) StateFn {
	l.next()
	var peek = l.peek()

	if !IsAlphaNumeric(peek) || peek == EOF {
		l.backup()
		return lexOperator
	}

	if IsNumber(peek) {
		return lexNumber
	}

	if IsAlphaNumeric(peek) {
		l.readUntil(func(r rune) bool {
			return !IsAlphaNumeric(r) && r != '-'
		})
		l.emit(token.Tag)
		return lexStart
	}

	return lexOperator
}
