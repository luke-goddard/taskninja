package services

import "github.com/luke-goddard/taskninja/db"

func (handler *ServiceHandler) CreateTask(task *db.Task) (*db.Task, error){
	return handler.Store.CreateTask(task)
}
