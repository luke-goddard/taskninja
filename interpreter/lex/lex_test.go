package lex

import (
	"testing"

	"github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/token"
	"github.com/stretchr/testify/assert"
)

func TestLexStart(t *testing.T) {

	var tc = []struct {
		input     string
		tokenType token.TokenType
		value     string
	}{
		{"", token.Eof, ""},
		{"hello hello", token.String, "hello"},
		{"add hello", token.Command, "add"},
		{"all", token.Command, "all"},
		{"delete", token.Command, "delete"},
		{"done", token.Command, "done"},
		{"list", token.Command, "list"},
		{"modify", token.Command, "modify"},
		{"ready", token.Command, "ready"},
		{"start", token.Command, "start"},
		{"stop", token.Command, "stop"},
		{"tags", token.Command, "tags"},
		{"1", token.Number, "1"},
		{"1.1", token.Number, "1.1"},
		{"-1.1", token.Number, "-1.1"},
		{"+", token.Plus, "+"},
		{"-", token.Minus, "-"},
		{"/", token.Slash, "/"},
		{"/2", token.Slash, "/"},
		{"*", token.Star, "*"},
		{`"string"`, token.String, "string"},
		{`'string'`, token.String, "string"},
		{`"string\""`, token.String, `string\"`},
		{`+Tag`, token.Tag, "+Tag"},
		{`-Tag`, token.Tag, "-Tag"},
		{`project:home`, token.Key, "project"},
		{`:home`, token.Colon, ":"},
		{`(`, token.LeftParen, "("},
		{`)`, token.RightParen, ")"},
		{`<`, token.LessThan, "<"},
	}

	for _, c := range tc {
		t.Run(c.input, func(t *testing.T) {
			var manager = manager.NewErrorManager()
			var tokens, errs = NewLexer(manager).SetInput(c.input).Tokenize()
			assert.Len(t, errs, 0, "Expected no errors")
			assert.True(t, len(tokens) != 0, "Expected to recive a token")
			assert.Equal(t, c.tokenType, tokens[0].Type, "Expected token type to be %s, got %s", c.tokenType.String(), tokens[0].Type.String())
			assert.Equal(t, c.value, tokens[0].Value, "Expected token value to be %s, got %s", c.value, tokens[0].Value)
		})
	}
}
