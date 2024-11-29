package services

import "github.com/luke-goddard/taskninja/db"

func (serv *ServiceHandler) GetDependenciesForServices(taskId int64) ([]db.TaskDependency, error) {
	return serv.Store.GetDependenciesForTask(taskId)
}
