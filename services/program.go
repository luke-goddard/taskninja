package services

import (
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/interpreter/ast"
)

func (handler *ServiceHandler) RunProgram(e *events.RunProgram) (*ast.Command, error) {
	var sql, args, err = handler.Interprete.Execute(e.Program)
	if err != nil {
		return nil, err
	}
	var lastCmd = handler.Interprete.GetLastCmd()
	assert.NotNil(lastCmd, "last command is nil")

	_, err = handler.Store.Con.Exec(string(sql), args...)

	return lastCmd, err
}
