package ast

type LiteralType int // LiteralType is an enum for literal types.

const (
	LiteralTypeString LiteralType = iota // e.g "buy dog"
	LiteralTypeNumber LiteralType = iota // e.g 5
)

// Literal represents a literal value in the AST.
// Example (string): "buy dog"
// Example (number): 5
type Literal struct {
	Value       string
	LiteralType LiteralType
  NodePosition
}

func (l *Literal) Type() NodeType {
	return NodeTypeLiteral
}
