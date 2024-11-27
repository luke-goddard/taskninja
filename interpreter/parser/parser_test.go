package parser

import (
	"testing"

	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

func FuzzParser(f *testing.F) {
	var errManager = manager.NewErrorManager()
	var parser = NewParser(errManager)
	f.Fuzz(func(t *testing.T, input []byte) {
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
		parser.Reset().Parse(tokens)
	})
}
