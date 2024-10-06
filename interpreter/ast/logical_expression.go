package ast

type LogicalOperator int // LogicalOperator is an enum for logical operators.

const (
	LogicalOperatorAnd LogicalOperator = iota // e.g 1 and 2
	LogicalOperatorOr  LogicalOperator = iota // e.g 1 or 2
)

// LogicalExpression represents a logical expression.
// EXAMPLE (and):   1 and 2
// EXAMPLE (or):    1 or 2
type LogicalExpression struct {
	Left     Node
	Operator LogicalOperator
	Right    Node
}

func (l *LogicalExpression) Type() NodeType {
  return NodeTypeLogicalExpression
}
