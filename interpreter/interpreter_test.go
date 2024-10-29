package interpreter

import (
	"testing"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	var tc = []struct {
		input        string
		expectedSql  string
		expectedArgs interface{}
	}{
		{
			input:        `add "do the dishes"`,
			expectedSql:  "INSERT INTO tasks (description) VALUES (?)",
			expectedArgs: ast.SqlArgs{"do the dishes"},
		},
		{
			input:        `add "cook" priority:High`,
			expectedSql:  `INSERT INTO tasks (description, priority) VALUES (?, ?)`,
			expectedArgs: ast.SqlArgs{"cook", "High"},
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
