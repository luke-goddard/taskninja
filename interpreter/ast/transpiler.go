package ast

import (
	"github.com/huandu/go-sqlbuilder"
)

type SqlStatement string
type SqlArgs []interface{}

type TranspileError struct {
	Message error
	Node    Node
}

type Transpiler struct {
	errors   []TranspileError
	values   []interface{}
	cols     []string
	Selecter *sqlbuilder.SelectBuilder
	Inserter *sqlbuilder.InsertBuilder
}

func NewTranspiler() *Transpiler {
	return &Transpiler{
		errors: make([]TranspileError, 0),
		values: make([]interface{}, 0),
		cols:   make([]string, 0),
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
	return transpiler
}

func (transpiler *Transpiler) Transpile(
	command *Command,
) (SqlStatement, SqlArgs, []TranspileError) {
	switch command.Kind {
	case CommandKindList:
		return transpiler.transpileCommandList(command)
	case CommandKindAdd:
		return transpiler.transpileCommandAdd(command)
	}
	return "", nil, transpiler.errors
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
	var sql, args = transpiler.Inserter.Build()
	return SqlStatement(sql), SqlArgs(args), transpiler.errors
}
