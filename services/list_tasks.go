package services

import (
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/events"
)

func (handler *ServiceHandler) ListTasks(e *events.ListTasks) ([]db.Task, error) {
	return handler.store.ListTasks()
}
