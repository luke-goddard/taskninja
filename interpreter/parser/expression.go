package parser

import (
	"fmt"
	"runtime/debug"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

func parsePrimaryExpression(p *Parser) ast.Expression {
	switch p.current().Type {
	case token.String:
		return &ast.Literal{
			Kind:  ast.LiteralKindString,
			Value: p.consume().Value,
		}
	case token.Number:
		return &ast.Literal{
			Kind:  ast.LiteralKindNumber,
			Value: p.consume().Value,
		}
	case token.Key:
		var k = p.consume().Value
		p.consume() // Get past colon
		return &ast.Key{Key: k}
	}
	var current = p.current()
	var err = fmt.Errorf("Unknown primary expression: %s", current.String())
	panic(err)
}

func parseExpression(p *Parser, bp BindingPower) ast.Expression {
	if p.hasNoTokens() || p.current().Type == token.Eof {
		return nil
	}
	var tokenKind = p.current().Type
	var nudHandler, exists = NudTable[tokenKind]
	if !exists {
		var current = p.current()
		debug.PrintStack()
		var err = fmt.Errorf("Nud handler does not exist for token: %s", current.String())
		p.errors.add(err, *current)
		return nil
	}

	var left = nudHandler(p)

	for BindingPowerTable[p.current().Type] > bp {
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
	if !p.expectOneOf(
		token.Equal, token.LessThan, token.Plus,
		token.Minus, token.Star, token.Slash,
	) {
		return nil
	}
	var op = p.consume()
	var binop ast.BinaryOperator
	switch op.Type {
	case token.Equal:
		binop = ast.BinaryOperatorEq
	case token.LessThan:
		binop = ast.BinaryOperatorLt
	case token.Plus:
		binop = ast.BinaryOperatorAdd
	case token.Minus:
		binop = ast.BinaryOperatorSub
	case token.Star:
		binop = ast.BinaryOperatorMul
	default:
		var err = fmt.Errorf("Unknown binary operator: %s", op.String())
		panic(err)
	}

	var right = parseExpression(p, bp)
	if right == nil {
		return nil
	}
	return &ast.BinaryExpression{
		Operator: binop,
		Left:     left,
		Right:    right,
	}
}

func parseGroupedExpression(p *Parser) ast.Expression {
	p.consume() // Get past left paren
	var expression = parseExpression(p, BP_DEFAULT)
	if expression == nil {
		return nil
	}
	p.consume() // Get past right paren
	return expression
}
