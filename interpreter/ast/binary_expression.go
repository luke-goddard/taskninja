package ast

type BinaryOperator int // BinaryOperator is an enum for binary operators.

const (
	BinaryOperatorAdd BinaryOperator = iota // e.g 1 + 2
	BinaryOperatorSub BinaryOperator = iota // e.g 1 - 2
	BinaryOperatorMul BinaryOperator = iota // e.g 1 * 2
	BinaryOperatorDiv BinaryOperator = iota // e.g 1 / 2
	BinaryOperatorMod BinaryOperator = iota // e.g 1 % 2
	BinaryOperatorEq  BinaryOperator = iota // e.g 1 == 2
	BinaryOperatorNe  BinaryOperator = iota // e.g 1 != 2
	BinaryOperatorLt  BinaryOperator = iota // e.g 1 < 2
	BinaryOperatorLe  BinaryOperator = iota // e.g 1 <= 2
	BinaryOperatorGt  BinaryOperator = iota // e.g 1 > 2
	BinaryOperatorGe  BinaryOperator = iota // e.g 1 >= 2
)

// BinaryExpression represents a binary expression.
// EXAMPLE: 1 + 2
// Left: 1
// Operator: +
// Right: 2
type BinaryExpression struct {
	Left     *Node
	Operator BinaryOperator
	Right    *Node
}

func (b *BinaryExpression) Type() NodeType {
  return NodeTypeBinaryExpression
}

