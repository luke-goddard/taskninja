// ========================================================
// Examples:
// ========================================================
// 1.   task add "Buy milk" due:2016-01-02 priority:high
// 2.   task 1 modify due:2016-01-02
// 3.   task list (project:home and priority:high)
// ========================================================
// Grammar:
// ========================================================
// COMMAND -> (
//  COMMAND_ADD | // e.g add "Buy milk"
// )
// COMMAND_ADD -> (
//  add PARAM | // e.g add "Buy milk"
//  add PARAM EXPRESSION_STATEMENTS | // e.g add "Buy milk" due:2016-01-02 priority:high
// )
// EXPRESSION_STATEMENTS -> EXPRESSION_STATEMENT | EXPRESSION_STATEMENT EXPRESSION_STATEMENTS
// EXPRESSION_STATEMENT -> EXPRESSION | EXPRESSION EXPRESSION_STATEMENT
// EXPRESSION -> BINARY_EXPRESSION | LOGICAL_EXPRESSION | TAG | PAIR
// BINARY_EXPRESSION -> (EXPRESSION) BINARY_OPERATOR (EXPRESSION) | (EXPRESSION) BINARY_OPERATOR (EXPRESSION) BINARY_EXPRESSION
// LOGICAL_EXPRESSION -> (EXPRESSION) LOGICAL_OPERATOR (EXPRESSION) | (EXPRESSION) LOGICAL_OPERATOR (EXPRESSION) LOGICAL_EXPRESSION
// LOGICAL_OPERATOR -> and | or
// BINARY_OPERATOR -> + | - | * | / | %
// PARAM -> TASKID | STRING
// TAG -> +TAG | -TAG
// PAIR -> key:EXPRESSION | key:EXPRESSION_STATEMENTS
// TASKID -> number

package parser

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

type ParseError struct {
	message error
	token   lex.Token
}

type ParseErrorList []ParseError

func (e *ParseErrorList) add(message error, token lex.Token) {
	*e = append(*e, ParseError{message, token})
}

type Parser struct {
	tokens   []lex.Token
	position int
	errors   ParseErrorList
}

func NewParser(tokens []lex.Token) *Parser {
	createLookupTable()
	return &Parser{tokens: tokens, position: 0}
}

func (p *Parser) Parse() (*ast.Command, ParseErrorList) {
	if p.tokens == nil || len(p.tokens) == 0{
		p.errors.add(fmt.Errorf("no tokens to parse"), lex.Token{})
		return nil, p.errors
	}
	fmt.Printf("parsing\n")
	return parseCommand(p), p.errors
}

func (p *Parser) current() *lex.Token {
	return &p.tokens[p.position]
}

func (p *Parser) consume() lex.Token {
	p.position++
	return p.tokens[p.position-1]
}
