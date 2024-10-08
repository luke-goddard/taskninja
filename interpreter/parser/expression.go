package parser

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

func parsePrimaryExpression(p *Parser) ast.Expression {
	fmt.Printf("parse primary expression\n")
	switch p.current().Type {
	case lex.TokenString:
		return &ast.Literal{
			Kind:  ast.LiteralKindString,
			Value: p.consume().Value,
		}
	case lex.TokenNumber:
		return &ast.Literal{
			Kind:  ast.LiteralKindNumber,
			Value: p.consume().Value,
		}
	case lex.TokenWord:
		return &ast.Literal{
			Kind:  ast.LiteralKindString,
			Value: p.consume().Value,
		}
	case lex.TokenKey:

		var k = p.consume().Value
		p.consume() // Get past colon
		return &ast.Key{Key: k}
	}
	var current = p.current()
	var err = fmt.Errorf("Unknown primary expression: %s", current.String())
	panic(err)
}

func parseExpression(p *Parser, bp BindingPower) ast.Expression {
	var tokenKind = p.current().Type
	var current = p.current()
	var nudHandler, exists = NudTable[tokenKind]
	fmt.Printf(
		"parse expression kind: %s value: %s\n",
		current.String(),
		p.current().Value,
	)
	if !exists {
		var current = p.current()
		var err = fmt.Errorf("Nud handler does not exist for token: %s", current.String())
		panic(err)
	}

	var left = nudHandler(p)

	for BindingPowerTable[p.current().Type] > bp {
		var current = p.current()
		fmt.Printf(
			"Current: %v, %d > %d\n",
			current.String(),
			BindingPowerTable[current.Type],
			bp,
		)
		var tokenKind = p.current().Type
		var ledHandler, exists = LedTable[tokenKind]
		if !exists {
			var current = p.current()
			var err = fmt.Errorf("Missing led handler Unknown token: %s", current.String())
			panic(err)
		}
		left = ledHandler(p, left, bp)
	}
	return left
}

func parseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	// var operator = p.current().Type
	p.consume()
	var right = parseExpression(p, bp)
	return &ast.BinaryExpression{
		Operator: ast.BinaryOperatorEq, // TODO: THIS NEEDS TO BE WORKED OUT
		Left:     left,
		Right:    right,
	}
}

func parseGroupedExpression(p *Parser) ast.Expression {
	fmt.Printf("parse grouped expression\n")

	fmt.Printf("Current: %v\n", p.current().String())
	p.consume() // Get past left paren

	fmt.Printf("Current: %v\n", p.current().String())
	var expression = parseExpression(p, BP_DEFAULT)

	fmt.Printf("Current: %v\n", p.current().String())
	p.consume() // Get past right paren

	return expression
}
