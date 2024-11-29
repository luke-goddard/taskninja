package interpreter

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
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

var _ = Describe("Transpiler should transpile add commands", func() {
	var interpreter *Interpreter
	var store *db.Store
	var tx *sqlx.Tx
	var err error

	BeforeEach(func() {
		store = db.NewInMemoryStore()
		interpreter = NewInterpreter(store)
		tx, err = store.Con.BeginTxx(context.TODO(), &sql.TxOptions{ReadOnly: false})
		Expect(err).To(BeNil())
	})

	DescribeTable("good",
		func(input string, expectedSql string, expectedArgs interface{}) {
			sql, args, err := interpreter.Execute(input, tx)
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
		Entry(
			`add "cook" project:Home priority:L`,
			`add "cook" project:Home priority:L`,
			`INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			ast.SqlArgs{"cook", db.TaskPriorityLow},
		),
	)

	DescribeTable("bad",
		func(input string, expectedErr string) {
			_, _, err := interpreter.Execute(input, tx)
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

	Describe("When adding a task with a project", func() {
		BeforeEach(func() {
			_, _, err := interpreter.Execute(`add "cook" project:Home`, tx)
			Expect(err).To(BeNil())
		})
		It("should should add a new task", func() {
			var tasks, err = store.ListTasks(context.Background())
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal("cook"))
			Expect(tasks[0].ProjectNames.Value()).To(Equal("home"))
		})
	})

	Describe("When adding a task with a project and priority", func() {
		BeforeEach(func() {
			_, _, err := interpreter.Execute(`add "cook" project:Home priority:High`, tx)
			Expect(err).To(BeNil())
		})
		It("should should add a new task", func() {
			var tasks, err = store.ListTasks(context.Background())
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal("cook"))
			Expect(tasks[0].ProjectNames.Value()).To(Equal("home"))
			Expect(tasks[0].Priority).To(Equal(db.TaskPriorityHigh))
		})
	})

	Describe("When adding a task with multiple projects", func() {
		BeforeEach(func() {
			_, _, err := interpreter.Execute(`add "cook" project:Home project:Work`, tx)
			Expect(err).To(BeNil())
		})
		It("should should add a new task", func() {
			var tasks, err = store.ListTasks(context.Background())
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal("cook"))
			Expect(tasks[0].ProjectNames.Value()).To(Equal("home,work"))
		})
	})

	Describe("When adding a task with multiple projects of the same project", func() {
		BeforeEach(func() {
			_, _, err := interpreter.Execute(`add "cook" project:Home project:home`, tx)
			Expect(err).To(BeNil())
		})
		It("should should add a new task", func() {
			var tasks, err = store.ListTasks(context.Background())
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal("cook"))
			Expect(tasks[0].ProjectNames.Value()).To(Equal("home"))
		})
	})

	Describe("When adding a task with a project in uppercase", func() {
		BeforeEach(func() {
			_, _, err := interpreter.Execute(`add "cook" project:HOME`, tx)
			Expect(err).To(BeNil())
		})
		It("should should add a new task but with a lowercase project", func() {
			var tasks, err = store.ListTasks(context.Background())
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal("cook"))
			Expect(tasks[0].ProjectNames.Value()).To(Equal("home"))
		})
	})

	Describe("When adding a task with a dependency that exists", func() {
		BeforeEach(func() {
			var _, _, err = interpreter.Execute(`add "1"`, tx)
			Expect(err).To(BeNil())

			tx, err = store.Con.BeginTxx(context.TODO(), &sql.TxOptions{ReadOnly: false})
			Expect(err).To(BeNil())

			_, _, err = interpreter.Execute(`add "2" deps:1`, tx)
			Expect(err).To(BeNil())
		})

		It("should add a dependency", func() {
			var tasks, err = store.ListTasks(context.Background())
			Expect(err).To(BeNil())
			Expect(tasks).To(HaveLen(2))
			tasks, _ = store.ListTasks(context.Background())
			var t2 = store.FilterByTaskId(2, tasks)
			Expect(t2.Dependencies.Value()).To(Equal("1"))
		})
	})

	Describe("When adding a task with a dependency that is it's self", func() {
		It("should error", func() {
			var _, _, err = interpreter.Execute(`add "1" deps:4`, tx)
			Expect(err).NotTo(BeNil())
		})
	})
})

var _ = Describe("Should be able to add dependencies commands", func() {

	var interpreter *Interpreter
	var store *db.Store
	var tx *sqlx.Tx
	var tasks []db.TaskDetailed

	BeforeEach(func() {
		store = db.NewInMemoryStore()
		interpreter = NewInterpreter(store)

		// https://www.youtube.com/watch?v=o7NyNnwrm70
		interpreter.Execute(`add "get the money"`, store.MustCreateTxTodo())
		interpreter.Execute(`add "get the power"`, store.MustCreateTxTodo())
		tx = store.MustCreateTxTodo()
		tasks, _ = store.ListTasks(context.Background())
	})

	It("should add a dependency", func() {
		_, _, err := interpreter.Execute(`depends 1 on 2`, tx)
		Expect(err).To(BeNil())

		tasks, err = store.ListTasks(context.Background())
		Expect(err).To(BeNil())
		Expect(tasks).To(HaveLen(2))
		for _, task := range tasks {
			if task.ID == 1 {
				Expect(task.Dependencies.Value()).To(Equal("2"))
			}
		}
	})

	It("should add a dependency", func() {
		_, _, err := interpreter.Execute(`depends 1 2`, tx)
		Expect(err).To(BeNil())

		tasks, err = store.ListTasks(context.Background())
		Expect(err).To(BeNil())
		Expect(tasks).To(HaveLen(2))
		for _, task := range tasks {
			if task.ID == 1 {
				Expect(task.Dependencies.Value()).To(Equal("2"))
			}
		}
	})

	It("should not allow a cyclical dependency", func() {
		_, _, err := interpreter.Execute(`depends 1 1`, tx)
		Expect(err).NotTo(BeNil())
	})
	It("should not allow a negative taskId", func() {
		_, _, err := interpreter.Execute(`depends -1 1`, tx)
		Expect(err).NotTo(BeNil())
	})
	It("should not allow a negative dependencyId", func() {
		_, _, err := interpreter.Execute(`depends 1 -1`, tx)
		Expect(err).NotTo(BeNil())
	})
})
