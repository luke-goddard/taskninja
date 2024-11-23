package services

import (
	"github.com/luke-goddard/taskninja/db"
)

func (handler *ServiceHandler) StartTimeToggleById(id int64) (*db.Task, error) {
	return handler.Store.StartTimeToggleById(id)
}
