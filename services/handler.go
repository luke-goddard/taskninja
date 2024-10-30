package services

import (
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter"
)

type ServiceHandler struct {
	interpreter *interpreter.Interpreter
	store       *db.Store
}

func NewServiceHandler(
	interpreter *interpreter.Interpreter,
	store *db.Store,
) *ServiceHandler {
	return &ServiceHandler{
		interpreter: interpreter,
		store:       store,
	}
}
