package ast

//=============================================================================
// Binary Expression
//=============================================================================

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
	Left     Node
	Operator BinaryOperator
	Right    Node
	NodePosition
}

func (b *BinaryExpression) Expression() {}
func (b *BinaryExpression) Type() NodeType {
	return NodeTypeBinaryExpression
}

//=============================================================================
// Expression Statement
//=============================================================================

// ExpressionStatement represents an expression statement.
// Example (Binary):    1 + 2
// Example (Logical):   1 and 2
// Example (Pair):      priority:<High
// Example (Pair):      priority:High
// Example (Literal):   "high"
// Example (Tag):       +HOME
type ExpressionStatement struct {
	Expr Expression // binary, logical, tag, pair, literal
	NodePosition
}

func (e *ExpressionStatement) Type() NodeType {
	return NodeTypeExpressionStatement
}

func (c *ExpressionStatement) Statement()  {}
func (c *ExpressionStatement) Expression() {}

//=============================================================================
// Literal Expression
//=============================================================================

type LiteralKind int // LiteralType is an enum for literal types.

const (
	LiteralKindString LiteralKind = iota // e.g "buy dog"
	LiteralKindNumber LiteralKind = iota // e.g 5
)

// Literal represents a literal value in the AST.
// Example (string): "buy dog"
// Example (number): 5
type Literal struct {
	Value string
	Kind  LiteralKind
	NodePosition
}

func (l *Literal) Type() NodeType {
	return NodeTypeLiteral
}

func (l *Literal) Expression() {}

func (l *LiteralKind) String() string {
	switch *l {
	case LiteralKindString:
		return "String"
	case LiteralKindNumber:
		return "Number"
	}
	return "Unknown"
}


//=============================================================================
// Logical Expression
//=============================================================================

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
	NodePosition
}

func (l *LogicalExpression) Type() NodeType {
	return NodeTypeLogicalExpression
}

func (l *LogicalExpression) Expression() {}
