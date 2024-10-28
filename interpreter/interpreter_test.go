package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	var tc = []struct {
		input        string
		expectedSql  string
		expectedArgs interface{}
	}{
		{input: `list "do the dishes"`, expectedSql: "SELECT id FROM tasks WHERE description = ?", expectedArgs: "do the dishes"},
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
