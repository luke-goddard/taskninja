package services

import (
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/events"
)

func (handler *ServiceHandler) StartTasks(e *events.StartTask) ([]db.Task, error) {
	// return handler.Store.StartTasks(e.Id)
	return nil, nil
}
