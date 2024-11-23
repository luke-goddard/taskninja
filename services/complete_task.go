package services

func (handler *ServiceHandler) CompleteTaskById(taskId int64) (bool, error) {
	return handler.Store.CompleteTaskById(taskId)
}
