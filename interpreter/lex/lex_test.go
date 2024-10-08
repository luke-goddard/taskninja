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
		{"hello hello", TokenString, "hello"},
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
		{"/2", TokenSlash, "/"},
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
			go lexer.Tokenize()
			var tokens = make([]Token, 0)
			for {
				var token = <-lexer.Items
				if token == nil {
					break
				}
				tokens = append(tokens, *token)
				t.Logf("Token: '%v'", token.String())

				if token.Type == TokenEOF || token.Type == TokenError {
					break
				}
			}
			assert.True(t, len(tokens) != 0, "Expected to recive a token")
			assert.Equal(t, c.tokenType, tokens[0].Type)
			assert.Equal(t, c.value, tokens[0].Value)
		})
	}
}

func FuzzLex(f *testing.F) {
	f.Add(`add "do" project:home +Home`)
	f.Fuzz(func(t *testing.T, input string) {
		var lexer = NewLexer(input)
		go lexer.Tokenize()
		for {
			var token = <-lexer.Items
			if token == nil {
				break
			}
			if token.Type == TokenError || token.Type == TokenEOF {
				break
			}
		}
	})
}
