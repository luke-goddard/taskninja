package parser

import (
	"fmt"
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

func parsePrimaryExpression(parser *Parser) ast.Expression {
	switch parser.current().Type {
	case token.String:
		return &ast.Literal{
			Kind:  ast.LiteralKindString,
			Value: parser.consume().Value,
		}
	case token.Number:
		return &ast.Literal{
			Kind:  ast.LiteralKindNumber,
			Value: parser.consume().Value,
		}
	case token.Key:
		var key = parser.consume().Value
		if !parser.expectCurrent(token.Colon){
			panic("TODO err handle")
		}
		parser.consume()
		return &ast.Key{
			Key: key,
			Expr: parseExpression(parser, BP_PRIMARY),
		}
	}
	var current = parser.current()
	var err = fmt.Errorf("Unknown primary expression: %s", current.String())
	panic(err)
}

func parseExpression(parser *Parser, bp BindingPower) ast.Expression {
	if parser.hasNoTokens() || parser.current().Type == token.Eof {
		return nil
	}
	var tokenKind = parser.current().Type
	var nudHandler, exists = NudTable[tokenKind]
	if !exists {
		var current = parser.current()
		var message = fmt.Sprintf("Nud handler does not exist for token: %s", current.String())
		parser.errors.EmitParse(message, current)
		return nil
	}

	var left = nudHandler(parser)

	for BindingPowerTable[parser.current().Type] > bp {
		var tokenKind = parser.current().Type
		var ledHandler, exists = LedTable[tokenKind]
		if !exists {
			var current = parser.current()
			var message = fmt.Sprintf("Missing led handler Unknown token: %s", current.String())
			parser.errors.EmitParse(message, current)
			return nil
		}
		left = ledHandler(parser, left, bp)
	}
	return left
}

func parseBinaryExpression(parser *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	if !parser.expectOneOf(
		token.Equal, token.LessThan, token.Plus,
		token.Minus, token.Star, token.Slash,
	) {
		return nil
	}
	var op = parser.consume()
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

	var right = parseExpression(parser, bp)
	if right == nil {
		return nil
	}
	return &ast.BinaryExpression{
		Operator: binop,
		Left:     left,
		Right:    right,
	}
}

func parseGroupedExpression(parser *Parser) ast.Expression {
	parser.consume() // Get past left paren
	var expression = parseExpression(parser, BP_DEFAULT)
	if expression == nil {
		return nil
	}
	parser.consume() // Get past right paren
	return expression
}
