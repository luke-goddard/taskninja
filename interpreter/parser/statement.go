package parser

import (
	"strings"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

func parseCommand(parser *Parser) *ast.Command {
	if parser.hasNoTokens() {
		parser.errors.EmitParse("no tokens to parse", &token.Token{})
		return nil
	}
	if parser.current().Type == token.Command &&
		strings.ToLower(parser.current().Value) == "add" {
		return parseAddCommand(parser)
	}

	if parser.current().Type == token.Command &&
		strings.ToLower(parser.current().Value) == "list" {
		return parseListCommand(parser)
	}

	parser.errors.EmitParse("Unknown command", parser.current())
	return nil
}

func parseAddCommand(parser *Parser) *ast.Command {
	parser.consume()
	if parser.hasNoTokens() {
		parser.errors.EmitParse("Expected a param", &token.Token{})
		return nil
	}
	if !parser.expectCurrent(token.String) {
		return nil
	}
	var param = parseParam(parser)
	if param == nil {
		return nil
	}
	var options = parseStatments(parser)
	return &ast.Command{
		Kind:    ast.CommandKindAdd,
		Param:   param,
		Options: options,
	}
}

func parseListCommand(parser *Parser) *ast.Command {
	parser.consume()
	var options = parseStatments(parser)
	return &ast.Command{
		Kind:    ast.CommandKindList,
		Options: options,
	}
}

func parseExpressionStatement(parser *Parser) *ast.ExpressionStatement {
	var expression = parseExpression(parser, BP_DEFAULT)
	if expression == nil {
		parser.errors.EmitParse("Expected an expression", parser.current())
		return nil
	}
	return &ast.ExpressionStatement{Expr: expression}
}

func parseParam(parser *Parser) *ast.Param {
	if parser.hasNoTokens() {
		parser.errors.EmitParse("Expected a param", &token.Token{})
		return nil
	}
	if parser.current().Type == token.String {
		return &ast.Param{
			Kind:  ast.ParamTypeDescription,
			Value: parser.consume().Value,
		}
	}
	if parser.current().Type == token.Number {
		return &ast.Param{
			Kind:  ast.ParamTypeTaskId,
			Value: parser.consume().Value,
		}
	}
	panic("Unknown param")
}

func parseTagDecStatement(parser *Parser) ast.Statement {
	var value = parser.current().Value
	if len(value) == 0 {
		panic("Expected tag")
	}

	var op ast.TagOperator

	if value[0] == '+' {
		op = ast.TagOperatorPlus
	} else if value[0] == '-' {
		op = ast.TagOperatorMinus
	} else {
		panic("Expected tag operator")
	}

	parser.consume()
	return &ast.Tag{
		Operator: op,
		Value:    value[1:],
	}
}

func parseStatments(parser *Parser) []ast.Statement {
	var statements []ast.Statement

	for {
		if parser.hasNoTokens() {
			break
		}
		var statement = parseStatment(parser)
		if statement == nil {
			break
		}
		statements = append(statements, statement)
	}

	return statements
}

func parseStatment(parser *Parser) ast.Statement {
	var handler, exists = StatementTable[parser.current().Type]
	if exists {
		return handler(parser)
	}
	var statement = parseExpressionStatement(parser)
	if statement == nil {
		parser.errors.EmitParse("Expected a statement", parser.current())
		return nil
	}
	return statement
}
