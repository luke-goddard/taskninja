package lex

import "fmt"

func lexNumber(l *Lexer) StateFn {
	var last = l.readUntil(func(r rune) bool {
		return !IsNumber(r) && r != '.'
	})
	fmt.Printf("last %c\n", last)

		l.emit(TokenNumber)

	return lexStart
}
