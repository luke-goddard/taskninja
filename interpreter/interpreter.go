package interpreter

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/lex"
	"github.com/luke-goddard/taskninja/interpreter/parser"
	"github.com/sanity-io/litter"
)

type Interpreter struct {
	input string
	lexer *lex.Lexer
  parser *parser.Parser
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
	var parser = parser.NewParser(tokens)
	var command, _ = parser.Parse()
	litter.Dump(command)
}
