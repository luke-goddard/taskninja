package lex

func lexNumber(l *Lexer) StateFn {
	l.readUntil(func(r rune) bool {
		return !IsNumber(r) && r != '.'
	})

	l.emit(TokenNumber)

	return lexStart
}
