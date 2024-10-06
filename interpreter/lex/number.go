package lex

func lexNumber(l *Lexer) StateFn {
	var last = l.readUntil(func(r rune) bool {
		return !IsNumber(r) && r != '.'
	})

	if last == EOF || last == ' ' {
		l.emit(TokenNumber)
		return lexStart
	}

	return lexStart
}
