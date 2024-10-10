package parser

import (
	"fmt"
	"strings"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

func parseCommand(p *Parser) *ast.Command {
	if p.hasNoTokens() {
		p.errors.add(fmt.Errorf("Expected a command"), lex.Token{})
		return nil
	}
	if p.current().Type == lex.TokenCommand &&
		strings.ToLower(p.current().Value) == "add" {
		return parseAddCommand(p)
	}
	p.errors.add(fmt.Errorf("Unknown command"), *p.current())
	return nil
}

func parseAddCommand(p *Parser) *ast.Command {
	p.consume()
	if p.hasNoTokens() {
		p.errors.add(fmt.Errorf("Expected a param"), lex.Token{})
		return nil
	}
	if !p.expectCurrent(lex.TokenString) {
		return nil
	}
	var param = parseParam(p)
	if param == nil {
		return nil
	}
	var options = parseStatments(p)
	return &ast.Command{
		Kind:    ast.CommandKindAdd,
		Param:   param,
		Options: options,
	}
}

func parseExpressionStatements(p *Parser) []*ast.ExpressionStatement {
	var statements []*ast.ExpressionStatement
	for p.current().Type != lex.TokenEOF {
		statements = append(statements, parseExpressionStatement(p))
	}
	return statements
}

func parseExpressionStatement(p *Parser) *ast.ExpressionStatement {
	return &ast.ExpressionStatement{Expr: parseExpression(p, BP_DEFAULT)}
}

func parseParam(p *Parser) *ast.Param {
	if p.hasNoTokens() {
		p.errors.add(fmt.Errorf("Expected a param"), lex.Token{})
		return nil
	}
	if p.current().Type == lex.TokenString {
		return &ast.Param{
			Kind:  ast.ParamTypeDescription,
			Value: p.consume().Value,
		}
	}
	if p.current().Type == lex.TokenNumber {
		return &ast.Param{
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
	var statements []ast.Statement

	for {
		fmt.Printf(
			"parseStatments: %v\n",
			p.current().String(),
		)
		if p.hasNoTokens() || p.current().Type == lex.TokenEOF {
			break
		}
		statements = append(statements, parseStatment(p))
	}

	return statements
}

func parseStatment(p *Parser) ast.Statement {
	var handler, exists = StatementTable[p.current().Type]
	if exists {
		return handler(p)
	}
	var statement = parseExpressionStatement(p)
	return statement
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
