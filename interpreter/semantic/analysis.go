package semantic

import (
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/manager"
)

type Analyzer struct {
	errors *manager.ErrorManager
}

func NewAnalyzer(errors *manager.ErrorManager) *Analyzer {
	return &Analyzer{errors: errors}
}

func (a *Analyzer) Analyze(node ast.Node) *manager.ErrorTranspiler {
	a.Visit(node)
	if a.errors.HasErrors() {
		return &a.errors.Errors()[0]
	}
	return nil
}

func (a *Analyzer) Visit(node ast.Node) *Analyzer {
	if node == nil {
		return nil
	}

	switch node.(type) {
	case *ast.Command:
		return a.VisitCommand(node.(*ast.Command))
	}

	return a
}

func (a *Analyzer) EmitError(message string, node ast.Node) *Analyzer {
	a.errors.EmitSemantic(message, node)
	return nil // Return nil to stop the visitor
}

func (a *Analyzer) Reset() *Analyzer {
	a.errors.Reset()
	return a
}
