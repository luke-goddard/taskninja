package interpreter

import (
	"testing"

	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/stretchr/testify/assert"
)

func TestInterpreterGood(t *testing.T) {
	var tc = []struct {
		input        string
		expectedSql  string
		expectedArgs interface{}
	}{
		{
			input:        `add "do the dishes"`,
			expectedSql:  "INSERT INTO tasks (title) VALUES (?)",
			expectedArgs: ast.SqlArgs{"do the dishes"},
		},
		{
			input:        `add "cook" priority:High`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityHigh},
		},
		{
			input:        `add "cook" priority:Medium`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityMedium},
		},
		{
			input:        `add "cook" priority:Low`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityLow},
		},
		{
			input:        `add "cook" priority:None`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityNone},
		},
		{
			input:        `add "cook" priority:high`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityHigh},
		},
		{
			input:        `add "cook" priority:medium`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityMedium},
		},
		{
			input:        `add "cook" priority:low`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityLow},
		},
		{
			input:        `add "cook" priority:none`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityNone},
		},
		{
			input:        `add "cook" priority:h`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityHigh},
		},
		{
			input:        `add "cook" priority:m`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityMedium},
		},
		{
			input:        `add "cook" priority:l`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityLow},
		},
		{
			input:        `add "cook" priority:n`,
			expectedSql:  `INSERT INTO tasks (title, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", db.TaskPriorityNone},
		},
	}

	var interpreter = NewInterpreter()
	for _, test := range tc {
		t.Run(test.input, func(t *testing.T) {
			var sql, args, err = interpreter.Execute(test.input)
			assert.Nil(t, err)
			assert.Equal(t, test.expectedSql, string(sql))
			assert.Equal(t, test.expectedArgs, args)
		})
	}
}

func TestInterpreterBad(t *testing.T) {
	var tc = []struct {
		input       string
		expectedErr string
	}{
		{
			input:       `add "" project:Lol`,
			expectedErr: "(Fatal) Semantic: Description cannot be empty",
		},
		{
			input:       `add 1 project:Lol`,
			expectedErr: "(Fatal) Syntax: Expected token type String, got Number: Number(1)",
		},
	}

	var interpreter = NewInterpreter()
	for _, test := range tc {
		t.Run(test.input, func(t *testing.T) {
			var _, _, err = interpreter.Execute(test.input)
			assert.NotNil(t, err)
			assert.Equal(t, test.expectedErr, err.Error())
		})
	}
}
