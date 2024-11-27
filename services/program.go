package services

import (
	"context"

	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/interpreter/ast"
)

func (handler *ServiceHandler) RunProgram(program string) (*ast.Command, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	var sql, args, err = handler.Interprete.Execute(program)
	if err != nil {
		return nil, err
	}
	var lastCmd = handler.Interprete.GetLastCmd()
	assert.NotNil(lastCmd, "last command is nil")

	_, err = handler.Store.Con.ExecContext(ctx, string(sql), args...)

	return lastCmd, err
}
