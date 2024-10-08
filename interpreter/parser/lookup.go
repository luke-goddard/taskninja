package parser

import (
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
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

type StatementLookup map[lex.TokenType]StatementHandler
type NudLookup map[lex.TokenType]NudHandler
type LedLookup map[lex.TokenType]LedHandler
type BindingPowerLookup map[lex.TokenType]BindingPower

var BindingPowerTable = BindingPowerLookup{}
var NudTable = NudLookup{}
var LedTable = LedLookup{}
var StatementTable = StatementLookup{}

func led(kind lex.TokenType, bp BindingPower, handler LedHandler) {
	LedTable[kind] = handler
	BindingPowerTable[kind] = bp
}

func nud(kind lex.TokenType, bp BindingPower, handler NudHandler) {
	NudTable[kind] = handler
	BindingPowerTable[kind] = bp
}

func stmt(kind lex.TokenType, handler StatementHandler) {
	StatementTable[kind] = handler
	BindingPowerTable[kind] = BP_DEFAULT
}

func createLookupTable() {
	// Literal
	nud(lex.TokenString, BP_PRIMARY, parsePrimaryExpression)
	nud(lex.TokenNumber, BP_PRIMARY, parsePrimaryExpression)
	nud(lex.TokenTag, BP_PRIMARY, parsePrimaryExpression)
	nud(lex.TokenKey, BP_PRIMARY, parsePrimaryExpression)
	nud(lex.TokenLeftParen, BP_PRIMARY, parseGroupedExpression)

	// Logical
	led(lex.TokenOr, BP_LOGICAL, parseBinaryExpression)
	led(lex.TokenAnd, BP_LOGICAL, parseBinaryExpression)

	// Relational
	led(lex.TokenEQ, BP_RELATIONAL, parseBinaryExpression)
	led(lex.TokenLT, BP_RELATIONAL, parseBinaryExpression)

	// Additive
	led(lex.TokenPlus, BP_ADDITIVE, parseBinaryExpression)
	led(lex.TokenMinus, BP_ADDITIVE, parseBinaryExpression)

	// Multiplicative
	led(lex.TokenSlash, BP_MULTIPLICATIVE, parseBinaryExpression)
	led(lex.TokenStar, BP_MULTIPLICATIVE, parseBinaryExpression)

	stmt(lex.TokenTag, parseTagDecStatement)
}
