package lex

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/token"
)

func lexOperator(l *Lexer) StateFn {
	var opp = l.next()
	var peek = l.peek()
	if peek == EOF {
		l.emit(getOperator(opp))
		return nil
	}

	if IsWhitespace(peek) {
		l.emit(getOperator(opp))
		return lexStart
	}

	if opp == '-' && IsNumber(peek) {
		return lexNumber
	}

	if IsNumber(peek) {
		l.emit(getOperator(opp))
		return lexStart
	}

	// +test TODO: Lex tags here
	return lexWord
}

func getOperator(opp rune) token.TokenType {
	switch opp {
	case '+':
		return token.Plus
	case '-':
		return token.Minus
	case '*':
		return token.Star
	case '/':
		return token.Slash
	}
	var e = fmt.Errorf("Unknown operator: %c", opp)
	panic(e)
}
