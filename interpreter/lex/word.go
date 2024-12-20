package lex

import "github.com/luke-goddard/taskninja/interpreter/token"

// This is the default if everything else failse
func lexWord(l *Lexer) StateFn {
	for {
		var r = l.next()
		if r == EOF {
			break
		}
		if IsWhitespace(r) || !IsAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	l.emit(token.String)
	return lexStart
}
