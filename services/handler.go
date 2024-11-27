package services

import (
	"time"

	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter"
)

const DefaultTimeout = time.Duration(1500 * time.Millisecond)

type ServiceHandler struct {
	Interprete *interpreter.Interpreter
	Store      *db.Store
	Timeout    time.Duration
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
		Timeout:    DefaultTimeout,
	}
}

func (handler *ServiceHandler) timeout() time.Time {
	return time.Now().Add(handler.Timeout)
}
