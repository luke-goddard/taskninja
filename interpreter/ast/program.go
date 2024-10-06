package ast

type Program struct {
	Commands []*Command
}

func (p *Program) Type() NodeType {
  return NodeTypeProgram
}
