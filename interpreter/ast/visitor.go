package ast

type Visitor interface {
	Visit(node Node) Visitor
}

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
		panic("unreachable")
	}

}

func WalkList(v Visitor, nodes []Node) {
	for _, node := range nodes {
		Walk(v, node)
	}
}
