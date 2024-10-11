package errorManger

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/lex"
)

type ErrorTranspilerVariant string  // What part of the pipeline the error occurred in
type ErrorTranspilerSeverity string // How severe the error is

const (
	TranspilerErrorLex           ErrorTranspilerVariant = "Lexical Analysis Error"
	TranspilerErrorParse         ErrorTranspilerVariant = "Syntax Analysis Error"
	TranspilerErrorSemantic      ErrorTranspilerVariant = "Semantic Analysis Error"
	TranspilerErrorTranspilation ErrorTranspilerVariant = "Transpilation Error"
)

const (
	TranspilerErrorSeverityFatal   ErrorTranspilerSeverity = "Fatal"
	TranspilerErrorSeverityWarning ErrorTranspilerSeverity = "Warning"
)

type ErrorTranspiler struct {
	Variant  ErrorTranspilerVariant
	Severity ErrorTranspilerSeverity
	Message  string
	Token    *lex.Token
	Node     *ast.Node
}

func NewErrorTranspiler(
	variant ErrorTranspilerVariant,
	severity ErrorTranspilerSeverity,
	message string,
) *ErrorTranspiler {
	return &ErrorTranspiler{
		Variant:  variant,
		Severity: severity,
		Message:  message,
	}
}

func (e *ErrorTranspiler) Error() string {
	var baseMessage = fmt.Sprintf("(%s) %s: %s", e.Severity, e.Variant, e.Message)
  if e.hasToken() {
    baseMessage = fmt.Sprintf("%s at %s", baseMessage, e.Token.Position.String())
  }
  return baseMessage
}

func (e *ErrorTranspiler) SetToken(token *lex.Token) *ErrorTranspiler {
	e.Token = token
	return e
}

func (e *ErrorTranspiler) SetNode(node *ast.Node) *ErrorTranspiler {
	e.Node = node
	return e
}

func (e *ErrorTranspiler) hasToken() bool { return e.Token != nil }
func (e *ErrorTranspiler) hasNode() bool  { return e.Node != nil }
