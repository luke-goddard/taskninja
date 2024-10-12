package interpreter

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/lex"
	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/parser"
	"github.com/sanity-io/litter"
)

type Interpreter struct {
	input  string
	lexer  *lex.Lexer
	parser *parser.Parser
	errs   *manager.ErrorManager
}

func NewInterpreter() *Interpreter {
	var manager = manager.NewErrorManager()
	return &Interpreter{
		lexer:  lex.NewLexer(manager),
		parser: parser.NewParser(manager),
	}
}

func (interpreter *Interpreter) Execute(input string) {
	interpreter.input = input
	var tokens, errs = interpreter.lexer.
		Reset().
		SetInput(input).
		Tokenize()

  if len(errs) > 0 {
    fmt.Println(errs)
  }

	var ast, err = interpreter.parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
	}

	litter.Dump(ast)
}
