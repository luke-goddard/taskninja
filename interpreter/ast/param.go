package ast

type ParamType int

const (
	ParamTypeTaskId      ParamType = iota // e.g 1"
	ParamTypeDescription ParamType = iota // e.g "buy dog"
)

// Param represents a parameter in the AST.
// Some command require parameters like `task 1 modify`
// Here the parameter is 1
type Param struct {
	ParamType ParamType // e.g TaskId, Description
	Value     string    // e.g 1, "buy dog"
}

func (p *Param) Type() NodeType {
	return NodeTypeParam
}
