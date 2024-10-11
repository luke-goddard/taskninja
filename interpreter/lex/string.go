package lex

import "github.com/luke-goddard/taskninja/interpreter/token"

func lexString(l *Lexer) StateFn {
	var delim = l.next()
	l.ignore()
	l.readUntil(func(r rune) bool {
		return r == delim && l.prev() != '\\'
	})
	l.emit(token.String)
	l.next()
	l.ignore()
	return lexStart
}
