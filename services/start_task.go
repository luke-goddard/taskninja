package services

func (handler *ServiceHandler) StartTimeToggleById(id int64) error {
	// TODO: Convert to transaction
	var err = handler.Store.SetTaskStateToStarted(id)
	if err != nil {
		return err
	}
	_, err = handler.Store.StartTrackingTaskTime(id)
	return err
}

func (handler *ServiceHandler) StopTimeToggleById(id int64) error {
	// TODO: Convert to transaction
	var err = handler.Store.SetTaskStateToIncomplete(id)
	if err != nil {
		return err
	}
	_, err = handler.Store.StopTrackingTaskTime(id)
	return err
}
