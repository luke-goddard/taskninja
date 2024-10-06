package lex

// The sequence is specific, and must follow these rules:
//   - date < duration < uuid < identifier
//   - dom < uuid
//   - uuid < hex < number
//   - url < pair < identifier
//   - hex < number
//   - separator < tag < operator
//   - path < substitution < pattern
//   - set < number
//   - word last

func lexStart(l *Lexer) StateFn {
  l.depth++
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
