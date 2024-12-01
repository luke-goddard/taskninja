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
	case ast.CommandKindDepends:
		return a.VisitDependsCommand(cmd)
	case ast.CommandKindNext:
		return a.VisitNextCommand(cmd)
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

	var description = cmd.Param.Value.(string)
	if len(description) == 0 {
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

func (a *Analyzer) VisitDependsCommand(cmd *ast.Command) *Analyzer {
	var param = cmd.Param.Value.(ast.ParamDependency)
	if param.TaskId < 0 {
		return a.EmitError("Task ID cannot be negative", cmd.Param)
	}
	if param.DependsOnId < 0 {
		return a.EmitError("DependsOn ID cannot be negative", cmd.Param)
	}
	if param.TaskId == param.DependsOnId {
		return a.EmitError("Task ID and DependsOn ID cannot be the same", cmd.Param)
	}
	return a
}

func (a *Analyzer) VisitNextCommand(cmd *ast.Command) *Analyzer {
	var tid = cmd.Param.Value.(int64)
	if tid <= 0 {
		return a.EmitError("Task ID cannot be zero or negative", cmd.Param)
	}
	return a
}
