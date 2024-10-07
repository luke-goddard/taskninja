package parser

import (
	"strings"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

func parseCommand(p *Parser) *ast.Command {
	if p.current().Type == lex.TokenCommand &&
		strings.ToLower(p.current().Value) == "add" {
		return parseAddCommand(p)
	}
	return nil
}

func parseAddCommand(p *Parser) *ast.Command {
	p.consume()
	return &ast.Command{
		Kind:    ast.CommandKindAdd,
		Param:   parseParam(p),
		Options: parseExpressionStatements(p),
	}
}

func parseExpressionStatements(p *Parser) []ast.ExpressionStatement {
	var statements []ast.ExpressionStatement
	for p.current().Type != lex.TokenEOF {
		statements = append(statements, parseExpressionStatement(p))
	}
	return statements
}

func parseExpressionStatement(p *Parser) ast.ExpressionStatement {
	var expression = parseExpression(p, BP_DEFAULT)
	return ast.ExpressionStatement{
		Expression: expression,
	}
}

func parseParam(p *Parser) ast.Param {
	if p.current().Type == lex.TokenString {
		return ast.Param{
			Kind:  ast.ParamTypeDescription,
			Value: p.consume().Value,
		}
	}
	if p.current().Type == lex.TokenNumber {
		return ast.Param{
			Kind:  ast.ParamTypeTaskId,
			Value: p.consume().Value,
		}
	}
	panic("Unknown param")
}

func parseCommandKind(p *Parser) ast.CommandKind {
	if p.current().Type != lex.TokenCommand {
		panic("Expected command")
	}
	p.consume()
	switch strings.ToLower(p.current().Value) {
	case "add":
		return ast.CommandKindAdd
	}
	panic("Unknown command")
}
