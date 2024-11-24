package services

func (handler *ServiceHandler) CompleteTaskById(taskId int64) (bool, error) {
	// TODO: Convert to transaction
	var _, err = handler.Store.StopTrackingTaskTime(taskId)
	if err != nil {
		return false, err
	}
	return handler.Store.CompleteTaskById(taskId)
}
