package ast

type CommandKind int

const (
	CommandKindAdd CommandKind = iota // e.g add "buy dog"
)

// Command represents a command in the AST.
// Example: add "buy dog" priority:high
// -----------------------^^^^^^^^^^^^^ options
// -------------^^^^^^^^^ parameter
type Command struct {
	Kind   CommandKind           // Kind represents the type of command. e.g add
	Param  Param                 // Param represents a parameter in the command. e.g "buy dog"
	Options []ExpressionStatement // Option represents an option in the command. e.g priority:high
	NodePosition
}

func (c *Command) Type() NodeType {
	return NodeTypeCommand
}

func (c *Command) Statement() {}
