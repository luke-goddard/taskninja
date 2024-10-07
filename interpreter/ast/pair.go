package ast

// Pair represents a key-value pair in the AST.
// Example: priority:high
// Example: priority:<high
type Pair struct {
	Key   string              // e.g. priority
	Value ExpressionStatement // e.g. high
	NodePosition
}

func (p *Pair) Type() NodeType {
	return NodeTypePair
}

func (p *Pair) Expression() {}
