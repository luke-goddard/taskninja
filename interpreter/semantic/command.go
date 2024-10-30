package semantic

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
)

func (a *Analyzer) VisitCommand(cmd *ast.Command) *Analyzer {
	switch cmd.Kind {
	case ast.CommandKindAdd:
		return a.VisitAddCommand(cmd)
	case ast.CommandKindList:
		return a.VisitListCommand(cmd)
	}
	return a.EmitError(fmt.Sprintf("Unknown command kind: %d", cmd.Kind), cmd)
}

func (a *Analyzer) VisitAddCommand(cmd *ast.Command) *Analyzer {
	if cmd.Param == nil {
		return a.EmitError("Add command requies a description", cmd)
	}

	if cmd.Param.Kind != ast.ParamTypeDescription {
		return a.EmitError("Add command requires a description", cmd.Param)
	}

	if len(cmd.Param.Value) == 0 {
		return a.EmitError("Description cannot be empty", cmd.Param)
	}

	return a
}

func (a *Analyzer) VisitListCommand(cmd *ast.Command) *Analyzer {
	if cmd.Param != nil {
		return nil
	}
	return a
}
