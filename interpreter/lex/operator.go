package lex

func lexOperator(l *Lexer) StateFn {
	var opp = l.next()
	var peek = l.peek()
	if peek == EOF {
		l.emit(TokenOperator)
		return nil
	}

	if IsWhitespace(peek) {
		l.emit(TokenOperator)
		return lexStart
	}

	if opp == '-' && IsNumber(peek) {
		return lexNumber
	}

	// +test TODO: Lex tags here
	return lexWord
}
