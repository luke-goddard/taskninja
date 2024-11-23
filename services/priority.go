package services

func (handler *ServiceHandler) IncreasePriority(id int64) (bool, error) {
	return handler.Store.IncreasePriority(id)
}

func (handler *ServiceHandler) DecreasePriority(id int64) (bool, error) {
	return handler.Store.DecreasePriority(id)
}
