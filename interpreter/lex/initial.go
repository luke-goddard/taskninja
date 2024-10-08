package lex

func lexStart(l *Lexer) StateFn {
	l.skipWhitespace()

	var peek = l.peek()
	if peek == EOF {
		l.emit(TokenEOF)
		return nil
	}

	if peek == ':' {
		l.next()
		l.emit(TokenColon)
		return lexStart
	}

	if peek == '=' {
		l.next()
		l.emit(TokenEQ)
		return lexStart
	}

	if peek == '(' {
		l.next()
		l.emit(TokenLeftParen)
		return lexStart
	}

	if peek == ')' {
		l.next()
		l.emit(TokenRightParen)
		return lexStart
	}

	if peek == '<' {
		l.next()
		l.emit(TokenLT)
		return lexStart
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

	if IsAlphaNumeric(peek) {
		return lexCommand
	}

	return l.errorf("Unknown character: %c", peek)
}
