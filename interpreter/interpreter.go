package interpreter

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/parser"
	"github.com/luke-goddard/taskninja/interpreter/token"
	"github.com/luke-goddard/taskninja/interpreter/transpiler"
	"github.com/sanity-io/litter"
)

type Interpreter struct {
	input      string
	lexer      *lex.Lexer
	parser     *parser.Parser
	transpiler *transpiler.Transpiler
	errs       *manager.ErrorManager
}

func NewInterpreter() *Interpreter {
	var manager = manager.NewErrorManager()
	return &Interpreter{
		lexer:      lex.NewLexer(manager),
		parser:     parser.NewParser(manager),
		transpiler: transpiler.NewTranspiler(manager),
	}
}

func (interpreter *Interpreter) Execute(input string) {
	interpreter.input = input

	var tokens []token.Token
	var cmd *ast.Command
	var sql transpiler.SqlStatement
	var args transpiler.SqlArgs
	var errs []manager.ErrorTranspiler

	tokens, errs = interpreter.lexer.
		Reset().
		SetInput(input).
		Tokenize()

	if len(errs) > 0 {
		fmt.Println(errs)
		return
	}

	cmd, errs = interpreter.parser.
		Reset().
		Parse(tokens)

	if len(errs) > 0 {
		fmt.Println(errs)
		return
	}

	sql, args, errs = interpreter.transpiler.Reset().Transpile(cmd)
	fmt.Println(sql)
	fmt.Println(args)
	litter.Dump(cmd)
}
