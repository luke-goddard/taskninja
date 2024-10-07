package lex

import "fmt"

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

	// +test TODO: Lex tags here
	return lexWord
}

func getOperator(opp rune) TokenType {
	switch opp {
	case '+':
		return TokenPlus
	case '-':
		return TokenMinus
	case '*':
		return TokenStar
	case '/':
		return TokenSlash
	case '<':
		return TokenLT
	}
	var e = fmt.Errorf("Unknown operator: %d", opp)
	panic(e)
}
