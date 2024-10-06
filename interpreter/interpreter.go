package interpreter

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/lex"
)

type Interpreter struct {
	lexer *lex.Lexer
}

func NewInterpreter(input string) *Interpreter {
	return &Interpreter{
		lexer: lex.NewLexer(input),
	}
}

func (i *Interpreter) Execute() {
	var tokens = i.lexer.Tokenize()
	for _, token := range tokens {
		fmt.Println(token.String())
	}
}
