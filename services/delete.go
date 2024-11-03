package services

func (handler *ServiceHandler) DeleteTasks(id int) (bool, error) {
	return handler.store.DeleteTaskById(id)
}
