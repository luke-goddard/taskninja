package services

func (handler *ServiceHandler) DeleteTaskById(id int64) (bool, error) {
	return handler.Store.DeleteTaskById(id)
}
