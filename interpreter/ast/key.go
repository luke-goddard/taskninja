package ast

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
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
	default:
		transpiler.AddError(fmt.Errorf("Unknown key: %s", key.Key), key)
		return nil
	}
}
