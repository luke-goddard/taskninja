package ast

// Command represents a command in the AST.
// Example: add "buy dog" priority:high
// -----------------------^^^^^^^^^^^^^ options
// -------------^^^^^^^^^ parameters
type Command struct {
	Param  []*Param               // Param represents a parameter in the command. e.g "buy dog"
	Option []*ExpressionStatement // Option represents an option in the command. e.g priority:high
}

func (c *Command) Type() NodeType {
	return NodeTypeCommand
}
