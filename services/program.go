package services

import (
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/interpreter/ast"
)

func (handler *ServiceHandler) RunProgram(program string) (*ast.Command, error) {
	var sql, args, err = handler.Interprete.Execute(program)
	if err != nil {
		return nil, err
	}
	var lastCmd = handler.Interprete.GetLastCmd()
	assert.NotNil(lastCmd, "last command is nil")

	_, err = handler.Store.Con.Exec(string(sql), args...)

	return lastCmd, err
}
