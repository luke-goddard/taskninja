package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// =============================================================================
// Key
// =============================================================================
type Key struct {
	Key  string
	Expr Expression
	NodePosition
}

func (key *Key) Expression() {}
func (key *Key) Type() NodeType {
	return NodeTypeBinaryExpression
}

func (key *Key) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	return key.Expr.EvalSelect(builder, addError)
}

func (key *Key) EvalInsert(transpiler *Transpiler) interface{} {
	var lowerK = strings.ToLower(key.Key)
	switch lowerK {
	case "priority", "p":
		transpiler.AddCol("priority")
		transpiler.setContext(TranspilerContext{isPriorityKey: true})
		return key.Expr.EvalInsert(transpiler)
	case "proj", "project":
		return key.handleProjectKey(transpiler)
	case "deps", "dends", "dependencies":
		return key.handleDependencies(transpiler)
	default:
		transpiler.AddError(fmt.Errorf("Unknown key: %s", key.Key), key)
		return nil
	}
}

func (key *Key) handleProjectKey(transpiler *Transpiler) interface{} {
	if key.Expr.Type() != NodeTypeLiteral {
		transpiler.AddError(fmt.Errorf("Expected literal value for project key"), key)
		return nil
	}
	var lit = key.Expr.(*Literal)
	if lit.Kind != LiteralKindString {
		transpiler.AddError(fmt.Errorf("Expected string value for project key"), key)
		return nil
	}
	var projectName = strings.ToLower(lit.Value)
	var projectId, err = transpiler.store.ProjectGetIDByNameOrCreateTx(transpiler.tx, projectName)
	if err != nil {
		err = fmt.Errorf("Failed to get or create project with name: %s -> %w", projectName, err)
		transpiler.AddError(err, key)
		return nil
	}

	transpiler.addCallback(func(tx *sqlx.Tx, taskId int64) error {
		var err = transpiler.store.ProjectLinkTaskTx(tx, projectId, taskId)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				log.Warn().Msg("Project already linked to task")
				return nil
			}
			return fmt.Errorf(
				"Failed to link project with ID: %d to task with ID: %d -> %w",
				projectId, taskId, err,
			)
		}
		return nil
	})

	return nil

}

func (key *Key) handleDependencies(trans *Transpiler) interface{} {
	if key.Expr.Type() != NodeTypeLiteral {
		trans.AddError(fmt.Errorf("Expected literal value for dependencies"), key)
		return nil
	}
	var lit = key.Expr.(*Literal)
	if lit.Kind != LiteralKindNumber {
		trans.AddError(fmt.Errorf("Expected the dependency value to be a TaskID"), key)
		return nil
	}

	var depOnTaskId = lit.Value
	var depOnTaskIdInt64, err = strconv.ParseInt(depOnTaskId, 10, 64)
	if err != nil {
		trans.AddError(fmt.Errorf("Failed to parse taskId as Int64"), key)
		return nil
	}
	var exists = trans.store.TaskIdExistsAndNotCompleted(trans.tx, depOnTaskIdInt64)
	if !exists {
		trans.AddError(fmt.Errorf("Task dependency does not exist"), key)
	}

	trans.addCallback(func(tx *sqlx.Tx, taskId int64) error {
		var err = trans.store.TaskDependsOnTx(tx, taskId, depOnTaskIdInt64)
		if err != nil {
			trans.AddError(fmt.Errorf("Failed to insert task dependency: %w", err), key)
			return err
		}
		return err
	})

	return nil
}
