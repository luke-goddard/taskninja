package services

import (
	"github.com/luke-goddard/taskninja/db"
)

func (handler *ServiceHandler) ListTasks() ([]db.TaskDetailed, error) {
	return handler.Store.ListTasks()
}
