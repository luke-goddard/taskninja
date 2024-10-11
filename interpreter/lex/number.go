package lex

import "github.com/luke-goddard/taskninja/interpreter/token"

func lexNumber(l *Lexer) StateFn {
	l.readUntil(func(r rune) bool {
		return !IsNumber(r) && r != '.'
	})

	l.emit(token.Number)

	return lexStart
}
