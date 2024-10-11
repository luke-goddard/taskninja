package interpreter

import (
	"github.com/luke-goddard/taskninja/interpreter/lex"
	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/parser"
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
		lexer: lex.NewLexer(manager),
	}
}

func (interpreter *Interpreter) Execute(input string) {
	interpreter.input = input
	var tokens = interpreter.lexer.
		Reset().
		SetInput(input).
		Tokenize()

	interpreter.parser.Parse(tokens)
}

