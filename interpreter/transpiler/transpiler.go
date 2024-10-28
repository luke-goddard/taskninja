package transpiler

import (
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/manager"
)

type SqlStatement string
type SqlArgs []interface{}

type Transpiler struct {
	manager *manager.ErrorManager
	join    *JoinTranspiler
}

func NewTranspiler(manager *manager.ErrorManager) *Transpiler {
	return &Transpiler{
		manager: manager,
	}
}

func (transpiler *Transpiler) Reset() *Transpiler {
	transpiler.manager.Reset()
	return transpiler
}

func (transpiler *Transpiler) Transpile(
	command *ast.Command,
) (SqlStatement, SqlArgs, []manager.ErrorTranspiler) {
	switch command.Kind {
	case ast.CommandKindList:
		var builder = sqlbuilder.
			Select("id").
			From("tasks")
		var whereClauses = command.EvalSelect(builder, nil)
		builder.Where(fmt.Sprint(whereClauses))
		var sql, args = builder.Build()
		return SqlStatement(sql), SqlArgs(args), transpiler.manager.Errors()

	}
	return "", nil, transpiler.manager.Errors()
}
