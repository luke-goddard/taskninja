package services

import "github.com/luke-goddard/taskninja/db"


func (handler *ServiceHandler) GetTaskTimes(id int64) ([]db.TaskTime, error) {
	return handler.Store.GetTaskTimes(id)
}
