package services

func (handler *ServiceHandler) DeleteTasks(id int64) (bool, error) {
	return handler.store.DeleteTaskById(id)
}
