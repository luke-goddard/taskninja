package ast

// ExpressionStatement represents an expression statement.
// Example (Binary):    1 + 2
// Example (Logical):   1 and 2
// Example (Pair):      priority:<High
// Example (Pair):      priority:High
// Example (Literal):   "high"
// Example (Tag):       +HOME
type ExpressionStatement struct {
	Expression *Node // binary, logical, tag, pair, literal
}

func (e *ExpressionStatement) Type() NodeType {
	return NodeTypeExpressionStatement
}
