package services

func (handler *ServiceHandler) StartTimeToggleById(id int64) error {
	return handler.Store.StartTrackingTaskTime(id)
}

func (handler *ServiceHandler) StopTimeToggleById(id int64) error {
	return handler.Store.StopTrackingTaskTime(id)
}
