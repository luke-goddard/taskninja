package services

import (
	"github.com/luke-goddard/taskninja/db"
)

func (handler *ServiceHandler) StartTasksById(id int64) (*db.Task, error) {
	return handler.Store.StartTaskById(id)
}
