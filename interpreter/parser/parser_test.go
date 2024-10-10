package parser

import (
	"testing"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
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
		{
			input: `add "do the dishes" project:=home`,
			expected: ast.Command{
				Kind:  ast.CommandKindAdd,
				Param: &ast.Param{Kind: ast.ParamTypeDescription, Value: "do the dishes"},
				Options: []ast.Statement{
					&ast.ExpressionStatement{
						Expr: &ast.BinaryExpression{
							Left:     &ast.Key{Key: "project"},
							Operator: ast.BinaryOperatorEq,
							Right:    &ast.Literal{Value: "home"},
						},
					},
				},
			},
		},
	}

	for _, test := range tc {
		t.Run(test.input, func(t *testing.T) {
			var lexer = lex.NewLexer(test.input)
			go lexer.Tokenize()

			var tokens = make([]lex.Token, 0)

			for {
				var token = <-lexer.Items
				if token == nil {
					break
				}
				tokens = append(tokens, *token)
				if token.Type == lex.TokenEOF || token.Type == lex.TokenError {
					break
				}
			}

			var parser = NewParser(tokens)
			var tree, errs = parser.Parse()
			if !assert.Empty(t, errs) {
				for _, err := range errs {
					t.Logf("Error: %v", err)
				}
				t.Skip()
			}
			litter.Dump(tree)
			assert.Equal(t, &test.expected, tree)
		})
	}
}

func FuzzParser(f *testing.F) {
	f.Fuzz(func(t *testing.T, input []byte) {
		t.Logf("Fuzzing: %s", input)
		var tokens = make([]lex.Token, 0)
		tokens = append(tokens, lex.Token{Type: lex.TokenCommand, Value: "add"})
		for i, b := range input {
			if b >= 18 {
				t.Skip()
			}
			var token = lex.Token{
				Type:  lex.TokenType(b),
				Value: string(input[i:]),
			}
			tokens = append(tokens, token)
		}
		for _, token := range tokens {
			t.Logf("Token: %v", token.Type.String())
		}
		tokens = append(tokens, lex.Token{Type: lex.TokenEOF, Value: ""})
		var parser = NewParser(tokens)
		parser.Parse()
	})
}
