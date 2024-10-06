package lex

func lexString(l *Lexer) StateFn {
	var delim = l.next()
	l.ignore()
	l.readUntil(func(r rune) bool {
		return r == delim && l.prev() != '\\'
	})
	l.emit(TokenString)
	l.next()
	l.ignore()
	return lexStart
}
