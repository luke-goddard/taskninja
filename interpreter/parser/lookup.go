package parser

import (
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

type BindingPower int

const (
	BP_DEFAULT = iota
	BP_COMMA
	BP_ASSIGNMENT
	BP_LOGICAL
	BP_RELATIONAL
	BP_ADDITIVE
	BP_MULTIPLICATIVE
	BP_UNARY
	BP_CALL
	BP_MEMBER
	BP_PRIMARY
)

type StatementHandler func(*Parser) ast.Statement
type NudHandler func(*Parser) ast.Expression
type LedHandler func(*Parser, ast.Expression, BindingPower) ast.Expression

type StatementLookup map[token.TokenType]StatementHandler
type NudLookup map[token.TokenType]NudHandler
type LedLookup map[token.TokenType]LedHandler
type BindingPowerLookup map[token.TokenType]BindingPower

var BindingPowerTable = BindingPowerLookup{}
var NudTable = NudLookup{}
var LedTable = LedLookup{}
var StatementTable = StatementLookup{}

func led(kind token.TokenType, bp BindingPower, handler LedHandler) {
	LedTable[kind] = handler
	BindingPowerTable[kind] = bp
}

func nud(kind token.TokenType, bp BindingPower, handler NudHandler) {
	NudTable[kind] = handler
	BindingPowerTable[kind] = bp
}

func stmt(kind token.TokenType, handler StatementHandler) {
	StatementTable[kind] = handler
	BindingPowerTable[kind] = BP_DEFAULT
}

func createLookupTable() {
	// Literal
	nud(token.String, BP_PRIMARY, parsePrimaryExpression)
	nud(token.Number, BP_PRIMARY, parsePrimaryExpression)
	nud(token.Tag, BP_PRIMARY, parsePrimaryExpression)
	nud(token.Key, BP_PRIMARY, parsePrimaryExpression)
	nud(token.LeftParen, BP_PRIMARY, parseGroupedExpression)

	// Logical
	led(token.Or, BP_LOGICAL, parseBinaryExpression)
	led(token.And, BP_LOGICAL, parseBinaryExpression)

	// Relational
	led(token.Equal, BP_RELATIONAL, parseBinaryExpression)
	led(token.LessThan, BP_RELATIONAL, parseBinaryExpression)

	// Additive
	led(token.Plus, BP_ADDITIVE, parseBinaryExpression)
	led(token.Minus, BP_ADDITIVE, parseBinaryExpression)

	// Multiplicative
	led(token.Slash, BP_MULTIPLICATIVE, parseBinaryExpression)
	led(token.Star, BP_MULTIPLICATIVE, parseBinaryExpression)

	stmt(token.Tag, parseTagDecStatement)
}
