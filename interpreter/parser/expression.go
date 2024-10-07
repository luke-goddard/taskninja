package parser

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

func parsePrimaryExpression(p *Parser) ast.Expression {
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
	case lex.TokenTag:
		return parseTag(p)

	case lex.TokenKey:
		return parsePair(p)
	}
	var current = p.current()
	var err = fmt.Errorf("Unknown primary expression: %s", current.String())
	panic(err)
}

func parseExpression(p *Parser, bp BindingPower) ast.Expression {
	var tokenKind = p.current().Type
	var nudHandler, exists = NudTable[tokenKind]
	if !exists {
		var current = p.current()
		var err = fmt.Errorf("Nud handler does not exist for token: %s", current.String())
		panic(err)
	}

	var left = nudHandler(p)

	// We don't have a semicolon in the language,
	// so we have to check if the expression is a tag
	if left.Type() == ast.NodeTypeTag {
		return left
	}

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

func parseTag(p *Parser) ast.Expression {
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

func parsePair(p *Parser) ast.Expression {
	var key = p.consume().Value
	if p.current().Type != lex.TokenColon {
		panic("Expected colon")
	}
	p.consume()
	return &ast.Pair{
		Key:   key,
		Value: parseExpressionStatement(p),
	}
}
