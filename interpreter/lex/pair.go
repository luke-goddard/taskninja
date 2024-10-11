package lex

import "github.com/luke-goddard/taskninja/interpreter/token"

func lexPair(l *Lexer) StateFn {
	if l.peek() != ':' {
		panic("lexPair called without a colon")
	}
	l.emit(token.Key)
	return lexStart
}
