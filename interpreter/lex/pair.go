package lex

func lexPair(l *Lexer) StateFn {
	if l.peek() != ':' {
		panic("lexPair called without a colon")
	}
	l.emit(TokenKey)
	return lexStart
}
