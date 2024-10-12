package transpiler

import (
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/manager"
)

type SqlStatement string
type SqlArgs []interface{}

type Transpiler struct {
	manager *manager.ErrorManager
	join    *JoinTranspiler
	builder *SqlBuilder
}

func NewTranspiler(manager *manager.ErrorManager) *Transpiler {
	return &Transpiler{
		manager: manager,
		builder: NewSqlBuilder(),
	}
}

func (transpiler *Transpiler) Transpile(comand *ast.Command) (SqlStatement, SqlArgs, []manager.ErrorTranspiler) {
	panic("!TODO: Implement Transpiler.Transpile()")
}
