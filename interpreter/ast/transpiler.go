package ast

import (
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/luke-goddard/taskninja/db"
	"github.com/rs/zerolog/log"
)

type SqlStatement string
type SqlArgs []interface{}

type TranspileError struct {
	Message error
	Node    Node
}

type TranspilerContext struct {
	isPriorityKey bool
}

type TranspileCallback func(tx *sqlx.Tx, taskId int64) error

type Transpiler struct {
	errors    []TranspileError
	values    []interface{}
	cols      []string
	Selecter  *sqlbuilder.SelectBuilder
	Inserter  *sqlbuilder.InsertBuilder
	ctx       *TranspilerContext
	tx        *sqlx.Tx
	store     *db.Store
	callbacks []TranspileCallback // When multiple transactions are needed
}

func NewTranspiler(store *db.Store) *Transpiler {
	return &Transpiler{
		errors:    make([]TranspileError, 0),
		values:    make([]interface{}, 0),
		cols:      make([]string, 0),
		ctx:       &TranspilerContext{},
		callbacks: make([]TranspileCallback, 0),
		store:     store,
	}
}

func (transpiler *Transpiler) AddValue(value interface{}) {
	transpiler.values = append(transpiler.values, value)
}

func (transpiler *Transpiler) AddCol(col string) {
	transpiler.cols = append(transpiler.cols, col)
}

func (transpiler *Transpiler) AddError(message error, node Node) {
	transpiler.errors = append(transpiler.errors, TranspileError{
		Message: message,
		Node:    node,
	})
}

func (transpiler *Transpiler) Reset() *Transpiler {
	transpiler.errors = make([]TranspileError, 0)
	transpiler.values = make([]interface{}, 0)
	transpiler.cols = make([]string, 0)
	transpiler.Selecter = nil
	transpiler.Inserter = nil
	transpiler.ctx = &TranspilerContext{}
	transpiler.tx = nil
	transpiler.callbacks = make([]TranspileCallback, 0)
	return transpiler
}

func (transpiler *Transpiler) addCallback(fn TranspileCallback) {
	transpiler.callbacks = append(transpiler.callbacks, fn)
}

func (transpiler *Transpiler) getContext() TranspilerContext {
	var ctx = *transpiler.ctx
	transpiler.ctx = &TranspilerContext{}
	return ctx
}

func (transpiler *Transpiler) setContext(ctx TranspilerContext) {
	transpiler.ctx = &ctx
}

func (transpiler *Transpiler) Transpile(
	command *Command,
	tx *sqlx.Tx,
) (SqlStatement, SqlArgs, []TranspileError) {
	transpiler.tx = tx
	switch command.Kind {
	case CommandKindAdd:
		return transpiler.transpileCommandAdd(command)
	// case CommandKindList:
	// 	return transpiler.transpileCommandList(command)
	case CommandKindDepends:
		return "", nil, transpiler.transpileCommandDepends(command)
	case CommandKindNext:
		return "", nil, transpiler.transpileCommandNext(command)
	default:
		transpiler.AddError(fmt.Errorf("Unknown command kind: %s", command.Kind.String()), command)
		return "", nil, transpiler.errors
	}
}

func (transpiler *Transpiler) transpileCommandList(command *Command) (SqlStatement, SqlArgs, []TranspileError) {
	var builder = sqlbuilder.
		Select("id").
		From("tasks")
	var whereClauses = command.EvalSelect(builder, nil)
	var whereBuilder = whereClauses.(*sqlbuilder.SelectBuilder)
	var sql, args = whereBuilder.Build()
	return SqlStatement(sql), SqlArgs(args), transpiler.errors
}

func (transpiler *Transpiler) transpileCommandAdd(command *Command) (SqlStatement, SqlArgs, []TranspileError) {
	transpiler.Inserter = sqlbuilder.InsertInto("tasks")
	command.EvalInsert(transpiler)
	transpiler.Inserter.Cols(transpiler.cols...)
	transpiler.Inserter.Values(transpiler.values...)
	if len(transpiler.errors) != 0 {
		return "", nil, transpiler.errors
	}
	var sql, args = transpiler.Inserter.Build()
	log.Info().Str("sql", sql).Interface("args", args).Msg("Transpiler produced")
	var res, err = transpiler.tx.Exec(sql, args...)
	var taskId int64
	if err != nil {
		transpiler.AddError(fmt.Errorf("Failed to insert task: %w", err), command)
		return "", nil, transpiler.errors
	}

	taskId, err = res.LastInsertId()
	for _, callback := range transpiler.callbacks {
		var err = callback(transpiler.tx, taskId)
		if err != nil {
			transpiler.AddError(fmt.Errorf("Failed to execute postprocessing callback: %w", err), command)
			return "", nil, transpiler.errors
		}
	}
	return SqlStatement(sql), SqlArgs(args), transpiler.errors
}

func (tran *Transpiler) transpileCommandDepends(command *Command) []TranspileError {
	var param = command.Param.Value.(ParamDependency)
	var err = tran.store.TaskDependsOnTx(tran.tx, param.TaskId, param.DependsOnId)
	if err != nil {
		tran.AddError(fmt.Errorf("Failed to insert task dependency: %w", err), command)
		return tran.errors
	}
	return tran.errors
}

func (tran *Transpiler) transpileCommandNext(command *Command) []TranspileError {
	var taskId = command.Param.Value.(int64)
	if taskId < 0 {
		tran.AddError(fmt.Errorf("TaskId must be greater than zero"), command)
		return tran.errors
	}
	var err = tran.store.TaskMarkAsNextTx(tran.tx, taskId)
	if err != nil {
		tran.AddError(fmt.Errorf("Failed to mark task as next: %w", err), command)
		return tran.errors
	}
	return tran.errors

}
