package lex

func lexStart(l *Lexer) StateFn {
	l.skipWhitespace()

	var peek = l.peek()
	if peek == EOF {
		l.emit(TokenEOF)
		return nil
	}

	if peek == '"' || peek == '\'' {
		return lexString
	}

	if peek == '+' || peek == '-' {
		return lexTag
	}

	if peek == '*' || peek == '/' {
		return lexOperator
	}

	if IsNumber(peek) {
		return lexNumber
	}

	return lexCommand
}
