package manager

import (
	"fmt"

	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/luke-goddard/taskninja/interpreter/token"
)

type ErrorTranspilerVariant string  // What part of the pipeline the error occurred in
type ErrorTranspilerSeverity string // How severe the error is

const (
	TranspilerErrorLex           ErrorTranspilerVariant = "Lexical"
	TranspilerErrorParse         ErrorTranspilerVariant = "Syntax"
	TranspilerErrorSemantic      ErrorTranspilerVariant = "Semantic"
	TranspilerErrorTranspilation ErrorTranspilerVariant = "Transpilation"
)

const (
	TranspilerErrorSeverityFatal   ErrorTranspilerSeverity = "Fatal"
	TranspilerErrorSeverityWarning ErrorTranspilerSeverity = "Warning"
)

const DEFAULT_ERROR_SEVERITY = TranspilerErrorSeverityFatal

type ErrorTranspiler struct {
	Variant  ErrorTranspilerVariant
	Severity ErrorTranspilerSeverity
	Message  string
	Token    *token.Token
	Node     ast.Node // Node is an interface so no pointer
}

// NewErrorTranspiler creates a new error transpiler
func NewErrorTranspiler(
	variant ErrorTranspilerVariant,
	message string,
) *ErrorTranspiler {
	return &ErrorTranspiler{
		Variant:  variant,
		Severity: DEFAULT_ERROR_SEVERITY,
		Message:  message,
	}
}

func (e *ErrorTranspiler) Error() string {
	var baseMessage = fmt.Sprintf("(%s) %s: %s", e.Severity, e.Variant, e.Message)
	if e.hasToken() {
		baseMessage = fmt.Sprintf("%s: %s", baseMessage, e.Token.String())
	}
	return baseMessage
}

func (e *ErrorTranspiler) SetToken(token *token.Token) *ErrorTranspiler {
	e.Token = token
	return e
}

func (e *ErrorTranspiler) SetNode(node ast.Node) *ErrorTranspiler {
	e.Node = node
	return e
}

func (e *ErrorTranspiler) SetSeverityFatal() *ErrorTranspiler {
	e.Severity = TranspilerErrorSeverityFatal
	return e
}

func (e *ErrorTranspiler) SetSeverityWarning() *ErrorTranspiler {
	e.Severity = TranspilerErrorSeverityWarning
	return e
}

func (e *ErrorTranspiler) hasToken() bool { return e.Token != nil }
func (e *ErrorTranspiler) hasNode() bool  { return e.Node != nil }

// Occurred during the Lexical Analysis phase
func NewLexError(message string) *ErrorTranspiler {
	return NewErrorTranspiler(TranspilerErrorLex, message)
}

// Occurred during the Syntax Analysis phase
func NewParseError(message string) *ErrorTranspiler {
	return NewErrorTranspiler(TranspilerErrorParse, message)
}

// Occurred during the Semantic Analysis phase
func NewSemanticError(message string) *ErrorTranspiler {
	return NewErrorTranspiler(TranspilerErrorSemantic, message)
}

// Occurred during the Transpilation phase
func NewTranspilationError(message string) *ErrorTranspiler {
	return NewErrorTranspiler(TranspilerErrorTranspilation, message)
}
