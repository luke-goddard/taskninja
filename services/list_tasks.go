package services

import (
	"github.com/luke-goddard/taskninja/db"
)

func (handler *ServiceHandler) ListTasks() ([]db.TaskDetailed, error) {
	var tasks, err = handler.Store.ListTasks()
	if err != nil {
		return nil, err
	}
	handler.SortTasksByUrgency(tasks)
	return tasks, nil

}

func (handler *ServiceHandler) SortTasksByUrgency(tasks []db.TaskDetailed) {
	for i := 0; i < len(tasks); i++ {
		for j := 0; j < len(tasks)-1; j++ {
			if tasks[j].Urgency() < tasks[j+1].Urgency() {
				tasks[j], tasks[j+1] = tasks[j+1], tasks[j]
			}
		}
	}
}

func (handler *ServiceHandler) CountTasks() (int64, error) {
	return handler.Store.CountTasks()
}

func (handler *ServiceHandler) GetTaskById(taskId int64) (*db.Task, error) {
	return handler.Store.GetTaskById(taskId)
}
