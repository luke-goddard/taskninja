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
		// 	{
		// 		input: `add "do the dishes"`,
		// 		expected: ast.Command{
		// 			Kind: ast.CommandKindAdd,
		// 			Param: ast.Param{
		// 				Kind:  ast.ParamTypeDescription,
		// 				Value: "do the dishes",
		// 			},
		// 		},
		// 	},
		// 	{
		// 		input: `add "do the dishes" +Home -Kitchen`,
		// 		expected: ast.Command{
		// 			Kind: ast.CommandKindAdd,
		// 			Param: ast.Param{
		// 				Kind:  ast.ParamTypeDescription,
		// 				Value: "do the dishes",
		// 			},
		// 			Options: []ast.ExpressionStatement{
		// 				ast.ExpressionStatement{
		// 					Expr: &ast.Tag{
		// 						Operator: ast.TagOperatorPlus,
		// 						Value:    "Home",
		// 					},
		// 				},
		// 				ast.ExpressionStatement{
		// 					Expr: &ast.Tag{
		// 						Operator: ast.TagOperatorMinus,
		// 						Value:    "Kitchen",
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	{
		// 		input: `add "do the dishes" project:home`,
		// 		expected: ast.Command{
		// 			Kind: ast.CommandKindAdd,
		// 			Param: ast.Param{
		// 				Kind:  ast.ParamTypeDescription,
		// 				Value: "do the dishes",
		// 			},
		// 			Options: []ast.ExpressionStatement{
		// 				ast.ExpressionStatement{
		// 					Expr: &ast.Key{
		// 						Key: "project",
		// 						Value: ast.ExpressionStatement{
		// 							Expr: &ast.Literal{
		// 								Kind:  ast.LiteralKindString,
		// 								Value: "home",
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	{
		// 		input: `add "do the dishes" priority:HIGH`,
		// 		expected: ast.Command{
		// 			Kind: ast.CommandKindAdd,
		// 			Param: ast.Param{
		// 				Kind:  ast.ParamTypeDescription,
		// 				Value: "do the dishes",
		// 			},
		// 			Options: []ast.ExpressionStatement{
		// 				ast.ExpressionStatement{
		// 					Expr: &ast.Key{
		// 						Key: "priority",
		// 						Value: ast.ExpressionStatement{
		// 							Expr: &ast.Literal{
		// 								Kind:  ast.LiteralKindString,
		// 								Value: "HIGH",
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
	}

	for _, test := range tc {
		t.Run(test.input, func(t *testing.T) {
			var lexer = lex.NewLexer(test.input)
			var tokens = lexer.Tokenize()

			for _, token := range tokens {
				t.Logf("Token: %v", token.String())
			}

			var parser = NewParser(tokens)
			var tree = parser.Parse()
			litter.Dump(tree)
			assert.Equal(t, &test.expected, tree)
		})
	}
}
