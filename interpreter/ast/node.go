package ast

type Column int   // Column represents a column in the source code.
type Line int     // Line represents a line in the source code.
type NodeType int // NodeType represents the type of a node. e.g BinaryExpression, Literal, Command, etc.

const (
	NodeTypeBinaryExpression    NodeType = iota // NodeTypeBinaryExpression represents a binary expression. e.g 1 + 2
	NodeTypeCommand                             // NodeTypeCommand represents a command. e.g add "buy dog" priority:High
	NodeTypeExpressionStatement                 // NodeTypeExpressionStatement represents an expression statement. e.g 1 + 2
	NodeTypeLiteral                             // NodeTypeLiteral represents a literal. e.g "buy dog"
	NodeTypeLogicalExpression                   // NodeTypeLogicalExpression represents a logical expression. e.g 1 and 2
	NodeTypeOption                              // NodeTypeOption represents an option. e.g priority:High
	NodeTypePair                                // NodeTypePair represents a pair. e.g priority:High
	NodeTypeParam                               // NodeTypeParam represents a param. e.g "buy dog"
	NodeTypeProgram                             // NodeTypeProgram represents a program. e.g add "buy dog" priority:High
	NodeTypeTag                                 // NodeTypeTag represents a tag. e.g +tag
)

// Node represents a node in the AST.
type Node interface {
	Type() NodeType      // Type returns the type of the node.
	StartColumn() Column // StartColumn returns the starting column of the node.
	EndColumn() Column   // EndColumn returns the ending column of the node.
	StartLine() Line     // StartLine returns the starting line of the node.
	EndLine() Line       // EndLine returns the ending line of the node.
}
