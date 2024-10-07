package lex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexStart(t *testing.T) {
	var tc = []struct {
		input     string
		tokenType TokenType
		value     string
	}{
		{"", TokenEOF, ""},
		{"hello", TokenWord, "hello"},
		{"hello hello", TokenWord, "hello"},
		{"add hello", TokenCommand, "add"},
		{"all", TokenCommand, "all"},
		{"delete", TokenCommand, "delete"},
		{"done", TokenCommand, "done"},
		{"list", TokenCommand, "list"},
		{"modify", TokenCommand, "modify"},
		{"ready", TokenCommand, "ready"},
		{"start", TokenCommand, "start"},
		{"stop", TokenCommand, "stop"},
		{"tags", TokenCommand, "tags"},
		{"1", TokenNumber, "1"},
		{"1.1", TokenNumber, "1.1"},
		{"-1.1", TokenNumber, "-1.1"},
		{"+", TokenPlus, "+"},
		{"-", TokenMinus, "-"},
		{"/", TokenSlash, "/"},
		{"*", TokenStar, "*"},
		{`"string"`, TokenString, "string"},
		{`'string'`, TokenString, "string"},
		{`"string\""`, TokenString, `string\"`},
		{`+Tag`, TokenTag, "+Tag"},
		{`-Tag`, TokenTag, "-Tag"},
		{`project:home`, TokenKey, "project"},
		{`:home`, TokenColon, ":"},
		{`(`, TokenLeftParen, "("},
		{`)`, TokenRightParen, ")"},
		{`<`, TokenLT, "<"},
	}

	for _, c := range tc {
		t.Run(c.input, func(t *testing.T) {
			var lexer = NewLexer(c.input)
			var tokens = lexer.Tokenize()
			for _, token := range tokens {
				t.Logf("Token: %v", token.String())
			}
			assert.True(t, len(tokens) != 0, "Expected to recive a token")
			assert.Equal(t, c.tokenType, tokens[0].Type)
			assert.Equal(t, c.value, tokens[0].Value)
		})
	}
}
