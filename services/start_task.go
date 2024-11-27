package services

import "context"

func (handler *ServiceHandler) StartTimeToggleById(id int64) error {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.StartTrackingTaskTime(ctx, id)
}

func (handler *ServiceHandler) StopTimeToggleById(id int64) error {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.StopTrackingTaskTime(ctx, id)
}
