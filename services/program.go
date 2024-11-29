package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/rs/zerolog/log"
)

func (handler *ServiceHandler) RunProgram(program string) (*ast.Command, error) {
	log.Info().Msg("Running program")
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()

	var tx *sqlx.Tx
	var err error

	tx, err = handler.Store.Con.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Error starting transaction when running transpiler: %v", err)
	}
	_, _, err = handler.Interprete.Execute(program, tx)
	if err != nil {
		log.Error().Err(err).Msg("Error executing program")
		return nil, err
	}
	var lastCmd = handler.Interprete.GetLastCmd()
	assert.NotNil(lastCmd, "last command is nil")
	return lastCmd, err
}
