package interpreter

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/lex"
)

type Interpreter struct {
	input string
	lexer *lex.Lexer
}

func NewInterpreter(input string) *Interpreter {
	return &Interpreter{
		input: input,
		lexer: lex.NewLexer(input),
	}
}

func (i *Interpreter) Execute() {
	fmt.Printf("executing: %s\n", i.input)
	var tokens = make([]lex.Token, 0)
	go i.lexer.Tokenize()
	for {
		var token = <-i.lexer.Items
		tokens = append(tokens, *token)
		if token == nil || token.Type == lex.TokenEOF || token.Type == lex.TokenError {
			break
		}
	}
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	// var parser = parser.NewParser(tokens)
	// var command = parser.Parse()
	// litter.Dump(command)
}
