package services

import (
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter"
)

type ServiceHandler struct {
	Interprete *interpreter.Interpreter
	Store      *db.Store
}

func NewServiceHandler(
	interpreter *interpreter.Interpreter,
	store *db.Store,
) *ServiceHandler {
	assert.NotNil(interpreter, "Interpreter is nil")
	assert.NotNil(store, "Store is nil")
	return &ServiceHandler{
		Interprete: interpreter,
		Store:      store,
	}
}
