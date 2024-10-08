package lex

// +HOME or -HOME
func lexTag(l *Lexer) StateFn {
	l.next()
	var peek = l.peek()

	if !IsAlphaNumeric(peek) {
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
		l.emit(TokenTag)
		return lexStart
	}

	return lexOperator
}
