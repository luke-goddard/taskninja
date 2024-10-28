package ast

import "github.com/huandu/go-sqlbuilder"

type TagOperator int // TagOperator is an enum for tag operators.

const (
	TagOperatorPlus  TagOperator = iota // e.g. +HOME
	TagOperatorMinus TagOperator = iota // e.g. -HOME
)

// Tag represents a tag in the AST.
// Example: +HOME
// Example: -HOME
type Tag struct {
	Operator TagOperator // e.g. + or -
	Value    string      // e.g. HOME
	NodePosition
}

func (t *Tag) Type() NodeType {
	return NodeTypeTag
}

func (t *Tag) Statement() {}

func (t *Tag) EvalSelect(builder *sqlbuilder.SelectBuilder, addError AddError) interface{} {
	return ""
}
