package services

import (
	"context"

	"github.com/luke-goddard/taskninja/db"
)

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

func (handler *ServiceHandler) SetPriority(id int64, priority db.TaskPriority) (bool, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.SetPriority(ctx, id, priority)
}
