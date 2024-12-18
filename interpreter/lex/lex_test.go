package lex

import (
	"testing"

	man "github.com/luke-goddard/taskninja/interpreter/manager"
	"github.com/luke-goddard/taskninja/interpreter/token"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLexer(t *testing.T) {
	log.Logger = log.Output(zerolog.Nop())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lexer Suite")
}

var _ = Describe("Lexer", func() {
	var manager *man.ErrorManager
	var lexer *Lexer

	BeforeEach(func() {
		manager = man.NewErrorManager()
		lexer = NewLexer(manager)
	})

	DescribeTable("Tokenize should return tokens for",
		func(input string, firstTyp token.TokenType, length int) {
			var tokens, _ = lexer.SetInput(input).Tokenize()
			Expect(tokens).To(HaveLen(length))
			Expect(tokens[0].Type).To(Equal(firstTyp))

		},
		Entry("String", "hello hello", token.String, 2),
		Entry("double String", "hello hello", token.String, 2),
		Entry("Command", "add hello", token.Command, 2),
		Entry("Command", "depends 1 on 2", token.Command, 4),
		Entry("Number", "1", token.Number, 1),
		Entry("Number", "1.1", token.Number, 1),
		Entry("Number", "-1.1", token.Number, 1),
		Entry("Plus", "+", token.Plus, 1),
		Entry("Minus", "-", token.Minus, 1),
		Entry("Slash", "/", token.Slash, 1),
		Entry("Slash with numb", "/1", token.Slash, 2),
		Entry("Star", "*", token.Star, 1),
		Entry("Double Quoted String", `"string"`, token.String, 1),
		Entry("Single Quoted String", `'string'`, token.String, 1),
		Entry("Double Quoted String with escape", `"string\""`, token.String, 1),
		Entry("Tag", `+Tag`, token.Tag, 1),
		Entry("Tag", `-Tag`, token.Tag, 1),
		Entry("Key", `project:home`, token.Key, 3),
		Entry("Colon", `:home`, token.Colon, 2),
		Entry("LeftParen", `(`, token.LeftParen, 1),
		Entry("RightParen", `)`, token.RightParen, 1),
		Entry("LessThan", `<`, token.LessThan, 1),
	)

})
