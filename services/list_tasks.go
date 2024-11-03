package services

import (
	"github.com/luke-goddard/taskninja/db"
)

func (handler *ServiceHandler) ListTasks() ([]db.Task, error) {
	return handler.store.ListTasks()
}
