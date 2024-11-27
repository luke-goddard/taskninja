package services

import (
	"context"

	"github.com/luke-goddard/taskninja/db"
)

func (handler *ServiceHandler) GetTaskTimes(id int64) ([]db.TaskTime, error) {

	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.GetTaskTimes(ctx, id)
}
