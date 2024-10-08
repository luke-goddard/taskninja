package parser

import (
	"fmt"
	"strings"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

func parseCommand(p *Parser) *ast.Command {
	fmt.Printf("parseing command\n")
	if p.current().Type == lex.TokenCommand &&
		strings.ToLower(p.current().Value) == "add" {
		return parseAddCommand(p)
	}
  p.errors.add(fmt.Errorf("Unknown command"), *p.current())
  return nil
}

func parseAddCommand(p *Parser) *ast.Command {
	fmt.Printf("parse add command\n")
	p.consume()
	return &ast.Command{
		Kind:    ast.CommandKindAdd,
		Param:   parseParam(p),
		Options: parseStatments(p),
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
		Expr: expression,
	}
}

func parseParam(p *Parser) ast.Param {
	fmt.Printf("parse param\n")
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

func parseTagDecStatement(p *Parser) ast.Statement {
	var value = p.current().Value
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

	p.consume()
	return &ast.Tag{
		Operator: op,
		Value:    value[1:],
	}
}

func parseStatments(p *Parser) []ast.Statement {
	fmt.Printf("parse statements\n")
	var statements []ast.Statement
	for p.current().Type != lex.TokenEOF {
		statements = append(statements, parseStatment(p))
	}
	return statements
}

func parseStatment(p *Parser) ast.Statement {
	var current = p.current()
	fmt.Printf("parse statement: %s\n", current.String())
	var handler, exists = StatementTable[p.current().Type]
	if exists {
		fmt.Printf("handler exists for type: %s\n", current.String())
		return handler(p)
	}
	var statement = parseExpressionStatement(p)
	return &statement
}

func parsePairDeclStatement(p *Parser) ast.Statement {
	var key = p.consume().Value
	if p.current().Type != lex.TokenColon {
		panic("Expected colon")
	}
	p.consume()
	return &ast.Key{
		Key: key,
	}
}
