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
}

func NewInterpreter(input string) *Interpreter {
	return &Interpreter{
    input: input,
		lexer: lex.NewLexer(input),
	}
}

func (i *Interpreter) Execute() {
  fmt.Printf("executing: %s\n", i.input)
	var tokens = i.lexer.Tokenize()
	for _, token := range tokens {
		fmt.Println(token.String())
	}
  var parser = parser.NewParser(tokens)
  var command = parser.Parse()
  litter.Dump(command)
}
