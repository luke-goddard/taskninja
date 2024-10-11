package lex

import "github.com/luke-goddard/taskninja/interpreter/token"

func lexStart(l *Lexer) StateFn {
	l.skipWhitespace()

	var peek = l.peek()
	if peek == EOF {
		l.emit(token.Eof)
		return nil
	}

	if peek == ':' {
		l.next()
		l.emit(token.Colon)
		return lexStart
	}

	if peek == '=' {
		l.next()
		l.emit(token.Equal)
		return lexStart
	}

	if peek == '(' {
		l.next()
		l.emit(token.LeftParen)
		return lexStart
	}

	if peek == ')' {
		l.next()
		l.emit(token.RightParen)
		return lexStart
	}

	if peek == '<' {
		l.next()
		l.emit(token.LessThan)
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
