package lex

// +HOME or -HOME
func lexTag(l *Lexer) StateFn {
	l.next()
	var peek = l.peek()

	if peek == EOF || IsWhitespace(peek) {
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
