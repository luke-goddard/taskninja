package ast

import "github.com/luke-goddard/taskninja/assert"

type Visitor interface {
	Visit(node Node) Visitor // Visit returns the visitor to use for the children of the node.
}

// Walk traverses an AST in depth-first order: It starts by calling v.Visit(node); node must not be nil.
func Walk(v Visitor, node Node) {
	v.Visit(node)
	switch n := node.(type) {

	case *Literal:
	case *Tag:

	case *Command:
		Walk(v, n.Param)
		for _, node := range n.Options {
			Walk(v, node)
		}

	case *BinaryExpression:
		Walk(v, n.Left)
		Walk(v, n.Right)

	case *ExpressionStatement:
		Walk(v, n.Expr)

	case *LogicalExpression:
		Walk(v, n.Left)
		Walk(v, n.Right)

	default:
		assert.Fail("unexpected node type %T", n)
	}

}

// WalkList traverses a list of nodes in depth-first order: It starts by calling v.Visit(node) for each node in the list.
func WalkList(v Visitor, nodes []Node) {
	for _, node := range nodes {
		Walk(v, node)
	}
}
