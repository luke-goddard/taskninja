package ast

// Key represents a key-value pair in the AST.
// Example: priority:high
// Example: priority:<high
type Key struct {
	Key string // e.g. priority
	NodePosition
}

func (p *Key) Type() NodeType {
	return NodeTypePair
}

func (p *Key) Statement()  {}
func (p *Key) Expression() {}

