package services

import "context"

func (handler *ServiceHandler) DeleteTaskById(id int64) (bool, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.DeleteTaskById(ctx, id)
}
