package ast

type Program struct {
	Commands []*Command
  NodePosition
}

func (p *Program) Type() NodeType {
	return NodeTypeProgram
}
