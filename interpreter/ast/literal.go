package ast

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
