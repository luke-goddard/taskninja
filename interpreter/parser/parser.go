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
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

type ParseError struct {
	Message error
	Token   lex.Token
}

type ParseErrorList []ParseError

func (e *ParseErrorList) add(message error, token lex.Token) {
	if message == nil {
		panic("message cannot be nil")
	}
	fmt.Println(message)
	*e = append(*e, ParseError{message, token})
}

type Parser struct {
	tokens           []lex.Token
	position         int
	errors           ParseErrorList
	hasCheckedExists bool
	context          context.Context
	cancel           context.CancelFunc
}

func NewParser(tokens []lex.Token) *Parser {
	createLookupTable()
	var ctx, cancel = context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	var p = NewParserWithContext(tokens, ctx)
	p.cancel = cancel
	return p
}

func NewParserWithContext(tokens []lex.Token, ctx context.Context) *Parser {
	createLookupTable()
	return &Parser{
		tokens:           tokens,
		position:         0,
		hasCheckedExists: false,
		context:          ctx,
		errors:           ParseErrorList{},
	}
}

func (p *Parser) Parse() (*ast.Command, ParseErrorList) {
	// if p.cancel != nil {
	// 	defer p.cancel()
	// }
	if p.tokens == nil || len(p.tokens) == 0 {
		p.errors.add(fmt.Errorf("no tokens to parse"), lex.Token{})
		return nil, p.errors
	}
	return parseCommand(p), p.errors
}

func (p *Parser) hasTokens() bool {
	// select {
	// case <-p.context.Done():
	// 	p.errors.add(fmt.Errorf("Parsing took too long"), lex.Token{})
	// 	return false
	// default:
	// }
	if p.hasCheckedExists {
		return true
	}
	p.hasCheckedExists = true
	if p.position > len(p.tokens)-1 {
		return false
	}
	if p.current().Type == lex.TokenEOF {
		p.consume()
		return false
	}
	return true
}

func (p *Parser) hasNoTokens() bool {
	return !p.hasTokens()
}

func (p *Parser) current() *lex.Token {
	if !p.hasCheckedExists {
		panic("Must call hasTokens before calling current")
	}
	if p.position >= len(p.tokens) {
		var err = fmt.Errorf(
			"Position is greater than the number of tokens Position: %d total: %d\n",
			p.position,
			len(p.tokens),
		)
		panic(err)
	}
	if p.position < 0 {
		panic("Position is less than 0")
	}
	return &p.tokens[p.position]
}

func (p *Parser) consume() *lex.Token {
	if !p.hasCheckedExists {
		panic("Must call hasTokens before calling consume")
	}
	if p.position >= len(p.tokens) {
		debug.PrintStack()
		panic("Position is greater than the number of tokens")
	}
	if p.position < 0 {
		panic("Position is less than 0")
	}
	p.position++
	return &p.tokens[p.position-1]
}

func (p *Parser) expectCurrent(tokenType lex.TokenType) bool {
	if p.current().Type != tokenType {
		p.errors.add(
			fmt.Errorf(
				"Expected token type %s, got %s",
				tokenType.String(),
				p.current().Type.String(),
			),
			*p.current(),
		)
		return false
	}
	return true
}

func (p *Parser) expectOneOf(t ...lex.TokenType) bool {
	for _, tokenType := range t {
		if p.current().Type == tokenType {
			return true
		}
	}
	p.errors.add(
		fmt.Errorf(
			"Expected one of token types %v, got %s",
			t,
			p.current().Type.String(),
		),
		*p.current(),
	)
	return false
}
