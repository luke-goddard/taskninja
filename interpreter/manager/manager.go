package manager

import "github.com/luke-goddard/taskninja/interpreter/token"

// ErrorManager is used to store and process errors
type ErrorManager struct {
	errors    []ErrorTranspiler
	hasErrors bool
}

// NewErrorManager creates a new error manager
func NewErrorManager() *ErrorManager {
	return &ErrorManager{errors: make([]ErrorTranspiler, 0), hasErrors: false}
}

// GetErrors returns all the errors that have been stored
func (manager *ErrorManager) GetErrors() []ErrorTranspiler {
	return manager.errors
}

// HasErrors returns true if there are any errors stored
func (manager *ErrorManager) HasErrors() bool {
	return manager.hasErrors
}

// EmitLex emits a lexical error
func (manager *ErrorManager) EmitLex(message string, token *token.Token) {
	var err = NewLexError(message).
		SetToken(token).
		SetSeverityFatal()
	manager.emit(err)
}

// EmitParse emits a parsing error
func (manager *ErrorManager) EmitParse(message string, token *token.Token) {
	var err = NewParseError(message).
		SetToken(token).
		SetSeverityFatal()
	manager.emit(err)
}

// EmitSemantic emits a semantic error
func (manager *ErrorManager) EmitSemantic(message string, token *token.Token) {
	var err = NewSemanticError(message).
		SetToken(token).
		SetSeverityFatal()
	manager.emit(err)
}

// EmitTranspilation emits a transpilation error
func (manager *ErrorManager) EmitTranspilation(message string, token *token.Token) {
	var err = NewTranspilationError(message).
		SetToken(token).
		SetSeverityFatal()
	manager.emit(err)
}

func (manager *ErrorManager) emit(e *ErrorTranspiler) {
	manager.hasErrors = true
	manager.errors = append(manager.errors, *e)
}
