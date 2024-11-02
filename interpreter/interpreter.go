// Converts the query into a SQL statement
package interpreter

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/parser"
	"github.com/luke-goddard/taskninja/interpreter/semantic"
	"github.com/luke-goddard/taskninja/interpreter/token"
	"github.com/rs/zerolog/log"
)

type Interpreter struct {
	input      string
	lexer      *lex.Lexer
	parser     *parser.Parser
	semantic   *semantic.Analyzer
	transpiler *ast.Transpiler
	errs       *manager.ErrorManager
	lastCmd    *ast.Command
}

func NewInterpreter() *Interpreter {
	var manager = manager.NewErrorManager()
	return &Interpreter{
		lexer:      lex.NewLexer(manager),
		parser:     parser.NewParser(manager),
		semantic:   semantic.NewAnalyzer(manager),
		transpiler: ast.NewTranspiler(),
		errs:       manager,
	}
}

func (interpreter *Interpreter) Reset() *Interpreter {
	interpreter.input = ""
	interpreter.lexer.Reset()
	interpreter.parser.Reset()
	interpreter.semantic.Reset()
	interpreter.transpiler.Reset()
	interpreter.errs.Reset()
	return interpreter
}

func (interpreter *Interpreter) Lex(input string) ([]token.Token, []manager.ErrorTranspiler) {
	interpreter.Reset()
	interpreter.input = input
	return interpreter.lexer.
		Reset().
		SetInput(input).
		Tokenize()
}

func (interpreter *Interpreter) Parse(tokens []token.Token) (*ast.Command, []manager.ErrorTranspiler) {
	interpreter.Reset()
	return interpreter.parser.
		Reset().
		Parse(tokens)
}

func (interpreter *Interpreter) ParserString(input string) (*ast.Command, []manager.ErrorTranspiler) {
	interpreter.Reset()
	tokens, errs := interpreter.Lex(input)
	if len(errs) > 0 {
		return nil, errs
	}
	return interpreter.Parse(tokens)
}

func (interpreter *Interpreter) GetLastCmd() *ast.Command {
	return interpreter.lastCmd
}

func (interpreter *Interpreter) Execute(input string) (ast.SqlStatement, ast.SqlArgs, error) {
	interpreter.input = input
	interpreter.lastCmd = nil

	var tokens []token.Token
	var cmd *ast.Command
	var sql ast.SqlStatement
	var args ast.SqlArgs
	var errs []manager.ErrorTranspiler

	tokens, errs = interpreter.lexer.
		Reset().
		SetInput(input).
		Tokenize()

	if len(errs) > 0 {
		return "", nil, fmt.Errorf("failed to tokenize input")
	}

	cmd, errs = interpreter.parser.
		Reset().
		Parse(tokens)

	if len(errs) > 0 {
		var err = errs[0]
		return "", nil, &err
	}

	var semErr = interpreter.semantic.Analyze(cmd)
	if semErr != nil {
		return "", nil, semErr
	}

	var tranErrors []ast.TranspileError
	sql, args, tranErrors = interpreter.transpiler.Reset().Transpile(cmd)
	if len(tranErrors) > 0 {
		var err = fmt.Errorf("failed to transpile input: %v", tranErrors)
		return "", nil, err
	}

	log.Info().
		Str("sql", string(sql)).
		Interface("args", args).
		Msg("transpiled sql statement")

	interpreter.lastCmd = cmd
	return sql, args, nil
}
