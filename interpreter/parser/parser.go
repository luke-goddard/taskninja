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
	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

type Parser struct {
	tokens           []token.Token
	position         int
	hasCheckedExists bool // Token, not column/table
	errors           *manager.ErrorManager
}

func NewParser(manager *manager.ErrorManager) *Parser {
	createLookupTable()
	return &Parser{
		position:         0,
		hasCheckedExists: false,
		errors:           manager,
	}
}

func (p *Parser) Parse(tokens []token.Token) (*ast.Command, []manager.ErrorTranspiler) {
	p.tokens = tokens
	if len(tokens) == 0 {
		p.errors.EmitParse("no tokens to parse", &token.Token{})
		return nil, p.errors.ParseErrors()
	}
	return parseCommand(p), p.errors.ParseErrors()
}

func (p *Parser) Reset() *Parser {
	p.position = 0
	p.hasCheckedExists = false
	p.errors.Reset()
	return p
}

func (p *Parser) hasTokens() bool {
	if p.hasCheckedExists {
		return true
	}
	p.hasCheckedExists = true
	if p.position > len(p.tokens)-1 {
		return false
	}
	if p.current().Type == token.Eof {
		p.consume()
		return false
	}
	return true
}

func (parser *Parser) hasNoTokens() bool {
	return !parser.hasTokens()
}

func (parser *Parser) current() *token.Token {
	if !parser.hasCheckedExists {
		panic("Must call hasTokens before calling current")
	}
	if parser.position >= len(parser.tokens) {
		var err = fmt.Errorf(
			"Position is greater than the number of tokens Position: %d total: %d\n",
			parser.position,
			len(parser.tokens),
		)
		panic(err)
	}
	if parser.position < 0 {
		panic("Position is less than 0")
	}
	return &parser.tokens[parser.position]
}

func (parser *Parser) consume() *token.Token {
	if !parser.hasCheckedExists {
		panic("Must call hasTokens before calling consume")
	}
	if parser.position >= len(parser.tokens) {
		panic("Position is greater than the number of tokens")
	}
	if parser.position < 0 {
		panic("Position is less than 0")
	}
	parser.position++
	return &parser.tokens[parser.position-1]
}

func (parser *Parser) expectCurrent(tokenType token.TokenType) bool {
	if parser.current().Type != tokenType {
		var message = fmt.Sprintf(
			"Expected token type %s, got %s",
			tokenType.String(),
			parser.current().Type.String(),
		)
		parser.errors.EmitParse(message, parser.current())
		return false
	}
	return true
}

func (parser *Parser) expectOneOf(t ...token.TokenType) bool {
	for _, tokenType := range t {
		if parser.current().Type == tokenType {
			return true
		}
	}
	var current = parser.current().Type.String()
	var message = fmt.Sprintf("Expected one of token types %v, got %s", t, current)
	parser.errors.EmitLex(message, parser.current())
	return false
}
