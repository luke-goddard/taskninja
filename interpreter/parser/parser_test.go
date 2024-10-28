package parser

import (
	"testing"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/token"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	var tc = []struct {
		input    string
		expected ast.Command
	}{
		{
			input: `add "do the dishes"`,
			expected: ast.Command{
				Kind:  ast.CommandKindAdd,
				Param: &ast.Param{Kind: ast.ParamTypeDescription, Value: "do the dishes"},
			},
		},
		{
			input: `add "do the dishes" +Home -Kitchen`,
			expected: ast.Command{
				Kind:  ast.CommandKindAdd,
				Param: &ast.Param{Kind: ast.ParamTypeDescription, Value: "do the dishes"},
				Options: []ast.Statement{
					&ast.Tag{Operator: ast.TagOperatorPlus, Value: "Home"},
					&ast.Tag{Operator: ast.TagOperatorMinus, Value: "Kitchen"},
				},
			},
		},
	}

	for _, test := range tc {
		t.Run(test.input, func(t *testing.T) {
			var errManager = manager.NewErrorManager()
			var tokens []token.Token
			var errs []manager.ErrorTranspiler
			var ast *ast.Command

			tokens, errs = lex.NewLexer(errManager).SetInput(test.input).Tokenize()
			ast, errs = NewParser(errManager).Parse(tokens)

			if !assert.Empty(t, errs) {
				for _, err := range errs {
					t.Logf("Error: %v", err)
				}
				t.Skip()
			}
			litter.Dump(ast)
			assert.Equal(t, &test.expected, ast)
		})
	}
}

func FuzzParser(f *testing.F) {
	var errManager = manager.NewErrorManager()
	var parser = NewParser(errManager)
	f.Fuzz(func(t *testing.T, input []byte) {
		t.Logf("Fuzzing: %s", input)
		var tokens = make([]token.Token, 0)
		tokens = append(tokens, token.Token{Type: token.Command, Value: "add"})
		for i, b := range input {
			if b >= 18 {
				t.Skip()
			}
			var token = token.Token{
				Type:  token.TokenType(b),
				Value: string(input[i:]),
			}
			tokens = append(tokens, token)
		}
		for _, token := range tokens {
			t.Logf("Token: %v", token.Type.String())
		}
		tokens = append(tokens, token.Token{Type: token.Eof, Value: ""})
		parser.Reset().Parse(tokens)
	})
}
