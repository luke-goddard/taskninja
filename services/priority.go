package services

import "context"

func (handler *ServiceHandler) IncreasePriority(id int64) (bool, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.IncreasePriority(ctx, id)
}

func (handler *ServiceHandler) DecreasePriority(id int64) (bool, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.DecreasePriority(ctx, id)
}
