package ast

import (
	"github.com/huandu/go-sqlbuilder"
)

type ParamType int

const (
	ParamTypeTaskId      ParamType = iota // e.g 1"
	ParamTypeDescription                  // e.g "buy dog"
	ParamTypeDependency                   // e.g 1
)

// Param represents a parameter in the AST.
// Some command require parameters like `task 1 modify`
// Here the parameter is 1
type Param struct {
	Kind  ParamType   // e.g TaskId, Description
	Value interface{} // e.g 1, "buy dog"
	NodePosition
}

// ParamDependency represents a dependency parameter in the AST.
// Some command require dependencies like `task 1 depends 2`
type ParamDependency struct {
	TaskId      int64
	DependsOnId int64
}

func (p *Param) Type() NodeType {
	return NodeTypeParam
}

func (p *Param) Expression() {}

func (p *Param) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	return ""
}

func (p *Param) EvalInsert(transpiler *Transpiler) interface{} {
	return ""
}
