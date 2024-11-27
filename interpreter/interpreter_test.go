package interpreter

import (
	"testing"

	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter/ast"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestServices(t *testing.T) {
	log.Logger = log.Output(zerolog.Nop())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Interpreter Suite")
}

var _ = Describe("Transpiler should transpile", func() {
	var interpreter = NewInterpreter()

	DescribeTable("good",
		func(input string, expectedSql string, expectedArgs interface{}) {
			sql, args, err := interpreter.Execute(input)
			Expect(err).To(BeNil())
			Expect(string(sql)).To(Equal(expectedSql))
			Expect(args).To(Equal(expectedArgs))
		},
		Entry(
			"add 'do the dishes'",
			`add "do the dishes"`,
			"INSERT INTO tasks (title) VALUES (?)",
			ast.SqlArgs{"do the dishes"},
		),
		Entry(
			"add 'cook' priority:High",
			`add "cook" priority:High`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityHigh},
		),
		Entry(
			"add 'cook' priority:Medium",
			`add "cook" priority:Medium`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityMedium},
		),
		Entry(
			"add 'cook' priority:Low",
			`add "cook" priority:Low`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityLow},
		),
		Entry(
			"add 'cook' priority:None",
			`add "cook" priority:None`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityNone},
		),
		Entry(
			"add 'cook' priority:high",
			`add "cook" priority:high`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityHigh},
		),
		Entry(
			"add 'cook' priority:medium",
			`add "cook" priority:medium`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityMedium},
		),
		Entry(
			"add 'cook' priority:low",
			`add "cook" priority:low`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityLow},
		),
		Entry(
			"add 'cook' priority:none",
			`add "cook" priority:none`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityNone},
		),
		Entry(
			"add 'cook' priority:h",
			`add "cook" priority:h`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityHigh},
		),
		Entry(
			"add 'cook' priority:m",
			`add "cook" priority:m`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityMedium},
		),
		Entry(
			"add 'cook' priority:l",
			`add "cook" priority:l`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityLow},
		),
		Entry(
			"add 'cook' priority:n",
			`add "cook" priority:n`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityNone},
		),
	)
})

var _ = Describe("Transpiler should fail", func() {
	var interpreter = NewInterpreter()

	DescribeTable("bad",
		func(input string, expectedErr string) {
			_, _, err := interpreter.Execute(input)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal(expectedErr))
		},
		Entry(
			`add "" project:Lol`,
			`add "" project:Lol`,
			`(Fatal) Semantic: Description cannot be empty`,
		),
		Entry(
			`add 1 project:Lol`,
			`add 1 project:Lol`,
			`(Fatal) Syntax: Expected token type String, got Number: Number(1)`,
		),
	)
})
